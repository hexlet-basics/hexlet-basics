package handlers_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/blogpost"
	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// A state is not an update: the Publication State of a blog post moves only
// through publish, unpublish and archive (ADR-0020, state transitions addendum).

// The second post is a draft in the fixtures; each test puts it in the state it needs.
var statePostID = fixture.ID(fixture.BlogPostSecondPost)

type transitionCall struct {
	name string
	// run calls the operation on the post and returns the post it answered with.
	run func(h *testsupport.Harness) (*api.BlogPost, error)
	to  string
}

var transitionCalls = []transitionCall{
	{"publish", func(h *testsupport.Harness) (*api.BlogPost, error) {
		return h.Client.AdminPublishBlogPost(context.Background(), api.AdminPublishBlogPostParams{ID: int32(statePostID)})
	}, "published"},
	{"unpublish", func(h *testsupport.Harness) (*api.BlogPost, error) {
		return h.Client.AdminUnpublishBlogPost(context.Background(), api.AdminUnpublishBlogPostParams{ID: int32(statePostID)})
	}, "draft"},
	{"archive", func(h *testsupport.Harness) (*api.BlogPost, error) {
		return h.Client.AdminArchiveBlogPost(context.Background(), api.AdminArchiveBlogPostParams{ID: int32(statePostID)})
	}, "archived"},
}

func setState(h *testsupport.Harness, state string) {
	h.DB.BlogPost.UpdateOneID(statePostID).SetState(state).ExecX(context.Background())
}

func stateOf(t *testing.T, h *testsupport.Harness) string {
	t.Helper()
	row := h.DB.BlogPost.Query().Where(blogpost.IDEQ(statePostID)).OnlyX(t.Context())
	require.NotNil(t, row.State)
	return *row.State
}

// legal is the transition table as the issue states it; every other pair of
// (operation, current state) is refused, repeating the current state included.
var legal = map[string][]string{
	"publish":   {"draft", "archived"},
	"unpublish": {"published"},
	"archive":   {"draft", "published"},
}

func TestEveryLegalTransitionMovesThePostAndEveryOtherIsRefusedWith409(t *testing.T) {
	for _, call := range transitionCalls {
		for _, from := range []string{"draft", "published", "archived"} {
			t.Run(call.name+" from "+from, func(t *testing.T) {
				h := testsupport.NewHarness(t)
				setState(h, from)

				post, err := call.run(h)

				allowed := false
				for _, f := range legal[call.name] {
					allowed = allowed || f == from
				}
				if allowed {
					require.NoError(t, err)
					assert.Equal(t, call.to, string(post.State.Value))
					assert.Equal(t, call.to, stateOf(t, h))
					return
				}
				require.Error(t, err)
				assert.Equal(t, http.StatusConflict, h.LastStatus())
				assert.Equal(t, from, stateOf(t, h), "a refused move changes nothing")
				assert.NotEmpty(t, problemOf(t, err).Detail.Value, "the refusal explains itself")
			})
		}
	}
}

func TestARefusedTransitionExplainsItselfInTheCallersLanguage(t *testing.T) {
	for locale, want := range map[string]string{
		"ru": "«опубликовано»",
		"en": "published to published",
		"es": "«publicado»",
	} {
		t.Run(locale, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			testsupport.SpeakTo(h, locale)
			setState(h, "published")

			_, err := h.Client.AdminPublishBlogPost(t.Context(), api.AdminPublishBlogPostParams{ID: int32(statePostID)})

			require.Error(t, err)
			assert.Contains(t, problemOf(t, err).Detail.Value, want)
		})
	}
}

func TestAPostWithoutAStateCountsAsADraft(t *testing.T) {
	h := testsupport.NewHarness(t)
	h.DB.BlogPost.UpdateOneID(statePostID).ClearState().ExecX(t.Context())

	post, err := h.Client.AdminPublishBlogPost(t.Context(), api.AdminPublishBlogPostParams{ID: int32(statePostID)})

	require.NoError(t, err)
	assert.Equal(t, api.BlogPostStatePublished, post.State.Value)
}

func TestATransitionOfAMissingPostIs404(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminPublishBlogPost(t.Context(), api.AdminPublishBlogPostParams{ID: 999999})

	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestUpdateDoesNotChangeThePublicationState(t *testing.T) {
	h := testsupport.NewHarness(t)
	setState(h, "archived")

	_, err := h.Client.AdminUpdateBlogPost(t.Context(), blogInput("Edited", "edited"), api.AdminUpdateBlogPostParams{ID: int32(statePostID)})

	require.NoError(t, err)
	assert.Equal(t, "archived", stateOf(t, h))
}

func TestPublishingMakesAPostPublicAndUnpublishingHidesIt(t *testing.T) {
	h := testsupport.NewHarness(t)
	setState(h, "draft")
	slug := *h.DB.BlogPost.GetX(t.Context(), statePostID).Slug
	locale := *h.DB.BlogPost.GetX(t.Context(), statePostID).Locale
	visible := func() bool {
		return h.DB.BlogPost.Query().Where(blogpost.IDEQ(statePostID), blogpost.State("published"), blogpost.Locale(locale), blogpost.Slug(slug)).ExistX(t.Context())
	}
	require.False(t, visible())

	_, err := h.Client.AdminPublishBlogPost(t.Context(), api.AdminPublishBlogPostParams{ID: int32(statePostID)})
	require.NoError(t, err)
	assert.True(t, visible())

	_, err = h.Client.AdminUnpublishBlogPost(t.Context(), api.AdminUnpublishBlogPostParams{ID: int32(statePostID)})
	require.NoError(t, err)
	assert.False(t, visible())
}

// actAs gives the signed-in editor the fixture role with the label.
func actAs(h *testsupport.Harness, roleLabel fixture.Label) {
	h.DB.StaffMember.Update().SetRoleID(fixture.ID(roleLabel)).ExecX(context.Background())
}

// A role that may edit may not publish, and the reverse: the state is not an
// update.
func TestEditingAndPublishingAreGrantedSeparately(t *testing.T) {
	t.Run("update only", func(t *testing.T) {
		h := testsupport.NewHarnessAs(t, editorLabel)
		actAs(h, fixture.StaffRoleBlogUpdater)

		_, err := h.Client.AdminUpdateBlogPost(t.Context(), permissionsPostInput("hello-world"), api.AdminUpdateBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		require.NoError(t, err)

		_, err = h.Client.AdminPublishBlogPost(t.Context(), api.AdminPublishBlogPostParams{ID: int32(statePostID)})
		require.Error(t, err)
		assert.Equal(t, http.StatusForbidden, h.LastStatus())
	})
	t.Run("publish only", func(t *testing.T) {
		h := testsupport.NewHarnessAs(t, editorLabel)
		actAs(h, fixture.StaffRoleBlogPublisher)

		_, err := h.Client.AdminPublishBlogPost(t.Context(), api.AdminPublishBlogPostParams{ID: int32(statePostID)})
		require.NoError(t, err)

		_, err = h.Client.AdminUpdateBlogPost(t.Context(), permissionsPostInput("hello-world"), api.AdminUpdateBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		require.Error(t, err)
		assert.Equal(t, http.StatusForbidden, h.LastStatus())
	})
}

func TestEachTransitionLeavesAnAuditRecordWithTheOldAndNewState(t *testing.T) {
	for _, call := range transitionCalls {
		from := legal[call.name][0]
		t.Run(call.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			setState(h, from)

			_, err := call.run(h)
			require.NoError(t, err)

			operation := map[string]string{"publish": "AdminPublishBlogPost", "unpublish": "AdminUnpublishBlogPost", "archive": "AdminArchiveBlogPost"}[call.name]
			rows := records(t, h, operation)
			require.Len(t, rows, 1)
			assert.Equal(t, strconv.Itoa(statePostID), *rows[0].EntityID)
			state, ok := fieldOf(changesOf(t, rows[0]), "state")
			require.True(t, ok)
			assert.JSONEq(t, `"`+from+`"`, string(state.Old))
			assert.JSONEq(t, `"`+call.to+`"`, string(state.New))
		})
	}
}

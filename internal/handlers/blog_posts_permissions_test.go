package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// editorLabel is the fixtures' editor: he holds blog_posts.read and nothing else
// on blog posts, and works in ru and en.
var editorLabel = fixture.UserStaffEditor

// permissionsPostInput is a complete post body: every field is required by the
// contract, nullable or not.
func permissionsPostInput(slug string) *api.BlogPostInput {
	return &api.BlogPostInput{
		Name:              api.NewNilString("Notes"),
		Slug:              api.NewNilString(slug),
		Description:       api.NilString{Null: true},
		RichBody:          "<p>x</p>",
		CoverAttachmentId: api.NilInt32{Null: true},
	}
}

// blogPostOperations calls each blog post admin operation once and reports the
// HTTP status it got, so a test can compare a role with and without the grant.
var blogPostOperations = map[string]struct {
	permission string
	call       func(h *testsupport.Harness) error
}{
	"AdminListBlogPosts": {"blog_posts.read", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminListBlogPosts(context.Background(), api.AdminListBlogPostsParams{})
		return err
	}},
	"AdminGetBlogPost": {"blog_posts.read", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminGetBlogPost(context.Background(), api.AdminGetBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		return err
	}},
	"AdminCreateBlogPost": {"blog_posts.create", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminCreateBlogPost(context.Background(), permissionsPostInput("permissions-create"))
		return err
	}},
	"AdminUpdateBlogPost": {"blog_posts.update", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminUpdateBlogPost(context.Background(), permissionsPostInput("hello-world"), api.AdminUpdateBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		return err
	}},
	"AdminSetBlogPostRelatedCourses": {"blog_posts.update", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminSetBlogPostRelatedCourses(context.Background(),
			&api.BlogPostRelatedCoursesInput{CourseIds: []int32{}},
			api.AdminSetBlogPostRelatedCoursesParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		return err
	}},
	"AdminSuggestBlogPostRelatedCourses": {"blog_posts.update", func(h *testsupport.Harness) error {
		err := h.Client.AdminSuggestBlogPostRelatedCourses(context.Background(),
			api.AdminSuggestBlogPostRelatedCoursesParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		return err
	}},
	// the second post is a draft and the first is published, so each transition is legal for an
	// administrator.
	"AdminPublishBlogPost": {"blog_posts.publish", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminPublishBlogPost(context.Background(), api.AdminPublishBlogPostParams{ID: fixture.ID32(fixture.BlogPostSecondPost)})
		return err
	}},
	"AdminUnpublishBlogPost": {"blog_posts.unpublish", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminUnpublishBlogPost(context.Background(), api.AdminUnpublishBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		return err
	}},
	"AdminArchiveBlogPost": {"blog_posts.archive", func(h *testsupport.Harness) error {
		_, err := h.Client.AdminArchiveBlogPost(context.Background(), api.AdminArchiveBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
		return err
	}},
	"AdminDeleteBlogPost": {"blog_posts.delete", func(h *testsupport.Harness) error {
		return h.Client.AdminDeleteBlogPost(context.Background(), api.AdminDeleteBlogPostParams{ID: fixture.ID32(fixture.BlogPostToDelete)})
	}},
}

func TestARoleWithoutTheGrantIsRefusedEachBlogPostOperation(t *testing.T) {
	t.Parallel()
	for name, op := range blogPostOperations {
		if op.permission == "blog_posts.read" {
			continue // the editor holds it, see the next test
		}
		t.Run(name, func(t *testing.T) {
			h := testsupport.NewHarnessAs(t, editorLabel)

			require.Error(t, op.call(h))

			assert.Equal(t, http.StatusForbidden, h.LastStatus())
		})
	}
}

func TestARoleWithTheGrantMayUseTheOperation(t *testing.T) {
	t.Parallel()
	for name, op := range blogPostOperations {
		t.Run(name, func(t *testing.T) {
			h := testsupport.NewHarnessAs(t, editorLabel)
			// The editor already holds blog_posts.read; every other grant is new.
			h.Grant(t, int(roleEditorID), op.permission)

			err := op.call(h)

			require.NoError(t, err)
			assert.Less(t, h.LastStatus(), 300)
		})
	}
}

// A grant opens only the operations that declare it: holding the right to read
// is not the right to change.
func TestAGrantDoesNotOpenAnotherPermissionsOperations(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)

	require.NoError(t, blogPostOperations["AdminListBlogPosts"].call(h))
	require.Error(t, blogPostOperations["AdminUpdateBlogPost"].call(h))
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestAnAdministratorMayUseEveryBlogPostOperation(t *testing.T) {
	t.Parallel()
	for name, op := range blogPostOperations {
		t.Run(name, func(t *testing.T) {
			h := testsupport.NewHarness(t)

			require.NoError(t, op.call(h))
			assert.Less(t, h.LastStatus(), 300)
		})
	}
}

// The grant does not carry across locales: the staff record still decides
// where the role may act.
func TestAGrantDoesNotWidenTheStaffMembersLocales(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	testsupport.SpeakTo(h, "es")

	require.Error(t, blogPostOperations["AdminListBlogPosts"].call(h))
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

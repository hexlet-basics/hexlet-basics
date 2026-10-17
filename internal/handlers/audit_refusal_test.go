package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/blogpost"
	"hexletbasics/internal/api"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// A refused call is recorded: the editor lacks blog_posts.create, the security
// handler turns the call away before any handler or transaction exists, and the
// refusal still leaves an Audit Record of its own.
func TestARefusedChangeIsRecordedAsForbidden(t *testing.T) {
	h := testsupport.NewHarnessAs(t, editorLabel)
	testsupport.ArriveFrom(h, "198.51.100.9")
	testsupport.Identify(h, "req-refused-1", "Agent/2.0")

	_, err := h.Client.AdminCreateBlogPost(context.Background(), permissionsPostInput("refused-post"))
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())

	rows := records(t, h, "AdminCreateBlogPost")
	require.Len(t, rows, 1)
	row := rows[0]
	assert.Equal(t, "forbidden", row.Outcome)
	assert.Equal(t, editorID, *row.ActorUserID)
	assert.Equal(t, "web", row.Channel)
	assert.Equal(t, "req-refused-1", *row.RequestID)
	assert.Equal(t, "198.51.100.9", *row.Address)
	assert.Equal(t, "Agent/2.0", *row.UserAgent)

	exists, err := h.DB.BlogPost.Query().Where(blogpost.Slug("refused-post")).Exist(t.Context())
	require.NoError(t, err)
	assert.False(t, exists)
}

// What a refusal holds is who, where and which operation, never a value: the
// caller was refused because they may not read or change it, and the body they
// sent was never decoded into anything the record could keep.
func TestARefusalRecordHoldsNoFieldValuesNorEntity(t *testing.T) {
	h := testsupport.NewHarnessAs(t, editorLabel)

	input := permissionsPostInput("refused-secret")
	input.RichBody = "<p>a value the record must not keep</p>"
	_, err := h.Client.AdminUpdateBlogPost(context.Background(), input, api.AdminUpdateBlogPostParams{ID: fixtures.ID32(fixtures.BlogPostHelloWorld)})
	require.Error(t, err)

	rows := records(t, h, "AdminUpdateBlogPost")
	require.Len(t, rows, 1)
	assert.Equal(t, "forbidden", rows[0].Outcome)
	assert.Empty(t, rows[0].Changes)
	assert.Nil(t, rows[0].EntityType)
	assert.Nil(t, rows[0].EntityID)
}

// A read the role may not do is a refusal too, not only a change.
func TestARefusedReadIsRecordedAsForbidden(t *testing.T) {
	h := testsupport.NewHarnessAs(t, editorLabel)

	_, err := h.Client.AdminListUsers(context.Background(), api.AdminListUsersParams{})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())

	rows := records(t, h, "AdminListUsers")
	require.Len(t, rows, 1)
	assert.Equal(t, "forbidden", rows[0].Outcome)
}

// An administrator is not refused, and is not recorded as if they had been.
func TestAnAllowedCallLeavesNoForbiddenRecord(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminListBlogPosts(context.Background(), api.AdminListBlogPostsParams{})
	require.NoError(t, err)

	total, err := h.DB.AuditRecord.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, total)
}

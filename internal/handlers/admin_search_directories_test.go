package handlers_test

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// The `q` filter on the lists a use-case module owns (see admin_search_test.go
// for the generated ones).

func TestAdminListBlogPostsSearchMatchesTheNameAndTheSlug(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	slugs := func(page *api.BlogPostPage) []string {
		return lo.Map(page.Items, func(p api.BlogPost, _ int) string { return p.Slug.Value })
	}

	byName, err := h.Client.AdminListBlogPosts(ctx, api.AdminListBlogPostsParams{Q: api.NewOptString("HELLO WORLD")})
	require.NoError(t, err)
	bySlug, err := h.Client.AdminListBlogPosts(ctx, api.AdminListBlogPostsParams{Q: api.NewOptString("archived-")})
	require.NoError(t, err)

	assert.Equal(t, []string{"hello-world"}, slugs(byName))
	assert.EqualValues(t, 1, byName.Total)
	assert.Equal(t, []string{"archived-notes"}, slugs(bySlug))
}

func TestAdminListBlogPostsSearchKeepsTheNewestFirstOrder(t *testing.T) {
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListBlogPosts(context.Background(), api.AdminListBlogPostsParams{Q: api.NewOptString("post")})

	require.NoError(t, err)
	require.Greater(t, len(page.Items), 1)
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminListCourseLessonsSearchMatchesTheNameTheDescriptionAndTheSlug(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	slugs := func(page *api.CourseLessonListItemPage) []string {
		return lo.Map(page.Items, func(l api.CourseLessonListItem, _ int) string { return l.Slug })
	}

	byName, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{Q: api.NewOptString("VARIABLES")})
	require.NoError(t, err)
	byDescription, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{Q: api.NewOptString("STORING values")})
	require.NoError(t, err)
	bySlug, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{Q: api.NewOptString("hello-world")})
	require.NoError(t, err)

	assert.Equal(t, []string{"variables"}, slugs(byName))
	assert.Equal(t, []string{"variables"}, slugs(byDescription))
	assert.Contains(t, slugs(bySlug), "hello-world")
	assert.EqualValues(t, len(bySlug.Items), bySlug.Total)
}

func TestAdminListCourseLessonsSearchStaysInsideTheListedLessons(t *testing.T) {
	h := testsupport.NewHarness(t)

	// "Old Hello" is a translation of a previous course version: the list does
	// not show it, and neither does a search for its name.
	page, err := h.Client.AdminListCourseLessons(context.Background(), api.AdminListCourseLessonsParams{Q: api.NewOptString("Old Hello")})

	require.NoError(t, err)
	assert.Empty(t, page.Items)
	assert.Zero(t, page.Total)
}

func TestAdminListLeadsSearchMatchesEveryContactField(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	emails := func(q string) []string {
		page, err := h.Client.AdminListLeads(ctx, api.AdminListLeadsParams{Q: api.NewOptString(q)})
		require.NoError(t, err)
		return lo.Map(page.Items, func(l api.Lead, _ int) string { return l.Email.Value })
	}

	assert.Equal(t, []string{"anna@example.com"}, emails("ANNA@"))
	assert.Equal(t, []string{"anna@example.com"}, emails("@anna"), "telegram")
	assert.Equal(t, []string{"anna@example.com"}, emails("79990000001"), "phone")
	assert.Equal(t, []string{"boris@example.com"}, emails("79990000002"), "whatsapp")
}

func TestAdminListMessagesSearchMatchesTheText(t *testing.T) {
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListMessages(context.Background(), api.AdminListMessagesParams{Q: api.NewOptString("CONSOLE.LOG")})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "Use console.log to print a string.", page.Items[0].Content)
	assert.EqualValues(t, 1, page.Total)
}

func TestAdminListUsersSearchMatchesTheNameAndTheEmail(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	byEmail, err := h.Client.AdminListUsers(ctx, api.AdminListUsersParams{Q: api.NewOptString("SECOND-ADMIN@")})
	require.NoError(t, err)
	byName, err := h.Client.AdminListUsers(ctx, api.AdminListUsersParams{Q: api.NewOptString("passwordless")})
	require.NoError(t, err)
	management, err := h.Client.AdminListManagementUsers(ctx, api.AdminListManagementUsersParams{Q: api.NewOptString("SECOND-ADMIN@")})
	require.NoError(t, err)

	require.Len(t, byEmail.Items, 1)
	assert.Equal(t, "second-admin@example.com", byEmail.Items[0].Email.Value)
	require.Len(t, byName.Items, 1)
	assert.Equal(t, "removed-passwordless@example.com", byName.Items[0].Email.Value)
	require.Len(t, management.Items, 1, "the management list over users shares the filter")
}

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

// The `q` filter of the admin lists: a case-insensitive substring over the
// resource's human-readable fields, in the same query as the window, so the
// total and the page both follow it. The rows come from the committed fixtures.

func TestAdminListBannersSearchIsCaseInsensitiveAndKeepsTheTotal(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListBanners(context.Background(), api.AdminListBannersParams{Q: api.NewOptString("BLACK friday")})

	require.NoError(t, err)
	assert.Equal(t, []string{"Black Friday sale"}, lo.Map(page.Items, func(b api.Banner, _ int) string { return b.Body }))
	assert.EqualValues(t, 1, page.Total)
}

func TestAdminListBannersSearchComposesWithThePageWindow(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	// "fr" is in two English banners ("Free", "Friday"): the match is wider than
	// one row, the window cuts it, the total still counts it all.
	page, err := h.Client.AdminListBanners(context.Background(), api.AdminListBannersParams{
		Q:       api.NewOptString("FR"),
		PerPage: api.NewOptInt32(1),
	})

	require.NoError(t, err)
	assert.Len(t, page.Items, 1)
	assert.EqualValues(t, 2, page.Total)
}

func TestAdminListBannersBlankSearchListsEverything(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	all, err := h.Client.AdminListBanners(ctx, api.AdminListBannersParams{})
	require.NoError(t, err)
	blank, err := h.Client.AdminListBanners(ctx, api.AdminListBannersParams{Q: api.NewOptString("  ")})

	require.NoError(t, err)
	assert.Equal(t, all.Total, blank.Total)
}

func TestAdminListBannersSearchTreatsLikeWildcardsAsText(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListBanners(context.Background(), api.AdminListBannersParams{Q: api.NewOptString("%")})

	require.NoError(t, err)
	assert.EqualValues(t, 1, page.Total, "only the banner that says 50% matches a percent sign")
}

func TestAdminListCoursesSearchMatchesTheSlug(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListCourses(context.Background(), api.AdminListCoursesParams{Q: api.NewOptString("PYTHON")})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "python", page.Items[0].Slug)
}

func TestAdminListCourseLandingPagesSearchMatchesTheSlug(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListCourseLandingPages(context.Background(), api.AdminListCourseLandingPagesParams{Q: api.NewOptString("RUBY-")})

	require.NoError(t, err)

	slugs := lo.Map(page.Items, func(p api.CourseLandingPage, _ int) string { return p.Slug })
	assert.ElementsMatch(t, []string{"ruby-ru", "ruby-en"}, slugs)
}

func TestAdminListReviewsSearchMatchesTheReviewerAndTheText(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	byName, err := h.Client.AdminListReviews(ctx, api.AdminListReviewsParams{Q: api.NewOptString("petrov")})
	require.NoError(t, err)
	byText, err := h.Client.AdminListReviews(ctx, api.AdminListReviewsParams{Q: api.NewOptString("QUICK practice")})
	require.NoError(t, err)

	require.Len(t, byName.Items, 1)
	assert.Equal(t, "Petrov", byName.Items[0].LastName.Value)
	require.Len(t, byText.Items, 1)
	assert.Equal(t, "Jane", byText.Items[0].FirstName.Value)
}

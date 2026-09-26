package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// getSitemap reads the sitemap rows and fails the test unless they come back.
func getSitemap(t *testing.T, h *testsupport.Harness) *api.Sitemap {
	t.Helper()
	sitemap, err := h.Client.GetSitemap(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	return sitemap
}

func sitemapLandingSlugs(sitemap *api.Sitemap) []string {
	slugs := make([]string, len(sitemap.LandingPages))
	for i, lp := range sitemap.LandingPages {
		slugs[i] = lp.Slug
	}
	return slugs
}

// Only published, listed pages whose course has lessons in the page's locale
// are listed, ru first, then en, by id. In the fixtures only the javascript
// course has lessons; its archived and draft pages stay out.
func TestGetSitemapListsLandingPagesWithLessons(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	sitemap := getSitemap(t, h)

	assert.Equal(t, []string{"javascript-ru", "frontend-ru", "frontend-en"}, sitemapLandingSlugs(sitemap))
	assert.Equal(t, api.LocaleRu, sitemap.LandingPages[0].Locale)
	assert.Equal(t, api.LocaleEn, sitemap.LandingPages[2].Locale)
}

// The rows do not depend on the request locale: the page lists every locale.
func TestGetSitemapIgnoresTheRequestLocale(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	sitemap := getSitemap(t, h)

	assert.Equal(t, []string{"javascript-ru", "frontend-ru", "frontend-en"}, sitemapLandingSlugs(sitemap))
}

// Unlisted and es pages drop off; an unfinished course does not, as the
// sitemap has no completed-course filter.
func TestGetSitemapLeavesOutUnlistedAndSpanishPages(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	ctx := t.Context()

	require.NoError(t, h.DB.LandingPage.Update().
		Where(landingpage.SlugEQ("frontend-ru")).
		SetListed(false).
		Exec(ctx))
	require.NoError(t, h.DB.LandingPage.Update().
		Where(landingpage.SlugEQ("frontend-en")).
		SetLocale("es").
		Exec(ctx))
	require.NoError(t, h.DB.Course.Update().
		Where(course.SlugEQ("javascript")).
		SetProgress("in_development").
		Exec(ctx))

	sitemap := getSitemap(t, h)

	assert.Equal(t, []string{"javascript-ru"}, sitemapLandingSlugs(sitemap))
}

// The lessons have to be in the page's own locale and in the course's current
// version: en lessons left only in an older version do not count.
func TestGetSitemapNeedsLessonsInTheCurrentVersionAndLocale(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	ctx := t.Context()

	javascript, err := h.DB.Course.Query().Where(course.SlugEQ("javascript")).Only(ctx)
	require.NoError(t, err)
	require.NotNil(t, javascript.CurrentVersionID)

	// Moved to es rather than deleted: lesson reviews hold them by foreign key.
	require.NoError(t, h.DB.CourseLessonTranslation.Update().
		Where(
			courselessontranslation.LocaleEQ("en"),
			courselessontranslation.CourseVersionID(*javascript.CurrentVersionID),
		).
		SetLocale("es").
		Exec(ctx))
	olderEn, err := h.DB.CourseLessonTranslation.Query().
		Where(courselessontranslation.LocaleEQ("en")).
		Exist(ctx)
	require.NoError(t, err)
	require.True(t, olderEn, "fixtures keep en lessons in an older version")

	sitemap := getSitemap(t, h)

	assert.Equal(t, []string{"javascript-ru", "frontend-ru"}, sitemapLandingSlugs(sitemap))
}

// Published posts of every locale, newest first; drafts and archived posts
// stay out.
func TestGetSitemapListsPublishedPostsNewestFirst(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	sitemap := getSitemap(t, h)

	slugs := make([]string, len(sitemap.BlogPosts))
	for i, post := range sitemap.BlogPosts {
		slugs[i] = post.Slug
	}
	assert.Equal(t, []string{"english-post", "learning-ruby", "hello-world", "oldest-post"}, slugs)
}

// Every category of every locale is listed.
func TestGetSitemapListsAllCategories(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	sitemap := getSitemap(t, h)

	slugs := make([]string, len(sitemap.Categories))
	for i, category := range sitemap.Categories {
		slugs[i] = category.Slug.Value
	}
	assert.ElementsMatch(t, []string{
		"layouting-ru", "layouting-en", "programming-en",
		"frontend-en", "frontend-ru", "programming-ru",
	}, slugs)
}

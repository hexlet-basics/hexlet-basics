package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// getCategoryView reads a category page in the given locale and fails the test
// unless it is found.
func getCategoryView(t *testing.T, h *testsupport.Harness, slug string) *api.CourseCategoryView {
	t.Helper()
	res, err := h.Client.GetPublicCourseCategory(t.Context(), api.GetPublicCourseCategoryParams{Slug: slug})
	require.NoError(t, err)
	view, ok := res.(*api.CourseCategoryView)
	require.True(t, ok, "got %T", res)
	return view
}

func categoryLandingSlugs(view *api.CourseCategoryView) []string {
	slugs := make([]string, len(view.LandingPages))
	for i, lp := range view.LandingPages {
		slugs[i] = lp.Slug
	}
	return slugs
}

// The index is the request locale's categories only.
func TestListPublicCourseCategoriesShowsTheRequestLocale(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	categories, err := h.Client.ListPublicCourseCategories(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	slugs := make([]string, len(categories))
	for i, c := range categories {
		slugs[i] = c.Slug.Value
		assert.Equal(t, "ru", c.Locale.Value)
	}
	assert.ElementsMatch(t, []string{"layouting-ru", "programming-ru", "frontend-ru"}, slugs)
}

// The page lists the category's ru pages by course order; ruby, python and go
// share an order and fall back to course id. The en page grouped under the
// category is left out.
func TestGetPublicCourseCategoryListsItsLandingPagesInCourseOrder(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	view := getCategoryView(t, h, "programming-ru")

	assert.Equal(t, "programming-ru", view.Category.Slug.Value)
	assert.Equal(t,
		[]string{"typescript-ru", "ruby-ru", "python-ru", "go-ru"},
		categoryLandingSlugs(view))
	assert.Empty(t, view.QnaItems)
}

// Legacy `web.where(listed: true)`: an unlisted page, an unpublished one and a
// page whose course is not completed drop off the category page.
func TestGetPublicCourseCategoryLeavesOutHiddenPagesAndUnfinishedCourses(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")
	ctx := t.Context()

	require.NoError(t, h.DB.LandingPage.Update().
		Where(landingpage.SlugEQ("go-ru")).
		SetListed(false).
		Exec(ctx))
	require.NoError(t, h.DB.LandingPage.Update().
		Where(landingpage.SlugEQ("ruby-ru")).
		SetState("archived").
		Exec(ctx))
	require.NoError(t, h.DB.Course.Update().
		Where(course.SlugEQ("typescript")).
		SetProgress("in_development").
		Exec(ctx))

	view := getCategoryView(t, h, "programming-ru")
	assert.Equal(t, []string{"python-ru"}, categoryLandingSlugs(view))
}

// The questions come with the page, oldest first.
func TestGetPublicCourseCategoryCarriesItsQuestions(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	view := getCategoryView(t, h, "layouting-ru")

	questions := make([]string, len(view.QnaItems))
	for i, q := range view.QnaItems {
		questions[i] = q.Question
	}
	assert.Equal(t, []string{"Что такое вёрстка?", "Нужен ли опыт?"}, questions)
	assert.Empty(t, view.LandingPages)
}

func TestGetPublicCourseCategoryIsNotFoundOutsideTheLocale(t *testing.T) {
	cases := map[string]struct{ locale, slug string }{
		"another locale": {"en", "programming-ru"},
		"missing":        {"ru", "no-such-category"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			h := testsupport.NewAnonymousHarness(t)
			testsupport.SpeakTo(h, c.locale)

			_, err := h.Client.GetPublicCourseCategory(t.Context(), api.GetPublicCourseCategoryParams{Slug: c.slug})
			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, h.LastStatus())
		})
	}
}

package landingpage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/testsupport"
)

func newCatalog(t *testing.T) *landingpage.Catalog {
	t.Helper()

	return landingpage.NewCatalog(testsupport.NewClient(t))
}

// The catalog is the published, listed pages of one locale: drafts, archived
// pages and pages of another locale never appear. Fixtures carry both a draft
// and an archived ru page.
func TestCatalogServesOnlyPublishedPagesOfTheLocale(t *testing.T) {
	catalog := newCatalog(t)

	for _, locale := range []string{"ru", "en"} {
		pages, err := catalog.Pages(t.Context(), locale)
		require.NoError(t, err)
		require.NotEmpty(t, pages, locale)

		for _, page := range pages {
			assert.Equal(t, locale, *page.Locale, *page.Slug)
			assert.Equal(t, landingpage.StatePublished, *page.State, *page.Slug)
			assert.NotEqual(t, "nodejs-ru", *page.Slug, "draft")
			assert.NotEqual(t, "javascript-ru-archived", *page.Slug, "archived")
		}
	}
}

// Order is the Course display order, not the id: typescript (order 1) leads
// even though its Course id is larger than javascript's. Each page arrives with
// the Course and current Version the converter reads.
func TestCatalogOrdersByCourseOrderAndLoadsTheCourse(t *testing.T) {
	pages, err := newCatalog(t).Pages(t.Context(), "ru")
	require.NoError(t, err)
	require.NotEmpty(t, pages)

	assert.Equal(t, "typescript-ru", *pages[0].Slug)

	for _, page := range pages {
		require.NotNil(t, page.Edges.Course, *page.Slug)
	}

	for i := 1; i < len(pages); i++ {
		if pages[i].CourseID == pages[i-1].CourseID {
			assert.Less(t, pages[i-1].ID, pages[i].ID, "same-course pages ascend by id")
		}
	}
}

// within narrows the catalog and cannot widen it: a predicate matching every
// page still leaves the draft out.
func TestCatalogWithinNeverWidensTheScope(t *testing.T) {
	catalog := newCatalog(t)

	all, err := catalog.Pages(t.Context(), "ru")
	require.NoError(t, err)

	widened, err := catalog.Pages(t.Context(), "ru", courselandingpage.Or(courselandingpage.SlugEQ("nodejs-ru"), courselandingpage.SlugNEQ("")))
	require.NoError(t, err)
	assert.Len(t, widened, len(all))

	one, err := catalog.Pages(t.Context(), "ru", courselandingpage.SlugEQ("typescript-ru"))
	require.NoError(t, err)
	require.Len(t, one, 1)
}

// MainPages answers every course in one call and leaves out a course with no
// public main page; a draft main page is no page.
func TestMainPagesKeyedByCourseWithoutDraftsOrForeignLocales(t *testing.T) {
	catalog := newCatalog(t)

	listed, err := catalog.Pages(t.Context(), "ru")
	require.NoError(t, err)

	ids := make([]int, 0, len(listed))
	for _, page := range listed {
		ids = append(ids, page.CourseID)
	}

	mains, err := catalog.MainPages(t.Context(), "ru", ids...)
	require.NoError(t, err)
	require.NotEmpty(t, mains)

	for courseID, page := range mains {
		assert.Equal(t, courseID, page.CourseID)
		assert.True(t, *page.Main, *page.Slug)
		assert.Equal(t, "ru", *page.Locale)
		assert.Equal(t, landingpage.StatePublished, *page.State)
		require.NotNil(t, page.Edges.Course)
	}

	none, err := catalog.MainPages(t.Context(), "ru")
	require.NoError(t, err)
	assert.Empty(t, none)
}

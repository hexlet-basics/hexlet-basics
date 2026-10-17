package protosqlc_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/protosqlc"
	"hexletbasics/internal/protosqlc/db"
	"hexletbasics/internal/testsupport"
)

// The same rolled-back *sql.Tx the ent harness builds feeds sqlc's DBTX.
func newTx(t *testing.T) *sql.Tx {
	t.Helper()
	_, tx := testsupport.NewClientWithTx(t)
	return tx
}

func TestCatalogServesOnlyPublishedPagesOfTheLocale(t *testing.T) {
	catalog := protosqlc.NewCatalog(newTx(t))
	for _, locale := range []string{"ru", "en"} {
		pages, err := catalog.Pages(t.Context(), locale, nil)
		require.NoError(t, err)
		require.NotEmpty(t, pages, locale)
		for _, p := range pages {
			page := p.LanguageLandingPage
			assert.Equal(t, locale, page.Locale.String, page.Slug.String)
			assert.Equal(t, "published", page.State.String)
			assert.NotEqual(t, "nodejs-ru", page.Slug.String)
			assert.NotEqual(t, "javascript-ru-archived", page.Slug.String)
		}
	}
}

func TestCatalogOrdersByCourseOrderAndLoadsTheCourse(t *testing.T) {
	pages, err := protosqlc.NewCatalog(newTx(t)).Pages(t.Context(), "ru", nil)
	require.NoError(t, err)
	require.NotEmpty(t, pages)
	assert.Equal(t, "typescript-ru", pages[0].LanguageLandingPage.Slug.String)
	for i := 1; i < len(pages); i++ {
		a, b := pages[i-1].LanguageLandingPage, pages[i].LanguageLandingPage
		if a.LanguageID == b.LanguageID {
			assert.Less(t, a.ID, b.ID)
		}
	}
}

func TestCatalogCategoryNarrowing(t *testing.T) {
	catalog := protosqlc.NewCatalog(newTx(t))
	missing := int64(-1)
	none, err := catalog.Pages(t.Context(), "ru", &missing)
	require.NoError(t, err)
	assert.Empty(t, none)
}

func TestVersionsLoadedForListedCourses(t *testing.T) {
	catalog := protosqlc.NewCatalog(newTx(t))
	pages, err := catalog.Pages(t.Context(), "ru", nil)
	require.NoError(t, err)
	var courses []db.Language
	for _, p := range pages {
		courses = append(courses, p.Language)
	}
	versions, err := catalog.VersionsOf(t.Context(), courses)
	require.NoError(t, err)
	assert.NotEmpty(t, versions)
}

func TestMainPagesKeyedByCourse(t *testing.T) {
	catalog := protosqlc.NewCatalog(newTx(t))
	listed, err := catalog.Pages(t.Context(), "ru", nil)
	require.NoError(t, err)
	var ids []int64
	for _, p := range listed {
		ids = append(ids, p.LanguageLandingPage.LanguageID)
	}
	mains, err := catalog.MainPages(t.Context(), "ru", ids...)
	require.NoError(t, err)
	require.NotEmpty(t, mains)
	for id, p := range mains {
		assert.Equal(t, id, p.LanguageLandingPage.LanguageID)
		assert.True(t, p.LanguageLandingPage.Main.Bool)
	}
}

func TestDirectoryListNewestFirstAndPagesDoNotOverlap(t *testing.T) {
	d := protosqlc.NewDirectory(newTx(t))
	first, total, err := d.List(t.Context(), 0, 2)
	require.NoError(t, err)
	second, _, err := d.List(t.Context(), 2, 2)
	require.NoError(t, err)
	require.Len(t, first, 2)
	require.NotEmpty(t, second)
	assert.Greater(t, total, int64(2))
	assert.Greater(t, first[0].ID, first[1].ID)
	assert.Greater(t, first[1].ID, second[0].ID)
}

func TestDirectorySearch(t *testing.T) {
	d := protosqlc.NewDirectory(newTx(t))
	byEmail, err := d.Search(t.Context(), "SECOND-ADMIN@")
	require.NoError(t, err)
	require.Len(t, byEmail, 1)
	assert.Equal(t, int64(fixtures.ID(fixtures.UserSecondAdmin)), byEmail[0].ID)
	blank, err := d.Search(t.Context(), "")
	require.NoError(t, err)
	assert.Empty(t, blank)
}

func TestDirectoryDynamicFilterCombinations(t *testing.T) {
	d := protosqlc.NewDirectory(newTx(t))
	everyone, err := d.Filter(t.Context(), nil, nil, nil, 0, 1000)
	require.NoError(t, err)
	yes, no := true, false
	admins, err := d.Filter(t.Context(), nil, nil, &yes, 0, 1000)
	require.NoError(t, err)
	others, err := d.Filter(t.Context(), nil, nil, &no, 0, 1000)
	require.NoError(t, err)
	assert.NotEmpty(t, admins)
	assert.LessOrEqual(t, len(admins)+len(others), len(everyone)) // admin may be NULL
	term := "second"
	both, err := d.Filter(t.Context(), &term, nil, &yes, 0, 1000)
	require.NoError(t, err)
	assert.NotEmpty(t, both)
}

func TestEnrollUpsertIsIdempotent(t *testing.T) {
	tx := newTx(t)
	u, c := int64(fixtures.ID(fixtures.UserSecondAdmin)), int64(fixtures.ID(fixtures.CourseGo))
	var existing int64
	_ = tx.QueryRowContext(t.Context(), `select count(*) from language_members where user_id=$1 and language_id=$2`, u, c).Scan(&existing)
	id, created, err := protosqlc.Enroll(t.Context(), tx, u, c)
	require.NoError(t, err)
	if existing == 0 {
		assert.True(t, created)
		assert.NotZero(t, id)
	}
	_, created, err = protosqlc.Enroll(t.Context(), tx, u, c)
	require.NoError(t, err)
	assert.False(t, created, "second call loses the conflict quietly")
}

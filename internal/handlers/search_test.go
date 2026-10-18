package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

func searchIDs(page *api.SearchResultPage) []string {
	ids := make([]string, len(page.Items))
	for i, item := range page.Items {
		ids[i] = item.ID
	}
	return ids
}

// An anonymous visitor searches in the language of the request and gets compact
// results with the page of each record.
func TestSearchCatalogueAnswersAVisitorInTheRequestLocale(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	page, err := h.Client.SearchCatalogue(t.Context(), api.SearchCatalogueParams{Q: "шаблонным кодом"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, []string{"post:learning-ruby"}, searchIDs(page))
	got := page.Items[0]
	assert.Equal(t, api.SearchResultTypePost, got.Type)
	assert.Equal(t, "/ru/blog_posts/learning-ruby", got.Path)
	assert.Equal(t, "Learning Ruby", got.Title)
	assert.Contains(t, got.Snippet, "шаблонным")
	assert.Equal(t, int32(1), page.Total)
	assert.Equal(t, int32(1), page.Page)
	assert.Equal(t, int32(20), page.PerPage)
}

func TestSearchCatalogueNarrowsToOneType(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	page, err := h.Client.SearchCatalogue(t.Context(), api.SearchCatalogueParams{
		Q: "Привет", Type: api.NewOptSearchResultType(api.SearchResultTypeLesson),
	})
	require.NoError(t, err)

	require.NotEmpty(t, page.Items)
	for _, item := range page.Items {
		assert.Equal(t, api.SearchResultTypeLesson, item.Type)
	}
}

func TestSearchCataloguePagesTheAnswer(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")
	whole, err := h.Client.SearchCatalogue(t.Context(), api.SearchCatalogueParams{Q: "Привет"})
	require.NoError(t, err)
	require.Greater(t, whole.Total, int32(2))

	second, err := h.Client.SearchCatalogue(t.Context(), api.SearchCatalogueParams{
		Q: "Привет", Page: api.NewOptInt32(2), PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, searchIDs(whole)[2:4], searchIDs(second))
	assert.Equal(t, whole.Total, second.Total)
	assert.Equal(t, int32(2), second.Page)
}

// Nothing typed is an error: unbalanced quotes and operators are just words.
func TestSearchCatalogueTakesAnyTextButAnEmptyOne(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	page, err := h.Client.SearchCatalogue(t.Context(), api.SearchCatalogueParams{Q: `"unclosed & | ! (`})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Empty(t, page.Items)

	_, err = h.Client.SearchCatalogue(t.Context(), api.SearchCatalogueParams{Q: ""})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, h.LastStatus())
}

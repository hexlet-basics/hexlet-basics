package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// reviewNames projects reviews onto the reviewer's display name, the business
// fact the order assertions compare (fixture ids are opaque).
func reviewNames(reviews []api.Review) []string {
	names := make([]string, len(reviews))
	for i, r := range reviews {
		names[i] = r.FullName.Value
	}
	return names
}

// The ru list is the published ru reviews, newest first: the newer unpinned
// review sorts above the pinned one (legacy ordered by id only), and the ru
// draft is not shown.
func TestListPublicReviewsShowsPublishedReviewsOfTheRequestLocaleNewestFirst(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	page, err := h.Client.ListPublicReviews(t.Context(), api.ListPublicReviewsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, []string{"Maria Ivanova", "Ivan Petrov"}, reviewNames(page.Items))
	assert.Equal(t, int32(2), page.Total)
	assert.Equal(t, int32(1), page.Page)
	assert.Equal(t, int32(20), page.PerPage, "legacy pagy default")
}

func TestListPublicReviewsInEnglishListsOnlyEnglishReviews(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "en")

	page, err := h.Client.ListPublicReviews(t.Context(), api.ListPublicReviewsParams{})
	require.NoError(t, err)

	// John Doe is en but a draft.
	assert.Equal(t, []string{"Jane Smith"}, reviewNames(page.Items))
}

// Reviews are written in ru or en only, so the es page is empty, not missing.
func TestListPublicReviewsInSpanishIsEmpty(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "es")

	page, err := h.Client.ListPublicReviews(t.Context(), api.ListPublicReviewsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Empty(t, page.Items)
	assert.Equal(t, int32(0), page.Total)
}

// Pages slice the same ordered list; the total counts every page.
func TestListPublicReviewsPaginates(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	page, err := h.Client.ListPublicReviews(t.Context(), api.ListPublicReviewsParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(1),
	})
	require.NoError(t, err)

	assert.Equal(t, []string{"Ivan Petrov"}, reviewNames(page.Items))
	assert.Equal(t, int32(2), page.Total)
	assert.Equal(t, int32(2), page.Page)

	// Past the last page is an empty page, as legacy's pagy (overflow rescue
	// commented out) rendered it.
	page, err = h.Client.ListPublicReviews(t.Context(), api.ListPublicReviewsParams{
		Page:    api.NewOptInt32(3),
		PerPage: api.NewOptInt32(1),
	})
	require.NoError(t, err)
	assert.Empty(t, page.Items)
}

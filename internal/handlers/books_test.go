package handlers_test

import (
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/bookrequest"
	"hexletbasics/internal/api"
	"hexletbasics/internal/books"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// The harness user, Alice, has no book request; the fixture gives one to Bob.

func TestGetBookAnswersAVisitor(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	view, err := h.Client.GetBook(t.Context())
	require.NoError(t, err)
	assert.False(t, view.Requested)
	assert.Equal(t, "https://code-basics.com/ru/book", view.URL)
}

func TestGetBookIsNotRequestedByAnotherUsersRequest(t *testing.T) {
	h := testsupport.NewHarness(t)

	view, err := h.Client.GetBook(t.Context())
	require.NoError(t, err)
	assert.False(t, view.Requested)
}

func TestCreateBookRequestStoresTheRequestAndPublishesIt(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
	assert.IsType(t, &api.CreateBookRequestNoContent{}, res)

	stored := h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).OnlyX(t.Context())
	assert.Equal(t, books.StateRequested, lo.FromPtr(stored.State))

	require.Len(t, h.Events.Published, 1)
	published, ok := h.Events.Published[0].(events.BookRequested)
	require.True(t, ok, "expected BookRequested, got %T", h.Events.Published[0])
	assert.Equal(t, "ru", published.Locale)

	view, err := h.Client.GetBook(t.Context())
	require.NoError(t, err)
	assert.True(t, view.Requested)
}

func TestCreateBookRequestTwiceKeepsOneRequestAndOneEvent(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	_, err = h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	assert.Equal(t, 1, h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).CountX(t.Context()))
	assert.Len(t, h.Events.Published, 1)
}

func TestCreateBookRequestRequiresASignedInUser(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
	assert.IsType(t, &api.ProblemDetails{}, res)
	assert.Empty(t, h.Events.Published)
}

func TestDownloadBookRedirectsToThePDFAndMarksItDownloaded(t *testing.T) {
	h := testsupport.NewHarness(t)
	_, err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)

	res, err := h.Client.DownloadBook(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, h.LastStatus())
	found, ok := res.(*api.DownloadBookFound)
	require.True(t, ok, "expected a redirect, got %T", res)
	// The harness bucket cannot presign, so the PDF is served by the storage
	// read route; on S3 this is a presigned bucket URL.
	assert.Equal(t, "http://localhost:3001/api/storage/book.pdf", found.Location)

	stored := h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).OnlyX(t.Context())
	assert.Equal(t, books.StateDownloaded, lo.FromPtr(stored.State))
}

func TestDownloadBookWithoutARequestReturnsToTheBookPage(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.DownloadBook(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, h.LastStatus())
	found, ok := res.(*api.DownloadBookFound)
	require.True(t, ok, "expected a redirect, got %T", res)
	assert.Equal(t, "https://code-basics.com/ru/book", found.Location)

	assert.False(t, h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).ExistX(t.Context()))
}

func TestDownloadBookRequiresASignedInUser(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.DownloadBook(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
	assert.IsType(t, &api.ProblemDetails{}, res)
}

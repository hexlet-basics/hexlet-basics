package handlers_test

import (
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/bookrequest"
	"hexletbasics/internal/book"
	"hexletbasics/internal/testsupport"
)

// The harness user (user_admin_learner) has no book request; the fixture gives one to user_staff_editor.

func TestGetBookAnswersAVisitor(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)

	view, err := h.Client.GetBook(t.Context())
	require.NoError(t, err)
	assert.False(t, view.Requested)
}

func TestGetBookIsNotRequestedByAnotherUsersRequest(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	view, err := h.Client.GetBook(t.Context())
	require.NoError(t, err)
	assert.False(t, view.Requested)
}

func TestCreateBookRequestStoresTheRequest(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	stored := h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).OnlyX(t.Context())
	assert.Equal(t, book.StateRequested, lo.FromPtr(stored.State))

	view, err := h.Client.GetBook(t.Context())
	require.NoError(t, err)
	assert.True(t, view.Requested)
}

func TestCreateBookRequestTwiceKeepsOneRequest(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	err = h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	assert.Equal(t, 1, h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).CountX(t.Context()))
}

func TestCreateBookRequestRequiresASignedInUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	before := h.DB.BookRequest.Query().CountX(t.Context())

	err := h.Client.CreateBookRequest(t.Context())
	requireProblem(t, err, http.StatusUnauthorized)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
	assert.Equal(t, before, h.DB.BookRequest.Query().CountX(t.Context()))
}

func TestDownloadBookRedirectsToThePDFAndMarksItDownloaded(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	err := h.Client.CreateBookRequest(t.Context())
	require.NoError(t, err)

	res, err := h.Client.DownloadBook(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, h.LastStatus())

	found := res
	// The harness bucket cannot presign, so the PDF is served by the storage
	// read route; on S3 this is a presigned bucket URL.
	assert.Equal(t, "http://localhost:3001/api/storage/book.pdf", found.Location)

	stored := h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).OnlyX(t.Context())
	assert.Equal(t, book.StateDownloaded, lo.FromPtr(stored.State))
}

func TestDownloadBookWithoutARequestReturnsToTheBookPage(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	res, err := h.Client.DownloadBook(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, h.LastStatus())

	found := res
	assert.Equal(t, "https://code-basics.com/ru/book", found.Location)

	assert.False(t, h.DB.BookRequest.Query().Where(bookrequest.UserID(h.UserID)).ExistX(t.Context()))
}

func TestDownloadBookRequiresASignedInUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.DownloadBook(t.Context())
	requireProblem(t, err, http.StatusUnauthorized)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
}

package books_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/books"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
	"hexletbasics/internal/testsupport/testdb"
)

func TestMain(m *testing.M) {
	os.Exit(testdb.Run(m))
}

// user_blank has no book request in the fixtures; user_staff_editor already has one.
var learnerID = fixtures.ID(fixtures.UserBlank)

func newRecorder(t *testing.T) *books.Recorder {
	t.Helper()
	_, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	return books.NewRecorder(transactor, testsupport.NewServices(t, tx, transactor).Events)
}

func TestFunnelRunsFromRequestToDownload(t *testing.T) {
	recorder := newRecorder(t)

	requested, err := recorder.Requested(t.Context(), learnerID)
	require.NoError(t, err)
	assert.False(t, requested, "a visitor has no request")

	marked, err := recorder.MarkDownloaded(t.Context(), learnerID)
	require.NoError(t, err)
	assert.False(t, marked, "nothing to mark before the request")

	require.NoError(t, recorder.Request(t.Context(), learnerID, "ru"))
	require.NoError(t, recorder.Request(t.Context(), learnerID, "ru"), "a repeat changes nothing")

	requested, err = recorder.Requested(t.Context(), learnerID)
	require.NoError(t, err)
	assert.True(t, requested)

	marked, err = recorder.MarkDownloaded(t.Context(), learnerID)
	require.NoError(t, err)
	assert.True(t, marked)
}

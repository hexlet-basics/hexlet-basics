package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/handlers"
	"hexletbasics/internal/testsupport"
)

// goneClient is a connection whose reader has left: every write and flush fails.
type goneClient struct {
	header http.Header
	status int
}

var errGone = errors.New("client gone")

func (c *goneClient) Header() http.Header {
	if c.header == nil {
		c.header = http.Header{}
	}
	return c.header
}
func (c *goneClient) WriteHeader(status int)    { c.status = status }
func (c *goneClient) Write([]byte) (int, error) { return 0, errGone }

// flakyFlusher takes the bytes but cannot push them out.
type flakyFlusher struct{ *httptest.ResponseRecorder }

func (flakyFlusher) FlushError() error { return errGone }

// OPTIONS is answered with no content, whatever the route allows; any other
// method is a 405 that names what is allowed.
func TestMethodNotAllowedHandlerAnswersOptionsWithNoContent(t *testing.T) {
	handler := handlers.NewMethodNotAllowedHandler(testsupport.NewTranslator(t))

	options := httptest.NewRecorder()
	handler(options, httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/api/me", nil), "GET")
	assert.Equal(t, http.StatusNoContent, options.Code)

	refused := httptest.NewRecorder()
	handler(refused, httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/me", nil), "GET")
	assert.Equal(t, http.StatusMethodNotAllowed, refused.Code)
	assert.Equal(t, "GET", refused.Header().Get("Allow"))
}

// The error answer is still sent when the client has gone: the handler logs the
// failed write and returns instead of panicking.
func TestErrorHandlerSurvivesAClientThatLeft(t *testing.T) {
	handler := testsupport.NewAPIErrorHandler(t, testsupport.NewTranslator(t))
	client := &goneClient{}

	assert.NotPanics(t, func() {
		handler.Write(t.Context(), client, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil), io.ErrUnexpectedEOF)
	})
	assert.Equal(t, http.StatusInternalServerError, client.status)
}

// A streamed answer is flushed as it is written, and a client that is gone
// stops the copy feeding it.
func TestStreamResponsesFlushesEachWriteAndReportsAGoneClient(t *testing.T) {
	var written int
	var writeErr error
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		written, writeErr = w.Write([]byte("chunk"))
	})
	serve := func(w http.ResponseWriter) {
		handlers.StreamResponses(next).ServeHTTP(w, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))
	}

	recorder := httptest.NewRecorder()
	serve(recorder)
	require.NoError(t, writeErr)
	assert.Equal(t, 5, written)
	assert.True(t, recorder.Flushed, "the chunk is on the wire, not buffered")
	assert.Equal(t, "no", recorder.Header().Get("X-Accel-Buffering"), "nginx must not buffer it either")

	serve(&goneClient{})
	require.ErrorIs(t, writeErr, errGone)

	serve(flakyFlusher{httptest.NewRecorder()})
	require.ErrorIs(t, writeErr, errGone, "a flush that fails means the client is gone too")
}

// A status written after the body has begun is a failure mid-answer: the
// connection is cut rather than a JSON error glued to a short answer.
func TestStreamResponsesCutsAnAnswerThatFailsMidway(t *testing.T) {
	failMidway := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("partial"))
		w.WriteHeader(http.StatusInternalServerError)
	})

	assert.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handlers.StreamResponses(failMidway).ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil),
		)
	})
}

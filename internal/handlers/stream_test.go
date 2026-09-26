package handlers_test

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/handlers"
	"hexletbasics/internal/telemetry"
)

// A streamed response reaches the client chunk by chunk, through the same
// Sentry middleware that wraps every production request, and outlives the
// server's write timeout: an answer takes as long as the model takes.
func TestStreamResponsesFlushesEachWriteAndOutlivesTheWriteTimeout(t *testing.T) {
	release := make(chan struct{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "first\n")
		<-release
		_, _ = io.WriteString(w, "second\n")
	})
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)

	server := httptest.NewUnstartedServer(telemetry.NewSentryHTTPHandler(sentryClient, handlers.StreamResponses(inner)))
	server.Config.WriteTimeout = 100 * time.Millisecond
	server.Start()
	t.Cleanup(server.Close)
	// Registered after Close so it runs first: a failed assertion must still
	// let the handler finish, or Close waits on it forever.
	var releaseOnce sync.Once
	unblock := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(unblock)

	// Unflushed, not even the headers leave the server, so the request itself
	// is what has to be bounded in time.
	type firstChunk struct {
		res  *http.Response
		body *bufio.Reader
		line string
		err  error
	}
	arrived := make(chan firstChunk, 1)
	go func() {
		res, err := http.Get(server.URL)
		if err != nil {
			arrived <- firstChunk{err: err}
			return
		}
		body := bufio.NewReader(res.Body)
		line, err := body.ReadString('\n')
		arrived <- firstChunk{res: res, body: body, line: line, err: err}
	}()
	var first firstChunk
	select {
	case first = <-arrived:
	case <-time.After(5 * time.Second):
		t.Fatal("the first chunk was buffered")
	}
	require.NoError(t, first.err)
	t.Cleanup(func() { _ = first.res.Body.Close() })
	assert.Equal(t, "first\n", first.line, "the first chunk arrives while the handler is still writing")
	assert.Equal(t, "no", first.res.Header.Get("X-Accel-Buffering"), "nginx must not buffer the stream either")

	// Past the write timeout before the handler writes again.
	time.Sleep(300 * time.Millisecond)
	unblock()
	rest, err := io.ReadAll(first.body)
	require.NoError(t, err)
	assert.Equal(t, "second\n", string(rest))
}

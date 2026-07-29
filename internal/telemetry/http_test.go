package telemetry_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/telemetry"
)

func TestSentryHTTPHandlerReportsAndRepanics(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)

	handler := telemetry.NewSentryHTTPHandler(client, http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		hub := sentry.GetHubFromContext(r.Context())
		require.Same(t, client, hub.Client())
		panic("boom")
	}))
	request := httptest.NewRequest(http.MethodPost, "/panic?source=test", nil)

	assert.PanicsWithValue(t, "boom", func() {
		handler.ServeHTTP(httptest.NewRecorder(), request)
	})

	events := transport.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "boom", events[0].Message)
	assert.Equal(t, http.MethodPost, events[0].Request.Method)
	assert.Equal(t, "http://example.com/panic", events[0].Request.URL)
	assert.Equal(t, "source=test", events[0].Request.QueryString)
}

func TestSentryHTTPHandlerIsolatesRequestHubs(t *testing.T) {
	client, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)

	hubs := make(chan *sentry.Hub, 2)
	handler := telemetry.NewSentryHTTPHandler(client, http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		hub := sentry.GetHubFromContext(r.Context())
		require.NotNil(t, hub)
		require.Same(t, client, hub.Client())
		hubs <- hub
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/first", nil))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/second", nil))

	first := <-hubs
	second := <-hubs
	assert.NotSame(t, first, second)
	assert.NotSame(t, first.Scope(), second.Scope())
}

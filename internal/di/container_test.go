package di

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/events"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/testsupport/testdb"
)

func configureTestEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", testdb.DatabaseURL())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("JWT_SECRET", "test-secret")
}

// TestServerContainerResolvesHTTPWithoutAsyncRuntime proves that the synchronous
// graph has everything needed to serve and enqueue while exposing neither
// Watermill consumers nor a startable River worker runtime.
func TestServerContainerResolvesHTTPWithoutAsyncRuntime(t *testing.T) {
	configureTestEnvironment(t)

	injector := newServerContainer()
	t.Cleanup(func() { _ = injector.Shutdown() })

	srv, err := do.Invoke[*http.Server](injector)
	require.NoError(t, err)
	assert.NotNil(t, srv.Handler)

	_, err = do.Invoke[*events.Runtime](injector)
	require.Error(t, err)
}

func TestWorkerContainerResolvesAsyncRuntimeWithoutHTTP(t *testing.T) {
	configureTestEnvironment(t)

	injector := newWorkerContainer()
	t.Cleanup(func() { _ = injector.Shutdown() })

	riverClient, err := do.Invoke[*river.Client[*sql.Tx]](injector)
	require.NoError(t, err)
	assert.NotNil(t, riverClient)

	eventRuntime, err := do.Invoke[*events.Runtime](injector)
	require.NoError(t, err)
	assert.NotNil(t, eventRuntime)

	_, err = do.Invoke[*http.Server](injector)
	require.Error(t, err)
}

func TestDevCORSAllowsCredentialedXSRFRequests(t *testing.T) {
	configureTestEnvironment(t)

	injector := newServerContainer()
	t.Cleanup(func() { _ = injector.Shutdown() })

	srv, err := do.Invoke[*http.Server](injector)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodOptions, "/session", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodDelete)
	req.Header.Set("Access-Control-Request-Headers", "content-type,x-xsrf-token")
	rec := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers")), "x-xsrf-token")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodDelete)
}

func TestHTTPHandlerReportsAndRepanics(t *testing.T) {
	configureTestEnvironment(t)

	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)

	injector := newServerContainer()
	t.Cleanup(func() { _ = injector.Shutdown() })
	do.OverrideValue(injector, client)
	do.OverrideNamedValue[http.Handler](injector, routerServiceName, http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			panic("router panic")
		},
	))

	handler, err := do.InvokeNamed[http.Handler](injector, httpHandlerServiceName)
	require.NoError(t, err)

	assert.PanicsWithValue(t, "router panic", func() {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/panic", nil))
	})

	events := transport.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "router panic", events[0].Message)
	assert.Equal(t, "http://example.com/panic", events[0].Request.URL)
}

func TestServerProviderReturnsDependencyError(t *testing.T) {
	configureTestEnvironment(t)

	wantErr := errors.New("load localization")
	injector := newServerContainer()
	do.Override(injector, func(do.Injector) (*localization.Translator, error) {
		return nil, wantErr
	})

	_, err := do.Invoke[*http.Server](injector)
	require.ErrorIs(t, err, wantErr)
}

func TestBuildServerReturnsStartupError(t *testing.T) {
	configureTestEnvironment(t)
	t.Setenv("BLOB_BUCKET_URL", "unsupported://bucket")

	_, err := BuildServer()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

func TestBuildWorkerReturnsStartupError(t *testing.T) {
	configureTestEnvironment(t)
	t.Setenv("BLOB_BUCKET_URL", "unsupported://bucket")

	_, err := BuildWorker()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

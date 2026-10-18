package di

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
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

// Faults a test injects; each is told apart from the others by identity.
var (
	errWant = errors.New("load localization")
)

// configureWorkerEnvironment is everything the worker is given: no cookie
// secret and no public host, which only the server needs.
func configureWorkerEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", testdb.DatabaseURL())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("PUBLIC_URL", "http://localhost:3001")
	t.Setenv("SITE_URL", "http://localhost:5173")
	t.Setenv("EMAIL_TOKEN_SECRET", "test-email-secret")
}

func configureTestEnvironment(t *testing.T) {
	t.Helper()
	configureWorkerEnvironment(t)
	t.Setenv("GUEST_COOKIE_SECRET", "test-secret")
	t.Setenv("APP_HOST", "code-basics.com")
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
	configureWorkerEnvironment(t)

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
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/panic", nil))
	})

	events := transport.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "router panic", events[0].Message)
	assert.Equal(t, "http://example.com/panic", events[0].Request.URL)
}

func TestServerProviderReturnsDependencyError(t *testing.T) {
	configureTestEnvironment(t)

	injector := newServerContainer()
	do.Override(injector, func(do.Injector) (*localization.Translator, error) {
		return nil, errWant
	})

	_, err := do.Invoke[*http.Server](injector)
	require.ErrorIs(t, err, errWant)
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

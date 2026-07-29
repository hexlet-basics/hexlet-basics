package di_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/di"
	"hexletbasics/internal/events"
)

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return defaultTestDSN
}

// TestServerContainerResolvesHTTPWithoutAsyncRuntime proves that the synchronous
// graph has everything needed to serve and enqueue while exposing neither
// Watermill consumers nor a startable River worker runtime.
func TestServerContainerResolvesHTTPWithoutAsyncRuntime(t *testing.T) {
	t.Setenv("DATABASE_URL", testDSN())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())

	injector := di.NewServer()
	t.Cleanup(func() { _ = injector.Shutdown() })

	srv, err := do.Invoke[*http.Server](injector)
	require.NoError(t, err)
	assert.NotNil(t, srv.Handler)

	_, err = do.Invoke[*events.Runtime](injector)
	require.Error(t, err)
}

func TestWorkerContainerResolvesAsyncRuntimeWithoutHTTP(t *testing.T) {
	t.Setenv("DATABASE_URL", testDSN())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())

	injector := di.NewWorker()
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
	t.Setenv("DATABASE_URL", testDSN())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())

	injector := di.NewServer()
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

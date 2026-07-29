package di_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/di"
)

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return defaultTestDSN
}

// TestContainerResolvesHTTPServer proves the whole provider graph wires up: the
// *http.Server depends (transitively) on the api server, handlers.Server, the
// river client, the exercise Loader, the git fetcher, the blob bucket and the ent
// client. Resolving it exercises every provider added for course loading and
// catches a missing binding or a `do` cycle — which would otherwise only surface
// at server boot in production. Construction is lazy (no queries run here); the
// bucket points at a throwaway dir so nothing touches real storage.
func TestContainerResolvesHTTPServer(t *testing.T) {
	t.Setenv("DATABASE_URL", testDSN())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())

	injector := di.New()
	t.Cleanup(func() { _ = injector.Shutdown() })

	srv, err := do.Invoke[*http.Server](injector)
	require.NoError(t, err)
	assert.NotNil(t, srv.Handler)
}

func TestDevCORSAllowsCredentialedXSRFRequests(t *testing.T) {
	t.Setenv("DATABASE_URL", testDSN())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())

	injector := di.New()
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

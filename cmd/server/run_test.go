package main

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport/testdb"
)

// serverEnv is the environment a deployment gives the process, pointed at the
// test database and a scratch blob bucket.
func serverEnv(t *testing.T, addr string) {
	t.Helper()
	t.Setenv("ADDR", addr)
	t.Setenv("DATABASE_URL", testdb.DatabaseURL())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("GUEST_COOKIE_SECRET", "test-secret")
	t.Setenv("PUBLIC_URL", "http://localhost:3001")
	t.Setenv("SITE_URL", "http://localhost:5173")
	t.Setenv("APP_HOST", "code-basics.com")
	t.Setenv("EMAIL_TOKEN_SECRET", "test-email-secret")
}

func TestRunServesUntilSignalledAndStopsCleanly(t *testing.T) {
	addr := freeAddr(t)
	serverEnv(t, addr)

	signalCtx, signal := context.WithCancel(t.Context())
	defer signal()

	done := make(chan error, 1)
	go func() { done <- run(signalCtx, signal) }()

	require.Eventually(t, func() bool {
		res, err := get(t.Context(), "http://"+addr+"/api/languages")
		if err != nil {
			return false
		}

		_ = res.Body.Close()

		return res.StatusCode == http.StatusOK
	}, 30*time.Second, 50*time.Millisecond, "the process answers requests")

	signal()

	select {
	case err := <-done:
		require.NoError(t, err, "a signal is the normal way to stop, not a failure")
	case <-time.After(30 * time.Second):
		t.Fatal("the process did not stop after the signal")
	}
}

func TestRunFailsWhenTheDependenciesCannotBeBuilt(t *testing.T) {
	serverEnv(t, freeAddr(t))
	t.Setenv("BLOB_BUCKET_URL", "unsupported://bucket")

	err := run(t.Context(), func() {})

	require.Error(t, err)
	assert.ErrorContains(t, err, "building server dependencies")
}

func TestRunFailsWhenTheAddressIsTaken(t *testing.T) {
	taken, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	defer func() { assert.NoError(t, taken.Close()) }()

	serverEnv(t, taken.Addr().String())

	runErr := run(t.Context(), func() {})

	require.Error(t, runErr)
	assert.ErrorContains(t, runErr, "serve HTTP")
}

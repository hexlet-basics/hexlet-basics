package main

import (
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestMainServesUntilTheProcessIsAskedToStop runs main itself, as the platform
// does: it serves requests, and SIGTERM makes it drain, release its resources
// and return. A failed start would call os.Exit and take the test binary down.
func TestMainServesUntilTheProcessIsAskedToStop(t *testing.T) {
	addr := freeAddr(t)
	serverEnv(t, addr)

	returned := make(chan struct{})
	go func() {
		main()
		close(returned)
	}()

	require.Eventually(t, func() bool {
		res, err := get(t.Context(), "http://"+addr+"/api/languages")
		if err != nil {
			return false
		}
		_ = res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 30*time.Second, 50*time.Millisecond, "the process answers requests")

	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
	select {
	case <-returned:
	case <-time.After(30 * time.Second):
		t.Fatal("main did not return after SIGTERM")
	}
}

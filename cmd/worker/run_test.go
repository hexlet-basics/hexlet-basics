package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport/testdb"
)

func workerEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", testdb.DatabaseURL())
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("EMAIL_TOKEN_SECRET", "test-email-secret")
}

func TestRunConsumesUntilSignalledAndStopsCleanly(t *testing.T) {
	workerEnv(t)
	signalCtx, signal := context.WithCancel(t.Context())
	defer signal()

	done := make(chan error, 1)
	go func() { done <- run(signalCtx, signal) }()

	// The worker has no endpoint to probe; give River and the event router the
	// time to start, then stop it the way the platform does.
	time.Sleep(2 * time.Second)
	signal()

	select {
	case err := <-done:
		assert.NoError(t, err, "a signal is the normal way to stop, not a failure")
	case <-time.After(60 * time.Second):
		t.Fatal("the worker did not stop after the signal")
	}
}

func TestRunFailsWhenTheDependenciesCannotBeBuilt(t *testing.T) {
	workerEnv(t)
	t.Setenv("BLOB_BUCKET_URL", "unsupported://bucket")

	err := run(t.Context(), func() {})

	require.Error(t, err)
	assert.ErrorContains(t, err, "building worker dependencies")
}

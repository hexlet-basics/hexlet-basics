package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func workerEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://nobody@127.0.0.1:1/none")
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("GUEST_COOKIE_SECRET", "test-secret")
	t.Setenv("EMAIL_TOKEN_SECRET", "test-email-secret")
}

func TestRunFailsWhenTheDependenciesCannotBeBuilt(t *testing.T) {
	workerEnv(t)
	t.Setenv("BLOB_BUCKET_URL", "unsupported://bucket")

	err := run(t.Context(), func() {})

	require.Error(t, err)
	assert.ErrorContains(t, err, "building worker dependencies")
}

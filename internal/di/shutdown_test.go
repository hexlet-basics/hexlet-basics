package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerShutdownReleasesEveryResourceTheProcessOwns(t *testing.T) {
	configureTestEnvironment(t)
	deps, err := BuildServer()
	require.NoError(t, err)

	require.NoError(t, deps.Shutdown(t.Context()))

	require.Error(t, deps.Bucket.WriteAll(t.Context(), "key", []byte("x"), nil), "the bucket is closed")
	_, err = deps.Database.User.Query().Count(t.Context())
	assert.Error(t, err, "the database pool is closed")
}

func TestWorkerShutdownReleasesEveryResourceTheProcessOwns(t *testing.T) {
	configureTestEnvironment(t)
	deps, err := BuildWorker()
	require.NoError(t, err)

	require.NoError(t, deps.Shutdown(t.Context()))

	require.Error(t, deps.Bucket.WriteAll(t.Context(), "key", []byte("x"), nil), "the bucket is closed")
	_, err = deps.Database.User.Query().Count(t.Context())
	assert.Error(t, err, "the database pool is closed")
}

// A resource that cannot close is reported, not swallowed.
func TestShutdownReportsWhatFailsToClose(t *testing.T) {
	configureTestEnvironment(t)
	deps, err := BuildServer()
	require.NoError(t, err)
	require.NoError(t, deps.Bucket.Close())

	err = deps.Shutdown(t.Context())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "blob bucket")
}

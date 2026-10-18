package jobs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Discarded rows are the dead-letter store, so how long they stay is a
// decision of ours, not River's default (which may change under us). River
// exposes no accessor for a built client's config, hence this test of the
// config the worker client is built from.
func TestWorkerConfigKeepsDiscardedJobsForAnExplicitRetention(t *testing.T) {
	cfg := workerConfig(nil, nil, nil, nil, nil)

	require.Equal(t, DiscardedRetention, cfg.DiscardedJobRetentionPeriod)
	require.Equal(t, 30*24*time.Hour, DiscardedRetention)
}

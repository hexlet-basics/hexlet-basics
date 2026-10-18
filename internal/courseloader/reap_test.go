package courseloader_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// The loader's clock reads testsupport.Today (2026-07-27 12:00 UTC); the
// fixtures say how long ago each version last changed.
func TestReapStuckFailsOnlyBuildsIdlePastThreshold(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	stuck := fixture.ID(fixture.VersionLoadableStuck)
	running := fixture.ID(fixture.VersionLoadableBuilding)
	queued := fixture.ID(fixture.VersionLoadableQueued)
	built := fixture.ID(fixture.VersionLoadableBuilt)

	// The result quotes the timestamp as stored, before the reap bumps it.
	staleAt := db.CourseVersion.GetX(ctx, stuck).UpdatedAt

	loader := newLoaderWith(t, db, txStore, unreachableFetcher())
	reaped, err := loader.ReapStuck(ctx)
	require.NoError(t, err)
	require.Len(t, reaped, 1)
	assert.Equal(t, stuck, reaped[0].ID)

	version := db.CourseVersion.GetX(ctx, stuck)
	assert.Equal(t, "failed", derefStr(version.State))
	assert.Equal(t,
		"Build reaped: stuck in 'building' since "+staleAt.UTC().Format(time.RFC3339)+
			" (worker likely killed, e.g. OOM)",
		derefStr(version.Result),
	)

	for id, state := range map[int]string{running: "building", queued: "created", built: "built"} {
		version := db.CourseVersion.GetX(ctx, id)
		assert.Equal(t, state, derefStr(version.State))
		assert.Equal(t, id == built, version.Result != nil, "only the built version has a result")
	}
}

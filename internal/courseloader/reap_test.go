package courseloader_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/testsupport"
)

func TestReapStuckFailsOnlyBuildsIdlePastThreshold(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()

	course := db.Course.Create().SetSlug("reap-lang").SetName("Reap").SaveX(ctx)
	stale := time.Now().Add(-courseloader.StuckBuildAfter - time.Hour)
	fresh := time.Now().Add(-time.Hour)

	create := func(state string, updatedAt time.Time) int {
		return db.CourseVersion.Create().
			SetCourseID(course.ID).
			SetState(state).
			SetUpdatedAt(updatedAt).
			SaveX(ctx).ID
	}
	stuck := create("building", stale)
	running := create("building", fresh)
	queued := create("created", stale)
	built := create("built", stale)

	// The result quotes the timestamp as stored, before the reap bumps it.
	staleAt := db.CourseVersion.GetX(ctx, stuck).UpdatedAt

	loader := newLoaderWith(t, db, txStore, panicFetcher{})
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
		assert.Nil(t, version.Result)
	}
}

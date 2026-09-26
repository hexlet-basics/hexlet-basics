package jobs_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"hexletbasics/ent/course"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// TestWorkerReapsStuckVersionBuildsPeriodically proves the worker client
// schedules the reaper through River's periodic jobs: starting it runs the
// sweep (RunOnStart) without anything enqueuing it, failing the stuck build and
// leaving a healthy one alone. It uses a real pool, so it removes its rows.
func TestWorkerReapsStuckVersionBuildsPeriodically(t *testing.T) {
	ctx := context.Background()

	sqlDB, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)

	ruby, err := db.Course.Query().Where(course.Slug("ruby")).Only(ctx)
	require.NoError(t, err)
	stuck := db.CourseVersion.Create().
		SetCourseID(ruby.ID).
		SetState("building").
		SetUpdatedAt(time.Now().Add(-courseloader.StuckBuildAfter - time.Hour)).
		SaveX(ctx)
	healthy := db.CourseVersion.Create().
		SetCourseID(ruby.ID).
		SetState("building").
		SaveX(ctx)
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = sqlDB.ExecContext(cleanupCtx, "DELETE FROM river_job WHERE kind = 'reap_stuck_version_builds'")
		_, _ = sqlDB.ExecContext(cleanupCtx,
			"DELETE FROM language_versions WHERE id IN ($1, $2)", stuck.ID, healthy.ID)
		_ = sqlDB.Close()
	})

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	// The reaper only touches version rows, so the loader needs no store,
	// assets, fetcher or progress module.
	workerClient, err := jobs.NewWorkerClient(
		sqlDB,
		courseloader.NewLoader(db, nil, nil, nil, nil),
		nil,
		nil,
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		jobs.NewErrorHandler(sentryClient),
		nooptrace.NewTracerProvider(),
		noopmetric.NewMeterProvider(),
	)
	require.NoError(t, err)

	completed, cancel := workerClient.Subscribe(river.EventKindJobCompleted)
	defer cancel()

	require.NoError(t, workerClient.Start(ctx))
	t.Cleanup(func() {
		stopCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = workerClient.Stop(stopCtx)
	})

	timeout := time.After(15 * time.Second)
	for reaped := false; !reaped; {
		select {
		case ev := <-completed:
			reaped = ev.Job.Kind == jobs.ReapStuckVersionBuildsArgs{}.Kind()
		case <-timeout:
			t.Fatal("periodic reaper did not complete within 15s")
		}
	}

	stuck = db.CourseVersion.GetX(ctx, stuck.ID)
	assert.Equal(t, "failed", *stuck.State)
	require.NotNil(t, stuck.Result)
	assert.True(t, strings.HasPrefix(*stuck.Result, "Build reaped: stuck in 'building' since "))

	healthy = db.CourseVersion.GetX(ctx, healthy.ID)
	assert.Equal(t, "building", *healthy.State)
	assert.Nil(t, healthy.Result)
}

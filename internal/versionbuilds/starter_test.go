package versionbuilds_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
	"hexletbasics/internal/versionbuilds"
)

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

func testDSN() string {
	if value := os.Getenv("TEST_DATABASE_URL"); value != "" {
		return value
	}
	return defaultTestDSN
}

// TestStarterCommitsVersionAndJobTogether exercises ent and River over their
// shared database/sql transaction against the real test database.
func TestStarterCommitsVersionAndJobTogether(t *testing.T) {
	ctx := context.Background()
	sqlDB, err := sql.Open("pgx", testDSN())
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	db := store.NewClient(sqlDB)
	// A non-nil loader registers the exercise_loader worker. It is never started
	// in this test, so its own dependencies are intentionally nil.
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	riverClient, err := jobs.NewClient(
		sqlDB,
		courseloader.NewLoader(nil, nil, nil, nil),
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		jobs.NewErrorHandler(sentryClient),
	)
	require.NoError(t, err)
	starter := versionbuilds.NewStarter(store.New(sqlDB), riverClient)

	ruby, err := db.Course.Query().Where(course.Slug("ruby")).Only(ctx)
	require.NoError(t, err)

	version, err := starter.Start(ctx, ruby.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = sqlDB.ExecContext(context.Background(), "DELETE FROM river_job WHERE kind = $1 AND args->>'version_id' = $2", "exercise_loader", strconv.Itoa(version.ID))
		_, _ = sqlDB.ExecContext(context.Background(), "DELETE FROM language_versions WHERE id = $1", version.ID)
	})

	duplicate, err := riverClient.Insert(ctx, jobs.ExerciseLoaderArgs{VersionID: version.ID}, nil)
	require.NoError(t, err)
	assert.True(t, duplicate.UniqueSkippedAsDuplicate)

	saved, err := db.CourseVersion.Get(ctx, version.ID)
	require.NoError(t, err)
	require.NotNil(t, saved.State)
	assert.Equal(t, "created", *saved.State)

	var jobsCount int
	err = sqlDB.QueryRowContext(ctx,
		"SELECT count(*) FROM river_job WHERE kind = $1 AND args->>'version_id' = $2",
		"exercise_loader", strconv.Itoa(version.ID),
	).Scan(&jobsCount)
	require.NoError(t, err)
	assert.Equal(t, 1, jobsCount)
}

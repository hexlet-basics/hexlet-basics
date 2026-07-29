package jobs_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/jobs"
)

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return defaultTestDSN
}

// TestQueueBackbone proves the river backbone end to end against the test DB:
// a ping job is enqueued, the worker runs it, and a completion event fires.
// It uses a real pool (not the harness's rolled-back transaction), so it cleans up the
// river_job rows it creates.
func TestQueueBackbone(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("pgx", testDSN())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM river_job WHERE kind = 'ping'")
		_ = db.Close()
	})

	// nil loader: this smoke test only exercises the ping worker, so the
	// exercise-build worker (and its db/blob deps) is intentionally not registered.
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	client, err := jobs.NewWorkerClient(
		db,
		nil,
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		jobs.NewErrorHandler(sentryClient),
	)
	require.NoError(t, err)

	// Subscribe before Start so the completion event can't be missed.
	events, cancel := client.Subscribe(river.EventKindJobCompleted)
	defer cancel()

	require.NoError(t, client.Start(ctx))
	t.Cleanup(func() {
		stopCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = client.Stop(stopCtx)
	})

	_, err = client.Insert(ctx, jobs.PingArgs{}, nil)
	require.NoError(t, err)

	select {
	case ev := <-events:
		require.Equal(t, rivertype.JobStateCompleted, ev.Job.State)
		require.Equal(t, "ping", ev.Job.Kind)
	case <-time.After(10 * time.Second):
		t.Fatal("ping job did not complete within 10s")
	}
}

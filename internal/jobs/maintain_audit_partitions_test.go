package jobs_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"hexletbasics/internal/audit"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

func auditPartitionExists(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var exists bool
	require.NoError(t, db.QueryRowContext(t.Context(),
		"SELECT to_regclass($1) IS NOT NULL", name).Scan(&exists))
	return exists
}

// TestWorkerMaintainsAuditPartitionsPeriodically proves the worker runs the
// upkeep on its own (RunOnStart): with a one-month retention it drops an old
// partition and its record, keeps the month in progress, and has created the
// next month's partition before anything needs it. It uses a real pool, so it
// removes what it creates.
func TestWorkerMaintainsAuditPartitionsPeriodically(t *testing.T) {
	ctx := context.Background()

	sqlDB, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = sqlDB.ExecContext(cleanupCtx, "DELETE FROM river_job WHERE kind = 'maintain_audit_partitions'")
		_ = sqlDB.Close()
	})

	// A year-old partition holding a record: far outside a one-month retention.
	old := time.Now().UTC().AddDate(-1, 0, 0)
	partitions := audit.NewPartitions(sqlDB)
	require.NoError(t, partitions.EnsureAhead(ctx, old))
	oldRecord := db.AuditRecord.Create().
		SetOccurredAt(old).SetChannel("web").SetOperation("AdminUpdateBlogPost").SetOutcome("succeeded").
		SaveX(ctx)
	oldName := "audit_records_p" + old.Format("200601")
	require.True(t, auditPartitionExists(t, sqlDB, oldName))

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	workerClient, err := jobs.NewWorkerClient(
		sqlDB,
		logger,
		jobs.NewErrorHandler(sentryClient),
		nooptrace.NewTracerProvider(),
		noopmetric.NewMeterProvider(),
		jobs.MaintainAuditPartitions(partitions, 1, logger),
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
	for done := false; !done; {
		select {
		case ev := <-completed:
			done = ev.Job.Kind == jobs.MaintainAuditPartitionsArgs{}.Kind()
		case <-timeout:
			t.Fatal("periodic audit maintenance did not complete within 15s")
		}
	}

	assert.False(t, auditPartitionExists(t, sqlDB, oldName), "the expired partition is dropped")
	_, err = db.AuditRecord.Get(ctx, oldRecord.ID)
	assert.Error(t, err, "the record went with its partition")

	now := time.Now().UTC()
	assert.True(t, auditPartitionExists(t, sqlDB, "audit_records_p"+now.Format("200601")), "the month in progress is kept")
	assert.True(t, auditPartitionExists(t, sqlDB, "audit_records_p"+time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC).Format("200601")),
		"next month's partition exists before it is needed")
}

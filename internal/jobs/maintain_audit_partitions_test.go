package jobs_test

import (
	"database/sql"
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/audit"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

func auditPartitionExists(t *testing.T, tx *sql.Tx, name string) bool {
	t.Helper()
	var exists bool
	require.NoError(t, tx.QueryRowContext(t.Context(),
		"SELECT to_regclass($1) IS NOT NULL", name).Scan(&exists))
	return exists
}

// With a one-month retention the upkeep drops an old partition and its record,
// keeps the month in progress, and creates the next month's partition before
// anything needs it. Partitions are DDL, which is transactional, so the whole
// run is worked in the test's transaction and rolled back with it.
func TestMaintainAuditPartitionsJobDropsTheExpiredAndCreatesTheNext(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)

	// A year-old partition holding a record: far outside a one-month retention.
	old := time.Now().UTC().AddDate(-1, 0, 0)
	partitions := audit.NewPartitions(tx)
	require.NoError(t, partitions.EnsureAhead(t.Context(), old))
	//nolint:forbidigo // the row lands in a partition this test has just created, which no fixture can hold
	oldRecord := db.AuditRecord.Create().
		SetOccurredAt(old).SetChannel("web").SetOperation("AdminUpdateBlogPost").SetOutcome("succeeded").
		SaveX(t.Context())
	oldName := "audit_records_p" + old.Format("200601")
	require.True(t, auditPartitionExists(t, tx, oldName))

	result, err := work(t, jobs.MaintainAuditPartitionsWorker(partitions, 1, slog.New(slog.DiscardHandler)),
		jobs.MaintainAuditPartitionsArgs{})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
	assert.False(t, auditPartitionExists(t, tx, oldName), "the expired partition is dropped")
	_, err = db.AuditRecord.Get(t.Context(), oldRecord.ID)
	assert.Error(t, err, "the record went with its partition")

	now := time.Now().UTC()
	assert.True(t, auditPartitionExists(t, tx, "audit_records_p"+now.Format("200601")), "the month in progress is kept")
	assert.True(t, auditPartitionExists(t, tx, "audit_records_p"+time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC).Format("200601")),
		"next month's partition exists before it is needed")
}

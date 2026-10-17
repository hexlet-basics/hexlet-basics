package audit_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/auditrecord"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/testsupport"
)

// The tests pin "now" to a past date so they never collide with the partitions
// the migration created around the real current month, and they run inside the
// harness transaction, so every partition they create or drop rolls back.
var partitionNow = time.Date(2024, 10, 17, 12, 0, 0, 0, time.UTC)

func partitionNames(t *testing.T, tx *sql.Tx) []string {
	t.Helper()
	rows, err := tx.QueryContext(t.Context(), `SELECT c.relname FROM pg_inherits i
		JOIN pg_class c ON c.oid = i.inhrelid
		WHERE i.inhparent = 'audit_records'::regclass ORDER BY c.relname`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()
	var names []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	return names
}

// addRecordAt writes an Audit Record that happened at an instant, through the
// production writer, and returns the row.
func addRecordAt(t *testing.T, client *ent.Client, at time.Time) *ent.AuditRecord {
	t.Helper()
	require.NoError(t, audit.Write(t.Context(), client, audit.Entry{
		OccurredAt: at,
		Source:     audit.Source{Channel: audit.ChannelWeb},
		Operation:  "AdminUpdateBlogPost",
		Outcome:    audit.OutcomeSucceeded,
	}, nil))
	return client.AuditRecord.Query().Order(ent.Desc(auditrecord.FieldID)).FirstX(t.Context())
}

func partitionOf(t *testing.T, tx *sql.Tx, id int) string {
	t.Helper()
	var name string
	require.NoError(t, tx.QueryRowContext(t.Context(),
		"SELECT tableoid::regclass::text FROM audit_records WHERE id = $1", id).Scan(&name))
	return name
}

func TestEnsureAheadCreatesThisAndTheNextMonthsOnce(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	partitions := audit.NewPartitions(tx)

	require.NoError(t, partitions.EnsureAhead(t.Context(), partitionNow))
	require.NoError(t, partitions.EnsureAhead(t.Context(), partitionNow), "running twice is harmless")

	names := partitionNames(t, tx)
	for _, want := range []string{"audit_records_p202410", "audit_records_p202411", "audit_records_p202412", "audit_records_p202501"} {
		assert.Contains(t, names, want)
	}
	assert.NotContains(t, names, "audit_records_p202502")
	assert.Contains(t, names, "audit_records_default", "the safety net is never dropped or replaced")
}

func TestARecordLandsInThePartitionOfItsMonth(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	require.NoError(t, audit.NewPartitions(tx).EnsureAhead(t.Context(), partitionNow))

	next := addRecordAt(t, client, time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC))
	lastOfMonth := addRecordAt(t, client, time.Date(2024, 10, 31, 23, 59, 59, 0, time.UTC))
	firstOfNext := addRecordAt(t, client, time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, "audit_records_p202411", partitionOf(t, tx, next.ID))
	assert.Equal(t, "audit_records_p202410", partitionOf(t, tx, lastOfMonth.ID), "the month ends at midnight UTC")
	assert.Equal(t, "audit_records_p202411", partitionOf(t, tx, firstOfNext.ID))
}

func TestARecordForAMonthWithoutAPartitionIsKeptAndMovedWhenItsMonthIsCreated(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	orphan := addRecordAt(t, client, time.Date(2024, 11, 5, 8, 0, 0, 0, time.UTC))
	assert.Equal(t, "audit_records_default", partitionOf(t, tx, orphan.ID), "an insert never fails for a missing month")

	require.NoError(t, audit.NewPartitions(tx).EnsureAhead(t.Context(), partitionNow))

	assert.Equal(t, "audit_records_p202411", partitionOf(t, tx, orphan.ID))
	moved := client.AuditRecord.GetX(t.Context(), orphan.ID)
	assert.Equal(t, orphan.OccurredAt.UTC(), moved.OccurredAt.UTC(), "the record is carried over unchanged")
	assert.Contains(t, partitionNames(t, tx), "audit_records_default", "the default partition is attached again")

	// Re-attaching clones the append-only trigger back: a record that still
	// lives in the default partition cannot be changed either.
	stranger := addRecordAt(t, client, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	require.Equal(t, "audit_records_default", partitionOf(t, tx, stranger.ID))
	_, err := tx.ExecContext(t.Context(), "DELETE FROM audit_records_default WHERE id = $1", stranger.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "append-only")
}

func TestARecordInAPartitionCannotBeUpdatedOrDeleted(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	require.NoError(t, audit.NewPartitions(tx).EnsureAhead(t.Context(), partitionNow))
	record := addRecordAt(t, client, time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC))

	for name, statement := range map[string]string{
		"update through the parent": "UPDATE audit_records SET outcome = 'failed' WHERE id = $1",
		"delete through the parent": "DELETE FROM audit_records WHERE id = $1",
		"update on the partition":   "UPDATE audit_records_p202410 SET outcome = 'failed' WHERE id = $1",
		"delete from the partition": "DELETE FROM audit_records_p202410 WHERE id = $1",
	} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			_, err := tx.ExecContext(ctx, "SAVEPOINT refusal")
			require.NoError(t, err)
			_, err = tx.ExecContext(ctx, statement, record.ID)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "append-only")
			// Roll back to the savepoint so the aborted statement does not
			// poison the transaction the other cases share.
			_, err = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT refusal")
			require.NoError(t, err)
		})
	}
	assert.Equal(t, "succeeded", client.AuditRecord.GetX(t.Context(), record.ID).Outcome)
}

func TestDropExpiredDropsWholeOldPartitionsAndKeepsTheRest(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	partitions := audit.NewPartitions(tx)
	// May to October 2024 exist, each with one record.
	records := map[time.Month]*ent.AuditRecord{}
	for month := time.May; month <= time.October; month++ {
		require.NoError(t, partitions.EnsureAhead(t.Context(), time.Date(2024, month, 1, 0, 0, 0, 0, time.UTC)))
		records[month] = addRecordAt(t, client, time.Date(2024, month, 15, 0, 0, 0, 0, time.UTC))
	}

	// Retention of one month from 17 October: the cutoff is 17 September.
	dropped, err := partitions.DropExpired(t.Context(), partitionNow, 1)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		"audit_records_p202405", "audit_records_p202406", "audit_records_p202407", "audit_records_p202408",
	}, dropped)
	for _, month := range []time.Month{time.May, time.June, time.July, time.August} {
		_, err = client.AuditRecord.Get(t.Context(), records[month].ID)
		assert.True(t, ent.IsNotFound(err), "the record of %s went with its partition", month)
	}
	// September holds the cutoff: it is not over for retention purposes, so the
	// unfinished month is kept whole, as is the month in progress.
	for _, month := range []time.Month{time.September, time.October} {
		_, err = client.AuditRecord.Get(t.Context(), records[month].ID)
		assert.NoError(t, err, "the record of %s is kept", month)
	}
	assert.Contains(t, partitionNames(t, tx), "audit_records_default")
}

func TestDropExpiredNeverDropsTheMonthInProgressOrTheDefaultPartition(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	partitions := audit.NewPartitions(tx)
	require.NoError(t, partitions.EnsureAhead(t.Context(), partitionNow))
	// An ancient record parked in the default partition is not a month to drop.
	ancient := addRecordAt(t, client, time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC))

	dropped, err := partitions.DropExpired(t.Context(), partitionNow, 1)
	require.NoError(t, err)

	assert.Empty(t, dropped)
	assert.Contains(t, partitionNames(t, tx), "audit_records_p202410")
	assert.Equal(t, "audit_records_default", partitionOf(t, tx, ancient.ID))
}

func TestRetentionBelowOneMonthIsRefused(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	_, err := audit.NewPartitions(tx).DropExpired(t.Context(), partitionNow, 0)
	assert.Error(t, err, "zero months would drop the month in progress")
}

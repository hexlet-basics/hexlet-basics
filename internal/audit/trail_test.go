package audit_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

func TestTheWindowDefaultsToThirtyDaysBackFromNow(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

	from, to, err := audit.Filter{}.Window(now)

	require.NoError(t, err)
	assert.Equal(t, now, to)
	assert.Equal(t, now.Add(-audit.DefaultWindow), from)
}

func TestTheWindowCountsBackFromAGivenEnd(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	end := now.Add(-100 * 24 * time.Hour)

	from, to, err := audit.Filter{To: &end}.Window(now)

	require.NoError(t, err)
	assert.Equal(t, end, to)
	assert.Equal(t, end.Add(-audit.DefaultWindow), from)
}

func TestTheWindowIsRefusedWhenReversedOrTooWide(t *testing.T) {
	now := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

	_, _, err := audit.Filter{From: new(now), To: new(now.Add(-time.Hour))}.Window(now)
	require.ErrorIs(t, err, audit.ErrWindowReversed)

	_, _, err = audit.Filter{From: new(now), To: new(now)}.Window(now)
	require.ErrorIs(t, err, audit.ErrWindowReversed, "an empty window reads nothing")

	_, _, err = audit.Filter{From: new(now.Add(-audit.MaxWindow - time.Second)), To: new(now)}.Window(now)
	require.ErrorIs(t, err, audit.ErrWindowTooWide)

	_, _, err = audit.Filter{From: new(now.Add(-audit.MaxWindow)), To: new(now)}.Window(now)
	assert.NoError(t, err, "exactly the maximum is allowed")
}

func TestAnEntityIdNeedsItsType(t *testing.T) {
	_, err := audit.Filter{EntityID: new("5")}.Query(nil, time.Now())

	require.ErrorIs(t, err, audit.ErrEntityIDWithoutType)
}

// recorder is the sql.Tx the ent client reads through, remembering every query
// it was asked to run so a test can ask PostgreSQL how it would execute them.
type recorder struct {
	tx      *sql.Tx
	queries []recorded
}

type recorded struct {
	query string
	args  []any
}

func (r *recorder) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.tx.ExecContext(ctx, query, args...)
}

func (r *recorder) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	r.queries = append(r.queries, recorded{query: query, args: args})

	return r.tx.QueryContext(ctx, query, args...)
}

// planOf is PostgreSQL's plan for a query the recorder saw, as one string.
func planOf(t *testing.T, tx *sql.Tx, q recorded) string {
	t.Helper()
	rows, err := tx.QueryContext(t.Context(), "EXPLAIN "+q.query, q.args...)
	require.NoError(t, err)

	defer func() { _ = rows.Close() }()

	var plan strings.Builder

	for rows.Next() {
		var line string
		require.NoError(t, rows.Scan(&line))
		plan.WriteString(line + "\n")
	}

	require.NoError(t, rows.Err())

	return plan.String()
}

// The table is partitioned by month, so the listing must carry a bound on
// occurred_at or PostgreSQL visits every partition (and the default one). The
// proof is the plan: with partitions for several months, a filtered query with
// a window inside one month names that month's partition and no other, however
// the other filters are combined.
func TestAFilteredQueryReadsOnlyThePartitionsOfItsWindow(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	partitions := audit.NewPartitions(store.NewTxDriver(tx))
	// Pinned in the past so these partitions never collide with the real ones.
	now := time.Date(2024, 10, 17, 12, 0, 0, 0, time.UTC)
	for _, month := range []time.Time{
		now.AddDate(0, -3, 0), now.AddDate(0, -2, 0), now.AddDate(0, -1, 0), now,
	} {
		require.NoError(t, partitions.EnsureAhead(t.Context(), month))
	}

	rec := &recorder{tx: tx}
	client := ent.NewClient(ent.Driver(entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: rec})))
	window := audit.Filter{
		From: new(time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC)),
		To:   new(time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC)),
	}

	filters := map[string]audit.Filter{
		"window only": window,
		"by actor":    {From: window.From, To: window.To, ActorUserID: new(7)},
		"by entity":   {From: window.From, To: window.To, EntityType: new("BlogPost"), EntityID: new("3")},
		"by channel and operation": {
			From: window.From, To: window.To,
			Channel: new(audit.ChannelMCP), Operation: new("AdminUpdateBlogPost"),
		},
	}
	for name, filter := range filters {
		t.Run(name, func(t *testing.T) {
			rec.queries = nil
			query, err := filter.Query(client, now)
			require.NoError(t, err)
			_, err = query.Limit(20).All(t.Context())
			require.NoError(t, err)
			require.Len(t, rec.queries, 1)

			plan := planOf(t, tx, rec.queries[0])

			assert.Contains(t, plan, "audit_records_p202409", plan)

			for _, other := range []string{"p202407", "p202408", "p202410", "audit_records_default"} {
				assert.NotContains(t, plan, other, plan)
			}
		})
	}
}

// The control: without a bound the same table is read in every partition, which
// is the cost the window exists to prevent.
func TestWithoutABoundEveryPartitionIsRead(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	rec := &recorder{tx: tx}
	unbounded := ent.NewClient(ent.Driver(entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: rec})))

	require.NoError(t, audit.NewPartitions(store.NewTxDriver(tx)).EnsureAhead(t.Context(), time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC)))

	_, err := unbounded.AuditRecord.Query().Limit(20).All(t.Context())
	require.NoError(t, err)

	plan := planOf(t, tx, rec.queries[0])
	assert.Contains(t, plan, "audit_records_p202410")
	assert.Contains(t, plan, "audit_records_default")
}

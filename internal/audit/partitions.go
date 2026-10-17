package audit

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PartitionsAhead is how many months past the current one the maintenance job
// keeps ready. A job that runs daily needs one; three leaves room for a worker
// that is down over a month boundary or two without any record reaching the
// default partition.
const PartitionsAhead = 3

// DefaultRetentionMonths is the retention when nothing configures one: records
// are kept for two years.
const DefaultRetentionMonths = 24

// Queryer is what *sql.DB and *sql.Tx share. Tests hand in a transaction so
// the partitions they create and drop roll back with it.
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// Partitions maintains the monthly partitions of audit_records (ADR-0017,
// ticket #846). It is the only code that may remove Audit Records, and it does
// so by dropping a whole partition; there is no path that updates or deletes a
// single record (the database trigger refuses both).
//
// The DDL is written once, in the SQL functions of the partitioning migration;
// this type only calls them with fixed statements. That is also the answer to
// the "no raw SQL in production Go" rule: ent has no partition DDL, so the
// statements are confined to these two calls, with the logic reviewed as SQL.
type Partitions struct {
	db Queryer
}

// NewPartitions builds the maintainer over a pool (or a transaction in tests).
func NewPartitions(db Queryer) *Partitions { return &Partitions{db: db} }

// EnsureAhead creates the partition of now's month and of the PartitionsAhead
// months after it. It is idempotent: an existing partition is left alone.
// Rows that reached the default partition for one of those months are moved
// into the new partition by the SQL function.
func (p *Partitions) EnsureAhead(ctx context.Context, now time.Time) error {
	month := firstOfMonth(now)
	for i := 0; i <= PartitionsAhead; i++ {
		target := month.AddDate(0, i, 0)
		if _, err := p.call(ctx, "SELECT audit_records_ensure_partition($1)", target); err != nil {
			return fmt.Errorf("ensure audit partition for %s: %w", target.Format("2006-01"), err)
		}
	}
	return nil
}

// DropExpired drops every monthly partition that ended at or before
// now minus retentionMonths and returns the names dropped. The month holding
// the cutoff is kept whole, so a record lives at least retentionMonths and up
// to a month more; partial months are never trimmed (that would be a DELETE).
func (p *Partitions) DropExpired(ctx context.Context, now time.Time, retentionMonths int) ([]string, error) {
	if retentionMonths < 1 {
		return nil, fmt.Errorf("audit retention must be at least one month, got %d", retentionMonths)
	}
	cutoff := now.UTC().AddDate(0, -retentionMonths, 0)
	dropped, err := p.call(ctx, "SELECT audit_records_drop_partitions_ended_before($1)", cutoff)
	if err != nil {
		return nil, fmt.Errorf("drop audit partitions before %s: %w", cutoff.Format(time.DateOnly), err)
	}
	return dropped, nil
}

func firstOfMonth(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// call runs one of the maintenance functions and collects the text it returns
// (the partition names created or dropped).
func (p *Partitions) call(ctx context.Context, statement string, month time.Time) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, statement, month) //nolint:forbidigo // fixed calls to the migration's partition functions; ent has no partition DDL
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

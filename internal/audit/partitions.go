package audit

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

// PartitionsAhead is how many months past the current one the maintenance job
// keeps ready. A job that runs daily needs one; three leaves room for a worker
// that is down over a month boundary or two without any record reaching the
// default partition.
const PartitionsAhead = 3

// DefaultRetentionMonths is the retention when nothing configures one: records
// are kept for two years.
const DefaultRetentionMonths = 24

// Driver is the part of ent's dialect.Driver the maintenance needs: a query run
// through the layer ent itself uses for every statement. The store builds it
// (store.NewDriver over the pool, store.NewTxDriver over a transaction a test
// owns and rolls back).
type Driver interface {
	Query(ctx context.Context, query string, args, v any) error
}

// Partitions maintains the monthly partitions of audit_records (ADR-0020,
// ticket #846). It is the only code that may remove Audit Records, and it does
// so by dropping a whole partition; there is no path that updates or deletes a
// single record.
//
// The DDL is written once, in the SQL functions of the partitioning migration;
// this type only calls them. That is also the answer to the "no raw SQL in
// production Go" rule: ent has no partition DDL, so the call is built with ent's
// own SQL builder (a select of the function, the month as a bound argument) and
// run on ent's driver, never through database/sql.
type Partitions struct {
	db Driver
}

// NewPartitions builds the maintainer over ent's driver.
func NewPartitions(db Driver) *Partitions { return &Partitions{db: db} }

// EnsureAhead creates the partition of now's month and of the PartitionsAhead
// months after it. It is idempotent: an existing partition is left alone.
// Rows that reached the default partition for one of those months are moved
// into the new partition by the SQL function.
func (p *Partitions) EnsureAhead(ctx context.Context, now time.Time) error {
	month := firstOfMonth(now)
	for i := 0; i <= PartitionsAhead; i++ {
		target := month.AddDate(0, i, 0)

		_, err := p.call(ctx, "audit_records_ensure_partition", target)
		if err != nil {
			return errs.With("month", target.Format("2006-01")).Wrapf(err, "ensure audit partition")
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
		return nil, errs.With("retention_months", retentionMonths).Errorf("audit retention must be at least one month")
	}

	cutoff := now.UTC().AddDate(0, -retentionMonths, 0)

	dropped, err := p.call(ctx, "audit_records_drop_partitions_ended_before", cutoff)
	if err != nil {
		return nil, errs.With("cutoff", cutoff.Format(time.DateOnly)).Wrapf(err, "drop audit partitions")
	}

	return dropped, nil
}

func firstOfMonth(t time.Time) time.Time {
	u := t.UTC()

	return time.Date(u.Year(), u.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// call runs one of the two maintenance functions, named by a constant of this
// file and given the month as a bound argument, and collects the text it
// returns (the partition names created or dropped).
func (p *Partitions) call(ctx context.Context, function string, month time.Time) (_ []string, err error) {
	query, args := entsql.Dialect(dialect.Postgres).
		Select().
		AppendSelectExprAs(entsql.ExprFunc(func(b *entsql.Builder) {
			b.WriteString(function).WriteByte('(').Arg(month).WriteByte(')')
		}), "name").
		Query()

	rows := &entsql.Rows{}

	err = p.db.Query(ctx, query, args, rows)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = errors.Join(err, errs.Wrapf(closeErr, "close partition names"))
		}
	}()

	var names []string

	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			return nil, errs.Wrapf(err, "scan partition name")
		}

		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		return nil, errs.Wrapf(err, "read partition names")
	}

	return names, nil
}

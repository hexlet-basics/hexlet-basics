package jobs

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/internal/store"
)

// Enqueuer is how the HTTP process schedules jobs. It wraps the insert-only
// River client and owns the one rule every caller needs: jobs are inserted
// through a SQL transaction (so they commit or roll back with whatever the
// caller wrote, and a test's rollback-only transaction sees them) and a batch
// goes in all together or not at all. Each kind's queue, attempts and
// uniqueness come from its Args' InsertOpts, never from the caller; a
// duplicate skipped by uniqueness is not an error.
type Enqueuer struct {
	store store.Transactor
	river *river.Client[*sql.Tx]
}

// NewEnqueuer wires the shared pool and the insert-only River client.
func NewEnqueuer(txStore store.Transactor, riverClient *river.Client[*sql.Tx]) *Enqueuer {
	return &Enqueuer{store: txStore, river: riverClient}
}

// Enqueue inserts args in one transaction of its own.
func (e *Enqueuer) Enqueue(ctx context.Context, args ...river.JobArgs) error {
	if len(args) == 0 {
		return nil
	}
	return e.store.WithinTx(ctx, func(tx *sql.Tx, _ *ent.Client) error {
		return e.EnqueueTx(ctx, tx, args...)
	})
}

// EnqueueTx inserts args inside the caller's transaction, for a business write
// that must not exist without its job (a course version and its build).
func (e *Enqueuer) EnqueueTx(ctx context.Context, tx *sql.Tx, args ...river.JobArgs) error {
	if len(args) == 0 {
		return nil
	}
	params := lo.Map(args, func(a river.JobArgs, _ int) river.InsertManyParams {
		return river.InsertManyParams{Args: a}
	})
	if _, err := e.river.InsertManyTx(ctx, tx, params); err != nil {
		return oops.Wrapf(err, "enqueue %s", args[0].Kind())
	}
	return nil
}

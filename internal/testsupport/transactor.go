package testsupport

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"

	"hexletbasics/ent"
)

// savepointTransactor adapts the outer transaction owned by a test to the same
// transaction interface production modules use. Savepoints preserve callback
// rollback semantics without committing the fixture-isolation transaction.
type savepointTransactor struct {
	tx     *sql.Tx
	client *ent.Client
	next   atomic.Uint64
}

func newSavepointTransactor(tx *sql.Tx, client *ent.Client) *savepointTransactor {
	return &savepointTransactor{tx: tx, client: client}
}

func (s *savepointTransactor) WithinTx(
	ctx context.Context,
	fn func(*sql.Tx, *ent.Client) error,
) (err error) {
	name := fmt.Sprintf("test_within_tx_%d", s.next.Add(1))
	if _, err := s.tx.ExecContext(ctx, "SAVEPOINT "+name); err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	finished := false
	defer func() {
		if finished {
			return
		}
		rollbackCtx := context.WithoutCancel(ctx)
		_, _ = s.tx.ExecContext(rollbackCtx, "ROLLBACK TO SAVEPOINT "+name)
		_, _ = s.tx.ExecContext(rollbackCtx, "RELEASE SAVEPOINT "+name)
	}()

	if err := fn(s.tx, s.client); err != nil {
		rollbackCtx := context.WithoutCancel(ctx)
		if _, rollbackErr := s.tx.ExecContext(rollbackCtx, "ROLLBACK TO SAVEPOINT "+name); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
		}
		if _, releaseErr := s.tx.ExecContext(rollbackCtx, "RELEASE SAVEPOINT "+name); releaseErr != nil {
			return errors.Join(err, fmt.Errorf("rollback transaction: %w", releaseErr))
		}
		finished = true
		return err
	}
	if _, err := s.tx.ExecContext(ctx, "RELEASE SAVEPOINT "+name); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	finished = true
	return nil
}

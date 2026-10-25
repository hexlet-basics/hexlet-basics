package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errCallbackBusinessWriteFailed = errors.New("business write failed")
	errBegin                       = errors.New("begin failed")
	errCommit                      = errors.New("commit failed")
	errCallback                    = errors.New("callback failed")
	errRollback                    = errors.New("rollback failed")
	errUseConnector                = errors.New("use connector")
	errPrepareNotSupported         = errors.New("prepare not supported")
)

func TestWithinTxCommitsSuccessfulCallback(t *testing.T) {
	state := &transactionState{
		mu:          sync.Mutex{},
		begins:      0,
		execs:       0,
		commits:     0,
		rollbacks:   0,
		beginErr:    nil,
		commitErr:   nil,
		rollbackErr: nil,
	}
	txStore := New(openTransactionDB(t, state))

	err := txStore.WithinTx(t.Context(), func(tx *sql.Tx, txClient *ent.Client) error {
		require.NotNil(t, txClient)
		_, execErr := tx.ExecContext(t.Context(), "write")

		return errs.Wrapf(execErr, "exec in transaction")
	})

	require.NoError(t, err)
	assert.Equal(t, 1, state.begins)
	assert.Equal(t, 1, state.execs)
	assert.Equal(t, 1, state.commits)
	assert.Zero(t, state.rollbacks)
}

func TestWithinTxRollsBackAndPreservesCallbackError(t *testing.T) {
	state := &transactionState{
		mu:          sync.Mutex{},
		begins:      0,
		execs:       0,
		commits:     0,
		rollbacks:   0,
		beginErr:    nil,
		commitErr:   nil,
		rollbackErr: nil,
	}
	txStore := New(openTransactionDB(t, state))

	err := txStore.WithinTx(t.Context(), func(*sql.Tx, *ent.Client) error {
		return errCallbackBusinessWriteFailed
	})

	require.ErrorIs(t, err, errCallbackBusinessWriteFailed)
	assert.Zero(t, state.commits)
	assert.Equal(t, 1, state.rollbacks)
}

func TestWithinTxRollsBackPanic(t *testing.T) {
	state := &transactionState{
		mu:          sync.Mutex{},
		begins:      0,
		execs:       0,
		commits:     0,
		rollbacks:   0,
		beginErr:    nil,
		commitErr:   nil,
		rollbackErr: nil,
	}
	txStore := New(openTransactionDB(t, state))

	assert.PanicsWithValue(t, "boom", func() {
		assert.NoError(t, txStore.WithinTx(t.Context(), func(*sql.Tx, *ent.Client) error {
			panic("boom")
		}))
	})
	assert.Zero(t, state.commits)
	assert.Equal(t, 1, state.rollbacks)
}

func TestWithinTxNormalizesLifecycleErrors(t *testing.T) {
	t.Run("begin", func(t *testing.T) {
		state := &transactionState{
			beginErr:    errBegin,
			mu:          sync.Mutex{},
			begins:      0,
			execs:       0,
			commits:     0,
			rollbacks:   0,
			commitErr:   nil,
			rollbackErr: nil,
		}

		err := New(openTransactionDB(t, state)).WithinTx(
			t.Context(),
			func(*sql.Tx, *ent.Client) error { return nil },
		)

		require.ErrorIs(t, err, errBegin)
		assert.ErrorContains(t, err, "begin transaction")
	})

	t.Run("commit", func(t *testing.T) {
		state := &transactionState{
			commitErr:   errCommit,
			mu:          sync.Mutex{},
			begins:      0,
			execs:       0,
			commits:     0,
			rollbacks:   0,
			beginErr:    nil,
			rollbackErr: nil,
		}

		err := New(openTransactionDB(t, state)).WithinTx(
			t.Context(),
			func(*sql.Tx, *ent.Client) error { return nil },
		)

		require.ErrorIs(t, err, errCommit)
		assert.ErrorContains(t, err, "commit transaction")
	})

	t.Run("rollback", func(t *testing.T) {
		state := &transactionState{
			rollbackErr: errRollback,
			mu:          sync.Mutex{},
			begins:      0,
			execs:       0,
			commits:     0,
			rollbacks:   0,
			beginErr:    nil,
			commitErr:   nil,
		}

		err := New(openTransactionDB(t, state)).WithinTx(
			t.Context(),
			func(*sql.Tx, *ent.Client) error { return errCallback },
		)

		require.ErrorIs(t, err, errCallback)
		require.ErrorIs(t, err, errRollback)
		assert.ErrorContains(t, err, "rollback transaction")
	})
}

// The ent client handed to the callback is bound to the caller's transaction, so
// what the callback does with it must not end that transaction or close the pool:
// a nested Tx joins the outer one and Close leaves the connection alone.
func TestWithinTxClientNeitherOwnsTheTransactionNorClosesThePool(t *testing.T) {
	state := &transactionState{
		mu:          sync.Mutex{},
		begins:      0,
		execs:       0,
		commits:     0,
		rollbacks:   0,
		beginErr:    nil,
		commitErr:   nil,
		rollbackErr: nil,
	}
	db := openTransactionDB(t, state)

	err := New(db).WithinTx(t.Context(), func(tx *sql.Tx, txClient *ent.Client) error {
		nested, err := txClient.Tx(t.Context())
		require.NoError(t, err)
		require.NoError(t, nested.Commit())
		require.NoError(t, txClient.Close())

		_, execErr := tx.ExecContext(t.Context(), "write")

		return errs.Wrapf(execErr, "exec in transaction")
	})

	require.NoError(t, err)
	assert.Equal(t, 1, state.begins, "the nested Tx joined the outer transaction")
	assert.Equal(t, 1, state.commits, "only the outer transaction committed")
	assert.Equal(t, 1, state.execs)
	require.NoError(t, db.PingContext(t.Context()), "the pool is still open")
}

type transactionState struct {
	mu          sync.Mutex
	begins      int
	execs       int
	commits     int
	rollbacks   int
	beginErr    error
	commitErr   error
	rollbackErr error
}

func openTransactionDB(t *testing.T, state *transactionState) *sql.DB {
	t.Helper()

	db := sql.OpenDB(transactionConnector{state: state})

	t.Cleanup(func() { require.NoError(t, db.Close()) })

	return db
}

type transactionConnector struct {
	state *transactionState
}

func (c transactionConnector) Connect(context.Context) (driver.Conn, error) {
	return &transactionConn{state: c.state}, nil
}

func (transactionConnector) Driver() driver.Driver {
	return transactionDriver{}
}

type transactionDriver struct{}

func (transactionDriver) Open(string) (driver.Conn, error) {
	return nil, errUseConnector
}

type transactionConn struct {
	state *transactionState
}

func (*transactionConn) Prepare(string) (driver.Stmt, error) {
	return nil, errPrepareNotSupported
}

func (*transactionConn) Close() error { return nil }

func (c *transactionConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *transactionConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.state.mu.Lock()
	defer c.state.mu.Unlock()

	if c.state.beginErr != nil {
		return nil, c.state.beginErr
	}

	c.state.begins++

	return transactionTx{state: c.state}, nil
}

func (c *transactionConn) ExecContext(
	context.Context,
	string,
	[]driver.NamedValue,
) (driver.Result, error) {
	c.state.mu.Lock()
	defer c.state.mu.Unlock()

	c.state.execs++

	return driver.RowsAffected(1), nil
}

type transactionTx struct {
	state *transactionState
}

func (tx transactionTx) Commit() error {
	tx.state.mu.Lock()
	defer tx.state.mu.Unlock()

	tx.state.commits++

	return tx.state.commitErr
}

func (tx transactionTx) Rollback() error {
	tx.state.mu.Lock()
	defer tx.state.mu.Unlock()

	tx.state.rollbacks++

	return tx.state.rollbackErr
}

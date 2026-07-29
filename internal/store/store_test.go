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

func TestWithinTxCommitsSuccessfulCallback(t *testing.T) {
	state := &transactionState{}
	txStore := New(openTransactionDB(t, state))

	err := txStore.WithinTx(t.Context(), func(tx *sql.Tx, txClient *ent.Client) error {
		require.NotNil(t, txClient)
		_, execErr := tx.ExecContext(t.Context(), "write")
		return execErr
	})

	require.NoError(t, err)
	assert.Equal(t, 1, state.begins)
	assert.Equal(t, 1, state.execs)
	assert.Equal(t, 1, state.commits)
	assert.Zero(t, state.rollbacks)
}

func TestWithinTxRollsBackAndPreservesCallbackError(t *testing.T) {
	callbackErr := errors.New("business write failed")
	state := &transactionState{}
	txStore := New(openTransactionDB(t, state))

	err := txStore.WithinTx(t.Context(), func(*sql.Tx, *ent.Client) error {
		return callbackErr
	})

	require.ErrorIs(t, err, callbackErr)
	assert.Zero(t, state.commits)
	assert.Equal(t, 1, state.rollbacks)
}

func TestWithinTxRollsBackPanic(t *testing.T) {
	state := &transactionState{}
	txStore := New(openTransactionDB(t, state))

	assert.PanicsWithValue(t, "boom", func() {
		_ = txStore.WithinTx(t.Context(), func(*sql.Tx, *ent.Client) error {
			panic("boom")
		})
	})
	assert.Zero(t, state.commits)
	assert.Equal(t, 1, state.rollbacks)
}

func TestWithinTxNormalizesLifecycleErrors(t *testing.T) {
	t.Run("begin", func(t *testing.T) {
		beginErr := errors.New("begin failed")
		state := &transactionState{beginErr: beginErr}

		err := New(openTransactionDB(t, state)).WithinTx(
			t.Context(),
			func(*sql.Tx, *ent.Client) error { return nil },
		)

		require.ErrorIs(t, err, beginErr)
		assert.ErrorContains(t, err, "begin transaction")
	})

	t.Run("commit", func(t *testing.T) {
		commitErr := errors.New("commit failed")
		state := &transactionState{commitErr: commitErr}

		err := New(openTransactionDB(t, state)).WithinTx(
			t.Context(),
			func(*sql.Tx, *ent.Client) error { return nil },
		)

		require.ErrorIs(t, err, commitErr)
		assert.ErrorContains(t, err, "commit transaction")
	})

	t.Run("rollback", func(t *testing.T) {
		callbackErr := errors.New("callback failed")
		rollbackErr := errors.New("rollback failed")
		state := &transactionState{rollbackErr: rollbackErr}

		err := New(openTransactionDB(t, state)).WithinTx(
			t.Context(),
			func(*sql.Tx, *ent.Client) error { return callbackErr },
		)

		require.ErrorIs(t, err, callbackErr)
		require.ErrorIs(t, err, rollbackErr)
		assert.ErrorContains(t, err, "rollback transaction")
	})
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
	t.Cleanup(func() { _ = db.Close() })
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
	return nil, errors.New("use connector")
}

type transactionConn struct {
	state *transactionState
}

func (*transactionConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
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

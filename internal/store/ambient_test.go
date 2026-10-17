package store_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/auditrecord"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// These run on a real pool (not the test's shared transaction) because what is
// under test is which connection a statement uses. Every write is rolled back,
// so nothing outlives the test.

func openPool(t *testing.T) (*store.Store, *ent.Client) {
	t.Helper()
	db, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return store.New(db), store.NewClient(db)
}

func countOf(ctx context.Context, t *testing.T, client *ent.Client, operation string) int {
	t.Helper()
	n, err := client.AuditRecord.Query().Where(auditrecord.Operation(operation)).Count(ctx)
	require.NoError(t, err)
	return n
}

func TestAClientStatementJoinsTheAmbientTransaction(t *testing.T) {
	txStore, client := openPool(t)
	const operation = "AmbientJoinProbe"
	boom := errors.New("operation failed")

	err := txStore.WithinTx(t.Context(), func(tx *sql.Tx, txClient *ent.Client) error {
		ctx := store.WithAmbientTx(t.Context(), tx, txClient)

		// Written through the pooled client, yet inside the transaction: the
		// handler's plain `s.db` takes part without knowing.
		require.NoError(t, client.AuditRecord.Create().
			SetChannel("web").SetOperation(operation).SetOutcome("succeeded").Exec(ctx))
		// A read in the same operation sees the write ...
		assert.Equal(t, 1, countOf(ctx, t, client, operation))
		// ... while another connection does not, because nothing committed.
		assert.Zero(t, countOf(t.Context(), t, client, operation))
		return boom
	})

	require.ErrorIs(t, err, boom)
	assert.Zero(t, countOf(t.Context(), t, client, operation), "rolled back with the operation")
}

func TestWithinTxJoinsTheAmbientTransactionInsteadOfOpeningAnother(t *testing.T) {
	txStore, client := openPool(t)
	const operation = "AmbientNestedProbe"
	boom := errors.New("operation failed")

	err := txStore.WithinTx(t.Context(), func(tx *sql.Tx, txClient *ent.Client) error {
		ctx := store.WithAmbientTx(t.Context(), tx, txClient)

		// A module that owns an atomic step commits it on its own when it is
		// called alone; called inside an operation it must not.
		require.NoError(t, txStore.WithinTx(ctx, func(_ *sql.Tx, inner *ent.Client) error {
			return inner.AuditRecord.Create().
				SetChannel("web").SetOperation(operation).SetOutcome("succeeded").Exec(ctx)
		}))
		return boom
	})

	require.ErrorIs(t, err, boom)
	assert.Zero(t, countOf(t.Context(), t, client, operation), "the module's step is undone with the operation")
}

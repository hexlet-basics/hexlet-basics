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
	"hexletbasics/internal/audit"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// These run on a real pool (not the test's shared transaction) because what is
// under test is which connection a statement uses. Every write is rolled back,
// so nothing outlives the test.

func openPool(t *testing.T) (*store.Store, *ent.Client) {
	t.Helper()
	db := testdb.Clone(t)
	return store.New(db), store.NewClient(db)
}

// probe writes an Audit Record, through the production writer, as the row whose
// visibility the cases look at: a small insert the ambient transaction has to
// capture without the caller knowing.
func probe(ctx context.Context, client *ent.Client, operation string) error {
	return audit.Write(ctx, client, audit.Entry{
		Source:    audit.Source{Channel: audit.ChannelWeb},
		Operation: operation,
		Outcome:   audit.OutcomeSucceeded,
	}, nil)
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
		require.NoError(t, probe(ctx, client, operation))
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
			return probe(ctx, inner, operation)
		}))
		return boom
	})

	require.ErrorIs(t, err, boom)
	assert.Zero(t, countOf(t.Context(), t, client, operation), "the module's step is undone with the operation")
}

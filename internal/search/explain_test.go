package search_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// explain is the plan PostgreSQL makes for a statement the search sent.
func explain(t *testing.T, tx *sql.Tx, q capturedQuery) string {
	t.Helper()
	rows, err := tx.QueryContext(t.Context(), "EXPLAIN "+q.query, q.args...)
	require.NoError(t, err, q.query)
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

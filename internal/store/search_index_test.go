package store_test

import (
	"strings"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport"
)

// The admin lists filter on `q` with ent's ContainsFold, which PostgreSQL gets
// as ILIKE '%term%'. Only a trigram index can serve that, so the tables that
// grow without bound have one per searched column (the migration says which
// and why). A plan with a sequential scan would still be correct, only slow at
// production size, so this reads the plan: the planner is told a sequential
// scan is as dear as it can be, because on the handful of fixture rows it would
// pick one regardless and say nothing about the index.
func TestTheSearchOfEachLargeTableIsServedByATrigramIndex(t *testing.T) {
	cases := []struct{ table, column string }{
		{"users", "first_name"},
		{"users", "last_name"},
		{"users", "email"},
		{"leads", "email"},
		{"leads", "phone"},
		{"leads", "whatsapp"},
		{"leads", "telegram"},
		{"ai_messages", "content"},
		{"language_lesson_version_infos", "name"},
		{"language_lesson_version_infos", "description"},
	}
	_, tx := testsupport.NewClientWithTx(t)
	_, err := tx.ExecContext(t.Context(), `SET LOCAL enable_seqscan = off`)
	require.NoError(t, err)

	for _, c := range cases {
		t.Run(c.table+"."+c.column, func(t *testing.T) {
			// The predicate is the one ent's generated `<Field>ContainsFold` builds.
			query, args := sql.Dialect(dialect.Postgres).
				Select(sql.Count("*")).
				From(sql.Table(c.table)).
				Where(sql.ContainsFold(c.column, "term")).
				Query()

			rows, err := tx.QueryContext(t.Context(), "EXPLAIN "+query, args...)
			require.NoError(t, err)

			defer func() { _ = rows.Close() }()

			var plan []string

			for rows.Next() {
				var line string
				require.NoError(t, rows.Scan(&line))
				plan = append(plan, line)
			}

			require.NoError(t, rows.Err())

			// users.email also has the unique btree, which the planner may scan
			// whole for want of anything better, so the plan alone cannot say the
			// trigram index is the one that serves it: the index has to exist.
			var trigram int
			require.NoError(t, tx.QueryRowContext(t.Context(),
				`SELECT count(*) FROM pg_indexes
				 WHERE tablename = $1 AND indexdef LIKE '%USING gin (' || $2 || ' gin_trgm_ops)'`,
				c.table, c.column).Scan(&trigram))
			assert.Equal(t, 1, trigram, "a trigram index on %s.%s", c.table, c.column)

			joined := strings.Join(plan, "\n")
			assert.NotContains(t, joined, "Seq Scan", joined)
		})
	}
}

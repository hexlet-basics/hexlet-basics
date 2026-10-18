package blog_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

const likeIndex = "index_blog_post_likes_on_blog_post_id_and_user_id"

// The unique-like migration must keep a person's earliest like when legacy left
// several for one post, and leave likes without a user alone. The migration file
// runs verbatim, so what is verified is what ships; the test database already
// carries the index, so it is dropped first to arrange legacy's duplicates.
func TestUniqueLikeMigrationKeepsTheEarliestLikePerPair(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	ctx := t.Context()
	post := fixture.ID(fixture.BlogPostSecondPost)
	author := fixture.ID(fixture.UserAdminLearner)
	other := fixture.ID(fixture.UserLearnerHiddenCards)

	// Raw SQL is the arrange here: ent cannot write the duplicates the index
	// forbids, and the legacy rows are the case under test.
	_, err := tx.ExecContext(ctx, `DROP INDEX "`+likeIndex+`"`)
	require.NoError(t, err)
	insert := func(userID *int, createdAt string) {
		t.Helper()
		_, err := tx.ExecContext(ctx,
			`INSERT INTO "blog_post_likes" ("blog_post_id", "user_id", "created_at", "updated_at") VALUES ($1, $2, $3, $3)`,
			post, userID, createdAt)
		require.NoError(t, err)
	}
	insert(&author, "2026-01-02 10:00:00")
	insert(&author, "2026-01-01 10:00:00") // the earliest, inserted second
	insert(&author, "2026-01-03 10:00:00")
	insert(&other, "2026-01-01 09:00:00")
	insert(nil, "2026-01-01 09:00:00")
	insert(nil, "2026-01-02 09:00:00")

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	statements, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations", "20261017210100_unique_blog_post_like.sql"))
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, string(statements))
	require.NoError(t, err)

	created := func(userID *int) []string {
		t.Helper()
		var user sql.NullInt64
		if userID != nil {
			user = sql.NullInt64{Int64: int64(*userID), Valid: true}
		}
		rows, err := tx.QueryContext(ctx,
			`SELECT to_char("created_at", 'YYYY-MM-DD HH24:MI:SS') FROM "blog_post_likes"
			 WHERE "blog_post_id" = $1 AND "user_id" IS NOT DISTINCT FROM $2 ORDER BY "created_at"`, post, user)
		require.NoError(t, err)
		t.Cleanup(func() { _ = rows.Close() })
		var got []string
		for rows.Next() {
			var at string
			require.NoError(t, rows.Scan(&at))
			got = append(got, at)
		}
		require.NoError(t, rows.Err())
		return got
	}
	assert.Equal(t, []string{"2026-01-01 10:00:00"}, created(&author), "only the earliest like of the pair stays")
	assert.Equal(t, []string{"2026-01-01 09:00:00"}, created(&other), "a single like is untouched")
	assert.Len(t, created(nil), 2, "likes without a user are not pairs")

	var indexes int
	require.NoError(t, tx.QueryRowContext(ctx,
		`SELECT count(*) FROM pg_indexes WHERE indexname = $1`, likeIndex).Scan(&indexes))
	assert.Equal(t, 1, indexes, "the unique index exists")
}

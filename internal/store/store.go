// Package store opens the ent client backed by pgx.
package store

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"hexletbasics/ent"
)

// NewPool opens a pgx connection pool for river, which needs a native
// *pgxpool.Pool (riverpgxv5) rather than the database/sql handle ent uses. Both
// point at the same database; the pool is a small, separate connection set for
// the job queue.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
}

// NewClient opens an ent client over pgx (database/sql). It never runs
// migrations: the schema is owned by the legacy Rails database.
func NewClient(dsn string) (*ent.Client, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

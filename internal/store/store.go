// Package store opens the ent client backed by pgx.
package store

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gocloud.dev/blob"
	// Blank imports register the gocloud blob backends selected by URL scheme
	// (ADR-0005): fileblob (`file://…`) for dev, s3blob (`s3://…`) for prod.
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"

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

// NewBucket opens the gocloud blob bucket that stores uploaded assets (ADR-0005).
// The URL scheme selects the backend (registered via the blank imports above):
// `file://…` writes to local disk in dev, `s3://…` talks to S3 in prod. The
// caller owns Close (main.go drains it on shutdown, like the pool and ent client).
func NewBucket(ctx context.Context, bucketURL string) (*blob.Bucket, error) {
	return blob.OpenBucket(ctx, bucketURL)
}

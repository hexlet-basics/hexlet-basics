// Package store opens the ent client backed by pgx.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gocloud.dev/blob"
	// Blank imports register the gocloud blob backends selected by URL scheme
	// (ADR-0005): fileblob (`file://…`) for dev, s3blob (`s3://…`) for prod.
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"

	"hexletbasics/ent"
)

// Transactor is the transaction seam shared by business modules. The callback
// receives both views of the same PostgreSQL transaction: database/sql for
// River and Watermill adapters, and ent for application persistence.
type Transactor interface {
	WithinTx(context.Context, func(*sql.Tx, *ent.Client) error) error
}

// Store owns the application's transaction lifecycle over the shared SQL pool.
type Store struct {
	db *sql.DB
}

// New builds the transaction store over the application's shared SQL pool.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// NewDB opens the application's single database/sql pool over pgx. Both ent and
// River use it so a business write and its background job can share one sql.Tx.
func NewDB(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn)
}

// NewClient builds an ent client over the shared database/sql pool. It never
// runs migrations: atlas owns the schema. Closing the client closes the pool.
func NewClient(db *sql.DB) *ent.Client {
	return ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
}

// NewTxClient binds ent to an existing sql.Tx owned by the caller. The standard
// ent SQL driver already implements all query and scan behavior; NopTx only
// prevents ent from trying to own the caller's commit/rollback lifecycle.
func NewTxClient(tx *sql.Tx) *ent.Client {
	driver := entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: tx})
	return ent.NewClient(ent.Driver(txBoundDriver{Driver: driver}))
}

// WithinTx runs fn with SQL and ent adapters bound to one transaction. Callback
// errors keep their identity; lifecycle errors are normalized here so callers
// only need to add business context.
func (s *Store) WithinTx(ctx context.Context, fn func(*sql.Tx, *ent.Client) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	// This also protects panic paths. Explicit rollback below is what lets us
	// report rollback failures without obscuring the callback error.
	defer func() { _ = tx.Rollback() }()

	if err := fn(tx, NewTxClient(tx)); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr == nil || errors.Is(rollbackErr, sql.ErrTxDone) {
			return err
		}
		return errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

type txBoundDriver struct{ dialect.Driver }

func (d txBoundDriver) Tx(context.Context) (dialect.Tx, error) {
	return dialect.NopTx(d), nil
}

func (txBoundDriver) Close() error { return nil }

// NewBucket opens the gocloud blob bucket that stores uploaded assets (ADR-0005).
// The URL scheme selects the backend (registered via the blank imports above):
// `file://…` writes to local disk in dev, `s3://…` talks to S3 in prod. The
// caller owns Close (main.go drains it on shutdown, like the pool and ent client).
func NewBucket(ctx context.Context, bucketURL string) (*blob.Bucket, error) {
	return blob.OpenBucket(ctx, bucketURL)
}

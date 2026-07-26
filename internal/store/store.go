// Package store opens the ent client backed by pgx.
package store

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"hexletbasics/ent"
)

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

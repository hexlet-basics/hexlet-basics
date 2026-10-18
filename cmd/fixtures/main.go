// Command fixtures loads the committed fixtures/ into a database, replacing
// the data of the fixture tables. It is `mise run db:fixtures` and the
// end-to-end database setup; the Go tests load the same files through the same
// package, so every consumer sees the same starting data.
package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"hexletbasics/internal/fixture"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fixtures:", err)
		os.Exit(1)
	}
}

// errConnRequired is what run answers when no connection string is given.
var errConnRequired = errors.New("-conn is required")

func run() error {
	conn := flag.String("conn", "", "PostgreSQL connection string")
	dir := flag.String("dir", "fixtures", "directory of fixture files")
	anyDatabase := flag.Bool("any-database", false, `load into a database whose name does not contain "test" (the development database)`)

	flag.Parse()

	if *conn == "" {
		return errConnRequired
	}

	db, err := sql.Open("pgx", *conn)
	if err != nil {
		return errs.Wrapf(err, "open database")
	}
	defer func() { _ = db.Close() }()

	return fixture.Load(db, *dir, *anyDatabase)
}

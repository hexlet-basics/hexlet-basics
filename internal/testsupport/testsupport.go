// Package testsupport provides Rails-style integration-test plumbing: each test
// runs against a real Postgres transaction that is rolled back on cleanup
// (go-txdb), over a baseline prepared once by `make test-prepare` — the schema
// snapshot (db/structure.sql) plus the fixtures/ YAML loaded by the testfixtures
// CLI. Handlers hit a real database; every test's writes are discarded on
// rollback, so the shared baseline is never mutated and nothing is left behind.
package testsupport

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	txdb "github.com/DATA-DOG/go-txdb"
	_ "github.com/jackc/pgx/v5/stdlib"

	"hexletbasics/ent"
)

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

var registerOnce sync.Once

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return defaultTestDSN
}

// NewClient opens an ent client bound to a fresh go-txdb transaction that is
// rolled back when the test finishes. Fixtures are not loaded here: they are the
// pre-loaded baseline from `make test-prepare`, which every test sees and whose
// own writes vanish on rollback. Run `make test-prepare` before `go test`.
func NewClient(t *testing.T) *ent.Client {
	t.Helper()

	registerOnce.Do(func() {
		txdb.Register("txdb", "pgx", testDSN())
	})

	db, err := sql.Open("txdb", t.Name())
	if err != nil {
		t.Fatalf("open txdb: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() }) // rolls the transaction back

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

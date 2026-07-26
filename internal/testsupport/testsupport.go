// Package testsupport provides Rails-style integration-test plumbing: each test
// runs against a real Postgres transaction that is rolled back on cleanup
// (go-txdb), with the starting data loaded from YAML fixtures (testfixtures).
// Handlers are exercised against a real database; nothing is left behind.
package testsupport

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	txdb "github.com/DATA-DOG/go-txdb"
	testfixtures "github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/v5/stdlib"

	"hexletbasics/ent"
)

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

var (
	registerOnce sync.Once
	fixturesDir  string
)

func init() {
	// Fixtures live in the repo-root `fixtures/` dir; resolve it from this
	// file's location so tests work regardless of the working directory.
	_, thisFile, _, _ := runtime.Caller(0)
	fixturesDir = filepath.Join(filepath.Dir(thisFile), "..", "..", "fixtures")
}

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return defaultTestDSN
}

// NewClient opens an ent client bound to a fresh transaction, loads the
// fixtures into it, and rolls the transaction back when the test finishes.
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

	loadFixtures(t, db)

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func loadFixtures(t *testing.T, db *sql.DB) {
	t.Helper()

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgresql"),
		testfixtures.Directory(fixturesDir),
	)
	if err != nil {
		t.Fatalf("build fixtures: %v", err)
	}
	if err := fixtures.Load(); err != nil {
		t.Fatalf("load fixtures: %v", err)
	}
}

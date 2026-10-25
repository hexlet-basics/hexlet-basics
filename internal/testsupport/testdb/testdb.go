// Package testdb owns the lifecycle of the disposable database a test package
// uses. A package gets its own database cloned from a migrated, fixture-loaded
// template on the development Postgres daemon (`mise run db:up`), the way Rails
// clones `db_test-N` for parallel tests. The template is built once per
// migrations/fixtures content hash and shared by every package, so a clone costs
// milliseconds where a container per package cost seconds (ADR-0019).
package testdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"ariga.io/atlas/sql/migrate"
	atlaspostgres "ariga.io/atlas/sql/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mokevnin/dbcleaner"
	_ "github.com/mokevnin/dbcleaner/postgres"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/common"

	"hexletbasics/internal/fixture"
)

// The server tests clone from. The defaults are the db daemon of mise.toml; CI
// points TEST_DB_HOST/TEST_DB_PORT at its service container, as E2E_DB_* does
// for the end-to-end suite. The credentials are the daemon's superuser, which
// pgtestdb uses only to create the template, the clones and its own role.
const (
	defaultHost = "127.0.0.1"
	defaultPort = "54330"
	adminUser   = "postgres"
	adminPass   = "postgres"
)

var (
	databaseURL string
	// template is what Clone clones from: the server and the migrator Run used.
	template struct {
		conf pgtestdb.Config
		m    migrator
	}
	// cleaner hands each test the package's database inside a transaction that
	// is rolled back when the test ends. Its setup is empty: the clone already
	// holds the migrations and fixtures.
	cleaner *dbcleaner.Cleaner
	urlMu   sync.RWMutex
)

// Run clones the package's database before running its tests. TestMain must
// pass the returned status to os.Exit so the database is dropped before the
// process exits.
func Run(m *testing.M) (code int) {
	root, err := repositoryRoot()
	if err != nil {
		return setupFailed(err)
	}

	conf := serverConfig()
	if err := requireServer(conf); err != nil {
		return setupFailed(err)
	}

	// pgtestdb takes a testing.TB, and TestMain has none: this one records the
	// failure for Run to report and holds the cleanup until the tests are done.
	setup := &mainTB{}
	defer func() {
		setup.close()
		// A drop that failed is reported through Fatalf, after the tests passed.
		if setup.err != nil && code == 0 {
			fmt.Fprintf(os.Stderr, "test database cleanup failed: %v\n", setup.err)
			code = 1
		}
	}()

	mig := migrator{root: root, admin: conf}
	clone := pgtestdb.Custom(setup, conf, mig)
	if setup.err != nil {
		return setupFailed(setup.err)
	}

	urlMu.Lock()
	databaseURL = clone.URL()
	cleaner = dbcleaner.New(databaseURL)
	template.conf, template.m = conf, mig
	urlMu.Unlock()

	defer func() {
		if err := cleaner.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "test database cleaner close failed: %v\n", err)
		}
	}()
	return m.Run()
}

func serverConfig() pgtestdb.Config {
	return pgtestdb.Config{
		DriverName: "pgx",
		Host:       envOr("TEST_DB_HOST", defaultHost),
		Port:       envOr("TEST_DB_PORT", defaultPort),
		User:       adminUser,
		Password:   adminPass,
		Database:   "postgres",
		Options:    "sslmode=disable",
		// A package that leaks a connection must not turn into a failed cleanup
		// and a database left behind on the shared server.
		ForceTerminateConnections: true,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// requireServer fails fast with the fix when the daemon is down: pgtestdb's own
// error is a bare connection refusal.
func requireServer(conf pgtestdb.Config) error {
	addr := net.JoinHostPort(conf.Host, conf.Port)
	c, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return errs.With("addr", addr).Wrapf(err, "no Postgres (start it with `mise run db:up`, or set TEST_DB_HOST/TEST_DB_PORT)")
	}
	return c.Close()
}

// DatabaseURL returns the connection string for the current package's
// disposable database. It is only valid after Run has completed setup.
func DatabaseURL() string {
	urlMu.RLock()
	defer urlMu.RUnlock()
	if databaseURL == "" {
		panic("testdb: DatabaseURL called without testdb.Run in TestMain")
	}
	return databaseURL
}

// DB returns the package's database for one test. Everything the test and the
// code under test write, savepoints of their own included, is rolled back when
// the test ends (dbcleaner's Transaction strategy), so a test cleans nothing up.
func DB(t testing.TB) *sql.DB {
	t.Helper()
	DatabaseURL() // panics when Run did not set the package up
	urlMu.RLock()
	c := cleaner
	urlMu.RUnlock()
	return c.DB(t)
}

// migrator builds the template: the Atlas migrations, then the committed
// fixtures. Migrations run in-process through Atlas's library rather than
// pgtestdb's atlasmigrator, which shells out to an atlas binary the test run
// does not otherwise need. The fixture loader switches foreign-key triggers
// off, which takes a superuser, while migrations and the clones themselves run
// as pgtestdb's unprivileged role like the application does; so the fixtures
// go in through admin, the server's own credentials.
type migrator struct {
	root  string
	admin pgtestdb.Config
}

// Hash names the template. It covers what the template contains, so editing a
// migration or a fixture row builds a new template and the next run never
// sees a stale one. Old templates stay for pgtestdb to reuse if the content
// comes back.
func (mg migrator) Hash() (string, error) {
	h := common.NewRecursiveHash()
	if err := h.AddDirs(nil, "*.sql", filepath.Join(mg.root, "migrations")); err != nil {
		return "", errs.Wrapf(err, "hash migrations")
	}
	if err := h.AddDirs(nil, "*.yml", filepath.Join(mg.root, "fixtures")); err != nil {
		return "", errs.Wrapf(err, "hash fixtures")
	}
	return h.String(), nil
}

func (mg migrator) Migrate(ctx context.Context, db *sql.DB, template pgtestdb.Config) error {
	if err := applyMigrations(ctx, db, filepath.Join(mg.root, "migrations")); err != nil {
		return err
	}
	admin := mg.admin
	admin.Database = template.Database
	adminDB, err := admin.Connect()
	if err != nil {
		return errs.Wrapf(err, "connect to the template as admin")
	}
	defer func() { _ = adminDB.Close() }()
	return loadFixtures(adminDB, filepath.Join(mg.root, "fixtures"))
}

func applyMigrations(ctx context.Context, db *sql.DB, path string) error {
	dir, err := migrate.NewLocalDir(path)
	if err != nil {
		return errs.Wrapf(err, "open migration directory")
	}
	if err := migrate.Validate(dir); err != nil {
		return errs.Wrapf(err, "validate migration directory")
	}
	files, err := dir.Files()
	if err != nil {
		return errs.Wrapf(err, "read migration files")
	}
	driver, err := atlaspostgres.Open(db)
	if err != nil {
		return errs.Wrapf(err, "open Atlas Postgres driver")
	}
	executor, err := migrate.NewExecutor(driver, dir, migrate.NopRevisionReadWriter{})
	if err != nil {
		return errs.Wrapf(err, "create migration executor")
	}
	if err := executor.ExecuteFiles(ctx, files); err != nil {
		return errs.Wrapf(err, "apply test migrations")
	}
	return nil
}

func loadFixtures(db *sql.DB, path string) error {
	return fixture.Load(db, path, false)
}

func repositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errs.Wrapf(err, "get working directory")
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", errs.Wrapf(err, "inspect repository root")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errs.With("dir", dir).Errorf("find repository root: go.mod not found")
		}
		dir = parent
	}
}

// Clone gives the test a database of its own, cloned from the same template as
// the package's, and drops it when the test ends. Almost every test wants the
// rolled-back transaction testsupport hands out instead, which costs nothing;
// a clone is for the case where the commit is the thing under test, such as
// which of two connections sees an uncommitted write, or concurrent
// transactions racing for one row. What such a test commits dies with its
// clone, so it has nothing to clean up and cannot leak into another test.
func Clone(t testing.TB) *sql.DB {
	t.Helper()
	DatabaseURL() // panics when Run did not set the package up
	urlMu.RLock()
	conf, m := template.conf, template.m
	urlMu.RUnlock()
	return pgtestdb.New(t, conf, m)
}

func setupFailed(err error) int {
	fmt.Fprintf(os.Stderr, "test database setup failed: %v\n", err)
	return 1
}

// mainTB is the slice of testing.TB pgtestdb needs, for use outside a test. A
// Fatalf is recorded rather than stopping the process, since pgtestdb returns
// right after it; Run reads err once pgtestdb is back.
type mainTB struct {
	err      error
	cleanups []func()
}

func (b *mainTB) Cleanup(f func())    { b.cleanups = append(b.cleanups, f) }
func (b *mainTB) Failed() bool        { return b.err != nil }
func (b *mainTB) Helper()             {}
func (b *mainTB) Logf(string, ...any) {}
func (b *mainTB) Fatalf(format string, args ...any) {
	b.err = errors.Join(b.err, errs.Errorf(format, args...))
}

// close drops the package's database. Cleanups run last-in first-out, as in
// testing, and a panicking one does not stop the rest; its panic lands in err.
func (b *mainTB) close() {
	for i := len(b.cleanups) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					b.err = errors.Join(b.err, errs.With("panic", r).Errorf("cleanup panicked"))
				}
			}()
			b.cleanups[i]()
		}()
	}
}

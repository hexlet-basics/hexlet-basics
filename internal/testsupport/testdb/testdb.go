// Package testdb owns the lifecycle of the disposable Postgres used by a test
// package. Each package gets an independently migrated and seeded database, so
// `go test ./...` needs Docker but no separately prepared test service.
package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"ariga.io/atlas/sql/migrate"
	atlaspostgres "ariga.io/atlas/sql/postgres"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	postgresImage = "postgres:18-alpine"
	databaseName  = "code_basics_test"
	databaseUser  = "postgres"
	databasePass  = "postgres"
)

var (
	databaseURL string
	urlMu       sync.RWMutex
)

// Run starts, migrates, and seeds the package's Postgres before running its
// tests. TestMain must pass the returned status to os.Exit so the container can
// be terminated before the process exits.
func Run(m *testing.M) (code int) {
	root, err := repositoryRoot()
	if err != nil {
		return setupFailed(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	container, err := tcpostgres.Run(
		ctx,
		postgresImage,
		tcpostgres.WithDatabase(databaseName),
		tcpostgres.WithUsername(databaseUser),
		tcpostgres.WithPassword(databasePass),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		cancel()
		return setupFailed(fmt.Errorf("start Postgres container: %w", err))
	}
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			fmt.Fprintf(os.Stderr, "test database cleanup failed: %v\n", err)
			if code == 0 {
				code = 1
			}
		}
	}()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err == nil {
		err = prepareDatabase(ctx, dsn, root)
	}
	cancel()
	if err != nil {
		return setupFailed(err)
	}

	urlMu.Lock()
	databaseURL = dsn
	urlMu.Unlock()

	return m.Run()
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

func prepareDatabase(ctx context.Context, dsn, root string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open test database: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping test database: %w", err)
	}
	if err := applyMigrations(ctx, db, filepath.Join(root, "migrations")); err != nil {
		return err
	}
	if err := loadFixtures(db, filepath.Join(root, "fixtures")); err != nil {
		return err
	}
	return nil
}

func applyMigrations(ctx context.Context, db *sql.DB, path string) error {
	dir, err := migrate.NewLocalDir(path)
	if err != nil {
		return fmt.Errorf("open migration directory: %w", err)
	}
	if err := migrate.Validate(dir); err != nil {
		return fmt.Errorf("validate migration directory: %w", err)
	}
	files, err := dir.Files()
	if err != nil {
		return fmt.Errorf("read migration files: %w", err)
	}
	driver, err := atlaspostgres.Open(db)
	if err != nil {
		return fmt.Errorf("open Atlas Postgres driver: %w", err)
	}
	executor, err := migrate.NewExecutor(driver, dir, migrate.NopRevisionReadWriter{})
	if err != nil {
		return fmt.Errorf("create migration executor: %w", err)
	}
	if err := executor.ExecuteFiles(ctx, files); err != nil {
		return fmt.Errorf("apply test migrations: %w", err)
	}
	return nil
}

func loadFixtures(db *sql.DB, path string) error {
	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(path),
		testfixtures.SkipTableChecksumComputation(),
	)
	if err != nil {
		return fmt.Errorf("create fixture loader: %w", err)
	}
	if err := loader.Load(); err != nil {
		return fmt.Errorf("load fixtures: %w", err)
	}
	return nil
}

func repositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("inspect repository root: %w", err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("find repository root from %q: go.mod not found", dir)
		}
		dir = parent
	}
}

func setupFailed(err error) int {
	fmt.Fprintf(os.Stderr, "test database setup failed: %v\n", err)
	return 1
}

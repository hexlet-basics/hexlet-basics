package store_test

import (
	"os"
	"testing"

	"hexletbasics/internal/testsupport/testdb"
)

// The package's schema-level tests assert invariants PostgreSQL owns, so they
// need the migrated, fixture-loaded container. The in-package unit tests (the
// transaction seam, over a stub driver) share the binary and are unaffected.
func TestMain(m *testing.M) {
	os.Exit(testdb.Run(m))
}

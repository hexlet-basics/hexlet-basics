package transition_test

import (
	"os"
	"testing"

	"hexletbasics/internal/testsupport/testdb"
)

// The state move tests run the generated Move<Field> against the migrated,
// fixture-loaded database: the condition it writes under is SQL.
func TestMain(m *testing.M) {
	os.Exit(testdb.Run(m))
}

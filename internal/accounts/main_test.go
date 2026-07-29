package accounts

import (
	"os"
	"testing"

	"hexletbasics/internal/testsupport/testdb"
)

func TestMain(m *testing.M) {
	os.Exit(testdb.Run(m))
}

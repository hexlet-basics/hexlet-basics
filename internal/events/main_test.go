package events

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/goleak"

	"hexletbasics/internal/testsupport/testdb"
)

// goleak.VerifyTestMain exits the process itself, so it cannot wrap testdb.Run,
// which must drop the cloned database after the tests. The leak check runs
// after the clone is gone and only when the tests passed (a failure leaves
// goroutines of its own).
func TestMain(m *testing.M) {
	code := testdb.Run(m)
	if code == 0 {
		if err := goleak.Find(); err != nil {
			fmt.Fprintf(os.Stderr, "goleak: %v\n", err)

			code = 1
		}
	}

	os.Exit(code)
}

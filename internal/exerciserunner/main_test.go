package exerciserunner

import (
	"testing"

	"go.uber.org/goleak"
)

// The runner waits on containers and reports from goroutines; a leak after the
// tests means a run that never released its watcher.
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

package assistant_test

import (
	"testing"

	"go.uber.org/goleak"
)

// The streaming chat spawns a goroutine per answer; a test that leaves one
// behind (an unread pipe, a stream never drained) fails the package here.
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

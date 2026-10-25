package courseloader_test

import (
	"path/filepath"
	"testing"
)

// fixtureRepo is the committed minimal course (one module, one lesson, three
// locales), the canonical example of the on-disk format. Consumed by the loader
// integration test (loader_test.go).
func fixtureRepo(t *testing.T) string {
	t.Helper()

	return filepath.Join("testdata", "exercises")
}

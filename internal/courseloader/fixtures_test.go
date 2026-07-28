package courseloader_test

import (
	"path/filepath"
	"testing"
)

// fixtureRepo is the committed minimal course (one module, one lesson, three
// locales) shared with the legacy Rails loader test — the canonical example of
// the on-disk format. Consumed by the loader integration test (loader_test.go).
func fixtureRepo(t *testing.T) string {
	t.Helper()
	// internal/courseloader -> repo root is two up.
	return filepath.Join("..", "..", "legacy", "test", "fixtures", "files", "exercises")
}

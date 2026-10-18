package mcpserver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// as asserts that v holds a T and returns it: the checked form of v.(T) for
// walking decoded JSON, which fails the test with the types instead of panicking.
func as[T any](t *testing.T, v any) T {
	t.Helper()

	got, ok := v.(T)
	require.True(t, ok, "got %T, want %T", v, got)

	return got
}

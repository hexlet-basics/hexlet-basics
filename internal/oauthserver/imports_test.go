package oauthserver_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestThePackageDoesNotDependOnTheHTTPHandlers fails if internal/oauthserver
// starts importing internal/handlers, directly or through another package
// (ADR-0021): the server is meant to be liftable out of the application and
// reaches it only through the Application interface. `go list -deps` follows
// transitive imports, which a lint rule on direct imports would miss. The
// package's own tests are excluded on purpose: they use testsupport, which
// imports handlers.
func TestThePackageDoesNotDependOnTheHTTPHandlers(t *testing.T) {
	out, err := exec.CommandContext(t.Context(), "go", "list", "-deps", "hexletbasics/internal/oauthserver").CombinedOutput()
	require.NoError(t, err, string(out))

	deps := strings.Fields(string(out))
	require.Contains(t, deps, "hexletbasics/internal/oauthserver", "go list did not report the package itself")
	for _, forbidden := range []string{
		"hexletbasics/internal/handlers",
		"hexletbasics/internal/testsupport",
	} {
		require.NotContains(t, deps, forbidden, "internal/oauthserver must not depend on %s", forbidden)
	}
}

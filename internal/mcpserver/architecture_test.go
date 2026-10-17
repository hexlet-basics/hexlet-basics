package mcpserver_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestThePackageCannotReachAnOperationExceptThroughTheGeneratedServer is the
// structural half of ADR-0019's "one checkpoint" (the behavioural half is
// TestNoRegisteredAdminToolRunsWithoutTheCheck). The permission check lives in
// the generated server's security handler, so a tool that reached a handler
// method, the database or the generated server's constructor would skip it. The
// package therefore may not import internal/handlers or the store (and its admin files not ent), nor name
// api.Handler, api.UnimplementedHandler or api.NewServer: what it holds of the
// website is an http.Handler it is given, and the generated types it describes.
func TestThePackageCannotReachAnOperationExceptThroughTheGeneratedServer(t *testing.T) {
	out, err := exec.CommandContext(t.Context(), "go", "list", "-f", `{{join .Imports "\n"}}`, "hexletbasics/internal/mcpserver").CombinedOutput()
	require.NoError(t, err, string(out))
	for _, imported := range strings.Fields(string(out)) {
		for _, forbidden := range []string{"hexletbasics/internal/handlers", "hexletbasics/internal/store", "database/sql"} {
			assert.False(t, imported == forbidden || strings.HasPrefix(imported, forbidden+"/"),
				"internal/mcpserver must not import %s", imported)
		}
	}

	// ent is imported by the shared tool code for its error helpers (IsNotFound),
	// which reach no data; the admin files, which hold the dispatch, may not.
	for _, name := range []string{"admin.go", "operation.go", "principal.go"} {
		parsed, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
		require.NoError(t, err)
		for _, spec := range parsed.Imports {
			assert.NotContains(t, spec.Path.Value, "hexletbasics/ent", "%s must not reach the database", name)
		}
	}

	files, err := filepath.Glob("*.go")
	require.NoError(t, err)
	checked := 0
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		checked++
		src, err := os.ReadFile(name)
		require.NoError(t, err)
		parsed, err := parser.ParseFile(token.NewFileSet(), name, src, 0)
		require.NoError(t, err)
		ast.Inspect(parsed, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "api" {
					assert.NotContains(t, []string{"Handler", "UnimplementedHandler", "NewServer"}, sel.Sel.Name,
						"%s names api.%s", name, sel.Sel.Name)
				}
			}
			return true
		})
	}
	assert.Positive(t, checked)
}

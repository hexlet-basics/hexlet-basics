package permissioncatalog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissioncatalog"
)

// x-requires-capability must name a capability the contract lists, and sit on an
// operation that has an id, so a typo cannot silently leave an operation
// without the capability it was meant to need.
func TestRequiredCapabilityMustNameAListedCapability(t *testing.T) {
	t.Parallel()
	for name, doc := range map[string]string{
		"an unlisted capability": `
paths:
  /api/admin/x:
    get:
      operationId: adminX
      x-permission: {name: x.read, kind: read}
      x-requires-capability: mcp.nothing`,
		"without an operation id": `
x-capabilities:
  - {name: mcp.manage, kind: write}
paths:
  /api/admin/x:
    get:
      x-requires-capability: mcp.manage`,
	} {
		_, err := permissioncatalog.Parse([]byte(doc))
		require.Error(t, err, name)
	}

	_, err := permissioncatalog.Parse([]byte(`
x-capabilities:
  - {name: mcp.manage, kind: write}
paths:
  /api/admin/x:
    get:
      operationId: adminX
      x-permission: {name: x.read, kind: read}
      x-requires-capability: mcp.manage`))
	assert.NoError(t, err)
}

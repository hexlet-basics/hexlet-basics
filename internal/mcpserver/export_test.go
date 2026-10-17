package mcpserver

import (
	"maps"
	"slices"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AllowForTest replaces the allowlist for the test that follows, to exercise the
// pipeline with an operation the registry does not offer yet. It must run before
// the server is built, which reads the list once.
func AllowForTest(t *testing.T, operationIDs ...string) {
	t.Helper()
	was := adminOperations
	t.Cleanup(func() { adminOperations = was })
	adminOperations = operationIDs
}

// SessionCookie is the cookie name the dispatcher puts the placeholder under.
const SessionCookie = sessionCookie

// ContractAdminOperations are the operationIds of every operation of the
// contract that requires the AdminSession scheme, sorted.
func ContractAdminOperations(t *testing.T) []string {
	t.Helper()
	doc, err := contract()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, methods := range doc.Paths {
		for _, raw := range methods {
			node, _ := raw.(map[string]any)
			id, _ := node["operationId"].(string)
			if id == "" {
				continue
			}
			// Read from the raw security requirement, not from a built operation:
			// an operation a tool cannot yet describe must still be counted.
			for _, requirement := range asList(node["security"]) {
				if m, ok := requirement.(map[string]any); ok {
					if _, admin := m["AdminSession"]; admin {
						ids = append(ids, id)
						break
					}
				}
			}
		}
	}
	slices.Sort(ids)
	return ids
}

// AllowlistedAdminOperations is the allowlist as the registry reads it.
func AllowlistedAdminOperations() []string { return slices.Clone(adminOperations) }

// ExcludedAdminOperations are the operationIds deliberately not tools.
func ExcludedAdminOperations() map[string]string { return maps.Clone(excludedAdminOperations) }

// AdminToolsFromTheContract builds the tool of every allowlisted operation, the
// way the registry does at start-up.
func AdminToolsFromTheContract(t *testing.T) []*mcp.Tool {
	t.Helper()
	var tools []*mcp.Tool
	for _, id := range adminOperations {
		op, err := readOperation(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		tool, err := adminTool(op)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		tools = append(tools, tool)
	}
	return tools
}

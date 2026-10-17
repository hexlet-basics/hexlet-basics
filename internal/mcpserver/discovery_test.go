package mcpserver_test

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/audit"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// allCapabilities is every capability Permission, for a role that is to see
// every admin tool.
var allCapabilities = []string{"mcp.read", "mcp.write", "mcp.manage", "audit.read"}

// everything connects bob, whose role is granted every capability, with every
// Scope: the caller that sees the whole tool list.
func everything(t *testing.T, o *testsupport.OAuth) *mcp.ClientSession {
	t.Helper()
	o.Grant(editor, allCapabilities...)
	tokens := o.Connect(bob, oauthserver.ScopeCatalogRead, oauthserver.ScopeProgressRead, oauthserver.ScopeAdmin)
	return o.MCPClient(tokens.AccessToken)
}

func listTools(t *testing.T, session *mcp.ClientSession) *mcp.ListToolsResult {
	t.Helper()
	list, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)
	return list
}

func names(tools []*mcp.Tool) []string {
	out := make([]string, len(tools))
	for i, tool := range tools {
		out[i] = tool.Name
	}
	return out
}

// The instructions are what a client reads once, so they are short (Claude Code
// cuts at 4 KB), name the domains and say when to look tools up.
func TestInstructionsNameTheDomainsAndSayWhenToSearch(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := everything(t, o)

	instructions := session.InitializeResult().Instructions
	assert.NotEmpty(t, instructions)
	assert.Less(t, len(instructions), 4096)
	for _, domain := range []string{"catalogue", "learner progress", "content administration", "access management", "audit"} {
		assert.Contains(t, instructions, domain)
	}
	assert.Contains(t, instructions, "findTools")
	assert.Contains(t, instructions, "does not load tools on demand")
}

// Annotations are explicit on every tool: the protocol's defaults (destructive,
// open-world) are the unsafe ones, and a client that reads a missing hint as the
// default would ask its user to confirm a read.
func TestEveryToolSaysWhatItDoesToTheWorld(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tools := listTools(t, everything(t, o)).Tools

	require.Greater(t, len(tools), len(mcpserver.AllowlistedAdminOperations()))
	for _, tool := range tools {
		require.NotNil(t, tool.Annotations, tool.Name)
		assert.NotNil(t, tool.Annotations.DestructiveHint, "%s: destructiveHint", tool.Name)
		assert.NotNil(t, tool.Annotations.OpenWorldHint, "%s: openWorldHint", tool.Name)
		// readOnlyHint and idempotentHint are plain booleans the protocol always
		// writes; what is left to check is that they agree with each other.
		if tool.Annotations.ReadOnlyHint {
			assert.False(t, *tool.Annotations.DestructiveHint, "%s: a read destroys nothing", tool.Name)
		}
		assert.False(t, *tool.Annotations.OpenWorldHint, "%s: nothing here leaves the site's own data", tool.Name)
	}
}

// The tool list is built, then checked as a client would load it.
func TestTheToolListIsValidStableAndWithinItsBudget(t *testing.T) {
	// The size of the whole list for the caller who sees every tool, measured
	// when this budget was set (ADR-0021, #878); a tool or schema that grows it
	// past the budget is a decision to make on purpose.
	const budget = 100 * 1024

	o := testsupport.NewOAuth(t)
	session := everything(t, o)
	first := listTools(t, session)
	serverName := session.InitializeResult().ServerInfo.Name

	total := 0
	for _, tool := range first.Tools {
		raw, err := json.Marshal(tool)
		require.NoError(t, err)
		total += len(raw)

		// The client-side full name is mcp__<server>__<tool>, and Claude Code
		// refuses one of 128 characters or more.
		assert.Less(t, len("mcp__"+serverName+"__"+tool.Name), 128, tool.Name)
		assert.NotEmpty(t, tool.Description, tool.Name)
		assert.LessOrEqual(t, len([]rune(tool.Description)), 160, "%s: a description is read from its start", tool.Name)

		for kind, schema := range map[string]any{"input": tool.InputSchema, "output": tool.OutputSchema} {
			if schema == nil {
				continue
			}
			raw, err := json.Marshal(schema)
			require.NoError(t, err)
			var decoded map[string]any
			require.NoError(t, json.Unmarshal(raw, &decoded))
			assert.Equal(t, "object", decoded["type"], "%s %s: the root is an object", tool.Name, kind)
			for _, combined := range []string{"anyOf", "oneOf", "allOf"} {
				assert.NotContains(t, decoded, combined, "%s %s: Claude Code flattens a root %s and loses its enforcement", tool.Name, kind, combined)
			}

			var parsed jsonschema.Schema
			require.NoError(t, json.Unmarshal(raw, &parsed), "%s %s", tool.Name, kind)
			_, err = parsed.Resolve(nil)
			assert.NoError(t, err, "%s %s is a valid JSON Schema (2020-12)", tool.Name, kind)
		}
	}
	t.Logf("tool list: %d tools, %d bytes (budget %d)", len(first.Tools), total, budget)
	assert.Less(t, total, budget)

	// Deterministic order, the same on every call, and by name: a client that
	// caches the list keeps its cache when nothing changed.
	second := listTools(t, session)
	assert.Equal(t, names(first.Tools), names(second.Tools))
	assert.True(t, slices.IsSorted(names(first.Tools)), "tools are listed by name")
}

// The list is private to the caller that asked, because it differs by caller.
func TestToolListResultsArePrivateToTheCaller(t *testing.T) {
	o := testsupport.NewOAuth(t)
	list := listTools(t, everything(t, o))
	assert.Equal(t, "private", list.CacheScope)
}

// A token is shown what it may use. The checkpoint still judges each call: a
// tool that is hidden and called by name is refused as it always was, and
// nothing is recorded as done.
func TestAConnectionIsOnlyShownWhatItsGrantsAllow(t *testing.T) {
	o := testsupport.NewOAuth(t)

	// Without the admin Scope: the catalogue, the lookup, and no admin tool.
	catalogueOnly := o.MCPClient(o.Connect(bob, oauthserver.ScopeCatalogRead).AccessToken)
	shown := names(listTools(t, catalogueOnly).Tools)
	assert.Contains(t, shown, "listCourses")
	assert.Contains(t, shown, "findTools")
	assert.NotContains(t, shown, "getMyDashboard", "no progress:read Scope")
	for _, name := range shown {
		assert.False(t, strings.HasPrefix(name, "admin"), "%s needs the admin Scope", name)
	}

	// With the admin Scope but a role that holds no capability: nothing to offer
	// that the door would not close at once.
	adminScope := o.MCPClient(o.Connect(bob, oauthserver.ScopeAdmin).AccessToken)
	shown = names(listTools(t, adminScope).Tools)
	assert.NotContains(t, shown, "adminListBanners", "bob's role holds no mcp.read")

	// mcp.read opens the reads; mcp.write the writes; mcp.manage the access
	// operations; audit.read the audit trail.
	o.Grant(editor, "mcp.read")
	shown = names(listTools(t, adminScope).Tools)
	assert.Contains(t, shown, "adminListBanners")
	assert.NotContains(t, shown, "adminUpdateBanner")
	assert.NotContains(t, shown, "adminListUsers")
	assert.NotContains(t, shown, "adminListAuditRecords")

	o.Grant(editor, "mcp.write")
	assert.Contains(t, names(listTools(t, adminScope).Tools), "adminUpdateBanner")
	o.Grant(editor, "mcp.manage", "audit.read")
	shown = names(listTools(t, adminScope).Tools)
	assert.Contains(t, shown, "adminListUsers")
	assert.Contains(t, shown, "adminListAuditRecords")

	// Taking a grant away is seen by the next list.
	o.RevokeGrant(editor, "mcp.write")
	assert.NotContains(t, names(listTools(t, adminScope).Tools), "adminUpdateBanner")
}

func TestHidingAToolIsNotAuthorizingTheOthers(t *testing.T) {
	o := testsupport.NewOAuth(t)
	before := o.DB.Banner.Query().CountX(t.Context())

	// bob holds no mcp.write, and a token without the admin Scope would not be
	// offered adminCreateBanner either. Both call it by name anyway.
	o.Grant(editor, "mcp.read")
	adminScope := o.MCPClient(o.Connect(bob, oauthserver.ScopeAdmin).AccessToken)
	require.NotContains(t, names(listTools(t, adminScope).Tools), "adminCreateBanner")
	assertRefused(t, call(t, adminScope, "adminCreateBanner", map[string]any{
		"locale": "en", "body": bannerInput("created by name"),
	}))

	catalogueOnly := o.MCPClient(o.Connect(bob, oauthserver.ScopeCatalogRead).AccessToken)
	assertRefused(t, call(t, catalogueOnly, "adminListBanners", map[string]any{"locale": "en"}))

	assert.Equal(t, before, o.DB.Banner.Query().CountX(t.Context()), "nothing was created")
	for _, record := range o.Records("AdminCreateBanner") {
		assert.NotEqual(t, audit.OutcomeSucceeded, record.Outcome, "a call refused at the door is never a succeeded record")
	}
}

func findTools(t *testing.T, session *mcp.ClientSession, args map[string]any) (total int, tools []map[string]any) {
	t.Helper()
	content := structured(t, callTool(t, session, "findTools", args))
	total = int(content["total"].(float64))
	for _, row := range content["tools"].([]any) {
		tools = append(tools, row.(map[string]any))
	}
	return total, tools
}

func TestFindToolsLooksToolsUpAtEachLevelOfDetail(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := everything(t, o)

	total, rows := findTools(t, session, map[string]any{"query": "banner", "detail": "name"})
	assert.Equal(t, len(rows), total)
	for _, row := range rows {
		assert.Contains(t, strings.ToLower(row["name"].(string)), "banner")
		assert.Equal(t, []string{"name"}, sortedKeys(row), "name only")
	}
	assert.GreaterOrEqual(t, len(rows), 6, "a group of tools by resource: list, get, create, update, publish, delete...")

	_, rows = findTools(t, session, map[string]any{"query": "update banner"})
	require.Len(t, rows, 1, "every word must match")
	row := rows[0]
	assert.Equal(t, "adminUpdateBanner", row["name"])
	assert.NotEmpty(t, row["description"])
	assert.NotContains(t, row, "inputSchema", "the default level is without the schema")
	assert.Equal(t, map[string]any{
		"scope": "admin", "permission": "banners.update", "capabilities": []any{"mcp.write"},
	}, row["needs"])

	_, rows = findTools(t, session, map[string]any{"query": "update banner", "detail": "schema"})
	require.Len(t, rows, 1)
	schema, ok := rows[0]["inputSchema"].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, schema["properties"], "body")

	_, rows = findTools(t, session, map[string]any{"query": "banner", "limit": 2})
	assert.Len(t, rows, 2)
	total, _ = findTools(t, session, map[string]any{"query": "banner", "limit": 2})
	assert.Greater(t, total, 2, "the total counts what the limit left out")
}

// findTools shows a caller only the tools it may use, as tools/list does.
func TestFindToolsDoesNotNameWhatTheCallerIsNotOffered(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := o.MCPClient(o.Connect(bob, oauthserver.ScopeCatalogRead).AccessToken)

	total, rows := findTools(t, session, map[string]any{"query": "banner"})
	assert.Zero(t, total)
	assert.Empty(t, rows)
	_, rows = findTools(t, session, map[string]any{"query": "course"})
	assert.Contains(t, mapNames(rows), "listCourses")
	assert.NotContains(t, mapNames(rows), "adminListCourses")
}

// There is one way to look a tool up and none to run one by name: running goes
// through the tool's own registration, scope, rate limit and Audit Record.
func TestNoToolRunsAnotherToolByName(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := everything(t, o)
	list := listTools(t, session)

	for _, tool := range list.Tools {
		assert.NotRegexp(t, `(?i)^(call|execute|exec|run|invoke|dispatch|proxy)`, tool.Name, "a tool that runs others by name would bypass the checkpoint")
		raw, err := json.Marshal(tool.InputSchema)
		require.NoError(t, err)
		var schema struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		require.NoError(t, json.Unmarshal(raw, &schema))
		for _, argument := range []string{"tool", "toolName", "tool_name", "operationId", "code", "script"} {
			assert.NotContains(t, schema.Properties, argument, "%s takes %s", tool.Name, argument)
		}
	}

	// And findTools is read-only in fact: asked to run a tool, it refuses the
	// argument and leaves nothing behind.
	find := byName(list.Tools)["findTools"]
	assert.True(t, find.Annotations.ReadOnlyHint)
	before := o.DB.Banner.Query().CountX(t.Context())
	result := call(t, session, "findTools", map[string]any{
		"tool": "adminCreateBanner", "arguments": map[string]any{"locale": "en", "body": bannerInput("run by name")},
	})
	assert.True(t, result.IsError)
	assert.Equal(t, before, o.DB.Banner.Query().CountX(t.Context()))
	assert.Empty(t, o.Records("AdminCreateBanner"))
}

func byName(tools []*mcp.Tool) map[string]*mcp.Tool {
	out := map[string]*mcp.Tool{}
	for _, tool := range tools {
		out[tool.Name] = tool
	}
	return out
}

func sortedKeys(m map[string]any) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

func mapNames(rows []map[string]any) []string {
	var out []string
	for _, row := range rows {
		out = append(out, row["name"].(string))
	}
	return out
}

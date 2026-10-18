package mcpserver_test

import (
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// adminSearch and adminFetch are the content tools of the staff (#877): the
// catalogue search with the storefront's filters off, offered like any other
// admin tool, so they need the admin Scope and mcp.read, and a refusal is
// recorded like any other.

func searchItems(t *testing.T, result *mcp.CallToolResult) []map[string]any {
	t.Helper()
	require.False(t, result.IsError, "%+v", result.Content)
	listed := result.StructuredContent.(map[string]any)["items"].([]any)

	items := make([]map[string]any, 0, len(listed))
	for _, item := range listed {
		items = append(items, item.(map[string]any))
	}

	return items
}

func itemByID(t *testing.T, items []map[string]any, id string) map[string]any {
	t.Helper()

	for _, item := range items {
		if item["id"] == id {
			return item
		}
	}

	require.Failf(t, "not found", "%s is not among the results", id)

	return nil
}

func TestAnAdminFindsAndFetchesWhatTheStorefrontHides(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, bob)

	// A course not yet completed, on an unlisted page, is one more record the
	// storefront hides.
	o.DB.CourseLandingPage.Update().Where(courselandingpage.SlugEQ("ruby-ru")).SetListed(false).ExecX(t.Context())
	o.DB.Course.Update().Where(course.SlugEQ("ruby")).SetReadiness("in_development").ExecX(t.Context())

	for name, tc := range map[string]struct {
		locale, query, kind, id, state, readiness string
	}{
		"a draft landing page":                      {"ru", "Node.js", "course", "course:nodejs-ru", "draft", "completed"},
		"an archived landing page":                  {"ru", "архив", "course", "course:javascript-ru-archived", "archived", "completed"},
		"an unlisted course that is in development": {"ru", "Ruby", "course", "course:ruby-ru", "published", "in_development"},
		"a lesson of a course in development":       {"ru", "первая программа", "lesson", "lesson:ruby/hello-world-ruby", "", "in_development"},
		"an archived post":                          {"ru", "Archived notes", "post", "post:archived-notes", "archived", ""},
		"a draft post of the other locale":          {"en", "Second post", "post", "post:second-post", "draft", ""},
	} {
		t.Run(name, func(t *testing.T) {
			found := call(t, session, "adminSearch", map[string]any{"locale": tc.locale, "q": tc.query, "type": tc.kind})

			hit := itemByID(t, searchItems(t, found), tc.id)
			assert.Equal(t, tc.kind, hit["type"])
			assert.Equal(t, tc.readiness, hit["readiness"])

			fetched := call(t, session, "adminFetch", map[string]any{"locale": tc.locale, "id": tc.id})

			require.False(t, fetched.IsError, "%+v", fetched.Content)
			record := fetched.StructuredContent.(map[string]any)
			assert.Equal(t, tc.id, record["id"])
			assert.Equal(t, hit["title"], record["title"])
			assert.NotEmpty(t, record["text"])
			meta := record["metadata"].(map[string]any)
			assert.Equal(t, tc.kind, meta["type"])
			assert.Equal(t, tc.readiness, meta["readiness"], "the record says how ready its course is")

			if tc.state != "" {
				assert.Equal(t, tc.state, hit["state"])
				assert.Equal(t, tc.state, meta["state"], "the record says how published it is")
			}
		})
	}
}

// The locale is the one the caller asks for, as for every admin tool: a record of
// another locale is not found there.
func TestAdminSearchAndFetchStayWithinTheLocaleOfTheCall(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, bob)

	russian := call(t, session, "adminSearch", map[string]any{"locale": "ru", "q": "Second post"})
	english := call(t, session, "adminSearch", map[string]any{"locale": "en", "q": "Second post"})
	fetched := call(t, session, "adminFetch", map[string]any{"locale": "ru", "id": "post:second-post"})

	for _, item := range searchItems(t, russian) {
		assert.NotEqual(t, "post:second-post", item["id"])
	}

	itemByID(t, searchItems(t, english), "post:second-post")
	require.True(t, fetched.IsError)
	assert.Contains(t, textOf(t, fetched), "HTTP 404")
}

func TestAdminFetchOfAnIdThatIsNotOneSearchReturnsIsNotFound(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, bob)

	for _, id := range []string{"javascript-ru", "page:home", "lesson:javascript", "post:", "post:no-such-post", "course:no-such-course"} {
		result := call(t, session, "adminFetch", map[string]any{"locale": "ru", "id": id})

		require.True(t, result.IsError, id)
		assert.Contains(t, textOf(t, result), "HTTP 404", id)
	}
}

// The public search is unchanged by the admin one: the same draft is not found by
// a visitor.
func TestThePublicSearchStillHidesWhatTheAdminSearchFinds(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Do("GET", "/api/search?q=Node.js", nil, russian)
	defer func() { _ = res.Body.Close() }()

	assert.NotContains(t, body(t, res), "nodejs-ru")
}

func TestOnlyStaffWithMCPReadAndTheAdminScopeAreOfferedTheSearchTools(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tools := []string{"adminSearch", "adminFetch"}

	// A connection without the admin Scope, the public ones included.
	catalogueOnly := o.MCPClient(o.Connect(bob, oauthserver.ScopeCatalogRead, oauthserver.ScopeProgressRead).AccessToken)
	shown := names(listTools(t, catalogueOnly).Tools)
	assert.Contains(t, shown, "listCourses", "the catalogue tools keep working")
	assert.NotContains(t, shown, "search")
	assert.NotContains(t, shown, "fetch")

	for _, name := range tools {
		assert.NotContains(t, shown, name, "the admin Scope is missing")
	}

	// A staff member whose role does not hold mcp.read.
	adminScope := o.MCPClient(o.Connect(bob, oauthserver.ScopeAdmin).AccessToken)

	shown = names(listTools(t, adminScope).Tools)
	for _, name := range tools {
		assert.NotContains(t, shown, name, "the role holds no mcp.read")
	}

	// A user who is not staff at all.
	notStaff := o.MCPClient(o.Connect(alice, oauthserver.ScopeAdmin).AccessToken)
	for _, name := range tools {
		assert.NotContains(t, names(listTools(t, notStaff).Tools), name)
	}

	o.Grant(editor, "mcp.read")

	shown = names(listTools(t, adminScope).Tools)
	for _, name := range tools {
		assert.Contains(t, shown, name)
	}

	assert.NotContains(t, shown, "search")
	assert.NotContains(t, shown, "fetch")
}

func TestTheSearchToolsAreRefusedByNameAndRecordedWhenTheCallerMayNotUseThem(t *testing.T) {
	arguments := map[string]map[string]any{
		"adminSearch": {"locale": "ru", "q": "Node.js"},
		"adminFetch":  {"locale": "ru", "id": "course:nodejs-ru"},
	}
	operations := map[string]string{"adminSearch": "AdminSearch", "adminFetch": "AdminFetch"}

	t.Run("a token without the admin Scope", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		o.Grant(editor, "mcp.read")

		session := o.MCPClient(o.Connect(bob, oauthserver.ScopeCatalogRead).AccessToken)
		for tool, args := range arguments {
			result := call(t, session, tool, args)

			assertRefused(t, result)
			assert.NotContains(t, textOf(t, result), "Node.js", "nothing of the record is told")

			records := o.Records(operations[tool])
			require.Len(t, records, 1, tool)
			assert.Equal(t, "forbidden", records[0].Outcome)
			assert.Equal(t, "mcp", records[0].Channel)
		}
	})

	t.Run("a staff member whose role holds no mcp.read", func(t *testing.T) {
		o := testsupport.NewOAuth(t)

		session, _ := admin(t, o, bob)
		for tool, args := range arguments {
			assertRefused(t, call(t, session, tool, args))
			records := o.Records(operations[tool])
			require.Len(t, records, 1, tool)
			assert.Equal(t, "forbidden", records[0].Outcome)
		}
	})

	t.Run("an administrator who is not staff", func(t *testing.T) {
		o := testsupport.NewOAuth(t)

		session, _ := admin(t, o, alice)
		for tool, args := range arguments {
			assertRefused(t, call(t, session, tool, args))
			records := o.Records(operations[tool])
			require.Len(t, records, 1, tool)
			assert.Equal(t, "forbidden", records[0].Outcome)
		}
	})

	t.Run("a locale the staff role does not allow", func(t *testing.T) {
		o := testsupport.NewOAuth(t)
		o.Grant(editor, "mcp.read")
		o.DB.StaffMember.Update().SetAllowedLocales([]string{"en"}).ExecX(t.Context())
		session, _ := admin(t, o, bob)

		assertRefused(t, call(t, session, "adminSearch", arguments["adminSearch"]))
	})
}

// A convention that tells who a tool is for: whatever needs the admin Scope is
// called admin...; the lookup of tools is not an admin tool and stays as it is.
func TestEveryToolThatNeedsTheAdminScopeIsNamedAdmin(t *testing.T) {
	scopes := mcpserver.ToolScopes(t)
	require.Contains(t, scopes, "adminSearch")

	for name, scope := range scopes {
		if scope == oauthserver.ScopeAdmin {
			assert.True(t, strings.HasPrefix(name, "admin"), "%s needs the admin Scope and must be named for it", name)
		} else {
			assert.False(t, strings.HasPrefix(name, "admin"), "%s is named admin and needs the %q Scope", name, scope)
		}
	}

	assert.NotEqual(t, oauthserver.ScopeAdmin, scopes["findTools"], "findTools is not an admin tool")
	require.Contains(t, scopes, "findTools")
	assert.NotContains(t, scopes, "search")
	assert.NotContains(t, scopes, "fetch")
}

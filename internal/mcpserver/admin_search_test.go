package mcpserver_test

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/testsupport"
)

// The text filter `q` of the admin lists is a contract argument, so it reaches
// the tools with no MCP code of its own: these cases pin that it does, that it
// narrows what the caller may list and nothing else, and that searching is
// still reading (the Permission, and the record of a read of personal data).

// searchCase is one admin list tool and a term that finds some of its rows.
type searchCase struct {
	tool string
	// permission is what the operation's contract declares; the editor role
	// holds none of them for these resources but banners and blog posts, which
	// the refusal cases take away.
	permission string
	// extra are the grants the operation needs beside it (mcp.manage for users).
	extra []string
	// q is a term whose case differs from the stored text, so the match is
	// case-insensitive too.
	q string
}

var searchCases = []searchCase{
	{tool: "adminListBanners", permission: "banners.read", q: "BLACK"},
	{tool: "adminListBlogPosts", permission: "blog_posts.read", q: "HELLO WORLD"},
	{tool: "adminListReviews", permission: "reviews.read", q: "PETROV"},
	{tool: "adminListCourses", permission: "courses.read", q: "PYTHON"},
	{tool: "adminListCourseLandingPages", permission: "course_landing_pages.read", q: "RUBY-"},
	{tool: "adminListCourseLessons", permission: "course_lessons.read", q: "VARIABLES"},
	{tool: "adminListLeads", permission: "leads.read", q: "ANNA@"},
	{tool: "adminListMessages", permission: "messages.read", q: "CONSOLE.LOG"},
	{tool: "adminListUsers", permission: "users.read", extra: []string{"mcp.manage"}, q: "SECOND-ADMIN@"},
}

func (c searchCase) args(q string) map[string]any {
	args := map[string]any{"locale": "en"}
	if q != "" {
		args["q"] = q
	}
	return args
}

func TestEveryAdminSearchToolDeclaresQInItsInputSchema(t *testing.T) {
	o := testsupport.NewOAuth(t)
	// The list a caller gets holds the tools it may call, so it holds them all.
	o.Grant(editor, "mcp.read", "mcp.manage")
	for _, c := range searchCases {
		o.Grant(editor, c.permission)
	}
	session, _ := admin(t, o, bob)

	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)
	byName := map[string]*mcp.Tool{}
	for _, tool := range tools.Tools {
		byName[tool.Name] = tool
	}

	for _, c := range searchCases {
		tool, ok := byName[c.tool]
		require.True(t, ok, c.tool)
		raw, err := json.Marshal(tool.InputSchema)
		require.NoError(t, err)
		var schema struct {
			Properties map[string]struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"properties"`
			Required []string `json:"required"`
		}
		require.NoError(t, json.Unmarshal(raw, &schema))
		q, ok := schema.Properties["q"]
		assert.True(t, ok, "%s has no q", c.tool)
		assert.Equal(t, "string", q.Type, c.tool)
		assert.NotEmpty(t, q.Description, c.tool)
		assert.NotContains(t, schema.Required, "q", "%s: the filter is optional", c.tool)
	}
}

func TestAnAdminSearchToolReturnsOnlyTheRowsTheTermFinds(t *testing.T) {
	for _, c := range searchCases {
		t.Run(c.tool, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			o.Grant(editor, append([]string{c.permission, "mcp.read"}, c.extra...)...)
			session, _ := admin(t, o, bob)

			everything := structured(t, call(t, session, c.tool, c.args("")))
			found := structured(t, call(t, session, c.tool, c.args(c.q)))

			items := found["items"].([]any)
			require.NotEmpty(t, items)
			for _, item := range items {
				raw, err := json.Marshal(item)
				require.NoError(t, err)
				assert.Contains(t, strings.ToLower(string(raw)), strings.ToLower(c.q))
			}
			assert.Less(t, found["total"], everything["total"], "the total follows the filter")
		})
	}
}

func TestAnAdminSearchToolIsRefusedToARoleThatMayNotListTheResource(t *testing.T) {
	for _, c := range searchCases {
		t.Run(c.tool, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			o.RevokeGrant(editor, c.permission)
			o.Grant(editor, append([]string{"mcp.read"}, c.extra...)...)
			session, _ := admin(t, o, bob)

			assertRefused(t, call(t, session, c.tool, c.args(c.q)))
		})
	}
}

// A search over personal data is the read of those rows, recorded as any list
// is: the ids of the rows it answered, not of the rows it hid.
func TestASearchOverPersonalDataIsRecordedAsAReadOfTheRowsItAnswered(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "leads.read", "mcp.read")
	session, _ := admin(t, o, bob)

	page := structured(t, call(t, session, "adminListLeads", map[string]any{"locale": "en", "q": "anna@"}))
	testsupport.Deliver(t, o.Services, eventhandlers.StaffReadPersonalData(o.DB))

	records := o.Records("AdminListLeads")
	require.Len(t, records, 1)
	assert.Equal(t, "mcp", records[0].Channel)
	assert.Equal(t, "Lead", *records[0].EntityType)
	var reads struct {
		Count int      `json:"count"`
		IDs   []string `json:"ids"`
	}
	require.NoError(t, json.Unmarshal(records[0].Reads, &reads))
	require.Equal(t, 1, reads.Count)
	assert.Equal(t, []string{strconv.Itoa(itemIDs(t, page)[0])}, reads.IDs)
}

func TestARefusedSearchOverPersonalDataIsRecordedAsRefusedAndReadsNothing(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read") // leads.read is not held
	session, _ := admin(t, o, bob)

	assertRefused(t, call(t, session, "adminListLeads", map[string]any{"locale": "en", "q": "anna@"}))

	records := o.Records("AdminListLeads")
	require.Len(t, records, 1)
	assert.Equal(t, "forbidden", records[0].Outcome)
}

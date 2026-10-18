package mcpserver

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/internal/oauthserver"
)

// What a client needs to keep the tool list cheap (ADR-0021 addendum, #878).
// Searching and loading tools lazily is the client's feature, and a server
// cooperates: it says early what it has, keeps names and order stable, shows a
// caller what the caller may use, and offers one read-only way to look tools up
// for a client that cannot defer them. It never offers a way to run a tool by
// name, which would bypass the per-operation checkpoint and its Audit Record.

const (
	// serverName is the name the server announces; a client prefixes it to every
	// tool name, and a test keeps the whole under the 128 characters clients allow.
	serverName = "hexlet-basics"

	// findToolsName is the one discovery tool.
	findToolsName = "findTools"
)

// instructionsText is what the client reads once at the start. It stays well
// under 4 KB, where Claude Code cuts it. It names the domains and says when to
// search for tools, and it tells once what the cut-result hint contains, so the
// hint's description is not repeated on every tool.
const instructionsText = `Hexlet Basics teaches programming through courses of short lessons. This server reads the public catalogue, a learner's own progress, and lets staff administer the site.

Domains, by the start of the tool name:
- catalogue: listCourses, getCourse, getCourseLesson, listPublicCourseCategories, getPublicCourseCategory, listBlogPosts, getBlogPost. Needs the catalog:read scope.
- learner progress: getMyDashboard, getMyCourseProgress, getMyLessonProgress. Only the signed-in user's own data. Needs progress:read.
- content administration (admin...): banners, blog posts, courses, course versions, landing pages, categories, reviews, lessons, QnA items, attachments, leads, messages. Needs the admin scope and the matching staff Permission.
- finding content (adminSearch, adminFetch): full-text search of courses, lessons and posts in a locale, drafts, archived and unlisted pages included, and the read of one result by its id with its state. Needs the admin scope and mcp.read. There is no public search tool.
- access management (admin...): users, staff members, staff roles and their Permissions. Needs the admin scope and the mcp.manage Permission.
- audit (adminListAuditRecords): the trail of what staff and assistants did. Needs audit.read.

There are about 90 tools and you only see the ones your sign-in allows. If your client does not load tools on demand, do not guess names: call findTools with a word such as "banner" or "publish" and it answers with names, short descriptions and what each needs. findTools only looks tools up; to use one, call it by its name.

Admin tools take a required locale (en or ru), never assumed. Reading a record is a get tool and changing it is an update tool; publishing, unpublishing, archiving and marking a course's state are tools of their own, and an update never changes state. A refused call says why (HTTP 403 is a missing Permission or locale).

Lists come in pages (page, perPage). A result larger than the size limit is shortened and carries a truncation object: cuts says what was cut (path, kept, total), and for a list suggestedPage and suggestedPerPage continue reading without a gap. Ask for a narrower page or a single record instead of the whole.`

// toolNeeds is what a call to a tool is judged on, for showing a caller only the
// tools it may use. It is a convenience: the checkpoint judges every call again.
type toolNeeds struct {
	// scope is the OAuth Scope the connection must hold.
	scope string
	// permission is the operation's own Permission, when the tool mirrors an admin
	// operation, for findTools to say what a tool needs.
	permission string
	// capabilities are the capability Permissions (mcp.read, mcp.manage,
	// audit.read) the caller's Staff Role must hold.
	capabilities []string
}

// listedTool is a tool and what it needs, as the registry keeps it.
type listedTool struct {
	tool  *mcp.Tool
	needs toolNeeds
}

// note records a tool for the list filter and findTools. It is called where a
// tool is registered, so a tool the registry does not know cannot be offered.
func (s *Server) note(tool *mcp.Tool, needs toolNeeds) {
	s.registry = append(s.registry, listedTool{tool: tool, needs: needs})
}

// visibleTo is the registered tools the connection's token would not be refused
// at the door of, sorted by name so the order never depends on registration.
//
// A tool is shown when the token holds its Scope and the user holds each of its
// capabilities. When the capabilities cannot be read (none wired, or a failure
// of ours) the capability test is skipped rather than hiding everything: hiding
// is not authorization, so showing a tool that the call then refuses costs one
// refusal and nothing else.
func (s *Server) visibleTo(ctx context.Context, info *auth.TokenInfo) []listedTool {
	if info == nil {
		return nil
	}

	var held map[string]bool

	if s.deps.HeldCapabilities != nil && slices.Contains(info.Scopes, oauthserver.ScopeAdmin) {
		if userID, err := strconv.Atoi(info.UserID); err == nil && userID > 0 {
			names, err := s.deps.HeldCapabilities(ctx, userID)
			if err != nil {
				s.deps.Logger.ErrorContext(ctx, "read capabilities for the MCP tool list", slog.Any("err", err))
			} else {
				held = map[string]bool{}
				for _, name := range names {
					held[name] = true
				}
			}
		}
	}

	var shown []listedTool

	for _, t := range s.registry {
		if t.needs.scope != "" && !slices.Contains(info.Scopes, t.needs.scope) {
			continue
		}

		if held != nil && slices.ContainsFunc(t.needs.capabilities, func(c string) bool { return !held[c] }) {
			continue
		}

		shown = append(shown, t)
	}

	slices.SortFunc(shown, func(a, b listedTool) int { return strings.Compare(a.tool.Name, b.tool.Name) })

	return shown
}

// filterToolList narrows a `tools/list` answer to what visibleTo allows. The
// SDK builds the list from every registered tool; this is the one place it is
// made per caller.
func (s *Server) filterToolList(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		res, err := next(ctx, method, req)

		list, ok := res.(*mcp.ListToolsResult)
		if err != nil || method != "tools/list" || !ok {
			return res, err
		}

		var info *auth.TokenInfo
		if extra := req.GetExtra(); extra != nil {
			info = extra.TokenInfo
		}

		shown := map[string]bool{}
		for _, t := range s.visibleTo(ctx, info) {
			shown[t.tool.Name] = true
		}

		kept := make([]*mcp.Tool, 0, len(list.Tools))
		for _, tool := range list.Tools {
			if shown[tool.Name] {
				kept = append(kept, tool)
			}
		}

		list.Tools = kept

		return list, nil
	}
}

// markPrivate is SetCacheable: a tool list differs by caller, so a shared cache
// must not serve one caller's list to another. The time to live is left alone
// (zero: stale at once), so a client asks again and sees a change of grants.
func markPrivate(_ context.Context, _ mcp.Request, c *mcp.Cacheable) {
	c.CacheScope = "private"
}

// findToolsArgs are the arguments of findTools.
type findToolsArgs struct {
	Query  string `json:"query"`
	Detail string `json:"detail"`
	Limit  int    `json:"limit"`
}

// Detail levels of findTools.
const (
	detailName        = "name"
	detailDescription = "description"
	detailSchema      = "schema"

	defaultFindLimit = 20
	maxFindLimit     = 50
)

// findToolsOutput is the advertised result: small on purpose, the detail it
// carries is the caller's to ask for.
func findToolsOutput() map[string]any {
	return withTruncation(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"total": map[string]any{"type": "integer"},
			"tools": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":     "object",
					"required": []string{"name"},
					"properties": map[string]any{
						"name":        map[string]any{"type": "string"},
						"description": map[string]any{"type": "string"},
						"needs":       map[string]any{"type": "object"},
						"inputSchema": map[string]any{"type": "object"},
					},
				},
			},
		},
		"required": []string{"total", "tools"},
	})
}

// addFindTools registers the discovery tool. It only reads the registry: it
// holds no handler of any other tool, so there is nothing in it that could run
// one (a test pins that no tool executes another by name).
func (s *Server) addFindTools(server *mcp.Server) {
	input := object(map[string]any{
		"query": map[string]any{"type": "string", "description": "Words to look for in tool names and descriptions, all of which must match. Omit to list every tool you may call."},
		"detail": map[string]any{
			"type": "string", "enum": []string{detailName, detailDescription, detailSchema},
			"description": "name: names only. description: name, short description and what it needs (default). schema: also the full input schema.",
		},
		"limit": map[string]any{"type": "integer", "minimum": 1, "description": fmt.Sprintf("Most tools to return. Default %d, at most %d.", defaultFindLimit, maxFindLimit)},
	})
	tool := &mcp.Tool{
		Name:         findToolsName,
		Description:  "Find tools by a word, with what each needs. Looks up only; call a tool by its own name.",
		InputSchema:  input,
		OutputSchema: findToolsOutput(),
		Annotations:  hints(true, false, true),
	}
	server.AddTool(tool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info := req.GetExtra().TokenInfo
		if info == nil {
			return toolError("this connection is not signed in"), nil
		}

		userID, err := strconv.Atoi(info.UserID)
		if err != nil || userID <= 0 {
			return toolError("this connection is not signed in as a user"), nil
		}

		if _, denied := s.spend(req, userID, ClassCatalogue); denied != nil {
			return denied, nil
		}

		return s.findTools(ctx, req, info)
	})
	// No Scope: any signed-in caller may look, and what it finds is what it may use.
	s.note(tool, toolNeeds{})
}

func (s *Server) findTools(ctx context.Context, req *mcp.CallToolRequest, info *auth.TokenInfo) (*mcp.CallToolResult, error) {
	c := call{args: req.Params.Arguments}

	args := findToolsArgs{Detail: detailDescription}
	if err := c.decode(&args); err != nil {
		return toolError(err.Error()), nil
	}

	if !slices.Contains([]string{detailName, detailDescription, detailSchema}, args.Detail) {
		return toolError("detail must be name, description or schema"), nil
	}

	if args.Limit < 0 {
		return toolError("limit must be at least 1"), nil
	}

	limit := defaultFindLimit
	if args.Limit > 0 {
		limit = min(args.Limit, maxFindLimit)
	}

	words := strings.Fields(strings.ToLower(args.Query))

	var matched []listedTool

	for _, t := range s.visibleTo(ctx, info) {
		haystack := strings.ToLower(t.tool.Name + " " + t.tool.Description)
		if !slices.ContainsFunc(words, func(w string) bool { return !strings.Contains(haystack, w) }) {
			matched = append(matched, t)
		}
	}

	rows := []map[string]any{}

	for _, t := range matched[:min(len(matched), limit)] {
		row := map[string]any{"name": t.tool.Name}
		if args.Detail != detailName {
			row["description"] = t.tool.Description
			row["needs"] = t.needs.describe()
		}

		if args.Detail == detailSchema {
			row["inputSchema"] = t.tool.InputSchema
		}

		rows = append(rows, row)
	}

	text := fmt.Sprintf("%d of %d matching tools", len(rows), len(matched))

	return s.finish(findToolsName, result{text: text, data: map[string]any{"total": len(matched), "tools": rows}})
}

// describe is what a tool needs, as findTools says it.
func (n toolNeeds) describe() map[string]any {
	out := map[string]any{}
	if n.scope != "" {
		out["scope"] = n.scope
	}

	if n.permission != "" {
		out["permission"] = n.permission
	}

	if len(n.capabilities) > 0 {
		out["capabilities"] = n.capabilities
	}

	return out
}

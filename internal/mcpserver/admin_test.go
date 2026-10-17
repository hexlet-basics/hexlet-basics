package mcpserver_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/auditrecord"
	"hexletbasics/ent/banner"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/testsupport"
)

// The fixtures seed bob (user_staff_editor) as an editor, role staff_role_editor,
// who may read and create banners in ru and en; alice (user_admin_learner) is an
// administrator with no Staff Role.
var editor = fixtures.ID(fixtures.StaffRoleEditor)

// adminStaff is an administrator who is also a staff member, on the editor role
// (user_admin_staff): alice has no staff record, and must not (the cases lean on it);
// a record cannot be made for oneself.
var adminStaff = fixtures.ID(fixtures.UserAdminStaff)

// bannerInput is a valid body of adminUpdateBanner.
func bannerInput(body string) map[string]any {
	return map[string]any{
		"state": "draft", "background": "dark", "locale": "en", "body": body,
		"url": nil, "startsAt": nil, "finishesAt": nil,
	}
}

func call(t *testing.T, session *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: name, Arguments: args})
	require.NoError(t, err)
	return result
}

// admin connects an assistant for the user with the admin Scope.
func admin(t *testing.T, o *testsupport.OAuth, userID int) (*mcp.ClientSession, oauthserver.Token) {
	t.Helper()
	tokens := o.Connect(userID, oauthserver.ScopeAdmin)
	token, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
	return o.MCPClient(tokens.AccessToken), token
}

func firstBannerID(t *testing.T, o *testsupport.OAuth) int {
	t.Helper()
	return o.DB.Banner.Query().Order(banner.ByID()).FirstIDX(t.Context())
}

func assertRefused(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	require.True(t, result.IsError, "%+v", result.Content)
	assert.Contains(t, textOf(t, result), "HTTP 403")
}

func TestOnlyAllowlistedOperationsAreTools(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := admin(t, o, bob)

	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	var names []string
	for _, tool := range tools.Tools {
		if strings.HasPrefix(tool.Name, "admin") {
			names = append(names, tool.Name)
		}
	}
	slices.Sort(names)
	// The tools are exactly the allowlist: an excluded operation is in the
	// contract and is not a tool.
	assert.Equal(t, slices.Sorted(slices.Values(mcpserver.AllowlistedAdminOperations())), names)
	// Calling an excluded one anyway is refused by the endpoint itself: there is
	// no such tool.
	_, err = session.CallTool(t.Context(), &mcp.CallToolParams{
		Name: "adminGetMyCapabilities", Arguments: map[string]any{"locale": "en"},
	})
	require.Error(t, err)
	assert.Equal(t, 5, o.DB.Banner.Query().CountX(t.Context()))
}

func TestAdminToolSchemasComeFromTheContract(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := admin(t, o, bob)
	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	byName := map[string]*mcp.Tool{}
	for _, tool := range tools.Tools {
		byName[tool.Name] = tool
	}
	schemaOf := func(v any) map[string]any {
		raw, err := json.Marshal(v)
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal(raw, &m))
		return m
	}

	list := byName["adminListBanners"]
	assert.True(t, list.Annotations.ReadOnlyHint)
	input := schemaOf(list.InputSchema)
	assert.Equal(t, []any{"locale"}, input["required"], "the locale is required, the paging is not")
	properties := input["properties"].(map[string]any)
	assert.Equal(t, []any{"en", "ru"}, properties["locale"].(map[string]any)["enum"])
	assert.Contains(t, properties, "page")
	assert.Contains(t, properties, "perPage")
	assert.NotContains(t, properties, "body")

	update := byName["adminUpdateBanner"]
	assert.False(t, update.Annotations.ReadOnlyHint)
	assert.False(t, *update.Annotations.DestructiveHint)
	input = schemaOf(update.InputSchema)
	assert.ElementsMatch(t, []any{"locale", "id", "body"}, input["required"])
	body := input["properties"].(map[string]any)["body"].(map[string]any)
	assert.Equal(t, "object", body["type"])
	assert.Contains(t, body["properties"], "locale", "the banner's own locale sits in the body, apart from the call's")
	assert.NotContains(t, body, "$ref")
	assert.Equal(t, "object", schemaOf(update.OutputSchema)["type"])
}

func TestRoleWithTheGrantsListsBanners(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, bob)

	result := call(t, session, "adminListBanners", map[string]any{"locale": "en", "perPage": 2})

	require.False(t, result.IsError, "%+v", result.Content)
	page := result.StructuredContent.(map[string]any)
	assert.EqualValues(t, 5, page["total"])
	assert.Len(t, page["items"], 2, "the query parameter reached the operation")

	// The advertised output schema accepts what came back.
	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)
	for _, tool := range tools.Tools {
		if tool.Name != "adminListBanners" {
			continue
		}
		raw, err := json.Marshal(tool.OutputSchema)
		require.NoError(t, err)
		var schema jsonschema.Schema
		require.NoError(t, json.Unmarshal(raw, &schema))
		resolved, err := schema.Resolve(nil)
		require.NoError(t, err)
		assert.NoError(t, resolved.Validate(page))
	}

	// A read of banners is not recorded: it is not personal data.
	assert.Empty(t, o.Records("AdminListBanners"))
}

func TestWriteToolChangesTheBannerAndLeavesAnAuditRecordOfTheMCPChannel(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "banners.update", "mcp.write")
	session, token := admin(t, o, bob)
	id := firstBannerID(t, o)

	result := call(t, session, "adminUpdateBanner", map[string]any{
		"locale": "ru", "id": id, "body": bannerInput("Changed by an assistant"),
	})

	require.False(t, result.IsError, "%+v", result.Content)
	assert.Equal(t, "Changed by an assistant", result.StructuredContent.(map[string]any)["body"])
	assert.Equal(t, "Changed by an assistant", o.DB.Banner.GetX(t.Context(), id).Body)

	records := o.Records("AdminUpdateBanner")
	require.Len(t, records, 1)
	record := records[0]
	assert.Equal(t, audit.OutcomeSucceeded, audit.Outcome(record.Outcome))
	assert.Equal(t, "mcp", record.Channel)
	assert.Equal(t, bob, *record.ActorUserID)
	assert.Equal(t, "ru", *record.Locale, "the locale of the call is the argument")
	assert.Equal(t, strconv.Itoa(id), *record.EntityID)
	// The Connected App is the oauth_clients row the token was issued to, and the
	// token is its own id (not the bearer).
	app := o.DB.OAuthClient.Query().OnlyX(t.Context())
	assert.Equal(t, app.ID, *record.ConnectedAppID)
	assert.Equal(t, token.ID, *record.TokenID)
	assert.Equal(t, token.ConnectedAppID, *record.ConnectedAppID)
	assert.Contains(t, string(record.Changes), "Changed by an assistant")
}

func TestWriteIsRefusedWithOnlyMCPReadAndRecordedAsForbidden(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "banners.update", "mcp.read")
	session, token := admin(t, o, bob)
	id := firstBannerID(t, o)

	result := call(t, session, "adminUpdateBanner", map[string]any{
		"locale": "ru", "id": id, "body": bannerInput("Must not land"),
	})

	assertRefused(t, result)
	assert.NotEqual(t, "Must not land", o.DB.Banner.GetX(t.Context(), id).Body)
	records := o.Records("AdminUpdateBanner")
	require.Len(t, records, 1)
	assert.Equal(t, "forbidden", records[0].Outcome)
	assert.Equal(t, "mcp", records[0].Channel)
	assert.Equal(t, token.ID, *records[0].TokenID)
	assert.Equal(t, token.ConnectedAppID, *records[0].ConnectedAppID)
	assert.Empty(t, records[0].Changes, "a refusal records no values")
}

func TestOperationPermissionIsRequiredAlongsideTheMCPCapability(t *testing.T) {
	o := testsupport.NewOAuth(t)
	// The role may use MCP but was never given banners.update.
	o.Grant(editor, "mcp.read", "mcp.write")
	session, _ := admin(t, o, bob)

	result := call(t, session, "adminUpdateBanner", map[string]any{
		"locale": "ru", "id": firstBannerID(t, o), "body": bannerInput("Nope"),
	})

	assertRefused(t, result)
}

func TestReadIsRefusedWithoutMCPRead(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session, _ := admin(t, o, bob) // banners.read is held, mcp.read is not

	result := call(t, session, "adminListBanners", map[string]any{"locale": "en"})

	assertRefused(t, result)
	records := o.Records("AdminListBanners")
	require.Len(t, records, 1)
	assert.Equal(t, "forbidden", records[0].Outcome)
}

func TestAnAdministratorDoesNotGetMCPImplicitly(t *testing.T) {
	o := testsupport.NewOAuth(t)

	t.Run("no staff record", func(t *testing.T) {
		session, _ := admin(t, o, alice)
		assertRefused(t, call(t, session, "adminListBanners", map[string]any{"locale": "en"}))
	})

	// An administrator who is staff, on a role that holds mcp.read but not mcp.write.
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, adminStaff)

	assert.False(t, call(t, session, "adminListBanners", map[string]any{"locale": "en"}).IsError, "read is granted")
	write := call(t, session, "adminUpdateBanner", map[string]any{
		"locale": "en", "id": firstBannerID(t, o), "body": bannerInput("Administrator without mcp.write"),
	})
	assertRefused(t, write)
	assert.NotEqual(t, "Administrator without mcp.write", o.DB.Banner.Query().Order(banner.ByID()).FirstX(t.Context()).Body)

	o.Grant(editor, "mcp.write")
	write = call(t, session, "adminUpdateBanner", map[string]any{
		"locale": "en", "id": firstBannerID(t, o), "body": bannerInput("Administrator with mcp.write"),
	})
	assert.False(t, write.IsError, "%+v", write.Content)
}

func TestTokenWithoutTheAdminScopeIsRefused(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	tokens := o.Connect(bob, oauthserver.ScopeCatalogRead)
	session := o.MCPClient(tokens.AccessToken)

	result := call(t, session, "adminListBanners", map[string]any{"locale": "en"})

	assertRefused(t, result)
	records := o.Records("AdminListBanners")
	require.Len(t, records, 1)
	assert.Equal(t, "forbidden", records[0].Outcome)
	assert.Equal(t, "mcp", records[0].Channel)
}

func TestTakingAPermissionAwayTakesEffectOnTheNextCall(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, bob)
	args := map[string]any{"locale": "en"}

	require.False(t, call(t, session, "adminListBanners", args).IsError)

	o.RevokeGrant(editor, "mcp.read")
	assertRefused(t, call(t, session, "adminListBanners", args))

	o.Grant(editor, "mcp.read")
	o.RevokeGrant(editor, "banners.read")
	assertRefused(t, call(t, session, "adminListBanners", args))
}

func TestLocaleIsARequiredArgumentThatDecidesTheCheck(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	// Bob works in ru and en; narrow him to ru.
	o.DB.StaffMember.Update().SetAllowedLocales([]string{"ru"}).ExecX(t.Context())
	tokens := o.Connect(bob, oauthserver.ScopeAdmin)
	// The connection's own Accept-Language says en; the argument is what counts.
	session := o.MCPClient(tokens.AccessToken, func(r *http.Request) { r.Header.Set("Accept-Language", "en") })

	for name, args := range map[string]map[string]any{
		"missing":     {},
		"not a value": {"locale": 7},
		"unknown":     {"locale": "xx"},
		"unsupported": {"locale": "es"},
	} {
		result := call(t, session, "adminListBanners", args)
		assert.True(t, result.IsError, name)
		assert.Contains(t, textOf(t, result), "locale", name)
	}
	assert.Empty(t, o.Records("AdminListBanners"), "an invalid call never reaches the checkpoint")

	assertRefused(t, call(t, session, "adminListBanners", map[string]any{"locale": "en"}))
	allowed := call(t, session, "adminListBanners", map[string]any{"locale": "ru"})
	assert.False(t, allowed.IsError, "%+v", allowed.Content)

	records := o.Records("AdminListBanners")
	require.Len(t, records, 1)
	assert.Equal(t, "en", *records[0].Locale, "the refusal was judged on the argument, not on Accept-Language")
}

func TestArgumentsTheOperationDoesNotDeclareAreRefused(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "mcp.read")
	session, _ := admin(t, o, bob)

	extra := call(t, session, "adminListBanners", map[string]any{"locale": "en", "color": "red"})
	assert.True(t, extra.IsError)
	assert.Contains(t, textOf(t, extra), "unknown arguments: color")

	missing := call(t, session, "adminUpdateBanner", map[string]any{"locale": "en", "body": bannerInput("x")})
	assert.True(t, missing.IsError)
	assert.Contains(t, textOf(t, missing), "id argument is required")
}

func TestValidationFailureIsAToolErrorAndIsRecordedAsFailed(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "banners.update", "mcp.write")
	session, _ := admin(t, o, bob)

	invalid := bannerInput("")
	result := call(t, session, "adminUpdateBanner", map[string]any{"locale": "ru", "id": firstBannerID(t, o), "body": invalid})

	require.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "HTTP 400")
	records := o.Records("AdminUpdateBanner")
	// The generated server refuses an empty body while decoding, before the
	// handler and before the audit middleware, like over HTTP.
	for _, r := range records {
		assert.NotEqual(t, "succeeded", r.Outcome)
	}
}

func TestMissingBannerIsTheSameNotFoundHTTPAnswers(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "banners.update", "mcp.write")
	session, _ := admin(t, o, bob)

	result := call(t, session, "adminUpdateBanner", map[string]any{"locale": "ru", "id": 999999, "body": bannerInput("x")})

	require.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "HTTP 404")
	records := o.Records("AdminUpdateBanner")
	require.Len(t, records, 1)
	assert.Equal(t, "failed", records[0].Outcome)
	assert.Equal(t, "mcp", records[0].Channel)
}

// HTTP and MCP run one pipeline, so the same operation by the same user answers
// and records the same, apart from the channel.
func TestHTTPAndMCPBehaveTheSameForTheSameOperation(t *testing.T) {
	o := testsupport.NewOAuth(t)
	cookie := testsupport.SessionFor(t, o.DB, bob)
	session, _ := admin(t, o, bob)
	id := firstBannerID(t, o)
	path := "/api/admin/banners/" + strconv.Itoa(id)

	overHTTP := func(input map[string]any) (int, string) {
		raw, err := json.Marshal(input)
		require.NoError(t, err)
		res := o.Do(http.MethodPut, path, nil, func(r *http.Request) {
			r.Body = io.NopCloser(bytes.NewReader(raw))
			r.ContentLength = int64(len(raw))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Accept-Language", "ru")
			r.AddCookie(cookie.Cookie)
		})
		defer func() { _ = res.Body.Close() }()
		return res.StatusCode, body(t, res)
	}
	overMCP := func(input map[string]any) *mcp.CallToolResult {
		return call(t, session, "adminUpdateBanner", map[string]any{"locale": "ru", "id": id, "body": input})
	}

	// 1. Refused: banners.update is not held. The same status and the same
	// problem document, and a forbidden record from each channel.
	status, httpBody := overHTTP(bannerInput("refused"))
	o.Grant(editor, "mcp.write")
	refused := overMCP(bannerInput("refused"))
	assert.Equal(t, http.StatusForbidden, status)
	require.True(t, refused.IsError)
	assert.Equal(t, "HTTP 403: "+httpBody, textOf(t, refused))

	// 2. Invalid input: the same 400 and problem document.
	o.Grant(editor, "banners.update")
	status, httpBody = overHTTP(bannerInput(""))
	invalid := overMCP(bannerInput(""))
	assert.Equal(t, http.StatusBadRequest, status)
	require.True(t, invalid.IsError)
	assert.Equal(t, "HTTP 400: "+httpBody, textOf(t, invalid))

	// 3. Allowed: the same banner comes back and the same record is written.
	status, httpBody = overHTTP(bannerInput("same"))
	require.Equal(t, http.StatusOK, status, httpBody)
	viaMCP := overMCP(bannerInput("same"))
	require.False(t, viaMCP.IsError, "%+v", viaMCP.Content)
	var viaHTTP map[string]any
	require.NoError(t, json.Unmarshal([]byte(httpBody), &viaHTTP))
	mcpBanner := viaMCP.StructuredContent.(map[string]any)
	for _, k := range []string{"id", "body", "state", "background", "locale", "url"} {
		assert.Equal(t, viaHTTP[k], mcpBanner[k], k)
	}

	records := o.Records("AdminUpdateBanner")
	byOutcome := map[string]map[string]*string{}
	for _, r := range records {
		if byOutcome[r.Outcome] == nil {
			byOutcome[r.Outcome] = map[string]*string{}
		}
		byOutcome[r.Outcome][r.Channel] = r.Locale
	}
	for _, outcome := range []string{"forbidden", "succeeded"} {
		assert.Contains(t, byOutcome[outcome], "web", outcome)
		assert.Contains(t, byOutcome[outcome], "mcp", outcome)
		assert.Equal(t, "ru", *byOutcome[outcome]["web"], outcome)
		assert.Equal(t, "ru", *byOutcome[outcome]["mcp"], outcome)
	}
	// A refused or invalid attempt changed nothing on either channel.
	assert.Equal(t, 1, o.DB.AuditRecord.Query().
		Where(auditrecord.Operation("AdminUpdateBanner"), auditrecord.Channel("mcp"), auditrecord.Outcome("succeeded")).CountX(t.Context()))
}

// Every tool the registry offers is refused for a caller whose role holds
// nothing, whatever it asks: no registered tool reaches an operation without the
// checkpoint. The arguments are generated from the tool's own input schema, so a
// tool added to the allowlist is covered with no edit here.
func TestNoRegisteredAdminToolRunsWithoutTheCheck(t *testing.T) {
	// One call per tool must reach the checkpoint, so the budgets are not the
	// subject here.
	rates := mcpserver.Rates{Catalogue: 1000, Personal: 1000, Write: 1000}
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{User: rates, Client: rates}))
	// A role with nothing: bob's grants are removed, and no mcp.* is held.
	o.DB.StaffRoleGrant.Delete().ExecX(t.Context())
	session, _ := admin(t, o, bob)
	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	checked := 0
	for _, tool := range tools.Tools {
		if !strings.HasPrefix(tool.Name, "admin") {
			continue
		}
		checked++
		result := call(t, session, tool.Name, sampleArguments(t, tool))
		assertRefused(t, result)
	}
	assert.Positive(t, checked)
	assert.Zero(t, o.DB.AuditRecord.Query().Where(auditrecord.Outcome("succeeded")).CountX(t.Context()))
	assert.Equal(t, 5, o.DB.Banner.Query().CountX(t.Context()))
}

// sampleArguments fills the required properties of a tool's input schema with
// values of the right type; what they say does not matter, because the check
// refuses before the request is read.
func sampleArguments(t *testing.T, tool *mcp.Tool) map[string]any {
	t.Helper()
	raw, err := json.Marshal(tool.InputSchema)
	require.NoError(t, err)
	var schema struct {
		Required   []string                  `json:"required"`
		Properties map[string]map[string]any `json:"properties"`
	}
	require.NoError(t, json.Unmarshal(raw, &schema))
	args := map[string]any{}
	for _, name := range schema.Required {
		switch property := schema.Properties[name]; {
		case name == "locale":
			args[name] = "en"
		case name == "body" && tool.Name == "adminUploadAttachment":
			args[name] = map[string]any{"fileName": "x.png", "contentBase64": "eA=="}
		case property["type"] == "integer":
			args[name] = 1
		case property["type"] == "string":
			args[name] = "x"
		default:
			args[name] = map[string]any{}
		}
	}
	return args
}

// TestSessionPlaceholderOpensNothingOverHTTP pins the property the in-process
// request rests on (ADR-0021 addendum, #854): the session cookie a tool call's
// request carries is a placeholder that only satisfies the generated server's
// security code, and it is worth nothing to a client. Sent over HTTP, where the
// context holds no Principal, the session lookup finds no session under it and
// the request is refused before anything runs, for a read and for a write. It
// fails if the MCP branch of the security handler ever starts to trust a cookie
// value instead of a Principal in the context.
func TestSessionPlaceholderOpensNothingOverHTTP(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor, "banners.update", "banners.read", "mcp.write", "mcp.read")
	bannerID := fixtures.ID(fixtures.BannerRuPublished)
	placeholder := &http.Cookie{Name: mcpserver.SessionCookie, Value: mcpserver.CookiePlaceholder}

	read := o.Do(http.MethodGet, "/api/admin/banners", nil, func(r *http.Request) {
		r.AddCookie(placeholder)
		r.Header.Set("Accept-Language", "en")
	})
	defer func() { _ = read.Body.Close() }()
	assert.Equal(t, http.StatusUnauthorized, read.StatusCode)

	// What the in-process request carries, sent from outside.
	write := o.Do(http.MethodPut, "/api/admin/banners/"+strconv.Itoa(bannerID), nil, func(r *http.Request) {
		r.AddCookie(placeholder)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Language", "en")
		r.Body = io.NopCloser(strings.NewReader(`{}`))
	})
	defer func() { _ = write.Body.Close() }()
	assert.Equal(t, http.StatusUnauthorized, write.StatusCode)
	assert.NotEqual(t, "Nope", o.DB.Banner.GetX(t.Context(), bannerID).Body)
}

// TestThePlaceholderIsTheWebsitesSessionCookie keeps the dispatcher's cookie name
// equal to the website's session cookie: under any other name the generated
// server would refuse every tool call before the security handler ran.
func TestThePlaceholderIsTheWebsitesSessionCookie(t *testing.T) {
	assert.Equal(t, sessions.CookieName, mcpserver.SessionCookie)
}

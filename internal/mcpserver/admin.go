package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/internal/audit"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/permissions"
)

// adminOperations is the allowlist of admin operations exposed as tools: the
// operationIds, nothing else (ADR-0021). An operation added to the contract is
// not an agent capability until its operationId is listed here, so exposure is a
// reviewed line of code and never a side effect of a contract change. New fails
// for a name that is not an admin operation of the contract.
//
// Each operation still passes the one checkpoint when it is called (see
// callAdminTool); the list says what may be offered, not what may be done.
var adminOperations = []string{
	"adminArchiveBanner",
	"adminArchiveBlogPost",
	"adminArchiveCourseLandingPage",
	"adminArchiveReview",
	"adminCreateBanner",
	"adminCreateBlogPost",
	"adminCreateCategoryQnaItem",
	"adminCreateCourse",
	"adminCreateCourseCategory",
	"adminCreateCourseLandingPage",
	"adminCreateCourseVersion",
	"adminCreateLandingPageQnaItem",
	"adminCreateReview",
	"adminCreateRole",
	"adminCreateStaffMember",
	"adminCreateUser",
	"adminDeleteBanner",
	"adminDeleteBlogPost",
	"adminDeleteCategoryQnaItem",
	"adminDeleteCourseCategory",
	"adminDeleteCourseLandingPage",
	"adminDeleteLandingPageQnaItem",
	"adminDeleteReview",
	"adminDeleteRole",
	"adminDeleteStaffMember",
	"adminDeleteUser",
	"adminFetch",
	"adminGetBanner",
	"adminGetBlogPost",
	"adminGetCourse",
	"adminGetCourseCategory",
	"adminGetCourseLandingPage",
	"adminGetManagementUser",
	"adminGetReview",
	"adminGetRole",
	"adminGetRolePermissions",
	"adminGetStaffMember",
	"adminGetUser",
	"adminListAuditRecords",
	"adminListBanners",
	"adminListBlogPosts",
	"adminListCategoryQnaItems",
	"adminListCourseCategories",
	"adminListCourseLandingPages",
	"adminListCourseLessonReviews",
	"adminListCourseLessons",
	"adminListCourseOptions",
	"adminListCourses",
	"adminListLandingPageQnaItems",
	"adminListLeads",
	"adminListLessonProgress",
	"adminListManagementUsers",
	"adminListMessages",
	"adminListPermissions",
	"adminListReviews",
	"adminListRoles",
	"adminListStaffMembers",
	"adminListUsers",
	"adminMarkCourseCompleted",
	"adminMarkCourseDraft",
	"adminMarkCourseInDevelopment",
	"adminPublishBanner",
	"adminPublishBlogPost",
	"adminPublishCourseLandingPage",
	"adminPublishReview",
	"adminReviewCourse",
	"adminReviewCourseLesson",
	"adminSearch",
	"adminSearchUsers",
	"adminSetBlogPostRelatedCourses",
	"adminSuggestBlogPostRelatedCourses",
	"adminUnpublishBanner",
	"adminUnpublishBlogPost",
	"adminUnpublishCourseLandingPage",
	"adminUnpublishReview",
	"adminUpdateBanner",
	"adminUpdateBlogPost",
	"adminUpdateCategoryQnaItem",
	"adminUpdateCourse",
	"adminUpdateCourseCategory",
	"adminUpdateCourseLandingPage",
	"adminUpdateLandingPageQnaItem",
	"adminUpdateManagementUser",
	"adminUpdateReview",
	"adminUpdateRole",
	"adminUpdateRolePermissions",
	"adminUpdateStaffMember",
	"adminUpdateUser",
	"adminUploadAttachment",
}

// excludedAdminOperations are the admin operations (the ones that require the
// AdminSession scheme) that are deliberately not tools, each with its reason.
// A test requires every admin operation of the contract to be in
// adminOperations or here, so an operation added to the contract later cannot
// be forgotten: it is either offered to agents or excluded on the record.
var excludedAdminOperations = map[string]string{
	"adminGetMyCapabilities": "tells the admin menu which entries to show for the signed-in user; an agent learns what it may do from the tools' answers (a refusal names the status), and the call needs no Permission, so it is not an operation an agent acts through",
}

// adminLocales are the locales staff work in. The tool's locale argument is
// validated against them, because the website's locale matcher silently falls
// back to English for a language it does not know, which would judge a call in a
// locale the caller never named.
var adminLocales = []string{"en", "ru"}

const (
	// localeArgument and bodyArgument are reserved names of a tool's input. The
	// request body is nested under one key instead of being merged with the
	// parameters, so that its fields can never collide with them (a banner has a
	// `locale` of its own, which is not the locale the call is made in).
	localeArgument = "locale"
	bodyArgument   = "body"

	// sessionCookie is the name of the website's session cookie (auth.CookieName;
	// a test pins the two together). The package does not import the session
	// package, which would put the session store within reach of the tools.
	sessionCookie = "session"

	// cookiePlaceholder fills the session cookie the generated server's security
	// code insists on before it calls the security handler. The handler ignores
	// its value when the context carries a Principal, and when it does not (the
	// same request sent over HTTP) no session exists under it, so the session
	// lookup answers 401: the placeholder opens nothing.
	cookiePlaceholder = "mcp-in-process"
)

// addAdminTools registers a tool for each allowlisted operation.
func (s *Server) addAdminTools(server *mcp.Server) error {
	for _, id := range adminOperations {
		op, err := readOperation(id)
		if err != nil {
			return errs.With("tool", id).Wrapf(err, "admin tool")
		}
		tool, err := adminTool(op)
		if err != nil {
			return errs.With("tool", id).Wrapf(err, "admin tool")
		}
		s.note(tool, adminNeeds(op))
		server.AddTool(tool, s.adminHandler(op))
	}
	return nil
}

// adminNeeds is what the checkpoint will ask of a caller of the operation, as
// the list filter and findTools show it: the admin Scope, the operation's
// Permission, and the capability Permissions: the ones every MCP call needs
// (mcp.read or mcp.write by kind, mcp.manage where declared) and the operation's
// own Permission when it is one (the audit trail's audit.read).
func adminNeeds(op *operation) toolNeeds {
	needs := toolNeeds{scope: oauthserver.ScopeAdmin}
	if requirement, ok := permissions.For(op.ogenName()); ok {
		needs.permission = requirement.Permission
		if permissions.IsCapability(requirement.Permission) {
			needs.capabilities = append(needs.capabilities, requirement.Permission)
		}
	}
	for _, capability := range permissions.MCPCapabilities(op.ogenName()) {
		if !slices.Contains(needs.capabilities, capability) {
			needs.capabilities = append(needs.capabilities, capability)
		}
	}
	return needs
}

// adminTool describes an operation as a tool: named after its operationId, with
// input and output schemas derived from the contract.
func adminTool(op *operation) (*mcp.Tool, error) {
	if !op.Admin {
		return nil, errors.New("is not an admin operation (it does not require the AdminSession scheme), so the checkpoint would not run for it")
	}

	properties := map[string]any{
		localeArgument: map[string]any{
			"type":        "string",
			"enum":        adminLocales,
			"description": "The locale this call is made in: the one your Staff Role must allow. Required, never inferred.",
		},
	}
	required := []string{localeArgument}
	for _, p := range op.Params {
		if p.Name == localeArgument || p.Name == bodyArgument {
			return nil, errs.With("parameter", p.Name).Errorf("parameter collides with a reserved argument")
		}
		schema := map[string]any{}
		for k, v := range p.Schema {
			schema[k] = v
		}
		if p.Description != "" {
			schema["description"] = p.Description
		}
		properties[p.Name] = schema
		if p.Required {
			required = append(required, p.Name)
		}
	}
	if op.Body != nil {
		properties[bodyArgument] = op.Body
		if op.BodyRequired {
			required = append(required, bodyArgument)
		}
	}
	if op.Multipart != nil {
		properties[bodyArgument] = uploadSchema(op.Multipart)
		if op.BodyRequired {
			required = append(required, bodyArgument)
		}
	}
	input := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"required":             required,
		"additionalProperties": false,
	}

	tool := &mcp.Tool{
		Name:        op.ID,
		Description: describe(op),
		InputSchema: trimInput(input).(map[string]any),
		Annotations: adminAnnotations(op),
	}
	if op.Response != nil {
		tool.OutputSchema = withTruncation(outputShape(wrapObject(op.Response)))
	}
	return tool, nil
}

// describe is the tool's description: the contract's own text for the
// operation, or, for one that has none, its operationId read as words ("Admin
// list banners" becomes "List banners"). It ends with what the call needs, so
// an agent learns why a refusal came without making the call: the operation's
// Permission and, over MCP, mcp.manage where the contract declares it.
func describe(op *operation) string {
	// Only the contract's first sentence: the rest says how errors surface, which
	// a refusal tells an agent when it matters, and clients cut a description at a
	// fixed length from its start.
	text := firstSentence(op.Description, maxToolDescription)
	if text == "" {
		words := camelWords(strings.TrimPrefix(op.ID, "admin"))
		text = strings.ToUpper(words[:1]) + words[1:] + "."
	}
	needs := []string{}
	if requirement, ok := permissions.For(op.ogenName()); ok {
		needs = append(needs, requirement.Permission)
	}
	for _, capability := range permissions.MCPCapabilities(op.ogenName()) {
		if capability == "mcp.manage" {
			needs = append(needs, capability)
		}
	}
	if len(needs) > 0 {
		text += " Needs " + strings.Join(needs, " and ") + "."
	}
	return text
}

// camelWords splits "ListBlogPosts" into "list blog posts".
func camelWords(name string) string {
	var out strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out.WriteByte(' ')
		}
		out.WriteRune(r)
	}
	return strings.ToLower(out.String())
}

// wrapObject is the schema of structured content, which MCP requires to be an
// object: an object schema stands as it is, anything else sits under `result`.
func wrapObject(schema map[string]any) map[string]any {
	if schema["type"] == "object" {
		return schema
	}
	return map[string]any{
		"type":                 "object",
		"properties":           map[string]any{"result": schema},
		"required":             []string{"result"},
		"additionalProperties": false,
	}
}

// annotations are the hints an agent shows its user: reads are read-only, a
// delete is destructive (ADR-0021). The kind comes from the HTTP method, which
// is what the contract's own write classification falls back to as well.
func adminAnnotations(op *operation) *mcp.ToolAnnotations {
	if op.Method == http.MethodGet {
		return hints(true, false, true)
	}
	return hints(false, op.Method == http.MethodDelete, op.Method == http.MethodPut || op.Method == http.MethodDelete)
}

// adminHandler is a tool call: the arguments become an HTTP request that is
// served by the generated api.Server, in process.
//
// This is the whole of ADR-0021's "one checkpoint". The handler never reaches
// an operation's code. The only thing it holds is an http.Handler, the same
// generated server (with its security handler, validation, audit middleware and
// error mapping) the website calls, so the Permission, the capability, the
// transaction, the Audit Record and the error a refusal becomes are the code
// the HTTP call runs, not a copy of it. What the call adds is the Principal in
// the context, which the security handler honours in place of a session cookie.
func (s *Server) adminHandler(op *operation) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info := req.GetExtra().TokenInfo
		if info == nil {
			return toolError("this connection is not signed in"), nil
		}
		userID, err := strconv.Atoi(info.UserID)
		if err != nil || userID <= 0 {
			return toolError("this connection is not signed in as a user"), nil
		}
		// The rate limit every tool shares. An admin write is ClassWrite, the
		// strictest; an admin read is ClassPersonal (stricter than the public
		// catalogue, since it can show staff-only and personal data). The Scope
		// is not checked here but by the security handler, so a refusal is
		// recorded.
		class := ClassPersonal
		if permissions.IsAdminWrite(op.ogenName()) {
			class = ClassWrite
		}
		if _, denied := s.spend(req, userID, class); denied != nil {
			return denied, nil
		}
		ctx, cancel := context.WithTimeout(ctx, s.limits.ToolTimeout)
		defer cancel()

		args, problem := parseArguments(req.Params.Arguments)
		if problem != "" {
			return toolError(problem), nil
		}
		locale, problem := localeOf(args)
		if problem != "" {
			return toolError(problem), nil
		}
		principal := Principal{
			UserID: userID,
			Scopes: info.Scopes,
			Locale: locale,
		}
		callCtx := WithPrincipal(audit.WithSource(ctx, sourceOf(info)), principal)
		httpReq, problem := buildRequest(callCtx, op, args, locale, req.GetExtra().Header)
		if problem != "" {
			return toolError(problem), nil
		}

		rec := httptest.NewRecorder()
		s.deps.API.ServeHTTP(rec, httpReq)
		return s.answer(op, rec)
	}
}

// sourceOf names the channel, Connected App and token an Audit Record of this
// call carries. The verifier put them in the token info: the oauth_clients row
// id of the client the token was issued to, and the access token's own id (not
// the bearer, which is that id encrypted).
func sourceOf(info *auth.TokenInfo) audit.Source {
	source := audit.Source{Channel: audit.ChannelMCP}
	if id, ok := info.Extra[extraConnectedAppID].(int); ok && id > 0 {
		source.ConnectedAppID = &id
	}
	if id, ok := info.Extra[extraTokenID].(string); ok && id != "" {
		source.TokenID = &id
	}
	return source
}

// parseArguments reads the tool's arguments as an object whose values stay raw
// JSON, so a number reaches the request exactly as the agent wrote it.
func parseArguments(raw json.RawMessage) (map[string]json.RawMessage, string) {
	args := map[string]json.RawMessage{}
	if len(raw) == 0 {
		return args, ""
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, "the arguments are not a JSON object"
	}
	return args, ""
}

// localeOf is the required locale argument: present, a string and one of the
// staff locales. It is never defaulted from the request's Accept-Language.
func localeOf(args map[string]json.RawMessage) (string, string) {
	raw, ok := args[localeArgument]
	if !ok {
		return "", "the locale argument is required"
	}
	var locale string
	if err := json.Unmarshal(raw, &locale); err != nil || !slices.Contains(adminLocales, locale) {
		return "", fmt.Sprintf("the locale argument must be one of %s", strings.Join(adminLocales, ", "))
	}
	return locale, ""
}

// buildRequest turns arguments into the HTTP request the operation takes. A
// name the operation does not declare, or a required one that is missing, is
// the caller's mistake and is reported as a tool error; whatever else is wrong
// (a value of the wrong type, a body that breaks the schema) the generated
// server's own validation answers, as it does over HTTP.
func buildRequest(ctx context.Context, op *operation, args map[string]json.RawMessage, locale string, forwarded http.Header) (*http.Request, string) {
	known := map[string]bool{localeArgument: true}
	path := op.Path
	query := url.Values{}
	for _, p := range op.Params {
		known[p.Name] = true
		raw, ok := args[p.Name]
		if !ok {
			if p.Required {
				return nil, fmt.Sprintf("the %s argument is required", p.Name)
			}
			continue
		}
		values, err := scalars(raw)
		if err != nil {
			return nil, fmt.Sprintf("the %s argument must be a value or a list of values", p.Name)
		}
		switch p.In {
		case "path":
			path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(values[0]))
		default:
			for _, v := range values {
				query.Add(p.Name, v)
			}
		}
	}
	var body []byte
	contentType := "application/json"
	if op.Multipart != nil {
		known[bodyArgument] = true
		if raw, ok := args[bodyArgument]; ok {
			var problem string
			if body, contentType, problem = multipartBody(op.Multipart, raw); problem != "" {
				return nil, problem
			}
		} else if op.BodyRequired {
			return nil, "the body argument is required"
		}
	}
	if op.Body != nil {
		known[bodyArgument] = true
		if raw, ok := args[bodyArgument]; ok {
			body = raw
		} else if op.BodyRequired {
			return nil, "the body argument is required"
		}
	}
	names := make([]string, 0, len(args))
	for name := range args {
		if !known[name] {
			names = append(names, name)
		}
	}
	if len(names) > 0 {
		sort.Strings(names)
		return nil, fmt.Sprintf("unknown arguments: %s", strings.Join(names, ", "))
	}

	target := path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, op.Method, target, bytes.NewReader(body))
	if err != nil {
		return nil, "the arguments do not make a request"
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	// The locale the call is made in is the one the HTTP check reads: the
	// request's Accept-Language, here set from the argument and from nothing else.
	req.Header.Set("Accept-Language", locale)
	req.AddCookie(&http.Cookie{Name: sessionCookie, Value: cookiePlaceholder})
	// What an Audit Record keeps of the caller: the MCP client's user agent, and
	// the address and request id the proxy and client sent with the /mcp request.
	for _, h := range []string{"User-Agent", "X-Forwarded-For", "X-Real-IP", "X-Request-ID"} {
		if v := forwarded.Get(h); v != "" {
			req.Header.Set(h, v)
		}
	}
	return req, ""
}

// scalars reads a JSON value as the strings a path or query parameter carries:
// a string, number or boolean is one value, an array of them is several.
func scalars(raw json.RawMessage) ([]string, error) {
	var one any
	if err := json.Unmarshal(raw, &one); err != nil {
		return nil, err
	}
	items, isList := one.([]any)
	if !isList {
		items = []any{one}
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		switch v := item.(type) {
		case string:
			out = append(out, v)
		case float64, bool:
			// Re-encode, not format: 3 must stay 3, not 3.000000.
			b, _ := json.Marshal(v)
			out = append(out, string(b))
		default:
			return nil, errors.New("not a scalar")
		}
	}
	if len(out) == 0 {
		return nil, errors.New("empty")
	}
	return out, nil
}

// answer maps the recorded HTTP response to a tool result: a 2xx is the
// structured content, anything else is a tool error that carries the status and
// the problem document exactly as an HTTP caller would have received them.
func (s *Server) answer(op *operation, rec *httptest.ResponseRecorder) (*mcp.CallToolResult, error) {
	res := rec.Result()
	defer func() { _ = res.Body.Close() }()
	body := bytes.TrimSpace(rec.Body.Bytes())

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return toolError(fmt.Sprintf("HTTP %d: %s", res.StatusCode, body)), nil
	}
	if op.Response == nil || len(body) == 0 {
		return s.finish(op.ID, result{text: op.ID + ": done", data: map[string]any{}})
	}
	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, errs.With("operation", op.ID).Wrapf(err, "decode response")
	}
	if _, isObject := decoded.(map[string]any); !isObject {
		decoded = map[string]any{"result": decoded}
	}
	return s.finish(op.ID, result{text: fmt.Sprintf("%s: %s", op.ID, http.StatusText(res.StatusCode)), data: decoded})
}

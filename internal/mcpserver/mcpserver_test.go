package mcpserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// alice is a fixture user.
var alice = fixture.ID(fixture.UserAdminLearner)

func body(t *testing.T, res *http.Response) string {
	t.Helper()
	defer func() { _ = res.Body.Close() }()
	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return string(b)
}

func initialize(t *testing.T, o *testsupport.OAuth, mutate func(*http.Request)) *http.Response {
	t.Helper()
	return o.Do(http.MethodPost, "/mcp", nil, func(r *http.Request) {
		r.Body = io.NopCloser(strings.NewReader(
			`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"0"}}}`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json, text/event-stream")
		if mutate != nil {
			mutate(r)
		}
	})
}

func TestRequestWithoutATokenIsChallengedWithTheResourceMetadata(t *testing.T) {
	o := testsupport.NewOAuth(t)

	for name, mutate := range map[string]func(*http.Request){
		"no credentials": nil,
		"garbage token":  func(r *http.Request) { r.Header.Set("Authorization", "Bearer not-a-token") },
		"basic scheme":   func(r *http.Request) { r.Header.Set("Authorization", "Basic YTpi") },
	} {
		t.Run(name, func(t *testing.T) {
			res := initialize(t, o, mutate)
			defer func() { _ = res.Body.Close() }()

			assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
			assert.Contains(t, res.Header.Get("WWW-Authenticate"),
				`resource_metadata="`+testsupport.OAuthIssuer+`/.well-known/oauth-protected-resource/mcp"`)
		})
	}
}

func TestResourceMetadataIsServedAtBothWellKnownLocations(t *testing.T) {
	o := testsupport.NewOAuth(t)

	for _, path := range []string{
		"/.well-known/oauth-protected-resource",
		"/.well-known/oauth-protected-resource/mcp",
	} {
		res := o.Do(http.MethodGet, path, nil)
		defer func() { _ = res.Body.Close() }()
		require.Equal(t, http.StatusOK, res.StatusCode, path)

		var doc map[string]any
		require.NoError(t, json.Unmarshal([]byte(body(t, res)), &doc))
		assert.Equal(t, testsupport.MCPEndpoint, doc["resource"], path)
		assert.Equal(t, []any{testsupport.OAuthIssuer}, doc["authorization_servers"], path)
		assert.Equal(t, []any{"header"}, doc["bearer_methods_supported"], path)
	}
}

func TestWebsiteSessionCookieIsNotAcceptedOnTheEndpoint(t *testing.T) {
	o := testsupport.NewOAuth(t)
	session := testsupport.SessionFor(t, o.DB, alice)

	res := initialize(t, o, func(r *http.Request) {
		r.AddCookie(session.Cookie)
	})
	defer func() { _ = res.Body.Close() }()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Contains(t, res.Header.Get("WWW-Authenticate"), "resource_metadata=")
}

func TestBearerTokenIsNotAcceptedOnAnyAPIOperation(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead, oauthserver.ScopeAdmin)
	bearer := func(r *http.Request) { r.Header.Set("Authorization", "Bearer "+tokens.AccessToken) }

	// An operation that identifies the caller when it can: the cookie session
	// resolves the user, the bearer token does not, so a visitor comes back.
	session := testsupport.SessionFor(t, o.DB, alice)
	withCookie := o.Do(http.MethodGet, "/api/me", nil, func(r *http.Request) { r.AddCookie(session.Cookie) })
	defer func() { _ = withCookie.Body.Close() }()
	require.Equal(t, http.StatusOK, withCookie.StatusCode)
	assert.Contains(t, body(t, withCookie), strconv.Itoa(alice))

	withBearer := o.Do(http.MethodGet, "/api/me", nil, bearer)
	defer func() { _ = withBearer.Body.Close() }()
	require.Equal(t, http.StatusOK, withBearer.StatusCode)
	assert.NotContains(t, body(t, withBearer), strconv.Itoa(alice))

	// Operations that need a signed-in user refuse it.
	for _, path := range []string{"/api/my", "/api/admin/api/users"} {
		res := o.Do(http.MethodGet, path, nil, bearer)
		_ = res.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode, path)
	}
}

func TestClientListsCoursesThroughOneToolWithTheContractSchema(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	session := o.MCPClient(tokens.AccessToken)

	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)
	var tool *mcp.Tool
	for _, candidate := range tools.Tools {
		if candidate.Name == "listCourses" {
			tool = candidate
		}
	}
	require.NotNil(t, tool, "the catalogue tool is offered beside the admin ones")
	assert.True(t, tool.Annotations.ReadOnlyHint)

	result, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "listCourses"})
	require.NoError(t, err)
	require.False(t, result.IsError, "%+v", result.Content)

	// The catalogue is what the HTTP operation answers: fixtures carry courses.
	structured, ok := result.StructuredContent.(map[string]any)
	require.True(t, ok)
	courses, ok := structured["courses"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, courses)
	assert.Contains(t, courses[0], "slug")
	require.Len(t, result.Content, 1)
	text, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, text.Text, courses[0].(map[string]any)["slug"])

	// The advertised schema, which comes from the contract, accepts the content.
	raw, err := json.Marshal(tool.OutputSchema)
	require.NoError(t, err)
	var schema jsonschema.Schema
	require.NoError(t, json.Unmarshal(raw, &schema))
	resolved, err := schema.Resolve(nil)
	require.NoError(t, err)
	assert.NoError(t, resolved.Validate(structured))
}

func TestToolAnswersInTheLocaleOfTheRequest(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	locales := map[string]bool{}
	for _, language := range []string{"en", "ru"} {
		session := o.MCPClient(tokens.AccessToken, func(r *http.Request) { r.Header.Set("Accept-Language", language) })
		result, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "listCourses"})
		require.NoError(t, err)
		for _, c := range result.StructuredContent.(map[string]any)["courses"].([]any) {
			locales[language+":"+c.(map[string]any)["locale"].(string)] = true
		}
	}

	assert.Equal(t, map[string]bool{"en:en": true, "ru:ru": true}, locales)
}

func TestToolRefusesAConnectionWithoutTheCatalogReadScope(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeProgressRead)
	session := o.MCPClient(tokens.AccessToken)

	result, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "listCourses"})

	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Nil(t, result.StructuredContent)
}

// verifier answers with what a test sets, to reach tokens the real server never
// issues.
type verifier struct {
	token oauthserver.Token
	err   error
}

func (v verifier) VerifyAccessToken(context.Context, string) (oauthserver.Token, error) {
	return v.token, v.err
}
func (verifier) Resource() string { return testsupport.MCPEndpoint }

// emptyCatalogue answers an empty course list; the other methods are the nil
// embedded interface's, which no test of it reaches.
type emptyCatalogue struct{ mcpserver.Catalogue }

func (emptyCatalogue) ListCourses(context.Context) ([]api.CourseCatalogItem, error) { return nil, nil }

func status(t *testing.T, v verifier) int {
	t.Helper()
	server, err := mcpserver.New(mcpserver.Deps{Verifier: v, Catalogue: emptyCatalogue{}, Issuer: testsupport.OAuthIssuer})
	require.NoError(t, err)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/mcp", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer anything")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	return rec.Code
}

func TestTokenIsAcceptedOnlyForTheCanonicalResourceAndAUserSubject(t *testing.T) {
	good := oauthserver.Token{UserID: alice, Scopes: []string{"catalog:read"}, Audience: []string{testsupport.MCPEndpoint}}

	other := good
	other.Audience = []string{"https://elsewhere.example/mcp"}
	none := good
	none.Audience = nil
	client := good
	client.UserID = 0

	assert.NotEqual(t, http.StatusUnauthorized, status(t, verifier{token: good}), "a token for the endpoint passes")
	assert.Equal(t, http.StatusUnauthorized, status(t, verifier{token: other}), "another audience")
	assert.Equal(t, http.StatusUnauthorized, status(t, verifier{token: none}), "no audience")
	assert.Equal(t, http.StatusUnauthorized, status(t, verifier{token: client}), "a client subject")
	assert.Equal(t, http.StatusUnauthorized, status(t, verifier{err: oauthserver.ErrInvalidToken}), "an unknown token")
}

func TestVerifierFailureOfOursIsAServerErrorWithoutDetail(t *testing.T) {
	server, err := mcpserver.New(mcpserver.Deps{
		Verifier:  verifier{err: errors.New("pq: secret detail")},
		Catalogue: emptyCatalogue{},
		Issuer:    testsupport.OAuthIssuer,
	})
	require.NoError(t, err)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/mcp", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer anything")
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "secret detail")
}

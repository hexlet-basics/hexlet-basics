package oauthserver_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/oauthclient"
	"hexletbasics/internal/testsupport"
)

// registered is the RFC 7591 response of a successful registration.
type registered struct {
	ClientID                string   `json:"client_id"`
	ClientIDIssuedAt        int64    `json:"client_id_issued_at"`
	ClientSecret            string   `json:"client_secret"`
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
}

type registrationError struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func decodeJSON[T any](t *testing.T, res *http.Response) T {
	t.Helper()
	defer func() { _ = res.Body.Close() }()
	var v T
	require.NoError(t, json.NewDecoder(res.Body).Decode(&v))
	return v
}

// request is a valid registration an attacker edits one field of.
func request(edit func(map[string]any)) map[string]any {
	m := map[string]any{
		"client_name":   "Attacker",
		"redirect_uris": []string{"https://app.example.com/callback"},
	}
	if edit != nil {
		edit(m)
	}
	return m
}

func clientCount(t *testing.T, o *testsupport.OAuth) int {
	t.Helper()
	n, err := o.DB.OAuthClient.Query().Count(t.Context())
	require.NoError(t, err)
	return n
}

func TestTheMetadataAdvertisesTheRegistrationEndpointAtBothPaths(t *testing.T) {
	o := testsupport.NewOAuth(t)
	for _, path := range []string{"/.well-known/oauth-authorization-server", "/.well-known/openid-configuration"} {
		doc := decodeJSON[map[string]any](t, o.Do(http.MethodGet, path, nil))
		assert.Equal(t, testsupport.OAuthIssuer+"/oauth/register", doc["registration_endpoint"], path)
	}
}

func TestAPublicClientWithAnHTTPSRedirectRegistersAndCompletesTheFlow(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Register(request(func(m map[string]any) {
		m["scope"] = "catalog:read offline_access"
	}))

	require.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "no-store", res.Header.Get("Cache-Control"))
	got := decodeJSON[registered](t, res)
	require.NotEmpty(t, got.ClientID)
	assert.Empty(t, got.ClientSecret, "a public client has no secret")
	assert.Equal(t, "none", got.TokenEndpointAuthMethod)
	assert.Equal(t, []string{"authorization_code", "refresh_token"}, got.GrantTypes)
	assert.Equal(t, []string{"code"}, got.ResponseTypes)
	assert.Equal(t, []string{"https://app.example.com/callback"}, got.RedirectURIs)
	assert.ElementsMatch(t, []string{"catalog:read", "offline_access"}, strings.Fields(got.Scope))

	f := o.StartWith(alice, testsupport.StartOptions{
		ClientID: got.ClientID, RedirectURIs: got.RedirectURIs, Scopes: []string{"catalog:read"},
	})
	require.Equal(t, http.StatusFound, f.Authorize.StatusCode)
	require.NotEmpty(t, f.RequestID)
	code := f.Approve()
	tokens, exchange := f.Exchange(code, testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, exchange.StatusCode)
	verified, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, got.ClientID, verified.ClientID)
	assert.Contains(t, verified.Scopes, "catalog:read")
}

func TestALoopbackRedirectRegistersOnAnyPort(t *testing.T) {
	o := testsupport.NewOAuth(t)
	for i, uri := range []string{"http://127.0.0.1:53124/callback", "http://localhost:8000/cb", "http://[::1]:9/cb"} {
		res := o.Register(request(func(m map[string]any) { m["redirect_uris"] = []string{uri} }),
			testsupport.FromAddress(fmt.Sprintf("203.0.113.%d", i+1)))
		assert.Equal(t, http.StatusCreated, res.StatusCode, uri)
	}
}

func TestARegistrationThatNamesNoScopeGetsTheReadOnlyDefault(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Register(request(nil))

	require.Equal(t, http.StatusCreated, res.StatusCode)
	got := decodeJSON[registered](t, res)
	assert.ElementsMatch(t, []string{"catalog:read", "progress:read"}, strings.Fields(got.Scope))
	assert.NotContains(t, got.Scope, "admin", "admin must always be asked for by name")
	row, err := o.DB.OAuthClient.Query().Where(oauthclient.ClientID(got.ClientID)).Only(t.Context())
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"catalog:read", "progress:read"}, []string(row.Scopes))
}

func TestARegisteredClientCannotRequestAScopeItDidNotRegisterFor(t *testing.T) {
	o := testsupport.NewOAuth(t)
	got := decodeJSON[registered](t, o.Register(request(func(m map[string]any) { m["scope"] = "catalog:read" })))

	f := o.StartWith(alice, testsupport.StartOptions{
		ClientID: got.ClientID, RedirectURIs: got.RedirectURIs, Scopes: []string{"admin"},
	})

	require.Equal(t, http.StatusFound, f.Authorize.StatusCode)
	to, err := f.Authorize.Location()
	require.NoError(t, err)
	assert.Equal(t, "invalid_scope", to.Query().Get("error"))
}

// TestRegistrationRefusals is the attacker's list: every request here asks for
// something an anonymous caller must not get, and each leaves no client behind.
func TestRegistrationRefusals(t *testing.T) {
	redirects := func(uris ...string) func(map[string]any) {
		return func(m map[string]any) { m["redirect_uris"] = uris }
	}
	set := func(key string, value any) func(map[string]any) {
		return func(m map[string]any) { m[key] = value }
	}
	var manyRedirects []string
	for i := range 50 {
		manyRedirects = append(manyRedirects, fmt.Sprintf("https://app.example.com/cb%d", i))
	}

	cases := map[string]struct {
		body any
		code string
	}{
		"plain http non-loopback redirect":       {request(redirects("http://evil.example.com/cb")), "invalid_redirect_uri"},
		"http redirect beside a good one":        {request(redirects("https://ok.example.com/cb", "http://evil.example.com/cb")), "invalid_redirect_uri"},
		"host that merely starts with localhost": {request(redirects("http://localhost.evil.example.com/cb")), "invalid_redirect_uri"},
		"loopback lookalike with userinfo":       {request(redirects("http://127.0.0.1@evil.example.com/cb")), "invalid_redirect_uri"},
		"wildcard host":                          {request(redirects("https://*.example.com/cb")), "invalid_redirect_uri"},
		"wildcard path":                          {request(redirects("https://app.example.com/*")), "invalid_redirect_uri"},
		"encoded wildcard":                       {request(redirects("https://app.example.com/%2A")), "invalid_redirect_uri"},
		"redirect with a fragment":               {request(redirects("https://app.example.com/cb#x")), "invalid_redirect_uri"},
		"javascript scheme":                      {request(redirects("javascript:alert(1)")), "invalid_redirect_uri"},
		"custom scheme":                          {request(redirects("cursor://callback")), "invalid_redirect_uri"},
		"relative redirect":                      {request(redirects("/callback")), "invalid_redirect_uri"},
		"no redirect":                            {request(func(m map[string]any) { delete(m, "redirect_uris") }), "invalid_redirect_uri"},
		"empty redirect list":                    {request(redirects()), "invalid_redirect_uri"},
		"redirect list of the wrong type":        {request(set("redirect_uris", "https://app.example.com/cb")), "invalid_client_metadata"},
		"too many redirects":                     {request(set("redirect_uris", manyRedirects)), "invalid_redirect_uri"},
		"client_credentials grant":               {request(set("grant_types", []string{"client_credentials"})), "invalid_client_metadata"},
		"implicit grant":                         {request(set("grant_types", []string{"implicit"})), "invalid_client_metadata"},
		"extra grant beside the allowed ones":    {request(set("grant_types", []string{"authorization_code", "client_credentials"})), "invalid_client_metadata"},
		"password grant":                         {request(set("grant_types", []string{"password"})), "invalid_client_metadata"},
		"device grant":                           {request(set("grant_types", []string{"urn:ietf:params:oauth:grant-type:device_code"})), "invalid_client_metadata"},
		"refresh without the code grant":         {request(set("grant_types", []string{"refresh_token"})), "invalid_client_metadata"},
		"token response type":                    {request(set("response_types", []string{"token"})), "invalid_client_metadata"},
		"id_token response type":                 {request(set("response_types", []string{"code", "id_token"})), "invalid_client_metadata"},
		"caller-chosen audience":                 {request(set("audience", []string{"https://evil.example.com"})), "invalid_client_metadata"},
		"caller-chosen resource":                 {request(set("resource", "https://evil.example.com")), "invalid_client_metadata"},
		"scope admin lookalike":                  {request(set("scope", "catalog:read admin:write")), "invalid_client_metadata"},
		"unknown scope":                          {request(set("scope", "superuser")), "invalid_client_metadata"},
		"profile scope":                          {request(set("scope", "profile email")), "invalid_client_metadata"},
		"scope of the wrong type":                {request(set("scope", []string{"admin"})), "invalid_client_metadata"},
		"confidential client by method basic":    {request(set("token_endpoint_auth_method", "client_secret_basic")), "invalid_client_metadata"},
		"confidential client by method post":     {request(set("token_endpoint_auth_method", "client_secret_post")), "invalid_client_metadata"},
		"private_key_jwt client":                 {request(set("token_endpoint_auth_method", "private_key_jwt")), "invalid_client_metadata"},
		"client secret supplied":                 {request(set("client_secret", "chosen-by-the-caller")), "invalid_client_metadata"},
		"jwks supplied":                          {request(set("jwks_uri", "https://evil.example.com/jwks")), "invalid_client_metadata"},
		"skip_consent":                           {request(set("skip_consent", true)), "invalid_client_metadata"},
		"skip_logout_consent":                    {request(set("skip_logout_consent", true)), "invalid_client_metadata"},
		"caller-chosen client_id":                {request(set("client_id", "test-mcp-client")), "invalid_client_metadata"},
		"overlong name":                          {request(set("client_name", strings.Repeat("a", 500))), "invalid_client_metadata"},
		"name of the wrong type":                 {request(set("client_name", 7)), "invalid_client_metadata"},
		"malformed JSON":                         {`{"redirect_uris": [`, "invalid_client_metadata"},
		"not an object":                          {`["https://app.example.com/cb"]`, "invalid_client_metadata"},
		"null":                                   {`null`, "invalid_client_metadata"},
		"empty body":                             {``, "invalid_client_metadata"},
		"trailing garbage":                       {`{"redirect_uris":["https://app.example.com/cb"]} {"admin":true}`, "invalid_client_metadata"},
	}
	o := testsupport.NewOAuth(t)
	before := clientCount(t, o)
	i := 0
	for name, tc := range cases {
		i++
		t.Run(name, func(t *testing.T) {
			// A different address per case: this test is about validation, not
			// about the rate limit.
			res := o.Register(tc.body, testsupport.FromAddress(fmt.Sprintf("198.51.100.%d", i)))

			require.Equal(t, http.StatusBadRequest, res.StatusCode)
			assert.Equal(t, tc.code, decodeJSON[registrationError](t, res).Error)
		})
	}
	assert.Equal(t, before, clientCount(t, o), "no refused registration may leave a client behind")
}

func TestAnOversizedBodyIsRefusedWithoutBeingRead(t *testing.T) {
	o := testsupport.NewOAuth(t)
	before := clientCount(t, o)
	body := fmt.Sprintf(`{"redirect_uris":["https://app.example.com/cb"],"client_name":"%s"}`, strings.Repeat("a", 1<<20))

	res := o.Register(body)

	assert.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
	assert.Equal(t, "invalid_client_metadata", decodeJSON[registrationError](t, res).Error)
	assert.Equal(t, before, clientCount(t, o))
}

func TestOnlyJSONIsAccepted(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Register(`redirect_uris=https://app.example.com/cb`, func(r *http.Request) {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	assert.Equal(t, http.StatusUnsupportedMediaType, res.StatusCode)
}

func TestRegistrationAcceptsPostOnly(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Do(http.MethodGet, "/oauth/register", nil)

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	assert.Equal(t, "POST, OPTIONS", res.Header.Get("Allow"))
}

func TestARefusedRegistrationNeverRevealsAnExistingClient(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.SeedClient("catalog:read")

	res := o.Register(request(func(m map[string]any) { m["client_id"] = "test-mcp-client" }))

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.NotContains(t, decodeJSON[registrationError](t, res).Description, "exists")
}

func TestRegistrationIsRateLimitedPerAddress(t *testing.T) {
	o := testsupport.NewOAuth(t)
	attacker := testsupport.FromAddress("203.0.113.50")
	for i := range 10 {
		res := o.Register(request(nil), attacker)
		require.Equal(t, http.StatusCreated, res.StatusCode, "registration %d is within the limit", i+1)
	}

	res := o.Register(request(nil), attacker)

	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	assert.NotEmpty(t, res.Header.Get("Retry-After"))
	assert.Equal(t, "temporarily_unavailable", decodeJSON[registrationError](t, res).Error)
	assert.Equal(t, 10, clientCount(t, o), "the refused registration created nothing")

	other := o.Register(request(nil), testsupport.FromAddress("203.0.113.51"))
	assert.Equal(t, http.StatusCreated, other.StatusCode, "another address has its own budget")

	o.Advance(time.Hour + time.Minute)
	again := o.Register(request(nil), attacker)
	assert.Equal(t, http.StatusCreated, again.StatusCode, "the budget returns when the window passes")
}

func TestRefusedAttemptsSpendTheRateLimitToo(t *testing.T) {
	o := testsupport.NewOAuth(t)
	attacker := testsupport.FromAddress("203.0.113.60")
	for range 10 {
		res := o.Register(`{`, attacker)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	}

	res := o.Register(request(nil), attacker)

	assert.Equal(t, http.StatusTooManyRequests, res.StatusCode, "probing with bad requests must not be free")
}

func TestPrivateHopsInAForwardingChainDoNotChooseTheBucket(t *testing.T) {
	o := testsupport.NewOAuth(t)
	// Behind the ingress the chain ends with the real peer; private addresses a
	// caller prepends are skipped by the client-IP middleware and cannot mint a
	// fresh budget.
	for i := range 10 {
		res := o.Register(request(nil), func(r *http.Request) {
			r.Header.Set("X-Forwarded-For", fmt.Sprintf("10.9.9.%d, 203.0.113.70", i))
		})
		require.Equal(t, http.StatusCreated, res.StatusCode)
	}

	res := o.Register(request(nil), func(r *http.Request) {
		r.Header.Set("X-Forwarded-For", "10.8.8.8, 203.0.113.70")
	})

	assert.Equal(t, http.StatusTooManyRequests, res.StatusCode)
}

func TestARegisteredNameIsStoredAsPlainTrimmedText(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Register(request(func(m map[string]any) { m["client_name"] = "  <b>Claude</b>\n\x00  " }))

	require.Equal(t, http.StatusCreated, res.StatusCode)
	row, err := o.DB.OAuthClient.Query().Where(oauthclient.ClientID(decodeJSON[registered](t, res).ClientID)).Only(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "<b>Claude</b>", row.Name, "control characters are stripped and the name is trimmed")
}

func TestARegistrationThatNamesOnlyOfflineAccessHoldsNoCustomScope(t *testing.T) {
	o := testsupport.NewOAuth(t)

	res := o.Register(request(func(m map[string]any) { m["scope"] = "offline_access" }))

	require.Equal(t, http.StatusCreated, res.StatusCode)
	got := decodeJSON[registered](t, res)
	assert.Equal(t, "offline_access", got.Scope)
}

func TestAnUnknownDescriptiveFieldIsIgnoredButAnUnknownOneIsNot(t *testing.T) {
	o := testsupport.NewOAuth(t)

	ok := o.Register(request(func(m map[string]any) {
		m["client_uri"] = "https://app.example.com"
		m["application_type"] = "native"
	}), testsupport.FromAddress("203.0.113.80"))
	refused := o.Register(request(func(m map[string]any) { m["software_statement"] = "eyJ" }),
		testsupport.FromAddress("203.0.113.81"))

	assert.Equal(t, http.StatusCreated, ok.StatusCode)
	assert.Equal(t, http.StatusBadRequest, refused.StatusCode)
}

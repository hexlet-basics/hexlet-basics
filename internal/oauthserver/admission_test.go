package oauthserver_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

const mcpResource = testsupport.OAuthIssuer + "/mcp"

// jwtClaims decodes the payload of a JWT without verifying it, to read what an
// ID token says.
func jwtClaims(t *testing.T, jwt string) map[string]any {
	t.Helper()
	parts := strings.Split(jwt, ".")
	require.Len(t, parts, 3)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	var claims map[string]any
	require.NoError(t, json.Unmarshal(payload, &claims))
	return claims
}

// refusal is the redirect a client receives for a request that was refused.
func refusal(t *testing.T, f *testsupport.Flow) url.Values {
	t.Helper()
	require.Equal(t, http.StatusFound, f.Authorize.StatusCode)
	to, err := f.Authorize.Location()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:8976", to.Host, "the error goes to the client, not to the consent screen")
	return to.Query()
}

func TestTheThreeScopesAreAdmittedAndAppearInTheIssuedToken(t *testing.T) {
	o := testsupport.NewOAuth(t)
	all := []string{oauthserver.ScopeCatalogRead, oauthserver.ScopeProgressRead, oauthserver.ScopeAdmin}

	tokens := o.Connect(alice, all...)

	assert.ElementsMatch(t, append(all, "offline_access"), tokens.Scopes)
	verified, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
	assert.ElementsMatch(t, append(all, "offline_access"), verified.Scopes)
}

func TestTheFlowWorksWithoutTheOpenIDScope(t *testing.T) {
	o := testsupport.NewOAuth(t)

	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	assert.NotContains(t, tokens.Scopes, "openid")
	assert.Equal(t, "Bearer", tokens.TokenType)
	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
}

func TestAnUnknownScopeIsRefusedWithInvalidScopeNotDropped(t *testing.T) {
	o := testsupport.NewOAuth(t)
	for _, scope := range []string{
		"catalog:read superpower", // unknown next to a good one: the whole request is refused
		"profile",                 // a standard OpenID scope that discloses nothing here
		"email",
	} {
		f := o.StartWith(alice, testsupport.StartOptions{
			Scopes: []string{oauthserver.ScopeCatalogRead},
			Params: func(q url.Values) { q.Set("scope", scope) },
		})
		q := refusal(t, f)
		assert.Equal(t, "invalid_scope", q.Get("error"), scope)
		assert.Contains(t, q.Get("error_description"), "Supported scopes", scope)
		assert.Equal(t, f.State, q.Get("state"))
	}
	count, err := o.DB.OAuthAuthRequest.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, count, "a refused request leaves nothing to approve")
}

func TestAKnownScopeTheClientWasNotRegisteredForIsRefused(t *testing.T) {
	o := testsupport.NewOAuth(t)

	f := o.StartWith(alice, testsupport.StartOptions{
		Scopes: []string{oauthserver.ScopeCatalogRead},
		Params: func(q url.Values) { q.Set("scope", "catalog:read admin") },
	})

	assert.Equal(t, "invalid_scope", refusal(t, f).Get("error"))
}

func TestTheResourceBindsTheTokenAudienceToTheCanonicalMCPURL(t *testing.T) {
	o := testsupport.NewOAuth(t)
	assert.Equal(t, mcpResource, o.Server.Resource(), "the issuer from configuration plus /mcp")

	for name, resource := range map[string][]string{
		"named":   {mcpResource},
		"omitted": nil,
	} {
		f := o.StartWith(alice, testsupport.StartOptions{
			Scopes: []string{"openid", oauthserver.ScopeCatalogRead},
			Params: func(q url.Values) {
				if resource != nil {
					q["resource"] = resource
				}
			},
		})
		require.Equal(t, http.StatusFound, f.Authorize.StatusCode, name)
		code := f.Approve()
		tokens, res := f.ExchangeWith(code, testsupport.PKCEVerifier(), url.Values{"resource": {mcpResource}})
		require.Equal(t, http.StatusOK, res.StatusCode, name)

		verified, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, []string{mcpResource}, verified.Audience, name)
		assert.Contains(t, jwtClaims(t, tokens.IDToken)["aud"], mcpResource, name)

		// A refreshed token is for the same resource.
		res = o.Do(http.MethodPost, "/oauth/token", url.Values{
			"grant_type": {"refresh_token"}, "client_id": {testsupport.OAuthClientID}, "refresh_token": {tokens.RefreshToken},
		})
		require.Equal(t, http.StatusOK, res.StatusCode, name)
		var refreshed testsupport.Tokens
		decode(t, res, &refreshed)
		verified, err = o.Server.VerifyAccessToken(t.Context(), refreshed.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, []string{mcpResource}, verified.Audience, name)
	}
}

func TestAnotherResourceIsRefusedWithInvalidTarget(t *testing.T) {
	o := testsupport.NewOAuth(t)

	for _, resource := range [][]string{
		{"https://other.example/mcp"},
		{testsupport.OAuthIssuer}, // the origin is not the MCP endpoint
		{testsupport.OAuthIssuer + "/mcp/"},
		{mcpResource, "https://other.example/mcp"}, // one wrong value among right ones
	} {
		f := o.StartWith(alice, testsupport.StartOptions{
			Scopes: []string{oauthserver.ScopeCatalogRead},
			Params: func(q url.Values) { q["resource"] = resource },
		})
		q := refusal(t, f)
		assert.Equal(t, "invalid_target", q.Get("error"), resource)
		assert.Contains(t, q.Get("error_description"), mcpResource)
	}
}

func TestTheTokenRequestRefusesAnotherResource(t *testing.T) {
	o := testsupport.NewOAuth(t)
	f := o.Start(alice, oauthserver.ScopeCatalogRead)
	code := f.Approve()

	_, res := f.ExchangeWith(code, testsupport.PKCEVerifier(), url.Values{"resource": {"https://other.example/mcp"}})
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var body struct{ Error string }
	decode(t, res, &body)
	assert.Equal(t, "invalid_target", body.Error)

	// The refusal did not spend the code.
	_, res = f.Exchange(code, testsupport.PKCEVerifier())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRevokingTheRowStillRevokesATokenWithAnAudience(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	_, err := o.DB.OAuthAccessToken.Delete().Exec(t.Context())
	require.NoError(t, err)

	_, err = o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestTheAuthorizationResponseCarriesTheIssuer(t *testing.T) {
	o := testsupport.NewOAuth(t)

	// Success.
	f := o.Start(alice, oauthserver.ScopeCatalogRead)
	f.Approve()
	assert.Equal(t, testsupport.OAuthIssuer, f.Callback.Query().Get("iss"))

	// The user denies.
	f = o.Start(alice, oauthserver.ScopeCatalogRead)
	res := f.Decide(false)
	require.Equal(t, http.StatusOK, res.StatusCode)
	var decision struct {
		RedirectURL string `json:"redirectUrl"`
	}
	decode(t, res, &decision)
	denied, err := url.Parse(decision.RedirectURL)
	require.NoError(t, err)
	assert.Equal(t, "access_denied", denied.Query().Get("error"))
	assert.Equal(t, testsupport.OAuthIssuer, denied.Query().Get("iss"))

	// The request is refused.
	f = o.StartWith(alice, testsupport.StartOptions{
		Scopes: []string{oauthserver.ScopeCatalogRead},
		Params: func(q url.Values) { q.Set("scope", "superpower") },
	})
	assert.Equal(t, testsupport.OAuthIssuer, refusal(t, f).Get("iss"))

	// The hop to the consent screen is not a response to the client.
	f = o.Start(alice, oauthserver.ScopeCatalogRead)
	to, err := f.Authorize.Location()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(to.String(), testsupport.OAuthIssuer+"/connect"))
	assert.NotContains(t, to.Query(), "iss")
}

func TestTheMetadataDeclaresTheIssuerParameterAtBothPaths(t *testing.T) {
	o := testsupport.NewOAuth(t)

	for _, path := range []string{"/.well-known/oauth-authorization-server", "/.well-known/openid-configuration"} {
		res := o.Do(http.MethodGet, path, nil)
		require.Equal(t, http.StatusOK, res.StatusCode, path)
		var doc map[string]any
		decode(t, res, &doc)
		assert.Equal(t, true, doc["authorization_response_iss_parameter_supported"], path)
		assert.Equal(t, testsupport.OAuthIssuer, doc["issuer"], path)
		assert.ElementsMatch(t, []any{"catalog:read", "progress:read", "admin", "offline_access"}, doc["scopes_supported"], path)
	}
}

func TestALoopbackRedirectMatchesOnAnyPortWithTheSamePath(t *testing.T) {
	o := testsupport.NewOAuth(t)
	registered := []string{"http://127.0.0.1:8976/callback", "https://app.example/cb"}

	connects := func(redirect string) bool {
		f := o.StartWith(alice, testsupport.StartOptions{
			Scopes:       []string{oauthserver.ScopeCatalogRead},
			RedirectURIs: registered,
			Params:       func(q url.Values) { q.Set("redirect_uri", redirect) },
		})
		if f.Authorize.StatusCode != http.StatusFound {
			return false
		}
		to, err := f.Authorize.Location()
		require.NoError(t, err)
		return strings.HasPrefix(to.String(), testsupport.OAuthIssuer+"/connect")
	}

	// Any port, the same path: native clients pick a free port per run (RFC 8252).
	for _, redirect := range []string{
		"http://127.0.0.1:8976/callback",
		"http://127.0.0.1:1/callback",
		"http://127.0.0.1:54321/callback",
		"http://127.0.0.1/callback",
		// The library compares only path and query once both sides are loopback,
		// so another loopback address or name also matches. That reaches only the
		// user's own machine (RFC 8252 section 8.3); recorded in ADR-0022.
		"http://localhost:5000/callback",
		"http://[::1]:5000/callback",
	} {
		assert.True(t, connects(redirect), redirect)
	}

	// Anything else must match exactly.
	for _, redirect := range []string{
		"http://127.0.0.1:5000/other",                 // a different path
		"http://127.0.0.1:5000/callback/extra",        // a longer path
		"http://127.0.0.1:5000/callback?next=/x",      // a different query
		"http://evil.example:8976/callback",           // a non-loopback host on the registered port
		"http://evil.example/callback",                //
		"https://app.example/other",                   // a registered https host, another path
		"https://app.example:8443/cb",                 // no port freedom for a non-loopback host
		"http://app.example/cb",                       // plain http for a registered https host
		"https://app.example/cb",                      // (control: exact match is fine, asserted below)
		"https://evil.example/cb",                     // not registered at all
		"javascript:alert(1)",                         // not a web redirect
		"http://127.0.0.1:5000@evil.example/callback", // userinfo trick
	} {
		if redirect == "https://app.example/cb" {
			assert.True(t, connects(redirect), "an exact https match is accepted")
			continue
		}
		assert.False(t, connects(redirect), redirect)
	}
}

func TestAPlainHTTPRedirectOnANonLoopbackHostIsRefusedEvenWhenRegistered(t *testing.T) {
	o := testsupport.NewOAuth(t)
	redirect := "http://app.example/cb"

	f := o.StartWith(alice, testsupport.StartOptions{
		Scopes:       []string{oauthserver.ScopeCatalogRead},
		RedirectURIs: []string{redirect},
		Params:       func(q url.Values) { q.Set("redirect_uri", redirect) },
	})

	// There is no trusted redirect to send the error to: the browser stays here.
	assert.NotEqual(t, http.StatusFound, f.Authorize.StatusCode)
	count, err := o.DB.OAuthAuthRequest.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, count)
}

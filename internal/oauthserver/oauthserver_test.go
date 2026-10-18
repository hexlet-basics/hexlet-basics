package oauthserver_test

import (
	"encoding/base64"
	"encoding/json"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/oauthaccesstoken"
	"hexletbasics/ent/oauthrefreshtoken"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// alice is the fixture user who signs in and decides.
var alice = fixture.ID(fixture.UserAdminLearner)

func decode(t *testing.T, res *http.Response, into any) {
	t.Helper()

	defer func() { _ = res.Body.Close() }()

	require.NoError(t, json.NewDecoder(res.Body).Decode(into))
}

func TestMetadataIsServedAtTheRootWellKnownPathWithEndpointsUnderOAuth(t *testing.T) {
	o := testsupport.NewOAuth(t)

	for _, path := range []string{
		"/.well-known/oauth-authorization-server",
		// Clients that miss the RFC 8414 document fall back to this one.
		"/.well-known/openid-configuration",
	} {
		res := o.Do(http.MethodGet, path, nil, func(r *http.Request) {
			// The issuer is configuration, never the Host header.
			r.Host = "evil.example"
		})
		require.Equal(t, http.StatusOK, res.StatusCode, path)

		var doc map[string]any
		decode(t, res, &doc)
		assert.Equal(t, testsupport.OAuthIssuer, doc["issuer"], path)
		assert.Equal(t, testsupport.OAuthIssuer+"/oauth/authorize", doc["authorization_endpoint"])
		assert.Equal(t, testsupport.OAuthIssuer+"/oauth/token", doc["token_endpoint"])
		assert.Equal(t, testsupport.OAuthIssuer+"/oauth/revoke", doc["revocation_endpoint"])
		assert.Equal(t, testsupport.OAuthIssuer+"/oauth/keys", doc["jwks_uri"])
		// S256 is the only PKCE method, and code the only flow.
		assert.Equal(t, []any{"S256"}, doc["code_challenge_methods_supported"])
		assert.Equal(t, []any{"code"}, doc["response_types_supported"])
		assert.Equal(t, []any{"authorization_code", "refresh_token"}, doc["grant_types_supported"])
		assert.NotContains(t, doc, "device_authorization_endpoint")
	}
}

func TestTheLibrariesRootLevelRoutesAreNotReachable(t *testing.T) {
	o := testsupport.NewOAuth(t)

	for _, path := range []string{"/authorize", "/keys", "/revoke", "/userinfo", "/end_session", "/healthz", "/oauth/introspect", "/oauth/device_authorization"} {
		res := o.Do(http.MethodPost, path, url.Values{})
		assert.Equal(t, http.StatusNotFound, res.StatusCode, path)
	}
}

func TestPreRegisteredClientCompletesAuthorizationCodeWithPKCE(t *testing.T) {
	o := testsupport.NewOAuth(t)

	flow := o.Start(alice, oauthserver.ScopeCatalogRead, oauthserver.ScopeProgressRead)

	// The authorization request sends the browser to the consent screen on the
	// canonical origin, which is an SSR route outside /oauth.
	require.Equal(t, http.StatusFound, flow.Authorize.StatusCode)
	to, err := flow.Authorize.Location()
	require.NoError(t, err)
	assert.Equal(t, testsupport.OAuthIssuer+"/connect", to.Scheme+"://"+to.Host+to.Path)
	require.NotEmpty(t, flow.RequestID)

	// After cookie sign-in the consent screen's data names the application, the
	// host it returns to, the scopes asked for and the account.
	var consent struct {
		ClientName   string   `json:"clientName"`   //nolint:tagliatelle // decodes the contract's camelCase JSON
		RedirectHost string   `json:"redirectHost"` //nolint:tagliatelle // decodes the contract's camelCase JSON
		Scopes       []string `json:"scopes"`
		AccountName  string   `json:"accountName"` //nolint:tagliatelle // decodes the contract's camelCase JSON
	}

	session := testsupport.SessionFor(t, o.DB, alice)
	res := o.Do(http.MethodGet, "/api/oauth/consents/"+flow.RequestID, nil, func(r *http.Request) { r.AddCookie(session.Cookie) })
	require.Equal(t, http.StatusOK, res.StatusCode)
	decode(t, res, &consent)
	assert.Equal(t, "Test MCP client", consent.ClientName)
	assert.Equal(t, "127.0.0.1:8976", consent.RedirectHost)
	assert.ElementsMatch(t, []string{"catalog:read", "progress:read", "offline_access"}, consent.Scopes)
	assert.NotEmpty(t, consent.AccountName)

	// Approving yields a code, and the code yields tokens.
	code := flow.Approve()
	tokens, res := flow.Exchange(code, testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "Bearer", tokens.TokenType)
	assert.Positive(t, tokens.ExpiresIn)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.ElementsMatch(t, []string{"catalog:read", "progress:read", "offline_access"}, tokens.Scopes)

	granted, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, alice, granted.UserID)
	assert.Equal(t, testsupport.OAuthClientID, granted.ClientID)
	assert.ElementsMatch(t, []string{"catalog:read", "progress:read", "offline_access"}, granted.Scopes)
}

func TestVerificationNamesTheConnectedAppByItsRowAndRefusesATokenOfARemovedClient(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	app := o.DB.OAuthClient.Query().OnlyX(t.Context())

	verified, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, app.ID, verified.ConnectedAppID)
	assert.NotEmpty(t, verified.ID)

	o.DB.OAuthClient.DeleteOne(app).ExecX(t.Context())
	_, err = o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestDenyingTheConsentReturnsAccessDeniedToTheClient(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice, oauthserver.ScopeCatalogRead)

	res := flow.Decide(false)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var decision struct {
		RedirectURL string `json:"redirectUrl"` //nolint:tagliatelle // decodes the contract's camelCase JSON
	}
	decode(t, res, &decision)

	back, err := url.Parse(decision.RedirectURL)
	require.NoError(t, err)
	assert.Equal(t, testsupport.OAuthRedirectURI, back.Scheme+"://"+back.Host+back.Path)
	assert.Equal(t, "access_denied", back.Query().Get("error"))
	assert.Equal(t, flow.State, back.Query().Get("state"))
	assert.Empty(t, back.Query().Get("code"))

	// The request is gone: it cannot be approved afterwards.
	res = flow.Decide(true)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestConsentOperationsNeedTheSessionAndRefuseCrossSiteWrites(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice, oauthserver.ScopeCatalogRead)
	session := testsupport.SessionFor(t, o.DB, alice)
	target := "/api/oauth/consents/" + flow.RequestID

	// Signed out: no data, no decision.
	assert.Equal(t, http.StatusUnauthorized, o.Do(http.MethodGet, target, nil).StatusCode)
	assert.Equal(t, http.StatusUnauthorized, o.Do(http.MethodPost, target+"/approve", nil).StatusCode)

	// A browser's cross-site write cannot decide, session or not (ADR-0017:
	// cross-origin protection replaces the XSRF token).
	crossSite := func(r *http.Request) {
		r.AddCookie(session.Cookie)
		r.Header.Set("Sec-Fetch-Site", "cross-site")
	}
	res := o.Do(http.MethodPost, target+"/approve", nil, crossSite)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
	res = o.Do(http.MethodPost, target+"/deny", nil, crossSite)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	// And nothing was decided by those attempts.
	assert.Equal(t, http.StatusOK, o.Do(http.MethodGet, target, nil, func(r *http.Request) { r.AddCookie(session.Cookie) }).StatusCode)
}

func TestAnUnknownOrExpiredRequestIsNotFound(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice, oauthserver.ScopeCatalogRead)
	session := testsupport.SessionFor(t, o.DB, alice)
	get := func(id string) int {
		return o.Do(http.MethodGet, "/api/oauth/consents/"+id, nil, func(r *http.Request) { r.AddCookie(session.Cookie) }).StatusCode
	}

	assert.Equal(t, http.StatusNotFound, get("no-such-request"))
	assert.Equal(t, http.StatusOK, get(flow.RequestID))

	o.Advance(11 * time.Minute)
	assert.Equal(t, http.StatusNotFound, get(flow.RequestID))
	assert.Equal(t, http.StatusNotFound, flow.Decide(true).StatusCode)
}

func TestAuthorizationWithoutPKCES256IsRefused(t *testing.T) {
	cases := map[string]url.Values{
		"no challenge":    {},
		"plain method":    {"code_challenge": {testsupport.PKCEChallenge()}, "code_challenge_method": {"plain"}},
		"default method":  {"code_challenge": {testsupport.PKCEChallenge()}},
		"challenge only":  {"code_challenge": {"x"}, "code_challenge_method": {""}},
		"unknown method":  {"code_challenge": {testsupport.PKCEChallenge()}, "code_challenge_method": {"S512"}},
		"method no value": {"code_challenge_method": {"S256"}},
	}
	for name, pkce := range cases {
		t.Run(name, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			o.SeedClient(oauthserver.ScopeCatalogRead)

			query := url.Values{
				"client_id":     {testsupport.OAuthClientID},
				"redirect_uri":  {testsupport.OAuthRedirectURI},
				"response_type": {"code"},
				"scope":         {"catalog:read"},
				"state":         {"s"},
			}
			maps.Copy(query, pkce)

			res := o.Do(http.MethodGet, "/oauth/authorize?"+query.Encode(), nil)

			// A registered redirect URI gets the error; either way no consent
			// screen is reached and no request is stored.
			require.Equal(t, http.StatusFound, res.StatusCode)
			to, err := res.Location()
			require.NoError(t, err)
			assert.NotContains(t, to.Path, "/connect")
			assert.Equal(t, "invalid_request", to.Query().Get("error"))
			count, err := o.DB.OAuthAuthRequest.Query().Count(t.Context())
			require.NoError(t, err)
			assert.Zero(t, count)
		})
	}
}

func TestAnAuthorizationCodeWorksOnceAndOnlyWithItsVerifier(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice, oauthserver.ScopeCatalogRead)
	code := flow.Approve()

	_, res := flow.Exchange(code, "a-different-verifier-that-is-long-enough-0123456789-0123456789")
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "a wrong verifier is refused")

	_, res = flow.Exchange(code, testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, res.StatusCode, "the failed attempt did not burn the code")

	_, res = flow.Exchange(code, testsupport.PKCEVerifier())
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "a code is single-use")
}

func TestTokensAreOpaqueAndTheDatabaseHoldsNoUsableToken(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	// Opaque: not a JWT, no readable subject in it.
	// The library mints the bearer as an encrypted compact JWE (five parts, an
	// AES-256-GCM payload), not a signed JWT whose claims anyone can read.
	parts := strings.Split(tokens.AccessToken, ".")
	require.Len(t, parts, 5, "an access token is an encrypted JWE")
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	require.NoError(t, err)
	assert.Contains(t, string(header), "A256GCM")
	assert.NotContains(t, string(header), "sub")
	assert.NotContains(t, tokens.AccessToken, strconv.Itoa(alice))

	// The row holds an id. Presenting that id as a bearer token is refused.
	row, err := o.DB.OAuthAccessToken.Query().Only(t.Context())
	require.NoError(t, err)
	assert.NotContains(t, tokens.AccessToken, row.ID)
	_, err = o.Server.VerifyAccessToken(t.Context(), row.ID)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
	assert.NotEqual(t, tokens.AccessToken, row.ID)

	// A refresh token is kept only as a hash: the value the client holds is not
	// in the table, and its hash is not a token.
	count, err := o.DB.OAuthRefreshToken.Query().Where(oauthrefreshtoken.TokenHash(tokens.RefreshToken)).Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, count)
	stored, err := o.DB.OAuthRefreshToken.Query().Only(t.Context())
	require.NoError(t, err)

	res := o.Do(http.MethodPost, "/oauth/token", url.Values{
		"grant_type": {"refresh_token"}, "client_id": {testsupport.OAuthClientID}, "refresh_token": {stored.TokenHash},
	})
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "the stored hash cannot be presented as the refresh token")

	// A tampered token is refused.
	_, err = o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken+"x")
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestVerificationIsALookupSoDeletingTheRowRevokesAtOnce(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)

	_, err = o.DB.OAuthAccessToken.Delete().Where(oauthaccesstoken.UserID(alice)).Exec(t.Context())
	require.NoError(t, err)

	_, err = o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestAnAccessTokenExpiresAfterAnHour(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	o.Advance(59 * time.Minute)
	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)

	o.Advance(2 * time.Minute)
	_, err = o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestRefreshExchangesTheTokenForANewPair(t *testing.T) {
	o := testsupport.NewOAuth(t)
	first := o.Connect(alice, oauthserver.ScopeCatalogRead)
	refresh := func(token string) (testsupport.Tokens, int) {
		res := o.Do(http.MethodPost, "/oauth/token", url.Values{
			"grant_type": {"refresh_token"}, "client_id": {testsupport.OAuthClientID}, "refresh_token": {token},
		})

		var tokens testsupport.Tokens
		if res.StatusCode == http.StatusOK {
			decode(t, res, &tokens)
		}

		return tokens, res.StatusCode
	}

	second, status := refresh(first.RefreshToken)
	require.Equal(t, http.StatusOK, status)
	assert.NotEqual(t, first.RefreshToken, second.RefreshToken)
	_, err := o.Server.VerifyAccessToken(t.Context(), second.AccessToken)
	require.NoError(t, err)

	// The old pair is spent: the refresh token is refused and its access token
	// no longer verifies.
	_, status = refresh(first.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, status)
	_, err = o.Server.VerifyAccessToken(t.Context(), first.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestRevokingATokenEndsItAtOnce(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	res := o.Do(http.MethodPost, "/oauth/revoke", url.Values{
		"client_id": {testsupport.OAuthClientID}, "token": {tokens.AccessToken}, "token_type_hint": {"access_token"},
	})
	require.Equal(t, http.StatusOK, res.StatusCode)

	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

package oauthserver_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

func TestTheKeySetPublishesTheKeyIDTokensAreSignedWith(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), "openid", oauthserver.ScopeCatalogRead)
	require.NotEmpty(t, tokens.IDToken, "an openid request gets an ID token")

	res := o.Do(http.MethodGet, "/oauth/keys", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var set struct {
		Keys []struct {
			Kid string `json:"kid"`
			Alg string `json:"alg"`
			Use string `json:"use"`
		} `json:"keys"`
	}
	decode(t, res, &set)
	require.Len(t, set.Keys, 1)
	assert.Equal(t, "EdDSA", set.Keys[0].Alg)
	assert.Equal(t, "sig", set.Keys[0].Use)

	// The ID token names that key and algorithm.
	header, err := base64.RawURLEncoding.DecodeString(strings.Split(tokens.IDToken, ".")[0])
	require.NoError(t, err)

	var parsed struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	require.NoError(t, json.Unmarshal(header, &parsed))
	assert.Equal(t, "EdDSA", parsed.Alg)
	assert.Equal(t, set.Keys[0].Kid, parsed.Kid)

	payload, err := base64.RawURLEncoding.DecodeString(strings.Split(tokens.IDToken, ".")[1])
	require.NoError(t, err)

	var claims struct {
		Sub string `json:"sub"`
		Iss string `json:"iss"`
		Azp string `json:"azp"`
	}
	require.NoError(t, json.Unmarshal(payload, &claims))
	assert.Equal(t, strconv.Itoa(alice()), claims.Sub)
	assert.Equal(t, testsupport.OAuthIssuer, claims.Iss)
	assert.Equal(t, testsupport.OAuthClientID, claims.Azp)
}

func TestUserinfoAnswersWithTheSubjectOfALiveToken(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	bearer := func(r *http.Request) { r.Header.Set("Authorization", "Bearer "+tokens.AccessToken) }

	res := o.Do(http.MethodGet, "/oauth/userinfo", nil, bearer)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var info struct {
		Sub string `json:"sub"`
	}
	decode(t, res, &info)
	assert.Equal(t, strconv.Itoa(alice()), info.Sub)

	// Once the token is revoked the endpoint refuses it.
	_, err := o.DB.OAuthAccessToken.Delete().Exec(t.Context())
	require.NoError(t, err)

	res = o.Do(http.MethodGet, "/oauth/userinfo", nil, bearer)
	assert.NotEqual(t, http.StatusOK, res.StatusCode)
}

func TestRevokingARefreshTokenEndsItsPairToo(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)

	res := o.Do(http.MethodPost, "/oauth/revoke", url.Values{
		"client_id": {testsupport.OAuthClientID}, "token": {tokens.RefreshToken}, "token_type_hint": {"refresh_token"},
	})
	require.Equal(t, http.StatusOK, res.StatusCode)

	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)

	res = o.Do(http.MethodPost, "/oauth/token", url.Values{
		"grant_type": {"refresh_token"}, "client_id": {testsupport.OAuthClientID}, "refresh_token": {tokens.RefreshToken},
	})
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestAClientCannotRevokeAnotherClientsToken(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	// Another client registers itself, the way a client does.
	res := o.Register(request(nil))
	require.Equal(t, http.StatusCreated, res.StatusCode)
	other := decodeJSON[registered](t, res)

	for _, token := range []string{tokens.AccessToken, tokens.RefreshToken} {
		res := o.Do(http.MethodPost, "/oauth/revoke", url.Values{"client_id": {other.ClientID}, "token": {token}})
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	}

	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err, "the token survived")
}

func TestEndingTheSessionRevokesTheClientsGrants(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), "openid", oauthserver.ScopeCatalogRead)

	res := o.Do(http.MethodGet, "/oauth/end_session?"+url.Values{"id_token_hint": {tokens.IDToken}}.Encode(), nil)
	body, _ := io.ReadAll(res.Body)
	require.Less(t, res.StatusCode, 400, "end_session answered %d: %s", res.StatusCode, body)

	_, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.ErrorIs(t, err, oauthserver.ErrInvalidToken)
}

func TestAuthorizationRefusesWhatCannotBeServed(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.SeedClient(oauthserver.ScopeCatalogRead)

	base := func() url.Values {
		return url.Values{
			"client_id": {testsupport.OAuthClientID}, "redirect_uri": {testsupport.OAuthRedirectURI},
			"response_type": {"code"}, "scope": {"catalog:read"}, "state": {"s"},
			"code_challenge": {testsupport.PKCEChallenge()}, "code_challenge_method": {"S256"},
		}
	}

	// prompt=none: nobody can sign in without a screen, so there is no silent grant.
	values := base()
	values.Set("prompt", "none")
	res := o.Do(http.MethodGet, "/oauth/authorize?"+values.Encode(), nil)
	to, err := res.Location()
	require.NoError(t, err)
	assert.Equal(t, "login_required", to.Query().Get("error"))

	// An unknown client has no redirect URI to trust, so there is nobody to send the error to.
	values = base()
	values.Set("client_id", "nobody")
	res = o.Do(http.MethodGet, "/oauth/authorize?"+values.Encode(), nil)
	assert.NotEqual(t, http.StatusFound, res.StatusCode)

	// The implicit flow is not offered.
	values = base()
	values.Set("response_type", "id_token token")

	res = o.Do(http.MethodGet, "/oauth/authorize?"+values.Encode(), nil)

	loc, err := res.Location()
	if err == nil {
		assert.NotContains(t, loc.Path, "/connect")
	}

	count, err := o.DB.OAuthAuthRequest.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, count)
}

func TestACallbackIssuesOneCodeOnly(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice(), oauthserver.ScopeCatalogRead)
	res := flow.Decide(true)

	var decision struct {
		RedirectURL string `json:"redirectUrl"` //nolint:tagliatelle // decodes the contract's camelCase JSON
	}
	decode(t, res, &decision)
	callback, err := url.Parse(decision.RedirectURL)
	require.NoError(t, err)

	first := o.Do(http.MethodGet, callback.RequestURI(), nil)
	to, err := first.Location()
	require.NoError(t, err)
	assert.NotEmpty(t, to.Query().Get("code"))

	// Opening the same callback again cannot mint a second code.
	second := o.Do(http.MethodGet, callback.RequestURI(), nil)

	to, err = second.Location()
	if err == nil {
		assert.Empty(t, to.Query().Get("code"))
	}
}

func TestAnUnapprovedRequestCannotBeExchangedForACode(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice(), oauthserver.ScopeCatalogRead)

	res := o.Do(http.MethodGet, "/oauth/authorize/callback?id="+flow.RequestID, nil)

	to, err := res.Location()
	require.NoError(t, err)
	assert.Empty(t, to.Query().Get("code"))
	assert.Equal(t, "interaction_required", to.Query().Get("error"))
}

func TestApprovingTwiceIsRefused(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice(), oauthserver.ScopeCatalogRead)

	assert.Equal(t, http.StatusOK, flow.Decide(true).StatusCode)
	assert.Equal(t, http.StatusNotFound, flow.Decide(true).StatusCode)
}

func TestTheConsentScreenNamesTheAccount(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice(), oauthserver.ScopeCatalogRead)
	session := testsupport.SessionFor(t, o.DB, alice())
	accountName := func() string {
		res := o.Do(http.MethodGet, "/api/oauth/consents/"+flow.RequestID, nil, func(r *http.Request) { r.AddCookie(session.Cookie) })
		require.Equal(t, http.StatusOK, res.StatusCode)

		var consent struct {
			AccountName string `json:"accountName"` //nolint:tagliatelle // decodes the contract's camelCase JSON
		}
		decode(t, res, &consent)

		return consent.AccountName
	}

	_, err := o.DB.User.UpdateOneID(alice()).SetFirstName("Ada").SetLastName("Lovelace").Save(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "Ada Lovelace", accountName())

	// Without a name the email stands in.
	_, err = o.DB.User.UpdateOneID(alice()).ClearFirstName().ClearLastName().Save(t.Context())
	require.NoError(t, err)
	u, err := o.DB.User.Get(t.Context(), alice())
	require.NoError(t, err)
	assert.Equal(t, *u.Email, accountName())
}

func TestConsentOperationsOfTheServerNeedASignedInUser(t *testing.T) {
	o := testsupport.NewOAuth(t)
	flow := o.Start(alice(), oauthserver.ScopeCatalogRead)
	ctx := context.Background()

	_, err := o.Server.Pending(ctx, flow.RequestID)
	require.ErrorIs(t, err, oauthserver.ErrNotSignedIn)
	_, err = o.Server.Approve(ctx, flow.RequestID)
	require.ErrorIs(t, err, oauthserver.ErrNotSignedIn)
	_, err = o.Server.Deny(ctx, flow.RequestID)
	require.ErrorIs(t, err, oauthserver.ErrNotSignedIn)
}

func TestTheServerRefusesAConfigurationItCannotRunWith(t *testing.T) {
	_, err := oauthserver.New(oauthserver.Deps{Config: oauthserver.Config{Issuer: "https://code-basics.test",
		EncryptionKey:        "",
		RegistrationLimit:    0,
		RegistrationWindow:   0,
		AccessTokenLifetime:  0,
		RefreshTokenLifetime: 0,
	},
		DB:            nil,
		Transactor:    nil,
		SignedIn:      nil,
		ClientAddress: nil,
		Now:           nil,
	})
	require.ErrorContains(t, err, "encryption key")

	_, err = oauthserver.New(oauthserver.Deps{
		Config: oauthserver.Config{
			EncryptionKey:        "k",
			Issuer:               "not an origin",
			RegistrationLimit:    0,
			RegistrationWindow:   0,
			AccessTokenLifetime:  0,
			RefreshTokenLifetime: 0,
		},
		DB:            nil,
		Transactor:    nil,
		SignedIn:      nil,
		ClientAddress: nil,
		Now:           nil,
	})
	require.Error(t, err)
}

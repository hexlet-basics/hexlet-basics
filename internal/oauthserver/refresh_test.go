package oauthserver_test

import (
	"net/http"
	"net/url"
	"sync"
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

// bob is the fixtures' editor, a second user whose grants must stay out of
// alice's way.
func bob() int {
	return fixture.ID(fixture.UserStaffEditor)
}

func TestRefreshRotatesTheTokenAndRefusesTheOldOne(t *testing.T) {
	o := testsupport.NewOAuth(t)
	first := o.Connect(alice(), oauthserver.ScopeCatalogRead)

	second, res := o.Refresh(testsupport.OAuthClientID, first.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, first.RefreshToken, second.RefreshToken)
	assert.True(t, o.Accepts(second.AccessToken))
	assert.False(t, o.Accepts(first.AccessToken), "the access token of the spent pair is gone")

	_, res = o.Refresh(testsupport.OAuthClientID, first.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestReplayingAUsedRefreshTokenRevokesTheWholeFamily(t *testing.T) {
	o := testsupport.NewOAuth(t)
	first := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	second, res := o.Refresh(testsupport.OAuthClientID, first.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	third, res := o.Refresh(testsupport.OAuthClientID, second.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.True(t, o.Accepts(third.AccessToken))

	// A thief replays the first token, long since spent.
	_, res = o.Refresh(testsupport.OAuthClientID, first.RefreshToken)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	// Everything of the family is dead, the newest tokens included: the victim is
	// locked out together with the thief and must sign in again.
	assert.False(t, o.Accepts(third.AccessToken))
	_, res = o.Refresh(testsupport.OAuthClientID, third.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	count, err := o.DB.OAuthRefreshToken.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, count)
}

func TestARevokedFamilyLeavesOtherGrantsAlone(t *testing.T) {
	o := testsupport.NewOAuth(t)
	victim := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	laterSignIn := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	otherUser := o.Connect(bob(), oauthserver.ScopeCatalogRead)

	// Another client of the same user, which registered itself.
	res := o.Register(request(func(m map[string]any) {
		m["redirect_uris"] = []string{"http://127.0.0.1:1/cb"}
		m["scope"] = oauthserver.ScopeCatalogRead + " offline_access"
	}))
	require.Equal(t, http.StatusCreated, res.StatusCode)
	other := decodeJSON[registered](t, res)
	flow := o.StartWith(alice(), testsupport.StartOptions{
		Scopes: []string{oauthserver.ScopeCatalogRead}, ClientID: other.ClientID, RedirectURIs: []string{"http://127.0.0.1:1/cb"},
	})
	otherClient, res := flow.Exchange(flow.Approve(), testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, res.StatusCode)

	rotated, res := o.Refresh(testsupport.OAuthClientID, victim.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	_, res = o.Refresh(testsupport.OAuthClientID, victim.RefreshToken) // replay
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.False(t, o.Accepts(rotated.AccessToken))

	for name, tokens := range map[string]testsupport.Tokens{
		"another sign-in of the same user": laterSignIn, "another user": otherUser,
	} {
		assert.True(t, o.Accepts(tokens.AccessToken), name)
		_, res = o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode, name)
	}

	assert.True(t, o.Accepts(otherClient.AccessToken), "another client of the same user")
	_, res = o.Refresh(other.ClientID, otherClient.RefreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestTheFamilySurvivesSeveralRotations(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)

	spent := make([]string, 0, 5)
	for range 5 {
		spent = append(spent, tokens.RefreshToken)
		next, res := o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
		require.Equal(t, http.StatusOK, res.StatusCode)

		tokens = next
	}

	families, err := o.DB.OAuthRefreshToken.Query().Select(oauthrefreshtoken.FieldFamilyID).Strings(t.Context())
	require.NoError(t, err)
	require.Len(t, families, 6, "five spent tokens stay on file beside the live one")

	for _, f := range families {
		assert.Equal(t, families[0], f, "one authorization, one family")
	}

	// Replaying the very first token still reaches back to the newest one.
	_, res := o.Refresh(testsupport.OAuthClientID, spent[0])
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.False(t, o.Accepts(tokens.AccessToken))
	_, res = o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestRevokingARefreshTokenEndsTheFamilyAndAnAccessTokenOnlyItself(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	rotated, res := o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)

	// An access token goes at once; the client may still refresh.
	require.Equal(t, http.StatusOK, o.Revoke(testsupport.OAuthClientID, rotated.AccessToken).StatusCode)
	assert.False(t, o.Accepts(rotated.AccessToken))
	again, res := o.Refresh(testsupport.OAuthClientID, rotated.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)

	// A refresh token takes the whole connection with it.
	require.Equal(t, http.StatusOK, o.Revoke(testsupport.OAuthClientID, again.RefreshToken).StatusCode)
	assert.False(t, o.Accepts(again.AccessToken))
	_, res = o.Refresh(testsupport.OAuthClientID, again.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// Revoking what is already gone is not an error (RFC 7009).
	assert.Equal(t, http.StatusOK, o.Revoke(testsupport.OAuthClientID, again.RefreshToken).StatusCode)
}

func TestTokensExpireOnTheirOwnLifetimes(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	// The library counts expires_in against the real clock, not the test's.
	assert.InDelta(t, 3600, tokens.ExpiresIn, 5, "an access token lives an hour")

	o.Advance(61 * time.Minute)
	assert.False(t, o.Accepts(tokens.AccessToken))
	// The refresh token outlives it: refreshing buys a fresh pair.
	refreshed, res := o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, o.Accepts(refreshed.AccessToken))

	// Each rotation grants a full thirty days again: 29 days later the rotated
	// token is still good, and a refresh token expires 30 days after its issue.
	o.Advance(29 * 24 * time.Hour)
	again, res := o.Refresh(testsupport.OAuthClientID, refreshed.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	o.Advance(30*24*time.Hour + time.Minute)
	_, res = o.Refresh(testsupport.OAuthClientID, again.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "a refresh token expires after thirty days")
}

func TestConfiguredLifetimesApply(t *testing.T) {
	o := testsupport.NewOAuthWith(t, oauthserver.Config{
		AccessTokenLifetime: 5 * time.Minute, RefreshTokenLifetime: time.Hour,
	})
	tokens := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	assert.InDelta(t, 300, tokens.ExpiresIn, 5)

	o.Advance(6 * time.Minute)
	assert.False(t, o.Accepts(tokens.AccessToken))
	o.Advance(time.Hour)
	_, res := o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestPurgeDeletesWhatExpiredAndKeepsTheRest(t *testing.T) {
	o := testsupport.NewOAuth(t)
	old := o.Connect(alice(), oauthserver.ScopeCatalogRead)
	_, res := o.Refresh(testsupport.OAuthClientID, old.RefreshToken) // leaves a used tombstone
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Nothing has expired yet.
	n, err := oauthserver.PurgeExpired(t.Context(), o.DB, o.Now())
	require.NoError(t, err)
	assert.Zero(t, n)

	o.Advance(31 * 24 * time.Hour)
	fresh := o.Connect(bob(), oauthserver.ScopeCatalogRead)
	n, err = oauthserver.PurgeExpired(t.Context(), o.DB, o.Now())
	require.NoError(t, err)
	assert.Equal(t, 3, n, "the rotated pair's access token, and both refresh tokens (one used) of the old family")

	access, err := o.DB.OAuthAccessToken.Query().Where(oauthaccesstoken.UserID(bob())).Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, access)
	refresh, err := o.DB.OAuthRefreshToken.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, refresh)
	assert.True(t, o.Accepts(fresh.AccessToken))
}

// TestParallelRefreshesNeverBothSucceed runs on a real connection pool, not the
// rolled-back test transaction, because only separate transactions race for the
// row. Of N exchanges of one token at most one may succeed, and since every
// loser reads as a replay the family ends revoked: there are never two live
// children.
//
// The client registers itself and the person connects through the real flow;
// the pool commits all of it into a database of the test's own, dropped with it.
func TestParallelRefreshesNeverBothSucceed(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.OnPool())
	res := o.Register(request(func(m map[string]any) {
		m["redirect_uris"] = []string{"http://127.0.0.1:1/cb"}
		m["scope"] = oauthserver.ScopeCatalogRead + " offline_access"
	}))
	require.Equal(t, http.StatusCreated, res.StatusCode)
	racer := decodeJSON[registered](t, res)
	flow := o.StartWith(alice(), testsupport.StartOptions{
		Scopes: []string{oauthserver.ScopeCatalogRead}, ClientID: racer.ClientID, RedirectURIs: []string{"http://127.0.0.1:1/cb"},
	})
	tokens, res := flow.Exchange(flow.Approve(), testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, res.StatusCode)

	const attempts = 8

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	succeeded := 0

	for range attempts {
		wg.Go(func() {
			res := o.Do(http.MethodPost, "/oauth/token", url.Values{
				"grant_type": {"refresh_token"}, "client_id": {racer.ClientID}, "refresh_token": {tokens.RefreshToken},
			})

			_ = res.Body.Close()
			if res.StatusCode == http.StatusOK {
				mu.Lock()
				succeeded++
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	ctx := t.Context()
	assert.LessOrEqual(t, succeeded, 1, "one token buys at most one child")

	live := o.DB.OAuthRefreshToken.Query().Where(oauthrefreshtoken.ClientID(racer.ClientID), oauthrefreshtoken.UsedAtIsNil()).CountX(ctx)
	assert.Zero(t, live, "a losing exchange revoked the family, winner's child included")
	assert.Zero(t, o.DB.OAuthAccessToken.Query().Where(oauthaccesstoken.ClientID(racer.ClientID)).CountX(ctx))
}

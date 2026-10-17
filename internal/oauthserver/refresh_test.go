package oauthserver_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/oauthaccesstoken"
	"hexletbasics/ent/oauthrefreshtoken"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
	"hexletbasics/internal/testsupport/testdb"
)

// bob is the fixtures' editor, a second user whose grants must stay out of
// alice's way.
var bob = fixtures.ID("user_staff_editor")

func TestRefreshRotatesTheTokenAndRefusesTheOldOne(t *testing.T) {
	o := testsupport.NewOAuth(t)
	first := o.Connect(alice, oauthserver.ScopeCatalogRead)

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
	first := o.Connect(alice, oauthserver.ScopeCatalogRead)
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
	victim := o.Connect(alice, oauthserver.ScopeCatalogRead)
	laterSignIn := o.Connect(alice, oauthserver.ScopeCatalogRead)
	otherUser := o.Connect(bob, oauthserver.ScopeCatalogRead)

	// Another client of the same user.
	//nolint:forbidigo // the client or token under test: its lifetime, secrets and ownership are what the test chooses
	_, err := o.DB.OAuthClient.Create().SetClientID("other-client").SetName("Other").
		SetRedirectUris([]string{"http://127.0.0.1:1/cb"}).SetScopes([]string{oauthserver.ScopeCatalogRead}).Save(t.Context())
	require.NoError(t, err)
	flow := o.StartWith(alice, testsupport.StartOptions{
		Scopes: []string{oauthserver.ScopeCatalogRead}, ClientID: "other-client", RedirectURIs: []string{"http://127.0.0.1:1/cb"},
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
	_, res = o.Refresh("other-client", otherClient.RefreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestTheFamilySurvivesSeveralRotations(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	var spent []string
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
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
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
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
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
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	assert.InDelta(t, 300, tokens.ExpiresIn, 5)

	o.Advance(6 * time.Minute)
	assert.False(t, o.Accepts(tokens.AccessToken))
	o.Advance(time.Hour)
	_, res := o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestPurgeDeletesWhatExpiredAndKeepsTheRest(t *testing.T) {
	o := testsupport.NewOAuth(t)
	old := o.Connect(alice, oauthserver.ScopeCatalogRead)
	_, res := o.Refresh(testsupport.OAuthClientID, old.RefreshToken) // leaves a used tombstone
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Nothing has expired yet.
	n, err := oauthserver.PurgeExpired(t.Context(), o.DB, o.Now())
	require.NoError(t, err)
	assert.Zero(t, n)

	o.Advance(31 * 24 * time.Hour)
	fresh := o.Connect(bob, oauthserver.ScopeCatalogRead)
	n, err = oauthserver.PurgeExpired(t.Context(), o.DB, o.Now())
	require.NoError(t, err)
	assert.Equal(t, 3, n, "the rotated pair's access token, and both refresh tokens (one used) of the old family")

	access, err := o.DB.OAuthAccessToken.Query().Where(oauthaccesstoken.UserID(bob)).Count(t.Context())
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
func TestParallelRefreshesNeverBothSucceed(t *testing.T) {
	//nolint:forbidigo // the case is about separate connections (concurrent transactions), so the test needs a pool of its own
	sqlDB, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)
	ctx := t.Context()
	server, err := oauthserver.New(oauthserver.Deps{
		Config:     oauthserver.Config{Issuer: testsupport.OAuthIssuer, EncryptionKey: "test-oauth-encryption-key-test-only-test-only"},
		DB:         db,
		Transactor: store.New(sqlDB),
	})
	require.NoError(t, err)

	//nolint:forbidigo // the client or token under test: its lifetime, secrets and ownership are what the test chooses
	client := db.OAuthClient.Create().SetClientID("race-client").SetName("Race").
		SetRedirectUris([]string{"http://127.0.0.1:1/cb"}).SetScopes([]string{}).SaveX(ctx)
	t.Cleanup(func() {
		// Deleting the client cascades its tokens away.
		_ = db.OAuthClient.DeleteOne(client).Exec(context.Background())
		_ = sqlDB.Close()
	})

	const secret = "race-refresh-secret"
	sum := sha256.Sum256([]byte(secret))
	//nolint:forbidigo // the client or token under test: its lifetime, secrets and ownership are what the test chooses
	db.OAuthAccessToken.Create().SetID("race-access").SetClientID("race-client").SetUserID(alice).
		SetScopes([]string{}).SetRefreshTokenID("race-refresh").SetFamilyID("race-family").
		SetExpiresAt(time.Now().Add(time.Hour)).SaveX(ctx)
	//nolint:forbidigo // the client or token under test: its lifetime, secrets and ownership are what the test chooses
	db.OAuthRefreshToken.Create().SetID("race-refresh").SetTokenHash(hex.EncodeToString(sum[:])).
		SetClientID("race-client").SetUserID(alice).SetScopes([]string{}).SetAuthTime(time.Now()).
		SetAccessTokenID("race-access").SetFamilyID("race-family").SetExpiresAt(time.Now().Add(time.Hour)).SaveX(ctx)

	const attempts = 8
	var wg sync.WaitGroup
	var mu sync.Mutex
	succeeded := 0
	for range attempts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/oauth/token", strings.NewReader(url.Values{
				"grant_type": {"refresh_token"}, "client_id": {"race-client"}, "refresh_token": {secret},
			}.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			server.Handler().ServeHTTP(rec, req)
			if rec.Code == http.StatusOK {
				mu.Lock()
				succeeded++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	assert.LessOrEqual(t, succeeded, 1, "one token buys at most one child")
	live := db.OAuthRefreshToken.Query().Where(oauthrefreshtoken.FamilyID("race-family"), oauthrefreshtoken.UsedAtIsNil()).CountX(ctx)
	assert.Zero(t, live, "a losing exchange revoked the family, winner's child included")
	assert.Zero(t, db.OAuthAccessToken.Query().Where(oauthaccesstoken.FamilyID("race-family")).CountX(ctx))
}

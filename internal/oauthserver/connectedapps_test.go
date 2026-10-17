package oauthserver_test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// connectedApp mirrors one item of the list operation.
type connectedApp struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Scopes      []string  `json:"scopes"`
	ConnectedAt time.Time `json:"connectedAt"`
	ClientHost  string    `json:"clientHost"`
}

// listApps calls the contract operation as a signed-in browser.
func listApps(t *testing.T, o *testsupport.OAuth, userID int) []connectedApp {
	t.Helper()
	res := o.Do(http.MethodGet, "/api/oauth/connected_apps", nil, func(r *http.Request) {
		r.AddCookie(testsupport.SessionFor(t, o.DB, userID).Cookie)
	})
	require.Equal(t, http.StatusOK, res.StatusCode)
	return decodeJSON[struct {
		Items []connectedApp `json:"items"`
	}](t, res).Items
}

// revokeApp calls the revoke operation; xsrf false leaves the header out.
func revokeApp(t *testing.T, o *testsupport.OAuth, userID, appID int, xsrf bool) *http.Response {
	t.Helper()
	s := testsupport.SessionFor(t, o.DB, userID)
	return o.Do(http.MethodDelete, "/api/oauth/connected_apps/"+strconv.Itoa(appID), nil, func(r *http.Request) {
		r.AddCookie(s.Cookie)
		if xsrf {
			r.Header.Set("X-XSRF-TOKEN", s.XSRF)
		}
	})
}

func TestTheListShowsTheAppsOfTheSignedInUserOnly(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Connect(alice, oauthserver.ScopeCatalogRead, oauthserver.ScopeProgressRead)

	mine := listApps(t, o, alice)
	require.Len(t, mine, 1)
	assert.Equal(t, "Test MCP client", mine[0].Name)
	assert.ElementsMatch(t, []string{"catalog:read", "progress:read", "offline_access"}, mine[0].Scopes)
	assert.WithinDuration(t, o.Now(), mine[0].ConnectedAt, time.Minute)
	assert.Empty(t, mine[0].ClientHost, "a client with an opaque id has no host to show")

	assert.Empty(t, listApps(t, o, bob), "another user's apps are not mine")
}

func TestTheListIsAnEmptyArrayWithoutApps(t *testing.T) {
	o := testsupport.NewOAuth(t)

	assert.NotNil(t, listApps(t, o, alice))
	assert.Empty(t, listApps(t, o, alice))
}

func TestSeveralSignInsToOneClientAreOneApp(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Connect(alice, oauthserver.ScopeCatalogRead)
	o.Advance(time.Hour)
	o.Connect(alice, oauthserver.ScopeCatalogRead)

	assert.Len(t, listApps(t, o, alice), 1)
}

func TestAnAppOfAClientDocumentShowsTheHostOfItsClientID(t *testing.T) {
	o := testsupport.NewOAuth(t)
	doc := testsupport.ValidDocument(documentRedirect)
	doc["scope"] = "catalog:read offline_access"
	id := o.Documents.Publish("/oauth/client.json", doc)
	f := authorizeDocument(o, id, documentRedirect, "catalog:read")
	_, res := f.Exchange(f.Approve(), testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, res.StatusCode)

	apps := listApps(t, o, alice)

	require.Len(t, apps, 1)
	assert.Equal(t, "Document client", apps[0].Name)
	assert.Equal(t, testsupport.DocumentsHost, apps[0].ClientHost)
}

func TestRevokingAnAppEndsItsAccessAndRefreshTokensAtOnce(t *testing.T) {
	o := testsupport.NewOAuth(t)
	first := o.Connect(alice, oauthserver.ScopeCatalogRead)
	// A rotation leaves a used tombstone behind, and a second sign-in is a second
	// family: the revoke ends all of it.
	second, res := o.Refresh(testsupport.OAuthClientID, first.RefreshToken)
	require.Equal(t, http.StatusOK, res.StatusCode)
	otherFamily := o.Connect(alice, oauthserver.ScopeCatalogRead)
	require.True(t, o.Accepts(second.AccessToken))
	require.True(t, o.Accepts(otherFamily.AccessToken))

	apps := listApps(t, o, alice)
	require.Len(t, apps, 1)
	res = revokeApp(t, o, alice, apps[0].ID, true)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	assert.False(t, o.Accepts(second.AccessToken), "the token that worked a moment ago is refused")
	assert.False(t, o.Accepts(otherFamily.AccessToken))
	for _, refresh := range []string{second.RefreshToken, first.RefreshToken, otherFamily.RefreshToken} {
		_, res = o.Refresh(testsupport.OAuthClientID, refresh)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "no refresh token survives, used ones included")
	}
	assert.Empty(t, listApps(t, o, alice))
}

func TestRevokingAnAppLeavesOtherUsersAndOtherAppsAlone(t *testing.T) {
	o := testsupport.NewOAuth(t)
	mine := o.Connect(alice, oauthserver.ScopeCatalogRead)
	theirs := o.Connect(bob, oauthserver.ScopeCatalogRead)
	other, _ := connectOtherClient(t, o, alice)

	apps := listApps(t, o, alice)
	require.Len(t, apps, 2)
	var target connectedApp
	for _, a := range apps {
		if a.Name == "Test MCP client" {
			target = a
		}
	}
	require.NotZero(t, target.ID)
	require.Equal(t, http.StatusNoContent, revokeApp(t, o, alice, target.ID, true).StatusCode)

	assert.False(t, o.Accepts(mine.AccessToken))
	assert.True(t, o.Accepts(theirs.AccessToken), "the same client for another user is untouched")
	assert.True(t, o.Accepts(other.AccessToken), "another client of the same user is untouched")
	assert.Len(t, listApps(t, o, bob), 1)
}

func TestRevokingSomeoneElsesOrAnUnknownAppIsNotFound(t *testing.T) {
	o := testsupport.NewOAuth(t)
	theirs := o.Connect(bob, oauthserver.ScopeCatalogRead)
	bobsApp := listApps(t, o, bob)[0]

	// Alice has not connected it.
	assert.Equal(t, http.StatusNotFound, revokeApp(t, o, alice, bobsApp.ID, true).StatusCode)
	assert.True(t, o.Accepts(theirs.AccessToken), "a refused revoke ends nothing")
	assert.Equal(t, http.StatusNotFound, revokeApp(t, o, alice, 987654, true).StatusCode)

	// Revoking twice: the second finds nothing to revoke.
	require.Equal(t, http.StatusNoContent, revokeApp(t, o, bob, bobsApp.ID, true).StatusCode)
	assert.Equal(t, http.StatusNotFound, revokeApp(t, o, bob, bobsApp.ID, true).StatusCode)
}

func TestTheOperationsNeedASessionAndTheRevokeNeedsTheXSRFHeader(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)
	app := listApps(t, o, alice)[0]

	assert.Equal(t, http.StatusUnauthorized, o.Do(http.MethodGet, "/api/oauth/connected_apps", nil).StatusCode)
	assert.Equal(t, http.StatusUnauthorized, o.Do(http.MethodDelete, "/api/oauth/connected_apps/"+strconv.Itoa(app.ID), nil).StatusCode)
	assert.Equal(t, http.StatusUnauthorized, revokeApp(t, o, alice, app.ID, false).StatusCode)
	assert.True(t, o.Accepts(tokens.AccessToken))
}

func TestAnAppUnusedForThirtyDaysLapsesAndOneRefreshedStaysConnected(t *testing.T) {
	o := testsupport.NewOAuth(t)
	idle := o.Connect(alice, oauthserver.ScopeCatalogRead)
	// A second client that keeps refreshing.
	active, otherClientID := connectOtherClient(t, o, alice)
	require.Len(t, listApps(t, o, alice), 2)

	// Every few days the active app refreshes; the idle one does nothing, not
	// even a refresh (an MCP call with a live access token is not a use).
	refresh := active.RefreshToken
	for range 7 {
		o.Advance(5 * 24 * time.Hour)
		tokens, res := o.Refresh(otherClientID, refresh)
		require.Equal(t, http.StatusOK, res.StatusCode)
		refresh = tokens.RefreshToken
	}

	// 35 days in: the idle app's refresh token (30 days from its issue) is gone.
	apps := listApps(t, o, alice)
	require.Len(t, apps, 1)
	assert.Equal(t, "Other", apps[0].Name)
	assert.False(t, o.Accepts(idle.AccessToken))
	_, res := o.Refresh(testsupport.OAuthClientID, idle.RefreshToken)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "the app no longer works")
}

func TestAnAppListedJustBeforeTheThirtiethDayIsGoneJustAfter(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Connect(alice, oauthserver.ScopeCatalogRead)

	o.Advance(30*24*time.Hour - time.Minute)
	assert.Len(t, listApps(t, o, alice), 1)
	o.Advance(2 * time.Minute)
	assert.Empty(t, listApps(t, o, alice))
}

// connectOtherClient completes the flow for a second client, registered through
// the real registration endpoint and named "Other", so a test can tell apps of
// one user apart. It returns the tokens and the new client_id.
func connectOtherClient(t *testing.T, o *testsupport.OAuth, userID int) (testsupport.Tokens, string) {
	t.Helper()
	res := o.Register(map[string]any{
		"client_name":   "Other",
		"redirect_uris": []string{"https://other.example.com/callback"},
		"scope":         "catalog:read offline_access",
	})
	require.Equal(t, http.StatusCreated, res.StatusCode)
	got := decodeJSON[registered](t, res)
	f := o.StartWith(userID, testsupport.StartOptions{
		Scopes: []string{oauthserver.ScopeCatalogRead}, ClientID: got.ClientID, RedirectURIs: got.RedirectURIs,
	})
	tokens, exchange := f.Exchange(f.Approve(), testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, exchange.StatusCode)
	return tokens, got.ClientID
}

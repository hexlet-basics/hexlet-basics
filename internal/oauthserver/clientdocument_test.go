package oauthserver_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/oauthclient"
	"hexletbasics/internal/testsupport"
)

const documentRedirect = "https://app.example.com/callback"

// authorizeDocument sends the authorization request of a document client and
// returns the flow, whose Authorize response is the answer.
func authorizeDocument(o *testsupport.OAuth, clientID, redirect string, scopes ...string) *testsupport.Flow {
	return o.StartWith(alice(), testsupport.StartOptions{
		ClientID:     clientID,
		RedirectURIs: []string{redirect},
		Scopes:       scopes,
		Params:       nil,
	})
}

func storedClient(t *testing.T, o *testsupport.OAuth, clientID string) bool {
	t.Helper()
	ok, err := o.DB.OAuthClient.Query().Where(oauthclient.ClientID(clientID)).Exist(t.Context())
	require.NoError(t, err)

	return ok
}

func TestTheMetadataAdvertisesClientIDMetadataDocumentSupportAtBothPaths(t *testing.T) {
	o := testsupport.NewOAuth(t)
	for _, path := range []string{"/.well-known/oauth-authorization-server", "/.well-known/openid-configuration"} {
		doc := decodeJSON[map[string]any](t, o.Do(http.MethodGet, path, nil))
		assert.Equal(t, true, doc["client_id_metadata_document_supported"], path)
	}
}

func TestAClientWithAMetadataDocumentCompletesTheFlowAndRefreshes(t *testing.T) {
	o := testsupport.NewOAuth(t)
	doc := testsupport.ValidDocument(documentRedirect)
	doc["scope"] = "catalog:read offline_access"
	id := o.Documents.Publish("/oauth/client.json", doc)

	authFlow := authorizeDocument(o, id, documentRedirect, "catalog:read")
	require.Equal(t, http.StatusFound, authFlow.Authorize.StatusCode)
	require.NotEmpty(t, authFlow.RequestID)

	// The consent screen names the client from its document and shows where the
	// user will be sent back to.
	flow := o.Do(http.MethodGet, "/api/oauth/consents/"+authFlow.RequestID, nil, func(r *http.Request) {
		r.AddCookie(testsupport.SessionFor(t, o.DB, alice()).Cookie)
	})
	require.Equal(t, http.StatusOK, flow.StatusCode)
	consent := decodeJSON[map[string]any](t, flow)
	assert.Equal(t, "Document client", consent["clientName"])
	assert.Equal(t, "app.example.com", consent["redirectHost"])

	tokens, exchange := authFlow.Exchange(authFlow.Approve(), testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, exchange.StatusCode)
	verified, err := o.Server.VerifyAccessToken(t.Context(), tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, id, verified.ClientID)

	res := o.Do(http.MethodPost, "/oauth/token", url.Values{
		"grant_type": {"refresh_token"}, "client_id": {id}, "refresh_token": {tokens.RefreshToken},
	})
	assert.Equal(t, http.StatusOK, res.StatusCode, "the row of a document client behaves like any other")
}

func TestAnUnlistedRedirectOfADocumentClientIsRefused(t *testing.T) {
	o := testsupport.NewOAuth(t)
	id := o.Documents.Publish("/c.json", testsupport.ValidDocument(documentRedirect))

	f := authorizeDocument(o, id, "https://evil.example.com/callback", "catalog:read")

	assert.NotEqual(t, http.StatusFound, f.Authorize.StatusCode, "the code must not be sent to an unlisted redirect")
	assert.Empty(t, f.RequestID)
}

func TestADocumentIsRefusedUnlessItIsExactlyWhatARegistrationMayBe(t *testing.T) {
	cases := []struct {
		name string
		edit func(map[string]any)
	}{
		{"its client_id is another URL", func(m map[string]any) { m["client_id"] = "https://example.com/other.json" }},
		{"it has no client_id", func(m map[string]any) { m["client_id"] = nil }},
		{"it names a plain-http redirect", func(m map[string]any) { m["redirect_uris"] = []string{"http://app.example.com/cb"} }},
		{"it names a wildcard redirect", func(m map[string]any) { m["redirect_uris"] = []string{"https://*.example.com/cb"} }},
		{"it names a custom-scheme redirect", func(m map[string]any) { m["redirect_uris"] = []string{"myapp://cb"} }},
		{"it asks for the client_credentials grant", func(m map[string]any) { m["grant_types"] = []string{"client_credentials"} }},
		{"it asks for a secret method", func(m map[string]any) { m["token_endpoint_auth_method"] = "client_secret_basic" }},
		{"it asks for an unoffered scope", func(m map[string]any) { m["scope"] = "profile" }},
		{"it carries a secret", func(m map[string]any) { m["client_secret"] = "s3cret" }},
		{"it carries an audience", func(m map[string]any) { m["audience"] = "https://elsewhere.example" }},
		{"it carries a jwks_uri", func(m map[string]any) { m["jwks_uri"] = "https://app.example.com/jwks" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			doc := testsupport.ValidDocument(documentRedirect)
			tc.edit(doc)

			id := o.Documents.Publish("/c.json", doc)
			if doc["client_id"] == nil {
				delete(doc, "client_id")
			}

			f := authorizeDocument(o, id, documentRedirect, "catalog:read")

			assert.Equal(t, http.StatusBadRequest, f.Authorize.StatusCode)
			assert.False(t, storedClient(t, o, id), "a refused document stores nothing")
		})
	}
}

func TestADocumentThatIsNotAnAcceptableResponseIsRefused(t *testing.T) {
	cases := map[string]http.HandlerFunc{
		"not found": http.NotFound,
		"server error": func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "boom", http.StatusInternalServerError)
		},
		"a redirect": func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/elsewhere.json", http.StatusFound)
		},
		"html": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`{}`))
		},
		"not json": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not json`))
		},
		"an array": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
		},
		"two objects": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{} {}`))
		},
		"larger than the limit": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"client_name":"` + strings.Repeat("a", 17<<10) + `"}`))
		},
	}
	for name, handler := range cases {
		t.Run(name, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			o.Documents.Serve("/c.json", handler)
			id := o.Documents.URL("/c.json")

			f := authorizeDocument(o, id, documentRedirect, "catalog:read")

			assert.Equal(t, http.StatusBadRequest, f.Authorize.StatusCode)
			assert.False(t, storedClient(t, o, id))
		})
	}
}

func TestAClientIDURLThatIsNotAnOrdinaryPublicHTTPSResourceIsRefusedBeforeAnyFetch(t *testing.T) {
	for _, id := range []string{
		"https://user:pass@example.com/c.json",
		"https://example.com",
		"https://example.com/",
		"https://example.com/c.json#fragment",
		"https://example.com:8443/c.json",
		"https://example.com/a/../c.json",
		"https://example.com/a/./c.json",
		"https://example.com/" + strings.Repeat("a", 2100),
	} {
		t.Run(id[:min(len(id), 50)], func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			o.Documents.Serve("/c.json", func(_ http.ResponseWriter, _ *http.Request) {
				t.Error("an unacceptable URL must not be fetched")
			})

			f := authorizeDocument(o, id, documentRedirect, "catalog:read")

			assert.Equal(t, http.StatusBadRequest, f.Authorize.StatusCode)
			assert.False(t, storedClient(t, o, id))
		})
	}
}

// The production fetcher refuses a destination that is not on the public
// internet, whatever form its address takes. The check runs when the socket
// connects, so these literal addresses and a name that resolves to loopback are
// all stopped there; nothing is sent.
func TestTheProductionFetcherRefusesNonPublicTargets(t *testing.T) {
	for _, host := range []string{
		"127.0.0.1", "[::1]", "localhost", "0.0.0.0",
		"10.0.0.1", "172.16.0.1", "192.168.1.1", "100.64.0.1",
		"169.254.169.254", "[fe80::1]", "[fd00::1]", "[fc00::1]",
		"224.0.0.1", "[ff02::1]",
		"[::ffff:127.0.0.1]", "[::ffff:10.0.0.1]", "[::ffff:169.254.169.254]",
		"[64:ff9b::a00:1]", "[2002:a00:1::1]",
	} {
		t.Run(host, func(t *testing.T) {
			o := testsupport.NewOAuthWithProductionFetcher(t)
			id := "https://" + host + "/client.json"

			f := authorizeDocument(o, id, documentRedirect, "catalog:read")

			assert.Equal(t, http.StatusBadRequest, f.Authorize.StatusCode)
			assert.False(t, storedClient(t, o, id))
		})
	}
}

func cacheControl(v string) http.Header { return http.Header{"Cache-Control": {v}} }

func updatedName(t *testing.T, o *testsupport.OAuth, id string) string {
	t.Helper()
	row, err := o.DB.OAuthClient.Query().Where(oauthclient.ClientID(id)).Only(t.Context())
	require.NoError(t, err)

	return row.Name
}

func TestADocumentIsCachedAccordingToItsHeadersWithinBounds(t *testing.T) {
	cases := []struct {
		name    string
		headers http.Header
		// staleAfter is how long after the fetch the document is fetched again;
		// freshBefore is a moment at which it must still be the cached copy.
		freshBefore, staleAfter time.Duration
	}{
		{"max-age", cacheControl("max-age=7200"), time.Hour, 2*time.Hour + time.Minute},
		{"no header means an hour", nil, 30 * time.Minute, time.Hour + time.Minute},
		{"a tiny max-age is raised to the floor", cacheControl("max-age=1"), 2 * time.Minute, 6 * time.Minute},
		{"no-store is the floor too", cacheControl("no-store"), 2 * time.Minute, 6 * time.Minute},
		{"a huge max-age is cut to a day", cacheControl("max-age=31536000"), 23 * time.Hour, 25 * time.Hour},
		{"s-maxage wins", cacheControl("max-age=60, s-maxage=7200"), time.Hour, 2*time.Hour + time.Minute},
		{"Expires is used without Cache-Control", http.Header{"Expires": {time.Now().Add(3 * time.Hour).UTC().Format(http.TimeFormat)}},
			2 * time.Hour, 3*time.Hour + 5*time.Minute},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			publish := func(name string) string {
				doc := testsupport.ValidDocument(documentRedirect)
				doc["client_name"] = name

				return o.Documents.Publish("/c.json", doc, tc.headers)
			}
			id := publish("First")
			require.Equal(t, http.StatusFound, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode)
			publish("Second")

			o.Advance(tc.freshBefore)
			require.Equal(t, http.StatusFound, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode)
			assert.Equal(t, "First", updatedName(t, o, id), "still the cached copy")

			o.Advance(tc.staleAfter - tc.freshBefore)
			require.Equal(t, http.StatusFound, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode)
			assert.Equal(t, "Second", updatedName(t, o, id), "fetched again once the cache expired")
		})
	}
}

func TestARefreshThatFailsRefusesNewGrantsButKeepsTheOnesAlreadyHeld(t *testing.T) {
	o := testsupport.NewOAuth(t)
	doc := testsupport.ValidDocument(documentRedirect)
	doc["scope"] = "catalog:read offline_access"
	id := o.Documents.Publish("/c.json", doc)
	f := authorizeDocument(o, id, documentRedirect, "catalog:read")
	tokens, res := f.Exchange(f.Approve(), testsupport.PKCEVerifier())
	require.Equal(t, http.StatusOK, res.StatusCode)

	o.Documents.Serve("/c.json", http.NotFound)
	o.Advance(2 * time.Hour)

	assert.Equal(t, http.StatusBadRequest, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode,
		"a document that can no longer be verified does not start a new grant")
	assert.True(t, storedClient(t, o, id), "the client row is kept")
	refreshed := o.Do(http.MethodPost, "/oauth/token", url.Values{
		"grant_type": {"refresh_token"}, "client_id": {id}, "refresh_token": {tokens.RefreshToken},
	})
	assert.Equal(t, http.StatusOK, refreshed.StatusCode, "a grant already held does not depend on the remote host")
}

func TestAFailedFetchIsRememberedSoItIsNotRepeatedOnEveryRequest(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Documents.Serve("/c.json", http.NotFound)
	id := o.Documents.URL("/c.json")
	require.Equal(t, http.StatusBadRequest, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode)

	// The host fixes its document, but the failure is remembered for a while.
	o.Documents.Publish("/c.json", testsupport.ValidDocument(documentRedirect))
	assert.Equal(t, http.StatusBadRequest, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode)

	o.Advance(6 * time.Minute)
	assert.Equal(t, http.StatusFound, authorizeDocument(o, id, documentRedirect, "catalog:read").Authorize.StatusCode)
}

func TestFetchingMetadataDocumentsIsRateLimitedPerAddress(t *testing.T) {
	o := testsupport.NewOAuth(t)
	attacker := testsupport.FromAddress("203.0.113.70")

	start := func(path string, from func(*http.Request)) *http.Response {
		id := o.Documents.Publish(path, testsupport.ValidDocument(documentRedirect))
		query := url.Values{
			"client_id": {id}, "redirect_uri": {documentRedirect}, "response_type": {"code"}, "scope": {"catalog:read"},
			"state": {"s"}, "code_challenge": {testsupport.PKCEChallenge()}, "code_challenge_method": {"S256"},
		}

		return o.Do(http.MethodGet, "/oauth/authorize?"+query.Encode(), nil, from)
	}
	for i := range 10 {
		require.Equal(t, http.StatusFound, start(fmt.Sprintf("/c%d.json", i), attacker).StatusCode)
	}

	res := start("/c-more.json", attacker)

	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	assert.NotEmpty(t, res.Header.Get("Retry-After"))
	assert.False(t, storedClient(t, o, o.Documents.URL("/c-more.json")), "the refused request fetched nothing")
	assert.Equal(t, http.StatusFound, start("/c0.json", attacker).StatusCode, "a cached document costs nothing")
	assert.Equal(t, http.StatusFound, start("/c-more.json", testsupport.FromAddress("203.0.113.71")).StatusCode,
		"another address has its own budget")
}

func TestARegisteredClientIDThatIsNotAURLNeverStartsAFetch(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Documents.Serve("/c.json", func(http.ResponseWriter, *http.Request) { t.Error("nothing should be fetched") })

	f := authorizeDocument(o, "http://example.com/c.json", documentRedirect, "catalog:read")

	assert.NotEqual(t, http.StatusFound, f.Authorize.StatusCode)
}

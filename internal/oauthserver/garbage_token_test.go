package oauthserver_test

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// stranger is a string that never was a token of ours: random, base64url, the
// shape a client could send.
func stranger(t *testing.T) string {
	t.Helper()
	b := make([]byte, 32)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return base64.RawURLEncoding.EncodeToString(b)
}

// The library recovers a token's id by decrypting the string it is given, and its
// decryption falls back to AES in CFB mode, which does not authenticate: a
// stranger "decrypts" to arbitrary bytes, and in a few percent of cases those hold
// exactly one colon and are split into an id and a subject. The id is not valid
// UTF-8, so looking it up made PostgreSQL fail and the endpoint answer 500. A
// single request hits that rarely, so many are sent: against the unfixed code the
// chance that none of them does is negligible.
func TestPresentingAStrangerNeverFailsTheServer(t *testing.T) {
	o := testsupport.NewOAuth(t)
	// A client that has connected: revoking is a client's call, and an unknown
	// client is a different refusal.
	o.Connect(alice, oauthserver.ScopeCatalogRead)
	const strangers = 300

	for range strangers {
		token := stranger(t)

		// Revoking a token that is not there is not an error (RFC 7009).
		res := o.Revoke(testsupport.OAuthClientID, token)
		require.Equal(t, http.StatusOK, res.StatusCode, "revoking a stranger must answer 200")

		// Presenting one as a bearer is a refusal, never a server error.
		_, err := o.Server.VerifyAccessToken(t.Context(), token)
		require.ErrorIs(t, err, oauthserver.ErrInvalidToken, "a stranger is an invalid token")
	}
}

// A real token must still be revoked and verified: the shape check only turns
// away what could not be ours.
func TestTheShapeOfAnIdDoesNotTurnARealTokenAway(t *testing.T) {
	o := testsupport.NewOAuth(t)
	tokens := o.Connect(alice, oauthserver.ScopeCatalogRead)

	assert.True(t, o.Accepts(tokens.AccessToken))
	require.Equal(t, http.StatusOK, o.Revoke(testsupport.OAuthClientID, tokens.AccessToken).StatusCode)
	assert.False(t, o.Accepts(tokens.AccessToken), "revoking a real access token still ends it")
}

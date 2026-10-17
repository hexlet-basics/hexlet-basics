package jobs_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/oauthrefreshtoken"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

// The sweep deletes an expired, already used refresh token with its access
// token and keeps a live one. The worker is worked in the test's transaction
// (rivertest.Worker), so what it deletes is rolled back with the test.
func TestPurgeOAuthTokensJobDeletesWhatExpiredAndKeepsTheRest(t *testing.T) {
	db, _ := testsupport.NewClientWithTx(t)
	// The tokens are the act's input and carry the sweep's clock: they expire an
	// hour either side of the moment the worker reads, which no fixture can hold.
	//nolint:forbidigo // expiry is relative to the worker's clock (time.Now), so the rows cannot be fixed in fixtures/
	client := db.OAuthClient.Create().SetClientID("purge-client").SetName("Purge").
		SetRedirectUris([]string{"http://127.0.0.1:1/cb"}).SetScopes([]string{}).SaveX(t.Context())
	alice := fixtures.ID("user_admin_learner")
	past, future := time.Now().Add(-time.Hour), time.Now().Add(time.Hour)
	for id, expires := range map[string]time.Time{"purge-old": past, "purge-live": future} {
		//nolint:forbidigo // see above
		db.OAuthAccessToken.Create().SetID(id + "-access").SetClientID(client.ClientID).SetUserID(alice).
			SetScopes([]string{}).SetFamilyID(id).SetExpiresAt(expires).SaveX(t.Context())
		//nolint:forbidigo // see above
		db.OAuthRefreshToken.Create().SetID(id).SetTokenHash(id).SetClientID(client.ClientID).SetUserID(alice).
			SetScopes([]string{}).SetAuthTime(time.Now()).SetFamilyID(id).SetUsedAt(time.Now()).SetExpiresAt(expires).SaveX(t.Context())
	}

	result, err := work(t, jobs.PurgeOAuthTokensWorker(db, slog.New(slog.DiscardHandler)), jobs.PurgeOAuthTokensArgs{})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
	left, err := db.OAuthRefreshToken.Query().Where(oauthrefreshtoken.ClientID(client.ClientID)).
		Select(oauthrefreshtoken.FieldID).Strings(t.Context())
	require.NoError(t, err)
	assert.Equal(t, []string{"purge-live"}, left, "the expired token is gone, the live one stays")
	access, err := db.OAuthAccessToken.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, access)
}

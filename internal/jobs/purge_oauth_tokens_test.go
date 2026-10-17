package jobs_test

import (
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/oauthrefreshtoken"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// The sweep deletes an expired, already used refresh token with its access
// token and keeps a live one. The worker is worked in the test's transaction
// (rivertest.Worker), so what it deletes is rolled back with the test.
//
// The tokens are the ones the authorization server issues, signed in through
// the real flow: one person connects ninety days ago, before the sweep's clock
// (time.Now), and rotates their refresh token; another connects now and
// rotates theirs, so each has a used token and its child.
func TestPurgeOAuthTokensJobDeletesWhatExpiredAndKeepsTheRest(t *testing.T) {
	o := testsupport.NewOAuth(t)
	alice, bob := fixtures.ID("user_admin_learner"), fixtures.ID("user_staff_editor")
	connectAndRotate := func(userID int) {
		tokens := o.Connect(userID, oauthserver.ScopeCatalogRead)
		_, res := o.Refresh(testsupport.OAuthClientID, tokens.RefreshToken)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, res.Body.Close())
	}
	o.Advance(-90 * 24 * time.Hour)
	connectAndRotate(alice)
	o.Advance(90 * 24 * time.Hour)
	connectAndRotate(bob)

	before, err := o.DB.OAuthRefreshToken.Query().Where(oauthrefreshtoken.UserID(alice)).Count(t.Context())
	require.NoError(t, err)
	require.Equal(t, 2, before, "the expired person holds a used token and its child before the sweep")

	result, err := work(t, jobs.PurgeOAuthTokensWorker(o.DB, slog.New(slog.DiscardHandler)), jobs.PurgeOAuthTokensArgs{})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
	expired, err := o.DB.OAuthRefreshToken.Query().Where(oauthrefreshtoken.UserID(alice)).Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, expired, "the expired tokens are gone")
	live, err := o.DB.OAuthRefreshToken.Query().Where(oauthrefreshtoken.UserID(bob)).Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, live, "the live ones stay, the used one with its child")
	access, err := o.DB.OAuthAccessToken.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, access, "a rotation replaces the access token, so the live person has one and the expired one is gone")
}

package jobs_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"hexletbasics/ent/oauthrefreshtoken"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// TestWorkerPurgesExpiredOAuthTokensPeriodically proves the worker runs the
// cleanup on its own (RunOnStart): an expired, already used refresh token and
// its access token are deleted, a live one is kept. It uses a real pool, so it
// removes what it creates.
func TestWorkerPurgesExpiredOAuthTokensPeriodically(t *testing.T) {
	ctx := context.Background()

	sqlDB, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	db := store.NewClient(sqlDB)
	client := db.OAuthClient.Create().SetClientID("purge-client").SetName("Purge").
		SetRedirectUris([]string{"http://127.0.0.1:1/cb"}).SetScopes([]string{}).SaveX(ctx)
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_ = db.OAuthClient.DeleteOne(client).Exec(cleanupCtx) // cascades the tokens
		_, _ = sqlDB.ExecContext(cleanupCtx, "DELETE FROM river_job WHERE kind = 'purge_oauth_tokens'")
		_ = sqlDB.Close()
	})

	const alice = 1001
	past, future := time.Now().Add(-time.Hour), time.Now().Add(time.Hour)
	for id, expires := range map[string]time.Time{"purge-old": past, "purge-live": future} {
		db.OAuthAccessToken.Create().SetID(id + "-access").SetClientID("purge-client").SetUserID(alice).
			SetScopes([]string{}).SetFamilyID(id).SetExpiresAt(expires).SaveX(ctx)
		db.OAuthRefreshToken.Create().SetID(id).SetTokenHash(id).SetClientID("purge-client").SetUserID(alice).
			SetScopes([]string{}).SetAuthTime(time.Now()).SetFamilyID(id).SetUsedAt(time.Now()).SetExpiresAt(expires).SaveX(ctx)
	}

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	workerClient, err := jobs.NewWorkerClient(
		sqlDB,
		logger,
		jobs.NewErrorHandler(sentryClient),
		nooptrace.NewTracerProvider(),
		noopmetric.NewMeterProvider(),
		jobs.PurgeOAuthTokens(db, logger),
	)
	require.NoError(t, err)

	completed, cancel := workerClient.Subscribe(river.EventKindJobCompleted)
	defer cancel()

	require.NoError(t, workerClient.Start(ctx))
	t.Cleanup(func() {
		stopCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = workerClient.Stop(stopCtx)
	})

	timeout := time.After(15 * time.Second)
	for done := false; !done; {
		select {
		case ev := <-completed:
			done = ev.Job.Kind == jobs.PurgeOAuthTokensArgs{}.Kind()
		case <-timeout:
			t.Fatal("periodic OAuth purge did not complete within 15s")
		}
	}

	left, err := db.OAuthRefreshToken.Query().Where(oauthrefreshtoken.ClientID("purge-client")).Select(oauthrefreshtoken.FieldID).Strings(ctx)
	require.NoError(t, err)
	assert.Equal(t, []string{"purge-live"}, left, "the expired token is gone, the live one stays")
	access, err := db.OAuthAccessToken.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, access)
}

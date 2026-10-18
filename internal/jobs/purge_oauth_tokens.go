package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/riverqueue/river"

	"hexletbasics/ent"
	"hexletbasics/internal/oauthserver"
)

// purgeOAuthTokensInterval is hourly: an expired token is already refused by
// its expiry check, so the sweep only bounds table growth and its timing is
// irrelevant to security.
const purgeOAuthTokensInterval = time.Hour

// PurgeOAuthTokensArgs is the periodic cleanup of the OAuth server's tables
// (ADR-0022, ticket #839): rotated refresh tokens are kept to detect replays,
// so something has to delete them once they expire. It is enqueued only by the
// periodic scheduler.
type PurgeOAuthTokensArgs struct{}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (PurgeOAuthTokensArgs) Kind() string { return "purge_oauth_tokens" }

type purgeOAuthTokensWorker struct {
	river.WorkerDefaults[PurgeOAuthTokensArgs]

	db     *ent.Client
	logger *slog.Logger
	now    func() time.Time
}

func (w *purgeOAuthTokensWorker) Work(ctx context.Context, _ *river.Job[PurgeOAuthTokensArgs]) error {
	n, err := oauthserver.PurgeExpired(ctx, w.db, w.now())
	if err != nil {
		return errs.Wrapf(err, "purge oauth tokens")
	}

	if n > 0 {
		w.logger.InfoContext(ctx, "purged expired oauth rows", slog.Int("rows", n))
	}

	return nil
}

// PurgeOAuthTokens deletes expired OAuth tokens and abandoned authorization
// requests. It runs whether or not the OAuth server is enabled: on empty tables
// it deletes nothing.
func PurgeOAuthTokens(db *ent.Client, logger *slog.Logger) Module {
	return module(PurgeOAuthTokensWorker(db, logger), purgeOAuthTokensSchedule())
}

// PurgeOAuthTokensWorker is the worker alone, without its schedule, so that a
// test can work a sweep through it in a transaction (rivertest.Worker).
func PurgeOAuthTokensWorker(db *ent.Client, logger *slog.Logger) river.Worker[PurgeOAuthTokensArgs] {
	return &purgeOAuthTokensWorker{
		db:             db,
		logger:         logger,
		now:            time.Now,
		WorkerDefaults: river.WorkerDefaults[PurgeOAuthTokensArgs]{},
	}
}

// purgeOAuthTokensSchedule is the recurring insert, without unique options for
// the same reason as the other periodic jobs (empty args and River's default
// unique states would block every run after the first).
func purgeOAuthTokensSchedule() *river.PeriodicJob {
	return river.NewPeriodicJob(
		river.PeriodicInterval(purgeOAuthTokensInterval),
		func() (river.JobArgs, *river.InsertOpts) {
			return PurgeOAuthTokensArgs{}, nil
		},
		&river.PeriodicJobOpts{
			ID:         PurgeOAuthTokensArgs{}.Kind(),
			RunOnStart: true,
		},
	)
}

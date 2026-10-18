package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/riverqueue/river"

	"hexletbasics/internal/auth"
)

// reapExpiredSessionsInterval is how often sessions that ran out are removed.
// Nothing reads an expired session (the store treats it as absent), so this only
// bounds how long dead rows stay in the table.
const reapExpiredSessionsInterval = time.Hour

// ReapExpiredSessionsArgs is the periodic sweep of sessions whose time has
// passed. It replaces the cleanup goroutine session stores usually start for
// themselves, which a process owned by the errgroup cannot have (ADR-0010).
type ReapExpiredSessionsArgs struct{}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (ReapExpiredSessionsArgs) Kind() string { return "reap_expired_sessions" }

type reapExpiredSessionsWorker struct {
	river.WorkerDefaults[ReapExpiredSessionsArgs]

	sessions *auth.Authenticator
	logger   *slog.Logger
}

// Work removes the expired sessions and says how many it removed, so a sweep
// that suddenly removes far more than usual is visible in the worker logs.
func (w *reapExpiredSessionsWorker) Work(ctx context.Context, _ *river.Job[ReapExpiredSessionsArgs]) error {
	removed, err := w.sessions.DeleteExpired(ctx)
	if err != nil {
		return err
	}
	w.logger.InfoContext(ctx, "reaped expired sessions", slog.Int("removed", removed))
	return nil
}

// ReapExpiredSessions sweeps expired sessions, hourly.
func ReapExpiredSessions(manager *auth.Authenticator, logger *slog.Logger) Module {
	return module(ReapExpiredSessionsWorker(manager, logger), reapExpiredSessionsSchedule())
}

// ReapExpiredSessionsWorker is the worker alone, without its schedule, so that a
// test can work a sweep through it in a transaction (rivertest.Worker).
func ReapExpiredSessionsWorker(manager *auth.Authenticator, logger *slog.Logger) river.Worker[ReapExpiredSessionsArgs] {
	return &reapExpiredSessionsWorker{sessions: manager, logger: logger}
}

// reapExpiredSessionsSchedule is the sweep's recurring insert. Like the other
// periodic jobs it takes no unique options: with empty args River's default
// unique states include `completed`, which would block every run after the first.
func reapExpiredSessionsSchedule() *river.PeriodicJob {
	return river.NewPeriodicJob(
		river.PeriodicInterval(reapExpiredSessionsInterval),
		func() (river.JobArgs, *river.InsertOpts) {
			return ReapExpiredSessionsArgs{}, nil
		},
		&river.PeriodicJobOpts{
			ID:         ReapExpiredSessionsArgs{}.Kind(),
			RunOnStart: true,
		},
	)
}

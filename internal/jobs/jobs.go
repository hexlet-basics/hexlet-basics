// Package jobs wires the river background-job queue (ADR-0004): jobs are
// "do this" commands (exercise-version builds, AI lesson review, mail/SMS),
// distinct from watermill domain events.
//
// This is the queue BACKBONE only. The real domain jobs (AI lesson review for
// the course/lesson review actions, exercise-version builds for createVersion)
// land in later waves — each needs work on top of the queue (the lesson
// version-info graph, the AI integration, the exercise builder). The single
// worker here (`ping`) exists to prove the enqueue -> work -> complete pipeline
// end to end; it is not a domain job.
package jobs

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"

	"hexletbasics/internal/courseloader"
)

// defaultMaxWorkers bounds concurrency on the default queue. Tune per-queue when
// real jobs (some heavy, e.g. exercise builds) arrive.
const defaultMaxWorkers = 10

// PingArgs is a no-op job used to smoke-test the queue backbone end to end. It is
// NOT a domain job — see the package doc.
type PingArgs struct{}

// Kind is river's stable job discriminator.
func (PingArgs) Kind() string { return "ping" }

type pingWorker struct {
	river.WorkerDefaults[PingArgs]
}

func (*pingWorker) Work(_ context.Context, _ *river.Job[PingArgs]) error { return nil }

// Workers builds the worker registry. The exercise loader is registered when a
// loader is supplied; a nil loader (insert-only clients that never Start) skips
// it, since only the worker process needs the loader's db/blob dependencies.
func Workers(loader *courseloader.Loader, leadSyncer LeadSyncer) *river.Workers {
	w := river.NewWorkers()
	river.AddWorker(w, &pingWorker{})
	if loader != nil {
		river.AddWorker(w, &exerciseLoaderWorker{loader: loader})
	}
	if leadSyncer != nil {
		river.AddWorker(w, &amoCRMLeadWorker{syncer: leadSyncer})
	}
	return w
}

// NewInsertOnlyClient builds the River adapter used by synchronous callers.
// It deliberately has no queues or workers and must never be started. Sharing
// ent's database/sql pool lets InsertTx participate in a business transaction
// while job execution remains isolated in the worker process.
func NewInsertOnlyClient(
	db *sql.DB,
	logger *slog.Logger,
) (*river.Client[*sql.Tx], error) {
	return river.NewClient(riverdatabasesql.New(db), &river.Config{
		Logger: logger,
	})
}

// NewWorkerClient builds the River runtime used only by the worker process.
// The caller owns Start and Stop. Worker dependencies stay behind this
// constructor so an HTTP process cannot accidentally execute background jobs.
func NewWorkerClient(
	db *sql.DB,
	loader *courseloader.Loader,
	leadSyncer LeadSyncer,
	logger *slog.Logger,
	errorHandler *ErrorHandler,
) (*river.Client[*sql.Tx], error) {
	return river.NewClient(riverdatabasesql.New(db), &river.Config{
		ErrorHandler: errorHandler,
		Logger:       logger,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: defaultMaxWorkers},
		},
		Workers: Workers(loader, leadSyncer),
	})
}

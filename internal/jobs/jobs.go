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
func Workers(loader *courseloader.Loader) *river.Workers {
	w := river.NewWorkers()
	river.AddWorker(w, &pingWorker{})
	if loader != nil {
		river.AddWorker(w, &exerciseLoaderWorker{loader: loader})
	}
	return w
}

// NewClient builds the River client over the same database/sql pool ent uses.
// It is insert-and-work capable: the caller Start()s it to run workers and
// Stop()s it on shutdown. Sharing the pool also lets InsertTx participate in an
// ent business transaction. The loader backs the exercise-build worker; pass
// nil for an insert-only client.
func NewClient(
	db *sql.DB,
	loader *courseloader.Loader,
	logger *slog.Logger,
	errorHandler *ErrorHandler,
) (*river.Client[*sql.Tx], error) {
	return river.NewClient(riverdatabasesql.New(db), &river.Config{
		ErrorHandler: errorHandler,
		Logger:       logger,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: defaultMaxWorkers},
		},
		Workers: Workers(loader),
	})
}

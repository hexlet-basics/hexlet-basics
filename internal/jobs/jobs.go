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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
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

// Workers builds the worker registry. Real job workers register here as they are
// implemented.
func Workers() *river.Workers {
	w := river.NewWorkers()
	river.AddWorker(w, &pingWorker{})
	return w
}

// NewClient builds the river client over a pgx pool (riverpgxv5). It is
// insert-and-work capable: the caller Start()s it to run workers and Stop()s it
// on shutdown. Insert-only callers (e.g. HTTP handlers enqueuing work) can use
// the same client without starting it.
func NewClient(pool *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: defaultMaxWorkers},
		},
		Workers: Workers(),
	})
}

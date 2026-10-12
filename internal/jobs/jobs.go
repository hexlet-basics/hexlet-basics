// Package jobs wires the river background-job queue (ADR-0004): jobs are
// "do this" commands (exercise-version builds, AI lesson review, mail, amoCRM
// sync), distinct from watermill domain events.
//
// Each job kind lives in its own file: its Args (with the insert policy —
// queue, attempts, uniqueness), its worker, and a Module constructor that the
// worker process passes to NewWorkerClient.
package jobs

import (
	"database/sql"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertype"
	"github.com/riverqueue/rivercontrib/otelriver"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Queues isolate classes of work so a burst in one cannot starve another: a
// course-wide batch of LLM reviews must not delay the sign-in email a visitor
// is waiting for. Each job kind names its queue in its InsertOpts; the worker
// process runs every queue below.
const (
	// QueueMail carries transactional email the user is waiting on.
	QueueMail = "mail"
	// QueueBuilds carries course-version builds and their maintenance.
	QueueBuilds = "builds"
	// QueueLLM carries paid, rate-limited LLM calls.
	QueueLLM = "llm"
)

// queues is the worker's concurrency budget per queue. Builds are heavy (git,
// parsing, uploads) so few run at once; LLM concurrency stays under provider
// rate limits; River's default queue keeps the remaining light jobs.
func queues() map[string]river.QueueConfig {
	return map[string]river.QueueConfig{
		river.QueueDefault: {MaxWorkers: 10},
		QueueMail:          {MaxWorkers: 5},
		QueueBuilds:        {MaxWorkers: 2},
		QueueLLM:           {MaxWorkers: 3},
	}
}

// Module is one job kind's registration in the worker process: its worker and,
// for recurring work, its schedule. A module exists only when its dependencies
// do, so the process composition (di) decides what runs, for example no LLM
// modules without credentials, instead of nil checks inside this package.
type Module struct {
	register func(*river.Workers)
	periodic []*river.PeriodicJob
}

// module registers w and, optionally, the schedule that inserts its jobs.
// Periodic jobs travel with their worker: a periodic kind without a registered
// worker would be inserted and then fail every tick.
func module[A river.JobArgs](w river.Worker[A], periodic ...*river.PeriodicJob) Module {
	return Module{
		register: func(ws *river.Workers) { river.AddWorker(ws, w) },
		periodic: periodic,
	}
}

// NewInsertOnlyClient builds the River adapter used by synchronous callers.
// It deliberately has no queues or workers and must never be started. Sharing
// ent's database/sql pool lets InsertTx participate in a business transaction
// while job execution remains isolated in the worker process.
func NewInsertOnlyClient(
	db *sql.DB,
	logger *slog.Logger,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*river.Client[*sql.Tx], error) {
	return river.NewClient(riverdatabasesql.New(db), &river.Config{
		Logger:  logger,
		Plugins: openTelemetryPlugins(tracerProvider, meterProvider),
	})
}

// NewWorkerClient builds the River runtime used only by the worker process.
// The caller owns Start and Stop. Worker dependencies stay behind the modules
// so an HTTP process cannot accidentally execute background jobs. The client
// fetches only the kinds its modules register, so a job whose module is absent
// (for example an LLM job on a worker without credentials) waits in the queue
// for a worker that has it instead of failing with UnknownJobKindError.
func NewWorkerClient(
	db *sql.DB,
	logger *slog.Logger,
	errorHandler *ErrorHandler,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
	modules ...Module,
) (*river.Client[*sql.Tx], error) {
	workers := river.NewWorkers()
	var periodic []*river.PeriodicJob
	for _, m := range modules {
		m.register(workers)
		periodic = append(periodic, m.periodic...)
	}
	return river.NewClient(riverdatabasesql.New(db), &river.Config{
		ErrorHandler:        errorHandler,
		FetchOnlyKnownKinds: true,
		Logger:              logger,
		PeriodicJobs:        periodic,
		Plugins:             openTelemetryPlugins(tracerProvider, meterProvider),
		Queues:              queues(),
		Workers:             workers,
	})
}

// openTelemetryPlugins keeps River's observability policy consistent between
// producers and workers. Explicit providers bind instruments to the process SDK
// at construction time instead of relying on mutable OpenTelemetry globals.
func openTelemetryPlugins(
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) []rivertype.Plugin {
	return []rivertype.Plugin{
		otelriver.NewMiddleware(&otelriver.MiddlewareConfig{
			EnableTracePropagation: true,
			MeterProvider:          meterProvider,
			TracerProvider:         tracerProvider,
		}),
	}
}

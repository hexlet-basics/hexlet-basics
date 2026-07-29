// Package di wires the application's dependency-injection container.
package di

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/events"
)

// ServerDependencies is the process-level root of the synchronous application
// graph. cmd/server owns the lifecycle of every resource exposed here.
type ServerDependencies struct {
	Logger           *slog.Logger   `do:""`
	SentryClient     *sentry.Client `do:""`
	HTTPServer       *http.Server   `do:""`
	Database         *ent.Client    `do:""`
	Bucket           *blob.Bucket   `do:""`
	OpenTelemetrySDK *otelconf.SDK  `do:""`
}

// WorkerDependencies is the process-level root of the asynchronous application
// graph. cmd/worker owns the lifecycle of every resource exposed here.
type WorkerDependencies struct {
	Logger           *slog.Logger           `do:""`
	SentryClient     *sentry.Client         `do:""`
	Database         *ent.Client            `do:""`
	RiverClient      *river.Client[*sql.Tx] `do:""`
	EventRuntime     *events.Runtime        `do:""`
	Bucket           *blob.Bucket           `do:""`
	OpenTelemetrySDK *otelconf.SDK          `do:""`
}

// BuildServer resolves the complete synchronous application graph. Its River
// client can enqueue jobs but has no queues or workers, so this process cannot
// execute background work.
func BuildServer() (ServerDependencies, error) {
	return do.InvokeStruct[ServerDependencies](newServerContainer())
}

// BuildWorker resolves the complete asynchronous application graph. It owns
// the Watermill subscribers and River workers but exposes no HTTP server.
func BuildWorker() (WorkerDependencies, error) {
	return do.InvokeStruct[WorkerDependencies](newWorkerContainer())
}

func newServerContainer() *do.RootScope {
	return do.New(commonPackage, serverPackage)
}

func newWorkerContainer() *do.RootScope {
	return do.New(commonPackage, workerPackage)
}

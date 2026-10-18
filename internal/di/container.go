// Package di wires the application's dependency-injection container.
package di

import (
	"context"
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

	injector do.Injector
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

	injector do.Injector
}

// Shutdown releases every resource the process built, through the injector
// (ADR-0023). Call it once, after the HTTP server has drained.
func (d ServerDependencies) Shutdown(ctx context.Context) error {
	return shutdown(ctx, d.injector)
}

// Shutdown releases every resource the process built, through the injector
// (ADR-0023). Call it once, after Watermill and River have stopped.
func (d WorkerDependencies) Shutdown(ctx context.Context) error {
	return shutdown(ctx, d.injector)
}

// BuildServer resolves the complete synchronous application graph. Its River
// client can enqueue jobs but has no queues or workers, so this process cannot
// execute background work. When the build fails, what was already built is
// released before the error returns.
func BuildServer() (ServerDependencies, error) {
	injector := newServerContainer()
	deps, err := do.InvokeStruct[ServerDependencies](injector)
	if err == nil {
		err = registerClosers(injector)
	}
	if err != nil {
		_ = injector.Shutdown()
		return ServerDependencies{}, err
	}
	deps.injector = injector
	return deps, nil
}

// BuildWorker resolves the complete asynchronous application graph. It owns
// the Watermill subscribers and River workers but exposes no HTTP server.
// When the build fails, what was already built is released before the error
// returns.
func BuildWorker() (WorkerDependencies, error) {
	injector := newWorkerContainer()
	deps, err := do.InvokeStruct[WorkerDependencies](injector)
	if err == nil {
		err = registerClosers(injector)
	}
	if err != nil {
		_ = injector.Shutdown()
		return WorkerDependencies{}, err
	}
	deps.injector = injector
	return deps, nil
}

func newServerContainer() *do.RootScope {
	return do.New(commonPackage, closersPackage, serverPackage)
}

func newWorkerContainer() *do.RootScope {
	return do.New(commonPackage, closersPackage, workerPackage)
}

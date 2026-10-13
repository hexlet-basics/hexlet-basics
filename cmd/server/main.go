package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/contrib/otelconf"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/di"
)

// shutdownTimeout bounds graceful shutdown: in-flight requests get this long to
// finish before the process exits (k8s sends SIGKILL after its own grace period).
const shutdownTimeout = 15 * time.Second

// main owns the process exit code and nothing else: all the work is in run,
// which returns its failure instead of exiting, so a test can drive it. Exiting
// is the last thing main does, after stop has released the signal handlers, so
// no deferred cleanup is skipped.
func main() {
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	err := run(signalCtx, stop)
	stop()
	if err != nil {
		os.Exit(1)
	}
}

// run builds the dependency graph, runs the process until signalCtx is done or
// the runtime fails, then releases the resources. It reports a failure to build
// or to run; stopping on a signal is not one.
func run(signalCtx context.Context, stopSignals func()) error {
	dependencies, err := di.BuildServer()
	if err != nil {
		slog.Error("building server dependencies", "err", err)
		return fmt.Errorf("building server dependencies: %w", err)
	}

	app := application{
		http:     dependencies.HTTPServer,
		httpAddr: dependencies.HTTPServer.Addr,
		logger:   dependencies.Logger,
	}
	runtimeErr := app.run(signalCtx, stopSignals)

	// The signal has already cancelled signalCtx; cleanup keeps its values but not
	// that cancellation, and gets its own deadline.
	cleanupCtx, cancelCleanup := context.WithTimeout(context.WithoutCancel(signalCtx), shutdownTimeout)
	defer cancelCleanup()
	closeResources(
		cleanupCtx,
		dependencies.Logger,
		dependencies.Bucket,
		dependencies.Database,
		dependencies.OpenTelemetrySDK,
		dependencies.SentryClient,
	)

	if runtimeErr != nil && !errors.Is(runtimeErr, errRuntimeSignal) {
		dependencies.Logger.Error("runtime failed", "err", runtimeErr)
		return runtimeErr
	}
	return nil
}

func closeResources(
	ctx context.Context,
	logger *slog.Logger,
	bucket *blob.Bucket,
	db *ent.Client,
	otelSDK *otelconf.SDK,
	sentryClient *sentry.Client,
) {
	if err := bucket.Close(); err != nil {
		logger.Error("closing blob bucket", "err", err)
	}
	if err := db.Close(); err != nil {
		logger.Error("closing database", "err", err)
	}
	if err := otelSDK.Shutdown(ctx); err != nil {
		logger.Error("shutting down OpenTelemetry", "err", err)
	}
	if !sentryClient.FlushWithContext(ctx) {
		logger.Error("flushing Sentry events timed out")
	}
}

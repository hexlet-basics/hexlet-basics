package main

import (
	"context"
	"errors"
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

func main() {
	dependencies, err := di.BuildServer()
	if err != nil {
		slog.Error("building server dependencies", "err", err)
		os.Exit(1)
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	app := application{
		http:     dependencies.HTTPServer,
		httpAddr: dependencies.HTTPServer.Addr,
		logger:   dependencies.Logger,
	}
	runtimeErr := app.run(signalCtx, stop)

	cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), shutdownTimeout)
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
		os.Exit(1)
	}
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

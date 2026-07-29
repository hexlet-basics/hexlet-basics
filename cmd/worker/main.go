package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/di"
	"hexletbasics/internal/events"
)

// shutdownTimeout bounds process cleanup after River and Watermill have been
// asked to drain. Deployment termination grace must exceed this value.
const shutdownTimeout = 15 * time.Second

func main() {
	injector := di.NewWorker()

	logger := do.MustInvoke[*slog.Logger](injector)
	sentryClient := do.MustInvoke[*sentry.Client](injector)
	db := do.MustInvoke[*ent.Client](injector)
	riverClient := do.MustInvoke[*river.Client[*sql.Tx]](injector)
	eventRuntime := do.MustInvoke[*events.Runtime](injector)
	bucket := do.MustInvoke[*blob.Bucket](injector)
	otelSDK := do.MustInvoke[*otelconf.SDK](injector)

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	app := application{
		jobs:   riverClient,
		events: eventRuntime,
		logger: logger,
	}
	runtimeErr := app.run(signalCtx, stop)

	cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelCleanup()
	closeResources(cleanupCtx, logger, bucket, db, otelSDK, sentryClient)

	if runtimeErr != nil && !errors.Is(runtimeErr, errRuntimeSignal) {
		logger.Error("worker runtime failed", "err", runtimeErr)
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

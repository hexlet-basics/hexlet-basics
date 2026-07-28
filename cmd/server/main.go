package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/di"
)

// shutdownTimeout bounds graceful shutdown: in-flight requests get this long to
// finish before the process exits (k8s sends SIGKILL after its own grace period).
const shutdownTimeout = 15 * time.Second

func main() {
	injector := di.New()

	logger := do.MustInvoke[*slog.Logger](injector)
	srv := do.MustInvoke[*http.Server](injector)
	// Held before shutdown: ent.Client exposes Close, not a do Shutdowner, so it
	// is closed explicitly once the injector has drained the HTTP server.
	db := do.MustInvoke[*ent.Client](injector)
	// The river job queue is likewise stopped explicitly.
	riverClient := do.MustInvoke[*river.Client[*sql.Tx]](injector)
	// The blob bucket exposes Close, not a do Shutdowner, so drain it explicitly
	// once the HTTP server (its only user) has stopped serving.
	bucket := do.MustInvoke[*blob.Bucket](injector)

	// Start processing background jobs. Start returns once the client is running;
	// workers run until Stop drains them during shutdown.
	if err := riverClient.Start(context.Background()); err != nil {
		logger.Error("starting job queue", "err", err)
		os.Exit(1)
	}

	go func() {
		logger.Info("backend listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	<-ctx.Done()
	stop() // restore default signal handling so a second signal force-quits

	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// do stops services in reverse dependency order, calling http.Server.Shutdown
	// (graceful drain) with this context.
	if report := injector.ShutdownWithContext(shutdownCtx); !report.Succeed {
		logger.Error("shutdown reported errors", "err", report.Error())
	}
	// Drain in-flight jobs before closing ent, which owns their shared SQL pool.
	if err := riverClient.Stop(shutdownCtx); err != nil {
		logger.Error("stopping job queue", "err", err)
	}
	if err := bucket.Close(); err != nil {
		logger.Error("closing blob bucket", "err", err)
	}
	if err := db.Close(); err != nil {
		logger.Error("closing database", "err", err)
	}
}

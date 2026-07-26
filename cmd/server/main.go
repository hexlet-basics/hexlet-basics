package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samber/do/v2"

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
	if err := db.Close(); err != nil {
		logger.Error("closing database", "err", err)
	}
}

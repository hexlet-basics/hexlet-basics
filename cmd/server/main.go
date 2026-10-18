package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hexletbasics/internal/config"
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
	// Development convenience only: LoadDotenv reads .env unless APP_ENV names
	// another environment, and it runs here, not in run, so tests never read it.
	if err := config.LoadDotenv(); err != nil {
		slog.Error("loading .env", "err", err)
		os.Exit(1)
	}
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	err := run(signalCtx, stop)
	stop()
	if err != nil {
		// The one place a run failure is logged; run returns it unlogged.
		slog.Error("process failed", "err", err)
		os.Exit(1)
	}
}

// run builds the dependency graph, runs the process until signalCtx is done or
// the runtime fails, then releases the resources. It reports a failure to build
// or to run; stopping on a signal is not one.
func run(signalCtx context.Context, stopSignals func()) error {
	dependencies, err := di.BuildServer()
	if err != nil {
		return errs.Wrapf(err, "building server dependencies")
	}
	// main logs the error run returns; give it the configured handler.
	slog.SetDefault(dependencies.Logger)

	app := application{
		http:     dependencies.HTTPServer,
		httpAddr: dependencies.HTTPServer.Addr,
		logger:   dependencies.Logger,
		drain:    dependencies.Probes.Drain,
	}
	runtimeErr := app.run(signalCtx, stopSignals)

	// The signal has already cancelled signalCtx; cleanup keeps its values but not
	// that cancellation, and gets its own deadline.
	cleanupCtx, cancelCleanup := context.WithTimeout(context.WithoutCancel(signalCtx), shutdownTimeout)
	defer cancelCleanup()
	// Every actor has stopped by now, so the injector may release the resources
	// they used (ADR-0023).
	if err := dependencies.Shutdown(cleanupCtx); err != nil {
		dependencies.Logger.ErrorContext(cleanupCtx, "releasing resources", "err", err)
	}

	return failure(runtimeErr)
}

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	errRuntimeSignal     = errors.New("runtime interrupted by signal")
	errHTTPServerStopped = errors.New("http server stopped unexpectedly")
)

// httpRuntime stays an interface: the runtime tests substitute a hand-written
// fake (stuckServer) for the *http.Server.
type httpRuntime interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Close() error
}

// application owns only the synchronous process runtime. Durable publishers
// and enqueuers are ordinary request dependencies; all async consumers execute
// in cmd/worker and therefore cannot affect HTTP availability.
type application struct {
	http     httpRuntime
	httpAddr string
	logger   *slog.Logger
	// drain, when set, is called first on every stop, before the HTTP drain:
	// /readyz starts answering 503 so the pod leaves rotation (ADR-0023).
	drain func()
	// shutdownTimeout overrides the package default; zero means the default.
	shutdownTimeout time.Duration
}

func (a *application) run(signalCtx context.Context, stopSignals func()) error {
	// The process context keeps the signal context's values but not its
	// cancellation: shutdown has to run after the signal has fired.
	processCtx, cancelProcess := context.WithCancel(context.WithoutCancel(signalCtx))
	defer cancelProcess()

	group, groupCtx := errgroup.WithContext(processCtx)

	group.Go(func() error {
		select {
		case <-signalCtx.Done():
			stopSignals()
			return errRuntimeSignal
		case <-groupCtx.Done():
			return nil
		}
	})

	group.Go(func() error {
		if a.logger != nil {
			a.logger.InfoContext(groupCtx, "backend listening", "addr", a.httpAddr)
		}
		if err := a.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errs.Wrapf(err, "serve HTTP")
		}
		return errHTTPServerStopped
	})

	// The group keeps only its first error, and on a signal that is the
	// sentinel, so the shutdown's own failure is kept apart and joined in.
	var shutdownErr error
	group.Go(func() error {
		<-groupCtx.Done()
		shutdownErr = a.shutdown(processCtx, cancelProcess)
		return shutdownErr
	})

	waitErr := group.Wait()
	if shutdownErr == nil || errors.Is(waitErr, shutdownErr) {
		return waitErr
	}
	return errors.Join(waitErr, shutdownErr)
}

// failure is what the process reports from run's result: the shutdown failures
// and any runtime failure, without the sentinel that only says a signal asked
// the process to stop (a normal exit).
func failure(err error) error {
	if joined, ok := errors.AsType[interface {
		error
		Unwrap() []error
	}](err); ok {
		var rest []error
		for _, e := range joined.Unwrap() {
			if !errors.Is(e, errRuntimeSignal) {
				rest = append(rest, e)
			}
		}
		return errors.Join(rest...)
	}
	if errors.Is(err, errRuntimeSignal) {
		return nil
	}
	return err
}

// shutdown returns what failed while stopping and does not log it: the process
// entry point logs the one error run returns (ADR-0018, log or return).
func (a *application) shutdown(parent context.Context, cancelProcess context.CancelFunc) error {
	if a.drain != nil {
		a.drain()
	}
	timeout := a.shutdownTimeout
	if timeout == 0 {
		timeout = shutdownTimeout
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), timeout)
	defer cancel()

	if a.logger != nil {
		a.logger.InfoContext(ctx, "shutting down")
	}
	if err := a.http.Shutdown(ctx); err != nil {
		stopErr := errs.Wrapf(err, "stop HTTP server")
		if closeErr := a.http.Close(); closeErr != nil {
			forceCloseErr := errs.Wrapf(closeErr, "force-close HTTP server")
			cancelProcess()
			return errors.Join(stopErr, forceCloseErr)
		}
		cancelProcess()
		return stopErr
	}

	cancelProcess()
	return nil
}

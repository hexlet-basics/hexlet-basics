package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	errRuntimeSignal     = errors.New("runtime interrupted by signal")
	errHTTPServerStopped = errors.New("http server stopped unexpectedly")
)

type httpRuntime interface {
	ListenAndServe() error
	Shutdown(context.Context) error
	Close() error
}

// application owns only the synchronous process runtime. Durable publishers
// and enqueuers are ordinary request dependencies; all async consumers execute
// in cmd/worker and therefore cannot affect HTTP availability.
type application struct {
	http     httpRuntime
	httpAddr string
	logger   *slog.Logger
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
			a.logger.Info("backend listening", "addr", a.httpAddr)
		}
		if err := a.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return errHTTPServerStopped
	})

	group.Go(func() error {
		<-groupCtx.Done()
		return a.shutdown(processCtx, cancelProcess)
	})

	return group.Wait()
}

func (a *application) shutdown(parent context.Context, cancelProcess context.CancelFunc) error {
	timeout := a.shutdownTimeout
	if timeout == 0 {
		timeout = shutdownTimeout
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), timeout)
	defer cancel()

	if a.logger != nil {
		a.logger.Info("shutting down")
	}
	if err := a.http.Shutdown(ctx); err != nil {
		stopErr := fmt.Errorf("stop HTTP server: %w", err)
		if a.logger != nil {
			a.logger.Error("runtime shutdown failed", "err", stopErr)
		}
		if closeErr := a.http.Close(); closeErr != nil {
			forceCloseErr := fmt.Errorf("force-close HTTP server: %w", closeErr)
			if a.logger != nil {
				a.logger.Error("runtime shutdown failed", "err", forceCloseErr)
			}
			cancelProcess()
			return errors.Join(stopErr, forceCloseErr)
		}
		cancelProcess()
		return stopErr
	}

	cancelProcess()
	return nil
}

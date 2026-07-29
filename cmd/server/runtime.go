package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	errRuntimeSignal      = errors.New("runtime interrupted by signal")
	errHTTPServerStopped  = errors.New("http server stopped unexpectedly")
	errJobQueueStopped    = errors.New("job queue stopped unexpectedly")
	errEventRouterStopped = errors.New("domain event router stopped unexpectedly")
)

type httpRuntime interface {
	ListenAndServe() error
	Shutdown(context.Context) error
	Close() error
}

type jobRuntime interface {
	Start(context.Context) error
	Stop(context.Context) error
	Stopped() <-chan struct{}
}

type eventRuntime interface {
	Run(context.Context) error
	Running() <-chan struct{}
	Close() error
}

// application owns the process runtime lifecycle. The DI container constructs
// its dependencies, while this module alone controls readiness, supervision,
// and the ordered shutdown of long-lived work.
type application struct {
	http     httpRuntime
	httpAddr string
	jobs     jobRuntime
	events   eventRuntime
	logger   *slog.Logger
}

func (a *application) run(signalCtx context.Context, stopSignals func()) error {
	processCtx, cancelProcess := context.WithCancel(context.Background())
	defer cancelProcess()

	group, groupCtx := errgroup.WithContext(processCtx)
	riverReady := make(chan struct{})
	state := &runtimeState{}

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
		if err := a.jobs.Start(processCtx); err != nil {
			return fmt.Errorf("start job queue: %w", err)
		}
		state.markJobsStarted()
		close(riverReady)

		<-a.jobs.Stopped()
		return errJobQueueStopped
	})

	group.Go(func() error {
		select {
		case <-riverReady:
		case <-groupCtx.Done():
			return nil
		}

		state.markEventsStarted()
		if err := a.events.Run(processCtx); err != nil {
			return fmt.Errorf("run domain event router: %w", err)
		}
		return errEventRouterStopped
	})

	group.Go(func() error {
		select {
		case <-a.events.Running():
		case <-groupCtx.Done():
			return nil
		}

		state.markHTTPStarted()
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
		return a.shutdown(state, cancelProcess)
	})

	return group.Wait()
}

func (a *application) shutdown(state *runtimeState, cancelProcess context.CancelFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if a.logger != nil {
		a.logger.Info("shutting down")
	}

	var shutdownErrors []error
	if state.httpWasStarted() {
		if err := a.http.Shutdown(ctx); err != nil {
			stopErr := fmt.Errorf("stop HTTP server: %w", err)
			shutdownErrors = append(shutdownErrors, stopErr)
			a.logShutdownError(stopErr)
			if closeErr := a.http.Close(); closeErr != nil {
				forceCloseErr := fmt.Errorf("force-close HTTP server: %w", closeErr)
				shutdownErrors = append(shutdownErrors, forceCloseErr)
				a.logShutdownError(forceCloseErr)
			}
		}
	}
	if state.eventsWereStarted() {
		if err := a.events.Close(); err != nil {
			stopErr := fmt.Errorf("stop domain event router: %w", err)
			shutdownErrors = append(shutdownErrors, stopErr)
			a.logShutdownError(stopErr)
		}
	}
	if state.jobsWereStarted() {
		if err := a.jobs.Stop(ctx); err != nil {
			stopErr := fmt.Errorf("stop job queue: %w", err)
			shutdownErrors = append(shutdownErrors, stopErr)
			a.logShutdownError(stopErr)
			cancelProcess()
		}
	}

	// Graceful stops have completed. Cancel the root context so any remaining
	// vendor goroutines observe process termination before resources are closed.
	cancelProcess()
	return errors.Join(shutdownErrors...)
}

func (a *application) logShutdownError(err error) {
	if a.logger != nil {
		a.logger.Error("runtime shutdown failed", "err", err)
	}
}

type runtimeState struct {
	mu sync.RWMutex

	httpStarted   bool
	jobsStarted   bool
	eventsStarted bool
}

func (s *runtimeState) markHTTPStarted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.httpStarted = true
}

func (s *runtimeState) markJobsStarted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobsStarted = true
}

func (s *runtimeState) markEventsStarted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventsStarted = true
}

func (s *runtimeState) httpWasStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.httpStarted
}

func (s *runtimeState) jobsWereStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobsStarted
}

func (s *runtimeState) eventsWereStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.eventsStarted
}

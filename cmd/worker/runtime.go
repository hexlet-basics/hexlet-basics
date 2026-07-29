package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	errRuntimeSignal      = errors.New("runtime interrupted by signal")
	errJobQueueStopped    = errors.New("job queue stopped unexpectedly")
	errEventRouterStopped = errors.New("domain event router stopped unexpectedly")
)

type jobRuntime interface {
	Start(context.Context) error
	Stop(context.Context) error
	Stopped() <-chan struct{}
}

type eventRuntime interface {
	Run(context.Context) error
	Close() error
}

// application owns the asynchronous process runtime. River starts before
// Watermill so event handlers can enqueue immediately; shutdown reverses that
// order so no new jobs arrive while workers drain.
type application struct {
	jobs   jobRuntime
	events eventRuntime
	logger *slog.Logger
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
		if a.logger != nil {
			a.logger.Info("async worker started")
		}
		if err := a.events.Run(processCtx); err != nil {
			return fmt.Errorf("run domain event router: %w", err)
		}
		return errEventRouterStopped
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
		a.logger.Info("shutting down async worker")
	}

	var shutdownErrors []error
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
		}
	}

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

	jobsStarted   bool
	eventsStarted bool
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

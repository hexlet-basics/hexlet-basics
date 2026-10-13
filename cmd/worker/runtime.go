package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"

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
	// Running is closed once the router has started every handler.
	Running() <-chan struct{}
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
	// The process context keeps the signal context's values but not its
	// cancellation: shutdown has to run after the signal has fired.
	processCtx, cancelProcess := context.WithCancel(context.WithoutCancel(signalCtx))
	defer cancelProcess()

	group, groupCtx := errgroup.WithContext(processCtx)
	riverReady := make(chan struct{})
	// jobsStartSettled is closed once starting the job queue has either succeeded
	// or failed. A shutdown waits for it before asking whether the queue was
	// started: a signal that lands in the middle of the start would otherwise skip
	// stopping a queue that then finishes starting, and nothing would ever stop it.
	jobsStartSettled := make(chan struct{})
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
		err := a.jobs.Start(processCtx)
		if err == nil {
			state.markJobsStarted()
		}
		close(jobsStartSettled)
		if err != nil {
			return fmt.Errorf("start job queue: %w", err)
		}
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

		// The router counts as started only once it reports itself running.
		// Watermill's Close waits for the handlers to come up, so closing a
		// router that has not got there yet blocks for its whole close timeout;
		// a shutdown that arrives earlier relies on the cancelled process
		// context to stop it instead.
		go func() {
			select {
			case <-a.events.Running():
				state.markEventsStarted()
			case <-groupCtx.Done():
			}
		}()
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
		<-jobsStartSettled
		return a.shutdown(processCtx, state, cancelProcess)
	})

	return group.Wait()
}

func (a *application) shutdown(parent context.Context, state *runtimeState, cancelProcess context.CancelFunc) error {
	// The drain gets its own deadline: the process context is already
	// cancelled by the time shutdown runs.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), shutdownTimeout)
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
	jobsStarted   atomic.Bool
	eventsStarted atomic.Bool
}

func (s *runtimeState) markJobsStarted() { s.jobsStarted.Store(true) }

func (s *runtimeState) markEventsStarted() { s.eventsStarted.Store(true) }

func (s *runtimeState) jobsWereStarted() bool { return s.jobsStarted.Load() }

func (s *runtimeState) eventsWereStarted() bool { return s.eventsStarted.Load() }

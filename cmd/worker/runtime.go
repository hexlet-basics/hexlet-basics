package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"hexletbasics/internal/health"
)

var (
	errRuntimeSignal      = errors.New("runtime interrupted by signal")
	errJobQueueStopped    = errors.New("job queue stopped unexpectedly")
	errEventRouterStopped = errors.New("domain event router stopped unexpectedly")
)

// jobRuntime, eventRuntime and httpRuntime stay interfaces: the runtime tests
// substitute hand-written fakes (stuckJobs, stuckServer, healthServer).
type jobRuntime interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Stopped() <-chan struct{}
}

// httpRuntime is the health-only listener.
type httpRuntime interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Close() error
}

type eventRuntime interface {
	Run(ctx context.Context) error
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
	// health, when set, is the health-only listener (/healthz, /readyz). It
	// starts with the runtime and is stopped last, so /readyz can report 503
	// while the queue and router drain.
	health httpRuntime
	probes *health.Probes
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
	state := &runtimeState{
		jobsStarted:   atomic.Bool{},
		eventsStarted: atomic.Bool{},
	}

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
			return errs.Wrapf(err, "start job queue")
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

				if a.probes != nil {
					a.probes.Started()
				}
			case <-groupCtx.Done():
			}
		}()

		if a.logger != nil {
			a.logger.InfoContext(groupCtx, "async worker started")
		}

		if err := a.events.Run(processCtx); err != nil {
			return errs.Wrapf(err, "run domain event router")
		}

		return errEventRouterStopped
	})

	if a.health != nil {
		group.Go(func() error {
			if err := a.health.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return errs.Wrapf(err, "serve health")
			}
			// A closed listener is the shutdown's doing; wait for the group to
			// finish rather than ending it early.
			<-groupCtx.Done()

			return nil
		})
	}

	// The group keeps only its first error, and on a signal that is the
	// sentinel, so the shutdown's own failure is kept apart and joined in.
	var shutdownErr error

	group.Go(func() error {
		<-groupCtx.Done()
		<-jobsStartSettled

		shutdownErr = a.shutdown(processCtx, state, cancelProcess)

		return shutdownErr
	})

	waitErr := group.Wait()
	if shutdownErr == nil || errors.Is(waitErr, shutdownErr) {
		return errs.Wrapf(waitErr, "wait for the runtime")
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
func (a *application) shutdown(parent context.Context, state *runtimeState, cancelProcess context.CancelFunc) error {
	// The drain gets its own deadline: the process context is already
	// cancelled by the time shutdown runs.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), shutdownTimeout)
	defer cancel()

	if a.probes != nil {
		a.probes.Drain()
	}

	if a.logger != nil {
		a.logger.InfoContext(ctx, "shutting down async worker")
	}

	var shutdownErrors []error

	if state.eventsWereStarted() {
		if err := a.events.Close(); err != nil {
			stopErr := errs.Wrapf(err, "stop domain event router")
			shutdownErrors = append(shutdownErrors, stopErr)
		}
	}

	if state.jobsWereStarted() {
		if err := a.jobs.Stop(ctx); err != nil {
			stopErr := errs.Wrapf(err, "stop job queue")
			shutdownErrors = append(shutdownErrors, stopErr)
		}
	}

	if a.health != nil {
		if err := a.health.Shutdown(ctx); err != nil {
			stopErr := errs.Wrapf(err, "stop health listener")
			shutdownErrors = append(shutdownErrors, stopErr)
			_ = a.health.Close()
		}
	}

	cancelProcess()

	return errors.Join(shutdownErrors...)
}

type runtimeState struct {
	jobsStarted   atomic.Bool
	eventsStarted atomic.Bool
}

func (s *runtimeState) markJobsStarted() { s.jobsStarted.Store(true) }

func (s *runtimeState) markEventsStarted() { s.eventsStarted.Store(true) }

func (s *runtimeState) jobsWereStarted() bool { return s.jobsStarted.Load() }

func (s *runtimeState) eventsWereStarted() bool { return s.eventsStarted.Load() }

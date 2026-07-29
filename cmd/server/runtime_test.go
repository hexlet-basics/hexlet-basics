package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeStopsRiverAndDoesNotStartHTTPWhenWatermillFails(t *testing.T) {
	watermillErr := errors.New("subscriber startup failed")
	jobs := newFakeJobRuntime()
	events := &fakeEventRuntime{
		runErr:  watermillErr,
		running: make(chan struct{}),
	}
	httpServer := &fakeHTTPRuntime{}
	app := application{
		http:   httpServer,
		jobs:   jobs,
		events: events,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.run(context.Background(), func() {})
	}()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, watermillErr)
	case <-time.After(time.Second):
		t.Fatal("runtime did not stop after Watermill failed")
	}

	assert.True(t, jobs.wasStopped())
	assert.False(t, httpServer.wasStarted())
}

func TestRuntimeStartsAndStopsComponentsInOrder(t *testing.T) {
	recorder := &lifecycleRecorder{}
	httpServer := newOrderedHTTPRuntime(recorder)
	jobs := newOrderedJobRuntime(recorder)
	events := newOrderedEventRuntime(recorder)
	app := application{
		http:   httpServer,
		jobs:   jobs,
		events: events,
	}
	signalCtx, signal := context.WithCancel(context.Background())
	stopSignalsCalled := make(chan struct{})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.run(signalCtx, func() { close(stopSignalsCalled) })
	}()

	select {
	case <-httpServer.started:
	case <-time.After(time.Second):
		t.Fatal("HTTP server did not start")
	}
	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(time.Second):
		t.Fatal("runtime did not stop after signal")
	}

	assert.Equal(t, []string{
		"river-start",
		"watermill-start",
		"http-start",
		"http-stop",
		"watermill-stop",
		"river-stop",
	}, recorder.events())
	assertClosed(t, stopSignalsCalled)
}

func TestRuntimeContinuesStagedShutdownAfterHTTPDrainFails(t *testing.T) {
	recorder := &lifecycleRecorder{}
	httpServer := newOrderedHTTPRuntime(recorder)
	httpServer.shutdownErr = errors.New("drain timed out")
	jobs := newOrderedJobRuntime(recorder)
	events := newOrderedEventRuntime(recorder)
	app := application{
		http:   httpServer,
		jobs:   jobs,
		events: events,
	}
	state := &runtimeState{}
	state.markJobsStarted()
	state.markEventsStarted()
	state.markHTTPStarted()

	err := app.shutdown(state, func() {})

	require.ErrorIs(t, err, httpServer.shutdownErr)
	assert.Equal(t, []string{
		"http-stop",
		"http-force-stop",
		"watermill-stop",
		"river-stop",
	}, recorder.events())
}

func assertClosed(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	default:
		t.Fatal("channel was not closed")
	}
}

type lifecycleRecorder struct {
	mu     sync.Mutex
	values []string
}

func (r *lifecycleRecorder) record(value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values = append(r.values, value)
}

func (r *lifecycleRecorder) events() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.values...)
}

type orderedJobRuntime struct {
	recorder *lifecycleRecorder
	stopped  chan struct{}
	stopOnce sync.Once
}

func newOrderedJobRuntime(recorder *lifecycleRecorder) *orderedJobRuntime {
	return &orderedJobRuntime{recorder: recorder, stopped: make(chan struct{})}
}

func (r *orderedJobRuntime) Start(context.Context) error {
	r.recorder.record("river-start")
	return nil
}

func (r *orderedJobRuntime) Stop(context.Context) error {
	r.recorder.record("river-stop")
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

func (r *orderedJobRuntime) Stopped() <-chan struct{} { return r.stopped }

type orderedEventRuntime struct {
	recorder *lifecycleRecorder
	running  chan struct{}
	closed   chan struct{}
	runOnce  sync.Once
	stopOnce sync.Once
}

func newOrderedEventRuntime(recorder *lifecycleRecorder) *orderedEventRuntime {
	return &orderedEventRuntime{
		recorder: recorder,
		running:  make(chan struct{}),
		closed:   make(chan struct{}),
	}
}

func (r *orderedEventRuntime) Run(context.Context) error {
	r.recorder.record("watermill-start")
	r.runOnce.Do(func() { close(r.running) })
	<-r.closed
	return nil
}

func (r *orderedEventRuntime) Running() <-chan struct{} { return r.running }

func (r *orderedEventRuntime) Close() error {
	r.recorder.record("watermill-stop")
	r.stopOnce.Do(func() { close(r.closed) })
	return nil
}

type orderedHTTPRuntime struct {
	recorder    *lifecycleRecorder
	started     chan struct{}
	stopped     chan struct{}
	stopOnce    sync.Once
	shutdownErr error
}

func newOrderedHTTPRuntime(recorder *lifecycleRecorder) *orderedHTTPRuntime {
	return &orderedHTTPRuntime{
		recorder: recorder,
		started:  make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

func (r *orderedHTTPRuntime) ListenAndServe() error {
	r.recorder.record("http-start")
	close(r.started)
	<-r.stopped
	return http.ErrServerClosed
}

func (r *orderedHTTPRuntime) Shutdown(context.Context) error {
	r.recorder.record("http-stop")
	if r.shutdownErr != nil {
		return r.shutdownErr
	}
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

func (r *orderedHTTPRuntime) Close() error {
	r.recorder.record("http-force-stop")
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

type fakeJobRuntime struct {
	stopped  chan struct{}
	stopOnce sync.Once

	mu      sync.Mutex
	started bool
	stop    bool
}

func newFakeJobRuntime() *fakeJobRuntime {
	return &fakeJobRuntime{stopped: make(chan struct{})}
}

func (r *fakeJobRuntime) Start(context.Context) error {
	r.mu.Lock()
	r.started = true
	r.mu.Unlock()
	return nil
}

func (r *fakeJobRuntime) Stop(context.Context) error {
	r.mu.Lock()
	r.stop = true
	r.mu.Unlock()
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

func (r *fakeJobRuntime) Stopped() <-chan struct{} { return r.stopped }

func (r *fakeJobRuntime) wasStopped() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stop
}

type fakeEventRuntime struct {
	runErr  error
	running chan struct{}
}

func (r *fakeEventRuntime) Run(context.Context) error { return r.runErr }
func (r *fakeEventRuntime) Running() <-chan struct{}  { return r.running }
func (r *fakeEventRuntime) Close() error              { return nil }

type fakeHTTPRuntime struct {
	mu      sync.Mutex
	started bool
}

func (r *fakeHTTPRuntime) ListenAndServe() error {
	r.mu.Lock()
	r.started = true
	r.mu.Unlock()
	return http.ErrServerClosed
}

func (r *fakeHTTPRuntime) Shutdown(context.Context) error { return nil }
func (r *fakeHTTPRuntime) Close() error                   { return nil }

func (r *fakeHTTPRuntime) wasStarted() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.started
}

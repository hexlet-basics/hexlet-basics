package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeStopsRiverWhenWatermillFails(t *testing.T) {
	watermillErr := errors.New("subscriber startup failed")
	recorder := &lifecycleRecorder{}
	jobs := newFakeJobRuntime(recorder)
	events := &fakeEventRuntime{recorder: recorder, runErr: watermillErr}
	app := application{jobs: jobs, events: events}

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.run(context.Background(), func() {})
	}()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, watermillErr)
	case <-time.After(time.Second):
		t.Fatal("worker did not stop after Watermill failed")
	}

	assert.Equal(t, []string{
		"river-start",
		"watermill-start",
		"watermill-stop",
		"river-stop",
	}, recorder.events())
}

func TestRuntimeStartsAndStopsAsyncComponentsInOrder(t *testing.T) {
	recorder := &lifecycleRecorder{}
	jobs := newFakeJobRuntime(recorder)
	events := &fakeEventRuntime{
		recorder: recorder,
		started:  make(chan struct{}),
		closed:   make(chan struct{}),
	}
	app := application{jobs: jobs, events: events}
	signalCtx, signal := context.WithCancel(context.Background())
	stopSignalsCalled := make(chan struct{})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.run(signalCtx, func() { close(stopSignalsCalled) })
	}()

	select {
	case <-events.started:
	case <-time.After(time.Second):
		t.Fatal("Watermill did not start")
	}
	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(time.Second):
		t.Fatal("worker did not stop after signal")
	}

	assert.Equal(t, []string{
		"river-start",
		"watermill-start",
		"watermill-stop",
		"river-stop",
	}, recorder.events())
	assertClosed(t, stopSignalsCalled)
}

func TestRuntimeDoesNotStartWatermillWhenRiverFails(t *testing.T) {
	startErr := errors.New("database unavailable")
	recorder := &lifecycleRecorder{}
	jobs := newFakeJobRuntime(recorder)
	jobs.startErr = startErr
	events := &fakeEventRuntime{recorder: recorder}
	app := application{jobs: jobs, events: events}

	err := app.run(context.Background(), func() {})

	require.ErrorIs(t, err, startErr)
	assert.Equal(t, []string{"river-start"}, recorder.events())
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

type fakeJobRuntime struct {
	recorder *lifecycleRecorder
	stopped  chan struct{}
	stopOnce sync.Once
	startErr error
}

func newFakeJobRuntime(recorder *lifecycleRecorder) *fakeJobRuntime {
	return &fakeJobRuntime{recorder: recorder, stopped: make(chan struct{})}
}

func (r *fakeJobRuntime) Start(context.Context) error {
	r.recorder.record("river-start")
	return r.startErr
}

func (r *fakeJobRuntime) Stop(context.Context) error {
	r.recorder.record("river-stop")
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

func (r *fakeJobRuntime) Stopped() <-chan struct{} { return r.stopped }

type fakeEventRuntime struct {
	recorder  *lifecycleRecorder
	runErr    error
	started   chan struct{}
	closed    chan struct{}
	startOnce sync.Once
	closeOnce sync.Once
}

func (r *fakeEventRuntime) Run(context.Context) error {
	r.recorder.record("watermill-start")
	if r.started != nil {
		r.startOnce.Do(func() { close(r.started) })
	}
	if r.runErr != nil {
		return r.runErr
	}
	if r.closed != nil {
		<-r.closed
	}
	return nil
}

func (r *fakeEventRuntime) Close() error {
	r.recorder.record("watermill-stop")
	if r.closed != nil {
		r.closeOnce.Do(func() { close(r.closed) })
	}
	return nil
}

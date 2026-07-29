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

func TestRuntimeStartsHTTPAndStopsOnSignal(t *testing.T) {
	httpServer := newFakeHTTPRuntime()
	app := application{http: httpServer}
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

	assert.True(t, httpServer.wasStopped())
	assertClosed(t, stopSignalsCalled)
}

func TestRuntimeForceClosesHTTPAfterDrainFailure(t *testing.T) {
	drainErr := errors.New("drain timed out")
	httpServer := newFakeHTTPRuntime()
	httpServer.shutdownErr = drainErr
	app := application{http: httpServer}

	err := app.shutdown(func() {})

	require.ErrorIs(t, err, drainErr)
	assert.True(t, httpServer.wasForceClosed())
}

func assertClosed(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	default:
		t.Fatal("channel was not closed")
	}
}

type fakeHTTPRuntime struct {
	started chan struct{}
	stopped chan struct{}

	startOnce   sync.Once
	stopOnce    sync.Once
	mu          sync.Mutex
	shutdownErr error
	forceClosed bool
}

func newFakeHTTPRuntime() *fakeHTTPRuntime {
	return &fakeHTTPRuntime{
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (r *fakeHTTPRuntime) ListenAndServe() error {
	r.startOnce.Do(func() { close(r.started) })
	<-r.stopped
	return http.ErrServerClosed
}

func (r *fakeHTTPRuntime) Shutdown(context.Context) error {
	if r.shutdownErr != nil {
		return r.shutdownErr
	}
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

func (r *fakeHTTPRuntime) Close() error {
	r.mu.Lock()
	r.forceClosed = true
	r.mu.Unlock()
	r.stopOnce.Do(func() { close(r.stopped) })
	return nil
}

func (r *fakeHTTPRuntime) wasStopped() bool {
	select {
	case <-r.stopped:
		return true
	default:
		return false
	}
}

func (r *fakeHTTPRuntime) wasForceClosed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.forceClosed
}

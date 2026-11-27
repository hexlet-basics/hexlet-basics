package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/health"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errJobQueueCannotStart = errors.New("job queue cannot start")
)

const eventually = 10 * time.Second

// fakeJobs is the job queue as the runtime sees it: it starts, it stops, and it
// says when it has stopped. The runtime's ordering is what the tests are about,
// so nothing here needs a database.
type fakeJobs struct {
	startErr error
	started  chan struct{}
	stopOnce sync.Once
	stopped  chan struct{}
}

func newJobs() *fakeJobs {
	return &fakeJobs{
		started:  make(chan struct{}),
		stopped:  make(chan struct{}),
		startErr: nil,
		stopOnce: sync.Once{},
	}
}

func (f *fakeJobs) Start(context.Context) error {
	if f.startErr != nil {
		return f.startErr
	}

	close(f.started)

	return nil
}

func (f *fakeJobs) Stop(context.Context) error {
	f.stopOnce.Do(func() { close(f.stopped) })

	return nil
}

func (f *fakeJobs) Stopped() <-chan struct{} { return f.stopped }

func stopped(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func TestRuntimeStopsTheJobQueueOnSignal(t *testing.T) {
	jobQueue := newJobs()
	app := application{
		jobs:   jobQueue,
		logger: nil,
		health: nil,
		probes: nil,
	}
	signalCtx, signal := context.WithCancel(context.Background())

	var stopSignalsCalled bool

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() { stopSignalsCalled = true }) }()

	select {
	case <-jobQueue.started:
	case <-time.After(eventually):
		t.Fatal("the job queue did not start")
	}

	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(eventually):
		t.Fatal("worker did not stop after signal")
	}

	assert.True(t, stopped(jobQueue.Stopped()), "the job queue is stopped")
	assert.True(t, stopSignalsCalled)
}

// A queue that stops with no signal asking it to is a failure of the process,
// which then exits rather than serve health for a worker that runs nothing.
func TestRuntimeFailsWhenTheJobQueueStopsOnItsOwn(t *testing.T) {
	jobQueue := newJobs()
	app := application{
		jobs:   jobQueue,
		logger: nil,
		health: nil,
		probes: nil,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(context.Background(), func() {}) }()

	select {
	case <-jobQueue.started:
	case <-time.After(eventually):
		t.Fatal("the job queue did not start")
	}

	require.NoError(t, jobQueue.Stop(t.Context()))

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errJobQueueStopped)
	case <-time.After(eventually):
		t.Fatal("worker did not stop after the job queue stopped")
	}
}

func TestRuntimeReportsAJobQueueThatCannotStart(t *testing.T) {
	jobQueue := newJobs()
	jobQueue.startErr = errJobQueueCannotStart
	probes := health.NewAwaitingStart(func(context.Context) error { return nil })
	app := application{
		jobs:   jobQueue,
		logger: nil,
		health: nil,
		probes: probes,
	}

	err := app.run(context.Background(), func() {})

	require.ErrorContains(t, err, "job queue cannot start")
	assert.Equal(t, http.StatusServiceUnavailable, readyz(t, probes.Wrap(http.NotFoundHandler())), "never ready")
}

// readyz is the worker's readiness as the kubelet sees it: the health handler
// answering in-process.
func readyz(t *testing.T, handler http.Handler) int {
	t.Helper()

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil))

	return rec.Code
}

func TestReadyzIsUnavailableBeforeTheRuntimeStartsOKWhileItRunsAndUnavailableAfterItStops(t *testing.T) {
	probes := health.NewAwaitingStart(func(context.Context) error { return nil })
	handler := probes.Wrap(http.NotFoundHandler())
	app := application{
		jobs:   newJobs(),
		probes: probes,
		logger: nil,
		health: nil,
	}

	require.Equal(t, http.StatusServiceUnavailable, readyz(t, handler), "before start")

	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()

	require.Eventually(t, func() bool { return readyz(t, handler) == http.StatusOK }, eventually, 10*time.Millisecond, "while running")
	signal()
	require.ErrorIs(t, <-errCh, errRuntimeSignal)

	assert.Equal(t, http.StatusServiceUnavailable, readyz(t, handler), "after stop")
}

// healthServer is a health listener that records being shut down.
type healthServer struct {
	served   chan struct{}
	shutdown atomic.Bool
}

func (h *healthServer) ListenAndServe() error {
	close(h.served)

	for !h.shutdown.Load() {
		time.Sleep(time.Millisecond)
	}

	return http.ErrServerClosed
}

func (h *healthServer) Shutdown(context.Context) error {
	h.shutdown.Store(true)

	return nil
}
func (h *healthServer) Close() error {
	h.shutdown.Store(true)

	return nil
}

func TestHealthListenerStartsAndStopsWithTheWorker(t *testing.T) {
	listener := &healthServer{
		served:   make(chan struct{}),
		shutdown: atomic.Bool{},
	}
	app := application{
		jobs:   newJobs(),
		health: listener,
		logger: nil,
		probes: nil,
	}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()

	select {
	case <-listener.served:
	case <-time.After(eventually):
		t.Fatal("the health listener did not start")
	}

	signal()
	require.ErrorIs(t, <-errCh, errRuntimeSignal)

	assert.True(t, listener.shutdown.Load(), "the health listener is stopped")
}

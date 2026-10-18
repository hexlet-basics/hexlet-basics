package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/health"
)

const eventually = 10 * time.Second

// fakeJobs is the job queue as the runtime sees it: it starts, it stops, and it
// says when it has stopped. The runtime's ordering is what the tests are about,
// so nothing here needs a database.
type fakeJobs struct {
	startErr error
	stopOnce sync.Once
	stopped  chan struct{}
}

func newJobs() *fakeJobs { return &fakeJobs{stopped: make(chan struct{})} }

func (f *fakeJobs) Start(context.Context) error { return f.startErr }

func (f *fakeJobs) Stop(context.Context) error {
	f.stopOnce.Do(func() { close(f.stopped) })
	return nil
}

func (f *fakeJobs) Stopped() <-chan struct{} { return f.stopped }

// watermillRouter is a Watermill router as the application's event runtime: the
// router reports readiness on a bidirectional channel, the runtime on a
// receive-only one.
type watermillRouter struct{ *message.Router }

func (r watermillRouter) Running() <-chan struct{} { return r.Router.Running() }

// newRouter is a real Watermill router with one handler over an in-memory
// pub/sub.
func newRouter(t *testing.T) watermillRouter {
	t.Helper()
	logger := watermill.NewSlogLogger(slog.New(slog.DiscardHandler))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	router, err := message.NewRouter(message.RouterConfig{CloseTimeout: time.Second}, logger)
	require.NoError(t, err)
	router.AddConsumerHandler("test_consumer", "test_topic", pubsub, func(*message.Message) error { return nil })
	t.Cleanup(func() { _ = router.Close() })
	return watermillRouter{router}
}

func stopped(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func TestRuntimeStopsBothComponentsOnSignal(t *testing.T) {
	jobQueue := newJobs()
	router := newRouter(t)
	app := application{jobs: jobQueue, events: router}
	signalCtx, signal := context.WithCancel(context.Background())
	var stopSignalsCalled bool

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() { stopSignalsCalled = true }) }()

	select {
	case <-router.Running():
	case <-time.After(eventually):
		t.Fatal("the event router did not start")
	}
	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(eventually):
		t.Fatal("worker did not stop after signal")
	}

	assert.True(t, stopped(jobQueue.Stopped()), "the job queue is stopped")
	assert.True(t, router.IsClosed(), "the event router is closed")
	assert.True(t, stopSignalsCalled)
}

func TestRuntimeStopsRiverWhenWatermillFails(t *testing.T) {
	jobQueue := newJobs()
	router := newRouter(t)
	// A router that is already running cannot be run again, so Run fails.
	go func() { _ = router.Run(t.Context()) }()
	select {
	case <-router.Running():
	case <-time.After(eventually):
		t.Fatal("the event router did not start")
	}
	app := application{jobs: jobQueue, events: router}

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(context.Background(), func() {}) }()

	select {
	case err := <-errCh:
		require.ErrorContains(t, err, "router is already running")
	case <-time.After(eventually):
		t.Fatal("worker did not stop after the event router failed")
	}

	assert.True(t, stopped(jobQueue.Stopped()), "the job queue is stopped")
}

func TestRuntimeDoesNotStartWatermillWhenRiverFails(t *testing.T) {
	jobQueue := newJobs()
	jobQueue.startErr = errors.New("job queue cannot start")
	router := newRouter(t)
	app := application{jobs: jobQueue, events: router}

	err := app.run(context.Background(), func() {})

	require.ErrorContains(t, err, "job queue cannot start")
	select {
	case <-router.Running():
		t.Fatal("the event router started although the job queue failed")
	default:
	}
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
	router := newRouter(t)
	app := application{jobs: newJobs(), events: router, probes: probes}
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

func (h *healthServer) Shutdown(context.Context) error { h.shutdown.Store(true); return nil }
func (h *healthServer) Close() error                   { h.shutdown.Store(true); return nil }

func TestHealthListenerStartsAndStopsWithTheWorker(t *testing.T) {
	listener := &healthServer{served: make(chan struct{})}
	app := application{jobs: newJobs(), events: newRouter(t), health: listener}
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

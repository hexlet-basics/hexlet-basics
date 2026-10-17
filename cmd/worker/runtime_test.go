package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	logger := watermill.NewSlogLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
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

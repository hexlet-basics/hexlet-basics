package main

import (
	"bytes"
	"context"
	"errors"
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

// Faults a test injects; each is told apart from the others by identity.
var (
	errStop = errors.New("workers still busy")
)

// logSink is a log destination the runtime's goroutines and the test share.
type logSink struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *logSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.buf.Write(p)
	if err != nil {
		return result, errs.Wrapf(err, "write log")
	}

	return result, nil
}

func (s *logSink) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.buf.String()
}

func (s *logSink) logger() *slog.Logger { return slog.New(slog.NewTextHandler(s, nil)) }

func TestRuntimeAnnouncesTheStartAndTheShutdown(t *testing.T) {
	router := newRouter(t)
	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	app := application{
		jobs:   newJobs(),
		events: router,
		logger: logs.logger(),
		health: nil,
		probes: nil,
	}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()

	select {
	case <-router.Running():
	case <-time.After(eventually):
		t.Fatal("the event router did not start")
	}

	signal()

	require.ErrorIs(t, <-errCh, errRuntimeSignal)
	assert.Contains(t, logs.String(), "async worker started")
	assert.Contains(t, logs.String(), "shutting down async worker")
}

// TestRuntimeReportsAnEventRouterThatCannotDrain starts a router whose handler
// never returns, so closing it runs into its close timeout.
func TestRuntimeReportsAnEventRouterThatCannotDrain(t *testing.T) {
	logger := watermill.NopLogger{}
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	router, err := message.NewRouter(message.RouterConfig{CloseTimeout: 200 * time.Millisecond}, logger)
	require.NoError(t, err)

	release := make(chan struct{})

	t.Cleanup(func() { close(release) })

	inFlight := make(chan struct{})

	var once sync.Once

	router.AddConsumerHandler("stuck", "stuck_topic", pubsub, func(*message.Message) error {
		once.Do(func() { close(inFlight) })
		<-release

		return nil
	})

	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	app := application{
		jobs:   newJobs(),
		events: watermillRouter{router},
		logger: logs.logger(),
		health: nil,
		probes: nil,
	}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()

	select {
	case <-router.Running():
	case <-time.After(eventually):
		t.Fatal("the event router did not start")
	}

	require.NoError(t, pubsub.Publish("stuck_topic", message.NewMessage("m-1", nil)))

	select {
	case <-inFlight:
	case <-time.After(eventually):
		t.Fatal("the message did not reach the handler")
	}

	signal()

	select {
	case err := <-errCh:
		// The signal that started the shutdown and what failed while draining
		// are both returned; main logs them once.
		require.ErrorIs(t, err, errRuntimeSignal)
		require.ErrorContains(t, err, "stop domain event router")
	case <-time.After(eventually):
		t.Fatal("worker did not stop")
	}

	assert.NotContains(t, logs.String(), "stop domain event router", "returned, not logged")
}

// stuckJobs is a job queue that starts, then cannot be stopped.
type stuckJobs struct {
	stopErr error
	done    chan struct{}
}

func (stuckJobs) Start(context.Context) error { return nil }

func (j stuckJobs) Stop(context.Context) error {
	close(j.done)

	return j.stopErr
}

func (j stuckJobs) Stopped() <-chan struct{} { return j.done }

func TestRuntimeReportsAJobQueueThatCannotStop(t *testing.T) {
	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	jobs := stuckJobs{stopErr: errStop, done: make(chan struct{})}
	app := application{
		jobs:   jobs,
		events: newRouter(t),
		logger: logs.logger(),
		health: nil,
		probes: nil,
	}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()

	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
		require.ErrorIs(t, err, errStop)
		require.ErrorContains(t, err, "stop job queue")
	case <-time.After(eventually):
		t.Fatal("worker did not stop")
	}

	assert.NotContains(t, logs.String(), "workers still busy", "returned, not logged")
}

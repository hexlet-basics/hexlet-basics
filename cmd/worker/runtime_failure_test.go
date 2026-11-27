package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

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
	jobQueue := newJobs()
	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	app := application{
		jobs:   jobQueue,
		logger: logs.logger(),
		health: nil,
		probes: nil,
	}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()

	require.Eventually(t, func() bool { return strings.Contains(logs.String(), "async worker started") },
		eventually, 10*time.Millisecond, "the start is announced once the queue is up")

	signal()

	require.ErrorIs(t, <-errCh, errRuntimeSignal)
	assert.Contains(t, logs.String(), "shutting down async worker")
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

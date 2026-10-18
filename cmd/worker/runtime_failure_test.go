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

	"hexletbasics/internal/di"
)

// logSink is a log destination the runtime's goroutines and the test share.
type logSink struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *logSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *logSink) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *logSink) logger() *slog.Logger { return slog.New(slog.NewTextHandler(s, nil)) }

func TestRuntimeAnnouncesTheStartAndTheShutdown(t *testing.T) {
	router := newRouter(t)
	logs := &logSink{}
	app := application{jobs: newJobs(), events: router, logger: logs.logger()}
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
	logs := &logSink{}
	app := application{jobs: newJobs(), events: watermillRouter{router}, logger: logs.logger()}
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
		// The signal that started the shutdown is what the process reports; what
		// failed while draining is logged.
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(eventually):
		t.Fatal("worker did not stop")
	}
	assert.Contains(t, logs.String(), "runtime shutdown failed")
	assert.Contains(t, logs.String(), "stop domain event router")
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
	stopErr := errors.New("workers still busy")
	logs := &logSink{}
	jobs := stuckJobs{stopErr: stopErr, done: make(chan struct{})}
	app := application{jobs: jobs, events: newRouter(t), logger: logs.logger()}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()
	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(eventually):
		t.Fatal("worker did not stop")
	}
	assert.Contains(t, logs.String(), "stop job queue")
	assert.Contains(t, logs.String(), "workers still busy")
}

func workerDependencies(t *testing.T) di.WorkerDependencies {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://nobody@127.0.0.1:1/none")
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("PUBLIC_URL", "http://localhost:3001")
	t.Setenv("SITE_URL", "http://localhost:5173")
	t.Setenv("EMAIL_TOKEN_SECRET", "test-email-secret")
	dependencies, err := di.BuildWorker()
	require.NoError(t, err)
	return dependencies
}

func TestCloseResourcesReleasesEverythingTheProcessOwns(t *testing.T) {
	dependencies := workerDependencies(t)
	logs := &logSink{}

	closeResources(t.Context(), logs.logger(),
		dependencies.Bucket, dependencies.Database, dependencies.OpenTelemetrySDK, dependencies.SentryClient)

	assert.NotContains(t, logs.String(), "level=ERROR")
	assert.Error(t, dependencies.Bucket.WriteAll(t.Context(), "key", []byte("x"), nil), "the bucket is closed")
	_, err := dependencies.Database.User.Query().Count(t.Context())
	assert.Error(t, err, "the database pool is closed")
}

func TestCloseResourcesLogsWhatFailsToClose(t *testing.T) {
	dependencies := workerDependencies(t)
	closeResources(t.Context(), (&logSink{}).logger(),
		dependencies.Bucket, dependencies.Database, dependencies.OpenTelemetrySDK, dependencies.SentryClient)
	logs := &logSink{}

	closeResources(t.Context(), logs.logger(),
		dependencies.Bucket, dependencies.Database, dependencies.OpenTelemetrySDK, dependencies.SentryClient)

	assert.Contains(t, logs.String(), "closing blob bucket")
}

package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

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

func TestRuntimeAnnouncesItsAddressAndTheShutdown(t *testing.T) {
	addr := freeAddr(t)
	logs := &logSink{}
	app := application{
		http:     &http.Server{Addr: addr, Handler: http.NotFoundHandler()},
		httpAddr: addr,
		logger:   logs.logger(),
	}
	signalCtx, signal := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()
	require.Eventually(t, func() bool {
		res, err := get(t.Context(), "http://"+addr)
		if err != nil {
			return false
		}
		_ = res.Body.Close()
		return true
	}, eventually, 20*time.Millisecond)
	signal()

	require.ErrorIs(t, <-errCh, errRuntimeSignal)
	assert.Contains(t, logs.String(), "backend listening")
	assert.Contains(t, logs.String(), addr)
	assert.Contains(t, logs.String(), "shutting down")
}

func TestRuntimeReportsAnAddressThatCannotBeServed(t *testing.T) {
	taken, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = taken.Close() })
	addr := taken.Addr().String()
	app := application{http: &http.Server{Addr: addr}, httpAddr: addr}

	err = app.run(context.Background(), func() {})

	require.ErrorContains(t, err, "serve HTTP")
}

func TestRuntimeReportsAServerThatStopsOnItsOwn(t *testing.T) {
	addr := freeAddr(t)
	server := &http.Server{Addr: addr, Handler: http.NotFoundHandler()}
	app := application{http: server, httpAddr: addr}

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(context.Background(), func() {}) }()
	require.Eventually(t, func() bool {
		res, err := get(t.Context(), "http://"+addr)
		if err != nil {
			return false
		}
		_ = res.Body.Close()
		return true
	}, eventually, 20*time.Millisecond)
	require.NoError(t, server.Close())

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errHTTPServerStopped)
	case <-time.After(eventually):
		t.Fatal("runtime did not notice the server had stopped")
	}
}

// stuckServer is an HTTP runtime that can neither drain nor be closed.
type stuckServer struct{ shutdownErr, closeErr error }

func (stuckServer) ListenAndServe() error { select {} }

func (s stuckServer) Shutdown(context.Context) error { return s.shutdownErr }

func (s stuckServer) Close() error { return s.closeErr }

func TestShutdownReportsBothFailuresWhenTheServerCannotBeStopped(t *testing.T) {
	drainErr := errors.New("drain failed")
	closeErr := errors.New("close failed")
	logs := &logSink{}
	app := application{http: stuckServer{shutdownErr: drainErr, closeErr: closeErr}, logger: logs.logger()}

	err := app.shutdown(t.Context(), func() {})

	require.ErrorIs(t, err, drainErr)
	require.ErrorIs(t, err, closeErr)
	assert.Contains(t, logs.String(), "stop HTTP server")
	assert.Contains(t, logs.String(), "force-close HTTP server")
}

func TestShutdownReportsTheDrainFailureWhenForceCloseWorks(t *testing.T) {
	drainErr := errors.New("drain failed")
	logs := &logSink{}
	app := application{http: stuckServer{shutdownErr: drainErr}, logger: logs.logger()}

	err := app.shutdown(t.Context(), func() {})

	require.ErrorIs(t, err, drainErr)
	assert.Contains(t, logs.String(), "stop HTTP server")
	assert.NotContains(t, logs.String(), "force-close")
}

func serverDependencies(t *testing.T) di.ServerDependencies {
	t.Helper()
	// The pool opens lazily, so the process's resources can be built and released
	// without a database behind them.
	t.Setenv("DATABASE_URL", "postgres://nobody@127.0.0.1:1/none")
	t.Setenv("BLOB_BUCKET_URL", "file://"+t.TempDir())
	t.Setenv("GUEST_COOKIE_SECRET", "test-secret")
	t.Setenv("PUBLIC_URL", "http://localhost:3001")
	t.Setenv("SITE_URL", "http://localhost:5173")
	t.Setenv("APP_HOST", "code-basics.com")
	t.Setenv("EMAIL_TOKEN_SECRET", "test-email-secret")
	dependencies, err := di.BuildServer()
	require.NoError(t, err)
	return dependencies
}

func TestCloseResourcesReleasesEverythingTheProcessOwns(t *testing.T) {
	dependencies := serverDependencies(t)
	logs := &logSink{}

	closeResources(t.Context(), logs.logger(),
		dependencies.Bucket, dependencies.Database, dependencies.OpenTelemetrySDK, dependencies.SentryClient)

	assert.NotContains(t, logs.String(), "level=ERROR")
	assert.Error(t, dependencies.Bucket.WriteAll(t.Context(), "key", []byte("x"), nil), "the bucket is closed")
	_, err := dependencies.Database.User.Query().Count(t.Context())
	assert.Error(t, err, "the database pool is closed")
}

func TestCloseResourcesLogsWhatFailsToClose(t *testing.T) {
	dependencies := serverDependencies(t)
	closeResources(t.Context(), (&logSink{}).logger(),
		dependencies.Bucket, dependencies.Database, dependencies.OpenTelemetrySDK, dependencies.SentryClient)
	logs := &logSink{}

	closeResources(t.Context(), logs.logger(),
		dependencies.Bucket, dependencies.Database, dependencies.OpenTelemetrySDK, dependencies.SentryClient)

	assert.Contains(t, logs.String(), "closing blob bucket")
}

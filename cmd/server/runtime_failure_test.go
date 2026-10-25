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
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errClose = errors.New("close failed")
	errDrain = errors.New("drain failed")
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

func TestRuntimeAnnouncesItsAddressAndTheShutdown(t *testing.T) {
	addr := freeAddr(t)
	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	app := application{
		http:            &http.Server{ReadHeaderTimeout: time.Second, Addr: addr, Handler: http.NotFoundHandler()},
		httpAddr:        addr,
		logger:          logs.logger(),
		drain:           nil,
		shutdownTimeout: 0,
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
	t.Cleanup(func() { require.NoError(t, taken.Close()) })

	addr := taken.Addr().String()
	app := application{
		http:            &http.Server{ReadHeaderTimeout: time.Second, Addr: addr},
		httpAddr:        addr,
		logger:          nil,
		drain:           nil,
		shutdownTimeout: 0,
	}

	err = app.run(context.Background(), func() {})

	require.ErrorContains(t, err, "serve HTTP")
}

func TestRuntimeReportsAServerThatStopsOnItsOwn(t *testing.T) {
	addr := freeAddr(t)
	server := &http.Server{ReadHeaderTimeout: time.Second, Addr: addr, Handler: http.NotFoundHandler()}
	app := application{
		http:            server,
		httpAddr:        addr,
		logger:          nil,
		drain:           nil,
		shutdownTimeout: 0,
	}

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
	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	app := application{
		http:            stuckServer{shutdownErr: errDrain, closeErr: errClose},
		logger:          logs.logger(),
		httpAddr:        "",
		drain:           nil,
		shutdownTimeout: 0,
	}

	err := app.shutdown(t.Context(), func() {})

	require.ErrorIs(t, err, errDrain)
	require.ErrorIs(t, err, errClose)
	require.ErrorContains(t, err, "stop HTTP server")
	require.ErrorContains(t, err, "force-close HTTP server")
	assert.NotContains(t, logs.String(), "runtime shutdown failed", "shutdown returns the failure and does not log it")
}

func TestShutdownReportsTheDrainFailureWhenForceCloseWorks(t *testing.T) {
	logs := &logSink{
		mu:  sync.Mutex{},
		buf: bytes.Buffer{},
	}
	app := application{
		http: stuckServer{
			shutdownErr: errDrain,
			closeErr:    nil,
		},
		logger:          logs.logger(),
		httpAddr:        "",
		drain:           nil,
		shutdownTimeout: 0,
	}

	err := app.shutdown(t.Context(), func() {})

	require.ErrorIs(t, err, errDrain)
	require.ErrorContains(t, err, "stop HTTP server")
	assert.NotContains(t, err.Error(), "force-close")
}

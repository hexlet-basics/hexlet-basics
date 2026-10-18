package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/health"
)

func TestRuntimeReportsNotReadyAfterTheSignal(t *testing.T) {
	addr := freeAddr(t)
	probes := health.New(func(context.Context) error { return nil })
	handler := probes.Wrap(http.NotFoundHandler())
	server := &http.Server{Addr: addr, Handler: handler}
	app := application{http: server, httpAddr: addr, drain: probes.Drain}
	signalCtx, signal := context.WithCancel(context.Background())
	readyz := func() int {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil))
		return rec.Code
	}
	require.Equal(t, http.StatusOK, readyz())

	errCh := make(chan error, 1)
	go func() { errCh <- app.run(signalCtx, func() {}) }()
	signal()
	require.ErrorIs(t, <-errCh, errRuntimeSignal)

	assert.Equal(t, http.StatusServiceUnavailable, readyz())
}

const eventually = 5 * time.Second

// freeAddr is a loopback address nothing listens on yet.
func freeAddr(t *testing.T) string {
	t.Helper()
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())
	return addr
}

func get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	return client.Do(req)
}

func TestRuntimeServesRequestsAndStopsOnSignal(t *testing.T) {
	addr := freeAddr(t)
	server := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})}
	app := application{http: server, httpAddr: addr}
	signalCtx, signal := context.WithCancel(context.Background())
	var stopSignalsCalled atomic.Bool

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.run(signalCtx, func() { stopSignalsCalled.Store(true) })
	}()

	require.Eventually(t, func() bool {
		res, err := get(t.Context(), "http://"+addr)
		if err != nil {
			return false
		}
		_ = res.Body.Close()
		return res.StatusCode == http.StatusNoContent
	}, eventually, 20*time.Millisecond, "the server answers requests")

	signal()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errRuntimeSignal)
	case <-time.After(eventually):
		t.Fatal("runtime did not stop after signal")
	}

	res, err := get(t.Context(), "http://"+addr)
	if res != nil {
		_ = res.Body.Close()
	}
	require.Error(t, err, "the port no longer accepts connections")
	assert.True(t, stopSignalsCalled.Load())
}

func TestRuntimeForceClosesHTTPAfterDrainFailure(t *testing.T) {
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	inFlight := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		close(inFlight)
		<-release
	})}
	go func() { _ = server.Serve(listener) }()
	app := application{http: server, shutdownTimeout: 200 * time.Millisecond}

	clientErr := make(chan error, 1)
	go func() {
		res, err := get(t.Context(), "http://"+listener.Addr().String())
		if err == nil {
			_ = res.Body.Close()
		}
		clientErr <- err
	}()
	select {
	case <-inFlight:
	case <-time.After(eventually):
		t.Fatal("request did not reach the handler")
	}

	err = app.shutdown(t.Context(), func() {})

	require.ErrorIs(t, err, context.DeadlineExceeded, "the drain gave up")
	select {
	case err := <-clientErr:
		require.Error(t, err, "the in-flight request was cut off")
	case <-time.After(eventually):
		t.Fatal("the in-flight request was not closed")
	}
}

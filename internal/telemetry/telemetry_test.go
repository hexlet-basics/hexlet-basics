package telemetry_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/config"
	"hexletbasics/internal/telemetry"
)

func TestNewTracerProviderIsNonRecordingWithoutEndpoint(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")

	provider, err := telemetry.NewTracerProvider(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		&config.Config{},
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, provider.Shutdown(t.Context())) })

	_, span := provider.Tracer("test").Start(t.Context(), "disabled")
	assert.False(t, span.IsRecording())
	span.End()
}

func TestNewTracerProviderExportsAndFlushesSpans(t *testing.T) {
	requests := make(chan *http.Request, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		requests <- r.Clone(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(collector.Close)

	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", collector.URL)
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	t.Setenv("OTEL_TRACES_SAMPLER", "always_on")

	provider, err := telemetry.NewTracerProvider(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		&config.Config{OTLPEndpoint: collector.URL},
	)
	require.NoError(t, err)

	_, span := provider.Tracer("test").Start(t.Context(), "exported")
	require.True(t, span.IsRecording())
	span.End()
	require.NoError(t, provider.Shutdown(t.Context()))

	select {
	case request := <-requests:
		assert.Equal(t, "/v1/traces", request.URL.Path)
		assert.Equal(t, "application/x-protobuf", request.Header.Get("Content-Type"))
	case <-time.After(time.Second):
		t.Fatal("OTLP collector did not receive a span batch")
	}
}

func TestNewTracerProviderRejectsInvalidEndpoint(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "://invalid")

	_, err := telemetry.NewTracerProvider(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		&config.Config{OTLPEndpoint: "://invalid"},
	)
	require.Error(t, err)
}

func TestNewSentryClientRejectsInvalidDSN(t *testing.T) {
	_, err := telemetry.NewSentryClient(&config.Config{SentryDSN: "not a dsn"})
	require.Error(t, err)
}

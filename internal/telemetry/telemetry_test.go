package telemetry_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	collectortrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	common "go.opentelemetry.io/proto/otlp/common/v1"
	"google.golang.org/protobuf/proto"

	"hexletbasics/internal/config"
	"hexletbasics/internal/telemetry"
)

func TestNewOpenTelemetrySDKIsNoopWithoutConfigFile(t *testing.T) {
	t.Setenv("OTEL_CONFIG_FILE", "")
	require.NoError(t, os.Unsetenv("OTEL_CONFIG_FILE"))

	sdk, err := telemetry.NewOpenTelemetrySDK(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, sdk.Shutdown(t.Context())) })

	_, span := sdk.TracerProvider().Tracer("test").Start(t.Context(), "disabled")
	assert.False(t, span.IsRecording())
	span.End()

	assert.Equal(t, sdk.TracerProvider(), otel.GetTracerProvider())
	assert.Equal(t, sdk.MeterProvider(), otel.GetMeterProvider())
	assert.Equal(t, sdk.LoggerProvider(), global.GetLoggerProvider())
	assert.Equal(t, sdk.Propagator(), otel.GetTextMapPropagator())
}

func TestNewOpenTelemetrySDKExportsConfiguredSpans(t *testing.T) {
	requests := make(chan *collectortrace.ExportTraceServiceRequest, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		request := &collectortrace.ExportTraceServiceRequest{}
		require.NoError(t, proto.Unmarshal(body, request))
		requests <- request
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(collector.Close)

	configPath := writeOTelConfig(t, fmt.Sprintf(`
file_format: "1.0"
resource:
  attributes:
    - name: service.name
      value: hexlet-basics-api
    - name: service.version
      value: test-release
propagator:
  composite:
    - tracecontext:
    - baggage:
tracer_provider:
  processors:
    - batch:
        exporter:
          otlp_http:
            endpoint: %s/v1/traces
  sampler:
    always_on:
`, collector.URL))
	t.Setenv("OTEL_CONFIG_FILE", configPath)

	sdk, err := telemetry.NewOpenTelemetrySDK(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.NoError(t, err)

	_, span := sdk.TracerProvider().Tracer("test").Start(t.Context(), "exported")
	require.True(t, span.IsRecording())
	span.End()
	require.NoError(t, sdk.Shutdown(t.Context()))

	select {
	case request := <-requests:
		require.Len(t, request.ResourceSpans, 1)
		attributes := request.ResourceSpans[0].Resource.Attributes
		assert.Contains(t, attributes, attribute("service.name", "hexlet-basics-api"))
		assert.Contains(t, attributes, attribute("service.version", "test-release"))
	case <-time.After(time.Second):
		t.Fatal("OTLP collector did not receive a span batch")
	}
}

func TestNewOpenTelemetrySDKRejectsMissingConfigFile(t *testing.T) {
	t.Setenv("OTEL_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.yaml"))

	_, err := telemetry.NewOpenTelemetrySDK(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.ErrorContains(t, err, "initialize OpenTelemetry SDK")
}

func TestNewOpenTelemetrySDKRejectsInvalidEndpoint(t *testing.T) {
	configPath := writeOTelConfig(t, `
file_format: "1.0"
tracer_provider:
  processors:
    - batch:
        exporter:
          otlp_http:
            endpoint: "://invalid"
`)
	t.Setenv("OTEL_CONFIG_FILE", configPath)

	_, err := telemetry.NewOpenTelemetrySDK(
		t.Context(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	require.ErrorContains(t, err, "initialize OpenTelemetry SDK")
}

func TestNewSentryClientRejectsInvalidDSN(t *testing.T) {
	_, err := telemetry.NewSentryClient(&config.Config{SentryDSN: "not a dsn"})
	require.Error(t, err)
}

func writeOTelConfig(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "otel.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
	return path
}

func attribute(key, value string) *common.KeyValue {
	return &common.KeyValue{
		Key: key,
		Value: &common.AnyValue{
			Value: &common.AnyValue_StringValue{StringValue: value},
		},
	}
}

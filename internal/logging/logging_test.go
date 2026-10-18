package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"hexletbasics/internal/logging"
)

// record logs one line through a logger built by the package and returns it
// decoded: the test reads what an operator's log pipeline would read.
func record(t *testing.T, level slog.Level, log func(context.Context, *slog.Logger)) map[string]any {
	t.Helper()

	var out bytes.Buffer
	log(context.Background(), logging.NewTo(&out, level))

	if out.Len() == 0 {
		return nil
	}

	var line map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &line), out.String())

	return line
}

func TestLogLineCarriesTraceIDOfAnUnsampledSpan(t *testing.T) {
	traceID, err := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("0102030405060708")
	require.NoError(t, err)
	// No FlagsSampled: production samples 1%, and the other 99% must correlate too.
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceID, SpanID: spanID})

	var out bytes.Buffer

	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)
	logging.NewTo(&out, slog.LevelInfo).InfoContext(ctx, "request served")

	var line map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &line))
	require.Equal(t, "0102030405060708090a0b0c0d0e0f10", line["trace_id"])
	require.Equal(t, "0102030405060708", line["span_id"])
}

func TestLogLineWithoutASpanHasNoTraceID(t *testing.T) {
	line := record(t, slog.LevelInfo, func(ctx context.Context, l *slog.Logger) { l.InfoContext(ctx, "starting") })

	require.NotContains(t, line, "trace_id")
	require.NotContains(t, line, "span_id")
}

func TestSecretNamedAttributesAreRedacted(t *testing.T) {
	// A header name keeps its capitals; the variable is out of sloglint's reach on purpose.
	header := "Authorization"
	line := record(t, slog.LevelInfo, func(ctx context.Context, l *slog.Logger) {
		l.InfoContext(ctx, "signed in",
			"user_id", 7,
			"password", "hunter2",
			"access_token", "abc",
			header, "Bearer abc",
			"cookie", "session=1",
			"email", "learner@example.com",
			slog.Group("request", slog.String("api_key", "k"), slog.String("path", "/")))
	})

	require.EqualValues(t, 7, line["user_id"])

	for _, key := range []string{"password", "access_token", "Authorization", "cookie", "email"} {
		require.Equal(t, "[redacted]", line[key], key)
	}

	request, ok := line["request"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "[redacted]", request["api_key"])
	require.Equal(t, "/", request["path"])
}

func TestLevelFiltersLowerLines(t *testing.T) {
	log := func(ctx context.Context, l *slog.Logger) { l.InfoContext(ctx, "chatty") }

	require.Nil(t, record(t, slog.LevelWarn, log))
	require.Equal(t, "chatty", record(t, slog.LevelInfo, log)["msg"])
}

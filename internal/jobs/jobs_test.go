package jobs_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"slices"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport/testdb"
)

// TestQueueBackbone proves the river backbone and its telemetry end to end
// against the test DB: an insert-only client enqueues a ping job, the worker
// runs it, W3C trace context crosses the queue, and both clients emit signals.
// It uses a real pool (not the harness's rolled-back transaction), so it cleans up the
// river_job rows it creates.
func TestQueueBackbone(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM river_job WHERE kind = 'ping'")
		_ = db.Close()
	})

	// nil loader: this smoke test only exercises the ping worker, so the
	// exercise-build worker (and its db/blob deps) is intentionally not registered.
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)

	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	metricReader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(metricReader))
	t.Cleanup(func() {
		require.NoError(t, tracerProvider.Shutdown(context.Background()))
		require.NoError(t, meterProvider.Shutdown(context.Background()))
	})

	workerClient, err := jobs.NewWorkerClient(
		db,
		nil,
		nil,
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		jobs.NewErrorHandler(sentryClient),
		tracerProvider,
		meterProvider,
	)
	require.NoError(t, err)
	insertClient, err := jobs.NewInsertOnlyClient(
		db,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		tracerProvider,
		meterProvider,
	)
	require.NoError(t, err)

	// Subscribe before Start so the completion event can't be missed.
	events, cancel := workerClient.Subscribe(river.EventKindJobCompleted)
	defer cancel()

	require.NoError(t, workerClient.Start(ctx))
	t.Cleanup(func() {
		stopCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = workerClient.Stop(stopCtx)
	})

	_, err = insertClient.Insert(ctx, jobs.PingArgs{}, nil)
	require.NoError(t, err)

	select {
	case ev := <-events:
		require.Equal(t, rivertype.JobStateCompleted, ev.Job.State)
		require.Equal(t, "ping", ev.Job.Kind)

		var metadata map[string]any
		require.NoError(t, json.Unmarshal(ev.Job.Metadata, &metadata))
		require.NotEmpty(t, metadata["traceparent"])
	case <-time.After(10 * time.Second):
		t.Fatal("ping job did not complete within 10s")
	}

	var insertSpan, workSpan sdktrace.ReadOnlySpan
	for _, span := range spanRecorder.Ended() {
		switch span.Name() {
		case "river.insert_many":
			insertSpan = span
		case "river.work":
			workSpan = span
		}
	}
	require.NotNil(t, insertSpan)
	require.NotNil(t, workSpan)
	require.Len(t, workSpan.Links(), 1)
	workLink := workSpan.Links()[0].SpanContext
	require.Equal(t, insertSpan.SpanContext().TraceID(), workLink.TraceID())
	require.Equal(t, insertSpan.SpanContext().SpanID(), workLink.SpanID())
	require.True(t, workLink.IsRemote())

	var metrics metricdata.ResourceMetrics
	require.NoError(t, metricReader.Collect(ctx, &metrics))
	var metricNames []string
	for _, scopeMetrics := range metrics.ScopeMetrics {
		for _, metric := range scopeMetrics.Metrics {
			metricNames = append(metricNames, metric.Name)
		}
	}
	require.True(t, slices.Contains(metricNames, "river.insert_count"))
	require.True(t, slices.Contains(metricNames, "river.work_count"))
}

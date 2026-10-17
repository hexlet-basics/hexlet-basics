package jobs_test

import (
	"context"
	"encoding/json"
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

	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport/testdb"
)

// TestQueueBackbone proves the river backbone and its telemetry end to end
// against the test DB: an insert-only client enqueues a job, the worker
// runs it, W3C trace context crosses the queue, and both clients emit signals.
// It uses a real pool (not the harness's rolled-back transaction), so it cleans up the
// river_job rows it creates.
func TestQueueBackbone(t *testing.T) {
	ctx := context.Background()

	db := testdb.Open(t)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM river_job WHERE kind = 'account_email'")
	})

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
		slog.New(slog.DiscardHandler),
		jobs.NewErrorHandler(sentryClient),
		tracerProvider,
		meterProvider,
		jobs.AccountEmail(emailSenderStub{}),
	)
	require.NoError(t, err)
	insertClient, err := jobs.NewInsertOnlyClient(
		db,
		slog.New(slog.DiscardHandler),
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

	_, err = insertClient.Insert(ctx, jobs.AccountEmailArgs{Purpose: "magic_link", UserID: 1}, nil)
	require.NoError(t, err)

	select {
	case ev := <-events:
		require.Equal(t, rivertype.JobStateCompleted, ev.Job.State)
		require.Equal(t, "account_email", ev.Job.Kind)

		var metadata map[string]any
		require.NoError(t, json.Unmarshal(ev.Job.Metadata, &metadata))
		require.NotEmpty(t, metadata["traceparent"])
	case <-time.After(10 * time.Second):
		t.Fatal("account_email job did not complete within 10s")
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

// A worker without LLM credentials registers no review worker. The job must
// stay available for a configured worker instead of being fetched and failed
// with UnknownJobKindError, which would burn its attempts.
func TestWorkerLeavesKindsItCannotRunQueued(t *testing.T) {
	ctx := context.Background()

	db := testdb.Open(t)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM river_job WHERE kind IN ('review_lesson', 'account_email')")
	})

	logger := slog.New(slog.DiscardHandler)
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	workerClient, err := jobs.NewWorkerClient(
		db, logger, jobs.NewErrorHandler(sentryClient),
		sdktrace.NewTracerProvider(), sdkmetric.NewMeterProvider(),
		jobs.AccountEmail(emailSenderStub{}),
	)
	require.NoError(t, err)
	insertClient, err := jobs.NewInsertOnlyClient(db, logger, sdktrace.NewTracerProvider(), sdkmetric.NewMeterProvider())
	require.NoError(t, err)

	completed, cancel := workerClient.Subscribe(river.EventKindJobCompleted)
	defer cancel()
	require.NoError(t, workerClient.Start(ctx))
	t.Cleanup(func() {
		stopCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = workerClient.Stop(stopCtx)
	})

	review, err := insertClient.Insert(ctx, jobs.ReviewLessonArgs{LessonInfoID: 1}, nil)
	require.NoError(t, err)
	// The account_email job completing proves the worker has fetched after the review
	// was inserted, so the review's survival is not a race with the poll.
	_, err = insertClient.Insert(ctx, jobs.AccountEmailArgs{Purpose: "magic_link", UserID: 1}, nil)
	require.NoError(t, err)
	select {
	case <-completed:
	case <-time.After(10 * time.Second):
		t.Fatal("account_email job did not complete within 10s")
	}

	got, err := insertClient.JobGet(ctx, review.Job.ID)
	require.NoError(t, err)
	require.Equal(t, rivertype.JobStateAvailable, got.State)
	require.Zero(t, got.Attempt)
	require.Equal(t, 3, got.MaxAttempts)
}

// Each job kind lands on the queue that isolates its class of work, and the
// worker runs every one of those queues.
func TestJobKindsLandOnTheirQueues(t *testing.T) {
	ctx := context.Background()

	db := testdb.Open(t)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(),
			"DELETE FROM river_job WHERE kind IN ('account_email', 'exercise_loader', 'review_lesson', 'suggest_related_courses', 'reap_stuck_version_builds')")
	})
	client, err := jobs.NewInsertOnlyClient(db, slog.New(slog.DiscardHandler),
		sdktrace.NewTracerProvider(), sdkmetric.NewMeterProvider())
	require.NoError(t, err)

	for args, queue := range map[river.JobArgs]string{
		jobs.AccountEmailArgs{Purpose: "magic_link", UserID: 1}: jobs.QueueMail,
		jobs.ExerciseLoaderArgs{VersionID: 1}:                   jobs.QueueBuilds,
		jobs.ReapStuckVersionBuildsArgs{}:                       jobs.QueueBuilds,
		jobs.ReviewLessonArgs{LessonInfoID: 1}:                  jobs.QueueLLM,
		jobs.SuggestRelatedCoursesArgs{BlogPostID: 1}:           jobs.QueueLLM,
	} {
		res, err := client.Insert(ctx, args, nil)
		require.NoError(t, err)
		require.Equal(t, queue, res.Job.Queue, args.Kind())
	}
}

type emailSenderStub struct{}

func (emailSenderStub) Send(context.Context, emailtokens.Purpose, int, string) error { return nil }

package jobs_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"hexletbasics/ent"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport/testdb"
)

// runJob runs args through a real worker backed by module and returns the first
// outcome River reports for it: the completion, or the failed attempt. Jobs go
// through the real queue, so what is asserted is what River recorded about the
// job, which is what an operator would see.
func runJob(t *testing.T, kind string, module jobs.Module, args river.JobArgs) *river.Event {
	t.Helper()
	ctx := context.Background()

	db, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM river_job WHERE kind = $1", kind)
		_ = db.Close()
	})

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	worker, err := jobs.NewWorkerClient(
		db, logger, jobs.NewErrorHandler(sentryClient),
		sdktrace.NewTracerProvider(), sdkmetric.NewMeterProvider(), module,
	)
	require.NoError(t, err)
	inserter, err := jobs.NewInsertOnlyClient(db, logger, sdktrace.NewTracerProvider(), sdkmetric.NewMeterProvider())
	require.NoError(t, err)

	events, cancel := worker.Subscribe(river.EventKindJobCompleted, river.EventKindJobFailed)
	defer cancel()
	require.NoError(t, worker.Start(ctx))
	t.Cleanup(func() {
		stopCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
		defer c()
		_ = worker.Stop(stopCtx)
	})

	_, err = inserter.Insert(ctx, args, nil)
	require.NoError(t, err)

	select {
	case event := <-events:
		return event
	case <-time.After(15 * time.Second):
		t.Fatalf("%s job was not worked within 15s", kind)
		return nil
	}
}

func lastError(t *testing.T, job *rivertype.JobRow) string {
	t.Helper()
	require.NotEmpty(t, job.Errors)
	return job.Errors[len(job.Errors)-1].Error
}

type reviewerStub struct{ err error }

func (s reviewerStub) ReviewLesson(context.Context, int) error { return s.err }

type suggesterStub struct{ err error }

func (s suggesterStub) SuggestRelatedCourses(context.Context, int) error { return s.err }

func TestReviewLessonJobCompletesWhenTheReviewSucceeds(t *testing.T) {
	event := runJob(t, "review_lesson", jobs.ReviewLesson(reviewerStub{}), jobs.ReviewLessonArgs{LessonInfoID: 1})

	assert.Equal(t, river.EventKindJobCompleted, event.Kind)
}

func TestReviewLessonJobSurfacesAFailedReviewAndStaysRetryable(t *testing.T) {
	event := runJob(t, "review_lesson",
		jobs.ReviewLesson(reviewerStub{err: errors.New("llm unavailable")}), jobs.ReviewLessonArgs{LessonInfoID: 1})

	assert.Equal(t, river.EventKindJobFailed, event.Kind)
	assert.Contains(t, lastError(t, event.Job), "llm unavailable")
	// Paid LLM calls are retried a bounded number of times, not for weeks.
	assert.Equal(t, 3, event.Job.MaxAttempts)
}

func TestSuggestRelatedCoursesJobCompletesWhenTheSuggestionSucceeds(t *testing.T) {
	event := runJob(t, "suggest_related_courses",
		jobs.SuggestRelatedCourses(suggesterStub{}), jobs.SuggestRelatedCoursesArgs{BlogPostID: 1})

	assert.Equal(t, river.EventKindJobCompleted, event.Kind)
}

func TestSuggestRelatedCoursesJobSurfacesAFailedSuggestionAndStaysRetryable(t *testing.T) {
	event := runJob(t, "suggest_related_courses",
		jobs.SuggestRelatedCourses(suggesterStub{err: errors.New("llm unavailable")}),
		jobs.SuggestRelatedCoursesArgs{BlogPostID: 1})

	assert.Equal(t, river.EventKindJobFailed, event.Kind)
	assert.Contains(t, lastError(t, event.Job), "llm unavailable")
	assert.Equal(t, 3, event.Job.MaxAttempts)
}

// A build whose version is gone fails once and is discarded: a course build is
// re-triggered by an admin or the next webhook, never retried by the queue.
func TestExerciseLoaderJobDiscardsABuildOfAVersionThatDoesNotExist(t *testing.T) {
	db, err := sql.Open("pgx", testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	loader := courseloader.NewLoader(client, nil, nil, nil, nil)

	event := runJob(t, "exercise_loader", jobs.ExerciseLoader(loader), jobs.ExerciseLoaderArgs{VersionID: 987654321})

	assert.Equal(t, river.EventKindJobFailed, event.Kind)
	assert.Contains(t, lastError(t, event.Job), "load course version")
	assert.Equal(t, rivertype.JobStateDiscarded, event.Job.State)
}

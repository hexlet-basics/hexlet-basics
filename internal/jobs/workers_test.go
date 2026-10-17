package jobs_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/testsupport"
)

// work runs one job through worker inside the test's transaction and returns
// what River recorded about it: the event it ended in and the job row as an
// operator would see it, and the error the worker returned. Nothing is
// committed, so nothing is left behind.
func work[A river.JobArgs](t *testing.T, worker river.Worker[A], args A) (*rivertest.WorkResult, error) {
	t.Helper()
	_, tx := testsupport.NewClientWithTx(t)
	return rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{}, worker).
		Work(t.Context(), t, tx, args, nil)
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
	result, err := work(t, jobs.ReviewLessonWorker(reviewerStub{}), jobs.ReviewLessonArgs{LessonInfoID: 1})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
}

func TestReviewLessonJobSurfacesAFailedReviewAndStaysRetryable(t *testing.T) {
	result, err := work(t, jobs.ReviewLessonWorker(reviewerStub{err: errors.New("llm unavailable")}),
		jobs.ReviewLessonArgs{LessonInfoID: 1})

	require.ErrorContains(t, err, "llm unavailable")
	assert.Equal(t, river.EventKindJobFailed, result.EventKind)
	assert.Contains(t, lastError(t, result.Job), "llm unavailable")
	// Paid LLM calls are retried a bounded number of times, not for weeks.
	assert.Equal(t, 3, result.Job.MaxAttempts)
}

func TestSuggestRelatedCoursesJobCompletesWhenTheSuggestionSucceeds(t *testing.T) {
	result, err := work(t, jobs.SuggestRelatedCoursesWorker(suggesterStub{}), jobs.SuggestRelatedCoursesArgs{BlogPostID: 1})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
}

func TestSuggestRelatedCoursesJobSurfacesAFailedSuggestionAndStaysRetryable(t *testing.T) {
	result, err := work(t, jobs.SuggestRelatedCoursesWorker(suggesterStub{err: errors.New("llm unavailable")}),
		jobs.SuggestRelatedCoursesArgs{BlogPostID: 1})

	require.ErrorContains(t, err, "llm unavailable")
	assert.Equal(t, river.EventKindJobFailed, result.EventKind)
	assert.Contains(t, lastError(t, result.Job), "llm unavailable")
	assert.Equal(t, 3, result.Job.MaxAttempts)
}

// A build whose version is gone fails once and is discarded: a course build is
// re-triggered by an admin or the next webhook, never retried by the queue.
func TestExerciseLoaderJobDiscardsABuildOfAVersionThatDoesNotExist(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)
	loader := courseloader.NewLoader(db, nil, nil, nil, nil)
	worker := rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{}, jobs.ExerciseLoaderWorker(loader))

	result, err := worker.Work(t.Context(), t, tx, jobs.ExerciseLoaderArgs{VersionID: 987654321}, nil)

	require.ErrorContains(t, err, "load course version")
	assert.Equal(t, river.EventKindJobFailed, result.EventKind)
	assert.Contains(t, lastError(t, result.Job), "load course version")
	assert.Equal(t, rivertype.JobStateDiscarded, result.Job.State)
}

// The sweep fails the build whose worker died and leaves a running one alone.
// The loader's clock reads testsupport.Today; the fixtures say how long ago each
// version last changed.
func TestReapStuckVersionBuildsJobFailsTheStuckBuildOnly(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)
	loader := courseloader.NewLoader(db, nil, nil, nil, nil,
		courseloader.WithClock(func() time.Time { return testsupport.Today }))
	worker := rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{},
		jobs.ReapStuckVersionBuildsWorker(loader, slog.New(slog.DiscardHandler)))

	result, err := worker.Work(t.Context(), t, tx, jobs.ReapStuckVersionBuildsArgs{}, nil)

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
	stuck := db.CourseVersion.GetX(t.Context(), fixtures.ID(fixtures.VersionLoadableStuck))
	assert.Equal(t, "failed", *stuck.State)
	require.NotNil(t, stuck.Result)
	assert.Contains(t, *stuck.Result, "Build reaped: stuck in 'building' since ")
	running := db.CourseVersion.GetX(t.Context(), fixtures.ID(fixtures.VersionLoadableBuilding))
	assert.Equal(t, "building", *running.State)
	assert.Nil(t, running.Result)
}

// The sweep removes a session that has run out and leaves a live one alone.
func TestReapExpiredSessionsJobRemovesOnlyTheSessionsThatRanOut(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)
	_, err := sessions.New(db).Issue(t.Context(), db.User.GetX(t.Context(), fixtures.ID(fixtures.UserWithPassword)))
	require.NoError(t, err)
	sweep := func(now time.Time) {
		t.Helper()
		worker := rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{},
			jobs.ReapExpiredSessionsWorker(
				sessions.New(db, sessions.WithClock(func() time.Time { return now })),
				slog.New(slog.DiscardHandler),
			))
		result, err := worker.Work(t.Context(), t, tx, jobs.ReapExpiredSessionsArgs{}, nil)
		require.NoError(t, err)
		assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
	}

	sweep(time.Now())
	assert.Equal(t, 1, db.Session.Query().CountX(t.Context()), "a live session stays")

	sweep(time.Now().Add(sessions.MaxLifetime + time.Hour))
	assert.Zero(t, db.Session.Query().CountX(t.Context()))
}

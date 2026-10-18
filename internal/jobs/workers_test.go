package jobs_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/auth"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errLlmUnavailable = errors.New("llm unavailable")
)

// work runs one job through worker inside the test's transaction and returns
// what River recorded about it: the event it ended in and the job row as an
// operator would see it, and the error the worker returned. Nothing is
// committed, so nothing is left behind.
func work[A river.JobArgs](t *testing.T, worker river.Worker[A], args A) (*rivertest.WorkResult, error) {
	t.Helper()
	_, tx := testsupport.NewClientWithTx(t)

	//nolint:wrapcheck // the worker helper passes the job error through to the test
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
	result, err := work(t, jobs.ReviewLessonWorker(reviewerStub{
		err: nil,
	}), jobs.ReviewLessonArgs{LessonInfoID: 1})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
}

func TestReviewLessonJobSurfacesAFailedReviewAndStaysRetryable(t *testing.T) {
	result, err := work(t, jobs.ReviewLessonWorker(reviewerStub{err: errLlmUnavailable}),
		jobs.ReviewLessonArgs{LessonInfoID: 1})

	require.ErrorContains(t, err, "llm unavailable")
	assert.Equal(t, river.EventKindJobFailed, result.EventKind)
	assert.Contains(t, lastError(t, result.Job), "llm unavailable")
	// Paid LLM calls are retried a bounded number of times, not for weeks.
	assert.Equal(t, 3, result.Job.MaxAttempts)
}

func TestSuggestRelatedCoursesJobCompletesWhenTheSuggestionSucceeds(t *testing.T) {
	result, err := work(t, jobs.SuggestRelatedCoursesWorker(suggesterStub{
		err: nil,
	}), jobs.SuggestRelatedCoursesArgs{BlogPostID: 1})

	require.NoError(t, err)
	assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
}

func TestSuggestRelatedCoursesJobSurfacesAFailedSuggestionAndStaysRetryable(t *testing.T) {
	result, err := work(t, jobs.SuggestRelatedCoursesWorker(suggesterStub{err: errLlmUnavailable}),
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
	stuck := db.CourseVersion.GetX(t.Context(), fixture.ID(fixture.VersionLoadableStuck))
	assert.Equal(t, "failed", *stuck.State)
	require.NotNil(t, stuck.Result)
	assert.Contains(t, *stuck.Result, "Build reaped: stuck in 'building' since ")
	running := db.CourseVersion.GetX(t.Context(), fixture.ID(fixture.VersionLoadableBuilding))
	assert.Equal(t, "building", *running.State)
	assert.Nil(t, running.Result)
}

// The sweep removes a session that has run out and leaves a live one alone.
func TestReapExpiredSessionsJobRemovesOnlyTheSessionsThatRanOut(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)
	_, err := auth.New(db).Issue(t.Context(), db.User.GetX(t.Context(), fixture.ID(fixture.UserWithPassword)))
	require.NoError(t, err)

	sweep := func(now time.Time) {
		t.Helper()
		worker := rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{},
			jobs.ReapExpiredSessionsWorker(
				auth.New(db, auth.WithClock(func() time.Time { return now })),
				slog.New(slog.DiscardHandler),
			))
		result, err := worker.Work(t.Context(), t, tx, jobs.ReapExpiredSessionsArgs{}, nil)
		require.NoError(t, err)
		assert.Equal(t, river.EventKindJobCompleted, result.EventKind)
	}

	sweep(time.Now())
	assert.Equal(t, 1, db.Session.Query().CountX(t.Context()), "a live session stays")

	sweep(time.Now().Add(auth.MaxLifetime + time.Hour))
	assert.Zero(t, db.Session.Query().CountX(t.Context()))
}

type senderStub struct{ err error }

func (s senderStub) Send(context.Context, emailtoken.Purpose, int, string) error { return s.err }

// A row that is gone or a change the data refuses will be the same on the next
// try, so the job is cancelled instead of burning its attempts.
func TestJobsCancelWhenTheErrorIsOneRetryingCannotFix(t *testing.T) {
	kinds := map[string]func(err error) (*rivertest.WorkResult, error){
		"review lesson": func(err error) (*rivertest.WorkResult, error) {
			return work(t, jobs.ReviewLessonWorker(reviewerStub{err: err}), jobs.ReviewLessonArgs{LessonInfoID: 1})
		},
		"suggest related courses": func(err error) (*rivertest.WorkResult, error) {
			return work(t, jobs.SuggestRelatedCoursesWorker(suggesterStub{err: err}), jobs.SuggestRelatedCoursesArgs{BlogPostID: 1})
		},
		"account email": func(err error) (*rivertest.WorkResult, error) {
			return work(t, jobs.AccountEmailWorker(senderStub{err: err}), jobs.AccountEmailArgs{
				UserID:  1,
				Purpose: "",
				Locale:  "",
			})
		},
	}
	for kind, run := range kinds {
		for name, cause := range map[string]error{"not found": domain.ErrNotFound, "conflict": domain.ErrConflict} {
			t.Run(kind+" "+name, func(t *testing.T) {
				result, err := run(fmt.Errorf("work: %w", cause))

				require.NoError(t, err, "River settles a cancellation itself")
				assert.Equal(t, river.EventKindJobCancelled, result.EventKind)
				assert.Equal(t, rivertype.JobStateCancelled, result.Job.State)
				assert.Contains(t, lastError(t, result.Job), cause.Error())
			})
		}
	}
}

// A job can run twice (a crash after the write, before River records it). The
// sweep's second run finds nothing left to reap, so the version keeps the one
// result the first run wrote.
func TestReapStuckVersionBuildsJobRunTwiceLeavesOneResult(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)
	loader := courseloader.NewLoader(db, nil, nil, nil, nil,
		courseloader.WithClock(func() time.Time { return testsupport.Today }))
	worker := rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{},
		jobs.ReapStuckVersionBuildsWorker(loader, slog.New(slog.DiscardHandler)))
	stuckID := fixture.ID(fixture.VersionLoadableStuck)

	_, err := worker.Work(t.Context(), t, tx, jobs.ReapStuckVersionBuildsArgs{}, nil)
	require.NoError(t, err)
	first := *db.CourseVersion.GetX(t.Context(), stuckID).Result

	_, err = worker.Work(t.Context(), t, tx, jobs.ReapStuckVersionBuildsArgs{}, nil)
	require.NoError(t, err)

	stuck := db.CourseVersion.GetX(t.Context(), stuckID)
	assert.Equal(t, "failed", *stuck.State)
	assert.Equal(t, first, *stuck.Result)
}

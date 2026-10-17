package progress_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/events"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

func (r *rig) failTheRun() {
	r.mocks.Runner.RunFunc = func(context.Context, progress.Submission) (progress.Outcome, error) {
		return testsupport.FailingOutcome(), nil
	}
}

func TestCheckByAGuestAddsToTheCookieOnlyWhenThePassFinishesSomethingNew(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	before := r.db.LessonProgress.Query().CountX(ctx)

	passed, err := r.progress.CheckSolution(ctx, r.check(t, progress.Learner{}, firstLesson))
	require.NoError(t, err)
	require.NotNil(t, passed.Guest, "a first pass is worth a cookie")
	assert.True(t, passed.Guest.Has(jsSlug, firstLesson))
	assert.False(t, passed.LessonFinished, "a guest finishes nothing a database knows about")
	assert.Zero(t, r.db.LessonProgress.Query().CountX(ctx)-before, "and leaves no rows")

	// The same lesson again finishes nothing new, so there is nothing to write.
	again, err := r.progress.CheckSolution(ctx, r.check(t, progress.Learner{Guest: *passed.Guest}, firstLesson))
	require.NoError(t, err)
	assert.Nil(t, again.Guest)

	// A failing run never moves them.
	r.failTheRun()
	failed, err := r.progress.CheckSolution(ctx, r.check(t, progress.Learner{Guest: *passed.Guest}, secondLesson))
	require.NoError(t, err)
	assert.False(t, failed.Outcome.Passed)
	assert.Nil(t, failed.Guest)

	checked := testsupport.PublishedEvents[events.SolutionChecked](t, r.services)
	assert.Len(t, checked, 3, "every submission is a fact, a guest's included")
	assert.Empty(t, testsupport.PublishedEvents[events.CourseStarted](t, r.services), "but a guest has no enrollment to announce")
}

func TestCheckThatFailsLeavesTheLessonStartedAndNothingFinished(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, "user_blank")
	r.failTheRun()

	result, err := r.progress.CheckSolution(ctx, r.check(t, blankUser, firstLesson))
	require.NoError(t, err)

	assert.False(t, result.Outcome.Passed)
	assert.False(t, result.LessonFinished)
	assert.Nil(t, result.Guest)
	row := r.db.LessonProgress.Query().Where(lessonprogress.UserID(blankUser.UserID)).OnlyX(ctx)
	assert.Equal(t, progress.StateStarted, *row.State)
	assert.Empty(t, testsupport.PublishedEvents[events.LessonFinished](t, r.services))
}

func TestPassingEveryLessonFinishesTheCourseOnceAndRepeatsChangeNothing(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, "user_blank")

	var results []*progress.CheckResult
	for _, slug := range []string{firstLesson, secondLesson, thirdLesson} {
		result, err := r.progress.CheckSolution(ctx, r.check(t, blankUser, slug))
		require.NoError(t, err)
		results = append(results, result)
	}

	assert.Equal(t, []bool{false, false, true}, []bool{results[0].CourseFinished, results[1].CourseFinished, results[2].CourseFinished})
	for _, result := range results {
		assert.True(t, result.LessonFinished)
	}
	state, err := r.progress.CourseState(ctx, blankUser, r.course(t, jsSlug).ID)
	require.NoError(t, err)
	assert.Equal(t, progress.StateFinished, state.State)
	assert.Equal(t, 100, state.Completion)

	course := testsupport.PublishedEvents[events.CourseFinished](t, r.services)
	require.Len(t, course, 1)
	assert.Equal(t, 1, course[0].OccurrenceCount, "the first course this learner has finished")
	assert.Equal(t, 1, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID), enrollment.StateEQ(progress.StateFinished)).CountX(ctx))

	// Experimenting on a passed lesson is safe: nothing transitions, nothing is published.
	published := testsupport.PublishedCount(t, r.services)
	again, err := r.progress.CheckSolution(ctx, r.check(t, blankUser, thirdLesson))
	require.NoError(t, err)
	assert.True(t, again.Outcome.Passed)
	assert.False(t, again.LessonFinished)
	assert.False(t, again.CourseFinished)
	assert.Equal(t, published+1, testsupport.PublishedCount(t, r.services), "only the submission itself is a fact")
}

func TestPassingALessonThatLeavesAGapKeepsTheCourseOpen(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, "user_admin_learner") // finished the first lesson only

	// Close the third lesson behind the learner's back, as inherited data would.
	third := r.lesson(t, thirdLesson)
	r.db.LessonProgress.Create().
		SetUserID(adminLearner.UserID).
		SetCourseID(r.course(t, jsSlug).ID).
		SetEnrollmentID(r.db.Enrollment.Query().Where(enrollment.UserID(adminLearner.UserID)).OnlyX(ctx).ID).
		SetLessonID(third.ID).
		SetState(progress.StateFinished).
		ExecX(ctx)

	result, err := r.progress.CheckSolution(ctx, r.check(t, adminLearner, firstLesson))
	require.NoError(t, err)

	assert.False(t, result.CourseFinished, "the second lesson is still unfinished")
	assert.Empty(t, testsupport.PublishedEvents[events.CourseFinished](t, r.services))
}

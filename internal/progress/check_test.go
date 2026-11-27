package progress_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/fixture"
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
	enrollmentsBefore := r.db.Enrollment.Query().CountX(ctx)

	passed, err := r.progress.CheckSolution(ctx, r.check(t, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, firstLesson))
	require.NoError(t, err)
	require.NotNil(t, passed.Guest, "a first pass is worth a cookie")
	assert.True(t, passed.Guest.Has(jsSlug, firstLesson))
	assert.False(t, passed.LessonFinished, "a guest finishes nothing a database knows about")
	assert.Zero(t, r.db.LessonProgress.Query().CountX(ctx)-before, "and leaves no rows")

	// The same lesson again finishes nothing new, so there is nothing to write.
	again, err := r.progress.CheckSolution(ctx, r.check(t, progress.Learner{
		Guest:  *passed.Guest,
		UserID: 0,
	}, firstLesson))
	require.NoError(t, err)
	assert.Nil(t, again.Guest)

	// A failing run never moves them.
	r.failTheRun()
	failed, err := r.progress.CheckSolution(ctx, r.check(t, progress.Learner{
		Guest:  *passed.Guest,
		UserID: 0,
	}, secondLesson))
	require.NoError(t, err)
	assert.False(t, failed.Outcome.Passed)
	assert.Nil(t, failed.Guest)
	assert.Zero(t, r.db.Enrollment.Query().CountX(ctx)-enrollmentsBefore, "a guest has no enrollment")
}

func TestCheckThatFailsLeavesTheLessonStartedAndNothingFinished(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, fixture.UserBlank)
	r.failTheRun()

	result, err := r.progress.CheckSolution(ctx, r.check(t, blankUser, firstLesson))
	require.NoError(t, err)

	assert.False(t, result.Outcome.Passed)
	assert.False(t, result.LessonFinished)
	assert.Nil(t, result.Guest)

	row := r.db.LessonProgress.Query().Where(lessonprogress.UserID(blankUser.UserID)).OnlyX(ctx)
	assert.Equal(t, progress.StateStarted, *row.State)
}

func TestPassingEveryLessonFinishesTheCourseOnceAndRepeatsChangeNothing(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, fixture.UserBlank)

	slugs := []string{firstLesson, secondLesson, thirdLesson}

	results := make([]*progress.CheckResult, 0, len(slugs))
	for _, slug := range slugs {
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

	assert.Equal(t, 1, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID), enrollment.StateEQ(progress.StateFinished)).CountX(ctx))

	// Experimenting on a passed lesson is safe: nothing transitions.
	again, err := r.progress.CheckSolution(ctx, r.check(t, blankUser, thirdLesson))
	require.NoError(t, err)
	assert.True(t, again.Outcome.Passed)
	assert.False(t, again.LessonFinished)
	assert.False(t, again.CourseFinished)
}

func TestPassingALessonThatLeavesAGapKeepsTheCourseOpen(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	// Inherited data: the first and the third lessons are finished, the second is not.
	gapLearner := r.learner(t, fixture.UserLearnerWithAGap)

	result, err := r.progress.CheckSolution(ctx, r.check(t, gapLearner, firstLesson))
	require.NoError(t, err)

	assert.False(t, result.CourseFinished, "the second lesson is still unfinished")
	assert.Zero(t, r.db.Enrollment.Query().
		Where(enrollment.UserID(gapLearner.UserID), enrollment.StateEQ(progress.StateFinished)).
		CountX(ctx), "the enrollment stays open")
}

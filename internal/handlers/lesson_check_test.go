package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// A passing submission finishes the lesson. The learner never pressed start —
// the Lesson Progress is created lazily, so arriving by a direct link does not
// lose the work.
func TestCheckFinishesALesson(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.True(t, result.Response.Passed)
	assert.True(t, result.Response.LessonHasBeenFinished)
	assert.False(t, result.Response.CourseHasBeenFinished, "the third lesson is still unfinished")
	assert.Empty(t, result.SetCookie, "a signed-in learner's progress is rows, not a cookie")
	assert.Empty(t, h.ResponseCookies(), "and no empty header is written for the one they do not get")

	row, err := h.DB.LessonProgress.Query().
		Where(
			lessonprogress.LessonID(lesson.ID),
			lessonprogress.UserID(actingUserID(t, h)),
		).
		Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, "finished", *row.State)

	checked := publishedOf[events.SolutionChecked](t, h)
	require.Len(t, checked, 1)
	assert.True(t, checked[0].Passed)
	assert.Equal(t, secondLessonSlug, checked[0].LessonSlug)
	assert.Equal(t, jsCourseSlug, checked[0].CourseSlug)

	started := publishedLessonStarted(t, h)
	assert.Equal(t, secondLessonSlug, started.LessonSlug, "submitting starts the lesson it records")

	finished := publishedOf[events.LessonFinished](t, h)
	require.Len(t, finished, 1)
	assert.Equal(t, 2, finished[0].OccurrenceCount, "counts the Lesson Progress rows on this enrollment")
	assert.Empty(t, publishedOf[events.CourseFinished](t, h))

}

// A failed check is still work done on the lesson: it is recorded, and the
// lesson stays started rather than being left untouched.
func TestCheckLeavesAFailedLessonStarted(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	h.Runner.RunFunc = func(context.Context, progress.Submission) (progress.Outcome, error) {
		return testsupport.FailingOutcome(), nil
	}
	lesson := lessonBySlug(t, h, secondLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.False(t, result.Response.Passed)
	assert.Equal(t, api.LessonCheckingResponseResultFailed, result.Response.Result)
	assert.False(t, result.Response.LessonHasBeenFinished)

	row := h.DB.LessonProgress.Query().
		Where(
			lessonprogress.LessonID(lesson.ID),
			lessonprogress.UserID(actingUserID(t, h)),
		).
		OnlyX(ctx)
	assert.Equal(t, "started", *row.State)

	checked := publishedOf[events.SolutionChecked](t, h)
	require.Len(t, checked, 1)
	assert.False(t, checked[0].Passed, "the funnel sees the attempt as well as the pass")
	assert.Empty(t, publishedOf[events.LessonFinished](t, h))
}

// Lessons may be taken in any order (ADR-0016): the third before the first
// finishes that lesson and only that lesson.
func TestCheckFinishesALessonTakenOutOfOrder(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, thirdLessonSlug)
	second := lessonBySlug(t, h, secondLessonSlug)
	before := h.DB.LessonProgress.Query().CountX(ctx)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.Passed)
	assert.True(t, result.Response.LessonHasBeenFinished)
	assert.False(t, result.Response.CourseHasBeenFinished)

	assert.Equal(t, before+1, h.DB.LessonProgress.Query().CountX(ctx), "one row, for this lesson")
	assert.Equal(t, 0, h.DB.LessonProgress.Query().Where(lessonprogress.LessonID(second.ID),
		lessonprogress.UserID(actingUserID(t, h)), lessonprogress.StateEQ("finished")).CountX(ctx),
		"the lesson before it is not implied")
}

// The last unfinished lesson of a course finishes the course as well, and the
// learner is enrolled by the submission itself.
func TestCheckFinishesTheCourseOnItsLastLesson(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonLabeled(t, h, fixture.LessonSingleLessonPrinting)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.LessonHasBeenFinished)
	assert.True(t, result.Response.CourseHasBeenFinished)

	enrolled := h.DB.Enrollment.Query().
		Where(
			enrollment.UserID(actingUserID(t, h)),
			enrollment.CourseID(*lesson.CourseID),
		).
		OnlyX(ctx)
	assert.Equal(t, "finished", *enrolled.State)

	require.Len(t, publishedOf[events.CourseStarted](t, h), 1, "the submission enrolled them")
	courseFinished := publishedOf[events.CourseFinished](t, h)
	require.Len(t, courseFinished, 1)
	assert.Equal(t, "single-lesson", courseFinished[0].Slug)
	assert.Equal(t, 1, courseFinished[0].OccurrenceCount,
		"counts the learner's finished enrollments including the one just finished")
}

// Re-submitting to a lesson already passed is safe: experimenting must not
// double-count anything the CRM reads.
func TestCheckIsSafeToRepeatOnAFinishedLesson(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	before := h.DB.LessonProgress.Query().CountX(ctx)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.Passed)
	assert.False(t, result.Response.LessonHasBeenFinished, "it was already finished")
	assert.False(t, result.Response.CourseHasBeenFinished)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx))
	row := h.DB.LessonProgress.Query().
		Where(
			lessonprogress.LessonID(lesson.ID),
			lessonprogress.UserID(actingUserID(t, h)),
		).
		OnlyX(ctx)
	assert.Equal(t, "finished", *row.State)

	assert.Len(t, publishedOf[events.SolutionChecked](t, h), 1, "the check itself is still a fact")
	assert.Empty(t, publishedOf[events.LessonFinished](t, h))
	assert.Empty(t, publishedOf[events.CourseFinished](t, h))
}

// A learner whose lessons are all finished while their enrollment is not —
// inherited data, or a version that dropped what they had left — closes the
// course by passing anything, rather than having no submission left that would
// notice.
func TestCheckClosesACourseAlreadyEffectivelyFinished(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerAllLessonsDone)
	ctx := t.Context()

	first := lessonBySlug(t, h, firstLessonSlug)
	result := submitSolution(t, h, first)
	assert.False(t, result.Response.LessonHasBeenFinished, "it was already finished")
	assert.True(t, result.Response.CourseHasBeenFinished)

	enrolled := h.DB.Enrollment.Query().
		Where(
			enrollment.UserID(actingUserID(t, h)),
			enrollment.CourseID(*first.CourseID),
		).
		OnlyX(ctx)
	assert.Equal(t, "finished", *enrolled.State)

	assert.Empty(t, publishedOf[events.LessonFinished](t, h), "no lesson transitioned")
	require.Len(t, publishedOf[events.CourseFinished](t, h), 1)
}

// A gap in the learner's history — inherited from a system that allowed any
// order — does not block them, but keeps the course open until they close it.
// The learner has finished the first and the third of four lessons.
func TestCheckKeepsACourseWithAGapUnfinished(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerFourLessonsGap)
	ctx := t.Context()
	crs := h.DB.Course.GetX(ctx, fixture.ID(fixture.CourseWithFourLessons))
	owner := h.DB.Enrollment.Query().
		Where(enrollment.UserID(actingUserID(t, h)), enrollment.CourseID(crs.ID)).
		OnlyX(ctx)

	// The gap does not hold them back.
	result := submitSolution(t, h, lessonLabeled(t, h, fixture.LessonWithFourLessonsFour))
	assert.True(t, result.Response.LessonHasBeenFinished)
	assert.False(t, result.Response.CourseHasBeenFinished, "the second lesson is still open")

	owner = h.DB.Enrollment.Query().Where(enrollment.ID(owner.ID)).OnlyX(ctx)
	assert.Equal(t, "started", *owner.State)

	// Closing the gap finishes the course.
	closing := submitSolution(t, h, lessonLabeled(t, h, fixture.LessonWithFourLessonsTwo))
	assert.True(t, closing.Response.CourseHasBeenFinished)

	owner = h.DB.Enrollment.Query().Where(enrollment.ID(owner.ID)).OnlyX(ctx)
	assert.Equal(t, "finished", *owner.State)
}

// A lesson version belonging to a superseded build is refused: the tests the
// code was written against are no longer the ones that would run.
func TestCheckRefusesASupersededLessonVersion(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)

	_, err := h.Client.CheckLesson(ctx,
		&api.CheckLessonInput{Code: submittedCode, VersionId: fixture.ID32(fixture.LessonVersionJavascriptHelloWorldRetired)},
		api.CheckLessonParams{ID: int32(lesson.ID)},
	)
	requireProblem(t, err, http.StatusNotFound)
}

// ---- Guests ---------------------------------------------------------------

// A visitor with no cookie may take any lesson of a course. Their pass comes
// back as a signed cookie and leaves no rows behind.
func TestGuestCheckAddsToTheCookieAndWritesNoRows(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	enrollments := h.DB.Enrollment.Query().CountX(ctx)
	rows := h.DB.LessonProgress.Query().CountX(ctx)

	result := submitSolution(t, h, lesson)
	assert.False(t, result.Response.LessonHasBeenFinished,
		"a guest finishes nothing that a record could hold")

	guest := testsupport.DecodeGuestCookie(t, h.ResponseCookies())
	assert.True(t, guest.Has(jsCourseSlug, firstLessonSlug))

	assert.Equal(t, enrollments, h.DB.Enrollment.Query().CountX(ctx))
	assert.Equal(t, rows, h.DB.LessonProgress.Query().CountX(ctx))

	assert.Len(t, publishedOf[events.SolutionChecked](t, h), 1, "pre-signup activity stays visible")
	assert.Empty(t, publishedOf[events.CourseStarted](t, h))
	assert.Empty(t, publishedOf[events.LessonStarted](t, h))
	assert.Empty(t, publishedOf[events.LessonFinished](t, h))
}

// A guest may take a lesson out of order too: the cookie names that lesson and
// not the ones before it.
func TestGuestCheckFinishesALessonOutOfOrder(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	lesson := lessonBySlug(t, h, thirdLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.Passed)

	guest := testsupport.DecodeGuestCookie(t, h.ResponseCookies())
	assert.True(t, guest.Has(jsCourseSlug, thirdLessonSlug))
	assert.False(t, guest.Has(jsCourseSlug, firstLessonSlug), "the lessons before it are not implied")
}

// Replaying the cookie a previous pass returned and passing another lesson adds
// it to what the cookie carries.
func TestGuestCheckAddsToTheCarriedCookie(t *testing.T) {
	t.Parallel()
	h := testsupport.NewVisitorHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, firstLessonSlug))
	lesson := lessonBySlug(t, h, secondLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.Passed)

	guest := testsupport.DecodeGuestCookie(t, h.ResponseCookies())
	assert.True(t, guest.Has(jsCourseSlug, firstLessonSlug), "what it carried is kept")
	assert.True(t, guest.Has(jsCourseSlug, secondLessonSlug), "and the new pass is added")
}

// A failed guest check records the fact and nothing else: the cookie changes
// only on a pass.
func TestGuestCheckDoesNotAdvanceOnAFailure(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	h.Runner.RunFunc = func(context.Context, progress.Submission) (progress.Outcome, error) {
		return testsupport.FailingOutcome(), nil
	}
	lesson := lessonBySlug(t, h, firstLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.False(t, result.Response.Passed)
	assert.Empty(t, h.ResponseCookies(), "nothing to store")
	assert.Len(t, publishedOf[events.SolutionChecked](t, h), 1)
}

// A cookie naming a lesson the current version no longer contains resolves to
// nothing, so that course starts over for the visitor: the next pass is the only
// lesson the cookie counts.
func TestGuestCheckIgnoresALessonTheCurrentVersionDropped(t *testing.T) {
	t.Parallel()
	h := testsupport.NewVisitorHarness(t, progressCookie(jsCourseSlug, "removed-in-a-later-build"))

	first := lessonBySlug(t, h, firstLessonSlug)
	result := submitSolution(t, h, first)
	assert.True(t, result.Response.Passed, "the visitor resumes from the beginning")

	guest := testsupport.DecodeGuestCookie(t, h.ResponseCookies())
	assert.True(t, guest.Has(jsCourseSlug, firstLessonSlug))
}

// ---- Helpers --------------------------------------------------------------

const submittedCode = "console.log('hello, world!')"

// checkLesson submits against the lesson's version in the course's current
// build, which is the only version a client can legitimately have loaded.
func checkLesson(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) (*api.LessonCheckingResponseHeaders, error) {
	t.Helper()
	return h.Client.CheckLesson(t.Context(),
		&api.CheckLessonInput{
			Code:      submittedCode,
			VersionId: int32(currentLessonVersionID(t, h, lesson)),
		},
		api.CheckLessonParams{ID: int32(lesson.ID)},
	)
}

// submitSolution is checkLesson for the tests that expect it to succeed.
func submitSolution(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) *api.LessonCheckingResponseHeaders {
	t.Helper()
	res, err := checkLesson(t, h, lesson)
	require.NoError(t, err)
	result := res
	return result
}

func currentLessonVersionID(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) int {
	t.Helper()
	ctx := t.Context()
	crs := h.DB.Course.Query().Where(course.ID(*lesson.CourseID)).OnlyX(ctx)
	require.NotNil(t, crs.CurrentVersionID)
	return h.DB.CourseLessonVersion.Query().
		Where(
			courselessonversion.LessonID(lesson.ID),
			courselessonversion.CourseVersionID(*crs.CurrentVersionID),
		).
		OnlyX(ctx).ID
}

// A solution that never terminates is its own answer, not a plain failure: the
// learner is told to look for a loop, and the lesson stays started.
func TestCheckReportsASolutionThatNeverTerminatedAsItsOwnResult(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.Runner.RunFunc = func(context.Context, progress.Submission) (progress.Outcome, error) {
		return progress.Outcome{Result: progress.ResultFailedInfinity, Status: 124}, nil
	}
	lesson := lessonBySlug(t, h, secondLessonSlug)

	result := submitSolution(t, h, lesson)

	assert.False(t, result.Response.Passed)
	assert.Equal(t, api.LessonCheckingResponseResultFailedInfinity, result.Response.Result)
	assert.False(t, result.Response.LessonHasBeenFinished)
	row := h.DB.LessonProgress.Query().
		Where(lessonprogress.LessonID(lesson.ID), lessonprogress.UserID(actingUserID(t, h))).
		OnlyX(t.Context())
	assert.Equal(t, "started", *row.State)
}

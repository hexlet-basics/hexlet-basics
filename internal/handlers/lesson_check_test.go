package handlers_test

import (
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
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// A passing submission for the lesson one past the learner's furthest finished
// one finishes it. The learner never pressed start — the Lesson Progress is
// created lazily, so arriving by a direct link does not lose the work.
func TestCheckFinishesAnAvailableLesson(t *testing.T) {
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

	checked := publishedOf[events.SolutionChecked](h)
	require.Len(t, checked, 1)
	assert.True(t, checked[0].Passed)
	assert.Equal(t, secondLessonSlug, checked[0].LessonSlug)
	assert.Equal(t, jsCourseSlug, checked[0].CourseSlug)

	started := publishedLessonStarted(t, h)
	assert.Equal(t, secondLessonSlug, started.LessonSlug, "submitting starts the lesson it records")

	finished := publishedOf[events.LessonFinished](h)
	require.Len(t, finished, 1)
	assert.Equal(t, 2, finished[0].OccurrenceCount, "counts the Lesson Progress rows on this enrollment")
	assert.Empty(t, publishedOf[events.CourseFinished](h))

	require.Len(t, h.Runner.Submissions, 1)
	assert.Equal(t, actingUserID(t, h), h.Runner.Submissions[0].UserID)
}

// A failed check is still work done on the lesson: it is recorded, and the
// lesson stays started rather than being left untouched.
func TestCheckLeavesAFailedLessonStarted(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	h.Runner.Outcome = testsupport.FailingOutcome()
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

	checked := publishedOf[events.SolutionChecked](h)
	require.Len(t, checked, 1)
	assert.False(t, checked[0].Passed, "the funnel sees the attempt as well as the pass")
	assert.Empty(t, publishedOf[events.LessonFinished](h))
}

// Two past the furthest finished lesson is refused before anything runs: no
// container is started and nothing is written.
func TestCheckRefusesALessonBeyondTheGate(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, thirdLessonSlug)
	before := h.DB.LessonProgress.Query().CountX(ctx)

	res, err := checkLesson(t, h, lesson)
	require.NoError(t, err)
	require.IsType(t, &api.ProblemDetails{}, res)
	assert.Equal(t, http.StatusConflict, h.LastStatus())

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx))
	assert.Empty(t, h.Events.Published)
	assert.Empty(t, h.Runner.Submissions, "a refused submission is never run")
}

// The last unfinished lesson of a course finishes the course as well, and the
// learner is enrolled by the submission itself.
func TestCheckFinishesTheCourseOnItsLastLesson(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := seedCourse(t, h, "go-basics", "printing")

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

	require.Len(t, publishedOf[events.CourseStarted](h), 1, "the submission enrolled them")
	courseFinished := publishedOf[events.CourseFinished](h)
	require.Len(t, courseFinished, 1)
	assert.Equal(t, "go-basics", courseFinished[0].Slug)
	assert.Equal(t, 1, courseFinished[0].OccurrenceCount,
		"counts the learner's finished enrollments including the one just finished")
}

// Re-submitting to a lesson already passed is safe: experimenting must not
// double-count anything the CRM reads.
func TestCheckIsSafeToRepeatOnAFinishedLesson(t *testing.T) {
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

	assert.Len(t, publishedOf[events.SolutionChecked](h), 1, "the check itself is still a fact")
	assert.Empty(t, publishedOf[events.LessonFinished](h))
	assert.Empty(t, publishedOf[events.CourseFinished](h))
}

// A learner whose lessons are all finished while their enrollment is not —
// inherited data, or a version that dropped what they had left — closes the
// course by passing anything, rather than having no submission left that would
// notice.
func TestCheckClosesACourseAlreadyEffectivelyFinished(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	finishLesson(t, h, secondLessonSlug)
	finishLesson(t, h, thirdLessonSlug)

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

	assert.Empty(t, publishedOf[events.LessonFinished](h), "no lesson transitioned")
	require.Len(t, publishedOf[events.CourseFinished](h), 1)
}

// A gap in the learner's history — inherited from a system that allowed any
// order — does not block them, but keeps the course open until they close it.
func TestCheckKeepsACourseWithAGapUnfinished(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	user := actingUserID(t, h)
	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(ctx)
	owner := h.DB.Enrollment.Query().
		Where(enrollment.UserID(user), enrollment.CourseID(js.ID)).
		OnlyX(ctx)

	// Finished out of order: the third lesson is done, the second is not.
	third := lessonBySlug(t, h, thirdLessonSlug)
	h.DB.LessonProgress.Create().
		SetUserID(user).
		SetCourseID(js.ID).
		SetEnrollmentID(owner.ID).
		SetLessonID(third.ID).
		SetState("finished").
		SaveX(ctx)
	fourth := seedLesson(t, h, js, "functions", 4)

	// The gap does not hold them back: position 4 is one past their furthest.
	result := submitSolution(t, h, fourth)
	assert.True(t, result.Response.LessonHasBeenFinished)
	assert.False(t, result.Response.CourseHasBeenFinished, "the second lesson is still open")

	owner = h.DB.Enrollment.Query().Where(enrollment.ID(owner.ID)).OnlyX(ctx)
	assert.Equal(t, "started", *owner.State)

	// Closing the gap finishes the course.
	second := lessonBySlug(t, h, secondLessonSlug)
	closing := submitSolution(t, h, second)
	assert.True(t, closing.Response.CourseHasBeenFinished)

	owner = h.DB.Enrollment.Query().Where(enrollment.ID(owner.ID)).OnlyX(ctx)
	assert.Equal(t, "finished", *owner.State)
}

// A lesson version belonging to a superseded build is refused: the tests the
// code was written against are no longer the ones that would run.
func TestCheckRefusesASupersededLessonVersion(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)
	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(ctx)

	retired := h.DB.CourseLessonVersion.Create().
		SetCourseID(js.ID).
		SetCourseVersionID(retiredVersionID(t, h, js)).
		SetLessonID(lesson.ID).
		SetModuleVersionID(seedModuleVersion(t, h, js.ID, retiredVersionID(t, h, js))).
		SetNaturalOrder(2).
		SaveX(ctx)

	res, err := h.Client.CheckLesson(ctx,
		&api.CheckLessonInput{Code: submittedCode, VersionId: int32(retired.ID)},
		api.CheckLessonParams{ID: int32(lesson.ID)},
	)
	require.NoError(t, err)
	require.IsType(t, &api.NotFoundError{}, res)
	assert.Empty(t, h.Runner.Submissions)
}

// ---- Guests ---------------------------------------------------------------

// A visitor with no cookie may take the first lesson of a course. Their pass
// comes back as a signed cookie and leaves no rows behind.
func TestGuestCheckAdvancesTheCookieAndWritesNoRows(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	ctx := t.Context()
	lesson := lessonBySlug(t, h, firstLessonSlug)
	enrollments := h.DB.Enrollment.Query().CountX(ctx)
	rows := h.DB.LessonProgress.Query().CountX(ctx)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.LessonHasBeenFinished == false,
		"a guest finishes nothing that a record could hold")

	guest := testsupport.DecodeGuestCookie(t, h.ResponseCookies())
	furthest, ok := guest.Furthest(jsCourseSlug)
	require.True(t, ok)
	assert.Equal(t, firstLessonSlug, furthest)

	assert.Equal(t, enrollments, h.DB.Enrollment.Query().CountX(ctx))
	assert.Equal(t, rows, h.DB.LessonProgress.Query().CountX(ctx))

	assert.Len(t, publishedOf[events.SolutionChecked](h), 1, "pre-signup activity stays visible")
	assert.Empty(t, publishedOf[events.CourseStarted](h))
	assert.Empty(t, publishedOf[events.LessonStarted](h))
	assert.Empty(t, publishedOf[events.LessonFinished](h))
}

// The same rule gates a guest: the second lesson is out of reach until the
// first one is passed.
func TestGuestCheckRefusesALessonBeyondTheGate(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)

	res, err := checkLesson(t, h, lesson)
	require.NoError(t, err)
	require.IsType(t, &api.ProblemDetails{}, res)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
	assert.Empty(t, h.Runner.Submissions)
}

// Replaying the cookie a previous pass returned opens the next lesson, and
// passing it moves the stored position forward.
func TestGuestCheckAdvancesFromTheCarriedCookie(t *testing.T) {
	h := testsupport.NewVisitorHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, firstLessonSlug))
	lesson := lessonBySlug(t, h, secondLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.True(t, result.Response.Passed)

	guest := testsupport.DecodeGuestCookie(t, h.ResponseCookies())
	furthest, ok := guest.Furthest(jsCourseSlug)
	require.True(t, ok)
	assert.Equal(t, secondLessonSlug, furthest, "the cookie carries the furthest lesson finished")
}

// A failed guest check records the fact and nothing else: the position moves
// only on a pass.
func TestGuestCheckDoesNotAdvanceOnAFailure(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	h.Runner.Outcome = testsupport.FailingOutcome()
	lesson := lessonBySlug(t, h, firstLessonSlug)

	result := submitSolution(t, h, lesson)
	assert.False(t, result.Response.Passed)
	assert.Empty(t, h.ResponseCookies(), "nothing to store")
	assert.Len(t, publishedOf[events.SolutionChecked](h), 1)
}

// A cookie naming a lesson the current version no longer contains resolves to
// no position at all, so that course starts over for the visitor.
func TestGuestCheckResetsAPositionTheCurrentVersionDropped(t *testing.T) {
	h := testsupport.NewVisitorHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, "removed-in-a-later-build"))

	second := lessonBySlug(t, h, secondLessonSlug)
	res, err := checkLesson(t, h, second)
	require.NoError(t, err)
	assert.IsType(t, &api.ProblemDetails{}, res, "the stored lesson raises nothing")

	first := lessonBySlug(t, h, firstLessonSlug)
	result := submitSolution(t, h, first)
	assert.True(t, result.Response.Passed, "the visitor resumes from the beginning")
}

// ---- Helpers --------------------------------------------------------------

const submittedCode = "console.log('hello, world!')"

// checkLesson submits against the lesson's version in the course's current
// build, which is the only version a client can legitimately have loaded.
func checkLesson(t *testing.T, h *testsupport.Harness, lesson *ent.CourseLesson) (api.CheckLessonRes, error) {
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
	result, ok := res.(*api.LessonCheckingResponseHeaders)
	require.True(t, ok, "unexpected check response %T (status %d)", res, h.LastStatus())
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

// seedLesson adds a lesson to the course's current version at a position of the
// caller's choosing, for the histories fixtures cannot express.
func seedLesson(t *testing.T, h *testsupport.Harness, crs *ent.Course, slug string, position int) *ent.CourseLesson {
	t.Helper()
	ctx := t.Context()
	require.NotNil(t, crs.CurrentVersionID)

	lesson := h.DB.CourseLesson.Create().
		SetSlug(slug).
		SetCourseID(crs.ID).
		SetNaturalOrder(position).
		SaveX(ctx)
	h.DB.CourseLessonVersion.Create().
		SetCourseID(crs.ID).
		SetCourseVersionID(*crs.CurrentVersionID).
		SetLessonID(lesson.ID).
		SetModuleVersionID(seedModuleVersion(t, h, crs.ID, *crs.CurrentVersionID)).
		SetNaturalOrder(position).
		SetPathToCode("/exercises-javascript/modules/10-basics/" + slug).
		SaveX(ctx)
	return lesson
}

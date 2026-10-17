package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// The harness acts as the fixture admin (user_admin_learner), who has finished the first
// lesson of the JavaScript course and is enrolled in it. Lessons are addressed
// by slug throughout: fixture ids change whenever a row is added.
const (
	firstLessonSlug  = "hello-world" // position 1, finished by the acting user
	secondLessonSlug = "variables"   // position 2, the first one still to do
	thirdLessonSlug  = "strings"     // position 3
	jsCourseSlug     = "javascript"
	// jsLandingSlug is the public URL of the JavaScript course's Russian landing
	// page, the slug GetCourse resolves (legacy `/languages/javascript-ru`).
	jsLandingSlug = "javascript-ru"
	// jsEnglishLandingSlug is the English one, the locale the lesson fixtures are
	// written in.
	jsEnglishLandingSlug = "frontend-en"
)

// Starting a lesson records it and publishes the lesson-started fact. No
// course-started fact: the learner was already enrolled, and the CRM counts
// transitions, not requests.
func TestStartLessonRecordsProgressOnALesson(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	lesson := lessonBySlug(t, h, secondLessonSlug)

	res, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.StartLessonNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	row, err := h.DB.LessonProgress.Query().
		Where(lessonprogress.LessonID(lesson.ID), lessonprogress.UserID(actingUserID(t, h))).
		Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, "started", *row.State)

	owner := row.QueryEnrollment().OnlyX(ctx)
	assert.Equal(t, *lesson.CourseID, owner.CourseID, "the progress hangs off the learner's enrollment in this course")

	started := publishedLessonStarted(t, h)
	assert.Equal(t, secondLessonSlug, started.LessonSlug)
	assert.Equal(t, jsCourseSlug, started.CourseSlug)
	assert.Equal(t, 2, started.OccurrenceCount, "counts the Lesson Progress rows on this enrollment")
	assert.Empty(t, publishedOf[events.CourseStarted](t, h), "the learner was already enrolled")
}

// The first lesson of a course the learner has never touched enrolls them, and
// both facts are published — course-started counting their started enrollments.
func TestStartLessonEnrollsInAnUntouchedCourse(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	lesson := lessonLabeled(t, h, fixtures.LessonSingleLessonPrinting)

	res, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.StartLessonNoContent{}, res)

	count, err := h.DB.Enrollment.Query().
		Where(
			enrollment.UserID(actingUserID(t, h)),
			enrollment.CourseID(*lesson.CourseID),
		).
		Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "the enrollment is created by starting, never by viewing")

	enrolledCourse := h.DB.Course.GetX(ctx, *lesson.CourseID)
	assert.Equal(t, 1, enrolledCourse.EnrollmentsCount, "the catalogue figure counts the new enrollment")

	courseStarted := publishedOf[events.CourseStarted](t, h)
	require.Len(t, courseStarted, 1)
	assert.Equal(t, "single-lesson", courseStarted[0].Slug)
	assert.Equal(t, 2, courseStarted[0].OccurrenceCount, "the learner's started enrollments, fixture one included")

	started := publishedLessonStarted(t, h)
	assert.Equal(t, "printing", started.LessonSlug)
	assert.Equal(t, 1, started.OccurrenceCount)
}

// Starting is idempotent: pressing the button twice, or coming back to a lesson
// already passed, succeeds and changes nothing. Publishing a second fact would
// corrupt the occurrence counts the CRM consumes.
func TestStartLessonIsIdempotent(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	lesson := lessonBySlug(t, h, firstLessonSlug)

	before := h.DB.LessonProgress.Query().CountX(ctx)

	res, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.StartLessonNoContent{}, res)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx), "no row is written")
	assert.Zero(t, testsupport.PublishedCount(t, h.Services), "an already-finished lesson publishes nothing")

	row := h.DB.LessonProgress.Query().
		Where(lessonprogress.LessonID(lesson.ID), lessonprogress.UserID(actingUserID(t, h))).
		OnlyX(ctx)
	assert.Equal(t, "finished", *row.State, "the finished state is not rolled back to started")
}

// Lessons may be taken in any order (ADR-0016): the third lesson starts before
// the second, and only that lesson is recorded.
func TestStartLessonLetsALearnerTakeALessonOutOfOrder(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	lesson := lessonBySlug(t, h, thirdLessonSlug)
	before := h.DB.LessonProgress.Query().CountX(ctx)

	res, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.StartLessonNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	assert.Equal(t, before+1, h.DB.LessonProgress.Query().CountX(ctx), "one row, for this lesson")
	assert.Equal(t, thirdLessonSlug, publishedLessonStarted(t, h).LessonSlug)
}

// A lesson the current version does not contain counts for nothing, so a learner
// who finished one that was dropped has no more of the course done than before.
func TestStartLessonIgnoresLessonsDroppedFromTheCurrentVersion(t *testing.T) {
	// Their only finished lesson exists in the retired version alone.
	h := testsupport.NewHarnessAs(t, fixtures.UserLearnerDroppedLesson)
	ctx := context.Background()

	// The dropped lesson adds nothing to the course's completion.
	third := lessonBySlug(t, h, thirdLessonSlug)
	_, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(third.ID)})
	require.NoError(t, err)

	view := courseViewOf(t, h)
	assert.Equal(t, int32(0), view.Progress.Value.Completion, "a dropped lesson has no place in the course")
	require.Len(t, view.Lessons, 3)
}

// Starting a lesson is not finishing it: the page the learner lands on reads the
// course and still resumes at that lesson, with the first one checked.
func TestStartLessonLeavesTheLearnerAtTheFirstUnfinishedLesson(t *testing.T) {
	h := testsupport.NewHarness(t)
	lesson := lessonBySlug(t, h, secondLessonSlug)

	_, err := h.Client.StartLesson(t.Context(), api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)

	view := courseViewOf(t, h)
	assert.Equal(t, secondLessonSlug, view.Progress.Value.NextLesson.Value.Slug, "started is not finished")
	require.Len(t, view.Lessons, 3)
	assert.True(t, view.Lessons[0].Finished)
	assert.False(t, view.Lessons[1].Finished)
}

// A guest starts the same way and gets the same payload, with nothing stored:
// only a check changes what their cookie names.
func TestStartLessonAnswersAGuestWithoutStoringAnything(t *testing.T) {
	h := testsupport.NewVisitorHarness(t, progressCookie(jsCourseSlug, firstLessonSlug))
	ctx := t.Context()
	lesson := lessonBySlug(t, h, secondLessonSlug)
	rows := h.DB.LessonProgress.Query().CountX(ctx)
	enrollments := h.DB.Enrollment.Query().CountX(ctx)

	_, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)

	view := courseViewOf(t, h)
	assert.Equal(t, api.EnrollmentStateStarted, view.Progress.Value.State.Value, "derived from what they finished")
	assert.True(t, view.Lessons[0].Finished)
	assert.False(t, view.Lessons[1].Finished)

	assert.Equal(t, rows, h.DB.LessonProgress.Query().CountX(ctx))
	assert.Equal(t, enrollments, h.DB.Enrollment.Query().CountX(ctx))
	assert.Zero(t, testsupport.PublishedCount(t, h.Services), "starting is not a fact for someone with no account")
}

// A visitor with no cookie may start any lesson: nothing is stored, so there is
// nothing to refuse.
func TestStartLessonLetsAGuestTakeAnyLesson(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	lesson := lessonBySlug(t, h, thirdLessonSlug)

	res, err := h.Client.StartLesson(t.Context(), api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.StartLessonNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
}

// courseViewOf is the JavaScript course page as the harness's visitor reads it.
func courseViewOf(t *testing.T, h *testsupport.Harness) *api.CourseView {
	t.Helper()
	res, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)
	view, ok := res.(*api.CourseView)
	require.True(t, ok, "got %T", res)
	return view
}

// lessonBySlug is the JavaScript lesson with the slug: other fixture courses
// have a lesson of the same name (a course the loader built once has a
// hello-world too), so the slug alone does not say which.
func lessonBySlug(t *testing.T, h *testsupport.Harness, slug string) *ent.CourseLesson {
	t.Helper()
	return h.DB.CourseLesson.Query().
		Where(courselesson.SlugEQ(slug), courselesson.CourseID(fixtures.ID(fixtures.CourseJavascript))).
		OnlyX(t.Context())
}

// lessonLabeled is the fixture lesson with the label.
func lessonLabeled(t *testing.T, h *testsupport.Harness, label fixtures.Label) *ent.CourseLesson {
	t.Helper()
	return h.DB.CourseLesson.GetX(t.Context(), fixtures.ID(label))
}

func actingUserID(t *testing.T, h *testsupport.Harness) int {
	t.Helper()
	u, ok := testsupport.HarnessUser(h)
	require.True(t, ok)
	return u
}

func publishedLessonStarted(t *testing.T, h *testsupport.Harness) events.LessonStarted {
	t.Helper()
	found := publishedOf[events.LessonStarted](t, h)
	require.Len(t, found, 1)
	return found[0]
}

func publishedOf[T events.Event](t *testing.T, h *testsupport.Harness) []T {
	t.Helper()
	return testsupport.PublishedEvents[T](t, h.Services)
}

// A lesson that does not exist is a 404: the client can stop looking.
func TestStartLessonAnswersAnUnknownLessonAsNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.StartLesson(t.Context(), api.StartLessonParams{ID: 2_000_000_000})
	require.NoError(t, err)

	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

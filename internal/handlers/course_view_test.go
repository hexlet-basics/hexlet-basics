package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/enrollment"
	"hexletbasics/internal/api"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// A signed-in learner gets their position in the course: the state, the
// percentage, where they resume, and per-lesson locks and checkmarks. The
// acting user has finished the first of the three current lessons.
func TestGetCourseReturnsProgressForASignedInLearner(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view, ok := res.(*api.CourseView)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	require.False(t, view.Enrollment.Null)
	progress := view.Enrollment.Value.Progress

	assert.Equal(t, api.EnrollmentStateStarted, progress.State.Value)
	assert.Equal(t, int32(33), progress.Completion, "one of three current lessons")
	assert.Equal(t, secondLessonSlug, progress.NextLessonSlug.Value, "resume at the first unfinished lesson")
	assert.Equal(t, int32(1), progress.FurthestFinishedPosition)

	require.Len(t, progress.Lessons, 3)
	assert.Equal(t, []string{firstLessonSlug, secondLessonSlug, thirdLessonSlug},
		[]string{progress.Lessons[0].Slug, progress.Lessons[1].Slug, progress.Lessons[2].Slug},
		"in course order")

	assert.True(t, progress.Lessons[0].Finished)
	assert.True(t, progress.Lessons[0].Available)
	assert.False(t, progress.Lessons[1].Finished)
	assert.True(t, progress.Lessons[1].Available, "one past the furthest finished lesson")
	assert.False(t, progress.Lessons[2].Available, "two past it is still locked")
}

// A gap does not block: a learner who finished the third lesson but not the
// second sees everything up to one past their furthest, and the course is still
// unfinished.
func TestGetCourseTreatsAGapAsProgressNotAsABlock(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	finishLesson(t, h, thirdLessonSlug)

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view := res.(*api.CourseView)
	progress := view.Enrollment.Value.Progress

	assert.Equal(t, int32(66), progress.Completion, "two of three, never clamped")
	assert.Equal(t, int32(3), progress.FurthestFinishedPosition)
	assert.Equal(t, secondLessonSlug, progress.NextLessonSlug.Value, "the gap is where they resume")
	assert.True(t, progress.Lessons[1].Available)
	assert.Equal(t, api.EnrollmentStateStarted, progress.State.Value, "a gap leaves the course unfinished")
}

// A visitor with no session and no cookie has a position too: the beginning,
// with the first lesson open. It is returned rather than left out, because a
// client deriving "the first one is available" would be a second copy of the
// gate.
func TestGetCourseReturnsTheStartingPositionForAnonymousVisitors(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view, ok := res.(*api.CourseView)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.True(t, view.Enrollment.Null, "no enrollment for a visitor with no session")
	assert.Equal(t, jsCourseSlug, view.Course.Slug)
	assert.Len(t, view.Lessons, 3, "the public lesson list does not depend on a session")

	require.False(t, view.Progress.Null)
	progress := view.Progress.Value
	assert.True(t, progress.State.Null, "nothing started")
	assert.Equal(t, int32(0), progress.Completion)
	assert.Equal(t, int32(0), progress.FurthestFinishedPosition)
	assert.Equal(t, firstLessonSlug, progress.NextLessonSlug.Value)
	assert.True(t, progress.Lessons[0].Available)
	assert.False(t, progress.Lessons[1].Available)
}

// A guest carrying a cookie gets exactly the shape a signed-in learner gets,
// derived from the one lesson the cookie names: everything before it is
// finished, the next one is open.
func TestGetCourseReturnsTheSameShapeForAGuest(t *testing.T) {
	h := testsupport.NewVisitorHarness(t, progressCookie(jsCourseSlug, secondLessonSlug))
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view := res.(*api.CourseView)

	require.True(t, view.Enrollment.Null, "a guest has no enrollment")
	require.False(t, view.Progress.Null)
	progress := view.Progress.Value

	assert.Equal(t, api.EnrollmentStateStarted, progress.State.Value)
	assert.Equal(t, int32(66), progress.Completion, "two of three")
	assert.Equal(t, int32(2), progress.FurthestFinishedPosition)
	assert.Equal(t, thirdLessonSlug, progress.NextLessonSlug.Value)
	assert.True(t, progress.Lessons[0].Finished, "the prefix is implied by the furthest lesson")
	assert.True(t, progress.Lessons[1].Finished)
	assert.True(t, progress.Lessons[2].Available)
}

// A cookie naming a lesson the current version dropped resolves to no position,
// so the guest starts that course over rather than being stranded.
func TestGetCourseResetsAGuestPositionTheCurrentVersionDropped(t *testing.T) {
	h := testsupport.NewVisitorHarness(t, progressCookie(jsCourseSlug, "removed-in-a-later-build"))

	res, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view := res.(*api.CourseView)

	require.False(t, view.Progress.Null)
	assert.Equal(t, int32(0), view.Progress.Value.FurthestFinishedPosition)
	assert.True(t, view.Progress.Value.Lessons[0].Available)
	assert.False(t, view.Progress.Value.Lessons[1].Available)
}

// A tampered cookie is not progress: it is discarded on the way in, leaving the
// visitor where a first visit leaves them.
func TestGetCourseIgnoresAForgedGuestCookie(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.ForgeGuestCookie(t, h, progressCookie(jsCourseSlug, thirdLessonSlug))

	res, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view := res.(*api.CourseView)

	require.False(t, view.Progress.Null)
	assert.Equal(t, int32(0), view.Progress.Value.FurthestFinishedPosition,
		"an unsigned cookie cannot mint completions")
}

// An unknown slug is a 404, not an empty page.
func TestGetCourseIsNotFoundForAnUnknownSlug(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: "no-such-course"})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// progressCookie is a visitor's whole state: the furthest lesson they finished
// in one course.
func progressCookie(courseSlug, lessonSlug string) progress.GuestProgress {
	return progress.GuestProgress{}.Record(courseSlug, lessonSlug)
}

// finishLesson marks the acting learner's progress on a lesson finished,
// creating the row when they have not taken it.
func finishLesson(t *testing.T, h *testsupport.Harness, slug string) {
	t.Helper()
	ctx := t.Context()
	lesson := lessonBySlug(t, h, slug)
	owner := h.DB.Enrollment.Query().
		Where(
			enrollment.UserID(actingUserID(t, h)),
			enrollment.CourseID(*lesson.CourseID),
		).
		OnlyX(ctx)

	h.DB.LessonProgress.Create().
		SetUserID(actingUserID(t, h)).
		SetCourseID(*lesson.CourseID).
		SetEnrollmentID(owner.ID).
		SetLessonID(lesson.ID).
		SetState("finished").
		SaveX(ctx)
}

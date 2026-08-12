package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// The harness acts as the fixture admin (alice), who has finished the first
// lesson of the JavaScript course and is enrolled in it. Lessons are addressed
// by slug throughout: fixture ids are legacy crc32 values.
const (
	firstLessonSlug  = "hello-world" // position 1, finished by the acting user
	secondLessonSlug = "variables"   // position 2, the next available one
	thirdLessonSlug  = "strings"     // position 3, beyond the gate
	jsCourseSlug     = "javascript"
)

// Starting the lesson one past the furthest finished one records it and
// publishes the lesson-started fact. No course-started fact: the learner was
// already enrolled, and the CRM counts transitions, not requests.
func TestStartLessonRecordsProgressOnTheNextAvailableLesson(t *testing.T) {
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
	assert.Empty(t, publishedOf[events.CourseStarted](h), "the learner was already enrolled")
}

// The first lesson of a course the learner has never touched enrolls them, and
// both facts are published — course-started counting their started enrollments.
func TestStartLessonEnrollsInAnUntouchedCourse(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	lesson := seedCourse(t, h, "go-basics", "printing")

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

	courseStarted := publishedOf[events.CourseStarted](h)
	require.Len(t, courseStarted, 1)
	assert.Equal(t, "go-basics", courseStarted[0].Slug)
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
	assert.Empty(t, h.Events.Published, "an already-finished lesson publishes nothing")

	row := h.DB.LessonProgress.Query().
		Where(lessonprogress.LessonID(lesson.ID), lessonprogress.UserID(actingUserID(t, h))).
		OnlyX(ctx)
	assert.Equal(t, "finished", *row.State, "the finished state is not rolled back to started")
}

// Two past the furthest finished lesson is refused, and nothing is written.
func TestStartLessonRefusesALessonBeyondTheGate(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	lesson := lessonBySlug(t, h, thirdLessonSlug)
	before := h.DB.LessonProgress.Query().CountX(ctx)

	res, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(lesson.ID)})
	require.NoError(t, err)
	require.IsType(t, &api.StartLessonConflict{}, res)
	assert.Equal(t, http.StatusConflict, h.LastStatus())

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx))
	assert.Empty(t, h.Events.Published)
}

// A lesson the current version does not contain contributes no Position, so a
// learner whose furthest finished lesson was dropped is gated from their
// furthest surviving one — not from the dropped lesson's old position.
func TestStartLessonIgnoresLessonsDroppedFromTheCurrentVersion(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	user := actingUserID(t, h)
	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(ctx)
	owner := h.DB.Enrollment.Query().
		Where(enrollment.UserID(user), enrollment.CourseID(js.ID)).
		OnlyX(ctx)

	// A lesson that exists only in the retired version, finished by the learner.
	dropped := h.DB.CourseLesson.Create().
		SetSlug("removed-in-a-later-build").
		SetCourseID(js.ID).
		SetNaturalOrder(9).
		SaveX(ctx)
	retired := retiredVersionID(t, h, js)
	h.DB.CourseLessonVersion.Create().
		SetCourseID(js.ID).
		SetCourseVersionID(retired).
		SetLessonID(dropped.ID).
		SetModuleVersionID(seedModuleVersion(t, h, js.ID, retired)).
		SetNaturalOrder(9).
		SaveX(ctx)
	h.DB.LessonProgress.Create().
		SetUserID(user).
		SetCourseID(js.ID).
		SetEnrollmentID(owner.ID).
		SetLessonID(dropped.ID).
		SetState("finished").
		SaveX(ctx)

	// Position 9 must not raise the gate: the third lesson stays out of reach.
	third := lessonBySlug(t, h, thirdLessonSlug)
	res, err := h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(third.ID)})
	require.NoError(t, err)
	assert.IsType(t, &api.StartLessonConflict{}, res, "a dropped lesson has no position")

	// While the lesson one past their furthest SURVIVING one still opens.
	second := lessonBySlug(t, h, secondLessonSlug)
	res, err = h.Client.StartLesson(ctx, api.StartLessonParams{ID: int32(second.ID)})
	require.NoError(t, err)
	assert.IsType(t, &api.StartLessonNoContent{}, res)
}

func lessonBySlug(t *testing.T, h *testsupport.Harness, slug string) *ent.CourseLesson {
	t.Helper()
	return h.DB.CourseLesson.Query().
		Where(courselesson.SlugEQ(slug)).
		OnlyX(t.Context())
}

func actingUserID(t *testing.T, h *testsupport.Harness) int {
	t.Helper()
	u, ok := testsupport.HarnessUser(h)
	require.True(t, ok)
	return u
}

func retiredVersionID(t *testing.T, h *testsupport.Harness, c *ent.Course) int {
	t.Helper()
	require.NotNil(t, c.CurrentVersionID)
	retired, err := h.DB.CourseVersion.Query().
		Where(
			courseversion.CourseID(c.ID),
			courseversion.IDNEQ(*c.CurrentVersionID),
		).
		First(t.Context())
	require.NoError(t, err, "fixtures must carry a retired course version")
	return retired.ID
}

// seedModuleVersion creates the module and module version a lesson version's
// NOT NULL FK requires. Fixtures can point at a module version that does not
// exist (loading disables FKs); a live insert cannot.
func seedModuleVersion(t *testing.T, h *testsupport.Harness, courseID, versionID int) int {
	t.Helper()
	module := h.DB.CourseModule.Create().
		SetSlug("basics").
		SetCourseID(courseID).
		SaveX(t.Context())
	return h.DB.CourseModuleVersion.Create().
		SetCourseID(courseID).
		SetCourseVersionID(versionID).
		SetModuleID(module.ID).
		SaveX(t.Context()).ID
}

// seedCourse builds a minimal course the acting user has never touched: one
// promoted version carrying one lesson at position 1.
func seedCourse(t *testing.T, h *testsupport.Harness, courseSlug, lessonSlug string) *ent.CourseLesson {
	t.Helper()
	ctx := t.Context()

	created := h.DB.Course.Create().SetSlug(courseSlug).SaveX(ctx)
	version := h.DB.CourseVersion.Create().
		SetCourseID(created.ID).
		SetState("ready").
		SaveX(ctx)
	promoted := created.Update().SetCurrentVersionID(version.ID).SaveX(ctx)

	lesson := h.DB.CourseLesson.Create().
		SetSlug(lessonSlug).
		SetCourseID(promoted.ID).
		SetNaturalOrder(1).
		SaveX(ctx)
	h.DB.CourseLessonVersion.Create().
		SetCourseID(promoted.ID).
		SetCourseVersionID(version.ID).
		SetLessonID(lesson.ID).
		SetModuleVersionID(seedModuleVersion(t, h, promoted.ID, version.ID)).
		SetNaturalOrder(1).
		SaveX(ctx)
	return lesson
}

func publishedLessonStarted(t *testing.T, h *testsupport.Harness) events.LessonStarted {
	t.Helper()
	found := publishedOf[events.LessonStarted](h)
	require.Len(t, found, 1)
	return found[0]
}

func publishedOf[T events.Event](h *testsupport.Harness) []T {
	var found []T
	for _, e := range h.Events.Published {
		if typed, ok := e.(T); ok {
			found = append(found, typed)
		}
	}
	return found
}

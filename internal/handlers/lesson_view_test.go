package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// The player payload: the lesson's prose for the request locale, the code and
// tests of the course's current build, and where the learner stands.
func TestGetCourseLessonReturnsThePlayerPayload(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view, ok := res.(*api.CourseLessonView)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	lesson := view.Lesson
	assert.Equal(t, secondLessonSlug, lesson.Slug)
	assert.Equal(t, jsCourseSlug, lesson.Course.Slug)
	assert.Equal(t, "Variables", lesson.Name.Value)
	assert.Equal(t, "en", lesson.Locale.Value)
	assert.Equal(t, int32(2), lesson.NaturalOrder)
	assert.Equal(t, "A variable is a name bound to a value.", lesson.Theory.Value)
	assert.Equal(t, "Assign the string `hello` to `greeting`.", lesson.Instructions.Value)
	assert.Equal(t, "let greeting = '';\n", lesson.PreparedCode.Value)
	assert.Equal(t, "assert.equal(greeting, 'hello');\n", lesson.TestCode.Value)

	assert.Equal(t, []string{"Names are case sensitive", "Prefer const until you need to reassign"}, lesson.Tips,
		"decoded from the YAML array the loader writes")
	require.Len(t, lesson.Definitions, 1)
	assert.Equal(t, "variable", lesson.Definitions[0].Name)

	assert.Equal(t,
		"https://github.com/hexlet-basics/exercises-javascript/blob/main/modules/10-basics/20-variables/en/README.md",
		lesson.SourceCodeUrl.Value,
		"the README a reader can propose a fix to")

	// The version the client submits its solution against.
	assert.Equal(t, lesson.VersionId, lesson.Version.Value)
	assert.NotZero(t, lesson.VersionId)
}

// The player titles the page with the course's landing copy, not with the
// course's own name — two different strings, and the marketing one is what a
// learner recognises. The page flagged main wins, exactly as the course read
// picks the canonical one.
func TestGetCourseLessonCarriesTheCoursesMainLandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res.(*api.CourseLessonView)

	require.False(t, view.LandingPage.Null)
	landing := view.LandingPage.Value
	assert.Equal(t, "javascript-ru", landing.Slug,
		"the page flagged main wins, though the archived one has the lower id")
	assert.Equal(t, "JavaScript", landing.Name)
	assert.NotEqual(t, view.Lesson.Course.Name.Value, landing.Name,
		"the two names are what makes this field necessary")
}

// The navigation tab's rows: every lesson of the current version, in course
// order, named. The player joins this to the progress payload by slug, so the
// two arrays line up entry for entry.
func TestGetCourseLessonCarriesTheCoursesLessonList(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res.(*api.CourseLessonView)

	require.Len(t, view.Lessons, 3)
	assert.Equal(t,
		[]string{firstLessonSlug, secondLessonSlug, thirdLessonSlug},
		[]string{view.Lessons[0].Slug, view.Lessons[1].Slug, view.Lessons[2].Slug},
		"in course order")
	assert.Equal(t, "Variables", view.Lessons[1].Name.Value, "named, or the nav list has no labels")

	require.False(t, view.Progress.Null)
	assert.Equal(t, view.Lessons[2].Slug, view.Progress.Value.Lessons[2].Slug,
		"the two arrays are joined by slug, so they must agree entry for entry")
}

// A course nobody has written copy for is still readable: the landing page is
// null rather than the read failing.
func TestGetCourseLessonWithoutALandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()

	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(ctx)
	h.DB.LandingPage.Delete().Where(landingpage.CourseID(js.ID)).ExecX(ctx)

	res, err := h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view, ok := res.(*api.CourseLessonView)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.True(t, view.LandingPage.Null)
}

// Reading a lesson never starts it: the frontend preloads routes on hover, so a
// read with that effect would enroll a learner in every lesson they point at.
func TestGetCourseLessonStartsNothing(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	rows := h.DB.LessonProgress.Query().CountX(ctx)

	_, err := h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)

	assert.Equal(t, rows, h.DB.LessonProgress.Query().CountX(ctx))
	assert.Empty(t, h.Events.Published, "reading is not a transition")
}

// The read carries the same progress the course page carries, so the player
// renders locks and the "next" target from the server's answer.
func TestGetCourseLessonReturnsProgress(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res.(*api.CourseLessonView)

	require.False(t, view.Progress.Null)
	progress := view.Progress.Value
	assert.Equal(t, int32(1), progress.FurthestFinishedPosition)
	assert.Equal(t, secondLessonSlug, progress.NextLessonSlug.Value)
	require.Len(t, progress.Lessons, 3)
	assert.True(t, progress.Lessons[1].Available, "the learner may take this one")
	assert.False(t, progress.Lessons[2].Available)
}

// Theory stays public: a visitor reads any lesson, including one they cannot
// take yet, which is what keeps the page answering the search that led here.
func TestGetCourseLessonServesAVisitorALockedLesson(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       thirdLessonSlug,
	})
	require.NoError(t, err)
	view, ok := res.(*api.CourseLessonView)
	require.True(t, ok, "got %T", res)

	assert.Equal(t, thirdLessonSlug, view.Lesson.Slug)
	require.False(t, view.Progress.Null)
	assert.False(t, view.Progress.Value.Lessons[2].Available,
		"readable, and visibly out of reach")
	assert.True(t, view.Progress.Value.Lessons[0].Available)
}

// A lesson of another course, an unknown slug, and a lesson the current build
// dropped are all the same answer: this course has no such lesson.
func TestGetCourseLessonIsNotFoundOutsideTheCurrentBuild(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()

	res, err := h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       "no-such-lesson",
	})
	require.NoError(t, err)
	require.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	// A lesson that exists but is not part of the current version.
	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(ctx)
	h.DB.CourseLesson.Create().
		SetSlug("removed-in-a-later-build").
		SetCourseID(js.ID).
		SetNaturalOrder(9).
		SaveX(ctx)

	res, err = h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       "removed-in-a-later-build",
	})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
}

// A lesson with no translation for the request locale is not a lesson this
// visitor can read, which is how legacy treated it too.
func TestGetCourseLessonIsNotFoundInAnUntranslatedLocale(t *testing.T) {
	h := testsupport.NewHarness(t)

	testsupport.SpeakTo(h, "es")

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
}

package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
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
	view := res
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

	// The build's own name, which the page's meta description opens with.
	assert.Equal(t, "Основы JavaScript", lesson.Course.CurrentVersion.Value.Name.Value)
}

// The player titles the page with the course's landing copy, not with the
// course's own name — two different strings, and the marketing one is what a
// learner recognises. The page flagged main wins, exactly as the course read
// picks the canonical one, and it is the one of the request's locale: the course
// page is served only under its own locale, so the link back must name that one.
func TestGetCourseLessonCarriesTheCoursesMainLandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       firstLessonSlug,
	})
	require.NoError(t, err)
	view := res

	require.False(t, view.LandingPage.Null)
	landing := view.LandingPage.Value
	assert.Equal(t, "javascript-ru", landing.Slug,
		"the page flagged main wins, though the archived one has the lower id")
	assert.Equal(t, "JavaScript", landing.Name)
	assert.NotEqual(t, view.Lesson.Course.Name.Value, landing.Name,
		"the two names are what makes this field necessary")
}

func TestGetCourseLessonCarriesTheLandingPageOfTheRequestLocale(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res

	require.False(t, view.LandingPage.Null)
	assert.Equal(t, jsEnglishLandingSlug, view.LandingPage.Value.Slug)
}

// The navigation tab's rows: every lesson of the current version, in course
// order, named, each with the visitor's standing in it.
func TestGetCourseLessonCarriesTheCoursesLessonList(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res

	require.Len(t, view.Lessons, 3)
	assert.Equal(t,
		[]string{firstLessonSlug, secondLessonSlug, thirdLessonSlug},
		[]string{view.Lessons[0].Slug, view.Lessons[1].Slug, view.Lessons[2].Slug},
		"in course order")
	assert.Equal(t, "Variables", view.Lessons[1].Name.Value, "named, or the nav list has no labels")

	assert.Equal(t, []int32{1, 2, 3},
		[]int32{view.Lessons[0].Position, view.Lessons[1].Position, view.Lessons[2].Position},
		"a position is the lesson's place in the course")
}

// A course nobody has written copy for is still readable: the landing page is
// null rather than the read failing.
func TestGetCourseLessonWithoutALandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()

	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(ctx)
	h.DB.CourseLandingPage.Delete().Where(courselandingpage.CourseID(js.ID)).ExecX(ctx)

	res, err := h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res
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
	assert.Zero(t, testsupport.PublishedCount(t, h.Services), "reading is not a transition")
}

// The read carries the same progress the course page carries, so the player
// renders the checkmarks and where to resume from the server's answer.
func TestGetCourseLessonReturnsProgress(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res

	require.False(t, view.Progress.Null)
	assert.Equal(t, secondLessonSlug, view.Progress.Value.NextLesson.Value.Slug)
	require.Len(t, view.Lessons, 3)
	assert.True(t, view.Lessons[0].Finished)
	assert.False(t, view.Lessons[1].Finished)
}

// Back and forward follow course order, whatever the learner has finished: the
// first lesson has nothing before it, the last nothing after it.
func TestGetCourseLessonNamesTheLessonsOnEitherSide(t *testing.T) {
	h := testsupport.NewHarness(t)
	adjacent := func(slug string) api.AdjacentLessons {
		res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{CourseSlug: jsCourseSlug, Slug: slug})
		require.NoError(t, err)
		return res.Adjacent
	}

	first := adjacent(firstLessonSlug)
	assert.True(t, first.Previous.Null)
	assert.Equal(t, secondLessonSlug, first.Next.Value.Slug)

	middle := adjacent(secondLessonSlug)
	assert.Equal(t, firstLessonSlug, middle.Previous.Value.Slug)
	assert.Equal(t, thirdLessonSlug, middle.Next.Value.Slug)
	assert.Equal(t, "Variables", adjacent(firstLessonSlug).Next.Value.Name.Value, "named, not just addressed")

	last := adjacent(thirdLessonSlug)
	assert.Equal(t, secondLessonSlug, last.Previous.Value.Slug)
	assert.True(t, last.Next.Null)
}

// The Next Lesson is where the learner resumes, not the lesson after this one:
// on the last lesson, with the second still to do, "next" has nowhere to go but
// the learner still resumes at the second.
func TestGetCourseLessonKeepsTheNextLessonApartFromAdjacentLessons(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       thirdLessonSlug,
	})
	require.NoError(t, err)
	view := res

	assert.True(t, view.Adjacent.Next.Null)
	assert.Equal(t, secondLessonSlug, view.Progress.Value.NextLesson.Value.Slug)
}

// Theory is public: a visitor reads any lesson with no session, which is what
// keeps the page answering the search that led here.
func TestGetCourseLessonServesAVisitorAnyLesson(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       thirdLessonSlug,
	})
	require.NoError(t, err)
	view := res

	assert.Equal(t, thirdLessonSlug, view.Lesson.Slug)
	require.False(t, view.Progress.Null)
	assert.False(t, view.Lessons[2].Finished, "a first visit has finished nothing")
}

// A lesson of another course, an unknown slug, and a lesson the current build
// dropped are all the same answer: this course has no such lesson.
func TestGetCourseLessonIsNotFoundOutsideTheCurrentBuild(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()

	_, err := h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       "no-such-lesson",
	})
	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	// A lesson that exists but is not part of the current version: the fixtures
	// carry one that only the retired version has.
	_, err = h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       "removed-in-a-later-build",
	})
	requireProblem(t, err, http.StatusNotFound)
}

// A lesson with no translation for the request locale is not a lesson this
// visitor can read, which is how legacy treated it too.
func TestGetCourseLessonIsNotFoundInAnUntranslatedLocale(t *testing.T) {
	h := testsupport.NewHarness(t)

	testsupport.SpeakTo(h, "es")

	_, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	requireProblem(t, err, http.StatusNotFound)
}

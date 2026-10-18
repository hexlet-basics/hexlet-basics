package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// A signed-in learner gets their position in the course: the state, the
// percentage, where they resume, and per-lesson locks and checkmarks. The
// acting user has finished the first of the three current lessons.
func TestGetCourseReturnsProgressForASignedInLearner(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)

	view := res

	assert.Equal(t, http.StatusOK, h.LastStatus())

	require.False(t, view.Enrollment.Null)
	progress := view.Progress.Value

	assert.Equal(t, api.EnrollmentStateStarted, progress.State.Value)
	assert.Equal(t, int32(33), progress.Completion, "one of three current lessons")
	assert.Equal(t, secondLessonSlug, progress.NextLesson.Value.Slug, "resume at the first unfinished lesson")

	require.Len(t, view.Lessons, 3)
	assert.Equal(t, []string{firstLessonSlug, secondLessonSlug, thirdLessonSlug},
		[]string{view.Lessons[0].Slug, view.Lessons[1].Slug, view.Lessons[2].Slug},
		"in course order")

	assert.True(t, view.Lessons[0].Finished)
	assert.False(t, view.Lessons[1].Finished)
	assert.False(t, view.Lessons[2].Finished)
}

// Lessons may be taken in any order: a learner who finished the third lesson but
// not the second has exactly those marked, and the course is still unfinished.
func TestGetCourseTreatsAGapAsProgressNotAsABlock(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerWithAGap)
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)

	view := res
	progress := view.Progress.Value

	assert.Equal(t, int32(66), progress.Completion, "two of three, never clamped")
	assert.Equal(t, secondLessonSlug, progress.NextLesson.Value.Slug, "the gap is where they resume")
	assert.False(t, view.Lessons[1].Finished)
	assert.True(t, view.Lessons[2].Finished)
	assert.Equal(t, api.EnrollmentStateStarted, progress.State.Value, "a gap leaves the course unfinished")
}

// A visitor with no session and no cookie has a standing too: the beginning,
// with nothing finished. It is returned rather than left out, so the page
// renders from the same shape for everyone.
func TestGetCourseReturnsTheStartingPositionForAnonymousVisitors(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)

	view := res

	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.True(t, view.Enrollment.Null, "no enrollment for a visitor with no session")
	assert.Equal(t, jsCourseSlug, view.Course.Slug)
	assert.Len(t, view.Lessons, 3, "the public lesson list does not depend on a session")

	require.False(t, view.LandingPage.Null)
	assert.Equal(t, jsEnglishLandingSlug, view.LandingPage.Value.Slug, "the page the URL names")

	require.False(t, view.Progress.Null)
	progress := view.Progress.Value
	assert.True(t, progress.State.Null, "nothing started")
	assert.Equal(t, int32(0), progress.Completion)
	assert.Equal(t, firstLessonSlug, progress.NextLesson.Value.Slug)
	assert.False(t, view.Lessons[0].Finished)
	assert.False(t, view.Lessons[1].Finished)
}

// The learning program: the current version's modules in the request locale,
// in the order the version gives them — not by id, which the fixtures set to
// disagree — each naming its lessons in course order.
func TestGetCourseReturnsTheModulesInCourseOrder(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)

	view := res

	names := make([]string, 0, len(view.Modules))
	for _, m := range view.Modules {
		names = append(names, m.Name.Value)
	}

	assert.Equal(t, []string{"Основы", "Арифметика", "Строки"}, names,
		"ru translations only, ordered by the module version's order")

	require.NotEmpty(t, view.Modules)
	assert.Equal(t, []string{firstLessonSlug}, moduleLessonSlugs(view.Modules[0]),
		"a module lists the lessons translated into the locale: the others have no page here")
	assert.Empty(t, view.Modules[1].Lessons)
}

// The landing page's questions and answers travel with the course, oldest
// first, and only the course's own.
func TestGetCourseReturnsTheLandingPageQuestions(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: "go-ru"})
	require.NoError(t, err)

	view := res

	questions := make([]string, 0, len(view.QnaItems))
	for _, item := range view.QnaItems {
		questions = append(questions, item.Question)
	}

	assert.Equal(t, []string{"Сколько длится курс Go?", "Есть ли сертификат?"}, questions)

	res, err = h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)
	assert.Empty(t, res.QnaItems, "a landing page with no questions has none")
}

// A guest carrying a cookie gets exactly the shape a signed-in learner gets,
// derived from the lessons the cookie names, and only those.
func TestGetCourseReturnsTheSameShapeForAGuest(t *testing.T) {
	t.Parallel()
	h := testsupport.NewVisitorHarness(t, progressCookie(firstLessonSlug, secondLessonSlug))
	ctx := context.Background()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)

	view := res

	require.True(t, view.Enrollment.Null, "a guest has no enrollment")
	require.False(t, view.Progress.Null)
	progress := view.Progress.Value

	assert.Equal(t, api.EnrollmentStateStarted, progress.State.Value)
	assert.Equal(t, int32(66), progress.Completion, "two of three")
	assert.Equal(t, thirdLessonSlug, progress.NextLesson.Value.Slug)
	assert.True(t, view.Lessons[0].Finished)
	assert.True(t, view.Lessons[1].Finished)
	assert.False(t, view.Lessons[2].Finished)
}

// A visitor who took the third lesson first has that one marked and not the two
// before it: a later lesson no longer implies the earlier ones.
func TestGetCourseMarksOnlyTheLessonsAGuestFinished(t *testing.T) {
	t.Parallel()
	h := testsupport.NewVisitorHarness(t, progressCookie(thirdLessonSlug))

	res, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)

	view := res
	progress := view.Progress.Value

	assert.Equal(t, int32(33), progress.Completion)
	assert.Equal(t, firstLessonSlug, progress.NextLesson.Value.Slug)
	assert.False(t, view.Lessons[0].Finished)
	assert.False(t, view.Lessons[1].Finished)
	assert.True(t, view.Lessons[2].Finished)
}

// A cookie naming a lesson the current version dropped resolves to nothing, so
// the guest starts that course over rather than being stranded.
func TestGetCourseResetsAGuestPositionTheCurrentVersionDropped(t *testing.T) {
	t.Parallel()
	h := testsupport.NewVisitorHarness(t, progressCookie("removed-in-a-later-build"))

	res, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: jsEnglishLandingSlug})
	require.NoError(t, err)

	view := res

	require.False(t, view.Progress.Null)
	assert.Equal(t, int32(0), view.Progress.Value.Completion)
	assert.False(t, view.Lessons[0].Finished)
}

// A tampered cookie is not progress: it is discarded on the way in, leaving the
// visitor where a first visit leaves them.
func TestGetCourseIgnoresAForgedGuestCookie(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")
	testsupport.ForgeGuestCookie(t, h, progressCookie(thirdLessonSlug))

	res, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)

	view := res

	require.False(t, view.Progress.Null)
	assert.Equal(t, int32(0), view.Progress.Value.Completion,
		"an unsigned cookie cannot mint completions")
}

// An unknown slug is a 404, not an empty page.
func TestGetCourseIsNotFoundForAnUnknownSlug(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	_, err := h.Client.GetCourse(context.Background(), api.GetCourseParams{Slug: "no-such-course"})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// The slug names a landing page, as in legacy: the course's own slug is not a
// URL, and a landing page of another locale is not served under this one. The
// page that comes back is the one the slug names, whichever course page that is.
func TestGetCourseResolvesALandingPageSlugInTheRequestLocale(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	ctx := t.Context()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)

	view := res
	assert.Equal(t, jsCourseSlug, view.Course.Slug)
	assert.Equal(t, jsLandingSlug, view.LandingPage.Value.Slug)

	_, err = h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsCourseSlug})
	require.Error(t, err, "the course slug is not a landing page slug")
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	_, err = h.Client.GetCourse(ctx, api.GetCourseParams{Slug: "go-en"})
	require.Error(t, err, "an English landing page is not served to a Russian request")
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	res, err = h.Client.GetCourse(ctx, api.GetCourseParams{Slug: "frontend-ru"})
	require.NoError(t, err)
	assert.Equal(t, "frontend-ru", res.LandingPage.Value.Slug,
		"a secondary landing page keeps its own copy, not the course's main page")
}

// A landing page that was archived says where its URL moved, as legacy's
// redirect_archived_language chose: the page it names, else the course's main
// page in the request locale. A live page moves nowhere.
func TestGetCourseSaysWhereAnArchivedLandingPageMoved(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	ctx := t.Context()

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)
	assert.True(t, res.MovedTo.Null, "a live page stays where it is")

	res, err = h.Client.GetCourse(ctx, api.GetCourseParams{Slug: "javascript-ru-archived"})
	require.NoError(t, err)
	assert.Equal(t, jsLandingSlug, res.MovedTo.Value,
		"with no page named, it moves to the course's main page")

	named := h.DB.CourseLandingPage.Query().
		Where(courselandingpage.SlugEQ("nodejs-ru")).OnlyX(ctx)
	h.DB.CourseLandingPage.Update().
		Where(courselandingpage.SlugEQ("javascript-ru-archived")).
		SetLandingPageToRedirectID(named.ID).ExecX(ctx)

	res, err = h.Client.GetCourse(ctx, api.GetCourseParams{Slug: "javascript-ru-archived"})
	require.NoError(t, err)
	assert.Equal(t, "nodejs-ru", res.MovedTo.Value,
		"the page it names wins over the main one")
}

// An archived page with nowhere to go is gone, not an empty redirect.
func TestGetCourseIsNotFoundForAnArchivedPageWithNowhereToGo(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "ru")

	ctx := t.Context()
	h.DB.CourseLandingPage.Update().
		Where(courselandingpage.SlugEQ(jsLandingSlug)).
		SetMain(false).ExecX(ctx)

	_, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: "javascript-ru-archived"})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// progressCookie is a visitor's whole state: the lessons they finished in the
// JavaScript course.
func progressCookie(lessonSlugs ...string) progress.GuestProgress {
	guest := progress.GuestProgress{}
	for _, slug := range lessonSlugs {
		guest = guest.Record(jsCourseSlug, slug)
	}

	return guest
}

// moduleLessonSlugs are the slugs of a module's lessons, in the order the module
// lists them.
func moduleLessonSlugs(m api.CourseModule) []string {
	slugs := make([]string, 0, len(m.Lessons))
	for _, lesson := range m.Lessons {
		slugs = append(slugs, lesson.Slug)
	}

	return slugs
}

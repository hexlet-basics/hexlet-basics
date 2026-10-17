package coursepages_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/internal/coursepages"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/landingpages"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

const (
	jsCourse      = "javascript"
	firstLesson   = "hello-world"
	secondLesson  = "variables"
	thirdLesson   = "strings"
	enLanding     = "frontend-en"
	ruLanding     = "javascript-ru"
	archivedRu    = "javascript-ru-archived"
	repoBase      = "https://github.com/hexlet-basics"
	untranslated  = "ru" // lessons are written in English only
	secondLessonN = "Variables"
)

type rig struct {
	db    *ent.Client
	pages *coursepages.Pages
}

func newRig(t *testing.T) *rig {
	t.Helper()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)
	tracker := progress.New(db, transactor, services.Events, testsupport.NewMocks().Runner)
	return &rig{db: db, pages: coursepages.New(db, landingpages.NewCatalog(db), tracker, repoBase)}
}

func (r *rig) learner(t *testing.T) progress.Learner {
	t.Helper()
	return progress.Learner{UserID: r.db.User.GetX(t.Context(), fixtures.ID(fixtures.UserAdminLearner)).ID}
}

func lessonSlugs(infos []*ent.CourseLessonTranslation) []string {
	return lo.Map(infos, func(i *ent.CourseLessonTranslation, _ int) string { return *i.Edges.Lesson.Slug })
}

// The course page of a signed-in learner: the page the slug names, the lessons
// in course order with the learner's standing, and their enrollment.
func TestCoursePageOfALearner(t *testing.T) {
	r := newRig(t)

	page, err := r.pages.Course(t.Context(), r.learner(t), "en", enLanding)
	require.NoError(t, err)

	assert.Equal(t, jsCourse, *page.Course.Slug)
	assert.Equal(t, enLanding, *page.Landing.Slug)
	assert.Nil(t, page.MovedTo)
	assert.Equal(t, []string{firstLesson, secondLesson, thirdLesson}, lessonSlugs(page.Infos))
	assert.Equal(t, 33, page.State.Completion)
	assert.True(t, page.State.Lessons[0].Finished)
	require.NotNil(t, page.Enrollment)
	assert.Equal(t, page.Course.ID, page.Enrollment.CourseID)
	assert.NotEmpty(t, page.Modules)
}

// A guest has the same page with no enrollment and nothing finished.
func TestCoursePageOfAGuestHasNoEnrollment(t *testing.T) {
	r := newRig(t)

	page, err := r.pages.Course(t.Context(), progress.Learner{}, "en", enLanding)
	require.NoError(t, err)

	assert.Nil(t, page.Enrollment)
	assert.Equal(t, 0, page.State.Completion)
	assert.Len(t, page.Infos, 3)
}

// One ordering: the lessons follow the positions the progress state reports, so
// a lesson moved to the end of the course moves in both.
func TestCoursePageOrdersLessonsByTheProgressPositions(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	first := r.db.CourseLesson.Query().
		Where(courselesson.SlugEQ(firstLesson), courselesson.CourseID(fixtures.ID(fixtures.CourseJavascript))).
		OnlyX(ctx)
	r.db.CourseLessonVersion.Update().
		Where(courselessonversion.LessonID(first.ID)).
		SetNaturalOrder(9).ExecX(ctx)

	page, err := r.pages.Course(ctx, progress.Learner{}, "en", enLanding)
	require.NoError(t, err)

	assert.Equal(t, []string{secondLesson, thirdLesson, firstLesson}, lessonSlugs(page.Infos))
	assert.Equal(t, []string{secondLesson, thirdLesson, firstLesson},
		lo.Map(page.State.Lessons, func(l progress.LessonState, _ int) string { return l.Slug }))
}

func TestCoursePageIsNotFoundForAnUnknownOrForeignSlug(t *testing.T) {
	r := newRig(t)

	_, err := r.pages.Course(t.Context(), progress.Learner{}, "en", "no-such-page")
	require.ErrorIs(t, err, coursepages.ErrNotFound)

	_, err = r.pages.Course(t.Context(), progress.Learner{}, untranslated, enLanding)
	require.ErrorIs(t, err, coursepages.ErrNotFound, "a page of another locale")
}

// An archived page names where it moved: the page it points at, else the
// course's main page in the locale; a live page moves nowhere; an archived page
// with neither is gone.
func TestCoursePageSaysWhereAnArchivedPageMoved(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	live, err := r.pages.Course(ctx, progress.Learner{}, "ru", ruLanding)
	require.NoError(t, err)
	assert.Nil(t, live.MovedTo)

	moved, err := r.pages.Course(ctx, progress.Learner{}, "ru", archivedRu)
	require.NoError(t, err)
	require.NotNil(t, moved.MovedTo)
	assert.Equal(t, ruLanding, *moved.MovedTo, "the main page")

	named := r.db.CourseLandingPage.Query().Where(courselandingpage.SlugEQ("nodejs-ru")).OnlyX(ctx)
	r.db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ(archivedRu)).
		SetLandingPageToRedirectID(named.ID).ExecX(ctx)
	moved, err = r.pages.Course(ctx, progress.Learner{}, "ru", archivedRu)
	require.NoError(t, err)
	assert.Equal(t, "nodejs-ru", *moved.MovedTo, "the page it names wins")

	r.db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ(archivedRu)).
		ClearLandingPageToRedirectID().ExecX(ctx)
	r.db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ(ruLanding)).SetMain(false).ExecX(ctx)
	_, err = r.pages.Course(ctx, progress.Learner{}, "ru", archivedRu)
	require.ErrorIs(t, err, coursepages.ErrNotFound)
}

// The lesson page: prose of the locale, the current version, the source link,
// the course's main landing page, and the lessons in course order.
func TestLessonPage(t *testing.T) {
	r := newRig(t)

	page, err := r.pages.Lesson(t.Context(), r.learner(t), "en", jsCourse, secondLesson)
	require.NoError(t, err)

	assert.Equal(t, secondLesson, *page.Lesson.Slug)
	assert.Equal(t, secondLessonN, *page.Translation.Name)
	assert.Equal(t, *page.Course.CurrentVersionID, page.Version.CourseVersionID)
	assert.Equal(t,
		"https://github.com/hexlet-basics/exercises-javascript/blob/main/modules/10-basics/20-variables/en/README.md",
		page.SourceURL)
	require.NotNil(t, page.Landing)
	assert.Equal(t, enLanding, *page.Landing.Slug)
	assert.Equal(t, []string{firstLesson, secondLesson, thirdLesson}, lessonSlugs(page.Infos))
	assert.Equal(t, 33, page.State.Completion)
}

func TestLessonPageWithoutALandingPageHasNone(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	js := r.db.Course.Query().Where(course.SlugEQ(jsCourse)).OnlyX(ctx)
	r.db.CourseLandingPage.Delete().Where(courselandingpage.CourseID(js.ID)).ExecX(ctx)

	page, err := r.pages.Lesson(ctx, progress.Learner{}, "en", jsCourse, secondLesson)
	require.NoError(t, err)
	assert.Nil(t, page.Landing)
}

// Every way a lesson is missing is the one outcome: unknown course, unknown
// lesson, a lesson of another course, a lesson with no translation in the
// locale, a course with nothing built.
func TestLessonPageIsNotFound(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	for name, args := range map[string][3]string{
		"unknown course":   {"en", "no-such-course", secondLesson},
		"unknown lesson":   {"en", jsCourse, "no-such-lesson"},
		"untranslated":     {untranslated, jsCourse, secondLesson},
		"another course's": {"en", "typescript", "no-such-lesson"},
	} {
		_, err := r.pages.Lesson(ctx, progress.Learner{}, args[0], args[1], args[2])
		require.ErrorIs(t, err, coursepages.ErrNotFound, name)
	}

	r.db.Course.Update().Where(course.SlugEQ(jsCourse)).ClearCurrentVersionID().ExecX(ctx)
	_, err := r.pages.Lesson(ctx, progress.Learner{}, "en", jsCourse, secondLesson)
	require.ErrorIs(t, err, coursepages.ErrNotFound, "nothing built")
}

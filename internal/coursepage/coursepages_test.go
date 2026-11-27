package coursepage_test

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
	"hexletbasics/internal/coursepage"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/landingpage"
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
	pages *coursepage.Pages
}

func newRig(t *testing.T) *rig {
	t.Helper()
	db, _, transactor := testsupport.NewClientWithTxAndTransactor(t)
	tracker := progress.New(db, transactor, testsupport.NewMocks().Runner)

	return &rig{db: db, pages: coursepage.New(db, landingpage.NewCatalog(db), tracker, repoBase)}
}

func (r *rig) learner(t *testing.T) progress.Learner {
	t.Helper()

	return progress.Learner{
		UserID: r.db.User.GetX(t.Context(), fixture.ID(fixture.UserAdminLearner)).ID,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}
}

func lessonSlugs(lessons []coursepage.LearnerLesson) []string {
	return lo.Map(lessons, func(l coursepage.LearnerLesson, _ int) string { return l.Slug })
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
	assert.Equal(t, []string{firstLesson, secondLesson, thirdLesson}, lessonSlugs(page.Lessons))
	require.NotNil(t, page.Progress)
	assert.Equal(t, 33, page.Progress.Completion)
	assert.True(t, page.Lessons[0].Finished)
	require.NotNil(t, page.Progress.NextLesson, "the learner resumes at the first unfinished lesson")
	assert.Equal(t, secondLesson, page.Progress.NextLesson.Slug)
	require.NotNil(t, page.Enrollment)
	assert.Equal(t, page.Course.ID, page.Enrollment.CourseID)
	assert.NotEmpty(t, page.Modules)
}

// A guest has the same page with no enrollment and nothing finished.
func TestCoursePageOfAGuestHasNoEnrollment(t *testing.T) {
	r := newRig(t)

	page, err := r.pages.Course(t.Context(), progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "en", enLanding)
	require.NoError(t, err)

	assert.Nil(t, page.Enrollment)
	require.NotNil(t, page.Progress)
	assert.Equal(t, 0, page.Progress.Completion)
	assert.Len(t, page.Lessons, 3)
}

// One ordering: the lessons follow the positions the progress state reports, so
// a lesson moved to the end of the course moves in both.
func TestCoursePageOrdersLessonsByTheProgressPositions(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	first := r.db.CourseLesson.Query().
		Where(courselesson.SlugEQ(firstLesson), courselesson.CourseID(fixture.ID(fixture.CourseJavascript))).
		OnlyX(ctx)
	r.db.CourseLessonVersion.Update().
		Where(courselessonversion.LessonID(first.ID)).
		SetNaturalOrder(9).ExecX(ctx)

	page, err := r.pages.Course(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "en", enLanding)
	require.NoError(t, err)

	assert.Equal(t, []string{secondLesson, thirdLesson, firstLesson}, lessonSlugs(page.Lessons))
	assert.IsIncreasing(t, lo.Map(page.Lessons, func(l coursepage.LearnerLesson, _ int) int { return l.Position }),
		"each lesson carries the position it is listed by")
}

func TestCoursePageIsNotFoundForAnUnknownOrForeignSlug(t *testing.T) {
	r := newRig(t)

	_, err := r.pages.Course(t.Context(), progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "en", "no-such-page")
	require.ErrorIs(t, err, domain.ErrNotFound)

	_, err = r.pages.Course(t.Context(), progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, untranslated, enLanding)
	require.ErrorIs(t, err, domain.ErrNotFound, "a page of another locale")
}

// An archived page names where it moved: the page it points at, else the
// course's main page in the locale; a live page moves nowhere; an archived page
// with neither is gone.
func TestCoursePageSaysWhereAnArchivedPageMoved(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	live, err := r.pages.Course(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "ru", ruLanding)
	require.NoError(t, err)
	assert.Nil(t, live.MovedTo)

	moved, err := r.pages.Course(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "ru", archivedRu)
	require.NoError(t, err)
	require.NotNil(t, moved.MovedTo)
	assert.Equal(t, ruLanding, *moved.MovedTo, "the main page")

	named := r.db.CourseLandingPage.Query().Where(courselandingpage.SlugEQ("nodejs-ru")).OnlyX(ctx)
	r.db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ(archivedRu)).
		SetLandingPageToRedirectID(named.ID).ExecX(ctx)
	moved, err = r.pages.Course(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "ru", archivedRu)
	require.NoError(t, err)
	assert.Equal(t, "nodejs-ru", *moved.MovedTo, "the page it names wins")

	r.db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ(archivedRu)).
		ClearLandingPageToRedirectID().ExecX(ctx)
	r.db.CourseLandingPage.Update().Where(courselandingpage.SlugEQ(ruLanding)).SetMain(false).ExecX(ctx)
	_, err = r.pages.Course(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "ru", archivedRu)
	require.ErrorIs(t, err, domain.ErrNotFound)
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
	assert.Equal(t, []string{firstLesson, secondLesson, thirdLesson}, lessonSlugs(page.Lessons))
	require.NotNil(t, page.Progress)
	assert.Equal(t, 33, page.Progress.Completion)
	require.NotNil(t, page.Adjacent.Previous)
	require.NotNil(t, page.Adjacent.Next)
	assert.Equal(t, firstLesson, page.Adjacent.Previous.Slug)
	assert.Equal(t, thirdLesson, page.Adjacent.Next.Slug)
}

func TestLessonPageWithoutALandingPageHasNone(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	js := r.db.Course.Query().Where(course.SlugEQ(jsCourse)).OnlyX(ctx)
	r.db.CourseLandingPage.Delete().Where(courselandingpage.CourseID(js.ID)).ExecX(ctx)

	page, err := r.pages.Lesson(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "en", jsCourse, secondLesson)
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
		_, err := r.pages.Lesson(ctx, progress.Learner{
			UserID: 0,
			Guest: progress.GuestProgress{
				Entries: nil,
			},
		}, args[0], args[1], args[2])
		require.ErrorIs(t, err, domain.ErrNotFound, name)
	}

	r.db.Course.Update().Where(course.SlugEQ(jsCourse)).ClearCurrentVersionID().ExecX(ctx)
	_, err := r.pages.Lesson(ctx, progress.Learner{
		UserID: 0,
		Guest: progress.GuestProgress{
			Entries: nil,
		},
	}, "en", jsCourse, secondLesson)
	require.ErrorIs(t, err, domain.ErrNotFound, "nothing built")
}

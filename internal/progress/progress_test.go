package progress_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/ent/user"
	"hexletbasics/internal/events"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// The progress module over the test's own transaction: the real gate, the real
// transitions and the real event publisher, with the exercise runner the one
// mock. Asserted is what a learner or a guest is left with — the rows, the
// state read back, the facts published — never what was called.

const (
	jsSlug       = "javascript"
	firstLesson  = "hello-world"
	secondLesson = "variables"
	thirdLesson  = "strings"
	locale       = "en"
)

type rig struct {
	db       *ent.Client
	tx       *sql.Tx
	services *testsupport.Services
	mocks    *testsupport.Mocks
	progress *progress.Progress
}

func newRig(t *testing.T) *rig {
	t.Helper()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)
	mocks := testsupport.NewMocks()
	return &rig{
		db:       db,
		tx:       tx,
		services: services,
		mocks:    mocks,
		progress: progress.New(db, transactor, services.Events, mocks.Runner),
	}
}

func (r *rig) course(t *testing.T, slug string) *ent.Course {
	t.Helper()
	return r.db.Course.Query().Where(course.SlugEQ(slug)).OnlyX(t.Context())
}

func (r *rig) lesson(t *testing.T, slug string) *ent.CourseLesson {
	t.Helper()
	return r.db.CourseLesson.Query().Where(courselesson.SlugEQ(slug)).OnlyX(t.Context())
}

// learner is a fixture account by email: alice has finished the first lesson,
// bob has the second one started, carol has not begun.
func (r *rig) learner(t *testing.T, email string) progress.Learner {
	t.Helper()
	id := r.db.User.Query().Where(user.EmailEQ(email)).OnlyX(t.Context()).ID
	return progress.Learner{UserID: id}
}

func (r *rig) finishedSlugs(t *testing.T, learner progress.Learner) []string {
	t.Helper()
	rows := r.db.LessonProgress.Query().
		Where(lessonprogress.UserID(learner.UserID), lessonprogress.StateEQ(progress.StateFinished)).
		WithLesson().
		AllX(t.Context())
	slugs := make([]string, 0, len(rows))
	for _, row := range rows {
		slugs = append(slugs, *row.Edges.Lesson.Slug)
	}
	return slugs
}

func guestAt(lessonSlug string) progress.Learner {
	return progress.Learner{Guest: progress.GuestProgress{}.Record(jsSlug, lessonSlug)}
}

func lessonSlugs(state *progress.CourseState) []string {
	slugs := make([]string, 0, len(state.Lessons))
	for _, lesson := range state.Lessons {
		slugs = append(slugs, lesson.Slug)
	}
	return slugs
}

func TestStartLessonEnrollsTheLearnerAndStaysSilentWhenRepeated(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	carol := r.learner(t, "carol@example.com")
	hello := r.lesson(t, firstLesson)

	state, err := r.progress.StartLesson(ctx, carol, hello.ID, locale)
	require.NoError(t, err)
	assert.Equal(t, progress.StateStarted, state.State)
	require.NotNil(t, state.NextLessonSlug)
	assert.Equal(t, firstLesson, *state.NextLessonSlug)

	started := testsupport.PublishedEvents[events.CourseStarted](t, r.services)
	require.Len(t, started, 1)
	assert.Equal(t, jsSlug, started[0].Slug)
	assert.Len(t, testsupport.PublishedEvents[events.LessonStarted](t, r.services), 1)

	// Pressing start again writes nothing and publishes nothing: the counts the
	// CRM keeps count transitions, not requests.
	_, err = r.progress.StartLesson(ctx, carol, hello.ID, locale)
	require.NoError(t, err)
	assert.Equal(t, 2, testsupport.PublishedCount(t, r.services), "one course fact and one lesson fact, once")
	assert.Equal(t, 1, r.db.Enrollment.Query().Where(enrollment.UserID(carol.UserID)).CountX(ctx))
}

func TestStartLessonRefusesALessonBeyondTheGate(t *testing.T) {
	r := newRig(t)
	carol := r.learner(t, "carol@example.com")

	_, err := r.progress.StartLesson(t.Context(), carol, r.lesson(t, thirdLesson).ID, locale)

	assert.ErrorIs(t, err, progress.ErrLessonNotAvailable)
	assert.Zero(t, testsupport.PublishedCount(t, r.services))
}

func TestStartLessonReadsALessonTheCurrentVersionDroppedAsNotFound(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	lesson := r.lesson(t, secondLesson)
	// A build that no longer numbers the lesson has dropped it from the course.
	r.db.CourseLessonVersion.Update().
		Where(courselessonversion.LessonID(lesson.ID)).
		ClearNaturalOrder().
		ExecX(ctx)

	_, err := r.progress.StartLesson(ctx, r.learner(t, "alice@example.com"), lesson.ID, locale)

	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

func TestStartLessonReadsALessonWithNoCourseAsNotFound(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	orphan := r.db.CourseLesson.Create().SetSlug("orphan").SaveX(ctx)

	_, err := r.progress.StartLesson(ctx, r.learner(t, "alice@example.com"), orphan.ID, locale)

	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

func TestStartLessonReadsAnUnknownLessonAsNotFound(t *testing.T) {
	r := newRig(t)

	_, err := r.progress.StartLesson(t.Context(), r.learner(t, "alice@example.com"), -1, locale)

	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

func TestStartLessonLetsAGuestInAtTheGateAndWritesNothing(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	before := r.db.LessonProgress.Query().CountX(ctx)

	state, err := r.progress.StartLesson(ctx, guestAt(firstLesson), r.lesson(t, secondLesson).ID, locale)
	require.NoError(t, err)

	assert.Equal(t, 1, state.FurthestFinishedPosition)
	assert.Equal(t, before, r.db.LessonProgress.Query().CountX(ctx), "a guest leaves no rows")
	assert.Zero(t, testsupport.PublishedCount(t, r.services), "and no facts: there is no account")

	_, err = r.progress.StartLesson(ctx, progress.Learner{}, r.lesson(t, secondLesson).ID, locale)
	assert.ErrorIs(t, err, progress.ErrLessonNotAvailable, "a first-time visitor cannot skip ahead")
}

func TestGuestCourseStateIsDerivedFromTheFurthestLesson(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	courseID := r.course(t, jsSlug).ID

	cases := []struct {
		name      string
		learner   progress.Learner
		state     string
		completed int
		next      string
	}{
		{"a first visit has not started", progress.Learner{}, "", 0, firstLesson},
		{"a visitor part way is started", guestAt(secondLesson), progress.StateStarted, 66, thirdLesson},
		{"a visitor at the last lesson has finished", guestAt(thirdLesson), progress.StateFinished, 100, ""},
		{"a slug the version dropped counts as a first visit", guestAt("gone"), "", 0, firstLesson},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			state, err := r.progress.CourseState(ctx, c.learner, courseID)
			require.NoError(t, err)

			assert.Equal(t, c.state, state.State)
			assert.Equal(t, c.completed, state.Completion)
			assert.Equal(t, []string{firstLesson, secondLesson, thirdLesson}, lessonSlugs(state))
			if c.next == "" {
				assert.Nil(t, state.NextLessonSlug)
			} else {
				require.NotNil(t, state.NextLessonSlug)
				assert.Equal(t, c.next, *state.NextLessonSlug)
			}
		})
	}
}

func TestCourseStateOfAnUnknownCourseIsNotFound(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	_, err := r.progress.CourseState(ctx, progress.Learner{}, -1)
	assert.True(t, ent.IsNotFound(err), "a guest: %v", err)

	_, err = r.progress.CourseState(ctx, r.learner(t, "alice@example.com"), -1)
	assert.True(t, ent.IsNotFound(err), "an account: %v", err)
}

func TestCourseStateOfACourseWithNoVersionHasNoLessons(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	bare := r.db.Course.Create().SetSlug("bare").SaveX(ctx)

	guest, err := r.progress.CourseState(ctx, progress.Learner{}, bare.ID)
	require.NoError(t, err)
	assert.Empty(t, guest.Lessons)
	assert.Zero(t, guest.Completion)

	account, err := r.progress.CourseState(ctx, r.learner(t, "alice@example.com"), bare.ID)
	require.NoError(t, err)
	assert.Empty(t, account.Lessons)
}

func TestCourseStatesReadsEachLearnersStandingInOnePass(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	javascript := r.course(t, jsSlug)
	ruby := r.course(t, "ruby")

	states, err := r.progress.CourseStates(ctx, r.learner(t, "alice@example.com").UserID, []int{javascript.ID, ruby.ID, -1})
	require.NoError(t, err)

	assert.NotContains(t, states, -1, "a course that does not exist is absent")
	js := states[javascript.ID]
	require.NotNil(t, js)
	assert.Equal(t, progress.StateStarted, js.State)
	assert.Equal(t, 33, js.Completion)
	assert.Equal(t, 1, js.FurthestFinishedPosition)
	available := make([]bool, 0, len(js.Lessons))
	for _, lesson := range js.Lessons {
		available = append(available, lesson.Available)
	}
	assert.Equal(t, []bool{true, true, false}, available, "one past the furthest finished lesson is open")
	assert.Empty(t, states[ruby.ID].State, "a course the learner never joined has no state")
}

func TestCourseStatesSkipsLessonsWithoutAPosition(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	lesson := r.lesson(t, thirdLesson)
	r.db.CourseLessonVersion.Update().
		Where(courselessonversion.LessonID(lesson.ID)).
		ClearNaturalOrder().
		ExecX(ctx)

	states, err := r.progress.CourseStates(ctx, r.learner(t, "alice@example.com").UserID, []int{r.course(t, jsSlug).ID})
	require.NoError(t, err)

	state := states[r.course(t, jsSlug).ID]
	assert.Equal(t, []string{firstLesson, secondLesson}, lessonSlugs(state), "an unnumbered lesson is not part of the course order")
}

func TestCourseStatesWithNoCoursesIsEmpty(t *testing.T) {
	r := newRig(t)

	states, err := r.progress.CourseStates(t.Context(), 1, []int{})

	require.NoError(t, err)
	assert.Empty(t, states)
}

// Check.

func (r *rig) check(t *testing.T, learner progress.Learner, lessonSlug string) progress.Check {
	t.Helper()
	ctx := t.Context()
	lesson := r.lesson(t, lessonSlug)
	current := *r.course(t, jsSlug).CurrentVersionID
	version := r.db.CourseLessonVersion.Query().
		Where(courselessonversion.LessonID(lesson.ID), courselessonversion.CourseVersionID(current)).
		OnlyX(ctx)
	return progress.Check{Learner: learner, LessonID: lesson.ID, VersionID: version.ID, Code: "x", Locale: locale}
}

func TestCheckPassesThroughARunnerFailureWithoutRecordingAnything(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	boom := errors.New("docker is down")
	r.mocks.Runner.RunFunc = func(_ context.Context, _ progress.Submission) (progress.Outcome, error) {
		return progress.Outcome{}, boom
	}
	carol := r.learner(t, "carol@example.com")

	_, err := r.progress.CheckSolution(ctx, r.check(t, carol, firstLesson))

	assert.ErrorIs(t, err, boom)
	assert.Zero(t, testsupport.PublishedCount(t, r.services))
	assert.Zero(t, r.db.Enrollment.Query().Where(enrollment.UserID(carol.UserID)).CountX(ctx))
}

func TestCheckRefusesAVersionThatIsNotTheCurrentOneOfThisLesson(t *testing.T) {
	r := newRig(t)
	check := r.check(t, r.learner(t, "alice@example.com"), secondLesson)
	check.VersionID = r.check(t, r.learner(t, "alice@example.com"), firstLesson).VersionID

	_, err := r.progress.CheckSolution(t.Context(), check)

	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

func TestCheckRefusesACourseWithNoCurrentVersion(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	check := r.check(t, r.learner(t, "alice@example.com"), secondLesson)
	r.db.Course.Update().Where(course.SlugEQ(jsSlug)).ClearCurrentVersionID().ExecX(ctx)

	_, err := r.progress.CheckSolution(ctx, check)

	assert.True(t, ent.IsNotFound(err), "a course with no built content has no lesson to check: %v", err)
}

func TestCheckCannotGradeAVersionBuiltWithoutAnImage(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	check := r.check(t, r.learner(t, "alice@example.com"), secondLesson)
	r.db.CourseVersion.UpdateOneID(*r.course(t, jsSlug).CurrentVersionID).
		ClearDockerImage().
		ExecX(ctx)

	_, err := r.progress.CheckSolution(ctx, check)

	assert.True(t, ent.IsNotFound(err), "an ungradeable lesson reads as not found: %v", err)
}

func TestUnavailableRunnerFailsLoudly(t *testing.T) {
	_, err := progress.UnavailableRunner{}.Run(t.Context(), progress.Submission{})

	assert.ErrorIs(t, err, progress.ErrRunnerUnavailable)
}

func TestSubmissionWithImageChangesOnlyTheImage(t *testing.T) {
	original := progress.Submission{Image: "repo", VersionID: 7, Code: "x"}

	pinned := original.WithImage("repo:7")

	assert.Equal(t, "repo:7", pinned.Image)
	assert.Equal(t, 7, pinned.VersionID)
	assert.Equal(t, "x", pinned.Code)
	assert.Equal(t, "repo", original.Image, "the original is untouched")
}

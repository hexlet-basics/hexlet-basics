package courseloader_test

import (
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// recordingPublisher keeps the facts a promotion emitted, without touching the
// outbox — Watermill's write is production's job, not this test's subject.
type recordingPublisher struct {
	published []events.Event
}

func (p *recordingPublisher) Publish(_ context.Context, _ *sql.Tx, event events.Event) error {
	p.published = append(p.published, event)
	return nil
}

func (p *recordingPublisher) finishedCourses() []events.CourseFinished {
	var found []events.CourseFinished
	for _, event := range p.published {
		if finished, ok := event.(events.CourseFinished); ok {
			found = append(found, finished)
		}
	}
	return found
}

// A learner who has finished every lesson the newly promoted version contains
// is finished by the promotion itself. Completion is otherwise only evaluated
// on a submission, and this learner has nothing left to submit.
func TestPromotionFinishesALearnerWhoHasCompletedTheNewVersion(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	publisher := &recordingPublisher{}

	course := db.Course.Create().SetSlug("loader-completion-lang").SetName("Completion").SaveX(ctx)
	loader := newLoaderWithCompletion(t, db, txStore, fakeFetcher{dir: fixtureRepo(t)}, publisher)

	// A first build gives the course its lesson, which the learner finishes.
	first := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, first.ID))

	learner := enrollLearner(t, db, course.ID, seedLearner(t, db, "completion@example.com"))
	lesson := db.CourseLesson.Query().Where(courselesson.CourseID(course.ID)).OnlyX(ctx)
	finishLessonFor(t, db, learner, lesson)

	// Promoting a version with the same single lesson completes them.
	second := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, second.ID))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)

	finished := publisher.finishedCourses()
	require.Len(t, finished, 1)
	assert.Equal(t, "loader-completion-lang", finished[0].Slug)
	assert.Equal(t, 1, finished[0].OccurrenceCount, "the learner's finished enrollments")
}

// A learner with an unfinished lesson in the new version stays in progress: the
// re-evaluation only ever moves an enrollment forward.
func TestPromotionLeavesAnUnfinishedLearnerInProgress(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	publisher := &recordingPublisher{}

	course := db.Course.Create().SetSlug("loader-incomplete-lang").SetName("Incomplete").SaveX(ctx)
	loader := newLoaderWithCompletion(t, db, txStore, fakeFetcher{dir: fixtureRepo(t)}, publisher)

	first := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, first.ID))

	// Enrolled, nothing finished.
	learner := enrollLearner(t, db, course.ID, seedLearner(t, db, "incomplete@example.com"))

	second := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, second.ID))

	assert.Equal(t, "started", *db.Enrollment.GetX(ctx, learner.ID).State)
	assert.Empty(t, publisher.finishedCourses())
}

// A learner already finished is not selected again, so the fact is not
// published twice by a later build.
func TestPromotionDoesNotRepublishForAnAlreadyFinishedLearner(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	publisher := &recordingPublisher{}

	course := db.Course.Create().SetSlug("loader-refinish-lang").SetName("Refinish").SaveX(ctx)
	loader := newLoaderWithCompletion(t, db, txStore, fakeFetcher{dir: fixtureRepo(t)}, publisher)

	first := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, first.ID))

	learner := enrollLearner(t, db, course.ID, seedLearner(t, db, "refinish@example.com"))
	lesson := db.CourseLesson.Query().Where(courselesson.CourseID(course.ID)).OnlyX(ctx)
	finishLessonFor(t, db, learner, lesson)

	second := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, second.ID))
	require.Len(t, publisher.finishedCourses(), 1)

	third := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, third.ID))

	assert.Len(t, publisher.finishedCourses(), 1, "the second promotion finds nothing to finish")
}

// seedLearner creates the account the enrollment's FK requires.
func seedLearner(t *testing.T, db *ent.Client, email string) int {
	t.Helper()
	return db.User.Create().SetEmail(email).SaveX(t.Context()).ID
}

func enrollLearner(t *testing.T, db *ent.Client, courseID, userID int) *ent.Enrollment {
	t.Helper()
	return db.Enrollment.Create().
		SetUserID(userID).
		SetCourseID(courseID).
		SetState("started").
		SaveX(t.Context())
}

func finishLessonFor(t *testing.T, db *ent.Client, enrolled *ent.Enrollment, lesson *ent.CourseLesson) {
	t.Helper()
	db.LessonProgress.Create().
		SetUserID(enrolled.UserID).
		SetCourseID(enrolled.CourseID).
		SetEnrollmentID(enrolled.ID).
		SetLessonID(lesson.ID).
		SetState("finished").
		SaveX(t.Context())
}

// The scenario the ticket exists for: a build that REMOVES a lesson completes
// the learners who had finished everything that survives. They have nothing
// left to submit, so a submission — the only other thing that evaluates
// completion — can never notice.
func TestPromotionFinishesLearnersWhenAVersionDropsALesson(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	publisher := &recordingPublisher{}

	course := db.Course.Create().SetSlug("loader-shrink-lang").SetName("Shrink").SaveX(ctx)

	// First build carries two lessons; the learner finishes only the first.
	wide := newLoaderWithCompletion(t, db, txStore, fakeFetcher{dir: twoLessonRepo(t)}, publisher)
	first := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, wide.Run(ctx, first.ID))

	learner := enrollLearner(t, db, course.ID, seedLearner(t, db, "shrink@example.com"))
	kept := db.CourseLesson.Query().
		Where(courselesson.CourseID(course.ID), courselesson.SlugEQ("hello-world")).
		OnlyX(ctx)
	finishLessonFor(t, db, learner, kept)

	require.Equal(t, "started", *db.Enrollment.GetX(ctx, learner.ID).State,
		"still in progress while the second lesson exists")

	// The next build drops the second lesson, which completes them.
	narrow := newLoaderWithCompletion(t, db, txStore, fakeFetcher{dir: fixtureRepo(t)}, publisher)
	second := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, narrow.Run(ctx, second.ID))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
	require.Len(t, publisher.finishedCourses(), 1)
	assert.Equal(t, "loader-shrink-lang", publisher.finishedCourses()[0].Slug)
}

// A learner whose state column is NULL is started — the legacy baseline allows
// NULL and rows predating AASM's writes exist. They must be able to complete.
func TestPromotionFinishesALearnerWhoseStateIsNull(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	publisher := &recordingPublisher{}

	course := db.Course.Create().SetSlug("loader-nullstate-lang").SetName("Null").SaveX(ctx)
	loader := newLoaderWithCompletion(t, db, txStore, fakeFetcher{dir: fixtureRepo(t)}, publisher)

	first := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, first.ID))

	learner := enrollLearner(t, db, course.ID, seedLearner(t, db, "nullstate@example.com"))
	db.Enrollment.UpdateOneID(learner.ID).ClearState().ExecX(ctx)
	lesson := db.CourseLesson.Query().Where(courselesson.CourseID(course.ID)).OnlyX(ctx)
	finishLessonFor(t, db, learner, lesson)

	second := db.CourseVersion.Create().SetCourseID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, second.ID))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
	assert.Len(t, publisher.finishedCourses(), 1)
}

// twoLessonRepo is the committed fixture course with its single lesson
// duplicated under a second slug, so a later build of the unmodified fixture
// removes one.
func twoLessonRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, copyTree(fixtureRepo(t), dir))

	basics := filepath.Join(dir, "modules", "10-basics")
	require.NoError(t, copyTree(
		filepath.Join(basics, "10-hello-world"),
		filepath.Join(basics, "20-second-lesson"),
	))
	return dir
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, body, 0o644)
	})
}

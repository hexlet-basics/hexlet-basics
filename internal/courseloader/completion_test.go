package courseloader_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/events"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// A learner who has finished every lesson the newly promoted version contains
// is finished by the promotion itself. Completion is otherwise only evaluated
// on a submission, and this learner has nothing left to submit.
func TestPromotionFinishesALearnerWhoHasCompletedTheNewVersion(t *testing.T) {
	db, tx, txStore := testsupport.NewClientWithTxAndTransactor(t)
	ctx := context.Background()
	services := testsupport.NewServices(t, tx, txStore)
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)), services.Events)
	// The rebuildable course was built once and has its only lesson finished by
	// user_learner_finished_rebuildable.
	learner := db.Enrollment.GetX(ctx, fixtures.ID(fixtures.EnrollmentLearnerFinishedRebuildableRebuildable))

	// Promoting a version with the same single lesson completes them.
	require.NoError(t, loader.Run(ctx, fixtures.ID(fixtures.VersionRebuildableCreatedSecond)))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)

	finished := testsupport.PublishedEvents[events.CourseFinished](t, services)
	require.Len(t, finished, 1)
	assert.Equal(t, "rebuildable", finished[0].Slug)
	assert.Equal(t, 1, finished[0].OccurrenceCount, "the learner's finished enrollments")
}

// A learner with an unfinished lesson in the new version stays in progress: the
// re-evaluation only ever moves an enrollment forward.
func TestPromotionLeavesAnUnfinishedLearnerInProgress(t *testing.T) {
	db, tx, txStore := testsupport.NewClientWithTxAndTransactor(t)
	ctx := context.Background()
	services := testsupport.NewServices(t, tx, txStore)
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)), services.Events)
	// Enrolled in the rebuildable course, nothing finished.
	learner := db.Enrollment.GetX(ctx, fixtures.ID(fixtures.EnrollmentLearnerUnfinishedRebuildableRebuildable))

	require.NoError(t, loader.Run(ctx, fixtures.ID(fixtures.VersionRebuildableCreatedSecond)))

	assert.Equal(t, "started", *db.Enrollment.GetX(ctx, learner.ID).State)
	finished := testsupport.PublishedEvents[events.CourseFinished](t, services)
	assert.Len(t, finished, 1, "only the learner who had finished the lesson is closed")
}

// A learner already finished is not selected again, so the fact is not
// published twice by a later build.
func TestPromotionDoesNotRepublishForAnAlreadyFinishedLearner(t *testing.T) {
	db, tx, txStore := testsupport.NewClientWithTxAndTransactor(t)
	ctx := context.Background()
	services := testsupport.NewServices(t, tx, txStore)
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)), services.Events)

	require.NoError(t, loader.Run(ctx, fixtures.ID(fixtures.VersionRebuildableCreatedSecond)))
	require.Len(t, testsupport.PublishedEvents[events.CourseFinished](t, services), 1)

	require.NoError(t, loader.Run(ctx, fixtures.ID(fixtures.VersionRebuildableCreatedThird)))

	assert.Len(t, testsupport.PublishedEvents[events.CourseFinished](t, services), 1, "the second promotion finds nothing to finish")
}

// The scenario the ticket exists for: a build that REMOVES a lesson completes
// the learners who had finished everything that survives. They have nothing
// left to submit, so a submission — the only other thing that evaluates
// completion — can never notice.
func TestPromotionFinishesLearnersWhenAVersionDropsALesson(t *testing.T) {
	db, tx, txStore := testsupport.NewClientWithTxAndTransactor(t)
	ctx := context.Background()
	services := testsupport.NewServices(t, tx, txStore)
	// The shrinkable course was built with two lessons; the learner finished only
	// the first.
	learner := db.Enrollment.GetX(ctx, fixtures.ID(fixtures.EnrollmentLearnerFinishedKeptLessonShrinkable))
	require.Equal(t, "started", *learner.State, "still in progress while the second lesson exists")

	// The next build is the committed fixture course: one lesson, the second is dropped.
	narrow := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)), services.Events)
	require.NoError(t, narrow.Run(ctx, fixtures.ID(fixtures.VersionShrinkableCreatedNext)))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
	finished := testsupport.PublishedEvents[events.CourseFinished](t, services)
	require.Len(t, finished, 1)
	assert.Equal(t, "shrinkable", finished[0].Slug)
}

// A learner whose state column is NULL is started — the legacy baseline allows
// NULL and rows predating AASM's writes exist. They must be able to complete.
func TestPromotionFinishesALearnerWhoseStateIsNull(t *testing.T) {
	db, tx, txStore := testsupport.NewClientWithTxAndTransactor(t)
	ctx := context.Background()
	services := testsupport.NewServices(t, tx, txStore)
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)), services.Events)
	learner := db.Enrollment.GetX(ctx, fixtures.ID(fixtures.EnrollmentLearnerNullStateNullState))
	require.Nil(t, learner.State, "the fixture enrollment has no state")

	require.NoError(t, loader.Run(ctx, fixtures.ID(fixtures.VersionNullStateCreatedSecond)))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
	assert.Len(t, testsupport.PublishedEvents[events.CourseFinished](t, services), 1)
}

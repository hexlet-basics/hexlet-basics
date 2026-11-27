package courseloader_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// A learner who has finished every lesson the newly promoted version contains
// is finished by the promotion itself. Completion is otherwise only evaluated
// on a submission, and this learner has nothing left to submit.
func TestPromotionFinishesALearnerWhoHasCompletedTheNewVersion(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)))
	// The rebuildable course was built once and has its only lesson finished by
	// user_learner_finished_rebuildable.
	learner := db.Enrollment.GetX(ctx, fixture.ID(fixture.EnrollmentLearnerFinishedRebuildableRebuildable))

	// Promoting a version with the same single lesson completes them.
	require.NoError(t, loader.Run(ctx, fixture.ID(fixture.VersionRebuildableCreatedSecond)))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
}

// A learner with an unfinished lesson in the new version stays in progress: the
// re-evaluation only ever moves an enrollment forward.
func TestPromotionLeavesAnUnfinishedLearnerInProgress(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)))
	// Enrolled in the rebuildable course, nothing finished.
	learner := db.Enrollment.GetX(ctx, fixture.ID(fixture.EnrollmentLearnerUnfinishedRebuildableRebuildable))

	require.NoError(t, loader.Run(ctx, fixture.ID(fixture.VersionRebuildableCreatedSecond)))

	assert.Equal(t, "started", *db.Enrollment.GetX(ctx, learner.ID).State)
	finisher := fixture.ID(fixture.EnrollmentLearnerFinishedRebuildableRebuildable)
	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, finisher).State, "only the learner who had finished the lesson is closed")
}

// The scenario the ticket exists for: a build that REMOVES a lesson completes
// the learners who had finished everything that survives. They have nothing
// left to submit, so a submission — the only other thing that evaluates
// completion — can never notice.
func TestPromotionFinishesLearnersWhenAVersionDropsALesson(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	// The shrinkable course was built with two lessons; the learner finished only
	// the first.
	learner := db.Enrollment.GetX(ctx, fixture.ID(fixture.EnrollmentLearnerFinishedKeptLessonShrinkable))
	require.Equal(t, "started", *learner.State, "still in progress while the second lesson exists")

	// The next build is the committed fixture course: one lesson, the second is dropped.
	narrow := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)))
	require.NoError(t, narrow.Run(ctx, fixture.ID(fixture.VersionShrinkableCreatedNext)))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
}

// A learner whose state column is NULL is started — the legacy baseline allows
// NULL and rows predating AASM's writes exist. They must be able to complete.
func TestPromotionFinishesALearnerWhoseStateIsNull(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	loader := newLoaderWithCompletion(t, db, txStore, repoFetcher(fixtureRepo(t)))
	learner := db.Enrollment.GetX(ctx, fixture.ID(fixture.EnrollmentLearnerNullStateNullState))
	require.Nil(t, learner.State, "the fixture enrollment has no state")

	require.NoError(t, loader.Run(ctx, fixture.ID(fixture.VersionNullStateCreatedSecond)))

	assert.Equal(t, "finished", *db.Enrollment.GetX(ctx, learner.ID).State)
}

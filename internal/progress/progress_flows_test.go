package progress_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"slices"
	"testing"

	entgo "entgo.io/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/progress"
)

// Recording a pass.

func TestCheckByALearnerWhoseCourseIsAlreadyFinishedDoesNotFinishItTwice(t *testing.T) {
	r := newRig(t)
	staffEditor := r.learner(t, fixture.UserStaffEditor) // enrollment finished, second lesson only started

	result, err := r.progress.CheckSolution(t.Context(), r.check(t, staffEditor, firstLesson))
	require.NoError(t, err)

	assert.True(t, result.Outcome.Passed)
	assert.True(t, result.LessonFinished)
	assert.False(t, result.CourseFinished, "the enrollment was closed before")
}

// Merging a visitor's cookie into an account.

func TestMergeGuestCreditsEveryFinishedLessonToANewAccount(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, fixture.UserBlank)

	require.NoError(t, r.progress.MergeGuest(ctx, blankUser.UserID, guestWith(firstLesson, secondLesson, thirdLesson).Guest))

	assert.ElementsMatch(t, []string{firstLesson, secondLesson, thirdLesson}, r.finishedSlugs(t, blankUser))
	state, err := r.progress.CourseState(ctx, blankUser, r.course(t, jsSlug).ID)
	require.NoError(t, err)
	assert.Equal(t, 100, state.Completion)
	assert.Equal(t, 1, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID)).CountX(ctx), "the merge enrolls them")
}

// A visitor who took the third lesson first is credited that lesson and nothing
// else: with lessons in any order, a later one no longer implies the earlier.
func TestMergeGuestCreditsOnlyTheLessonsTheVisitorFinished(t *testing.T) {
	r := newRig(t)
	blankUser := r.learner(t, fixture.UserBlank)

	require.NoError(t, r.progress.MergeGuest(t.Context(), blankUser.UserID, guestWith(thirdLesson).Guest))

	assert.Equal(t, []string{thirdLesson}, r.finishedSlugs(t, blankUser))
}

func TestMergeGuestCreditsOnlyWhatTheAccountLacks(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, fixture.UserAdminLearner) // has finished the first lesson

	require.NoError(t, r.progress.MergeGuest(ctx, adminLearner.UserID, guestWith(firstLesson, secondLesson).Guest))

	assert.ElementsMatch(t, []string{firstLesson, secondLesson}, r.finishedSlugs(t, adminLearner))
	assert.Equal(t, 1, r.db.Enrollment.Query().
		Where(enrollment.UserID(adminLearner.UserID), enrollment.HasCourseWith(course.SlugEQ(jsSlug))).
		CountX(ctx), "she was already enrolled")
}

func TestMergeGuestFinishesALessonTheAccountOnlyStarted(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	staffEditor := r.learner(t, fixture.UserStaffEditor) // second lesson started, nothing finished

	require.NoError(t, r.progress.MergeGuest(ctx, staffEditor.UserID, guestWith(firstLesson, secondLesson).Guest))

	assert.ElementsMatch(t, []string{firstLesson, secondLesson}, r.finishedSlugs(t, staffEditor))
	assert.Equal(t, 2, r.db.LessonProgress.Query().Where(lessonprogress.UserID(staffEditor.UserID)).CountX(ctx),
		"the started row was finished in place, not duplicated")
}

func TestMergeGuestNeverMovesAnAccountBackwards(t *testing.T) {
	r := newRig(t)
	adminLearner := r.learner(t, fixture.UserAdminLearner)

	require.NoError(t, r.progress.MergeGuest(t.Context(), adminLearner.UserID, guestWith(firstLesson).Guest))

	assert.Equal(t, []string{firstLesson}, r.finishedSlugs(t, adminLearner))
}

func TestMergeGuestIgnoresWhatItCannotResolve(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, fixture.UserBlank)
	stale := progress.GuestProgress{
		Entries: nil,
	}.
		Record("no-such-course", firstLesson).
		Record(jsSlug, "a-lesson-the-build-dropped")

	require.NoError(t, r.progress.MergeGuest(ctx, blankUser.UserID, stale))
	require.NoError(t, r.progress.MergeGuest(ctx, blankUser.UserID, progress.GuestProgress{
		Entries: nil,
	}))

	assert.Empty(t, r.finishedSlugs(t, blankUser))
	assert.Zero(t, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID)).CountX(ctx))
}

// Re-evaluating completion when a Version is promoted.

func (r *rig) dropLessonsAfterTheFirst(t *testing.T) {
	t.Helper()

	for _, slug := range []string{secondLesson, thirdLesson} {
		r.db.CourseLessonVersion.Update().
			Where(courselessonversion.LessonID(r.lesson(t, slug).ID)).
			ClearNaturalOrder().
			ExecX(t.Context())
	}
}

func TestReevaluateCompletionFinishesTheCourseOfALearnerWithNothingLeftToDo(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	// Every lesson finished, the enrollment still open: the only fixture learner
	// the promotion has anything to close for.
	allDone := r.learner(t, fixture.UserLearnerAllLessonsDone)
	courseID := r.course(t, jsSlug).ID
	finishedBefore := r.finishedEnrollments(t, courseID)

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.db, courseID))

	state, err := r.progress.CourseState(ctx, allDone, courseID)
	require.NoError(t, err)
	assert.Equal(t, progress.StateFinished, state.State)
	assert.Equal(t, finishedBefore+1, r.finishedEnrollments(t, courseID), "nobody else has nothing left to do")
}

// finishedEnrollments counts the Course's finished Enrollments.
func (r *rig) finishedEnrollments(t *testing.T, courseID int) int {
	t.Helper()

	return r.db.Enrollment.Query().
		Where(enrollment.CourseID(courseID), enrollment.StateEQ(progress.StateFinished)).
		CountX(t.Context())
}

// A new version that drops the lessons a learner had left leaves them with
// nothing to do, so the promotion finishes their course.
func TestReevaluateCompletionFinishesTheCourseOfALearnerWhoseRemainingLessonsWereDropped(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, fixture.UserAdminLearner)
	courseID := r.course(t, jsSlug).ID
	r.dropLessonsAfterTheFirst(t) // adminLearner finished the only lesson that is left

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.db, courseID))

	state, err := r.progress.CourseState(ctx, adminLearner, courseID)
	require.NoError(t, err)
	assert.Equal(t, progress.StateFinished, state.State)
}

func TestReevaluateCompletionLeavesLearnersWithWorkLeftAlone(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.db, r.course(t, jsSlug).ID))

	// Work left: the second lesson, or a lesson not even started.
	for _, label := range []fixture.Label{fixture.UserAdminLearner, fixture.UserLearnerWithAGap, fixture.UserLearnerHiddenCards} {
		state, err := r.progress.CourseState(ctx, r.learner(t, label), r.course(t, jsSlug).ID)
		require.NoError(t, err)
		assert.Equal(t, progress.StateStarted, state.State, label.String())
	}
}

func TestReevaluateCompletionCompletesNobodyForAVersionWithNoLessons(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	bare := r.course(t, "bare")
	finishedBefore := r.finishedEnrollments(t, bare.ID)

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.db, bare.ID))

	assert.Equal(t, finishedBefore, r.finishedEnrollments(t, bare.ID))
}

func TestReevaluateCompletionRefusesAnUnknownCourse(t *testing.T) {
	r := newRig(t)

	err := r.progress.ReevaluateCompletion(t.Context(), r.db, -1)

	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

// Two requests at once resolve to one row. The race is staged by letting the
// "other request" land its insert just before this one's.

func TestStartLessonSharesAnEnrollmentWonByAConcurrentRequest(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, fixture.UserBlank)
	courseID := r.course(t, jsSlug).ID
	membersBefore := r.course(t, jsSlug).EnrollmentsCount
	raced := false

	r.db.Enrollment.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mod ent.Mutation) (ent.Value, error) {
			if !raced && mod.Op().Is(ent.OpCreate) {
				raced = true

				r.db.Enrollment.Create(). //nolint:forbidigo // the competing request that wins the race is the act
								SetUserID(blankUser.UserID).SetCourseID(courseID).SetState(progress.StateStarted).ExecX(ctx)
			}

			return next.Mutate(ctx, mod)
		})
	})

	_, err := r.progress.StartLesson(ctx, blankUser, r.lesson(t, firstLesson).ID)
	require.NoError(t, err)

	assert.Equal(t, 1, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID)).CountX(ctx))
	assert.Equal(t, membersBefore, r.course(t, jsSlug).EnrollmentsCount,
		"only the request whose insert landed counts the enrollment")
}

func TestStartLessonSharesLessonProgressWonByAConcurrentRequest(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, fixture.UserAdminLearner)
	second := r.lesson(t, secondLesson)
	courseID := r.course(t, jsSlug).ID
	enrolled := r.db.Enrollment.Query().Where(enrollment.UserID(adminLearner.UserID)).OnlyX(ctx)
	raced := false

	r.db.LessonProgress.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mod ent.Mutation) (ent.Value, error) {
			if !raced && mod.Op().Is(ent.OpCreate) {
				raced = true

				r.db.LessonProgress.Create(). //nolint:forbidigo // the competing request that wins the race is the act
								SetUserID(adminLearner.UserID).SetCourseID(courseID).SetEnrollmentID(enrolled.ID).
								SetLessonID(second.ID).SetState(progress.StateStarted).ExecX(ctx)
			}

			return next.Mutate(ctx, mod)
		})
	})

	_, err := r.progress.StartLesson(ctx, adminLearner, second.ID)
	require.NoError(t, err)

	assert.Equal(t, 1, r.db.LessonProgress.Query().
		Where(lessonprogress.UserID(adminLearner.UserID), lessonprogress.LessonID(second.ID)).CountX(ctx))
}

// The guest cookie.

func signed(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))

	return base64.RawURLEncoding.EncodeToString([]byte(body)) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func TestGuestCookieRejectsASignatureThatIsNotEncoded(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"entries":[]}`))

	_, err := codec.Decode(payload + ".!!!")

	assert.ErrorIs(t, err, progress.ErrGuestCookieInvalid)
}

func TestGuestCookieRejectsAGenuinelySignedBodyThatIsNotProgress(t *testing.T) {
	codec := progress.NewGuestCodec(testSecret)

	_, err := codec.Decode(signed(testSecret, "this is not json"))

	assert.ErrorIs(t, err, progress.ErrGuestCookieInvalid)
}

// A store that fails. Every read and write the module makes must surface its
// failure rather than read as an empty answer or half a transition.

var errStore = errors.New("store is down")

func (r *rig) failReads(entity string) {
	r.db.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
			if entgo.QueryFromContext(ctx).Type == entity {
				return nil, errStore
			}

			return next.Query(ctx, q)
		})
	}))
}

func (r *rig) failWrites(entity string) {
	r.db.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if m.Type() == entity {
				return nil, errStore
			}

			return next.Mutate(ctx, m)
		})
	})
}

type operation struct {
	name string
	// prepare stages whatever data the call needs and returns the call itself, so
	// a fault armed afterwards fails the call and not its set-up.
	prepare func(t *testing.T, r *rig) func() error
	// reads and writes are the entities the call must touch; a failure of any of
	// them has to come back as an error.
	reads, writes []string
}

func operations() []operation {
	return slices.Concat(lessonOperations(), courseOperations())
}

// lessonOperations are the calls that start and check a lesson.
func lessonOperations() []operation {
	return []operation{
		{
			name: "starting a lesson",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				blankUser, lesson := r.learner(t, fixture.UserBlank), r.lesson(t, firstLesson)

				return func() error {
					_, err := r.progress.StartLesson(t.Context(), blankUser, lesson.ID)

					return err
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLesson, ent.TypeCourseLessonVersion, ent.TypeEnrollment, ent.TypeLessonProgress},
			writes: []string{ent.TypeEnrollment, ent.TypeLessonProgress, ent.TypeCourse},
		},
		{
			name: "checking a solution",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				check := r.check(t, r.learner(t, fixture.UserBlank), firstLesson)

				return func() error {
					_, err := r.progress.CheckSolution(t.Context(), check)

					return err
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLesson, ent.TypeCourseLessonVersion, ent.TypeCourseVersion, ent.TypeEnrollment, ent.TypeLessonProgress},
			writes: []string{ent.TypeEnrollment, ent.TypeLessonProgress, ent.TypeCourse},
		},
		{
			name: "checking the last lesson of a course",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				r.dropLessonsAfterTheFirst(t)
				check := r.check(t, r.learner(t, fixture.UserAdminLearner), firstLesson)

				return func() error {
					_, err := r.progress.CheckSolution(t.Context(), check)

					return err
				}
			},
			reads:  []string{ent.TypeLessonProgress, ent.TypeEnrollment},
			writes: []string{ent.TypeEnrollment},
		},
	}
}

// courseOperations are the calls that read or settle a whole course.
func courseOperations() []operation {
	return []operation{
		{
			name: "reading a learner's courses",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				adminLearner, courseID := r.learner(t, fixture.UserAdminLearner), r.course(t, jsSlug).ID

				return func() error {
					_, err := r.progress.CourseStates(t.Context(), adminLearner.UserID, []int{courseID})

					return err
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLessonVersion, ent.TypeLessonProgress, ent.TypeEnrollment},
			writes: nil,
		},
		{
			name: "reading a guest's course",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				courseID := r.course(t, jsSlug).ID

				return func() error {
					_, err := r.progress.CourseState(t.Context(), guestWith(secondLesson), courseID)

					return err
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLessonVersion},
			writes: nil,
		},
		{
			name: "merging a guest's progress",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				staffEditor := r.learner(t, fixture.UserStaffEditor)

				return func() error {
					return r.progress.MergeGuest(t.Context(), staffEditor.UserID, guestWith(thirdLesson).Guest)
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLessonVersion, ent.TypeLessonProgress, ent.TypeEnrollment},
			writes: []string{ent.TypeLessonProgress},
		},
		{
			name: "re-evaluating completion",
			prepare: func(t *testing.T, r *rig) func() error {
				t.Helper()
				r.dropLessonsAfterTheFirst(t)
				courseID := r.course(t, jsSlug).ID

				return func() error {
					return r.progress.ReevaluateCompletion(t.Context(), r.db, courseID)
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLessonVersion, ent.TypeLessonProgress, ent.TypeEnrollment},
			writes: []string{ent.TypeEnrollment},
		},
	}
}

func TestAFailingReadIsNeverAnEmptyAnswer(t *testing.T) {
	for _, op := range operations() {
		for _, entity := range op.reads {
			t.Run(op.name+"/"+entity, func(t *testing.T) {
				r := newRig(t)
				call := op.prepare(t, r)
				r.failReads(entity)

				require.ErrorIs(t, call(), errStore)
			})
		}
	}
}

func TestAFailingWriteLeavesNoPartialTransition(t *testing.T) {
	for _, op := range operations() {
		for _, entity := range op.writes {
			t.Run(op.name+"/"+entity, func(t *testing.T) {
				r := newRig(t)
				call := op.prepare(t, r)
				before := r.progressRows(t)
				r.failWrites(entity)

				require.ErrorIs(t, call(), errStore)
				assert.Equal(t, before, r.progressRows(t), "no row outlives a failed transition")
			})
		}
	}
}

// progressRows is what a transition writes, counted: Enrollments and Lesson
// Progress, all of them and the finished ones.
func (r *rig) progressRows(t *testing.T) [4]int {
	t.Helper()
	ctx := t.Context()

	return [4]int{
		r.db.Enrollment.Query().CountX(ctx),
		r.db.Enrollment.Query().Where(enrollment.StateEQ(progress.StateFinished)).CountX(ctx),
		r.db.LessonProgress.Query().CountX(ctx),
		r.db.LessonProgress.Query().Where(lessonprogress.StateEQ(progress.StateFinished)).CountX(ctx),
	}
}

package progress_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"

	entgo "entgo.io/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/events"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// Recording a pass.

func TestCheckByALearnerWhoseCourseIsAlreadyFinishedDoesNotFinishItTwice(t *testing.T) {
	r := newRig(t)
	staffEditor := r.learner(t, "user_staff_editor") // enrollment finished, second lesson only started

	result, err := r.progress.CheckSolution(t.Context(), r.check(t, staffEditor, firstLesson))
	require.NoError(t, err)

	assert.True(t, result.Outcome.Passed)
	assert.True(t, result.LessonFinished)
	assert.False(t, result.CourseFinished, "the enrollment was closed before")
	assert.Empty(t, testsupport.PublishedEvents[events.CourseFinished](t, r.services))
}

// Merging a visitor's cookie into an account.

func TestMergeGuestCreditsEveryFinishedLessonToANewAccount(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, "user_blank")

	require.NoError(t, r.progress.MergeGuest(ctx, blankUser.UserID, guestWith(firstLesson, secondLesson, thirdLesson).Guest, locale))

	assert.ElementsMatch(t, []string{firstLesson, secondLesson, thirdLesson}, r.finishedSlugs(t, blankUser))
	state, err := r.progress.CourseState(ctx, blankUser, r.course(t, jsSlug).ID)
	require.NoError(t, err)
	assert.Equal(t, 100, state.Completion)

	assert.Len(t, testsupport.PublishedEvents[events.CourseStarted](t, r.services), 1)
	assert.Len(t, testsupport.PublishedEvents[events.LessonStarted](t, r.services), 3)
	finished := testsupport.PublishedEvents[events.LessonFinished](t, r.services)
	require.Len(t, finished, 3)
	assert.Equal(t, []int{1, 2, 3}, []int{finished[0].OccurrenceCount, finished[1].OccurrenceCount, finished[2].OccurrenceCount},
		"each fact counts the lesson rows as of that lesson")
}

// A visitor who took the third lesson first is credited that lesson and nothing
// else: with lessons in any order, a later one no longer implies the earlier.
func TestMergeGuestCreditsOnlyTheLessonsTheVisitorFinished(t *testing.T) {
	r := newRig(t)
	blankUser := r.learner(t, "user_blank")

	require.NoError(t, r.progress.MergeGuest(t.Context(), blankUser.UserID, guestWith(thirdLesson).Guest, locale))

	assert.Equal(t, []string{thirdLesson}, r.finishedSlugs(t, blankUser))
	assert.Len(t, testsupport.PublishedEvents[events.LessonFinished](t, r.services), 1)
}

func TestMergeGuestCreditsOnlyWhatTheAccountLacks(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, "user_admin_learner") // has finished the first lesson

	require.NoError(t, r.progress.MergeGuest(ctx, adminLearner.UserID, guestWith(firstLesson, secondLesson).Guest, locale))

	assert.ElementsMatch(t, []string{firstLesson, secondLesson}, r.finishedSlugs(t, adminLearner))
	assert.Empty(t, testsupport.PublishedEvents[events.CourseStarted](t, r.services), "she was already enrolled")
	started := testsupport.PublishedEvents[events.LessonStarted](t, r.services)
	require.Len(t, started, 1)
	assert.Equal(t, secondLesson, started[0].LessonSlug)
}

func TestMergeGuestFinishesALessonTheAccountOnlyStarted(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	staffEditor := r.learner(t, "user_staff_editor") // second lesson started, nothing finished

	require.NoError(t, r.progress.MergeGuest(ctx, staffEditor.UserID, guestWith(firstLesson, secondLesson).Guest, locale))

	assert.ElementsMatch(t, []string{firstLesson, secondLesson}, r.finishedSlugs(t, staffEditor))
	assert.Equal(t, 2, r.db.LessonProgress.Query().Where(lessonprogress.UserID(staffEditor.UserID)).CountX(ctx),
		"the started row was finished in place, not duplicated")
}

func TestMergeGuestNeverMovesAnAccountBackwards(t *testing.T) {
	r := newRig(t)
	adminLearner := r.learner(t, "user_admin_learner")

	require.NoError(t, r.progress.MergeGuest(t.Context(), adminLearner.UserID, guestWith(firstLesson).Guest, locale))

	assert.Equal(t, []string{firstLesson}, r.finishedSlugs(t, adminLearner))
	assert.Zero(t, testsupport.PublishedCount(t, r.services))
}

func TestMergeGuestIgnoresWhatItCannotResolve(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, "user_blank")
	stale := progress.GuestProgress{}.
		Record("no-such-course", firstLesson).
		Record(jsSlug, "a-lesson-the-build-dropped")

	require.NoError(t, r.progress.MergeGuest(ctx, blankUser.UserID, stale, locale))
	require.NoError(t, r.progress.MergeGuest(ctx, blankUser.UserID, progress.GuestProgress{}, locale))

	assert.Empty(t, r.finishedSlugs(t, blankUser))
	assert.Zero(t, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID)).CountX(ctx))
	assert.Zero(t, testsupport.PublishedCount(t, r.services))
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
	allDone := r.learner(t, "user_learner_all_lessons_done")
	courseID := r.course(t, jsSlug).ID

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.tx, r.db, courseID, locale))

	state, err := r.progress.CourseState(ctx, allDone, courseID)
	require.NoError(t, err)
	assert.Equal(t, progress.StateFinished, state.State)
	finished := testsupport.PublishedEvents[events.CourseFinished](t, r.services)
	require.Len(t, finished, 1, "nobody else has nothing left to do")
	assert.Equal(t, jsSlug, finished[0].Slug)
	assert.Equal(t, 1, finished[0].OccurrenceCount)

	// A learner already finished is not selected again, so no fact repeats.
	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.tx, r.db, courseID, locale))
	assert.Len(t, testsupport.PublishedEvents[events.CourseFinished](t, r.services), 1)
}

// A new version that drops the lessons a learner had left leaves them with
// nothing to do, so the promotion finishes their course.
func TestReevaluateCompletionFinishesTheCourseOfALearnerWhoseRemainingLessonsWereDropped(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, "user_admin_learner")
	courseID := r.course(t, jsSlug).ID
	r.dropLessonsAfterTheFirst(t) // adminLearner finished the only lesson that is left

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.tx, r.db, courseID, locale))

	state, err := r.progress.CourseState(ctx, adminLearner, courseID)
	require.NoError(t, err)
	assert.Equal(t, progress.StateFinished, state.State)
	announced := testsupport.PublishedCount(t, r.services)
	require.NotZero(t, announced)

	// A learner already finished is not selected again, so no fact repeats.
	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.tx, r.db, courseID, locale))
	assert.Equal(t, announced, testsupport.PublishedCount(t, r.services))
}

func TestReevaluateCompletionLeavesLearnersWithWorkLeftAlone(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.tx, r.db, r.course(t, jsSlug).ID, locale))

	// Work left: the second lesson, or a lesson not even started.
	for _, label := range []string{"user_admin_learner", "user_learner_with_a_gap", "user_learner_hidden_cards"} {
		state, err := r.progress.CourseState(ctx, r.learner(t, label), r.course(t, jsSlug).ID)
		require.NoError(t, err)
		assert.Equal(t, progress.StateStarted, state.State, label)
	}
}

func TestReevaluateCompletionCompletesNobodyForAVersionWithNoLessons(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	bare := r.course(t, "bare")

	require.NoError(t, r.progress.ReevaluateCompletion(ctx, r.tx, r.db, bare.ID, locale))

	assert.Zero(t, testsupport.PublishedCount(t, r.services))
}

func TestReevaluateCompletionRefusesAnUnknownCourse(t *testing.T) {
	r := newRig(t)

	err := r.progress.ReevaluateCompletion(t.Context(), r.tx, r.db, -1, locale)

	assert.True(t, ent.IsNotFound(err), "got %v", err)
}

// Slugs the data may lack.

func TestFactsOfAContentRowWithoutASlugCarryAnEmptySlug(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, "user_blank")
	r.db.Course.Update().Where(course.SlugEQ(jsSlug)).ClearSlug().ExecX(ctx)
	hello := r.lesson(t, firstLesson)
	r.db.CourseLesson.UpdateOneID(hello.ID).ClearSlug().ExecX(ctx)

	_, err := r.progress.StartLesson(ctx, blankUser, hello.ID, locale)
	require.NoError(t, err)

	started := testsupport.PublishedEvents[events.LessonStarted](t, r.services)
	require.Len(t, started, 1)
	assert.Empty(t, started[0].CourseSlug)
	assert.Empty(t, started[0].LessonSlug)
}

// Two requests at once resolve to one row. The race is staged by letting the
// "other request" land its insert just before this one's.

func TestStartLessonSharesAnEnrollmentWonByAConcurrentRequest(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	blankUser := r.learner(t, "user_blank")
	courseID := r.course(t, jsSlug).ID
	raced := false
	r.db.Enrollment.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !raced && m.Op().Is(ent.OpCreate) {
				raced = true
				r.db.Enrollment.Create().SetUserID(blankUser.UserID).SetCourseID(courseID).SetState(progress.StateStarted).ExecX(ctx)
			}
			return next.Mutate(ctx, m)
		})
	})

	_, err := r.progress.StartLesson(ctx, blankUser, r.lesson(t, firstLesson).ID, locale)
	require.NoError(t, err)

	assert.Equal(t, 1, r.db.Enrollment.Query().Where(enrollment.UserID(blankUser.UserID)).CountX(ctx))
	assert.Empty(t, testsupport.PublishedEvents[events.CourseStarted](t, r.services),
		"only the request whose insert landed announces the enrollment")
}

func TestStartLessonSharesLessonProgressWonByAConcurrentRequest(t *testing.T) {
	r := newRig(t)
	ctx := t.Context()
	adminLearner := r.learner(t, "user_admin_learner")
	second := r.lesson(t, secondLesson)
	courseID := r.course(t, jsSlug).ID
	enrolled := r.db.Enrollment.Query().Where(enrollment.UserID(adminLearner.UserID)).OnlyX(ctx)
	raced := false
	r.db.LessonProgress.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if !raced && m.Op().Is(ent.OpCreate) {
				raced = true
				r.db.LessonProgress.Create().
					SetUserID(adminLearner.UserID).SetCourseID(courseID).SetEnrollmentID(enrolled.ID).
					SetLessonID(second.ID).SetState(progress.StateStarted).ExecX(ctx)
			}
			return next.Mutate(ctx, m)
		})
	})

	_, err := r.progress.StartLesson(ctx, adminLearner, second.ID, locale)
	require.NoError(t, err)

	assert.Equal(t, 1, r.db.LessonProgress.Query().
		Where(lessonprogress.UserID(adminLearner.UserID), lessonprogress.LessonID(second.ID)).CountX(ctx))
	assert.Empty(t, testsupport.PublishedEvents[events.LessonStarted](t, r.services))
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
	return []operation{
		{
			name: "starting a lesson",
			prepare: func(t *testing.T, r *rig) func() error {
				blankUser, lesson := r.learner(t, "user_blank"), r.lesson(t, firstLesson)
				return func() error {
					_, err := r.progress.StartLesson(t.Context(), blankUser, lesson.ID, locale)
					return err
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLesson, ent.TypeCourseLessonVersion, ent.TypeEnrollment, ent.TypeLessonProgress},
			writes: []string{ent.TypeEnrollment, ent.TypeLessonProgress, ent.TypeCourse},
		},
		{
			name: "checking a solution",
			prepare: func(t *testing.T, r *rig) func() error {
				check := r.check(t, r.learner(t, "user_blank"), firstLesson)
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
				r.dropLessonsAfterTheFirst(t)
				check := r.check(t, r.learner(t, "user_admin_learner"), firstLesson)
				return func() error {
					_, err := r.progress.CheckSolution(t.Context(), check)
					return err
				}
			},
			reads:  []string{ent.TypeLessonProgress, ent.TypeEnrollment},
			writes: []string{ent.TypeEnrollment},
		},
		{
			name: "reading a learner's courses",
			prepare: func(t *testing.T, r *rig) func() error {
				adminLearner, courseID := r.learner(t, "user_admin_learner"), r.course(t, jsSlug).ID
				return func() error {
					_, err := r.progress.CourseStates(t.Context(), adminLearner.UserID, []int{courseID})
					return err
				}
			},
			reads: []string{ent.TypeCourse, ent.TypeCourseLessonVersion, ent.TypeLessonProgress, ent.TypeEnrollment},
		},
		{
			name: "reading a guest's course",
			prepare: func(t *testing.T, r *rig) func() error {
				courseID := r.course(t, jsSlug).ID
				return func() error {
					_, err := r.progress.CourseState(t.Context(), guestWith(secondLesson), courseID)
					return err
				}
			},
			reads: []string{ent.TypeCourse, ent.TypeCourseLessonVersion},
		},
		{
			name: "merging a guest's progress",
			prepare: func(t *testing.T, r *rig) func() error {
				staffEditor := r.learner(t, "user_staff_editor")
				return func() error {
					return r.progress.MergeGuest(t.Context(), staffEditor.UserID, guestWith(thirdLesson).Guest, locale)
				}
			},
			reads:  []string{ent.TypeCourse, ent.TypeCourseLessonVersion, ent.TypeLessonProgress, ent.TypeEnrollment},
			writes: []string{ent.TypeLessonProgress},
		},
		{
			name: "re-evaluating completion",
			prepare: func(t *testing.T, r *rig) func() error {
				r.dropLessonsAfterTheFirst(t)
				courseID := r.course(t, jsSlug).ID
				return func() error {
					return r.progress.ReevaluateCompletion(t.Context(), r.tx, r.db, courseID, locale)
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

				assert.ErrorIs(t, call(), errStore)
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
				before := testsupport.PublishedCount(t, r.services)
				r.failWrites(entity)

				assert.ErrorIs(t, call(), errStore)
				assert.Equal(t, before, testsupport.PublishedCount(t, r.services), "no fact outlives a failed transition")
			})
		}
	}
}

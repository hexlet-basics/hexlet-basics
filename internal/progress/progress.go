// Package progress owns what a learner has done in a Course: the state
// transitions and domain facts that follow from starting and passing a Lesson,
// and the read model of where they stand (ADR-0012, ADR-0016).
//
// Lessons may be taken in any order. The rules for what counts as finished and
// as a completed Course live here rather than in a handler because three
// surfaces need the same answer — starting a Lesson, checking a solution, and
// reading a Course.
package progress

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/events"
	"hexletbasics/internal/store"
)

// State values on an Enrollment and on Lesson Progress. They are the legacy
// AASM strings, stored as they always were.
const (
	StateStarted  = "started"
	StateFinished = "finished"
)

// Tracker is the seam handlers depend on, so no handler learns where progress
// is stored.
type Tracker interface {
	StartLesson(ctx context.Context, learner Learner, lessonID int, locale string) (*CourseState, error)
	CheckSolution(ctx context.Context, check Check) (*CheckResult, error)
	CourseState(ctx context.Context, learner Learner, courseID int) (*CourseState, error)
	CourseStates(ctx context.Context, userID int, courseIDs []int) (map[int]*CourseState, error)
	MergeGuest(ctx context.Context, userID int, guest GuestProgress, locale string) error
}

// CourseState is where a learner stands in one Course, computed against its
// current Version.
type CourseState struct {
	State          string
	Completion     int
	NextLessonSlug *string
	Lessons        []LessonState
}

// LessonState is one row of that list, in course order.
type LessonState struct {
	Slug     string
	Position int
	Finished bool
}

// Progress records transitions and their domain facts in one transaction, and
// answers the same questions for reads.
type Progress struct {
	db        *ent.Client
	store     store.Transactor
	publisher events.TxPublisher
	runner    ExerciseRunner
	now       func() time.Time
}

// New builds the production tracker.
func New(db *ent.Client, txStore store.Transactor, publisher events.TxPublisher, runner ExerciseRunner) *Progress {
	return &Progress{db: db, store: txStore, publisher: publisher, runner: runner, now: time.Now}
}

// StartLesson enrolls the learner in the Lesson's Course if they are not
// enrolled yet, marks the Lesson started, and returns where they now stand.
//
// Idempotent by design: starting a Lesson that is already started, or already
// finished, succeeds and writes nothing — the frontend calls this from a button
// that a learner can press twice, and from "next" at the end of a Lesson.
//
// A guest stores nothing: their whole state is the Lessons they finished, which
// only a check can add to. They still get their state back, so the page they
// land on renders from the server's answer exactly as a learner's does.
//
// Returns an ent not-found error when the Lesson is not part of its Course's
// current Version — a Lesson a later build dropped is no longer startable.
func (p *Progress) StartLesson(
	ctx context.Context,
	learner Learner,
	lessonID int,
	locale string,
) (*CourseState, error) {
	courseID, err := p.startLesson(ctx, learner, lessonID, locale)
	if err != nil {
		return nil, err
	}
	return p.CourseState(ctx, learner, courseID)
}

// startLesson performs the write half and reports which Course was started, so
// the read that follows it runs outside the transaction rather than inside.
func (p *Progress) startLesson(ctx context.Context, learner Learner, lessonID int, locale string) (int, error) {
	if !learner.SignedIn() {
		// There is no guest equivalent of a started Lesson: the state would
		// carry no data, and the one feature anchored on Lesson Progress — the
		// in-lesson assistant — requires an account anyway. The Lesson is still
		// resolved, so a Lesson the Course does not have reads as not found.
		open, err := p.resolveLesson(ctx, p.db, lessonID)
		if err != nil {
			return 0, err
		}
		return open.course.ID, nil
	}

	courseID := 0
	err := p.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
		userID := learner.UserID
		open, err := p.resolveLesson(ctx, db, lessonID)
		if err != nil {
			return err
		}
		courseID = open.course.ID

		enrolled, createdEnrollment, err := p.enroll(ctx, db, userID, open.course.ID)
		if err != nil {
			return err
		}

		_, createdProgress, err := takeLesson(ctx, db, userID, open.course.ID, enrolled.ID, lessonID)
		if err != nil {
			return err
		}
		if !createdProgress {
			// Already started, or already finished: the learner pressed the
			// button twice, or came back to a Lesson they passed. Nothing
			// changes and no fact is published — the occurrence counts the CRM
			// consumes count transitions, not requests.
			return nil
		}

		occurredAt := p.now()
		if createdEnrollment {
			if err := p.publishCourseStarted(ctx, tx, db, userID, open.course, locale, occurredAt); err != nil {
				return err
			}
		}
		return p.publishLessonStarted(ctx, tx, db, enrolled.ID, open, locale, occurredAt)
	})
	return courseID, err
}

// resolveLesson finds a Lesson in its Course's current Version and the Position
// it has there. Any Lesson may be taken in any order (ADR-0016), so this is only
// a lookup — but it is the one place a Lesson is resolved on every write path,
// and it is why a Lesson a later build dropped reads as not found rather than as
// a Lesson at some stale Position.
func (p *Progress) resolveLesson(ctx context.Context, db *ent.Client, lessonID int) (*resolvedLesson, error) {
	lesson, crs, err := loadLessonCourse(ctx, db, lessonID)
	if err != nil {
		return nil, err
	}

	lessons, err := currentLessons(ctx, db, crs)
	if err != nil {
		return nil, err
	}
	target := positionOf(lessons, lessonID)
	if target == 0 {
		return nil, &ent.NotFoundError{}
	}

	return &resolvedLesson{
		lesson:   lesson,
		course:   crs,
		lessons:  lessons,
		position: target,
	}, nil
}

// resolvedLesson is a Lesson found in its Course, with everything resolving it
// already cost: these travel together through every step that follows, so they
// travel as one value.
type resolvedLesson struct {
	lesson  *ent.CourseLesson
	course  *ent.Course
	lessons []currentLesson

	// position is the Lesson's Position in the current Version.
	position int
}

// publishCourseStarted records the enrollment transition, counting the
// learner's started Enrollments as the legacy publisher did.
func (p *Progress) publishCourseStarted(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	userID int,
	crs *ent.Course,
	locale string,
	occurredAt time.Time,
) error {
	count, err := db.Enrollment.Query().
		Where(enrollment.UserID(userID), enrollment.StateEQ(StateStarted)).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("count started enrollments: %w", err)
	}
	if err := p.publisher.Publish(ctx, tx, events.CourseStarted{
		OccurrenceCount: count,
		Slug:            lo.FromPtr(crs.Slug),
		Locale:          locale,
		OccurredAt:      occurredAt,
	}); err != nil {
		return fmt.Errorf("publish course started: %w", err)
	}
	return nil
}

// publishLessonStarted records the Lesson transition, counting the Lesson
// Progress rows on the Enrollment.
func (p *Progress) publishLessonStarted(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	enrollmentID int,
	open *resolvedLesson,
	locale string,
	occurredAt time.Time,
) error {
	count, err := lessonProgressCount(ctx, db, enrollmentID)
	if err != nil {
		return err
	}
	if err := p.publisher.Publish(ctx, tx, events.LessonStarted{
		OccurrenceCount: count,
		LessonSlug:      lo.FromPtr(open.lesson.Slug),
		CourseSlug:      lo.FromPtr(open.course.Slug),
		Locale:          locale,
		OccurredAt:      occurredAt,
	}); err != nil {
		return fmt.Errorf("publish lesson started: %w", err)
	}
	return nil
}

// enroll returns the learner's Enrollment in the Course, creating it when
// absent, and reports whether this call created it — the course-started fact is
// published on creation only, so a learner returning to a Course they already
// began does not generate a second one.
//
// Two requests enrolling the same learner at once resolve to one row, and only
// the request whose insert landed reports it as created: the unique index
// decides, DO NOTHING hands the loser no row, and the loser reads the winner's.
// That keeps the course-started fact and the Course's member count to one per
// Enrollment.
func (p *Progress) enroll(ctx context.Context, db *ent.Client, userID, courseID int) (*ent.Enrollment, bool, error) {
	existing, err := db.Enrollment.Query().
		Where(enrollment.UserID(userID), enrollment.CourseID(courseID)).
		Only(ctx)
	switch {
	case err == nil:
		return existing, false, nil
	case !ent.IsNotFound(err):
		return nil, false, fmt.Errorf("load enrollment: %w", err)
	}

	// ON CONFLICT DO NOTHING rather than a bare insert: the read above cannot
	// exclude a concurrent request enrolling the same learner between it and the
	// write, and a constraint violation would abort the whole transaction.
	id, err := db.Enrollment.Create().
		SetUserID(userID).
		SetCourseID(courseID).
		SetState(StateStarted).
		OnConflictColumns(enrollment.FieldUserID, enrollment.FieldCourseID).
		DoNothing().
		ID(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		existing, err := db.Enrollment.Query().
			Where(enrollment.UserID(userID), enrollment.CourseID(courseID)).
			Only(ctx)
		if err != nil {
			return nil, false, fmt.Errorf("load concurrent enrollment: %w", err)
		}
		return existing, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("enroll: %w", err)
	}

	row, err := db.Enrollment.Get(ctx, id)
	if err != nil {
		return nil, false, fmt.Errorf("load new enrollment: %w", err)
	}

	// The Course's member count stays denormalized, as the legacy counter cache
	// kept it: it is the marketing figure on catalogue pages, which are read far
	// more often than anyone enrolls, and it is never used as a progress figure.
	if err := db.Course.UpdateOneID(courseID).AddEnrollmentsCount(1).Exec(ctx); err != nil {
		return nil, false, fmt.Errorf("count new enrollment on course %d: %w", courseID, err)
	}
	return row, true, nil
}

// takeLesson returns the learner's Lesson Progress for a Lesson, creating it as
// started when absent, and reports whether this call created it — the
// lesson-started fact is published on creation only.
//
// Two requests taking the same Lesson at once resolve to one row: the unique
// (user, lesson) index decides, DO NOTHING hands the loser no row instead of a
// constraint violation that would abort the transaction, and the loser reads
// the winner's.
func takeLesson(ctx context.Context, db *ent.Client, userID, courseID, enrollmentID, lessonID int) (*ent.LessonProgress, bool, error) {
	load := func() (*ent.LessonProgress, error) {
		return db.LessonProgress.Query().
			Where(lessonprogress.UserID(userID), lessonprogress.LessonID(lessonID)).
			Only(ctx)
	}

	existing, err := load()
	switch {
	case err == nil:
		return existing, false, nil
	case !ent.IsNotFound(err):
		return nil, false, fmt.Errorf("load lesson progress: %w", err)
	}

	id, err := db.LessonProgress.Create().
		SetUserID(userID).
		SetCourseID(courseID).
		SetEnrollmentID(enrollmentID).
		SetLessonID(lessonID).
		SetState(StateStarted).
		OnConflictColumns(lessonprogress.FieldUserID, lessonprogress.FieldLessonID).
		DoNothing().
		ID(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		existing, err := load()
		if err != nil {
			return nil, false, fmt.Errorf("load concurrent lesson progress: %w", err)
		}
		return existing, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("create lesson progress: %w", err)
	}

	row, err := db.LessonProgress.Get(ctx, id)
	if err != nil {
		return nil, false, fmt.Errorf("load new lesson progress: %w", err)
	}
	return row, true, nil
}

// loadLessonCourse resolves a Lesson and the Course it belongs to. A Lesson
// with no Course is unreachable content rather than a server fault, so it reads
// as not found, exactly like a Lesson id that does not exist.
func loadLessonCourse(ctx context.Context, db *ent.Client, lessonID int) (*ent.CourseLesson, *ent.Course, error) {
	lesson, err := db.CourseLesson.Query().
		Where(courselesson.ID(lessonID)).
		Only(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load lesson %d: %w", lessonID, err)
	}
	if lesson.CourseID == nil {
		return nil, nil, &ent.NotFoundError{}
	}

	crs, err := db.Course.Query().Where(course.ID(*lesson.CourseID)).Only(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load course %d: %w", *lesson.CourseID, err)
	}
	return lesson, crs, nil
}

// positionOf is a Lesson's Position in the current Version, or 0 when the
// Version does not contain it. Zero is not a Position — Positions start at one
// — so it reads as "this Lesson has none", which is the same answer lesson
// resolution and the guest merge both need for a Lesson a later build dropped.
func positionOf(lessons []currentLesson, lessonID int) int {
	i := slices.IndexFunc(lessons, func(lesson currentLesson) bool { return lesson.lessonID == lessonID })
	if i < 0 {
		return 0
	}
	return lessons[i].position
}

// currentLesson is a Lesson of the current Version with the Position and slug
// the read models need, ordered by Position.
type currentLesson struct {
	lessonID int
	slug     string
	position int
}

func currentLessons(ctx context.Context, db *ent.Client, crs *ent.Course) ([]currentLesson, error) {
	byCourse, err := currentLessonsByCourse(ctx, db, []*ent.Course{crs})
	if err != nil {
		return nil, err
	}
	return byCourse[crs.ID], nil
}

// lessonSlug is the slug of the Lesson a version row belongs to; a row whose
// lesson is missing reads as the empty slug, like an unnamed lesson.
func lessonSlug(v *ent.CourseLessonVersion) string {
	if v.Edges.Lesson == nil {
		return ""
	}
	return lo.FromPtr(v.Edges.Lesson.Slug)
}

// currentLessonsByCourse is currentLessons for many Courses in two queries, not
// two per Course. A Course with no current Version has no entry.
func currentLessonsByCourse(ctx context.Context, db *ent.Client, courses []*ent.Course) (map[int][]currentLesson, error) {
	courseOf := make(map[int]int, len(courses))
	versionIDs := make([]int, 0, len(courses))
	for _, crs := range courses {
		if crs.CurrentVersionID == nil {
			continue
		}
		courseOf[*crs.CurrentVersionID] = crs.ID
		versionIDs = append(versionIDs, *crs.CurrentVersionID)
	}
	if len(versionIDs) == 0 {
		return nil, nil
	}

	// The version rows carry the Position; the stable Lesson identity carries
	// the slug, reached through the version's lesson edge. Only the columns the
	// ordering needs are read: the rows also hold each lesson's test and starter
	// code, which a position lookup has no use for.
	versions, err := db.CourseLessonVersion.Query().
		Where(courselessonversion.CourseVersionIDIn(versionIDs...)).
		Select(
			courselessonversion.FieldCourseVersionID,
			courselessonversion.FieldLessonID,
			courselessonversion.FieldNaturalOrder,
		).
		WithLesson(func(q *ent.CourseLessonQuery) {
			q.Select(courselesson.FieldSlug)
		}).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load current lesson versions: %w", err)
	}

	byCourse := make(map[int][]currentLesson, len(courses))
	for _, v := range versions {
		if v.NaturalOrder == nil {
			continue
		}
		courseID := courseOf[v.CourseVersionID]
		byCourse[courseID] = append(byCourse[courseID], currentLesson{
			lessonID: v.LessonID,
			slug:     lessonSlug(v),
			position: *v.NaturalOrder,
		})
	}
	for _, lessons := range byCourse {
		slices.SortStableFunc(lessons, func(a, b currentLesson) int { return cmp.Compare(a.position, b.position) })
	}
	return byCourse, nil
}

// CourseState answers where the learner stands. A learner with no Enrollment
// gets the same shape with nothing finished — the absence of a record is the
// "not started" condition, not a state a record can be in, so the read does not
// have to distinguish them.
//
// A guest gets the same shape from the other storage: the Lessons their signed
// cookie names, resolved against the current Version.
func (p *Progress) CourseState(ctx context.Context, learner Learner, courseID int) (*CourseState, error) {
	if learner.SignedIn() {
		// One implementation for accounts, whether one Course is asked for or a
		// dashboard asks for many: a second copy of the read is a second place
		// for availability and completion to drift.
		states, err := p.CourseStates(ctx, learner.UserID, []int{courseID})
		if err != nil {
			return nil, err
		}
		state, ok := states[courseID]
		if !ok {
			return nil, fmt.Errorf("load course %d: %w", courseID, &ent.NotFoundError{})
		}
		return state, nil
	}
	return p.guestCourseState(ctx, learner, courseID)
}

// guestCourseState is CourseState for a visitor, whose position lives in a
// signed cookie rather than in rows.
func (p *Progress) guestCourseState(ctx context.Context, learner Learner, courseID int) (*CourseState, error) {
	crs, err := p.db.Course.Query().Where(course.ID(courseID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("load course %d: %w", courseID, err)
	}

	lessons, err := currentLessons(ctx, p.db, crs)
	if err != nil {
		return nil, err
	}

	finishedIDs := make(map[int]bool, len(lessons))
	for _, lesson := range lessons {
		finishedIDs[lesson.lessonID] = learner.Guest.Has(lo.FromPtr(crs.Slug), lesson.slug)
	}

	state, finishedCount := stateOf(lessons, finishedIDs)

	// A guest has no Enrollment to read a state from, so it is derived from
	// what they have finished. Deriving rather than leaving it null is what
	// makes the payload the client renders identical for both.
	switch {
	case finishedCount == 0:
		state.State = ""
	case finishedCount == len(lessons):
		state.State = StateFinished
	default:
		state.State = StateStarted
	}
	return state, nil
}

// CourseStates is CourseState for a signed-in learner across many Courses, in
// five queries however many there are — a dashboard needs every enrollment's
// state at once, and asking per Course made its cost grow with the learner's
// history. Each answer is exactly what CourseState gives for that Course; the
// two share stateOf, and a Course the learner has not joined comes back with
// nothing finished. A Course id that does not exist is absent from the result.
func (p *Progress) CourseStates(ctx context.Context, userID int, courseIDs []int) (map[int]*CourseState, error) {
	if len(courseIDs) == 0 {
		return map[int]*CourseState{}, nil
	}

	courses, err := p.db.Course.Query().Where(course.IDIn(courseIDs...)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load courses: %w", err)
	}

	lessons, err := currentLessonsByCourse(ctx, p.db, courses)
	if err != nil {
		return nil, err
	}

	finishedIDs, err := finishedLessonsByCourse(ctx, p.db, userID, courseIDs...)
	if err != nil {
		return nil, err
	}

	enrollments, err := p.db.Enrollment.Query().
		Where(enrollment.UserID(userID), enrollment.CourseIDIn(courseIDs...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load enrollments: %w", err)
	}
	enrollmentState := lo.SliceToMap(enrollments, func(e *ent.Enrollment) (int, string) {
		return e.CourseID, lo.FromPtr(e.State)
	})

	states := make(map[int]*CourseState, len(courses))
	for _, crs := range courses {
		state, _ := stateOf(lessons[crs.ID], finishedIDs[crs.ID])
		state.State = enrollmentState[crs.ID]
		states[crs.ID] = state
	}
	return states, nil
}

// stateOf folds a Course's current Lessons and the learner's finished set into
// the read model, and reports how many current Lessons are finished. It is the
// one place the next Lesson and completion are decided, whether the learner is a
// guest, one Course is asked for, or a dashboard asks for many.
func stateOf(lessons []currentLesson, finishedIDs map[int]bool) (*CourseState, int) {
	state := &CourseState{Lessons: make([]LessonState, 0, len(lessons))}
	finishedCount := 0
	for _, l := range lessons {
		finished := finishedIDs[l.lessonID]
		if finished {
			finishedCount++
		}
		state.Lessons = append(state.Lessons, LessonState{
			Slug:     l.slug,
			Position: l.position,
			Finished: finished,
		})
		if !finished && state.NextLessonSlug == nil {
			slug := l.slug
			state.NextLessonSlug = &slug
		}
	}

	// Counted from the current Version's finished Lessons, never from the
	// denormalized counter: that counter includes Lessons later Versions
	// dropped, which is why the legacy serializer had to clamp at 100%.
	if len(lessons) > 0 {
		state.Completion = finishedCount * 100 / len(lessons)
	}

	return state, finishedCount
}

// finishedLessonsByCourse is the set of Lessons a signed-in learner has
// finished in each of the Courses, keyed by Course id then Lesson id. It is the
// one place "finished" is read: every progress read and write that needs the
// finished set goes through it, so a Course the learner never joined and a
// Course they joined without finishing anything both come back empty.
func finishedLessonsByCourse(ctx context.Context, db *ent.Client, userID int, courseIDs ...int) (map[int]map[int]bool, error) {
	rows, err := db.LessonProgress.Query().
		Where(
			lessonprogress.UserID(userID),
			lessonprogress.CourseIDIn(courseIDs...),
			lessonprogress.StateEQ(StateFinished),
		).
		Select(lessonprogress.FieldCourseID, lessonprogress.FieldLessonID).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load finished lessons: %w", err)
	}

	finished := make(map[int]map[int]bool, len(courseIDs))
	for _, row := range rows {
		if finished[row.CourseID] == nil {
			finished[row.CourseID] = map[int]bool{}
		}
		finished[row.CourseID][row.LessonID] = true
	}
	return finished, nil
}

func finishedLessonIDs(ctx context.Context, db *ent.Client, userID, courseID int) (map[int]bool, error) {
	byCourse, err := finishedLessonsByCourse(ctx, db, userID, courseID)
	if err != nil {
		return nil, err
	}
	// A course with nothing finished has no entry; callers only index the
	// result, and indexing a nil map reads as "not finished".
	return byCourse[courseID], nil
}

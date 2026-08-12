// Package progress owns sequential progression: which Lesson a learner is
// allowed to take, and the state transitions and domain facts that follow from
// taking it (ADR-0012).
//
// The rule lives here rather than in a handler because three surfaces need the
// same answer — starting a Lesson, checking a solution, and reading a Course —
// and a gate evaluated in two places is a gate that will eventually disagree
// with itself.
package progress

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
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

// ErrLessonNotAvailable reports a Lesson beyond the learner's gate: its
// Position is more than one past the furthest Position they have finished.
// Distinct from a permission failure — the caller is allowed here, their
// progress is not far enough.
var ErrLessonNotAvailable = errors.New("progress: lesson is not available to this learner yet")

// Tracker is the seam handlers depend on, so no handler learns where progress
// is stored or how the gate is evaluated.
type Tracker interface {
	StartLesson(ctx context.Context, userID, lessonID int, locale string) error
	CourseState(ctx context.Context, userID, courseID int) (*CourseState, error)
	MergeGuest(ctx context.Context, userID int, guest GuestProgress, locale string) error
}

// CourseState is where a learner stands in one Course, computed against its
// current Version. The per-Lesson list carries availability rather than leaving
// the client to derive it: a client-side gate is a second implementation of the
// rule, free to drift from the server that enforces it.
type CourseState struct {
	State                    string
	Completion               int
	NextLessonSlug           *string
	FurthestFinishedPosition int
	Lessons                  []LessonState
}

// LessonState is one row of that list, in course order.
type LessonState struct {
	Slug      string
	Position  int
	Finished  bool
	Available bool
}

// Progress records transitions and their domain facts in one transaction, and
// answers the same questions for reads.
type Progress struct {
	db        *ent.Client
	store     store.Transactor
	publisher events.TxPublisher
	now       func() time.Time
}

// New builds the production tracker.
func New(db *ent.Client, txStore store.Transactor, publisher events.TxPublisher) *Progress {
	return &Progress{db: db, store: txStore, publisher: publisher, now: time.Now}
}

// StartLesson enrolls the learner in the Lesson's Course if they are not
// enrolled yet and marks the Lesson started.
//
// Idempotent by design: starting a Lesson that is already started, or already
// finished, succeeds and writes nothing — the frontend calls this from a button
// that a learner can press twice, and from "next" at the end of a Lesson.
//
// Returns ErrLessonNotAvailable when the gate refuses, and an ent not-found
// error when the Lesson is not part of its Course's current Version — a Lesson
// a later build dropped is no longer startable.
func (p *Progress) StartLesson(ctx context.Context, userID, lessonID int, locale string) error {
	return p.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
		lesson, err := db.CourseLesson.Query().
			Where(courselesson.ID(lessonID)).
			Only(ctx)
		if err != nil {
			return fmt.Errorf("load lesson %d: %w", lessonID, err)
		}
		if lesson.CourseID == nil {
			return &ent.NotFoundError{}
		}

		crs, err := db.Course.Query().Where(course.ID(*lesson.CourseID)).Only(ctx)
		if err != nil {
			return fmt.Errorf("load course %d: %w", *lesson.CourseID, err)
		}

		positions, err := currentPositions(ctx, db, crs)
		if err != nil {
			return err
		}
		target, inCurrentVersion := positions[lessonID]
		if !inCurrentVersion {
			return &ent.NotFoundError{}
		}

		furthest, err := furthestFinishedPosition(ctx, db, userID, crs.ID, positions)
		if err != nil {
			return err
		}
		if target > furthest+1 {
			return ErrLessonNotAvailable
		}

		enrolled, createdEnrollment, err := p.enroll(ctx, db, userID, crs.ID)
		if err != nil {
			return err
		}

		alreadyTaken, err := db.LessonProgress.Query().
			Where(
				lessonprogress.UserID(userID),
				lessonprogress.LessonID(lessonID),
			).
			Exist(ctx)
		if err != nil {
			return fmt.Errorf("load lesson progress: %w", err)
		}
		if alreadyTaken {
			// Already started, or already finished: the learner pressed the
			// button twice, or came back to a Lesson they passed. Nothing
			// changes and no fact is published — the occurrence counts the CRM
			// consumes count transitions, not requests.
			return nil
		}

		if _, err := db.LessonProgress.Create().
			SetUserID(userID).
			SetCourseID(crs.ID).
			SetEnrollmentID(enrolled.ID).
			SetLessonID(lessonID).
			SetState(StateStarted).
			Save(ctx); err != nil {
			return fmt.Errorf("create lesson progress: %w", err)
		}

		occurredAt := p.now()
		if createdEnrollment {
			count, err := db.Enrollment.Query().
				Where(enrollment.UserID(userID), enrollment.StateEQ(StateStarted)).
				Count(ctx)
			if err != nil {
				return fmt.Errorf("count started enrollments: %w", err)
			}
			if err := p.publisher.Publish(ctx, tx, events.CourseStarted{
				OccurrenceCount: count,
				Slug:            slugOf(crs),
				Locale:          locale,
				OccurredAt:      occurredAt,
			}); err != nil {
				return fmt.Errorf("publish course started: %w", err)
			}
		}

		lessonCount, err := db.LessonProgress.Query().
			Where(lessonprogress.EnrollmentID(enrolled.ID)).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("count lesson progress: %w", err)
		}

		if err := p.publisher.Publish(ctx, tx, events.LessonStarted{
			OccurrenceCount: lessonCount,
			LessonSlug:      slugOfLesson(lesson),
			CourseSlug:      slugOf(crs),
			Locale:          locale,
			OccurredAt:      occurredAt,
		}); err != nil {
			return fmt.Errorf("publish lesson started: %w", err)
		}
		return nil
	})
}

// enroll returns the learner's Enrollment in the Course, creating it when
// absent, and reports whether this call created it — the course-started fact is
// published on creation only, so a learner returning to a Course they already
// began does not generate a second one.
//
// The narrow race where two requests enroll the same learner at once resolves
// to one row (the unique index decides), and both may report it as created. A
// duplicate fact is the safe failure here; a lost row would not be.
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
		Ignore().
		ID(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("enroll: %w", err)
	}

	row, err := db.Enrollment.Get(ctx, id)
	if err != nil {
		return nil, false, fmt.Errorf("load new enrollment: %w", err)
	}
	return row, true, nil
}

// currentPositions maps every Lesson of the Course's current Version to its
// Position. natural_order numbers Lessons 1..N across all modules in build
// order, which is what course order means; `order` is per-module and cannot
// gate anything on its own.
//
// A Course with no current Version has no positions, so nothing in it is
// startable — which is correct: there is no built content to start.
func currentPositions(ctx context.Context, db *ent.Client, crs *ent.Course) (map[int]int, error) {
	lessons, err := currentLessons(ctx, db, crs)
	if err != nil {
		return nil, err
	}
	positions := make(map[int]int, len(lessons))
	for _, l := range lessons {
		positions[l.lessonID] = l.position
	}
	return positions, nil
}

// currentLesson is a Lesson of the current Version with the Position and slug
// the read models need, ordered by Position.
type currentLesson struct {
	lessonID int
	slug     string
	position int
}

func currentLessons(ctx context.Context, db *ent.Client, crs *ent.Course) ([]currentLesson, error) {
	if crs.CurrentVersionID == nil {
		return nil, nil
	}

	versions, err := db.CourseLessonVersion.Query().
		Where(courselessonversion.CourseVersionID(*crs.CurrentVersionID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load current lesson versions: %w", err)
	}

	// The version rows carry the Position; the stable Lesson identity carries
	// the slug. CourseLessonVersion declares no edge to it (the schema is the
	// loader's write shape), so they are joined here in one extra query.
	ids := make([]int, 0, len(versions))
	for _, v := range versions {
		if v.NaturalOrder != nil {
			ids = append(ids, v.LessonID)
		}
	}
	rows, err := db.CourseLesson.Query().Where(courselesson.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load current lessons: %w", err)
	}
	slugs := make(map[int]string, len(rows))
	for _, row := range rows {
		slugs[row.ID] = lo.FromPtr(row.Slug)
	}

	lessons := make([]currentLesson, 0, len(versions))
	for _, v := range versions {
		if v.NaturalOrder == nil {
			continue
		}
		lessons = append(lessons, currentLesson{
			lessonID: v.LessonID,
			slug:     slugs[v.LessonID],
			position: *v.NaturalOrder,
		})
	}
	sort.Slice(lessons, func(i, j int) bool { return lessons[i].position < lessons[j].position })
	return lessons, nil
}

// CourseState answers where the learner stands. A learner with no Enrollment
// gets the same shape with nothing finished — the absence of a record is the
// "not started" condition, not a state a record can be in, so the read does not
// have to distinguish them.
func (p *Progress) CourseState(ctx context.Context, userID, courseID int) (*CourseState, error) {
	crs, err := p.db.Course.Query().Where(course.ID(courseID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("load course %d: %w", courseID, err)
	}

	lessons, err := currentLessons(ctx, p.db, crs)
	if err != nil {
		return nil, err
	}

	finishedIDs, err := finishedLessonIDs(ctx, p.db, userID, courseID)
	if err != nil {
		return nil, err
	}

	state := &CourseState{Lessons: make([]LessonState, 0, len(lessons))}
	finishedCount := 0
	for _, l := range lessons {
		finished := finishedIDs[l.lessonID]
		if finished {
			finishedCount++
			if l.position > state.FurthestFinishedPosition {
				state.FurthestFinishedPosition = l.position
			}
		}
	}
	for _, l := range lessons {
		finished := finishedIDs[l.lessonID]
		state.Lessons = append(state.Lessons, LessonState{
			Slug:      l.slug,
			Position:  l.position,
			Finished:  finished,
			Available: l.position <= state.FurthestFinishedPosition+1,
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

	enrolled, err := p.db.Enrollment.Query().
		Where(enrollment.UserID(userID), enrollment.CourseID(courseID)).
		Only(ctx)
	switch {
	case err == nil:
		state.State = lo.FromPtr(enrolled.State)
	case !ent.IsNotFound(err):
		return nil, fmt.Errorf("load enrollment: %w", err)
	}
	return state, nil
}

func finishedLessonIDs(ctx context.Context, db *ent.Client, userID, courseID int) (map[int]bool, error) {
	rows, err := db.LessonProgress.Query().
		Where(
			lessonprogress.UserID(userID),
			lessonprogress.CourseID(courseID),
			lessonprogress.StateEQ(StateFinished),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load finished lessons: %w", err)
	}

	finished := make(map[int]bool, len(rows))
	for _, row := range rows {
		finished[row.LessonID] = true
	}
	return finished, nil
}

// furthestFinishedPosition is the highest Position the learner has finished in
// this Course, or 0 when they have finished nothing. Availability is measured
// from the furthest finished Lesson rather than the first unfinished one, so a
// gap in a learner's history does not block them — and a finished Lesson that a
// later Version dropped has no Position and therefore raises nothing.
func furthestFinishedPosition(ctx context.Context, db *ent.Client, userID, courseID int, positions map[int]int) (int, error) {
	finished, err := db.LessonProgress.Query().
		Where(
			lessonprogress.UserID(userID),
			lessonprogress.CourseID(courseID),
			lessonprogress.StateEQ(StateFinished),
		).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("load finished lessons: %w", err)
	}

	furthest := 0
	for _, row := range finished {
		if position, ok := positions[row.LessonID]; ok && position > furthest {
			furthest = position
		}
	}
	return furthest, nil
}

func slugOf(c *ent.Course) string {
	if c.Slug == nil {
		return ""
	}
	return *c.Slug
}

func slugOfLesson(l *ent.CourseLesson) string {
	if l.Slug == nil {
		return ""
	}
	return *l.Slug
}

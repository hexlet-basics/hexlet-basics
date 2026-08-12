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
	"time"

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
}

// Progress records transitions and their domain facts in one transaction.
type Progress struct {
	store     store.Transactor
	publisher events.TxPublisher
	now       func() time.Time
}

// New builds the production tracker.
func New(txStore store.Transactor, publisher events.TxPublisher) *Progress {
	return &Progress{store: txStore, publisher: publisher, now: time.Now}
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
	if crs.CurrentVersionID == nil {
		return map[int]int{}, nil
	}

	versions, err := db.CourseLessonVersion.Query().
		Where(courselessonversion.CourseVersionID(*crs.CurrentVersionID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load current lesson versions: %w", err)
	}

	positions := make(map[int]int, len(versions))
	for _, v := range versions {
		if v.NaturalOrder == nil {
			continue
		}
		positions[v.LessonID] = *v.NaturalOrder
	}
	return positions, nil
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

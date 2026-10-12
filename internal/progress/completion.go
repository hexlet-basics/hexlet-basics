package progress

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/events"
)

// ReevaluateCompletion finishes the Enrollments a newly promoted Version has
// completed, and publishes the fact for each.
//
// Completion is otherwise only ever evaluated when a learner submits a
// solution, which leaves a hole: dropping a Lesson can complete a Course for
// someone who has nothing left to submit. Evaluating it at the moment the
// Lesson set changes closes that hole without the hourly sweep the legacy
// system used — and therefore without a recurring scheduler the Go stack does
// not have.
//
// Learners who are not complete are untouched: this only ever moves an
// Enrollment to finished, never backwards, so a Version that inserts a Lesson
// leaves them where they are. A learner already finished is not selected, so no
// fact is published twice.
func (p *Progress) ReevaluateCompletion(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	courseID int,
	locale string,
) error {
	crs, err := db.Course.Query().Where(course.ID(courseID)).Only(ctx)
	if err != nil {
		return fmt.Errorf("load course %d: %w", courseID, err)
	}

	lessons, err := currentLessons(ctx, db, crs)
	if err != nil {
		return err
	}
	if len(lessons) == 0 {
		// A Version with no Lessons completes nobody; the alternative reading —
		// everybody has finished all zero of them — is the worse one.
		return nil
	}

	lessonIDs := make([]int, 0, len(lessons))
	for _, lesson := range lessons {
		lessonIDs = append(lessonIDs, lesson.lessonID)
	}

	completed, err := completedEnrollments(ctx, db, courseID, lessonIDs)
	if err != nil {
		return err
	}
	if len(completed) == 0 {
		return nil
	}

	enrollmentIDs := make([]int, 0, len(completed))
	userIDs := make([]int, 0, len(completed))
	for _, row := range completed {
		enrollmentIDs = append(enrollmentIDs, row.id)
		userIDs = append(userIDs, row.userID)
	}

	if err := db.Enrollment.Update().
		Where(enrollment.IDIn(enrollmentIDs...)).
		SetState(StateFinished).
		Exec(ctx); err != nil {
		return fmt.Errorf("finish completed enrollments: %w", err)
	}

	counts, err := finishedEnrollmentCounts(ctx, db, userIDs)
	if err != nil {
		return err
	}

	occurredAt := p.now()
	for _, row := range completed {
		if err := p.publisher.Publish(ctx, tx, events.CourseFinished{
			OccurrenceCount: counts[row.userID],
			Slug:            slugOf(crs),
			Locale:          locale,
			OccurredAt:      occurredAt,
		}); err != nil {
			return fmt.Errorf("publish course finished: %w", err)
		}
	}
	return nil
}

// completedEnrollment is one row of the sweep: the Enrollment to finish and the
// learner whose occurrence count the fact carries.
type completedEnrollment struct {
	id     int
	userID int
}

// completedEnrollments selects the Enrollments a promoted Version has just
// completed: not yet finished, and having finished every Lesson the new Version
// contains.
//
// A NULL state counts as started. AASM always wrote the column, but the legacy
// baseline allows NULL and rows predating it exist; the conversion layer
// already resolves NULL to started, and a learner whose row happens to be NULL
// must not be the one learner who can never complete.
//
// The shape is deliberate. A correlated per-Enrollment subquery would make the
// promotion pay for every Enrollment the Course has; this counts each learner's
// distinct finished Lessons of the new Version in one grouped pass, keeps those
// matching the Version's Lesson count, and then loads their unfinished
// Enrollments. Work is proportional to the finished Lesson Progress rows of the
// new Version's Lessons — the same order as the hundreds of content rows the
// promotion transaction is already writing.
func completedEnrollments(ctx context.Context, db *ent.Client, courseID int, lessonIDs []int) ([]completedEnrollment, error) {
	var finishedPerLearner []struct {
		UserID  int `json:"user_id"`
		Lessons int `json:"lessons"`
	}
	if err := db.LessonProgress.Query().
		Where(
			lessonprogress.State(StateFinished),
			lessonprogress.LessonIDIn(lessonIDs...),
		).
		GroupBy(lessonprogress.FieldUserID).
		Aggregate(func(s *entsql.Selector) string {
			return entsql.As(entsql.Count(entsql.Distinct(s.C(lessonprogress.FieldLessonID))), "lessons")
		}).
		Scan(ctx, &finishedPerLearner); err != nil {
		return nil, fmt.Errorf("count finished lessons per learner: %w", err)
	}

	var finishers []int
	for _, row := range finishedPerLearner {
		if row.Lessons == len(lessonIDs) {
			finishers = append(finishers, row.UserID)
		}
	}
	if len(finishers) == 0 {
		return nil, nil
	}

	enrollments, err := db.Enrollment.Query().
		Where(
			enrollment.CourseID(courseID),
			enrollment.UserIDIn(finishers...),
			enrollment.Or(enrollment.StateIsNil(), enrollment.StateNEQ(StateFinished)),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("select completed enrollments: %w", err)
	}

	completed := make([]completedEnrollment, 0, len(enrollments))
	for _, e := range enrollments {
		completed = append(completed, completedEnrollment{id: e.ID, userID: e.UserID})
	}
	return completed, nil
}

// finishedEnrollmentCounts is each learner's number of finished Enrollments,
// read after the update so the freshly finished Course is included — the legacy
// events counted the transition itself. One grouped query keeps the occurrence
// counts from costing a query per transitioning learner.
func finishedEnrollmentCounts(ctx context.Context, db *ent.Client, userIDs []int) (map[int]int, error) {
	var rows []struct {
		UserID int `json:"user_id"`
		Count  int `json:"count"`
	}
	if err := db.Enrollment.Query().
		Where(enrollment.UserIDIn(userIDs...), enrollment.State(StateFinished)).
		GroupBy(enrollment.FieldUserID).
		Aggregate(ent.Count()).
		Scan(ctx, &rows); err != nil {
		return nil, fmt.Errorf("count finished enrollments: %w", err)
	}

	counts := make(map[int]int, len(rows))
	for _, row := range rows {
		counts[row.UserID] = row.Count
	}
	return counts, nil
}

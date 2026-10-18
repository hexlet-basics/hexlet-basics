package progress

import (
	"context"
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/samber/lo"

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
		return errs.With("course_id", courseID).Wrapf(err, "load course")
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

	lessonIDs := lo.Map(lessons, func(lesson currentLesson, _ int) int { return lesson.lessonID })

	completed, err := completedEnrollments(ctx, db, courseID, lessonIDs)
	if err != nil {
		return err
	}

	if len(completed) == 0 {
		return nil
	}

	enrollmentIDs := lo.Map(completed, func(row completedEnrollment, _ int) int { return row.id })
	userIDs := lo.Map(completed, func(row completedEnrollment, _ int) int { return row.userID })

	if err := db.Enrollment.Update().
		Where(enrollment.IDIn(enrollmentIDs...)).
		SetState(StateFinished).
		Exec(ctx); err != nil {
		return errs.Wrapf(err, "finish completed enrollments")
	}

	counts, err := finishedEnrollmentCounts(ctx, db, userIDs)
	if err != nil {
		return err
	}

	occurredAt := p.now()
	for _, row := range completed {
		if err := p.publisher.Publish(ctx, tx, events.CourseFinished{
			OccurrenceCount: counts[row.userID],
			Slug:            lo.FromPtr(crs.Slug),
			Locale:          locale,
			OccurredAt:      occurredAt,
		}); err != nil {
			return errs.Wrapf(err, "publish course finished")
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
	type finishedPerLearnerRow struct {
		UserID  int `json:"user_id"`
		Lessons int `json:"lessons"`
	}

	var finishedPerLearner []finishedPerLearnerRow
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
		return nil, errs.Wrapf(err, "count finished lessons per learner")
	}

	finishers := lo.FilterMap(finishedPerLearner, func(row finishedPerLearnerRow, _ int) (int, bool) {
		return row.UserID, row.Lessons == len(lessonIDs)
	})
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
		return nil, errs.Wrapf(err, "select completed enrollments")
	}

	return lo.Map(enrollments, func(e *ent.Enrollment, _ int) completedEnrollment {
		return completedEnrollment{id: e.ID, userID: e.UserID}
	}), nil
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
		return nil, errs.Wrapf(err, "count finished enrollments")
	}

	counts := make(map[int]int, len(rows))
	for _, row := range rows {
		counts[row.UserID] = row.Count
	}

	return counts, nil
}

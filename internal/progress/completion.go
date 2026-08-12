package progress

import (
	"context"
	"database/sql"
	"fmt"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/events"
)

// completionSweepSQL selects the Enrollments a promoted Version has just
// completed: not yet finished, and having finished every Lesson the new Version
// contains.
//
// A NULL state counts as started. AASM always wrote the column, but the legacy
// baseline allows NULL and rows predating it exist; the conversion layer
// already resolves NULL to started, and a learner whose row happens to be NULL
// must not be the one learner who can never complete.
//
// The shape is deliberate. A correlated per-Enrollment subquery would make the
// promotion pay for every Enrollment the Course has; this is one pass — join
// the Course's Enrollments to their finished Lesson Progress rows restricted to
// the new Version's Lessons, and keep the learners whose distinct count matches
// the Version's Lesson count. Work is proportional to the Course's finished
// Lesson Progress rows and is a single statement, which is the same order as
// the hundreds of content rows the promotion transaction is already writing.
var completionSweepSQL = fmt.Sprintf(`
SELECT lm.id, lm.user_id
FROM language_members lm
JOIN language_lesson_members llm
  ON llm.user_id = lm.user_id
 AND llm.state = '%s'
 AND llm.lesson_id = ANY($2)
WHERE lm.language_id = $1
  AND (lm.state IS NULL OR lm.state <> '%s')
GROUP BY lm.id, lm.user_id
HAVING count(DISTINCT llm.lesson_id) = $3
`, StateFinished, StateFinished)

// finishedEnrollmentCountsSQL counts each learner's finished Enrollments in one
// statement, so the occurrence counts the facts carry do not cost a query per
// transitioning learner.
const finishedEnrollmentCountsSQL = `
SELECT user_id, count(*)
FROM language_members
WHERE user_id = ANY($1) AND state = $2
GROUP BY user_id
`

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

	lessonIDs := make([]int64, 0, len(lessons))
	for _, lesson := range lessons {
		lessonIDs = append(lessonIDs, int64(lesson.lessonID))
	}

	completed, err := completedEnrollments(ctx, tx, courseID, lessonIDs)
	if err != nil {
		return err
	}
	if len(completed) == 0 {
		return nil
	}

	enrollmentIDs := make([]int64, 0, len(completed))
	userIDs := make([]int64, 0, len(completed))
	for _, row := range completed {
		enrollmentIDs = append(enrollmentIDs, int64(row.id))
		userIDs = append(userIDs, int64(row.userID))
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE language_members SET state = $2, updated_at = now() WHERE id = ANY($1)`,
		enrollmentIDs, StateFinished,
	); err != nil {
		return fmt.Errorf("finish completed enrollments: %w", err)
	}

	counts, err := finishedEnrollmentCounts(ctx, tx, userIDs)
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

func completedEnrollments(ctx context.Context, tx *sql.Tx, courseID int, lessonIDs []int64) ([]completedEnrollment, error) {
	rows, err := tx.QueryContext(ctx, completionSweepSQL, courseID, lessonIDs, len(lessonIDs))
	if err != nil {
		return nil, fmt.Errorf("select completed enrollments: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var completed []completedEnrollment
	for rows.Next() {
		var row completedEnrollment
		if err := rows.Scan(&row.id, &row.userID); err != nil {
			return nil, fmt.Errorf("scan completed enrollment: %w", err)
		}
		completed = append(completed, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read completed enrollments: %w", err)
	}
	return completed, nil
}

// finishedEnrollmentCounts is each learner's number of finished Enrollments,
// read after the update so the freshly finished Course is included — the legacy
// events counted the transition itself.
func finishedEnrollmentCounts(ctx context.Context, tx *sql.Tx, userIDs []int64) (map[int]int, error) {
	rows, err := tx.QueryContext(ctx, finishedEnrollmentCountsSQL, userIDs, StateFinished)
	if err != nil {
		return nil, fmt.Errorf("count finished enrollments: %w", err)
	}
	defer func() { _ = rows.Close() }()

	counts := map[int]int{}
	for rows.Next() {
		var userID, count int
		if err := rows.Scan(&userID, &count); err != nil {
			return nil, fmt.Errorf("scan finished enrollment count: %w", err)
		}
		counts[userID] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read finished enrollment counts: %w", err)
	}
	return counts, nil
}

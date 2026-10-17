package lessonreviews

import (
	"context"

	"github.com/riverqueue/river"

	"hexletbasics/ent"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/predicate"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/jobs"
)

// Enqueuer schedules background jobs (jobs.Enqueuer in production and in tests,
// the latter over the test's transaction).
type Enqueuer interface {
	Enqueue(ctx context.Context, args ...river.JobArgs) error
}

// Scheduler asks for AI reviews of lessons: it picks the translations to
// summarize and schedules one job for each. Each job is independent and
// idempotent (the worker upserts by its natural key), so an administrator who
// is unsure whether it worked just asks again.
type Scheduler struct {
	db   *ent.Client
	jobs Enqueuer
}

// NewScheduler builds the scheduler over an ent client and the job enqueuer.
func NewScheduler(db *ent.Client, enqueuer Enqueuer) *Scheduler {
	return &Scheduler{db: db, jobs: enqueuer}
}

// ReviewCourse schedules a review of every translation of the course's current
// version, all locales included (legacy
// `language.current_lesson_infos.find_each { ReviewLessonJob.perform_later }`).
// A course without a current version schedules nothing. A missing course is
// domain.ErrNotFound.
func (s *Scheduler) ReviewCourse(ctx context.Context, courseID int) error {
	course, err := s.db.Course.Get(ctx, courseID)
	if err != nil {
		return domain.FromStore(err)
	}
	if course.CurrentVersionID == nil {
		return nil
	}
	return s.review(ctx, courselessontranslation.CourseVersionID(*course.CurrentVersionID))
}

// ReviewLesson schedules a review of every translation of one lesson: all
// locales and all versions, so the historical summaries refresh too (legacy
// `lesson.infos.find_each { ReviewLessonJob.perform_later }`). A missing lesson
// is domain.ErrNotFound.
func (s *Scheduler) ReviewLesson(ctx context.Context, lessonID int) error {
	lesson, err := s.db.CourseLesson.Get(ctx, lessonID)
	if err != nil {
		return domain.FromStore(err)
	}
	return s.review(ctx, courselessontranslation.CourseLessonID(lesson.ID))
}

// review schedules the jobs in one transaction, so a failure schedules none.
func (s *Scheduler) review(ctx context.Context, within ...predicate.CourseLessonTranslation) error {
	ids, err := s.db.CourseLessonTranslation.Query().
		Where(within...).
		Order(ent.Asc(courselessontranslation.FieldID)).
		IDs(ctx)
	if err != nil {
		return domain.FromStore(err)
	}

	args := make([]river.JobArgs, len(ids))
	for i, id := range ids {
		args[i] = jobs.ReviewLessonArgs{LessonInfoID: id}
	}
	return s.jobs.Enqueue(ctx, args...)
}

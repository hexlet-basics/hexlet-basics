package handlers

import (
	"context"

	"github.com/riverqueue/river"

	"hexletbasics/ent"
	"hexletbasics/internal/jobs"
)

// VersionBuildStarter is the atomic operation handlers need: create a course
// version and enqueue its loader job. Production performs both writes in one
// SQL transaction; tests use a moq mock whose default writes the version through their
// rollback-only DB.
type VersionBuildStarter interface {
	Start(ctx context.Context, courseID int) (*ent.CourseVersion, error)
}

// JobEnqueuer schedules background jobs from a handler (jobs.Enqueuer in
// production and in tests, over the test's savepoint transaction). The handler
// names what should run; each kind's queue, attempts and uniqueness live with
// its Args, so a request repeated within an email cooldown is skipped, not
// refused.
type JobEnqueuer interface {
	Enqueue(ctx context.Context, args ...river.JobArgs) error
}

// RelatedCoursesReplacer swaps a blog post's related-courses set atomically
// (relatedcourses.Replacer, over the production transactor or a test's
// savepoint transactor).
type RelatedCoursesReplacer interface {
	Replace(ctx context.Context, postID int, courseIDs []int) error
}

// enqueueLessonReviews schedules one review job per lesson version info, all in
// one transaction so a failure enqueues none; each job is independent and
// idempotent (the worker upserts by its natural key), so the admin just clicks
// again.
func (s *Server) enqueueLessonReviews(ctx context.Context, lessonInfoIDs []int) error {
	args := make([]river.JobArgs, len(lessonInfoIDs))
	for i, id := range lessonInfoIDs {
		args[i] = jobs.ReviewLessonArgs{LessonInfoID: id}
	}
	return s.jobs.Enqueue(ctx, args...)
}

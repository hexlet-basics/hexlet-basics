package handlers

import (
	"context"

	"github.com/riverqueue/river"

	"hexletbasics/ent"
)

// VersionBuildStarter is the atomic operation handlers need: create a course
// version and enqueue its loader job, for a course named by id (the admin) or by
// slug (the repository's webhook). Production performs both writes in one SQL
// transaction.
type VersionBuildStarter interface {
	Start(ctx context.Context, courseID int) (*ent.CourseVersion, error)
	StartBySlug(ctx context.Context, slug string) (*ent.CourseVersion, error)
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

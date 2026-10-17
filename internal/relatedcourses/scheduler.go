package relatedcourses

import (
	"context"

	"github.com/riverqueue/river"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/jobs"
)

// Enqueuer schedules background jobs (jobs.Enqueuer in production and in tests,
// the latter over the test's transaction).
type Enqueuer interface {
	Enqueue(ctx context.Context, args ...river.JobArgs) error
}

// Scheduler asks for the AI pick of a post's related courses. The worker
// replaces the set asynchronously, so a request answers nothing but whether the
// job was scheduled.
type Scheduler struct {
	db   *ent.Client
	jobs Enqueuer
}

// NewScheduler builds the scheduler over an ent client and the job enqueuer.
func NewScheduler(db *ent.Client, enqueuer Enqueuer) *Scheduler {
	return &Scheduler{db: db, jobs: enqueuer}
}

// Suggest schedules the pick for one post (legacy `related_courses` member
// action → FindRelatedCoursesForBlogPostJob). A missing post is
// domain.ErrNotFound, before anything is scheduled.
func (s *Scheduler) Suggest(ctx context.Context, postID int) error {
	if _, err := s.db.BlogPost.Query().Where(blogpost.IDEQ(postID)).Only(ctx); err != nil {
		return domain.FromStore(err)
	}
	return s.jobs.Enqueue(ctx, jobs.SuggestRelatedCoursesArgs{BlogPostID: postID})
}

package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// RelatedCoursesSuggester is the AI related-courses seam used by the River
// worker (implemented by relatedcourses.Suggester).
type RelatedCoursesSuggester interface {
	SuggestRelatedCourses(ctx context.Context, blogPostID int) error
}

// SuggestRelatedCoursesArgs identifies the blog post whose related courses the
// LLM should pick (legacy FindRelatedCoursesForBlogPostJob).
type SuggestRelatedCoursesArgs struct {
	BlogPostID int `json:"blog_post_id"`
}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (SuggestRelatedCoursesArgs) Kind() string { return "suggest_related_courses" }

type suggestRelatedCoursesWorker struct {
	river.WorkerDefaults[SuggestRelatedCoursesArgs]
	suggester RelatedCoursesSuggester
}

func (w *suggestRelatedCoursesWorker) Work(ctx context.Context, job *river.Job[SuggestRelatedCoursesArgs]) error {
	return w.suggester.SuggestRelatedCourses(ctx, job.Args.BlogPostID)
}

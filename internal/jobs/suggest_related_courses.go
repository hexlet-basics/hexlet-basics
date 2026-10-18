package jobs

import (
	"context"
	"time"

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

// InsertOpts bounds retries of a paid LLM call; see llmMaxAttempts.
func (SuggestRelatedCoursesArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: llmMaxAttempts, Queue: QueueLLM}
}

// SuggestRelatedCourses picks blog-post related courses with suggester.
func SuggestRelatedCourses(suggester RelatedCoursesSuggester) Module {
	return module(SuggestRelatedCoursesWorker(suggester))
}

// SuggestRelatedCoursesWorker is the worker alone, without its registration, so that a test
// can work a job through it in a transaction (rivertest.Worker).
func SuggestRelatedCoursesWorker(suggester RelatedCoursesSuggester) river.Worker[SuggestRelatedCoursesArgs] {
	return &suggestRelatedCoursesWorker{suggester: suggester}
}

type suggestRelatedCoursesWorker struct {
	river.WorkerDefaults[SuggestRelatedCoursesArgs]

	suggester RelatedCoursesSuggester
}

// Timeout overrides River's one-minute default, which an LLM answer can exceed.
func (*suggestRelatedCoursesWorker) Timeout(*river.Job[SuggestRelatedCoursesArgs]) time.Duration {
	return llmTimeout
}

func (w *suggestRelatedCoursesWorker) Work(ctx context.Context, job *river.Job[SuggestRelatedCoursesArgs]) error {
	return cancelUnfixable(w.suggester.SuggestRelatedCourses(ctx, job.Args.BlogPostID))
}

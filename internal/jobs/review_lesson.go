package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
)

const (
	// llmMaxAttempts caps retries of paid LLM jobs (review, related courses).
	llmMaxAttempts = 3
	// llmTimeout bounds one LLM job; River's default is one minute.
	llmTimeout = 5 * time.Minute
)

// LessonReviewer is the AI lesson-review seam used by the River worker
// (implemented by lessonreviews.Reviewer).
type LessonReviewer interface {
	ReviewLesson(ctx context.Context, lessonInfoID int) error
}

// ReviewLessonArgs identifies one lesson version info to summarize.
type ReviewLessonArgs struct {
	LessonInfoID int `json:"lesson_info_id"`
}

// Kind is River's stable job discriminator.
func (ReviewLessonArgs) Kind() string { return "review_lesson" }

// InsertOpts bounds retries of a paid LLM call: River's default of 25 attempts
// would keep re-billing a failing review for weeks. An admin can re-trigger it.
func (ReviewLessonArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: llmMaxAttempts, Queue: QueueLLM}
}

// ReviewLesson summarizes lessons with reviewer.
func ReviewLesson(reviewer LessonReviewer) Module {
	return module(ReviewLessonWorker(reviewer))
}

// ReviewLessonWorker is the worker alone, without its registration, so that a test
// can work a job through it in a transaction (rivertest.Worker).
func ReviewLessonWorker(reviewer LessonReviewer) river.Worker[ReviewLessonArgs] {
	return &reviewLessonWorker{reviewer: reviewer}
}

type reviewLessonWorker struct {
	river.WorkerDefaults[ReviewLessonArgs]
	reviewer LessonReviewer
}

// Timeout overrides River's one-minute default, which an LLM answer can exceed.
func (*reviewLessonWorker) Timeout(*river.Job[ReviewLessonArgs]) time.Duration {
	return llmTimeout
}

func (w *reviewLessonWorker) Work(ctx context.Context, job *river.Job[ReviewLessonArgs]) error {
	return cancelUnfixable(w.reviewer.ReviewLesson(ctx, job.Args.LessonInfoID))
}

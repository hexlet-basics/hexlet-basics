package jobs

import (
	"context"

	"github.com/riverqueue/river"
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

type reviewLessonWorker struct {
	river.WorkerDefaults[ReviewLessonArgs]
	reviewer LessonReviewer
}

func (w *reviewLessonWorker) Work(ctx context.Context, job *river.Job[ReviewLessonArgs]) error {
	return w.reviewer.ReviewLesson(ctx, job.Args.LessonInfoID)
}

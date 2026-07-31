package handlers

import (
	"context"

	"hexletbasics/ent"
)

// VersionBuildStarter is the atomic operation handlers need: create a course
// version and enqueue its loader job. Production performs both writes in one
// SQL transaction; tests use a recording adapter over their rollback-only DB.
type VersionBuildStarter interface {
	Start(ctx context.Context, courseID int) (*ent.CourseVersion, error)
}

// LessonReviewEnqueuer schedules the AI review jobs behind the admin review
// actions — one job per lesson version info (lessonreviews.Enqueuer in
// production, a recording adapter in tests).
type LessonReviewEnqueuer interface {
	EnqueueLessonReviews(ctx context.Context, lessonInfoIDs []int) error
}

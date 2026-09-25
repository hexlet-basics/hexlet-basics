package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/emailtokens"
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

// AccountEmailEnqueuer schedules a Magic Link or Password Reset email for the
// worker to render and send (accountemails.Enqueuer in production, a recording
// adapter in tests). A request repeated within the cooldown is skipped, not
// refused.
type AccountEmailEnqueuer interface {
	EnqueueAccountEmail(ctx context.Context, purpose emailtokens.Purpose, userID int, locale string) error
}

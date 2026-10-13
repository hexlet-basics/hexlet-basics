package handlers_test

import (
	"context"
	"testing"

	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessonreview"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// The three read-only admin lists over the lesson graph, and the action that
// enqueues a lesson's reviews, answered as a server error when the storage fails
// under them. Each list counts first and reads the page second, so a table gone
// fails the count and a column gone fails only the page read.
func TestAdminLessonListsAnswerAStorageFailureAsAServerError(t *testing.T) {
	lessons := func(ctx context.Context, h *testsupport.Harness) error {
		_, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{})
		return err
	}
	progress := func(ctx context.Context, h *testsupport.Harness) error {
		_, err := h.Client.AdminListLessonProgress(ctx, api.AdminListLessonProgressParams{})
		return err
	}
	reviews := func(ctx context.Context, h *testsupport.Harness) error {
		_, err := h.Client.AdminListCourseLessonReviews(ctx, api.AdminListCourseLessonReviewsParams{})
		return err
	}

	answerServerErrors(t, []failure{
		{"count the lessons", reads(courselessontranslation.Table), lessons},
		{"read the lessons' page", column(courselessontranslation.Table, courselessontranslation.FieldCreatedAt), lessons},
		{"count the progress rows", reads(lessonprogress.Table), progress},
		{"read the progress page", column(lessonprogress.Table, lessonprogress.FieldCreatedAt), progress},
		{"count the reviews", reads(courselessonreview.Table), reviews},
		{"read the reviews' page", column(courselessonreview.Table, courselessonreview.FieldCreatedAt), reviews},
	})
}

func TestAdminReviewCourseLessonAnswersAStorageFailureAsAServerError(t *testing.T) {
	review := func(ctx context.Context, h *testsupport.Harness) error {
		return h.Client.AdminReviewCourseLesson(ctx, api.AdminReviewCourseLessonParams{ID: 1001})
	}

	answerServerErrors(t, []failure{
		{"read the lesson", reads(courselesson.Table), review},
		{"read the lesson's infos", reads(courselessontranslation.Table), review},
		{"enqueue the reviews", writes("river_job", "INSERT"), review},
	})
}

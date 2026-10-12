package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// enqueuedLessonInfoIDs lists the lesson-info ids the handler asked the queue
// to review.
func enqueuedLessonInfoIDs(h *testsupport.Harness) []int {
	var ids []int
	for _, call := range h.Reviews.EnqueueLessonReviewsCalls() {
		ids = append(ids, call.LessonInfoIDs...)
	}
	return ids
}

func TestAdminReviewCourse(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// JavaScript (82481401): its CURRENT version's infos across every locale —
	// the en trio plus the ru hello-world info. The old-version info (2004)
	// must not be re-reviewed.
	err := h.Client.AdminReviewCourse(ctx, api.AdminReviewCourseParams{ID: 82481401})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	ids := enqueuedLessonInfoIDs(h)
	assert.ElementsMatch(t, []int{2001, 2002, 2003, 2005}, ids)
}

func TestAdminReviewCourseNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	err := h.Client.AdminReviewCourse(ctx, api.AdminReviewCourseParams{ID: 999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.Empty(t, h.Reviews.EnqueueLessonReviewsCalls())
}

func TestAdminReviewCourseLesson(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// hello-world (1001): every info of the lesson — current en, old-version
	// en, and ru — so historical summaries refresh too (legacy lesson.infos).
	err := h.Client.AdminReviewCourseLesson(ctx, api.AdminReviewCourseLessonParams{ID: 1001})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	ids := enqueuedLessonInfoIDs(h)
	assert.ElementsMatch(t, []int{2002, 2004, 2005}, ids)
}

func TestAdminReviewCourseLessonNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	err := h.Client.AdminReviewCourseLesson(ctx, api.AdminReviewCourseLessonParams{ID: 999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.Empty(t, h.Reviews.EnqueueLessonReviewsCalls())
}

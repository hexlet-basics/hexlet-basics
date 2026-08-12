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

// TestAdminListCourseLessons: the list is the en info of each lesson for the
// course's CURRENT version. The two negative fixtures (an old-version info and a
// ru info) must be filtered out, leaving the three current-en lessons.
func TestAdminListCourseLessons(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(3), page.Total, "old-version and ru infos must be excluded")
	require.Len(t, page.Items, 3)

	// Ordered by the Info id DESC (not the lesson id this row reports as `id`).
	// Fixture info ids: 2003(Variables/1002) > 2002(Hello/1001) > 2001(Strings/
	// 1003). So names come out Variables, Hello, Strings — and the reported ids
	// (1002, 1001, 1003) are NOT descending, proving we sort on the info id.
	assert.Equal(t, "Variables", page.Items[0].Name.Value)
	assert.Equal(t, "Hello, World!", page.Items[1].Name.Value)
	assert.Equal(t, "Strings", page.Items[2].Name.Value)

	// The reported id is the lesson id and the slug is joined from the lesson.
	assert.Equal(t, int32(1002), page.Items[0].ID)
	assert.Equal(t, "variables", page.Items[0].Slug)
	assert.Equal(t, "Storing values", page.Items[0].Description.Value)

	assert.Equal(t, int32(1001), page.Items[1].ID)
	assert.Equal(t, "hello-world", page.Items[1].Slug)
}

// TestAdminListLessonProgress: newest first, with course/lesson slugs and the
// lesson name enriched from the en info.
func TestAdminListLessonProgress(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListLessonProgress(ctx, api.AdminListLessonProgressParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(2), page.Total)
	require.Len(t, page.Items, 2)

	// id DESC: row 3002 first, then 3001.
	assert.Equal(t, int32(3002), page.Items[0].ID)
	assert.Equal(t, int32(3001), page.Items[1].ID)

	first := page.Items[0] // 3002: started, no messages, lesson "variables"
	assert.Equal(t, int32(1002), first.UserId)
	assert.Equal(t, api.EnrollmentStateStarted, first.State)
	assert.True(t, first.MessagesCount.Null, "blank messages_count should be null")
	assert.Equal(t, "javascript", first.CourseSlug)
	assert.Equal(t, "variables", first.CourseLessonSlug)
	assert.Equal(t, "Variables", first.CourseLessonName)

	second := page.Items[1] // 3001: finished, 5 messages, lesson "hello-world"
	assert.Equal(t, api.EnrollmentStateFinished, second.State)
	assert.Equal(t, int32(5), second.MessagesCount.Value)
	assert.Equal(t, "hello-world", second.CourseLessonSlug)
	assert.Equal(t, "Hello, World!", second.CourseLessonName)
}

// TestAdminListCourseLessonReviews: only en reviews with a non-empty summary
// (with_summary). The empty-summary and ru fixtures must be filtered out.
func TestAdminListCourseLessonReviews(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseLessonReviews(ctx, api.AdminListCourseLessonReviewsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(2), page.Total, "empty-summary and ru reviews must be excluded")
	require.Len(t, page.Items, 2)

	// id DESC: 4004 first, then 4001.
	assert.Equal(t, int32(4004), page.Items[0].ID)
	assert.Equal(t, int32(4001), page.Items[1].ID)

	first := page.Items[0] // 4004 -> lesson 1003 (strings, natural_order 3)
	assert.Equal(t, "en", first.Locale)
	assert.Equal(t, "Good grasp of concatenation overall.", first.Summary)
	assert.Equal(t, "strings", first.Slug)
	assert.Equal(t, int32(3), first.LessonNaturalOrder)
	assert.Equal(t, "javascript", first.CourseSlug)
	assert.Equal(t, int32(82481401), first.CourseId)
	assert.Equal(t, int32(1003), first.CourseLessonId)

	second := page.Items[1] // 4001 -> lesson 1001 (hello-world, natural_order 1)
	assert.Equal(t, "hello-world", second.Slug)
	assert.Equal(t, int32(1), second.LessonNaturalOrder)
}

// TestAdminListCourseLessonsPaginated exercises the paginated path with the
// filtered count intact.
func TestAdminListCourseLessonsPaginated(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(3), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 1) // 3 rows, page 2 of size 2 -> the 3rd row
	assert.Equal(t, "Strings", page.Items[0].Name.Value)
}

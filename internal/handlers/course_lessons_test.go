package handlers_test

import (
	"cmp"
	"context"
	"maps"
	"net/http"
	"slices"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// TestAdminListCourseLessons: the list is the en info of each lesson for its
// course's CURRENT version, across every course. The two negative fixtures (an
// old-version info and a ru info) must be filtered out of JavaScript's lessons,
// leaving its three current-en lessons; the other courses' starter lessons are
// listed beside them, so JavaScript's rows are picked out by slug.
func TestAdminListCourseLessons(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(len(page.Items)), page.Total)

	javascript := []string{"variables", "hello-world", "strings"}
	var items []api.CourseLessonListItem
	for _, item := range page.Items {
		assert.NotEqual(t, "Old Hello", item.Name.Value, "old-version infos must be excluded")
		assert.NotEqual(t, "Привет, мир!", item.Name.Value, "ru infos must be excluded")
		if slices.Contains(javascript, item.Slug) {
			items = append(items, item)
		}
	}
	require.Len(t, items, 3, "old-version and ru infos must be excluded")

	// Ordered by the Info id DESC (not the lesson id this row reports as `id`).
	infoOf := map[string]fixtures.Label{
		"variables":   fixtures.LessonVersionInfoJavascriptVariablesCurrentEn,
		"hello-world": fixtures.LessonVersionInfoJavascriptHelloWorldCurrentEn,
		"strings":     fixtures.LessonVersionInfoJavascriptStringsCurrentEn,
	}
	got := lo.Map(items, func(item api.CourseLessonListItem, _ int) string { return item.Slug })
	want := slices.SortedFunc(maps.Keys(infoOf), func(a, b string) int {
		return cmp.Compare(fixtures.ID(infoOf[b]), fixtures.ID(infoOf[a]))
	})
	assert.Equal(t, want, got)

	bySlug := lo.KeyBy(items, func(item api.CourseLessonListItem) string { return item.Slug })
	assert.Equal(t, "Variables", bySlug["variables"].Name.Value)
	assert.Equal(t, "Hello, World!", bySlug["hello-world"].Name.Value)
	assert.Equal(t, "Strings", bySlug["strings"].Name.Value)

	// The reported id is the lesson id and the slug is joined from the lesson.
	assert.Equal(t, fixtures.ID32(fixtures.LessonJavascriptVariables), bySlug["variables"].ID)
	assert.Equal(t, "Storing values", bySlug["variables"].Description.Value)
	assert.Equal(t, fixtures.ID32(fixtures.LessonJavascriptHelloWorld), bySlug["hello-world"].ID)
}

// TestAdminListLessonProgress: newest first, with course/lesson slugs and the
// lesson name enriched from the en info.
func TestAdminListLessonProgress(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListLessonProgress(ctx, api.AdminListLessonProgressParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(h.DB.LessonProgress.Query().CountX(ctx)), page.Total)
	require.Len(t, page.Items, int(page.Total))

	// id DESC: the newest row first.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
	byID := lo.KeyBy(page.Items, func(item api.LessonProgress) int32 { return item.ID })

	first := byID[fixtures.ID32(fixtures.LessonProgressStaffEditorJavascriptVariables)] // started, no messages, lesson "variables"
	assert.Equal(t, int32(fixtures.ID(fixtures.UserStaffEditor)), first.UserId)
	assert.Equal(t, api.EnrollmentStateStarted, first.State)
	assert.True(t, first.MessagesCount.Null, "blank messages_count should be null")
	assert.Equal(t, "javascript", first.CourseSlug)
	assert.Equal(t, "variables", first.CourseLessonSlug)
	assert.Equal(t, "Variables", first.CourseLessonName)

	second := byID[fixtures.ID32(fixtures.LessonProgressAdminLearnerJavascriptHelloWorld)] // finished, 5 messages, lesson "hello-world"
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

	// id DESC: the newest row first.
	assert.Greater(t, page.Items[0].ID, page.Items[1].ID)
	byID := lo.KeyBy(page.Items, func(item api.CourseLessonReview) int32 { return item.ID })

	first := byID[fixtures.ID32(fixtures.LessonReviewJavascriptStringsEn)] // strings, natural_order 3
	assert.Equal(t, "en", first.Locale)
	assert.Equal(t, "Good grasp of concatenation overall.", first.Summary)
	assert.Equal(t, "strings", first.Slug)
	assert.Equal(t, int32(3), first.LessonNaturalOrder)
	assert.Equal(t, "javascript", first.CourseSlug)
	assert.Equal(t, int32(fixtures.ID(fixtures.CourseJavascript)), first.CourseId)
	assert.Equal(t, int32(fixtures.ID(fixtures.LessonJavascriptStrings)), first.CourseLessonId)

	second := byID[fixtures.ID32(fixtures.LessonReviewJavascriptHelloWorldEn)] // hello-world, natural_order 1
	assert.Equal(t, "hello-world", second.Slug)
	assert.Equal(t, int32(1), second.LessonNaturalOrder)
}

// TestAdminListCourseLessonsPaginated exercises the paginated path with the
// filtered count intact.
func TestAdminListCourseLessonsPaginated(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	all, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	require.Greater(t, len(all.Items), 4)

	page, err := h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, all.Total, page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Equal(t, all.Items[2:4], page.Items) // page 2 of size 2 -> the 3rd and 4th rows
}

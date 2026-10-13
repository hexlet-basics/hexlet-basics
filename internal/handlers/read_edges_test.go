package handlers_test

import (
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/enrollment"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// A page past the last one is an empty page, not an error: the total still says
// how many posts there are.
func TestAdminListBlogPostsPastTheLastPage(t *testing.T) {
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListBlogPosts(t.Context(), api.AdminListBlogPostsParams{
		Page: api.NewOptInt32(999),
	})
	require.NoError(t, err)

	assert.Empty(t, page.Items)
	assert.Positive(t, page.Total)
	assert.Equal(t, int32(999), page.Page)
}

// A post whose state was never set reads as null, not as a made-up state.
func TestAdminGetBlogPostWithoutAState(t *testing.T) {
	h := testsupport.NewHarness(t)
	h.DB.BlogPost.Update().Where(blogpost.ID(6001)).ClearState().ExecX(t.Context())

	post, err := h.Client.AdminGetBlogPost(t.Context(), api.AdminGetBlogPostParams{ID: 6001})
	require.NoError(t, err)

	assert.True(t, post.State.Null)
}

// A course with no built version is still a page: no lessons, no modules and no
// position, rather than a failed read.
func TestGetCourseWithoutABuiltVersion(t *testing.T) {
	h := testsupport.NewHarness(t)
	h.DB.Course.Update().Where(course.SlugEQ(jsCourseSlug)).ClearCurrentVersionID().ExecX(t.Context())

	res, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view, ok := res.(*api.CourseView)
	require.True(t, ok, "got %T", res)

	assert.Empty(t, view.Lessons)
	assert.Empty(t, view.Modules)
}

// Without landing copy a course has no questions either.
func TestGetCourseWithoutALandingPageHasNoQuestions(t *testing.T) {
	h := testsupport.NewHarness(t)
	js := h.DB.Course.Query().Where(course.SlugEQ(jsCourseSlug)).OnlyX(t.Context())
	h.DB.CourseLandingPage.Delete().Where(courselandingpage.CourseID(js.ID)).ExecX(t.Context())

	res, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: jsCourseSlug})
	require.NoError(t, err)
	view := res.(*api.CourseView)

	assert.True(t, view.LandingPage.Null)
	assert.Empty(t, view.QnaItems)
}

// A lesson of an unknown course is the not-found answer, as is one of a course
// with nothing built: there is no content to serve.
func TestGetCourseLessonIsNotFoundWhereThereIsNothingToServe(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: "no-such-course",
		Slug:       firstLessonSlug,
	})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	h.DB.Course.Update().Where(course.SlugEQ(jsCourseSlug)).ClearCurrentVersionID().ExecX(t.Context())
	res, err = h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       firstLessonSlug,
	})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)

	res, err = h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       "no-such-lesson",
	})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
}

// A lesson whose version records no path in the course repository has no
// "propose a fix" link, rather than a broken one.
func TestGetCourseLessonWithoutASourcePathHasNoSourceLink(t *testing.T) {
	h := testsupport.NewHarness(t)
	h.DB.CourseLessonVersion.Update().ClearPathToCode().ExecX(t.Context())

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)
	view := res.(*api.CourseLessonView)

	assert.Empty(t, view.Lesson.SourceCodeUrl.Value)
}

// A signed-in learner who never started the course has no enrollment: the page
// says so with null, not with an empty record.
func TestGetCourseForALearnerWhoNeverStartedIt(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := t.Context()
	started := h.DB.Enrollment.Query().
		Where(enrollment.UserID(h.UserID)).
		Select(enrollment.FieldCourseID).
		IntsX(ctx)
	untouched := h.DB.Course.Query().Where(course.IDNotIn(started...)).FirstX(ctx)

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: lo.FromPtr(untouched.Slug)})
	require.NoError(t, err)
	view, ok := res.(*api.CourseView)
	require.True(t, ok, "got %T", res)

	assert.True(t, view.Enrollment.Null)
}

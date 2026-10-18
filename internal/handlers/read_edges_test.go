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
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// A page past the last one is an empty page, not an error: the total still says
// how many posts there are.
func TestAdminListBlogPostsPastTheLastPage(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.DB.BlogPost.Update().Where(blogpost.ID(fixture.ID(fixture.BlogPostHelloWorld))).ClearState().ExecX(t.Context())

	post, err := h.Client.AdminGetBlogPost(t.Context(), api.AdminGetBlogPostParams{ID: fixture.ID32(fixture.BlogPostHelloWorld)})
	require.NoError(t, err)

	assert.True(t, post.State.Null)
}

// A course with no built version is still a page: no lessons, no modules and no
// position, rather than a failed read.
func TestGetCourseWithoutABuiltVersion(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
	h.DB.Course.Update().Where(course.SlugEQ(jsCourseSlug)).ClearCurrentVersionID().ExecX(t.Context())

	res, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)

	view := res

	assert.Empty(t, view.Lessons)
	assert.Empty(t, view.Modules)
}

// A lesson of an unknown course is the not-found answer, as is one of a course
// with nothing built: there is no content to serve.
func TestGetCourseLessonIsNotFoundWhereThereIsNothingToServe(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	_, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: "no-such-course",
		Slug:       firstLessonSlug,
	})
	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	h.DB.Course.Update().Where(course.SlugEQ(jsCourseSlug)).ClearCurrentVersionID().ExecX(t.Context())
	_, err = h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       firstLessonSlug,
	})
	requireProblem(t, err, http.StatusNotFound)

	_, err = h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       "no-such-lesson",
	})
	requireProblem(t, err, http.StatusNotFound)
}

// A lesson whose version records no path in the course repository has no
// "propose a fix" link, rather than a broken one.
func TestGetCourseLessonWithoutASourcePathHasNoSourceLink(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.DB.CourseLessonVersion.Update().ClearPathToCode().ExecX(t.Context())

	res, err := h.Client.GetCourseLesson(t.Context(), api.GetCourseLessonParams{
		CourseSlug: jsCourseSlug,
		Slug:       secondLessonSlug,
	})
	require.NoError(t, err)

	view := res

	assert.Empty(t, view.Lesson.SourceCodeUrl.Value)
}

// A signed-in learner who never started the course has no enrollment: the page
// says so with null, not with an empty record.
func TestGetCourseForALearnerWhoNeverStartedIt(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	ctx := t.Context()
	started := h.DB.Enrollment.Query().
		Where(enrollment.UserID(h.UserID)).
		Select(enrollment.FieldCourseID).
		IntsX(ctx)
	untouched := h.DB.Course.Query().Where(course.IDNotIn(started...)).FirstX(ctx)

	landing := h.DB.CourseLandingPage.Query().
		Where(courselandingpage.CourseID(untouched.ID), courselandingpage.LocaleEQ("ru")).
		FirstX(ctx)

	res, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: lo.FromPtr(landing.Slug)})
	require.NoError(t, err)

	view := res

	assert.True(t, view.Enrollment.Null)
}

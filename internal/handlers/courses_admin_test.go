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

// languages.yml seeds this many courses; ruby is one of them.
const (
	totalCourses  = 13
	courseRubyIDA = 207281424
)

func TestAdminListCourses(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourses(ctx, api.AdminListCoursesParams{PerPage: api.NewOptInt32(100)})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(totalCourses), page.Total)

	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminGetCourse(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.AdminGetCourse(ctx, api.AdminGetCourseParams{ID: courseRubyIDA})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "ruby", got.Slug)
	// repositoryUrl is derived from the slug.
	assert.Equal(t, "https://github.com/hexlet-basics/exercises-ruby", got.RepositoryUrl.Value)
	// ruby's current version is embedded (fixtures link current_version_id).
	assert.False(t, got.CurrentVersion.Null)
}

func TestAdminGetCourseNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetCourse(ctx, api.AdminGetCourseParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateCourse(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	created, err := h.Client.AdminCreateCourse(ctx, &api.CourseInput{
		Slug:                     api.NewNilString("haskell"),
		LearnAs:                  api.NewNilCourseLearnAs(api.CourseLearnAsSecondLanguage),
		Progress:                 api.NewNilCourseReadiness(api.CourseReadinessInDevelopment),
		HexletProgramLandingPage: api.NilString{Null: true},
		RepositoryUrl:            api.NilString{Null: true},
		CoverAttachmentId:        api.NilInt32{Null: true},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.NotZero(t, created.ID)
	assert.Equal(t, "haskell", created.Slug)
	assert.Equal(t, api.CourseLearnAsSecondLanguage, created.LearnAs.Value)
	assert.Equal(t, api.CourseReadinessInDevelopment, created.Progress.Value)
	// repositoryUrl derives from the new slug.
	assert.Equal(t, "https://github.com/hexlet-basics/exercises-haskell", created.RepositoryUrl.Value)

	total, err := h.DB.Course.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalCourses+1, total)
}

// TestAdminCreateCourseDuplicateSlug proves the languages_slug_index uniqueness:
// a slug already present (elixir) surfaces as 409 via the central handler.
func TestAdminCreateCourseDuplicateSlug(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateCourse(ctx, &api.CourseInput{
		Slug:                     api.NewNilString("elixir"), // already seeded
		LearnAs:                  api.NilCourseLearnAs{Null: true},
		Progress:                 api.NilCourseReadiness{Null: true},
		HexletProgramLandingPage: api.NilString{Null: true},
		RepositoryUrl:            api.NilString{Null: true},
		CoverAttachmentId:        api.NilInt32{Null: true},
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestAdminUpdateCourse(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	updated, err := h.Client.AdminUpdateCourse(ctx, &api.CourseInput{
		Slug:                     api.NewNilString("ruby"),
		LearnAs:                  api.NewNilCourseLearnAs(api.CourseLearnAsFirstLanguage),
		Progress:                 api.NewNilCourseReadiness(api.CourseReadinessCompleted),
		HexletProgramLandingPage: api.NilString{Null: true}, // clears
		RepositoryUrl:            api.NilString{Null: true},
		CoverAttachmentId:        api.NilInt32{Null: true},
	}, api.AdminUpdateCourseParams{ID: courseRubyIDA})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, api.CourseReadinessCompleted, updated.Progress.Value)
	assert.Equal(t, api.CourseLearnAsFirstLanguage, updated.LearnAs.Value)

	row, err := h.DB.Course.Get(ctx, courseRubyIDA)
	require.NoError(t, err)
	assert.Nil(t, row.HexletProgramLandingPage, "null hexletProgramLandingPage should clear the column")
}

package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/testsupport"
)

func newServer(t *testing.T) *handlers.Server {
	t.Helper()
	return handlers.NewServer(testsupport.NewClient(t))
}

func TestAdminListCourseCategories(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{})
	require.NoError(t, err)

	assert.Equal(t, int32(3), page.Total)
	assert.Len(t, page.Items, 3)
	// Newest first (id desc): fixture id 3 leads.
	assert.Equal(t, int32(3), page.Items[0].ID)
}

func TestAdminListCourseCategoriesPaginated(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(3), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 1) // 3 rows, page 2 of size 2 -> 1 row
	assert.Equal(t, int32(1), page.Items[0].ID)
}

func TestAdminGetCourseCategory(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	got, err := srv.AdminGetCourseCategory(ctx, api.AdminGetCourseCategoryParams{ID: 1})
	require.NoError(t, err)

	assert.Equal(t, int32(1), got.ID)
	assert.Equal(t, "Programming", got.Name.Value)
	assert.Equal(t, "programming-en", got.Slug.Value)
}

func TestAdminCreateCourseCategory(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	res, err := srv.AdminCreateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Mobile",
		Header:      "Mobile courses",
		Slug:        "mobile-en",
		Description: api.NewNilString("Build apps"),
	})
	require.NoError(t, err)

	created, ok := res.(*api.CourseCategory)
	require.True(t, ok, "expected a created category, got %T", res)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "Mobile", created.Name.Value)
	assert.False(t, created.CreatedAt.IsZero())

	// It is now listable.
	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{})
	require.NoError(t, err)
	assert.Equal(t, int32(4), page.Total)
}

func TestAdminCreateCourseCategoryDuplicateName(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	res, err := srv.AdminCreateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming", // fixture id 1
		Header:      "Some other header",
		Slug:        "some-other-slug",
		Description: api.NilString{Null: true},
	})
	require.NoError(t, err)

	verr, ok := res.(*api.ValidationError)
	require.True(t, ok, "expected a validation error, got %T", res)
	assert.Contains(t, verr.Errors, "name")
}

func TestAdminUpdateCourseCategory(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	res, err := srv.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming Updated",
		Header:      "Programming courses",
		Slug:        "programming-en",
		Description: api.NilString{Null: true}, // clears the column
	}, api.AdminUpdateCourseCategoryParams{ID: 1})
	require.NoError(t, err)

	updated, ok := res.(*api.CourseCategory)
	require.True(t, ok, "expected an updated category, got %T", res)
	assert.Equal(t, "Programming Updated", updated.Name.Value)
	assert.True(t, updated.Description.Null, "null description should clear the column")
}

func TestAdminUpdateCourseCategoryDuplicateName(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	// Rename id 1 to id 2's name -> conflict.
	res, err := srv.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Frontend", // fixture id 2
		Header:      "Programming courses",
		Slug:        "programming-en",
		Description: api.NilString{Null: true},
	}, api.AdminUpdateCourseCategoryParams{ID: 1})
	require.NoError(t, err)

	verr, ok := res.(*api.ValidationError)
	require.True(t, ok, "expected a validation error, got %T", res)
	assert.Contains(t, verr.Errors, "name")
}

func TestAdminUpdateCourseCategoryKeepsOwnValues(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	// Updating a record with its own unchanged unique values must not conflict.
	res, err := srv.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming",
		Header:      "Programming courses",
		Slug:        "programming-en",
		Description: api.NewNilString("Now with a description"),
	}, api.AdminUpdateCourseCategoryParams{ID: 1})
	require.NoError(t, err)

	_, ok := res.(*api.CourseCategory)
	require.True(t, ok, "expected an updated category, got %T", res)
}

func TestAdminDeleteCourseCategory(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	err := srv.AdminDeleteCourseCategory(ctx, api.AdminDeleteCourseCategoryParams{ID: 3})
	require.NoError(t, err)

	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{})
	require.NoError(t, err)
	assert.Equal(t, int32(2), page.Total)
}

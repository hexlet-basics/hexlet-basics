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

// categoryBySlug resolves a fixture category by its stable business key. Fixture
// ids are Rails crc32 ids (large, not hand-picked), so tests key off slug and
// read the id back rather than hard-coding it.
func categoryBySlug(t *testing.T, srv *handlers.Server, slug string) api.CourseCategory {
	t.Helper()
	page, err := srv.AdminListCourseCategories(context.Background(), api.AdminListCourseCategoriesParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	for _, it := range page.Items {
		if it.Slug.Value == slug {
			return it
		}
	}
	t.Fatalf("category with slug %q not found", slug)
	return api.CourseCategory{}
}

// legacy language/categories.yml seeds these six.
const totalCategories = 6

func TestAdminListCourseCategories(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{})
	require.NoError(t, err)

	assert.Equal(t, int32(totalCategories), page.Total)
	assert.Len(t, page.Items, totalCategories)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminListCourseCategoriesPaginated(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(totalCategories), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 2) // 6 rows, page 2 of size 2 -> 2nd pair
}

func TestAdminGetCourseCategory(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	want := categoryBySlug(t, srv, "programming-en")

	res, err := srv.AdminGetCourseCategory(ctx, api.AdminGetCourseCategoryParams{ID: want.ID})
	require.NoError(t, err)

	got, ok := res.(*api.CourseCategory)
	require.True(t, ok, "expected a category, got %T", res)
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, "Programming", got.Name.Value)
	assert.Equal(t, "programming-en", got.Slug.Value)
}

func TestAdminGetCourseCategoryNotFound(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	res, err := srv.AdminGetCourseCategory(ctx, api.AdminGetCourseCategoryParams{ID: 999999999})
	require.NoError(t, err)

	_, ok := res.(*api.NotFoundError)
	require.True(t, ok, "expected a not-found error, got %T", res)
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
	assert.Equal(t, int32(totalCategories+1), page.Total)
}

func TestAdminCreateCourseCategoryDuplicateName(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	res, err := srv.AdminCreateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming", // programming-en's name
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

	target := categoryBySlug(t, srv, "programming-en")

	res, err := srv.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming Updated",
		Header:      target.Header.Value,
		Slug:        target.Slug.Value,
		Description: api.NilString{Null: true}, // clears the column
	}, api.AdminUpdateCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)

	updated, ok := res.(*api.CourseCategory)
	require.True(t, ok, "expected an updated category, got %T", res)
	assert.Equal(t, "Programming Updated", updated.Name.Value)
	assert.True(t, updated.Description.Null, "null description should clear the column")
}

func TestAdminUpdateCourseCategoryDuplicateName(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	target := categoryBySlug(t, srv, "programming-en")

	// Rename to another category's name -> conflict.
	res, err := srv.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Frontend", // frontend-en's name
		Header:      target.Header.Value,
		Slug:        target.Slug.Value,
		Description: api.NilString{Null: true},
	}, api.AdminUpdateCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)

	verr, ok := res.(*api.ValidationError)
	require.True(t, ok, "expected a validation error, got %T", res)
	assert.Contains(t, verr.Errors, "name")
}

func TestAdminUpdateCourseCategoryKeepsOwnValues(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	target := categoryBySlug(t, srv, "programming-en")

	// Updating a record with its own unchanged unique values must not conflict.
	res, err := srv.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        target.Name.Value,
		Header:      target.Header.Value,
		Slug:        target.Slug.Value,
		Description: api.NewNilString("Now with a description"),
	}, api.AdminUpdateCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)

	_, ok := res.(*api.CourseCategory)
	require.True(t, ok, "expected an updated category, got %T", res)
}

func TestAdminDeleteCourseCategory(t *testing.T) {
	srv := newServer(t)
	ctx := context.Background()

	target := categoryBySlug(t, srv, "layouting-en")

	err := srv.AdminDeleteCourseCategory(ctx, api.AdminDeleteCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)

	page, err := srv.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{})
	require.NoError(t, err)
	assert.Equal(t, int32(totalCategories-1), page.Total)
}

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

// categoryBySlug resolves a fixture category by its stable business key. Fixture
// ids are Rails crc32 ids (large, not hand-picked), so tests key off slug and
// read the id back rather than hard-coding it.
func categoryBySlug(t *testing.T, h *testsupport.Harness, slug string) api.CourseCategory {
	t.Helper()
	page, err := h.Client.AdminListCourseCategories(context.Background(), api.AdminListCourseCategoriesParams{
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
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalCategories), page.Total)
	assert.Len(t, page.Items, totalCategories)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminListCourseCategoriesPaginated(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseCategories(ctx, api.AdminListCourseCategoriesParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(totalCategories), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 2) // 6 rows, page 2 of size 2 -> 2nd pair
}

func TestAdminGetCourseCategory(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	want := categoryBySlug(t, h, "programming-en")

	got, err := h.Client.AdminGetCourseCategory(ctx, api.AdminGetCourseCategoryParams{ID: want.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, "Programming", got.Name.Value)
	assert.Equal(t, "programming-en", got.Slug.Value)
}

func TestAdminGetCourseCategoryNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// A missing id is an ent not-found error, mapped to 404 by the central
	// ErrorHandler; the client surfaces that undeclared status as an error.
	_, err := h.Client.AdminGetCourseCategory(ctx, api.AdminGetCourseCategoryParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateCourseCategory(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	created, err := h.Client.AdminCreateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Mobile",
		Header:      "Mobile courses",
		Slug:        "mobile-en",
		Description: api.NewNilString("Build apps"),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	assert.NotZero(t, created.ID)
	assert.Equal(t, "Mobile", created.Name.Value)
	assert.False(t, created.CreatedAt.IsZero())

	// It landed in the database (same transaction the handler wrote through).
	total, err := h.DB.CourseCategory.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalCategories+1, total)
}

func TestAdminCreateCourseCategoryDuplicateName(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// "Programming" is programming-en's name; the DB unique index rejects it,
	// surfaced as ent.IsConstraintError -> 409.
	_, err := h.Client.AdminCreateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming",
		Header:      "Some other header",
		Slug:        "some-other-slug",
		Description: api.NilString{Null: true},
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())

	// No DB assertion after the conflict on purpose: the failed INSERT aborts the
	// surrounding test transaction (Postgres 25P02), so any further query in it
	// errors. In production each request inserts in autocommit, so the failure is
	// isolated and no row persists — the 409 is the behavioural proof of that.
}

// TestAdminCreateCourseCategoryDuplicateSlug is the parity proof for the slug
// index. The input carries no locale, so new rows get locale=NULL; the fixtures
// are all en/ru, so a new (slug, NULL) can never collide with them. The ONLY way
// two rows conflict on slug here is the index's NULLS NOT DISTINCT clause, which
// reproduces Rails' `uniqueness: { scope: :locale }` for a nil locale. Without
// it this second create would wrongly succeed.
func TestAdminCreateCourseCategoryDuplicateSlug(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	newInput := func(name, header string) *api.CourseCategoryInput {
		return &api.CourseCategoryInput{
			Name:        name,
			Header:      header,
			Slug:        "dup-slug", // same slug, both with a null locale
			Description: api.NilString{Null: true},
		}
	}

	_, err := h.Client.AdminCreateCourseCategory(ctx, newInput("Name One", "Header One"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	// Distinct name and header, so only the slug can be the conflict.
	_, err = h.Client.AdminCreateCourseCategory(ctx, newInput("Name Two", "Header Two"))
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestAdminUpdateCourseCategory(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := categoryBySlug(t, h, "programming-en")

	updated, err := h.Client.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Programming Updated",
		Header:      target.Header.Value,
		Slug:        target.Slug.Value,
		Description: api.NilString{Null: true}, // clears the column
	}, api.AdminUpdateCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Programming Updated", updated.Name.Value)

	// A null description clears the column in the database.
	row, err := h.DB.CourseCategory.Get(ctx, int(target.ID))
	require.NoError(t, err)
	assert.Equal(t, "Programming Updated", *row.Name)
	assert.Nil(t, row.Description, "null description should clear the column")
}

func TestAdminUpdateCourseCategoryDuplicateName(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := categoryBySlug(t, h, "programming-en")

	// Rename to another category's name -> unique-index conflict -> 409.
	_, err := h.Client.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        "Frontend", // frontend-en's name
		Header:      target.Header.Value,
		Slug:        target.Slug.Value,
		Description: api.NilString{Null: true},
	}, api.AdminUpdateCourseCategoryParams{ID: target.ID})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestAdminUpdateCourseCategoryKeepsOwnValues(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := categoryBySlug(t, h, "programming-en")

	// Updating a record with its own unchanged unique values must not conflict:
	// the row excludes itself from the unique index at the same key.
	_, err := h.Client.AdminUpdateCourseCategory(ctx, &api.CourseCategoryInput{
		Name:        target.Name.Value,
		Header:      target.Header.Value,
		Slug:        target.Slug.Value,
		Description: api.NewNilString("Now with a description"),
	}, api.AdminUpdateCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
}

func TestAdminDeleteCourseCategory(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := categoryBySlug(t, h, "layouting-en")

	err := h.Client.AdminDeleteCourseCategory(ctx, api.AdminDeleteCourseCategoryParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.CourseCategory.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalCategories-1, total)
}

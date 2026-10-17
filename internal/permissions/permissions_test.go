package permissions_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/permissions"
)

func TestAnOperationReportsTheDeclaredPermissionAndKind(t *testing.T) {
	read, ok := permissions.For("AdminListBlogPosts")
	require.True(t, ok)
	assert.Equal(t, permissions.Requirement{Permission: "blog_posts.read", Kind: permissions.Read}, read)

	write, ok := permissions.For("AdminSuggestBlogPostRelatedCourses")
	require.True(t, ok)
	assert.Equal(t, permissions.Requirement{Permission: "blog_posts.update", Kind: permissions.Write}, write)
}

// Each admin operation that used to be opened by a bit of the retired
// create/read/update/delete matrix declares the Permission of that bit, with the
// legacy resource names replaced by the domain terms.
func TestEachRetiredMatrixBitBecameAPermissionOfTheSameOperations(t *testing.T) {
	for permission, operations := range map[string][]string{
		"banners.read":                {"AdminListBanners", "AdminGetBanner"},
		"banners.create":              {"AdminCreateBanner"},
		"banners.update":              {"AdminUpdateBanner"},
		"banners.delete":              {"AdminDeleteBanner"},
		"reviews.read":                {"AdminListReviews", "AdminGetReview"},
		"reviews.create":              {"AdminCreateReview"},
		"reviews.update":              {"AdminUpdateReview"},
		"reviews.delete":              {"AdminDeleteReview"},
		"leads.read":                  {"AdminListLeads"},
		"messages.read":               {"AdminListMessages"},
		"course_categories.read":      {"AdminListCourseCategories", "AdminGetCourseCategory"},
		"course_categories.create":    {"AdminCreateCourseCategory"},
		"course_categories.update":    {"AdminUpdateCourseCategory"},
		"course_categories.delete":    {"AdminDeleteCourseCategory"},
		"courses.read":                {"AdminListCourses", "AdminGetCourse"},
		"courses.create":              {"AdminCreateCourse"},
		"courses.update":              {"AdminUpdateCourse", "AdminReviewCourse"},
		"course_lessons.read":         {"AdminListCourseLessons"},
		"course_lessons.update":       {"AdminReviewCourseLesson"},
		"course_lesson_reviews.read":  {"AdminListCourseLessonReviews"},
		"lesson_progress.read":        {"AdminListLessonProgress"},
		"course_landing_pages.read":   {"AdminListCourseLandingPages", "AdminGetCourseLandingPage"},
		"course_landing_pages.create": {"AdminCreateCourseLandingPage"},
		"course_landing_pages.update": {"AdminUpdateCourseLandingPage"},
		"course_landing_pages.delete": {"AdminDeleteCourseLandingPage"},
	} {
		wantKind := permissions.Write
		if strings.HasSuffix(permission, ".read") {
			wantKind = permissions.Read
		}
		for _, operation := range operations {
			got, ok := permissions.For(operation)
			require.True(t, ok, operation)
			assert.Equal(t, permissions.Requirement{Permission: permission, Kind: wantKind}, got, operation)
		}
	}
}

func TestAnOperationWithoutADeclarationHasNoRequirement(t *testing.T) {
	_, ok := permissions.For("AdminListUsers")
	assert.False(t, ok)
}

func TestTheCatalogListsEachDeclaredPermissionOnceInOrder(t *testing.T) {
	catalog := permissions.Catalog()

	names := lo.Map(catalog, func(p permissions.Permission, _ int) string { return p.Name })
	assert.True(t, slices.IsSorted(names), "ordered by name")
	assert.Equal(t, lo.Uniq(names), names, "each Permission once")
	assert.Subset(t, catalog, []permissions.Permission{
		{Name: "blog_posts.create", Kind: permissions.Write},
		{Name: "blog_posts.delete", Kind: permissions.Write},
		{Name: "blog_posts.read", Kind: permissions.Read},
		{Name: "blog_posts.update", Kind: permissions.Write},
		{Name: "banners.read", Kind: permissions.Read},
	})
}

func TestAnOperationReportsTheAccessItDeclaresInsteadOfAPermission(t *testing.T) {
	assert.Equal(t, permissions.Administrators, permissions.AccessOf("AdminListRoles"))
	assert.Equal(t, permissions.AnyStaff, permissions.AccessOf("AdminUploadAttachment"))
	assert.Empty(t, permissions.AccessOf("AdminListBanners"), "a Permission is not an Access")

	assert.True(t, permissions.IsClassified("AdminListBanners"))
	assert.True(t, permissions.IsClassified("AdminListRoles"))
	assert.False(t, permissions.IsClassified("AdminInventedOperation"))
}

func TestKnownTellsCatalogNamesFromInventedOnes(t *testing.T) {
	assert.True(t, permissions.Known("blog_posts.delete"))
	assert.False(t, permissions.Known("blog_posts.publish"))
	assert.False(t, permissions.Known(""))
}

func TestMalformedDeclarationsAreRejected(t *testing.T) {
	for name, doc := range map[string]string{
		"not yaml": "paths: [",
		"unknown kind": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {name: x.read, kind: peek}`,
		"permission declared with two kinds": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {name: x.use, kind: read}
    post:
      operationId: adminCreateX
      x-permission: {name: x.use, kind: write}`,
		"unknown access": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-access: everyone`,
		"permission and access together": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {name: x.read, kind: read}
      x-access: administrators`,
		"declaration without a name": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {kind: read}`,
	} {
		_, err := permissions.Parse([]byte(doc))
		assert.Error(t, err, name)
	}
}

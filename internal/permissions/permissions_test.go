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
	t.Parallel()
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
	t.Parallel()
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
		"courses.update":              {"AdminUpdateCourse"},
		"courses.review":              {"AdminReviewCourse"},
		"course_lessons.read":         {"AdminListCourseLessons"},
		"course_lessons.review":       {"AdminReviewCourseLesson"},
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

// The operations that were administrators-only declare Permissions a role can
// hold, named with the GLOSSARY terms.
func TestFormerAdministratorsOnlyOperationsDeclarePermissions(t *testing.T) {
	t.Parallel()
	for permission, operations := range map[string][]string{
		"users.read":                           {"AdminListUsers", "AdminSearchUsers", "AdminGetUser", "AdminListManagementUsers", "AdminGetManagementUser"},
		"users.create":                         {"AdminCreateUser"},
		"users.update":                         {"AdminUpdateUser", "AdminUpdateManagementUser"},
		"users.delete":                         {"AdminDeleteUser"},
		"staff_members.read":                   {"AdminListStaffMembers", "AdminGetStaffMember"},
		"staff_members.create":                 {"AdminCreateStaffMember"},
		"staff_members.update":                 {"AdminUpdateStaffMember"},
		"staff_members.delete":                 {"AdminDeleteStaffMember"},
		"staff_roles.read":                     {"AdminListRoles", "AdminGetRole", "AdminListPermissions", "AdminGetRolePermissions"},
		"staff_roles.create":                   {"AdminCreateRole"},
		"staff_roles.update":                   {"AdminUpdateRole", "AdminUpdateRolePermissions"},
		"staff_roles.delete":                   {"AdminDeleteRole"},
		"course_landing_page_qna_items.read":   {"AdminListLandingPageQnaItems"},
		"course_landing_page_qna_items.create": {"AdminCreateLandingPageQnaItem"},
		"course_landing_page_qna_items.update": {"AdminUpdateLandingPageQnaItem"},
		"course_landing_page_qna_items.delete": {"AdminDeleteLandingPageQnaItem"},
		"course_category_qna_items.read":       {"AdminListCategoryQnaItems"},
		"course_category_qna_items.create":     {"AdminCreateCategoryQnaItem"},
		"course_category_qna_items.update":     {"AdminUpdateCategoryQnaItem"},
		"course_category_qna_items.delete":     {"AdminDeleteCategoryQnaItem"},
		"course_versions.create":               {"AdminCreateCourseVersion"},
	} {
		wantKind := permissions.Write
		if strings.HasSuffix(permission, ".read") {
			wantKind = permissions.Read
		}
		for _, operation := range operations {
			got, ok := permissions.For(operation)
			require.True(t, ok, operation)
			assert.Equal(t, permissions.Requirement{Permission: permission, Kind: wantKind}, got, operation)
			assert.Empty(t, permissions.AccessOf(operation), operation)
		}
	}
}

func TestAnOperationWithoutADeclarationHasNoRequirement(t *testing.T) {
	t.Parallel()
	_, ok := permissions.For("AdminInventedOperation")
	assert.False(t, ok)
}

// The capability Permissions are listed once in the contract, which no operation
// requires yet: they are in the catalog with their kinds, and are told apart
// from the operation Permissions.
func TestTheCatalogHoldsTheCapabilityPermissions(t *testing.T) {
	t.Parallel()
	assert.Subset(t, permissions.Catalog(), []permissions.Permission{
		{Name: "mcp.read", Kind: permissions.Read, Capability: true},
		{Name: "mcp.write", Kind: permissions.Write, Capability: true},
		{Name: "mcp.manage", Kind: permissions.Write, Capability: true},
		{Name: "audit.read", Kind: permissions.Read, Capability: true},
	})
	for _, name := range []string{"mcp.read", "mcp.write", "mcp.manage", "audit.read"} {
		assert.True(t, permissions.IsCapability(name), name)
		assert.True(t, permissions.Known(name), name)
	}
	assert.False(t, permissions.IsCapability("blog_posts.read"))
	assert.False(t, permissions.IsCapability("mcp.invented"))
}

func TestTheCatalogListsEachDeclaredPermissionOnceInOrder(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	assert.Empty(t, permissions.AccessOf("AdminListRoles"), "no operation is administrators-only any more")
	assert.Equal(t, permissions.AnyStaff, permissions.AccessOf("AdminUploadAttachment"))
	assert.Empty(t, permissions.AccessOf("AdminListBanners"), "a Permission is not an Access")

	assert.True(t, permissions.IsClassified("AdminListBanners"))
	assert.True(t, permissions.IsClassified("AdminListRoles"))
	assert.False(t, permissions.IsClassified("AdminInventedOperation"))
}

func TestKnownTellsCatalogNamesFromInventedOnes(t *testing.T) {
	t.Parallel()
	assert.True(t, permissions.Known("blog_posts.delete"))
	assert.True(t, permissions.Known("blog_posts.publish"))
	assert.False(t, permissions.Known("blog_posts.frobnicate"))
	assert.False(t, permissions.Known(""))
}

func TestMalformedDeclarationsAreRejected(t *testing.T) {
	t.Parallel()
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
		"the retired administrators access": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-access: administrators`,
		"capability without a kind": `
x-capabilities:
  - {name: mcp.read}`,
		"capability listed twice": `
x-capabilities:
  - {name: mcp.read, kind: read}
  - {name: mcp.read, kind: read}`,
		"capability and operation disagree on the kind": `
x-capabilities:
  - {name: audit.read, kind: read}
paths:
  /api/admin/x:
    post:
      operationId: adminCreateX
      x-permission: {name: audit.read, kind: write}`,
		"permission and access together": `
paths:
  /api/admin/x:
    get:
      operationId: adminListX
      x-permission: {name: x.read, kind: read}
      x-access: any-staff`,
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

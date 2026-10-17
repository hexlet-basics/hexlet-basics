package staffroles

import (
	"strings"

	"hexletbasics/internal/api"
)

type action int

const (
	index action = iota
	create
	update
	destroy
)

// rule is what an operation needs from a staff member: a bit of the matrix on
// one resource, or just being staff.
type rule struct {
	resource api.PermissionResource
	action   action
	anyStaff bool
}

// actionOf is the legacy controller-action → matrix-bit mapping: index and show
// read, new and create create, edit, update, review and related_courses update.
func actionOf(operation string) action {
	switch {
	case strings.HasPrefix(operation, "AdminCreate"):
		return create
	case strings.HasPrefix(operation, "AdminDelete"):
		return destroy
	case strings.HasPrefix(operation, "AdminList"), strings.HasPrefix(operation, "AdminGet"):
		return index
	default:
		return update
	}
}

// resourceOperations lists, per resource, the operations a staff member holding
// the matching bit may perform (legacy `STAFF_RESOURCE` of each admin
// controller). Operations that declare a Permission in the contract (blog posts
// so far, ADR-0017) are not listed: the declaration is their rule, and each
// resource leaves this table as it is migrated.
var resourceOperations = map[api.PermissionResource][]string{
	api.PermissionResourceBanners: {
		"AdminListBanners", "AdminGetBanner", "AdminCreateBanner", "AdminUpdateBanner", "AdminDeleteBanner",
	},
	api.PermissionResourceReviews: {
		"AdminListReviews", "AdminGetReview", "AdminCreateReview", "AdminUpdateReview", "AdminDeleteReview",
	},
	api.PermissionResourceLeads:    {"AdminListLeads"},
	api.PermissionResourceMessages: {"AdminListMessages"},
	api.PermissionResourceLanguageCategories: {
		"AdminListCourseCategories", "AdminGetCourseCategory", "AdminCreateCourseCategory",
		"AdminUpdateCourseCategory", "AdminDeleteCourseCategory",
	},
	api.PermissionResourceLanguages: {
		"AdminListCourses", "AdminGetCourse", "AdminCreateCourse", "AdminUpdateCourse", "AdminReviewCourse",
	},
	api.PermissionResourceLanguageLessons:       {"AdminListCourseLessons", "AdminReviewCourseLesson"},
	api.PermissionResourceLanguageLessonReviews: {"AdminListCourseLessonReviews"},
	api.PermissionResourceLanguageLessonMembers: {"AdminListLessonProgress"},
	api.PermissionResourceLanguageLandingPages: {
		"AdminListCourseLandingPages", "AdminGetCourseLandingPage", "AdminCreateCourseLandingPage",
		"AdminUpdateCourseLandingPage", "AdminDeleteCourseLandingPage",
	},
}

// anyStaffOperations need no matrix bit: the upload every form uses, like
// legacy's dashboard (`STAFF_RESOURCE = nil`).
var anyStaffOperations = []string{"AdminUploadAttachment"}

// adminOnly are the admin operations no matrix bit opens: users, staff, roles
// and the pieces legacy left to administrators (no `STAFF_RESOURCE`).
var adminOnly = map[string]bool{
	"AdminListUsers": true, "AdminGetUser": true, "AdminCreateUser": true, "AdminUpdateUser": true,
	"AdminDeleteUser": true, "AdminSearchUsers": true,
	"AdminListManagementUsers": true, "AdminGetManagementUser": true, "AdminUpdateManagementUser": true,
	"AdminListStaffMembers": true, "AdminGetStaffMember": true, "AdminCreateStaffMember": true,
	"AdminUpdateStaffMember": true, "AdminDeleteStaffMember": true,
	"AdminListPermissions": true,
	"AdminListRoles":       true, "AdminGetRole": true, "AdminCreateRole": true, "AdminUpdateRole": true,
	"AdminDeleteRole": true, "AdminGetRolePermissions": true, "AdminUpdateRolePermissions": true,
	"AdminCreateCourseVersion":  true,
	"AdminListCategoryQnaItems": true, "AdminCreateCategoryQnaItem": true, "AdminUpdateCategoryQnaItem": true,
	"AdminDeleteCategoryQnaItem":   true,
	"AdminListLandingPageQnaItems": true, "AdminCreateLandingPageQnaItem": true,
	"AdminUpdateLandingPageQnaItem": true, "AdminDeleteLandingPageQnaItem": true,
}

// rules is the operation → requirement table, built from the lists above.
var rules = buildRules()

func buildRules() map[string]rule {
	table := map[string]rule{}
	for resource, operations := range resourceOperations {
		for _, operation := range operations {
			table[operation] = rule{resource: resource, action: actionOf(operation)}
		}
	}
	for _, operation := range anyStaffOperations {
		table[operation] = rule{anyStaff: true}
	}
	return table
}

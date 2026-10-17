package staffroles

// anyStaffOperations need no Permission: the upload every form uses, like
// legacy's dashboard (`STAFF_RESOURCE = nil`).
var anyStaffOperations = []string{"AdminUploadAttachment"}

// adminOnly are the admin operations no Permission opens: users, staff, roles
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

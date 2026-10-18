package permission_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"hexletbasics/ent/staffmember"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/ent/user"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/permissioncatalog"
	"hexletbasics/internal/testsupport"
)

// The fixtures seed user_staff_editor as an editor holding the blog_posts.read,
// banners.read and banners.create grants, in ru and en; user_admin_learner is an administrator
// without a staff record, user_admin_staff one with it.
var (
	adminID  = fixture.ID(fixture.UserAdminLearner)
	editorID = fixture.ID(fixture.UserStaffEditor)
	// An administrator who is also staff, on the editor role: the harness user
	// has no staff record, and a staff record cannot be made for oneself.
	adminStaffID = fixture.ID(fixture.UserAdminStaff)
)

func TestPolicyFollowsTheRolesGrants(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := permission.NewPolicy(db)

	for name, c := range map[string]struct {
		user      int
		locale    string
		operation string
		allowed   bool
	}{
		"administrator may do anything":       {adminID, "es", "AdminDeleteUser", true},
		"granted bit":                         {editorID, "ru", "AdminListBlogPosts", true},
		"another grant of the same resource":  {editorID, "en", "AdminCreateBanner", true},
		"grant the role lacks, same resource": {editorID, "ru", "AdminDeleteBanner", false},
		"bit the role lacks":                  {editorID, "ru", "AdminCreateBlogPost", false},
		"resource the role has no row for":    {editorID, "ru", "AdminListReviews", false},
		"locale outside the staff record":     {editorID, "es", "AdminListBlogPosts", false},
		"operation the role does not hold":    {editorID, "ru", "AdminListUsers", false},
		"upload needs only being staff":       {editorID, "ru", "AdminUploadAttachment", true},
	} {
		u := db.User.GetX(t.Context(), c.user)

		err := policy.Allow(t.Context(), u, c.locale, c.operation)
		if c.allowed {
			assert.NoError(t, err, name)
		} else {
			assert.ErrorIs(t, err, permission.ErrForbidden, name)
		}
	}
}

func TestPolicyRefusesAUserWhoIsNotStaff(t *testing.T) {
	db := testsupport.NewClient(t)
	// user_blank is neither an administrator nor staff.
	u := db.User.GetX(t.Context(), fixture.ID(fixture.UserBlank))

	err := permission.NewPolicy(db).Allow(t.Context(), u, "ru", "AdminListBanners")

	assert.ErrorIs(t, err, permission.ErrForbidden)
}

// TestEveryAdminOperationIsClassified fails when a new admin operation is added
// to the contract without declaring a Permission (x-permission) or an Access
// (x-access): it would be closed to staff by default, which is safe but
// surprising.
func TestEveryAdminOperationIsClassified(t *testing.T) {
	raw, err := os.ReadFile("../../api-spec/dist/openapi.yaml")
	require.NoError(t, err)

	var spec struct {
		Paths map[string]map[string]struct {
			OperationID string                `yaml:"operationId"`
			Security    []map[string][]string `yaml:"security"`
		} `yaml:"paths"`
	}
	require.NoError(t, yaml.Unmarshal(raw, &spec))

	var unclassified []string

	for path, methods := range spec.Paths {
		if !strings.HasPrefix(path, "/api/admin/") {
			continue
		}

		for _, op := range methods {
			if op.OperationID != "" && !permission.IsClassified(permission.OperationName(op.OperationID)) {
				unclassified = append(unclassified, op.OperationID)
			}
		}
	}

	assert.Empty(t, unclassified,
		"declare x-permission or x-access on these in api-spec/admin.tsp")
}

// A staff record loaded with the user is used as it is; the answer is the same
// as when the policy loads it itself.
func TestPolicyUsesTheStaffRecordLoadedWithTheUser(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := permission.NewPolicy(db)
	editor := db.User.Query().Where(user.ID(editorID)).WithStaffMember().OnlyX(t.Context())

	require.NoError(t, policy.Allow(t.Context(), editor, "ru", "AdminListBlogPosts"))
	assert.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminCreateBlogPost"), permission.ErrForbidden)
}

// Each grant answers for its own operations: update and delete are refused
// until the role is given them, and then allowed — and only them.
func TestPolicyGrantsUpdateAndDeleteOnlyWhenTheRoleHoldsThem(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	policy := permission.NewPolicy(db)
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminUpdateBlogPost"), permission.ErrForbidden)
	require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminDeleteBlogPost"), permission.ErrForbidden)

	testsupport.Grant(t, db, transactor, staff.RoleID, "blog_posts.update")

	for _, operation := range []string{
		"AdminUpdateBlogPost", "AdminSetBlogPostRelatedCourses", "AdminSuggestBlogPostRelatedCourses",
	} {
		require.NoError(t, policy.Allow(t.Context(), editor, "ru", operation), operation)
	}

	require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminDeleteBlogPost"), permission.ErrForbidden)

	testsupport.Grant(t, db, transactor, staff.RoleID, "blog_posts.delete")
	assert.NoError(t, policy.Allow(t.Context(), editor, "ru", "AdminDeleteBlogPost"))
}

// A database that fails is an error, never a quiet refusal and never a quiet
// grant: the caller must be able to tell "not allowed" from "could not check".
func TestPolicyReportsAStorageFailureAsAnError(t *testing.T) {
	for name, table := range map[string]string{
		"staff record":  staffmember.Table,
		"role's grants": staffrolegrant.Table,
	} {
		t.Run(name, func(t *testing.T) {
			db, tx := testsupport.NewClientWithTx(t)
			editor := db.User.GetX(t.Context(), editorID)
			_, err := tx.ExecContext(t.Context(),
				`ALTER TABLE "`+table+`" RENAME TO "`+table+`_gone"`)
			require.NoError(t, err)

			// The grants are read for an operation that declares a Permission.
			operation := "AdminListBanners"
			err = permission.NewPolicy(db).Allow(t.Context(), editor, "ru", operation)

			require.Error(t, err)
			assert.NotErrorIs(t, err, permission.ErrForbidden)
		})
	}
}

// The operations that were administrators-only (users, staff members, roles,
// question-and-answer items, course versions) declare Permissions now: a staff
// member is refused until their role holds the Permission, and then gets exactly
// that operation. Nothing is open to a role that never held anything here, so the
// existing roles lose and gain nothing.
func TestFormerAdministratorsOnlyOperationsFollowTheirPermissions(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	policy := permission.NewPolicy(db)
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	former := map[string][]string{
		"users.read":                       {"AdminListUsers", "AdminGetManagementUser"},
		"users.update":                     {"AdminUpdateUser"},
		"staff_members.delete":             {"AdminDeleteStaffMember"},
		"staff_roles.read":                 {"AdminListRoles", "AdminListPermissions", "AdminGetRolePermissions"},
		"staff_roles.update":               {"AdminUpdateRolePermissions"},
		"course_category_qna_items.create": {"AdminCreateCategoryQnaItem"},
		"course_versions.create":           {"AdminCreateCourseVersion"},
	}
	// The seeded editor role holds none of them.
	for _, operations := range former {
		for _, operation := range operations {
			require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", operation), permission.ErrForbidden, operation)
		}
	}

	for permission, operations := range former {
		testsupport.Grant(t, db, transactor, staff.RoleID, permission)

		for _, operation := range operations {
			require.NoError(t, policy.Allow(t.Context(), editor, "ru", operation), operation)
		}
	}
	// A neighbour of a granted Permission stays closed.
	require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminDeleteUser"), permission.ErrForbidden)
	assert.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminCreateRole"), permission.ErrForbidden)
}

// An administrator keeps access to the user, staff and role operations, and to
// all the rest, with no grants at all.
func TestAdministratorsKeepTheFormerAdministratorsOnlyOperations(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := permission.NewPolicy(db)
	alice := db.User.GetX(t.Context(), adminID)

	for _, operation := range []string{
		"AdminListUsers", "AdminUpdateRolePermissions", "AdminListPermissions", "AdminCreateCourseVersion",
		"AdminCreateCategoryQnaItem", "AdminDeleteStaffMember", "AdminReviewCourse", "AdminReviewCourseLesson",
	} {
		assert.NoError(t, policy.Allow(t.Context(), alice, "ru", operation), operation)
	}
}

// The review of a course and of a lesson is a Permission of its own: holding
// update no longer opens it, holding review does not open update.
func TestReviewIsSeparateFromUpdate(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	policy := permission.NewPolicy(db)
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	testsupport.Grant(t, db, transactor, staff.RoleID, "courses.update")
	require.NoError(t, policy.Allow(t.Context(), editor, "ru", "AdminUpdateCourse"))
	require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminReviewCourse"), permission.ErrForbidden)

	db.StaffRoleGrant.Delete().Where(staffrolegrant.RoleID(staff.RoleID), staffrolegrant.Permission("courses.update")).ExecX(t.Context())
	testsupport.Grant(t, db, transactor, staff.RoleID, "courses.review")
	require.NoError(t, policy.Allow(t.Context(), editor, "ru", "AdminReviewCourse"))
	require.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminUpdateCourse"), permission.ErrForbidden)

	testsupport.Grant(t, db, transactor, staff.RoleID, "course_lessons.review")
	assert.NoError(t, policy.Allow(t.Context(), editor, "ru", "AdminReviewCourseLesson"))
}

// The administrator flag passes the operation Permissions but not the capability
// ones: an administrator whose role lacks mcp.read is refused the MCP check, and
// granting the capability to the role opens it. A user who is not staff has no
// role to hold it, whatever the flag says.
func TestAdministratorFlagDoesNotImplyCapabilities(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	policy := permission.NewPolicy(db)

	// Not staff: refused, though an administrator.
	alice := db.User.GetX(t.Context(), adminID)
	for _, capability := range []string{"mcp.read", "mcp.write", "mcp.manage", "audit.read"} {
		require.ErrorIs(t, policy.HasCapability(t.Context(), alice, capability), permission.ErrForbidden, capability)
	}

	// An administrator who is staff, on a role that holds no capability.
	staffAdmin := db.User.GetX(t.Context(), adminStaffID)
	require.ErrorIs(t, policy.HasCapability(t.Context(), staffAdmin, "mcp.read"), permission.ErrForbidden)
	// The administrator still passes ordinary operation Permissions.
	assert.NoError(t, policy.Allow(t.Context(), staffAdmin, "es", "AdminListBlogPosts"))
	assert.NoError(t, policy.Allow(t.Context(), staffAdmin, "ru", "AdminListUsers"))

	testsupport.Grant(t, db, transactor, staffAdmin.QueryStaffMember().OnlyX(t.Context()).RoleID, "mcp.read")
	require.NoError(t, policy.HasCapability(t.Context(), staffAdmin, "mcp.read"))
	// One capability does not open another.
	assert.ErrorIs(t, policy.HasCapability(t.Context(), staffAdmin, "mcp.write"), permission.ErrForbidden)
}

// A staff member who is not an administrator is given a capability the same way.
func TestAStaffRoleHoldsCapabilitiesThroughGrants(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	policy := permission.NewPolicy(db)
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	require.ErrorIs(t, policy.HasCapability(t.Context(), editor, "audit.read"), permission.ErrForbidden)
	testsupport.Grant(t, db, transactor, staff.RoleID, "audit.read")
	assert.NoError(t, policy.HasCapability(t.Context(), editor, "audit.read"))
}

// Asking about a Permission that is not a capability is the caller's mistake, not
// a refusal, so it cannot be mistaken for "not allowed".
func TestHasCapabilityRejectsANameThatIsNotACapability(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := db.User.GetX(t.Context(), editorID)

	err := permission.NewPolicy(db).HasCapability(t.Context(), editor, "blog_posts.read")

	require.Error(t, err)
	assert.NotErrorIs(t, err, permission.ErrForbidden)
}

// The upload is declared open to any staff member, whatever the role holds.
func TestUploadIsDeclaredOpenToAnyStaffMember(t *testing.T) {
	assert.Equal(t, permissioncatalog.AnyStaff, permissioncatalog.AccessOf("AdminUploadAttachment"))
}

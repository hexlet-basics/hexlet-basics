package staffroles_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"hexletbasics/ent/staffmember"
	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/ent/user"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/testsupport"
)

// The fixtures seed user_staff_editor as an editor holding the blog_posts.read grant and
// the banners matrix bits for reading and creating, in ru and en; user_admin_learner is an
// administrator.
var (
	adminID  = fixtures.ID(fixtures.UserAdminLearner)
	editorID = fixtures.ID(fixtures.UserStaffEditor)
)

func TestPolicyFollowsTheRolesGrantsAndMatrix(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := staffroles.NewPolicy(db)

	for name, c := range map[string]struct {
		user      int
		locale    string
		operation string
		allowed   bool
	}{
		"administrator may do anything":        {adminID, "es", "AdminDeleteUser", true},
		"granted bit":                          {editorID, "ru", "AdminListBlogPosts", true},
		"another granted bit on a resource":    {editorID, "en", "AdminCreateBanner", true},
		"bit the role lacks":                   {editorID, "ru", "AdminCreateBlogPost", false},
		"resource the role has no row for":     {editorID, "ru", "AdminListReviews", false},
		"locale outside the staff record":      {editorID, "es", "AdminListBlogPosts", false},
		"operation only administrators may do": {editorID, "ru", "AdminListUsers", false},
		"upload needs only being staff":        {editorID, "ru", "AdminUploadAttachment", true},
	} {
		u := db.User.GetX(t.Context(), c.user)
		err := policy.Allow(t.Context(), u, c.locale, c.operation)
		if c.allowed {
			assert.NoError(t, err, name)
		} else {
			assert.ErrorIs(t, err, staffroles.ErrForbidden, name)
		}
	}
}

func TestPolicyRefusesAUserWhoIsNotStaff(t *testing.T) {
	db := testsupport.NewClient(t)
	// user_blank is neither an administrator nor staff.
	u := db.User.GetX(t.Context(), fixtures.ID(fixtures.UserBlank))

	err := staffroles.NewPolicy(db).Allow(t.Context(), u, "ru", "AdminListBanners")

	assert.ErrorIs(t, err, staffroles.ErrForbidden)
}

// TestEveryAdminOperationIsClassified fails when a new admin operation is added
// to the contract without saying who may use it: it would be closed to staff
// by default, which is safe but surprising.
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
			if op.OperationID != "" && !staffroles.IsClassified(staffroles.OperationName(op.OperationID)) {
				unclassified = append(unclassified, op.OperationID)
			}
		}
	}
	assert.Empty(t, unclassified, "classify these in staffroles/rules.go")
}

// A staff record loaded with the user is used as it is; the answer is the same
// as when the policy loads it itself.
func TestPolicyUsesTheStaffRecordLoadedWithTheUser(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := staffroles.NewPolicy(db)
	editor := db.User.Query().Where(user.ID(editorID)).WithStaffMember().OnlyX(t.Context())

	assert.NoError(t, policy.Allow(t.Context(), editor, "ru", "AdminListBlogPosts"))
	assert.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminCreateBlogPost"), staffroles.ErrForbidden)
}

// Each grant answers for its own operations: update and delete are refused
// until the role is given them, and then allowed — and only them.
func TestPolicyGrantsUpdateAndDeleteOnlyWhenTheRoleHoldsThem(t *testing.T) {
	db := testsupport.NewClient(t)
	policy := staffroles.NewPolicy(db)
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	assert.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminUpdateBlogPost"), staffroles.ErrForbidden)
	assert.ErrorIs(t, policy.Allow(t.Context(), editor, "ru", "AdminDeleteBlogPost"), staffroles.ErrForbidden)

	db.StaffRoleGrant.Create().SetRoleID(staff.RoleID).SetPermission("blog_posts.update").ExecX(t.Context())
	for _, operation := range []string{
		"AdminUpdateBlogPost", "AdminSetBlogPostRelatedCourses", "AdminSuggestBlogPostRelatedCourses",
	} {
		assert.NoError(t, policy.Allow(t.Context(), bob, "ru", operation), operation)
	}
	assert.ErrorIs(t, policy.Allow(t.Context(), bob, "ru", "AdminDeleteBlogPost"), staffroles.ErrForbidden)

	db.StaffRoleGrant.Create().SetRoleID(staff.RoleID).SetPermission("blog_posts.delete").ExecX(t.Context())
	assert.NoError(t, policy.Allow(t.Context(), bob, "ru", "AdminDeleteBlogPost"))
}

// The matrix no longer decides blog posts: a stale blog_posts row granting
// everything changes nothing, because the declared Permission is the rule.
func TestPolicyIgnoresTheMatrixForBlogPosts(t *testing.T) {
	db := testsupport.NewClient(t)
	bob := db.User.GetX(t.Context(), editorID)
	staff := bob.QueryStaffMember().OnlyX(t.Context())
	db.StaffRolePermission.Create().SetRoleID(staff.RoleID).SetResource("blog_posts").
		SetCanIndex(true).SetCanCreate(true).SetCanUpdate(true).SetCanDestroy(true).ExecX(t.Context())

	err := staffroles.NewPolicy(db).Allow(t.Context(), bob, "ru", "AdminDeleteBlogPost")

	assert.ErrorIs(t, err, staffroles.ErrForbidden)
}

// The cutover migration turns each bit of a role's blog_posts matrix row into
// the grant that opens the same operations, so nobody gains or loses access.
// It runs here over every combination of the four bits.
func TestMigrationCarriesEveryBlogPostRightIntoGrants(t *testing.T) {
	migration, err := os.ReadFile("../../migrations/20261016120000_staff_role_grants.sql")
	require.NoError(t, err)
	data := string(migration[strings.Index(string(migration), "INSERT INTO"):])

	// What each bit of the old matrix opened, from the retired rules table.
	opened := map[string][]string{
		"index":   {"AdminListBlogPosts", "AdminGetBlogPost"},
		"create":  {"AdminCreateBlogPost"},
		"update":  {"AdminUpdateBlogPost", "AdminSetBlogPostRelatedCourses", "AdminSuggestBlogPostRelatedCourses"},
		"destroy": {"AdminDeleteBlogPost"},
	}
	bits := []string{"index", "create", "update", "destroy"}

	db, tx := testsupport.NewClientWithTx(t)
	bob := db.User.GetX(t.Context(), editorID)
	staff := bob.QueryStaffMember().OnlyX(t.Context())
	// Start from the old shape: no grants yet.
	db.StaffRoleGrant.Delete().ExecX(t.Context())

	for mask := 0; mask < 1<<len(bits); mask++ {
		role := db.StaffRole.Create().SetName(fmt.Sprintf("combo-%d", mask)).SaveX(t.Context())
		db.StaffRolePermission.Create().SetRoleID(role.ID).SetResource("blog_posts").
			SetCanIndex(mask&1 != 0).SetCanCreate(mask&2 != 0).
			SetCanUpdate(mask&4 != 0).SetCanDestroy(mask&8 != 0).ExecX(t.Context())
	}
	_, err = tx.ExecContext(t.Context(), data)
	require.NoError(t, err)

	remaining := db.StaffRolePermission.Query().Where(staffrolepermission.Resource("blog_posts")).CountX(t.Context())
	assert.Zero(t, remaining, "the matrix rows are replaced, not duplicated")

	policy := staffroles.NewPolicy(db)
	for mask := 0; mask < 1<<len(bits); mask++ {
		roleID := db.StaffRole.Query().Where(staffrole.Name(fmt.Sprintf("combo-%d", mask))).OnlyX(t.Context()).ID
		db.StaffMember.UpdateOneID(staff.ID).SetRoleID(roleID).ExecX(t.Context())
		for i, bit := range bits {
			for _, operation := range opened[bit] {
				err := policy.Allow(t.Context(), bob, "ru", operation)
				if mask&(1<<i) != 0 {
					assert.NoError(t, err, "bits %04b keep %s", mask, operation)
				} else {
					assert.ErrorIs(t, err, staffroles.ErrForbidden, "bits %04b still refuse %s", mask, operation)
				}
			}
		}
	}
}

// A database that fails is an error, never a quiet refusal and never a quiet
// grant: the caller must be able to tell "not allowed" from "could not check".
func TestPolicyReportsAStorageFailureAsAnError(t *testing.T) {
	for name, table := range map[string]string{
		"staff record":  staffmember.Table,
		"role's grants": staffrolegrant.Table,
		"role's matrix": staffrolepermission.Table,
	} {
		t.Run(name, func(t *testing.T) {
			db, tx := testsupport.NewClientWithTx(t)
			bob := db.User.GetX(t.Context(), editorID)
			_, err := tx.ExecContext(t.Context(),
				`ALTER TABLE "`+table+`" RENAME TO "`+table+`_gone"`)
			require.NoError(t, err)

			// Each table answers for its own operation: blog posts declare a
			// Permission, banners still use the matrix.
			operation := "AdminListBanners"
			if table == staffrolegrant.Table {
				operation = "AdminListBlogPosts"
			}
			err = staffroles.NewPolicy(db).Allow(t.Context(), bob, "ru", operation)

			require.Error(t, err)
			assert.NotErrorIs(t, err, staffroles.ErrForbidden)
		})
	}
}

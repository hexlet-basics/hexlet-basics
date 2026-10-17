package staffroles_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/staffrole"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/testsupport"
)

// opened is what each bit of a resource's retired matrix row opened: the
// operations the old policy let the bit perform (can_index: list and get,
// can_create: create, can_update: update and review, can_destroy: delete).
// Bits that opened nothing are absent, and a role holding them lost nothing by
// the migration not carrying them.
type opened struct {
	index, create, update, destroy []string
}

// cutover is one migration that moved some resources from the matrix to grants:
// the legacy `resource` value of each matrix row and what its bits opened.
type cutover struct {
	migration string
	resources map[string]opened
}

// legacyMatrixTable is the baseline's `staff_member_role_permissions`, which a
// later migration drops.
const legacyMatrixTable = `CREATE TABLE IF NOT EXISTS "staff_member_role_permissions" (
  "id" bigserial NOT NULL,
  "can_create" boolean NOT NULL DEFAULT false,
  "can_destroy" boolean NOT NULL DEFAULT false,
  "can_index" boolean NOT NULL DEFAULT false,
  "can_update" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL,
  "resource" character varying NOT NULL,
  "role_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
)`

var cutovers = []cutover{
	{
		migration: "20261016120000_staff_role_grants.sql",
		resources: map[string]opened{
			"blog_posts": {
				index:   []string{"AdminListBlogPosts", "AdminGetBlogPost"},
				create:  []string{"AdminCreateBlogPost"},
				update:  []string{"AdminUpdateBlogPost", "AdminSetBlogPostRelatedCourses", "AdminSuggestBlogPostRelatedCourses"},
				destroy: []string{"AdminDeleteBlogPost"},
			},
		},
	},
	{
		migration: "20261017100000_banners_reviews_leads_messages_into_grants.sql",
		resources: map[string]opened{
			"banners": {
				index:   []string{"AdminListBanners", "AdminGetBanner"},
				create:  []string{"AdminCreateBanner"},
				update:  []string{"AdminUpdateBanner"},
				destroy: []string{"AdminDeleteBanner"},
			},
			"reviews": {
				index:   []string{"AdminListReviews", "AdminGetReview"},
				create:  []string{"AdminCreateReview"},
				update:  []string{"AdminUpdateReview"},
				destroy: []string{"AdminDeleteReview"},
			},
			"leads":    {index: []string{"AdminListLeads"}},
			"messages": {index: []string{"AdminListMessages"}},
		},
	},
}

// TestMigrationsCarryEveryMatrixRightIntoGrants runs each cutover migration over
// a role for every combination of the four bits of every resource it moves, and
// asks the policy about each operation the bits used to open: nobody gains or
// loses an operation.
func TestMigrationsCarryEveryMatrixRightIntoGrants(t *testing.T) {
	for _, c := range cutovers {
		for resource, ops := range c.resources {
			t.Run(c.migration+"/"+resource, func(t *testing.T) {
				assertMigrationKeepsAccess(t, c.migration, resource, ops)
			})
		}
	}
}

func assertMigrationKeepsAccess(t *testing.T, file, resource string, ops opened) {
	t.Helper()
	migration, err := os.ReadFile("../../migrations/" + file)
	require.NoError(t, err)
	text := string(migration)
	data := text[strings.Index(text, "INSERT INTO"):]

	byBit := [][]string{ops.index, ops.create, ops.update, ops.destroy}

	db, tx := testsupport.NewClientWithTx(t)
	bob := db.User.GetX(t.Context(), bobID)
	staff := bob.QueryStaffMember().OnlyX(t.Context())
	// Start from the old shape: no grants, and the matrix table as the baseline
	// made it. The table is long gone from the migrated schema, so the test
	// recreates it inside its transaction (a no-op while it still exists).
	db.StaffRoleGrant.Delete().ExecX(t.Context())
	_, err = tx.ExecContext(t.Context(), legacyMatrixTable)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `DELETE FROM "staff_member_role_permissions"`)
	require.NoError(t, err)

	for mask := 0; mask < 1<<len(byBit); mask++ {
		role := db.StaffRole.Create().SetName(fmt.Sprintf("combo-%d", mask)).SaveX(t.Context())
		_, err = tx.ExecContext(t.Context(), `INSERT INTO "staff_member_role_permissions"
			("role_id", "resource", "can_index", "can_create", "can_update", "can_destroy", "created_at", "updated_at")
			VALUES ($1, $2, $3, $4, $5, $6, now(), now())`,
			role.ID, resource, mask&1 != 0, mask&2 != 0, mask&4 != 0, mask&8 != 0)
		require.NoError(t, err)
	}
	_, err = tx.ExecContext(t.Context(), data)
	require.NoError(t, err)

	var remaining int
	require.NoError(t, tx.QueryRowContext(t.Context(),
		`SELECT count(*) FROM "staff_member_role_permissions" WHERE "resource" = $1`, resource).Scan(&remaining))
	assert.Zero(t, remaining, "the matrix rows are replaced, not duplicated")

	policy := staffroles.NewPolicy(db)
	for mask := 0; mask < 1<<len(byBit); mask++ {
		roleID := db.StaffRole.Query().Where(staffrole.Name(fmt.Sprintf("combo-%d", mask))).OnlyX(t.Context()).ID
		db.StaffMember.UpdateOneID(staff.ID).SetRoleID(roleID).ExecX(t.Context())
		for i, operations := range byBit {
			for _, operation := range operations {
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

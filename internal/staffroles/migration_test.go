package staffroles_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/api"
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
		migration: "20261017200000_staff_role_grants.sql",
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
		migration: "20261017201000_banners_reviews_leads_messages_into_grants.sql",
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
	{
		migration: "20261017202000_course_resources_into_grants.sql",
		resources: map[string]opened{
			"language_categories": {
				index:   []string{"AdminListCourseCategories", "AdminGetCourseCategory"},
				create:  []string{"AdminCreateCourseCategory"},
				update:  []string{"AdminUpdateCourseCategory"},
				destroy: []string{"AdminDeleteCourseCategory"},
			},
			"languages": {
				index:  []string{"AdminListCourses", "AdminGetCourse"},
				create: []string{"AdminCreateCourse"},
				update: []string{"AdminUpdateCourse"},
			},
			"language_lessons":        {index: []string{"AdminListCourseLessons"}},
			"language_lesson_reviews": {index: []string{"AdminListCourseLessonReviews"}},
			"language_lesson_members": {index: []string{"AdminListLessonProgress"}},
			"language_landing_pages": {
				index:   []string{"AdminListCourseLandingPages", "AdminGetCourseLandingPage"},
				create:  []string{"AdminCreateCourseLandingPage"},
				update:  []string{"AdminUpdateCourseLandingPage"},
				destroy: []string{"AdminDeleteCourseLandingPage"},
			},
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

// createRole makes a Staff Role the way an administrator does, through the
// operation, and returns its id.
func createRole(t *testing.T, h *testsupport.Harness, name string) int {
	t.Helper()
	role, err := h.Client.AdminCreateRole(t.Context(), &api.RoleInput{Name: name, Description: api.NilString{Null: true}})
	require.NoError(t, err)
	return int(role.ID)
}

func assertMigrationKeepsAccess(t *testing.T, file, resource string, ops opened) {
	t.Helper()
	migration, err := os.ReadFile("../../migrations/" + file)
	require.NoError(t, err)
	text := string(migration)
	data := text[strings.Index(text, "INSERT INTO"):]

	byBit := [][]string{ops.index, ops.create, ops.update, ops.destroy}

	// The roles are made the way an administrator makes them, through the
	// operation; the harness's transaction is the one the migration runs in.
	h := testsupport.NewHarness(t)
	db, tx := h.DB, h.Tx
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())
	// Start from the old shape: no grants, and the matrix table as the baseline
	// made it. The table is long gone from the migrated schema, so the test
	// recreates it inside its transaction (a no-op while it still exists).
	db.StaffRoleGrant.Delete().ExecX(t.Context())
	_, err = tx.ExecContext(t.Context(), legacyMatrixTable)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `DELETE FROM "staff_member_role_permissions"`)
	require.NoError(t, err)

	for mask := 0; mask < 1<<len(byBit); mask++ {
		role := createRole(t, h, fmt.Sprintf("combo-%d", mask))
		_, err = tx.ExecContext(t.Context(), `INSERT INTO "staff_member_role_permissions"
			("role_id", "resource", "can_index", "can_create", "can_update", "can_destroy", "created_at", "updated_at")
			VALUES ($1, $2, $3, $4, $5, $6, now(), now())`,
			role, resource, mask&1 != 0, mask&2 != 0, mask&4 != 0, mask&8 != 0)
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
				err := policy.Allow(t.Context(), editor, "ru", operation)
				if mask&(1<<i) != 0 {
					assert.NoError(t, err, "bits %04b keep %s", mask, operation)
				} else {
					assert.ErrorIs(t, err, staffroles.ErrForbidden, "bits %04b still refuse %s", mask, operation)
				}
			}
		}
	}
}

// TestReviewMigrationKeepsEveryRolesAccess runs the review migration over a role
// for every combination of holding courses.update and course_lessons.update, and
// asks the policy about update and review: nobody gains or loses an operation,
// because review used to be opened by update.
func TestReviewMigrationKeepsEveryRolesAccess(t *testing.T) {
	migration, err := os.ReadFile("../../migrations/20261018100000_review_permissions.sql")
	require.NoError(t, err)

	h := testsupport.NewHarness(t)
	db, tx := h.DB, h.Tx
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())
	db.StaffRoleGrant.Delete().ExecX(t.Context())

	const holdsCourses, holdsLessons = 1, 2
	for mask := 0; mask < 4; mask++ {
		// A bystander grant must survive untouched.
		grants := []string{"courses.read"}
		if mask&holdsCourses != 0 {
			grants = append(grants, "courses.update")
		}
		if mask&holdsLessons != 0 {
			grants = append(grants, "course_lessons.update")
		}
		h.Grant(t, createRole(t, h, fmt.Sprintf("combo-%d", mask)), grants...)
	}
	_, err = tx.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)

	policy := staffroles.NewPolicy(db)
	for mask := 0; mask < 4; mask++ {
		roleID := db.StaffRole.Query().Where(staffrole.Name(fmt.Sprintf("combo-%d", mask))).OnlyX(t.Context()).ID
		db.StaffMember.UpdateOneID(staff.ID).SetRoleID(roleID).ExecX(t.Context())

		check := func(operation string, want bool) {
			err := policy.Allow(t.Context(), editor, "ru", operation)
			if want {
				assert.NoError(t, err, "combination %02b keeps %s", mask, operation)
			} else {
				assert.ErrorIs(t, err, staffroles.ErrForbidden, "combination %02b still refuses %s", mask, operation)
			}
		}
		check("AdminListCourses", true)
		check("AdminUpdateCourse", mask&holdsCourses != 0)
		check("AdminReviewCourse", mask&holdsCourses != 0)
		check("AdminReviewCourseLesson", mask&holdsLessons != 0)
	}

	dead, err := db.StaffRoleGrant.Query().Where(staffrolegrant.Permission("course_lessons.update")).Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, dead, "the grant that names no operation is replaced")
}

// TestBlogPostStateMigrationGrantsTheTransitionsToEveryUpdater runs the
// migration over a role for every combination of holding blog_posts.update, an
// already-held blog_posts.publish and a bystander blog_posts.read, and asks the
// policy about each operation: every holder of update gains the three
// transitions (the state used to be settable through the update), and no other
// role gains one.
func TestBlogPostStateMigrationGrantsTheTransitionsToEveryUpdater(t *testing.T) {
	migration, err := os.ReadFile("../../migrations/20261020090000_blog_post_state_permissions.sql")
	require.NoError(t, err)

	// The roles are made the way an administrator makes them, through the role
	// operations, so the scenario holds exactly the grants the API can give.
	h := testsupport.NewHarness(t)
	db := h.DB
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	const holdsUpdate, holdsPublish, holdsRead = 1, 2, 4
	grants := map[int]string{holdsUpdate: "blog_posts.update", holdsPublish: "blog_posts.publish", holdsRead: "blog_posts.read"}
	roles := make([]int, 8)
	for mask := range roles {
		role, err := h.Client.AdminCreateRole(t.Context(), &api.RoleInput{
			Name:        fmt.Sprintf("combo-%d", mask),
			Description: api.NilString{Null: true},
		})
		require.NoError(t, err)
		roles[mask] = int(role.ID)

		held := []string{}
		for bit, permission := range grants {
			if mask&bit != 0 {
				held = append(held, permission)
			}
		}
		_, err = h.Client.AdminUpdateRolePermissions(t.Context(), &api.RolePermissionsInput{Grants: held},
			api.AdminUpdateRolePermissionsParams{RoleId: role.ID})
		require.NoError(t, err)
	}
	// Running it twice proves it is idempotent for a role that already holds a grant.
	_, err = h.Tx.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)
	_, err = h.Tx.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)

	policy := staffroles.NewPolicy(db)
	for mask, roleID := range roles {
		db.StaffMember.UpdateOneID(staff.ID).SetRoleID(roleID).ExecX(t.Context())

		check := func(operation string, want bool) {
			err := policy.Allow(t.Context(), editor, "ru", operation)
			if want {
				assert.NoError(t, err, "combination %03b keeps %s", mask, operation)
			} else {
				assert.ErrorIs(t, err, staffroles.ErrForbidden, "combination %03b still refuses %s", mask, operation)
			}
		}
		updater := mask&holdsUpdate != 0
		check("AdminListBlogPosts", mask&holdsRead != 0)
		check("AdminUpdateBlogPost", updater)
		check("AdminPublishBlogPost", updater || mask&holdsPublish != 0)
		check("AdminUnpublishBlogPost", updater)
		check("AdminArchiveBlogPost", updater)
	}
}

// TestStateMigrationsGrantTheTransitionsToEveryUpdater is the blog post state
// migration's test for banners, reviews and landing pages: each migration runs
// over a role for every combination of holding the resource's update, an
// already-held publish and a bystander read, and the policy is asked about each
// operation. Every holder of update gains the three transitions (the state used
// to be settable through the update), and no other role gains one.
func TestStateMigrationsGrantTheTransitionsToEveryUpdater(t *testing.T) {
	for _, resource := range []struct {
		migration, permissions, operation, list string
	}{
		{"20261020100000_banner_state_permissions.sql", "banners", "Banner", "AdminListBanners"},
		{"20261020110000_review_state_permissions.sql", "reviews", "Review", "AdminListReviews"},
		{"20261020120000_course_landing_page_state_permissions.sql", "course_landing_pages", "CourseLandingPage", "AdminListCourseLandingPages"},
	} {
		t.Run(resource.permissions, func(t *testing.T) {
			migration, err := os.ReadFile("../../migrations/" + resource.migration)
			require.NoError(t, err)

			// The roles are made the way an administrator makes them, through the role
			// operations, so the scenario holds exactly the grants the API can give.
			h := testsupport.NewHarness(t)
			db := h.DB
			editor := db.User.GetX(t.Context(), editorID)
			staff := editor.QueryStaffMember().OnlyX(t.Context())

			const holdsUpdate, holdsPublish, holdsRead = 1, 2, 4
			grants := map[int]string{
				holdsUpdate:  resource.permissions + ".update",
				holdsPublish: resource.permissions + ".publish",
				holdsRead:    resource.permissions + ".read",
			}
			roles := make([]int, 8)
			for mask := range roles {
				role, err := h.Client.AdminCreateRole(t.Context(), &api.RoleInput{
					Name:        fmt.Sprintf("combo-%d", mask),
					Description: api.NilString{Null: true},
				})
				require.NoError(t, err)
				roles[mask] = int(role.ID)

				held := []string{}
				for bit, permission := range grants {
					if mask&bit != 0 {
						held = append(held, permission)
					}
				}
				_, err = h.Client.AdminUpdateRolePermissions(t.Context(), &api.RolePermissionsInput{Grants: held},
					api.AdminUpdateRolePermissionsParams{RoleId: role.ID})
				require.NoError(t, err)
			}
			// Running it twice proves it is idempotent for a role that already holds a grant.
			_, err = h.Tx.ExecContext(t.Context(), string(migration))
			require.NoError(t, err)
			_, err = h.Tx.ExecContext(t.Context(), string(migration))
			require.NoError(t, err)

			policy := staffroles.NewPolicy(db)
			for mask, roleID := range roles {
				db.StaffMember.UpdateOneID(staff.ID).SetRoleID(roleID).ExecX(t.Context())

				check := func(operation string, want bool) {
					err := policy.Allow(t.Context(), editor, "ru", operation)
					if want {
						assert.NoError(t, err, "combination %03b keeps %s", mask, operation)
					} else {
						assert.ErrorIs(t, err, staffroles.ErrForbidden, "combination %03b still refuses %s", mask, operation)
					}
				}
				updater := mask&holdsUpdate != 0
				check(resource.list, mask&holdsRead != 0)
				check("AdminUpdate"+resource.operation, updater)
				check("AdminPublish"+resource.operation, updater || mask&holdsPublish != 0)
				check("AdminUnpublish"+resource.operation, updater)
				check("AdminArchive"+resource.operation, updater)
			}
		})
	}
}

// TestCourseReadinessMigrationGrantsTheChangesToEveryUpdater is the state
// migrations' test for a Course's Readiness: the migration runs over a role for
// every combination of holding courses.update, an already-held mark_completed
// and a bystander read, and the policy is asked about each operation. Every
// holder of update gains the three changes (the readiness used to be settable
// through the update), and no other role gains one.
func TestCourseReadinessMigrationGrantsTheChangesToEveryUpdater(t *testing.T) {
	migration, err := os.ReadFile("../../migrations/20261020130000_course_readiness_permissions.sql")
	require.NoError(t, err)

	// The roles are made the way an administrator makes them, through the role
	// operations, so the scenario holds exactly the grants the API can give.
	h := testsupport.NewHarness(t)
	db := h.DB
	editor := db.User.GetX(t.Context(), editorID)
	staff := editor.QueryStaffMember().OnlyX(t.Context())

	const holdsUpdate, holdsCompleted, holdsRead = 1, 2, 4
	grants := map[int]string{
		holdsUpdate:    "courses.update",
		holdsCompleted: "courses.mark_completed",
		holdsRead:      "courses.read",
	}
	roles := make([]int, 8)
	for mask := range roles {
		role, err := h.Client.AdminCreateRole(t.Context(), &api.RoleInput{
			Name:        fmt.Sprintf("combo-%d", mask),
			Description: api.NilString{Null: true},
		})
		require.NoError(t, err)
		roles[mask] = int(role.ID)

		held := []string{}
		for bit, permission := range grants {
			if mask&bit != 0 {
				held = append(held, permission)
			}
		}
		_, err = h.Client.AdminUpdateRolePermissions(t.Context(), &api.RolePermissionsInput{Grants: held},
			api.AdminUpdateRolePermissionsParams{RoleId: role.ID})
		require.NoError(t, err)
	}
	// Running it twice proves it is idempotent for a role that already holds a grant.
	_, err = h.Tx.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)
	_, err = h.Tx.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)

	policy := staffroles.NewPolicy(db)
	for mask, roleID := range roles {
		db.StaffMember.UpdateOneID(staff.ID).SetRoleID(roleID).ExecX(t.Context())

		check := func(operation string, want bool) {
			err := policy.Allow(t.Context(), editor, "ru", operation)
			if want {
				assert.NoError(t, err, "combination %03b keeps %s", mask, operation)
			} else {
				assert.ErrorIs(t, err, staffroles.ErrForbidden, "combination %03b still refuses %s", mask, operation)
			}
		}
		updater := mask&holdsUpdate != 0
		check("AdminListCourses", mask&holdsRead != 0)
		check("AdminUpdateCourse", updater)
		check("AdminMarkCourseDraft", updater)
		check("AdminMarkCourseInDevelopment", updater)
		check("AdminMarkCourseCompleted", updater || mask&holdsCompleted != 0)
	}
}

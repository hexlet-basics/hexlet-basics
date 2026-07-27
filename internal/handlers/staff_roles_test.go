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

// staff_member_roles.yml seeds two roles; editor owns two permission rows.
const (
	totalRoles   = 2
	roleEditorID = 9001
)

func TestAdminListRoles(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListRoles(ctx, api.AdminListRolesParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(totalRoles), page.Total)

	// Ordered by name: editor before moderator.
	require.Len(t, page.Items, totalRoles)
	assert.Equal(t, "editor", page.Items[0].Name)
	assert.Equal(t, "moderator", page.Items[1].Name)
	// editor owns two permission rows.
	assert.Equal(t, int32(2), page.Items[0].PermissionsCount)
	assert.Equal(t, int32(0), page.Items[1].PermissionsCount)
}

func TestAdminGetRole(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.AdminGetRole(ctx, api.AdminGetRoleParams{ID: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "editor", got.Name)
	assert.Len(t, got.Permissions, 2)
}

func TestAdminCreateRole(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	created, err := h.Client.AdminCreateRole(ctx, &api.RoleInput{
		Name:        "support",
		Description: api.NewNilString("Support team"),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.Equal(t, "support", created.Name)
	assert.Empty(t, created.Permissions)
}

func TestAdminCreateRoleDuplicateName(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateRole(ctx, &api.RoleInput{
		Name:        "editor", // already seeded
		Description: api.NilString{Null: true},
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestAdminUpdateRole(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	updated, err := h.Client.AdminUpdateRole(ctx, &api.RoleInput{
		Name:        "editor-renamed",
		Description: api.NilString{Null: true}, // clears
	}, api.AdminUpdateRoleParams{ID: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "editor-renamed", updated.Name)
	assert.True(t, updated.Description.Null)
	// Permissions are untouched by a role update.
	assert.Len(t, updated.Permissions, 2)
}

func TestAdminDeleteRole(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// moderator (9002) has no permissions -> no FK children to block delete.
	err := h.Client.AdminDeleteRole(ctx, api.AdminDeleteRoleParams{ID: 9002})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.StaffRole.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalRoles-1, total)
}

func TestAdminGetRolePermissions(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.AdminGetRolePermissions(ctx, api.AdminGetRolePermissionsParams{RoleId: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "editor", got.Name)
	assert.Len(t, got.Permissions, 2)
}

// TestAdminUpdateRolePermissions upserts the matrix: updates an existing row
// (banners) and inserts a new one (reviews), leaving the untouched row
// (blog_posts) in place — matching the legacy find_or_initialize sync.
func TestAdminUpdateRolePermissions(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	updated, err := h.Client.AdminUpdateRolePermissions(ctx, &api.RolePermissionsInput{
		Permissions: []api.RolePermissionInput{
			{Resource: api.PermissionResourceBanners, CanIndex: true, CanCreate: true, CanUpdate: true, CanDestroy: true},
			{Resource: api.PermissionResourceReviews, CanIndex: true, CanCreate: false, CanUpdate: false, CanDestroy: false},
		},
	}, api.AdminUpdateRolePermissionsParams{RoleId: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	// blog_posts (untouched) + banners (updated) + reviews (inserted) = 3.
	assert.Len(t, updated.Permissions, 3)

	byResource := map[api.PermissionResource]api.StaffRolePermission{}
	for _, p := range updated.Permissions {
		byResource[p.Resource] = p
	}
	assert.True(t, byResource[api.PermissionResourceBanners].CanDestroy, "banners upgraded to full CRUD")
	assert.True(t, byResource[api.PermissionResourceReviews].CanIndex)
	assert.False(t, byResource[api.PermissionResourceReviews].CanCreate)
	_, hasBlog := byResource[api.PermissionResourceBlogPosts]
	assert.True(t, hasBlog, "unlisted blog_posts permission is left untouched")
}

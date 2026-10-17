package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// staff_member_roles.yml seeds two roles; editor owns three grants (banners.create, banners.read, blog_posts.read).
var (
	roleEditorID = fixtures.ID32(fixtures.StaffRoleEditor)
)

// editorGrants are the editor's grants in the fixtures, in the order the API
// lists them.
var editorGrants = []string{"banners.create", "banners.read", "blog_posts.read"}

func TestAdminListRoles(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalRoles := h.DB.StaffRole.Query().CountX(ctx)

	page, err := h.Client.AdminListRoles(ctx, api.AdminListRolesParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(totalRoles), page.Total)

	// Ordered by name: editor before moderator.
	require.Len(t, page.Items, totalRoles)
	assert.Equal(t, "editor", page.Items[0].Name)
	assert.Equal(t, "moderator", page.Items[1].Name)
	// editor owns three grants.
	assert.Equal(t, int32(3), page.Items[0].PermissionsCount)
	assert.Equal(t, int32(0), page.Items[1].PermissionsCount)
}

func TestAdminGetRole(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.AdminGetRole(ctx, api.AdminGetRoleParams{ID: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "editor", got.Name)
	assert.Equal(t, editorGrants, got.Grants)
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
	assert.Empty(t, created.Grants)
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
	// Grants are untouched by a role update.
	assert.Equal(t, editorGrants, updated.Grants)
}

func TestAdminDeleteRole(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalRoles := h.DB.StaffRole.Query().CountX(ctx)

	// moderator (9002) has no permissions -> no FK children to block delete.
	err := h.Client.AdminDeleteRole(ctx, api.AdminDeleteRoleParams{ID: fixtures.ID32(fixtures.StaffRoleModerator)})
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
	assert.Equal(t, editorGrants, got.Grants)
}

// TestAdminUpdateRolePermissionsReplacesTheGrants checks that the grants
// submitted are exactly the grants the role holds afterwards: the new ones are
// added, the ones left out are removed, the ones kept stay.
func TestAdminUpdateRolePermissionsReplacesTheGrants(t *testing.T) {
	h := testsupport.NewHarness(t)

	updated, err := h.Client.AdminUpdateRolePermissions(context.Background(), &api.RolePermissionsInput{
		Grants: []string{"blog_posts.update", "blog_posts.create", "blog_posts.create"},
	}, api.AdminUpdateRolePermissionsParams{RoleId: roleEditorID})
	require.NoError(t, err)

	assert.Equal(t, []string{"blog_posts.create", "blog_posts.update"}, updated.Grants)

	again, err := h.Client.AdminUpdateRolePermissions(context.Background(), &api.RolePermissionsInput{
		Grants: []string{},
	}, api.AdminUpdateRolePermissionsParams{RoleId: roleEditorID})
	require.NoError(t, err)
	assert.Empty(t, again.Grants)
}

// TestAdminUpdateRolePermissionsRefusesAnUnknownGrant keeps a typo from being
// stored as a Permission that opens nothing, and writes nothing else of the
// save either.
func TestAdminUpdateRolePermissionsRefusesAnUnknownGrant(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminUpdateRolePermissions(context.Background(), &api.RolePermissionsInput{
		Grants: []string{"blog_posts.publish"},
	}, api.AdminUpdateRolePermissionsParams{RoleId: roleEditorID})

	require.Error(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, h.LastStatus())
	got, err := h.Client.AdminGetRole(context.Background(), api.AdminGetRoleParams{ID: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, editorGrants, got.Grants)
}

// TestAdminListPermissions serves the catalog the roles screen offers, the one
// read from the operations' declarations in the contract.
func TestAdminListPermissions(t *testing.T) {
	h := testsupport.NewHarness(t)

	got, err := h.Client.AdminListPermissions(context.Background())
	require.NoError(t, err)

	// The whole catalog is pinned by the contract test; here is that the screen's
	// endpoint serves it, with each Permission's kind.
	assert.Subset(t, got.Permissions, []api.Permission{
		{Name: "banners.create", Kind: api.PermissionKindWrite},
		{Name: "banners.read", Kind: api.PermissionKindRead},
		{Name: "blog_posts.update", Kind: api.PermissionKindWrite},
		{Name: "leads.read", Kind: api.PermissionKindRead},
		{Name: "reviews.delete", Kind: api.PermissionKindWrite},
	})
}

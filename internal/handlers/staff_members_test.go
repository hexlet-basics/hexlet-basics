package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// staff_members.yml seeds the editor (staff_member_staff_editor) and an administrator who is staff (staff_member_admin_staff).

func TestAdminListStaffMembers(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalStaffMembers := h.DB.StaffMember.Query().CountX(ctx)

	page, err := h.Client.AdminListStaffMembers(ctx, api.AdminListStaffMembersParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(totalStaffMembers), page.Total)

	m := page.Items[0]
	assert.Equal(t, int32(fixture.ID(fixture.UserStaffEditor)), m.UserId)
	assert.Equal(t, "editor", m.Role.Name)
	assert.Equal(t, "staff-editor@example.com", m.User.Email.Value)
	// The native pg array round-trips to a plain []string.
	assert.Equal(t, []string{"ru", "en"}, m.AllowedLocales)
}

func TestAdminGetStaffMember(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.AdminGetStaffMember(ctx, api.AdminGetStaffMemberParams{ID: fixture.ID32(fixture.StaffMemberStaffEditor)})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "editor", got.Role.Name)
	assert.Equal(t, "staff-editor@example.com", got.User.Email.Value)
	assert.Equal(t, []string{"ru", "en"}, got.AllowedLocales)
}

func TestAdminGetStaffMemberNotFound(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetStaffMember(ctx, api.AdminGetStaffMemberParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateStaffMember(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalStaffMembers := h.DB.StaffMember.Query().CountX(ctx)

	created, err := h.Client.AdminCreateStaffMember(ctx, &api.StaffMemberInput{
		UserId:         fixture.ID32(fixture.UserWithPassword), // has no staff record yet; the caller (user_admin_learner) may not manage their own
		RoleId:         fixture.ID32(fixture.StaffRoleEditor),
		AllowedLocales: []string{"ru", "en", "es"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.NotZero(t, created.ID)
	assert.Equal(t, "with-password@example.com", created.User.Email.Value)
	assert.Equal(t, "editor", created.Role.Name)
	assert.Equal(t, []string{"ru", "en", "es"}, created.AllowedLocales)

	total, err := h.DB.StaffMember.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalStaffMembers+1, total)
}

// TestAdminCreateStaffMemberDuplicateUser proves the unique user_id index: a
// user already having a staff record surfaces as 409.
func TestAdminCreateStaffMemberDuplicateUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateStaffMember(ctx, &api.StaffMemberInput{
		UserId:         fixture.ID32(fixture.UserStaffEditor), // already a staff member (9201)
		RoleId:         fixture.ID32(fixture.StaffRoleEditor),
		AllowedLocales: []string{"ru"},
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestAdminUpdateStaffMember(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	updated, err := h.Client.AdminUpdateStaffMember(ctx, &api.StaffMemberInput{
		UserId:         fixture.ID32(fixture.UserStaffEditor),
		RoleId:         fixture.ID32(fixture.StaffRoleModerator), // reassign to moderator
		AllowedLocales: []string{"en"},
	}, api.AdminUpdateStaffMemberParams{ID: fixture.ID32(fixture.StaffMemberStaffEditor)})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "moderator", updated.Role.Name)
	assert.Equal(t, []string{"en"}, updated.AllowedLocales)
}

func TestAdminDeleteStaffMember(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalStaffMembers := h.DB.StaffMember.Query().CountX(ctx)

	err := h.Client.AdminDeleteStaffMember(ctx, api.AdminDeleteStaffMemberParams{ID: fixture.ID32(fixture.StaffMemberStaffEditor)})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.StaffMember.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalStaffMembers-1, total)
}

package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/staffmember"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// The fixtures: alice is the first administrator (the harness user), dana the second, bob the editor
// (role staff_role_editor), carol a plain user, and the moderator role holds no grants. A third
// administrator, user_admin_staff, is also a staff member (staff_member_admin_staff).
var (
	aliceID        = fixture.ID(fixture.UserAdminLearner)
	danaID         = fixture.ID(fixture.UserSecondAdmin)
	danaLabel      = fixture.UserSecondAdmin
	adminStaffID   = fixture.ID(fixture.UserAdminStaff)
	carolID        = fixture.ID(fixture.UserWithPassword)
	editorID       = fixture.ID(fixture.UserStaffEditor)
	moderator      = fixture.ID32(fixture.StaffRoleModerator)
	editorRecordID = fixture.ID32(fixture.StaffMemberStaffEditor)
)

// delegate gives the editor role Permissions it does not hold in the fixtures,
// the way an administrator delegating staff management would.
func delegate(t *testing.T, h *testsupport.Harness, permissions ...string) {
	t.Helper()
	h.Grant(t, int(roleEditorID), permissions...)
}

// staffRecordOf makes user a staff member with the role and returns the record
// id, the way an administrator does: through the operation, signed in as the
// second administrator on the harness's transaction (a record is never made for
// oneself, and the caller may be a staff member without the right to create).
func staffRecordOf(t *testing.T, h *testsupport.Harness, userID, roleID int) int {
	t.Helper()
	made, err := h.ActAs(t, danaLabel).Client.AdminCreateStaffMember(t.Context(), staffInput(userID, roleID))
	require.NoError(t, err)
	return int(made.ID)
}

// leaveOnlyAlice makes the first administrator the only one: the cases about
// the last administrator start from there. The fixtures hold three, dana and
// the staff member who is also an administrator being the others.
func leaveOnlyAlice(t *testing.T, h *testsupport.Harness) {
	t.Helper()
	h.DB.User.Update().Where(user.IDIn(danaID, adminStaffID)).SetAdmin(false).ExecX(t.Context())
}

func staffInput(userID, roleID int) *api.StaffMemberInput {
	return &api.StaffMemberInput{UserId: int32(userID), RoleId: int32(roleID), AllowedLocales: []string{"ru"}}
}

func userInput(email string, admin bool) *api.UserInput {
	return &api.UserInput{
		Email:     email,
		FirstName: api.NilString{Null: true},
		LastName:  api.NilString{Null: true},
		Admin:     api.NewNilBool(admin),
	}
}

// problemOf is the Problem Details the caller got for a refused call.
func problemOf(t *testing.T, err error) api.ProblemDetails {
	t.Helper()
	var problem *api.ProblemDetailsStatusCode
	require.ErrorAs(t, err, &problem)
	return problem.Response
}

func TestStaffMemberCannotChangeTheirOwnStaffRecord(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "staff_members.update", "staff_members.delete", "staff_members.create")

	_, err := h.Client.AdminUpdateStaffMember(t.Context(), staffInput(editorID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: editorRecordID})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())

	_, err = h.Client.AdminUpdateStaffMember(t.Context(), staffInput(editorID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: editorRecordID})
	assert.Contains(t, problemOf(t, err).Detail.Value, "your own staff record")

	own := h.DB.StaffMember.GetX(t.Context(), int(editorRecordID))
	assert.Equal(t, int(roleEditorID), own.RoleID, "the role is unchanged")

	err = h.Client.AdminDeleteStaffMember(t.Context(), api.AdminDeleteStaffMemberParams{ID: editorRecordID})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	exists, err := h.DB.StaffMember.Query().Where(staffmember.ID(int(editorRecordID))).Exist(t.Context())
	require.NoError(t, err)
	assert.True(t, exists)

	// Nor can a new record be made for oneself.
	h.DB.StaffMember.DeleteOneID(int(editorRecordID)).ExecX(t.Context())
	_, err = h.Client.AdminCreateStaffMember(t.Context(), staffInput(editorID, int(moderator)))
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestStaffMemberCannotMoveAnotherRecordOntoThemselves(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "staff_members.update")
	carol := staffRecordOf(t, h, carolID, int(moderator))

	_, err := h.Client.AdminUpdateStaffMember(t.Context(), staffInput(editorID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: int32(carol)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
}

func TestStaffMemberWithTheGrantsMayChangeSomeoneElsesRecord(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "staff_members.update", "staff_members.delete")
	carol := staffRecordOf(t, h, carolID, int(moderator))

	updated, err := h.Client.AdminUpdateStaffMember(t.Context(), staffInput(carolID, int(roleEditorID)), api.AdminUpdateStaffMemberParams{ID: int32(carol)})
	require.NoError(t, err)
	assert.Equal(t, "editor", updated.Role.Name)

	require.NoError(t, h.Client.AdminDeleteStaffMember(t.Context(), api.AdminDeleteStaffMemberParams{ID: int32(carol)}))
}

func TestAdministratorCannotChangeTheirOwnStaffRecordButAnotherCan(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserAdminStaff)
	record := fixture.ID(fixture.StaffMemberAdminStaff)

	_, err := h.Client.AdminUpdateStaffMember(t.Context(), staffInput(adminStaffID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: int32(record)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	err = h.Client.AdminDeleteStaffMember(t.Context(), api.AdminDeleteStaffMemberParams{ID: int32(record)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())

	// Dana, another administrator, can make the same changes to that record.
	other := h.ActAs(t, danaLabel)
	_, err = other.Client.AdminUpdateStaffMember(t.Context(), staffInput(adminStaffID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: int32(record)})
	require.NoError(t, err)
	require.NoError(t, other.Client.AdminDeleteStaffMember(t.Context(), api.AdminDeleteStaffMemberParams{ID: int32(record)}))
}

func TestARefusedSelfChangeIsRecordedAsFailed(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "staff_members.update")

	_, err := h.Client.AdminUpdateStaffMember(t.Context(), staffInput(editorID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: editorRecordID})
	require.Error(t, err)

	rows := records(t, h, "AdminUpdateStaffMember")
	require.Len(t, rows, 1)
	assert.Equal(t, "failed", rows[0].Outcome)
	assert.Equal(t, editorID, *rows[0].ActorUserID)
	assert.Empty(t, rows[0].Changes)
}

func TestTheRefusalIsExplainedInTheCallersLanguage(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "staff_members.update")
	testsupport.SpeakTo(h, "ru")

	_, err := h.Client.AdminUpdateStaffMember(t.Context(), staffInput(editorID, int(moderator)), api.AdminUpdateStaffMemberParams{ID: editorRecordID})
	assert.Contains(t, problemOf(t, err).Detail.Value, "собственную запись")
}

func TestUserCannotDeleteTheirOwnAccountThroughTheAdminOperations(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	err := h.Client.AdminDeleteUser(t.Context(), api.AdminDeleteUserParams{ID: int32(aliceID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	exists, err := h.DB.User.Query().Where(user.ID(aliceID)).Exist(t.Context())
	require.NoError(t, err)
	assert.True(t, exists)

	// Another administrator can delete alice, and alice can delete dana.
	require.NoError(t, h.Client.AdminDeleteUser(t.Context(), api.AdminDeleteUserParams{ID: int32(danaID)}))
}

func TestAnotherAdministratorMayDeleteAnAdministrator(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, danaLabel)
	// alice owns content rows that keep her from being deleted, so an administrator
	// of the case's own, made the way an administrator makes one, stands in for
	// the one dana removes.
	erin, err := h.Client.AdminCreateUser(t.Context(), userInput("erin@example.com", true))
	require.NoError(t, err)

	require.NoError(t, h.Client.AdminDeleteUser(t.Context(), api.AdminDeleteUserParams{ID: erin.ID}))
}

func TestTheLastAdministratorFlagCannotBeClearedWithOneAdministrator(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	leaveOnlyAlice(t, h)

	_, err := h.Client.AdminUpdateUser(t.Context(), userInput("alice@example.com", false), api.AdminUpdateUserParams{ID: int32(aliceID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
	assert.Contains(t, problemOf(t, err).Detail.Value, "last administrator")

	_, err = h.Client.AdminUpdateManagementUser(t.Context(), userInput("alice@example.com", false), api.AdminUpdateManagementUserParams{ID: int32(aliceID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())

	assert.True(t, *h.DB.User.GetX(t.Context(), aliceID).Admin)
}

func TestAnAdministratorFlagCanBeClearedWithTwoAdministrators(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.DB.User.UpdateOneID(adminStaffID).SetAdmin(false).ExecX(t.Context())

	_, err := h.Client.AdminUpdateUser(t.Context(), userInput("dana@example.com", false), api.AdminUpdateUserParams{ID: int32(danaID)})
	require.NoError(t, err)
	assert.False(t, *h.DB.User.GetX(t.Context(), danaID).Admin)
}

func TestTheLastAdministratorCannotBeDeletedByAStaffMember(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "users.delete")
	leaveOnlyAlice(t, h)

	err := h.Client.AdminDeleteUser(t.Context(), api.AdminDeleteUserParams{ID: int32(aliceID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
}

func TestStaffMemberCannotChangeTheGrantsOfTheirOwnRole(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "staff_roles.update")
	grantsBefore := h.DB.StaffRoleGrant.Query().CountX(t.Context())

	_, err := h.Client.AdminUpdateRolePermissions(t.Context(),
		&api.RolePermissionsInput{Grants: []string{"staff_roles.update", "users.update"}},
		api.AdminUpdateRolePermissionsParams{RoleId: roleEditorID})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	assert.Contains(t, problemOf(t, err).Detail.Value, "your own staff role")
	held := h.DB.StaffRoleGrant.Query().CountX(t.Context())
	assert.Equal(t, grantsBefore, held, "the refused saves added nothing")

	// The grants of another role are theirs to change.
	_, err = h.Client.AdminUpdateRolePermissions(t.Context(),
		&api.RolePermissionsInput{Grants: []string{"blog_posts.read"}},
		api.AdminUpdateRolePermissionsParams{RoleId: moderator})
	require.NoError(t, err)
}

func TestAdministratorMayChangeTheGrantsOfAnyRole(t *testing.T) {
	t.Parallel()
	// An administrator who is also staff, on the very role whose grants change.
	h := testsupport.NewHarnessAs(t, fixture.UserAdminStaff)

	updated, err := h.Client.AdminUpdateRolePermissions(t.Context(),
		&api.RolePermissionsInput{Grants: []string{"blog_posts.read"}},
		api.AdminUpdateRolePermissionsParams{RoleId: roleEditorID})
	require.NoError(t, err)
	assert.Equal(t, []string{"blog_posts.read"}, updated.Grants)
}

func TestOnlyAnAdministratorChangesTheAdministratorFlag(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, editorLabel)
	delegate(t, h, "users.update", "users.create")

	// Raising oneself, through either user surface.
	_, err := h.Client.AdminUpdateUser(t.Context(), userInput("bob@example.com", true), api.AdminUpdateUserParams{ID: int32(editorID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	_, err = h.Client.AdminUpdateManagementUser(t.Context(), userInput("bob@example.com", true), api.AdminUpdateManagementUserParams{ID: int32(editorID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	assert.False(t, *h.DB.User.GetX(t.Context(), editorID).Admin)

	// Raising someone else, creating an administrator, demoting one.
	_, err = h.Client.AdminUpdateUser(t.Context(), userInput("carol@example.com", true), api.AdminUpdateUserParams{ID: int32(carolID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	_, err = h.Client.AdminCreateUser(t.Context(), userInput("new@example.com", true))
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	_, err = h.Client.AdminUpdateUser(t.Context(), userInput("dana@example.com", false), api.AdminUpdateUserParams{ID: int32(danaID)})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())

	// Everything else about a user stays within users.update and users.create.
	_, err = h.Client.AdminUpdateUser(t.Context(), userInput("carol2@example.com", false), api.AdminUpdateUserParams{ID: int32(carolID)})
	require.NoError(t, err)
	_, err = h.Client.AdminCreateUser(t.Context(), userInput("new@example.com", false))
	require.NoError(t, err)
}

func TestAdministratorMayCreateAnotherAdministrator(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	created, err := h.Client.AdminCreateUser(t.Context(), userInput("boss@example.com", true))
	require.NoError(t, err)
	assert.True(t, created.Admin.Value)
}

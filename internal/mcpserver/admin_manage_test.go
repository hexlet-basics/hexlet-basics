package mcpserver_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/ent/user"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/permissioncatalog"
	"hexletbasics/internal/testsupport"
)

var ()

func userBody(email string, admin any) map[string]any {
	return map[string]any{"email": email, "firstName": nil, "lastName": nil, "admin": admin}
}

// accessOperations is one read and one write of each class of operation that
// decides who has access: users, staff members, Staff Roles and their grants.
// Each carries its own Permission and the MCP capability of its kind.
func accessOperations() []struct {
	class, tool string
	permission  string
	capability  string
	args        map[string]any
} {
	return []struct {
		class, tool string
		permission  string
		capability  string
		args        map[string]any
	}{
		{"users read", "adminListUsers", "users.read", "mcp.read", map[string]any{"locale": "en"}},
		{"users write", "adminUpdateUser", "users.update", "mcp.write", map[string]any{
			"locale": "en", "id": carol(), "body": userBody("with-password@example.com", false)}},
		{"staff members read", "adminListStaffMembers", "staff_members.read", "mcp.read", map[string]any{"locale": "en"}},
		{"staff members write", "adminUpdateStaffMember", "staff_members.update", "mcp.write", map[string]any{
			"locale": "en", "id": adminStaffSeat(), "body": map[string]any{
				"userId": adminStaff(), "roleId": editor(), "allowedLocales": []string{"en", "ru"}}}},
		{"roles read", "adminListRoles", "staff_roles.read", "mcp.read", map[string]any{"locale": "en"}},
		{"roles write", "adminCreateRole", "staff_roles.create", "mcp.write", map[string]any{
			"locale": "en", "body": map[string]any{"name": "Auditors of the assistant", "description": nil}}},
		{"grants read", "adminGetRolePermissions", "staff_roles.read", "mcp.read", map[string]any{"locale": "en", "roleId": editor()}},
		{"grants write", "adminUpdateRolePermissions", "staff_roles.update", "mcp.write", map[string]any{
			"locale": "en", "roleId": moderator(), "body": map[string]any{"grants": []string{"blog_posts.read"}}}},
		{"permission catalog", "adminListPermissions", "staff_roles.read", "mcp.read", map[string]any{"locale": "en"}},
	}
}

// The contract declares mcp.manage for every operation on users, staff members,
// Staff Roles and their grants; no other operation carries it.
func TestTheContractDeclaresMCPManageForTheAccessOperations(t *testing.T) {
	declared := permissioncatalog.RequiredCapabilityOperations()
	assert.Len(t, declared, 22)

	for _, name := range declared {
		assert.Contains(t, permissioncatalog.MCPCapabilities(name), "mcp.manage", name)
		requirement, ok := permissioncatalog.For(name)
		require.True(t, ok, name)
		assert.Regexp(t, `^(users|staff_members|staff_roles)\.`, requirement.Permission, name)
	}

	assert.Equal(t, []string{"mcp.read"}, permissioncatalog.MCPCapabilities("AdminListBanners"))
	assert.Equal(t, []string{"mcp.write", "mcp.manage"}, permissioncatalog.MCPCapabilities("AdminDeleteUser"))
	assert.Equal(t, []string{"mcp.read", "mcp.manage"}, permissioncatalog.MCPCapabilities("AdminGetRolePermissions"))
}

func TestAccessOperationsNeedMCPManageOnTopOfEverythingElse(t *testing.T) {
	rates := mcpserver.Rates{Catalogue: 1000, Personal: 1000, Write: 1000}

	for _, c := range accessOperations() {
		t.Run(c.class, func(t *testing.T) {
			o := newOAuthWithRates(t, rates)
			// Everything the operation needs except mcp.manage.
			o.Grant(editor(), c.permission, c.capability)
			session, token := admin(t, o, bob())

			refused := call(t, session, c.tool, c.args)

			assertRefused(t, refused)

			records := o.Records(ogenName(c.tool))
			require.Len(t, records, 1)
			assert.Equal(t, "forbidden", records[0].Outcome)
			assert.Equal(t, "mcp", records[0].Channel)
			assert.Equal(t, token.ID, *records[0].TokenID)

			o.Grant(editor(), "mcp.manage")

			allowed := call(t, session, c.tool, c.args)
			assert.False(t, allowed.IsError, "%+v", allowed.Content)
		})
	}
}

func TestAnAdministratorWithoutAnExplicitMCPManageGrantIsRefused(t *testing.T) {
	rates := mcpserver.Rates{Catalogue: 1000, Personal: 1000, Write: 1000}

	for _, c := range accessOperations() {
		t.Run(c.class, func(t *testing.T) {
			o := newOAuthWithRates(t, rates)
			// The administrator flag gives every operation Permission; the role
			// holds the MCP capabilities of the kind but not mcp.manage.
			o.Grant(editor(), "mcp.read", "mcp.write")
			session, _ := admin(t, o, adminStaff())

			args := c.args
			if c.tool == "adminUpdateStaffMember" {
				// Not the caller's own record: that is refused by self-protection.
				args = map[string]any{"locale": "en", "id": bobSeat(), "body": map[string]any{
					"userId": bob(), "roleId": editor(), "allowedLocales": []string{"en", "ru"}}}
			}

			assertRefused(t, call(t, session, c.tool, args))

			o.Grant(editor(), "mcp.manage")

			allowed := call(t, session, c.tool, args)
			assert.False(t, allowed.IsError, "%+v", allowed.Content)
		})
	}
}

func TestAnAccessWriteThroughAToolChangesTheRows(t *testing.T) {
	o := newOAuthWithRates(t, mcpserver.Rates{Catalogue: 100, Personal: 100, Write: 100})
	o.Grant(editor(), "staff_roles.update", "staff_roles.read", "mcp.write", "mcp.read", "mcp.manage")
	session, _ := admin(t, o, bob())

	result := call(t, session, "adminUpdateRolePermissions", map[string]any{
		"locale": "en", "roleId": moderator(), "body": map[string]any{"grants": []string{"blog_posts.read", "banners.read"}},
	})

	require.False(t, result.IsError, "%+v", result.Content)
	held := o.DB.StaffRoleGrant.Query().Where(staffrolegrant.RoleID(moderator())).CountX(t.Context())
	assert.Equal(t, 2, held)
}

func TestToolsAnnounceWhatTheyDo(t *testing.T) {
	byName := map[string]*mcp.Tool{}
	for _, tool := range mcpserver.AdminToolsFromTheContract(t) {
		byName[tool.Name] = tool
	}

	for _, name := range []string{"adminDeleteBanner", "adminDeleteUser", "adminDeleteStaffMember", "adminDeleteRole", "adminDeleteReview"} {
		require.Contains(t, byName, name)
		require.NotNil(t, byName[name].Annotations.DestructiveHint, name)
		assert.True(t, *byName[name].Annotations.DestructiveHint, name)
		assert.False(t, byName[name].Annotations.ReadOnlyHint, name)
	}

	for name, tool := range byName {
		switch {
		case permissioncatalog.IsAdminWrite(ogenName(name)):
			assert.False(t, tool.Annotations.ReadOnlyHint, name)

			if tool.Annotations.DestructiveHint != nil && *tool.Annotations.DestructiveHint {
				assert.Contains(t, name, "Delete", "only a delete is destructive: archive and unpublish keep the record")
			}
		default:
			assert.True(t, tool.Annotations.ReadOnlyHint, name)
		}

		assert.NotEmpty(t, tool.Description, name)
	}

	assert.False(t, *byName["adminArchiveBanner"].Annotations.DestructiveHint)
	assert.True(t, byName["adminUpdateBanner"].Annotations.IdempotentHint)
}

// newOAuthWithRates is NewOAuth with budgets large enough that a test of what
// the checkpoint refuses is not stopped by the rate limit first.
func newOAuthWithRates(t *testing.T, rates mcpserver.Rates) *testsupport.OAuth {
	t.Helper()

	return testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{User: rates, Client: rates}))
}

func ogenName(operationID string) string { return string(operationID[0]-32) + operationID[1:] }

// The self-protection rules (#844) live in the handlers, so a tool call is
// refused exactly as the HTTP call is: a 403 or 409 from the handler, recorded
// as failed, and no row changed.
func TestSelfProtectionHoldsOverMCP(t *testing.T) {
	rates := mcpserver.Rates{Catalogue: 1000, Personal: 1000, Write: 1000}
	staffManagement := []string{
		"users.read", "users.update", "users.delete", "users.create",
		"staff_members.update", "staff_members.delete", "staff_members.create",
		"staff_roles.update", "staff_roles.read", "mcp.read", "mcp.write", "mcp.manage",
	}

	assertRefusedWith := func(t *testing.T, result *mcp.CallToolResult, status string) {
		t.Helper()
		require.True(t, result.IsError, "%+v", result.Content)
		assert.Contains(t, textOf(t, result), status)
	}

	t.Run("own Staff Member record", func(t *testing.T) {
		o := newOAuthWithRates(t, rates)
		o.Grant(editor(), staffManagement...)
		session, _ := admin(t, o, bob())

		update := call(t, session, "adminUpdateStaffMember", map[string]any{
			"locale": "en", "id": bobSeat(), "body": map[string]any{
				"userId": bob(), "roleId": moderator(), "allowedLocales": []string{"en"}},
		})
		assertRefusedWith(t, update, "HTTP 403")
		remove := call(t, session, "adminDeleteStaffMember", map[string]any{"locale": "en", "id": bobSeat()})
		assertRefusedWith(t, remove, "HTTP 403")

		member := o.DB.StaffMember.GetX(t.Context(), bobSeat())
		assert.Equal(t, editor(), member.RoleID, "the role is unchanged")

		records := o.Records("AdminUpdateStaffMember")
		require.Len(t, records, 1)
		assert.Equal(t, "failed", records[0].Outcome)
		assert.Equal(t, "mcp", records[0].Channel)

		// Someone else's record is theirs to change.
		other := call(t, session, "adminUpdateStaffMember", map[string]any{
			"locale": "en", "id": adminStaffSeat(), "body": map[string]any{
				"userId": adminStaff(), "roleId": moderator(), "allowedLocales": []string{"en"}},
		})
		assert.False(t, other.IsError, "%+v", other.Content)
	})

	t.Run("own account", func(t *testing.T) {
		o := newOAuthWithRates(t, rates)
		o.Grant(editor(), "mcp.write", "mcp.manage")
		session, _ := admin(t, o, adminStaff())

		result := call(t, session, "adminDeleteUser", map[string]any{"locale": "en", "id": adminStaff()})

		assertRefusedWith(t, result, "HTTP 403")
		assert.True(t, o.DB.User.Query().Where(user.ID(adminStaff())).ExistX(t.Context()))
	})

	t.Run("last administrator", func(t *testing.T) {
		o := newOAuthWithRates(t, rates)
		o.Grant(editor(), staffManagement...)
		// adminStaff is left the only administrator.
		o.DB.User.Update().Where(user.IDNEQ(adminStaff())).SetAdmin(false).ExecX(t.Context())
		session, _ := admin(t, o, adminStaff())

		demoted := call(t, session, "adminUpdateUser", map[string]any{
			"locale": "en", "id": adminStaff(), "body": userBody("admin-staff@example.com", false)})
		assertRefusedWith(t, demoted, "HTTP 409")
		assert.True(t, o.DB.User.GetX(t.Context(), adminStaff()).Admin != nil && *o.DB.User.GetX(t.Context(), adminStaff()).Admin)

		// A staff member who may delete users cannot delete the last administrator.
		bobSession, _ := admin(t, o, bob())
		remove := call(t, bobSession, "adminDeleteUser", map[string]any{"locale": "en", "id": adminStaff()})
		assertRefusedWith(t, remove, "HTTP 409")
		assert.True(t, o.DB.User.Query().Where(user.ID(adminStaff())).ExistX(t.Context()))

		// With a second administrator the same flag change goes through.
		o.DB.User.UpdateOneID(secondAdmin()).SetAdmin(true).ExecX(t.Context())
		again := call(t, session, "adminUpdateUser", map[string]any{
			"locale": "en", "id": adminStaff(), "body": userBody("admin-staff@example.com", false)})
		assert.False(t, again.IsError, "%+v", again.Content)
	})

	t.Run("only an administrator changes the administrator flag", func(t *testing.T) {
		o := newOAuthWithRates(t, rates)
		o.Grant(editor(), staffManagement...)
		session, _ := admin(t, o, bob())

		promote := call(t, session, "adminUpdateUser", map[string]any{
			"locale": "en", "id": carol(), "body": userBody("with-password@example.com", true)})
		assertRefusedWith(t, promote, "HTTP 403")
		create := call(t, session, "adminCreateUser", map[string]any{
			"locale": "en", "body": userBody("new-admin@example.com", true)})
		assertRefusedWith(t, create, "HTTP 403")
		assert.False(t, *o.DB.User.GetX(t.Context(), carol()).Admin)

		adminSession, _ := admin(t, o, adminStaff())
		allowed := call(t, adminSession, "adminUpdateUser", map[string]any{
			"locale": "en", "id": carol(), "body": userBody("with-password@example.com", true)})
		assert.False(t, allowed.IsError, "%+v", allowed.Content)
	})

	t.Run("grants of the caller's own role", func(t *testing.T) {
		o := newOAuthWithRates(t, rates)
		o.Grant(editor(), staffManagement...)
		session, _ := admin(t, o, bob())
		before := o.DB.StaffRoleGrant.Query().CountX(t.Context())

		own := call(t, session, "adminUpdateRolePermissions", map[string]any{
			"locale": "en", "roleId": editor(), "body": map[string]any{"grants": []string{"staff_roles.update", "users.update"}}})
		assertRefusedWith(t, own, "HTTP 403")
		assert.Equal(t, before, o.DB.StaffRoleGrant.Query().CountX(t.Context()))

		other := call(t, session, "adminUpdateRolePermissions", map[string]any{
			"locale": "en", "roleId": moderator(), "body": map[string]any{"grants": []string{"blog_posts.read"}}})
		assert.False(t, other.IsError, "%+v", other.Content)
	})
}

func secondAdmin() int {
	return fixture.ID(fixture.UserSecondAdmin)
}

func moderator() int {
	return fixture.ID(fixture.StaffRoleModerator)
}

func adminStaffSeat() int {
	return fixture.ID(fixture.StaffMemberAdminStaff)
}

func bobSeat() int {
	return fixture.ID(fixture.StaffMemberStaffEditor)
}

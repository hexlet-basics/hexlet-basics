package handlers

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"hexletbasics/internal/api"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/permissioncatalog"
)

// roleDetail fetches a role with its grants and converts it to the
// StaffRoleDetail projection returned by create/update/get and the permission
// endpoints. A missing id surfaces as ent's not-found error (404 centrally).
func (s *Server) roleDetail(ctx context.Context, id int) (*api.StaffRoleDetail, error) {
	return s.db.StaffRole.AdminGet(ctx, int32(id), s.conv.ToStaffRoleDetail)
}

// --- Roles ------------------------------------------------------------------

// The role CRUD is generated (ent/admincrud.go): the list is ordered by name
// (legacy `order(:name)`) with each row's permission count, the single view is
// the StaffRoleDetail matrix. Permissions start empty on create and are set via
// the role-permissions endpoint; a duplicate name is a 409 centrally.

func (s *Server) AdminListRoles(ctx context.Context, params api.AdminListRolesParams) (*api.StaffRolePage, error) {
	return s.db.StaffRole.AdminList(ctx, params.Page, params.PerPage, s.conv.ToStaffRoles)
}

func (s *Server) AdminGetRole(ctx context.Context, params api.AdminGetRoleParams) (*api.StaffRoleDetail, error) {
	return s.db.StaffRole.AdminGet(ctx, params.ID, s.conv.ToStaffRoleDetail)
}

func (s *Server) AdminCreateRole(ctx context.Context, req *api.RoleInput) (*api.StaffRoleDetail, error) {
	return s.db.StaffRole.AdminCreate(ctx, req, s.conv.ToStaffRoleDetail)
}

func (s *Server) AdminUpdateRole(ctx context.Context, req *api.RoleInput, params api.AdminUpdateRoleParams) (*api.StaffRoleDetail, error) {
	return s.db.StaffRole.AdminUpdate(ctx, params.ID, req, s.conv.ToStaffRoleDetail)
}

func (s *Server) AdminDeleteRole(ctx context.Context, params api.AdminDeleteRoleParams) error {
	if err := s.db.StaffRole.AdminDelete(ctx, params.ID); err != nil {
		return err
	}
	return nil
}

// --- Role permissions -------------------------------------------------------

func (s *Server) AdminGetRolePermissions(ctx context.Context, params api.AdminGetRolePermissionsParams) (*api.StaffRoleDetail, error) {
	return s.roleDetail(ctx, int(params.RoleId))
}

// AdminUpdateRolePermissions applies the submitted set; the permission
// module writes them in one transaction. A missing role is a 404,
// before any write, and an undeclared grant a 422.
func (s *Server) AdminUpdateRolePermissions(ctx context.Context, req *api.RolePermissionsInput, params api.AdminUpdateRolePermissionsParams) (*api.StaffRoleDetail, error) {
	roleID := int(params.RoleId)

	if err := s.selfProtection.OwnRoleGrants(ctx, actorOfRequest(ctx), roleID); err != nil {
		return nil, err
	}

	if err := s.permissions.Check(ctx, roleID, req.Grants); err != nil {
		if errors.Is(err, permission.ErrUnknownPermission) {
			return nil, invalidText("grants", err.Error())
		}
		return nil, err
	}
	if err := s.permissions.Apply(ctx, roleID, req.Grants); err != nil {
		return nil, err
	}
	return s.roleDetail(ctx, roleID)
}

// --- Permission catalog -----------------------------------------------------

// AdminListPermissions returns the catalog of Permissions the roles screen
// offers. It is read from the operations' declarations in the contract, so a
// Permission appears here exactly when some operation requires it.
func (s *Server) AdminListPermissions(context.Context) (*api.PermissionCatalog, error) {
	catalog := permissioncatalog.Catalog()
	return &api.PermissionCatalog{Permissions: lo.Map(catalog, func(p permissioncatalog.Permission, _ int) api.Permission {
		return api.Permission{Name: p.Name, Kind: api.PermissionKind(p.Kind), Capability: p.Capability}
	})}, nil
}

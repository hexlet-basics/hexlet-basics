package handlers

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"hexletbasics/internal/api"
	"hexletbasics/internal/permissions"
	"hexletbasics/internal/staffroles"
)

// roleDetail fetches a role with its permission matrix and converts it to the
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

func (s *Server) AdminListRoles(ctx context.Context, params api.AdminListRolesParams) (api.AdminListRolesRes, error) {
	return s.db.StaffRole.AdminList(ctx, params.Page, params.PerPage, s.conv.ToStaffRoles)
}

func (s *Server) AdminGetRole(ctx context.Context, params api.AdminGetRoleParams) (api.AdminGetRoleRes, error) {
	return s.db.StaffRole.AdminGet(ctx, params.ID, s.conv.ToStaffRoleDetail)
}

func (s *Server) AdminCreateRole(ctx context.Context, req *api.RoleInput) (api.AdminCreateRoleRes, error) {
	return s.db.StaffRole.AdminCreate(ctx, req, s.conv.ToStaffRoleDetail)
}

func (s *Server) AdminUpdateRole(ctx context.Context, req *api.RoleInput, params api.AdminUpdateRoleParams) (api.AdminUpdateRoleRes, error) {
	return s.db.StaffRole.AdminUpdate(ctx, params.ID, req, s.conv.ToStaffRoleDetail)
}

func (s *Server) AdminDeleteRole(ctx context.Context, params api.AdminDeleteRoleParams) (api.AdminDeleteRoleRes, error) {
	if err := s.db.StaffRole.AdminDelete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.AdminDeleteRoleNoContent{}, nil
}

// --- Role permissions -------------------------------------------------------

func (s *Server) AdminGetRolePermissions(ctx context.Context, params api.AdminGetRolePermissionsParams) (api.AdminGetRolePermissionsRes, error) {
	return s.roleDetail(ctx, int(params.RoleId))
}

// AdminUpdateRolePermissions applies the submitted matrix and grants; the
// staffroles module writes them in one transaction. A missing role is a 404,
// before any write, and an undeclared grant a 422.
func (s *Server) AdminUpdateRolePermissions(ctx context.Context, req *api.RolePermissionsInput, params api.AdminUpdateRolePermissionsParams) (api.AdminUpdateRolePermissionsRes, error) {
	roleID := int(params.RoleId)

	if err := s.permissions.Apply(ctx, roleID, req.Permissions, req.Grants); err != nil {
		if errors.Is(err, staffroles.ErrUnknownPermission) {
			return validationError("grants", err.Error()), nil
		}
		return nil, err
	}
	return s.roleDetail(ctx, roleID)
}

// --- Permission catalog -----------------------------------------------------

// AdminListPermissions returns the catalog of Permissions the roles screen
// offers. It is read from the operations' declarations in the contract, so a
// Permission appears here exactly when some operation requires it.
func (s *Server) AdminListPermissions(context.Context) (api.AdminListPermissionsRes, error) {
	catalog := permissions.Catalog()
	return &api.PermissionCatalog{Permissions: lo.Map(catalog, func(p permissions.Permission, _ int) api.Permission {
		return api.Permission{Name: p.Name, Kind: api.PermissionKind(p.Kind)}
	})}, nil
}

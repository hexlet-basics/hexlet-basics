package handlers

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrole"
	"hexletbasics/internal/api"
	"hexletbasics/internal/permissions"
)

// roleDetail fetches a role with its grants and converts it to the
// StaffRoleDetail projection returned by create/update/get and the permission
// endpoints. A missing id surfaces as ent's not-found error (404 centrally).
func (s *Server) roleDetail(ctx context.Context, id int) (*api.StaffRoleDetail, error) {
	return getOne(ctx, id,
		func(ctx context.Context, id int) (*ent.StaffRole, error) {
			return s.db.StaffRole.Query().Where(staffrole.ID(id)).WithGrants().Only(ctx)
		},
		s.conv.ToStaffRoleDetail,
	)
}

// --- Roles ------------------------------------------------------------------

// AdminListRoles returns roles ordered by name (legacy `order(:name)`), each
// with its Permission count.
func (s *Server) AdminListRoles(ctx context.Context, params api.AdminListRolesParams) (api.AdminListRolesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.StaffRoleQuery {
			return s.db.StaffRole.Query().WithGrants().Order(ent.Asc(staffrole.FieldName))
		},
		s.conv.ToStaffRoles,
		func(items []api.StaffRole, total, page, perPage int32) *api.StaffRolePage {
			return &api.StaffRolePage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

func (s *Server) AdminGetRole(ctx context.Context, params api.AdminGetRoleParams) (api.AdminGetRoleRes, error) {
	return s.roleDetail(ctx, int(params.ID))
}

// AdminCreateRole creates a role (name + description); permissions start empty
// and are set via the role-permissions endpoint. A duplicate name violates the
// unique index and surfaces as 409 centrally.
func (s *Server) AdminCreateRole(ctx context.Context, req *api.RoleInput) (api.AdminCreateRoleRes, error) {
	row, err := s.db.StaffRole.Create().SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.roleDetail(ctx, row.ID)
}

func (s *Server) AdminUpdateRole(ctx context.Context, req *api.RoleInput, params api.AdminUpdateRoleParams) (api.AdminUpdateRoleRes, error) {
	row, err := s.db.StaffRole.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.roleDetail(ctx, row.ID)
}

func (s *Server) AdminDeleteRole(ctx context.Context, params api.AdminDeleteRoleParams) (api.AdminDeleteRoleRes, error) {
	if err := s.db.StaffRole.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteRoleNoContent{}, nil
}

// --- Role permissions -------------------------------------------------------

func (s *Server) AdminGetRolePermissions(ctx context.Context, params api.AdminGetRolePermissionsParams) (api.AdminGetRolePermissionsRes, error) {
	return s.roleDetail(ctx, int(params.RoleId))
}

// AdminUpdateRolePermissions replaces the role's grants with the submitted set;
// the permissions module writes them in one transaction.
func (s *Server) AdminUpdateRolePermissions(ctx context.Context, req *api.RolePermissionsInput, params api.AdminUpdateRolePermissionsParams) (api.AdminUpdateRolePermissionsRes, error) {
	roleID := int(params.RoleId)

	// Ensure the role exists first (404 for a missing id, before any write).
	if _, err := s.db.StaffRole.Query().Where(staffrole.ID(roleID)).Only(ctx); err != nil {
		return nil, err
	}

	// A grant must name a Permission the contract declares: a typo would
	// otherwise be stored and silently open nothing.
	if unknown := lo.Reject(req.Grants, func(name string, _ int) bool { return permissions.Known(name) }); len(unknown) > 0 {
		return validationError("grants", "unknown permission: "+strings.Join(unknown, ", ")), nil
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
func (s *Server) AdminListPermissions(context.Context) (api.AdminListPermissionsRes, error) {
	catalog := permissions.Catalog()
	return &api.PermissionCatalog{Permissions: lo.Map(catalog, func(p permissions.Permission, _ int) api.Permission {
		return api.Permission{Name: p.Name, Kind: api.PermissionKind(p.Kind)}
	})}, nil
}

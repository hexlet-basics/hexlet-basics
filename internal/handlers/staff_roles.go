package handlers

import (
	"context"
	"slices"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
)

// roleDetail fetches a role with its permission matrix and converts it to the
// StaffRoleDetail projection returned by create/update/get and the permission
// endpoints. A missing id surfaces as ent's not-found error (404 centrally).
func (s *Server) roleDetail(ctx context.Context, id int) (*api.StaffRoleDetail, error) {
	return getOne(ctx, id,
		func(ctx context.Context, id int) (*ent.StaffRole, error) {
			return s.db.StaffRole.Query().Where(staffrole.ID(id)).WithPermissions().Only(ctx)
		},
		s.conv.ToStaffRoleDetail,
	)
}

// --- Roles ------------------------------------------------------------------

// AdminListRoles returns roles ordered by name (legacy `order(:name)`), each
// with its permission count.
func (s *Server) AdminListRoles(ctx context.Context, params api.AdminListRolesParams) (api.AdminListRolesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.StaffRoleQuery {
			return s.db.StaffRole.Query().WithPermissions().Order(ent.Asc(staffrole.FieldName))
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

// AdminUpdateRolePermissions applies the submitted matrix, upserting each row by
// (role_id, resource) — mirrors the legacy StaffRoleMutator#sync_permissions
// (find_or_initialize_by(resource).update!), so unlisted resources are left
// untouched. The unique index on (role_id, resource) decides insert or update
// in one statement, so the matrix lands whole or not at all without a
// transaction, and two admins saving at once cannot both insert a row.
func (s *Server) AdminUpdateRolePermissions(ctx context.Context, req *api.RolePermissionsInput, params api.AdminUpdateRolePermissionsParams) (api.AdminUpdateRolePermissionsRes, error) {
	roleID := int(params.RoleId)

	// Ensure the role exists first (404 for a missing id, before any write).
	if _, err := s.db.StaffRole.Query().Where(staffrole.ID(roleID)).Only(ctx); err != nil {
		return nil, err
	}

	// A resource listed twice would make the upsert touch one row twice, which
	// PostgreSQL rejects; the last entry wins, as it did when rows were written
	// one after another.
	perms := slices.Clone(req.Permissions)
	slices.Reverse(perms)
	perms = lo.UniqBy(perms, func(p api.RolePermissionInput) api.PermissionResource { return p.Resource })
	slices.Reverse(perms)

	builders := lo.Map(perms, func(p api.RolePermissionInput, _ int) *ent.StaffRolePermissionCreate {
		return s.db.StaffRolePermission.Create().SetRoleID(roleID).SetInput(&p)
	})
	err := s.db.StaffRolePermission.CreateBulk(builders...).
		OnConflictColumns(staffrolepermission.FieldRoleID, staffrolepermission.FieldResource).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return s.roleDetail(ctx, roleID)
}

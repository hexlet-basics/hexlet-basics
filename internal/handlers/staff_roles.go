package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
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
func (s *Server) AdminListRoles(ctx context.Context, params api.AdminListRolesParams) (*api.StaffRolePage, error) {
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

func (s *Server) AdminGetRole(ctx context.Context, params api.AdminGetRoleParams) (*api.StaffRoleDetail, error) {
	return s.roleDetail(ctx, int(params.ID))
}

// AdminCreateRole creates a role (name + description); permissions start empty
// and are set via the role-permissions endpoint. A duplicate name violates the
// unique index and surfaces as 409 centrally.
func (s *Server) AdminCreateRole(ctx context.Context, req *api.RoleInput) (*api.StaffRoleDetail, error) {
	row, err := s.db.StaffRole.Create().
		SetName(req.Name).
		SetNillableDescription(inputconv.Ptr(req.Description)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.roleDetail(ctx, row.ID)
}

func (s *Server) AdminUpdateRole(ctx context.Context, req *api.RoleInput, params api.AdminUpdateRoleParams) (*api.StaffRoleDetail, error) {
	upd := s.db.StaffRole.UpdateOneID(int(params.ID)).SetName(req.Name)
	applyNil(req.Description.Null, req.Description.Value, upd.SetDescription, upd.ClearDescription)

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.roleDetail(ctx, row.ID)
}

func (s *Server) AdminDeleteRole(ctx context.Context, params api.AdminDeleteRoleParams) error {
	return s.db.StaffRole.DeleteOneID(int(params.ID)).Exec(ctx)
}

// --- Role permissions -------------------------------------------------------

func (s *Server) AdminGetRolePermissions(ctx context.Context, params api.AdminGetRolePermissionsParams) (*api.StaffRoleDetail, error) {
	return s.roleDetail(ctx, int(params.RoleId))
}

// AdminUpdateRolePermissions applies the submitted matrix, upserting each row by
// (role_id, resource) — mirrors the legacy StaffRoleMutator#sync_permissions
// (find_or_initialize_by(resource).update!), so unlisted resources are left
// untouched. ent's upsert feature is not enabled, so this is a per-row
// find-or-create.
func (s *Server) AdminUpdateRolePermissions(ctx context.Context, req *api.RolePermissionsInput, params api.AdminUpdateRolePermissionsParams) (*api.StaffRoleDetail, error) {
	roleID := int(params.RoleId)

	// Ensure the role exists first (404 for a missing id, before any write).
	if _, err := s.db.StaffRole.Query().Where(staffrole.ID(roleID)).Only(ctx); err != nil {
		return nil, err
	}

	for _, p := range req.Permissions {
		resource := string(p.Resource)
		existing, err := s.db.StaffRolePermission.Query().
			Where(
				staffrolepermission.RoleID(roleID),
				staffrolepermission.Resource(resource),
			).
			Only(ctx)
		switch {
		case ent.IsNotFound(err):
			_, err = s.db.StaffRolePermission.Create().
				SetRoleID(roleID).
				SetResource(resource).
				SetCanIndex(p.CanIndex).
				SetCanCreate(p.CanCreate).
				SetCanUpdate(p.CanUpdate).
				SetCanDestroy(p.CanDestroy).
				Save(ctx)
		case err != nil:
			return nil, err
		default:
			_, err = existing.Update().
				SetCanIndex(p.CanIndex).
				SetCanCreate(p.CanCreate).
				SetCanUpdate(p.CanUpdate).
				SetCanDestroy(p.CanDestroy).
				Save(ctx)
		}
		if err != nil {
			return nil, err
		}
	}

	return s.roleDetail(ctx, roleID)
}

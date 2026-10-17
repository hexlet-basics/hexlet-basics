package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// The management users surface (legacy `admin/management/users`) is a slim view
// over the same users as AdminUsers: list + get + update, no create or delete.
// It shares their read and write paths (see users.go).

func (s *Server) AdminListManagementUsers(ctx context.Context, params api.AdminListManagementUsersParams) (api.AdminListManagementUsersRes, error) {
	return s.userCrudPage(ctx, params.Page, params.PerPage, params.Q)
}

func (s *Server) AdminGetManagementUser(ctx context.Context, params api.AdminGetManagementUserParams) (api.AdminGetManagementUserRes, error) {
	return s.userCrud(ctx, int(params.ID))
}

// AdminUpdateManagementUser updates a user's editable fields. A null nullable
// field clears the column.
func (s *Server) AdminUpdateManagementUser(ctx context.Context, req *api.UserInput, params api.AdminUpdateManagementUserParams) (api.AdminUpdateManagementUserRes, error) {
	return s.updateUser(ctx, int(params.ID), req)
}

package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
)

// The management users surface (legacy `admin/management/users`) is a slim view
// over the same users table as AdminUsers: list + get + update, no create or
// delete. It reuses the User schema and UserCrud converter.

func (s *Server) AdminListManagementUsers(ctx context.Context, params api.AdminListManagementUsersParams) (*api.UserCrudPage, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.UserQuery { return s.db.User.Query().Order(ent.Desc(user.FieldID)) },
		s.conv.ToUserCruds,
		func(items []api.UserCrud, total, page, perPage int32) *api.UserCrudPage {
			return &api.UserCrudPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

func (s *Server) AdminGetManagementUser(ctx context.Context, params api.AdminGetManagementUserParams) (*api.UserCrud, error) {
	return getOne(ctx, int(params.ID), s.db.User.Get, s.conv.ToUserCrud)
}

// AdminUpdateManagementUser updates a user's editable fields. Email uniqueness
// is enforced by the baseline unique index (409 via the central handler); a
// null nullable field clears the column.
func (s *Server) AdminUpdateManagementUser(ctx context.Context, req *api.UserInput, params api.AdminUpdateManagementUserParams) (*api.UserCrud, error) {
	upd := s.db.User.UpdateOneID(int(params.ID)).SetEmail(req.Email)
	applyNil(req.FirstName.Null, req.FirstName.Value, upd.SetFirstName, upd.ClearFirstName)
	applyNil(req.LastName.Null, req.LastName.Value, upd.SetLastName, upd.ClearLastName)
	applyNil(req.Admin.Null, req.Admin.Value, upd.SetAdmin, upd.ClearAdmin)

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToUserCrud(row)
	return &item, nil
}

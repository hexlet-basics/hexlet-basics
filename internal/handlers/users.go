package handlers

import (
	"context"
	"errors"
	"net/http"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
)

// searchUsersLimit mirrors the legacy typeahead cap (`.limit(20)`).
const searchUsersLimit = 20

// AdminListUsers returns a page of users, newest first.
func (s *Server) AdminListUsers(ctx context.Context, params api.AdminListUsersParams) (api.AdminListUsersRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.UserQuery { return s.db.User.Query().Order(ent.Desc(user.FieldID)) },
		s.conv.ToUserCruds,
		func(items []api.UserCrud, total, page, perPage int32) *api.UserCrudPage {
			return &api.UserCrudPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminSearchUsers is the admin typeahead: a case-insensitive substring match on
// first name, last name, or email (legacy `first_name_or_last_name_or_email_cont`),
// capped at 20 rows. A blank term returns an empty list rather than every user.
func (s *Server) AdminSearchUsers(ctx context.Context, params api.AdminSearchUsersParams) (api.AdminSearchUsersRes, error) {
	if params.Q == "" {
		items := api.AdminSearchUsersOKApplicationJSON{}
		return &items, nil
	}

	rows, err := s.db.User.Query().
		Where(user.Or(
			user.FirstNameContainsFold(params.Q),
			user.LastNameContainsFold(params.Q),
			user.EmailContainsFold(params.Q),
		)).
		Order(ent.Desc(user.FieldID)).
		Limit(searchUsersLimit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := api.AdminSearchUsersOKApplicationJSON(s.conv.ToUserCruds(rows))
	return &items, nil
}

// AdminGetUser returns a single user by id. A missing id returns ent's not-found
// error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetUser(ctx context.Context, params api.AdminGetUserParams) (api.AdminGetUserRes, error) {
	return getOne(ctx, int(params.ID), s.db.User.Get, s.conv.ToUserCrud)
}

// AdminCreateUser creates a user.
//
// Email uniqueness is enforced by the baseline unique index, not a pre-check: a
// violation returns ent's constraint error, mapped to 409 centrally (race-free).
func (s *Server) AdminCreateUser(ctx context.Context, req *api.UserInput) (api.AdminCreateUserRes, error) {
	req.Email = accounts.NormalizeEmail(req.Email)
	row, err := s.db.User.Create().SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// AdminUpdateUser updates a user. A missing id returns ent's not-found error
// (mapped to 404 centrally). The generated SetInput clears a column on null,
// matching the legacy assign_attributes semantics.
func (s *Server) AdminUpdateUser(ctx context.Context, req *api.UserInput, params api.AdminUpdateUserParams) (api.AdminUpdateUserRes, error) {
	row, err := s.updateUser(ctx, int(params.ID), req)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// updateUser is the write shared by the two admin surfaces over users. The email
// is stored the way registration stores it, or the account could not be found
// at sign-in; clearing the administrator flag is refused for the last one.
func (s *Server) updateUser(ctx context.Context, id int, req *api.UserInput) (*ent.User, error) {
	req.Email = accounts.NormalizeEmail(req.Email)
	if req.Admin.Null || !req.Admin.Value {
		if err := accounts.EnsureAnotherAdmin(ctx, s.db, id); err != nil {
			return nil, lastAdminError(err)
		}
	}
	return s.db.User.UpdateOneID(id).SetInput(req).Save(ctx)
}

// lastAdminError gives the refusal its 409; any other error passes through.
func lastAdminError(err error) error {
	if errors.Is(err, accounts.ErrLastAdmin) {
		return withHTTPStatus(http.StatusConflict, err)
	}
	return err
}

// AdminDeleteUser removes a user by id. A missing id returns ent's not-found
// error (mapped to 404 centrally).
func (s *Server) AdminDeleteUser(ctx context.Context, params api.AdminDeleteUserParams) (api.AdminDeleteUserRes, error) {
	if err := accounts.EnsureAnotherAdmin(ctx, s.db, int(params.ID)); err != nil {
		return nil, lastAdminError(err)
	}
	if err := s.db.User.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteUserNoContent{}, nil
}

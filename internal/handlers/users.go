package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
)

// searchUsersLimit mirrors the legacy typeahead cap (`.limit(20)`).
const searchUsersLimit = 20

// AdminListUsers returns a page of users, newest first.
func (s *Server) AdminListUsers(ctx context.Context, params api.AdminListUsersParams) (*api.UserCrudPage, error) {
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
func (s *Server) AdminSearchUsers(ctx context.Context, params api.AdminSearchUsersParams) ([]api.UserCrud, error) {
	if params.Q == "" {
		return []api.UserCrud{}, nil
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

	return s.conv.ToUserCruds(rows), nil
}

// AdminGetUser returns a single user by id. A missing id returns ent's not-found
// error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetUser(ctx context.Context, params api.AdminGetUserParams) (*api.UserCrud, error) {
	return getOne(ctx, int(params.ID), s.db.User.Get, s.conv.ToUserCrud)
}

// AdminCreateUser creates a user.
//
// Email uniqueness is enforced by the baseline unique index, not a pre-check: a
// violation returns ent's constraint error, mapped to 409 centrally (race-free).
func (s *Server) AdminCreateUser(ctx context.Context, req *api.UserInput) (*api.UserCrud, error) {
	row, err := s.db.User.Create().
		SetEmail(req.Email).
		SetNillableFirstName(inputconv.Ptr(req.FirstName)).
		SetNillableLastName(inputconv.Ptr(req.LastName)).
		SetNillableAdmin(inputconv.Ptr(req.Admin)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// AdminUpdateUser updates a user. A missing id returns ent's not-found error
// (mapped to 404 centrally). A null nullable field clears the column, matching
// the legacy assign_attributes semantics.
func (s *Server) AdminUpdateUser(ctx context.Context, req *api.UserInput, params api.AdminUpdateUserParams) (*api.UserCrud, error) {
	upd := s.db.User.UpdateOneID(int(params.ID)).
		SetEmail(req.Email)

	if req.FirstName.Null {
		upd.ClearFirstName()
	} else {
		upd.SetFirstName(req.FirstName.Value)
	}
	if req.LastName.Null {
		upd.ClearLastName()
	} else {
		upd.SetLastName(req.LastName.Value)
	}
	if req.Admin.Null {
		upd.ClearAdmin()
	} else {
		upd.SetAdmin(req.Admin.Value)
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// AdminDeleteUser removes a user by id. A missing id returns ent's not-found
// error (mapped to 404 centrally).
func (s *Server) AdminDeleteUser(ctx context.Context, params api.AdminDeleteUserParams) error {
	return s.db.User.DeleteOneID(int(params.ID)).Exec(ctx)
}

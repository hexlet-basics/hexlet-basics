package handlers

import (
	"context"

	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
)

// The rules around editing users (the stored email, the last administrator)
// live in accounts.Editor, and the reads in accounts.Directory; these handlers
// only translate the contract to their commands and the rows to its models. A
// missing id, a duplicate email and the last administrator reach the central
// ErrorHandler as domain errors and become 404 and 409 there.

// AdminListUsers returns a page of users, newest first.
func (s *Server) AdminListUsers(ctx context.Context, params api.AdminListUsersParams) (api.AdminListUsersRes, error) {
	return s.userCrudPage(ctx, params.Page, params.PerPage)
}

// userCrudPage is the list the two admin surfaces over users share.
func (s *Server) userCrudPage(ctx context.Context, pageParam, perPageParam api.OptInt32) (*api.UserCrudPage, error) {
	return listFrom(ctx, pageParam, perPageParam, s.userDirectory.List, s.conv.ToUserCruds,
		func(items []api.UserCrud, total, page, perPage int32) *api.UserCrudPage {
			return &api.UserCrudPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminSearchUsers is the admin typeahead; a blank term returns an empty list.
func (s *Server) AdminSearchUsers(ctx context.Context, params api.AdminSearchUsersParams) (api.AdminSearchUsersRes, error) {
	rows, err := s.userDirectory.Search(ctx, params.Q)
	if err != nil {
		return nil, err
	}
	items := api.AdminSearchUsersOKApplicationJSON(s.conv.ToUserCruds(rows))
	return &items, nil
}

// AdminGetUser returns a single user by id.
func (s *Server) AdminGetUser(ctx context.Context, params api.AdminGetUserParams) (api.AdminGetUserRes, error) {
	return s.userCrud(ctx, int(params.ID))
}

// userCrud is the single view the two admin surfaces over users share.
func (s *Server) userCrud(ctx context.Context, id int) (*api.UserCrud, error) {
	row, err := s.userDirectory.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// AdminCreateUser creates a user.
func (s *Server) AdminCreateUser(ctx context.Context, req *api.UserInput) (api.AdminCreateUserRes, error) {
	if err := s.selfProtection.AdministratorFlag(actorOfRequest(ctx), false, !req.Admin.Null && req.Admin.Value); err != nil {
		return nil, protection(err)
	}
	row, err := s.userEditor.Create(ctx, userFieldsOf(req))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// AdminUpdateUser updates a user. A null nullable field clears its column.
func (s *Server) AdminUpdateUser(ctx context.Context, req *api.UserInput, params api.AdminUpdateUserParams) (api.AdminUpdateUserRes, error) {
	return s.updateUser(ctx, int(params.ID), req)
}

// updateUser is the write the two admin surfaces over users share. The editor
// stores the email the way registration does and refuses to take the flag from
// the last administrator (ErrLastAdmin); here, before it runs, only an
// administrator may change the administrator flag, and the refusal for the last
// one is given its localized explanation.
func (s *Server) updateUser(ctx context.Context, id int, req *api.UserInput) (*api.UserCrud, error) {
	current, err := s.userDirectory.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	admin := !req.Admin.Null && req.Admin.Value
	if err := s.selfProtection.AdministratorFlag(actorOfRequest(ctx), isAdministrator(current), admin); err != nil {
		return nil, protection(err)
	}
	row, err := s.userEditor.Update(ctx, id, userFieldsOf(req))
	if err != nil {
		return nil, lastAdminError(err)
	}
	item := s.conv.ToUserCrud(row)
	return &item, nil
}

// AdminDeleteUser removes a user by id. The actor's own account is not deleted
// through the admin operations, and the last administrator is refused with a
// localized explanation (the editor returns ErrLastAdmin).
func (s *Server) AdminDeleteUser(ctx context.Context, params api.AdminDeleteUserParams) (api.AdminDeleteUserRes, error) {
	if err := s.selfProtection.OwnAccountDeletion(actorOfRequest(ctx), int(params.ID)); err != nil {
		return nil, protection(err)
	}
	if err := s.userEditor.Delete(ctx, int(params.ID)); err != nil {
		return nil, lastAdminError(err)
	}
	return &api.AdminDeleteUserNoContent{}, nil
}

// userFieldsOf is the contract's UserInput as the editor's command.
func userFieldsOf(in *api.UserInput) accounts.UserFields {
	return accounts.UserFields{
		Email:     in.Email,
		FirstName: inputconv.StringPtr(in.FirstName),
		LastName:  inputconv.StringPtr(in.LastName),
		Admin:     inputconv.Ptr(in.Admin),
	}
}

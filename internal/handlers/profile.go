package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// GetProfile returns the signed-in user's editable profile. The contract
// guarantees a session, so the user is always in the context here.
func (s *Server) GetProfile(ctx context.Context) (*api.User, error) {
	u := MustAuthenticatedUser(ctx)
	profile := s.conv.ToUser(u)

	return &profile, nil
}

// UpdateProfile saves the first and last name. Their length and character
// rules are contract constraints, so the generated server has already refused
// anything legacy's validations would have. A null clears the column, as
// legacy's `update(struct.attributes)` did with a nil attribute.
func (s *Server) UpdateProfile(ctx context.Context, req *api.ProfileInput) (*api.User, error) {
	u := MustAuthenticatedUser(ctx)

	update := u.Update()
	if value, ok := req.FirstName.Get(); ok {
		update.SetFirstName(string(value))
	} else {
		update.ClearFirstName()
	}

	if value, ok := req.LastName.Get(); ok {
		update.SetLastName(string(value))
	} else {
		update.ClearLastName()
	}

	u, err := update.Save(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "save user")
	}

	profile := s.conv.ToUser(u)

	return &profile, nil
}

// DeleteAccount removes the signed-in user's account and signs them out. The
// row stays, in the removed state, as in legacy. Legacy ended only the current
// device's session; here every other device's session stops working too,
// because a session is only as good as the account it names and a removed
// account is refused whatever session it carries.
func (s *Server) DeleteAccount(ctx context.Context) error {
	u := MustAuthenticatedUser(ctx)
	if err := s.remover.Remove(ctx, u.ID); err != nil {
		return err
	}

	if err := s.sessions.End(ctx); err != nil {
		return err
	}

	return nil
}

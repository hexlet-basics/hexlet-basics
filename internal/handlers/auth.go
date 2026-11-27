package handlers

import (
	"context"
	"errors"

	"hexletbasics/internal/account"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/signin"
)

// The ways into an account. What each one does is signin.Gate's; a handler
// reads the request, words the refusals and sets the cookies.

// CreateSession signs in with an email and a password.
func (s *Server) CreateSession(ctx context.Context, req *api.SessionInput) (*api.UserHeaders, error) {
	entered, err := s.signIn.WithPassword(ctx, req.Email, req.Password, s.visitor(ctx))

	return s.entered(entered, err)
}

// CreateUser creates an account and signs it in.
func (s *Server) CreateUser(ctx context.Context, req *api.SignUpInput) (*api.UserHeaders, error) {
	var firstName *string
	if value, ok := req.FirstName.Get(); ok {
		firstName = &value
	}

	entered, err := s.signIn.SignUp(ctx, req.Email, req.Password, firstName, s.visitor(ctx))

	return s.entered(entered, err)
}

// DeleteSession ends the session; the session middleware expires its cookie.
func (s *Server) DeleteSession(ctx context.Context) error {
	err := s.sessions.End(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetCurrentUser answers who is signed in: the user the Identify middleware
// loaded, or anonymous. A missing, ended, stale or removed session is
// anonymous, not an error.
func (s *Server) GetCurrentUser(ctx context.Context, _ api.GetCurrentUserParams) (*api.CurrentUser, error) {
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return anonymousCurrentUser(), nil
	}

	return &api.CurrentUser{User: api.NewNilUser(s.conv.ToUser(u))}, nil
}

func anonymousCurrentUser() *api.CurrentUser {
	var user api.NilUser
	user.SetToNull()

	return &api.CurrentUser{User: user}
}

// visitor is who the request is: the Guest progress its cookie carried.
func (s *Server) visitor(ctx context.Context) signin.Visitor {
	guest, _ := GuestProgress(ctx)

	return signin.Visitor{Guest: guest}
}

// entered answers a way in: the account, and the cookie that clears the Guest
// progress once the gate has credited it to the account. A refusal is worded
// for the form it came from.
func (s *Server) entered(entered signin.Entered, err error) (*api.UserHeaders, error) {
	if err != nil {
		return nil, signInRefusal(err)
	}

	var cookies []string
	if entered.GuestMerged {
		cookies = []string{s.guests.cleared()}
	}

	return &api.UserHeaders{SetCookie: cookies, Response: s.conv.ToUser(entered.User)}, nil
}

// signInRefusal words the gate's refusals: on the form field they belong to,
// or as the dead-link answer of the emailed links. Anything else is a failure.
func signInRefusal(err error) error {
	switch {
	case errors.Is(err, account.ErrInvalidCredentials):
		return invalid("password", localization.WrongCredentials)
	case errors.Is(err, account.ErrPasswordProcessing):
		return invalid("password", localization.PasswordProcessingFailed)
	case errors.Is(err, emailtoken.ErrInvalid):
		return errEmailLinkInvalid
	}

	if field, ok := domain.FieldOf(err); ok && field == "email" {
		return invalid("email", localization.EmailTaken)
	}

	return err
}

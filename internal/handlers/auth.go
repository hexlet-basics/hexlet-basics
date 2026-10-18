package handlers

import (
	"context"
	"errors"
	"time"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/events"
	"hexletbasics/internal/localization"
)

// CreateSession authenticates once and reuses the loaded user for both the
// session and the response, keeping password login free of duplicate database
// queries.
func (s *Server) CreateSession(ctx context.Context, req *api.SessionInput) (*api.UserHeaders, error) {
	u, err := s.credentials.Authenticate(ctx, req.Email, req.Password)
	if errors.Is(err, account.ErrInvalidCredentials) {
		return nil, invalid("password", localization.WrongCredentials)
	}
	if err != nil {
		return nil, err
	}
	return s.signIn(ctx, u)
}

// signIn is what every way of proving who you are ends in for an account that
// already exists — password, Magic Link, Password Reset — so none of them can
// skip recording the sign-in. It then starts the session like sign-up does.
func (s *Server) signIn(ctx context.Context, u *ent.User) (*api.UserHeaders, error) {
	if err := s.events.PublishStandalone(ctx, events.UserSignedIn{
		UserID:          u.ID,
		OccurrenceCount: -1,
		Email:           u.Email,
		Locale:          s.i18n.Locale(ctx),
		OccurredAt:      time.Now().UTC(),
	}); err != nil {
		return nil, err
	}
	return s.startSession(ctx, u)
}

// startSession is the tail every way in shares, sign-up included: the session
// begins and the visitor's guest progress is credited to the account. Keeping it
// in one place is what stops sign-in and sign-up from issuing different
// sessions.
func (s *Server) startSession(ctx context.Context, u *ent.User) (*api.UserHeaders, error) {
	if err := s.sessions.Begin(ctx, u); err != nil {
		return nil, err
	}
	cookies, err := s.mergeGuestProgress(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: cookies, Response: s.conv.ToUser(u)}, nil
}

// mergeGuestProgress credits the visitor's cookie progress to the account that
// just signed in or signed up, and returns the cookie that clears it.
//
// The merge runs in its own transaction rather than the one that begins the
// session: beginning a session writes nothing the merge must be atomic with, and
// sign-up owns its transaction inside the registrar. The merge itself is atomic,
// and the cookie is cleared only after it succeeds, so a failed merge leaves
// the visitor's progress intact to be retried on the next sign-in.
func (s *Server) mergeGuestProgress(ctx context.Context, userID int) ([]string, error) {
	guest, ok := GuestProgress(ctx)
	if !ok || len(guest.Entries) == 0 {
		return nil, nil
	}

	if err := s.progress.MergeGuest(ctx, userID, guest, s.i18n.Locale(ctx)); err != nil {
		return nil, err
	}
	return []string{s.guests.cleared()}, nil
}

// CreateUser creates an account and signs it in.
func (s *Server) CreateUser(ctx context.Context, req *api.SignUpInput) (*api.UserHeaders, error) {
	var firstName *string
	if value, ok := req.FirstName.Get(); ok {
		firstName = &value
	}

	u, err := s.registrar.Register(ctx, account.Registration{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: firstName,
		Locale:    s.i18n.Locale(ctx),
	})
	if err != nil {
		if errors.Is(err, account.ErrPasswordProcessing) {
			return nil, invalid("password", localization.PasswordProcessingFailed)
		}
		if field, ok := domain.FieldOf(err); ok && field == "email" {
			return nil, invalid("email", localization.EmailTaken)
		}
		return nil, err
	}
	return s.startSession(ctx, u)
}

// DeleteSession ends the session; the session middleware expires its cookie.
func (s *Server) DeleteSession(ctx context.Context) error {
	if err := s.sessions.End(ctx); err != nil {
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

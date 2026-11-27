package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// CreateMagicLink emails a Magic Link to the User with this email. It answers
// the same whether or not one exists, so the form cannot be used to learn which
// emails are registered, and it never creates an account.
func (s *Server) CreateMagicLink(ctx context.Context, req *api.EmailInput) error {
	return s.signIn.RequestMagicLink(ctx, req.Email, s.i18n.Locale(ctx))
}

// ConsumeMagicLink signs in the User the link was issued to.
func (s *Server) ConsumeMagicLink(
	ctx context.Context,
	params api.ConsumeMagicLinkParams,
) (*api.UserHeaders, error) {
	entered, err := s.signIn.WithMagicLink(ctx, params.Token, s.visitor(ctx))

	return s.entered(entered, err)
}

// CreatePasswordReminder emails a Password Reset link, answering the same for
// an unknown email.
func (s *Server) CreatePasswordReminder(
	ctx context.Context,
	req *api.EmailInput,
) error {
	return s.signIn.RequestPasswordReset(ctx, req.Email, s.i18n.Locale(ctx))
}

// CheckPasswordResetToken lets the reset page refuse a dead link before the
// learner types a new password into it.
func (s *Server) CheckPasswordResetToken(
	ctx context.Context,
	params api.CheckPasswordResetTokenParams,
) error {
	err := s.signIn.CheckPasswordReset(ctx, params.Token)
	if err != nil {
		return signInRefusal(err)
	}

	return nil
}

// UpdatePassword replaces the password and signs the User in.
func (s *Server) UpdatePassword(
	ctx context.Context,
	req *api.ResetPasswordInput,
	params api.UpdatePasswordParams,
) (*api.UserHeaders, error) {
	entered, err := s.signIn.ResetPassword(ctx, params.Token, req.Password, s.visitor(ctx))

	return s.entered(entered, err)
}

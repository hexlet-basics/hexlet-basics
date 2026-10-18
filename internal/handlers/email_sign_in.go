package handlers

import (
	"context"
	"errors"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
	"hexletbasics/internal/accountemail"
	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/localization"
)

// CreateMagicLink emails a Magic Link to the User with this email. It answers
// the same whether or not one exists, so the form cannot be used to learn which
// emails are registered, and it never creates an account.
func (s *Server) CreateMagicLink(ctx context.Context, req *api.EmailInput) (api.CreateMagicLinkRes, error) {
	if err := s.requestAccountEmail(ctx, emailtoken.MagicLink, req.Email); err != nil {
		return nil, err
	}
	return &api.CreateMagicLinkNoContent{}, nil
}

// ConsumeMagicLink signs in the User the link was issued to. The link may be
// followed again while it lives, as in legacy.
func (s *Server) ConsumeMagicLink(
	ctx context.Context,
	params api.ConsumeMagicLinkParams,
) (api.ConsumeMagicLinkRes, error) {
	u, err := s.emailTokens.Verify(ctx, s.db, emailtoken.MagicLink, params.Token)
	if errors.Is(err, emailtoken.ErrInvalid) {
		return s.linkInvalid(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	return s.signIn(ctx, u)
}

// CreatePasswordReminder emails a Password Reset link. Unlike legacy it no
// longer says that an email is unknown: that would undo the Magic Link form's
// protection against learning which emails are registered.
func (s *Server) CreatePasswordReminder(
	ctx context.Context,
	req *api.EmailInput,
) (api.CreatePasswordReminderRes, error) {
	if err := s.requestAccountEmail(ctx, emailtoken.PasswordReset, req.Email); err != nil {
		return nil, err
	}
	return &api.CreatePasswordReminderNoContent{}, nil
}

// CheckPasswordResetToken lets the reset page refuse a dead link before the
// learner types a new password into it.
func (s *Server) CheckPasswordResetToken(
	ctx context.Context,
	params api.CheckPasswordResetTokenParams,
) (api.CheckPasswordResetTokenRes, error) {
	_, err := s.emailTokens.Verify(ctx, s.db, emailtoken.PasswordReset, params.Token)
	if errors.Is(err, emailtoken.ErrInvalid) {
		return s.linkInvalid(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	return &api.CheckPasswordResetTokenNoContent{}, nil
}

// UpdatePassword replaces the password and signs the User in: they have just
// proved they own the email, so asking for the new password again adds
// nothing. Changing the digest is what kills the link — its fingerprint no
// longer matches — so the link cannot be used twice.
func (s *Server) UpdatePassword(
	ctx context.Context,
	req *api.ResetPasswordInput,
	params api.UpdatePasswordParams,
) (api.UpdatePasswordRes, error) {
	u, err := s.emailTokens.Verify(ctx, s.db, emailtoken.PasswordReset, params.Token)
	if errors.Is(err, emailtoken.ErrInvalid) {
		return s.linkInvalid(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	u, err = s.credentials.SetPassword(ctx, u, req.Password)
	if errors.Is(err, account.ErrPasswordProcessing) {
		return validationError("password", s.i18n.Text(ctx, localization.PasswordProcessingFailed)), nil
	}
	if err != nil {
		return nil, err
	}
	return s.signIn(ctx, u)
}

// requestAccountEmail enqueues the email for a registered address and does
// nothing for an unknown one. The lookup normalizes the email the way registration
// stores it, as legacy's did.
func (s *Server) requestAccountEmail(ctx context.Context, purpose emailtoken.Purpose, email string) error {
	u, err := account.ActiveUserByEmail(ctx, s.db, email)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.jobs.Enqueue(ctx, jobs.AccountEmailArgs{
		Purpose: purpose,
		UserID:  u.ID,
		Locale:  accountemail.Locale(s.i18n.Locale(ctx)),
	})
}

func (s *Server) linkInvalid(ctx context.Context) *api.NotFoundError {
	return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.EmailLinkInvalid)}
}

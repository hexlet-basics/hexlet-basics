package handlers

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/accountemails"
	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/localization"
)

// CreateMagicLink emails a Magic Link to the User with this email. It answers
// the same whether or not one exists, so the form cannot be used to learn which
// emails are registered, and it never creates an account.
func (h *AuthHandler) CreateMagicLink(ctx context.Context, req *api.EmailInput) (api.CreateMagicLinkRes, error) {
	if err := h.requestAccountEmail(ctx, emailtokens.MagicLink, req.Email); err != nil {
		return nil, err
	}
	return &api.CreateMagicLinkNoContent{}, nil
}

// ConsumeMagicLink signs in the User the link was issued to. The link may be
// followed again while it lives, as in legacy.
func (h *AuthHandler) ConsumeMagicLink(
	ctx context.Context,
	params api.ConsumeMagicLinkParams,
) (api.ConsumeMagicLinkRes, error) {
	u, err := h.emailTokens.Verify(ctx, h.db, emailtokens.MagicLink, params.Token)
	if errors.Is(err, emailtokens.ErrInvalid) {
		return h.linkInvalid(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	return h.signIn(ctx, u)
}

// CreatePasswordReminder emails a Password Reset link. Unlike legacy it no
// longer says that an email is unknown: that would undo the Magic Link form's
// protection against learning which emails are registered.
func (h *AuthHandler) CreatePasswordReminder(
	ctx context.Context,
	req *api.EmailInput,
) (api.CreatePasswordReminderRes, error) {
	if err := h.requestAccountEmail(ctx, emailtokens.PasswordReset, req.Email); err != nil {
		return nil, err
	}
	return &api.CreatePasswordReminderNoContent{}, nil
}

// CheckPasswordResetToken lets the reset page refuse a dead link before the
// learner types a new password into it.
func (h *AuthHandler) CheckPasswordResetToken(
	ctx context.Context,
	params api.CheckPasswordResetTokenParams,
) (api.CheckPasswordResetTokenRes, error) {
	_, err := h.emailTokens.Verify(ctx, h.db, emailtokens.PasswordReset, params.Token)
	if errors.Is(err, emailtokens.ErrInvalid) {
		return h.linkInvalid(ctx), nil
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
func (h *AuthHandler) UpdatePassword(
	ctx context.Context,
	req *api.ResetPasswordInput,
	params api.UpdatePasswordParams,
) (api.UpdatePasswordRes, error) {
	u, err := h.emailTokens.Verify(ctx, h.db, emailtokens.PasswordReset, params.Token)
	if errors.Is(err, emailtokens.ErrInvalid) {
		return h.linkInvalid(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return validationError("password", h.i18n.Text(ctx, localization.PasswordProcessingFailed)), nil
	}
	u, err = u.Update().SetPasswordDigest(string(hash)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return h.signIn(ctx, u)
}

// requestAccountEmail enqueues the email for a registered address and does
// nothing for an unknown one. The lookup ignores case, as legacy's did; the
// contract's email format already rejects surrounding spaces.
func (h *AuthHandler) requestAccountEmail(ctx context.Context, purpose emailtokens.Purpose, email string) error {
	u, err := h.db.User.Query().
		Where(user.EmailEqualFold(email)).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return h.emails.EnqueueAccountEmail(ctx, purpose, u.ID, accountemails.Locale(h.i18n.Locale(ctx)))
}

func (h *AuthHandler) linkInvalid(ctx context.Context) *api.NotFoundError {
	return &api.NotFoundError{Message: h.i18n.Text(ctx, localization.EmailLinkInvalid)}
}

// Forwarding methods for the generated api.Handler seam, like the ones in
// auth.go.

func (s *Server) CreateMagicLink(ctx context.Context, req *api.EmailInput) (api.CreateMagicLinkRes, error) {
	return s.auth.CreateMagicLink(ctx, req)
}

func (s *Server) ConsumeMagicLink(
	ctx context.Context,
	params api.ConsumeMagicLinkParams,
) (api.ConsumeMagicLinkRes, error) {
	return s.auth.ConsumeMagicLink(ctx, params)
}

func (s *Server) CreatePasswordReminder(
	ctx context.Context,
	req *api.EmailInput,
) (api.CreatePasswordReminderRes, error) {
	return s.auth.CreatePasswordReminder(ctx, req)
}

func (s *Server) CheckPasswordResetToken(
	ctx context.Context,
	params api.CheckPasswordResetTokenParams,
) (api.CheckPasswordResetTokenRes, error) {
	return s.auth.CheckPasswordResetToken(ctx, params)
}

func (s *Server) UpdatePassword(
	ctx context.Context,
	req *api.ResetPasswordInput,
	params api.UpdatePasswordParams,
) (api.UpdatePasswordRes, error) {
	return s.auth.UpdatePassword(ctx, req, params)
}

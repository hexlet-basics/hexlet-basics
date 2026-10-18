package handlers

import (
	"context"
	"errors"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/selfprotection"
)

// Self-protection (ADR-0020 addendum). Whoever may manage staff must not be
// able to remove their own footing or raise themselves. The rules live in the
// selfprotection module, not in the screen, because the website and the MCP
// server both reach the admin operations through these handlers: a check in the
// screen would bind only the website. The handlers only map its refusals.
//
// A refusal is a 403 (the caller may not do this to themselves) or, for the
// last administrator, a 409 (the state of the system forbids it, whoever asks).
// Both carry a localized Problem Details `detail` and, as errors returned by a
// handler, are recorded by the audit middleware as failed calls.

// refusal is a handler error that carries its own status and a localized,
// client-safe explanation for the Problem Details `detail`. Its status is its
// cause's: a domain error in the chain (ErrForbidden, ErrConflict) that the one
// errorStatus table maps.
type refusal struct {
	message localization.Message
	// data fills the placeholders of message; a value that is itself a
	// localization.Message is translated first (a state's name, say).
	data  map[string]any
	cause error
}

func (e *refusal) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return "refused"
}

func (e *refusal) Unwrap() error { return e.cause }

func forbid(message localization.Message) error {
	return &refusal{message: message, cause: domain.ErrForbidden}
}

// lastAdminError gives the last-administrator refusal its 409 and message; any
// other error passes through.
func lastAdminError(err error) error {
	if errors.Is(err, account.ErrLastAdmin) {
		return &refusal{message: localization.LastAdministrator, cause: err}
	}
	return err
}

// actorOf is the signed-in user an admin operation acts as.
func actorOf(ctx context.Context) *ent.User {
	return MustAuthenticatedUser(ctx)
}

func isAdministrator(u *ent.User) bool {
	return u.Admin != nil && *u.Admin
}

// actorOfRequest is the actor the self-protection rules judge: the signed-in
// user the operation acts as.
func actorOfRequest(ctx context.Context) selfprotection.Actor {
	user := actorOf(ctx)
	return selfprotection.Actor{ID: user.ID, Administrator: isAdministrator(user)}
}

// protection maps the self-protection module's refusals to the 403 the
// handlers answer with, each with its localized explanation; any other error
// passes through.
func protection(err error) error {
	switch {
	case errors.Is(err, selfprotection.ErrOwnStaffRecord):
		return forbid(localization.OwnStaffRecord)
	case errors.Is(err, selfprotection.ErrOwnRoleGrants):
		return forbid(localization.OwnRoleGrants)
	case errors.Is(err, selfprotection.ErrAdministratorFlagRestricted):
		return forbid(localization.AdministratorFlagRestricted)
	case errors.Is(err, selfprotection.ErrOwnAccountDeletion):
		return forbid(localization.OwnAccountDeletion)
	default:
		return err
	}
}

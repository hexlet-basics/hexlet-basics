package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/selfprotection"
)

// Self-protection (ADR-0020 addendum). Whoever may manage staff must not be
// able to remove their own footing or raise themselves. The rules live in the
// selfprotection module, not in the screen, because the website and the MCP
// server both reach the admin operations through these handlers: a check in the
// screen would bind only the website. The handlers only map its refusals.
//
// A refusal is a domain.ErrForbidden (403: the caller may not do this to
// themselves) or, for the last administrator, a domain.ErrConflict (409: the
// state of the system forbids it, whoever asks). The handlers return them as
// they are; detailFor words each with its localized Problem Details `detail`,
// and, as errors returned by a handler, the audit middleware records them as
// failed calls.

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

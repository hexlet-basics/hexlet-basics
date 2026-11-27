package handlers

import (
	"context"

	"hexletbasics/internal/domain"
)

// Self-protection (ADR-0020). Whoever may manage staff must not be able to
// remove their own footing or raise themselves. The rules live in the modules
// that write (account, permission), next to the write each guards, so no
// adapter writes past them; the MCP admin tools reach them through these
// handlers like the website does. A handler passes the actor and maps the
// refusals.
//
// A refusal is a domain.ErrForbidden (403: the caller may not do this to
// themselves) or, for the last administrator, a domain.ErrConflict (409: the
// state of the system forbids it, whoever asks). detailFor words each with its
// localized Problem Details `detail`, and, as errors returned by a handler, the
// audit middleware records them as failed calls.

// actorOfRequest is the signed-in user an admin operation acts as.
func actorOfRequest(ctx context.Context) domain.Actor {
	user := MustAuthenticatedUser(ctx)

	return domain.Actor{UserID: user.ID, Administrator: user.Admin != nil && *user.Admin}
}

package handlers

import (
	"context"
	"errors"
	"net/http"

	"hexletbasics/ent"
	"hexletbasics/ent/staffmember"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/localization"
)

// Self-protection (ADR-0017 addendum). Whoever may manage staff must not be
// able to remove their own footing or raise themselves. The rules live here,
// in the operation layer, because the website and the MCP server both reach the
// admin operations through these handlers: a check in the screen would bind
// only the website.
//
// A refusal is a 403 (the caller may not do this to themselves) or, for the
// last administrator, a 409 (the state of the system forbids it, whoever asks).
// Both carry a localized Problem Details `detail` and, as errors returned by a
// handler, are recorded by the audit middleware as failed calls.

// refusal is a handler error that carries its own status and a localized,
// client-safe explanation for the Problem Details `detail`.
type refusal struct {
	status  int
	message localization.Message
	cause   error
}

func (e *refusal) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return "refused"
}

func (e *refusal) Unwrap() error { return e.cause }

// HTTPStatus lets the central error writer pick the status.
func (e *refusal) HTTPStatus() int { return e.status }

func forbid(message localization.Message) error {
	return &refusal{status: http.StatusForbidden, message: message}
}

// lastAdminError gives the last-administrator refusal its 409 and message; any
// other error passes through.
func lastAdminError(err error) error {
	if errors.Is(err, accounts.ErrLastAdmin) {
		return &refusal{status: http.StatusConflict, message: localization.LastAdministrator, cause: err}
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

// guardStaffRecordChange refuses a change to the actor's own Staff Member
// record (their role, locales or the record itself). ownerOf is the user the
// record belongs to before the change (0 for a record being created) and
// newOwner the user it will belong to after (0 for a delete), so moving
// someone else's record onto oneself is caught as well.
func guardStaffRecordChange(ctx context.Context, ownerOf, newOwner int) error {
	actor := actorOf(ctx)
	if ownerOf == actor.ID || newOwner == actor.ID {
		return forbid(localization.OwnStaffRecord)
	}
	return nil
}

// guardOwnRoleGrants refuses a non-administrator the change of the grants of
// the Staff Role they hold. Delegating staff_roles.update would otherwise let
// its holder grant their own role anything. An administrator may: they can
// already do everything.
func (s *Server) guardOwnRoleGrants(ctx context.Context, roleID int) error {
	actor := actorOf(ctx)
	if isAdministrator(actor) {
		return nil
	}
	holds, err := s.db.StaffMember.Query().
		Where(staffmember.UserID(actor.ID), staffmember.RoleID(roleID)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if holds {
		return forbid(localization.OwnRoleGrants)
	}
	return nil
}

// guardAdministratorFlag lets only an administrator change the administrator
// flag: a staff member who holds users.update (or users.create) could
// otherwise make themselves, or anyone, an administrator. was is the flag
// before the change (false for a new user), will the flag after.
func guardAdministratorFlag(ctx context.Context, was, will bool) error {
	if was == will || isAdministrator(actorOf(ctx)) {
		return nil
	}
	return forbid(localization.AdministratorFlagRestricted)
}

// guardOwnAccountDeletion refuses the deletion of the actor's own account
// through the admin operations.
func guardOwnAccountDeletion(ctx context.Context, id int) error {
	if actorOf(ctx).ID == id {
		return forbid(localization.OwnAccountDeletion)
	}
	return nil
}

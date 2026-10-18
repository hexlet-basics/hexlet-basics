// Package selfprotection holds the rules that keep whoever may manage staff
// from removing their own footing or raising themselves (ADR-0020 addendum).
// They live in a module, not in a screen or a handler body, because the website
// and the MCP server both reach the admin operations through the same handlers:
// a check in one adapter would bind only that one.
package selfprotection

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/staffmember"
	"hexletbasics/internal/domain"
)

// Actor is the signed-in user an admin operation acts as.
type Actor struct {
	ID            int
	Administrator bool
}

// Guard judges an actor's change to staff management against the rules.
type Guard struct {
	db *ent.Client
}

// NewGuard builds the production guard.
func NewGuard(db *ent.Client) *Guard {
	return &Guard{db: db}
}

// StaffRecordChange refuses a change to the actor's own Staff Member record
// (their role, locales or the record itself). recordID is the record being
// changed (0 for one being created) and newOwner the user it will belong to
// after the change (0 for a delete), so moving someone else's record onto
// oneself is caught as well. A record that is not there is domain.ErrNotFound.
func (g *Guard) StaffRecordChange(ctx context.Context, actor Actor, recordID, newOwner int) error {
	owner := 0

	if recordID != 0 {
		record, err := g.db.StaffMember.Get(ctx, recordID)
		if err != nil {
			return domain.FromStore(err)
		}

		owner = record.UserID
	}

	if owner == actor.ID || newOwner == actor.ID {
		return ErrOwnStaffRecord
	}

	return nil
}

// OwnRoleGrants refuses a non-administrator the change of the grants of the
// Staff Role they hold. Delegating staff_roles.update would otherwise let its
// holder grant their own role anything. An administrator may: they can already
// do everything.
func (g *Guard) OwnRoleGrants(ctx context.Context, actor Actor, roleID int) error {
	if actor.Administrator {
		return nil
	}

	holds, err := g.db.StaffMember.Query().
		Where(staffmember.UserID(actor.ID), staffmember.RoleID(roleID)).
		Exist(ctx)
	if err != nil {
		return domain.FromStore(err)
	}

	if holds {
		return ErrOwnRoleGrants
	}

	return nil
}

// AdministratorFlag lets only an administrator change the administrator flag: a
// staff member who holds users.update (or users.create) could otherwise make
// themselves, or anyone, an administrator. was is the flag before the change
// (false for a new user), will the flag after.
func (g *Guard) AdministratorFlag(actor Actor, was, will bool) error {
	if was == will || actor.Administrator {
		return nil
	}

	return ErrAdministratorFlagRestricted
}

// OwnAccountDeletion refuses the deletion of the actor's own account through
// the admin operations.
func (g *Guard) OwnAccountDeletion(actor Actor, userID int) error {
	if actor.ID == userID {
		return ErrOwnAccountDeletion
	}

	return nil
}

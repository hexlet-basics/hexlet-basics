package account

import (
	"context"
	"errors"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/domain"
)

// ErrLastAdmin refuses a change that would leave the system with no
// administrator: nobody could then manage users, staff or roles.
var ErrLastAdmin = errors.Join(domain.ErrConflict, errLastAdmin)

// errLastAdmin is the message of ErrLastAdmin.
var errLastAdmin = errors.New("the last administrator cannot be removed")

// EnsureAnotherAdmin is called before userID stops being an administrator
// (demoted, deleted or removing their own account). It returns ErrLastAdmin when
// userID is one and no other active administrator exists; for anyone else it
// returns nil.
//
// It is a check, not a lock: two administrators demoting each other at the same
// moment can both pass. Administrator changes are rare and made by people, so
// the cost of closing that gap would not be paid back.
func EnsureAnotherAdmin(ctx context.Context, db *ent.Client, userID int) error {
	isAdmin, err := db.User.Query().Where(user.ID(userID), user.Admin(true), NotRemoved()).Exist(ctx)
	if err != nil || !isAdmin {
		return err
	}
	others, err := db.User.Query().Where(user.IDNEQ(userID), user.Admin(true), NotRemoved()).Exist(ctx)
	if err != nil {
		return err
	}
	if !others {
		return ErrLastAdmin
	}
	return nil
}

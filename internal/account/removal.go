package account

import (
	"context"
	"database/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/predicate"
	"hexletbasics/ent/user"
	"hexletbasics/ent/useraccount"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/store"
)

// StateRemoved is the legacy AASM state of a deleted account. The row is kept
// (its progress, reviews and memberships still point at it); only the state
// says it is gone.
const StateRemoved = "removed"

// NotRemoved matches the users who still have an account. The legacy initial
// state is written lazily, so most rows carry NULL rather than "active", and a
// bare state <> 'removed' would drop them all: NULL never compares unequal.
func NotRemoved() predicate.User {
	return user.Or(user.StateIsNil(), user.StateNEQ(StateRemoved))
}

// Remover removes accounts the way legacy `mark_as_removed!` did. Removing an
// account has one meaning whoever asks (ADR-0030): the row stays, with what
// identifies the person erased, because its progress, reviews and memberships
// still point at it.
type Remover struct {
	store store.Transactor
}

// NewRemover builds the production account remover.
func NewRemover(txStore store.Transactor) *Remover {
	return &Remover{store: txStore}
}

// RemoveOwn moves the user to the removed state and erases what identifies them:
// names, email, phone, password and confirmation token (legacy clean_fields),
// and the linked social accounts (legacy remove_accounts). Nulling the email
// and phone is also what lets the same address register again, since the
// unique indexes on both skip NULL. One transaction, so a failure never leaves
// an account half-erased. It is the person's own removal, from their profile;
// the last administrator is refused with ErrLastAdmin.
func (r *Remover) RemoveOwn(ctx context.Context, userID int) error {
	return r.remove(ctx, userID)
}

// Remove removes another user's account on an administrator's behalf, the
// same way. The actor's own account is ErrOwnAccountRemoval (their profile is
// the way to it), a missing one domain.ErrNotFound, and the last administrator
// ErrLastAdmin.
func (r *Remover) Remove(ctx context.Context, actor domain.Actor, userID int) error {
	if actor.UserID == userID {
		return ErrOwnAccountRemoval
	}

	return r.remove(ctx, userID)
}

// remove erases the account in one transaction. An account that is not there,
// or already removed, is domain.ErrNotFound: a removed account is gone.
func (r *Remover) remove(ctx context.Context, userID int) error {
	return r.store.WithinTx(ctx, func(_ *sql.Tx, tx *ent.Client) error {
		active, err := tx.User.Query().Where(user.ID(userID), NotRemoved()).Exist(ctx)
		if err != nil {
			return errs.Wrapf(err, "check user")
		}

		if !active {
			return domain.ErrNotFound
		}

		err = EnsureAnotherAdmin(ctx, tx, userID)
		if err != nil {
			return err
		}

		_, err = tx.UserAccount.Delete().Where(useraccount.UserID(userID)).Exec(ctx)
		if err != nil {
			return errs.Wrapf(err, "delete linked accounts")
		}

		err = tx.User.UpdateOneID(userID).
			SetState(StateRemoved).
			ClearFirstName().
			ClearLastName().
			ClearNickname().
			ClearPasswordDigest().
			ClearConfirmationToken().
			ClearEmail().
			ClearPhone().
			Exec(ctx)
		if err != nil {
			return errs.Wrapf(err, "mark user removed")
		}

		return nil
	})
}

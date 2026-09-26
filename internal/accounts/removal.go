package accounts

import (
	"context"
	"database/sql"
	"fmt"

	"hexletbasics/ent"
	"hexletbasics/ent/predicate"
	"hexletbasics/ent/user"
	"hexletbasics/ent/useraccount"
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

// AccountRemover is the account-deletion seam the handlers depend on.
type AccountRemover interface {
	Remove(ctx context.Context, userID int) error
}

// Remover deletes accounts the way legacy `mark_as_removed!` did.
type Remover struct {
	store store.Transactor
}

// NewRemover builds the production account remover.
func NewRemover(txStore store.Transactor) *Remover {
	return &Remover{store: txStore}
}

// Remove moves the user to the removed state and erases what identifies them:
// names, email, phone, password and confirmation token (legacy clean_fields),
// and the linked social accounts (legacy remove_accounts). Nulling the email
// and phone is also what lets the same address register again, since the
// unique indexes on both skip NULL. One transaction, so a failure never leaves
// an account half-erased.
func (r *Remover) Remove(ctx context.Context, userID int) error {
	return r.store.WithinTx(ctx, func(_ *sql.Tx, tx *ent.Client) error {
		if _, err := tx.UserAccount.Delete().Where(useraccount.UserID(userID)).Exec(ctx); err != nil {
			return fmt.Errorf("delete linked accounts: %w", err)
		}
		err := tx.User.UpdateOneID(userID).
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
			return fmt.Errorf("mark user removed: %w", err)
		}
		return nil
	})
}

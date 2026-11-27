// Package account owns account-creation business operations.
package account

import (
	"context"
	"database/sql"
	"errors"

	"hexletbasics/ent"
	"hexletbasics/internal/store"
)

// ErrPasswordProcessing identifies passwords bcrypt cannot encode.
var ErrPasswordProcessing = errors.New("password processing failed")

// userUniques names the field behind each unique index of the users table, so a
// taken email is an invalid input naming "email" rather than a bare conflict.
func userUniques() []store.Unique {
	return []store.Unique{
		store.UniqueOn("index_users_on_email", "email"),
		store.UniqueOn("index_users_on_LOWER_email", "email"),
	}
}

// Registration is the validated public sign-up input.
type Registration struct {
	Email     string
	Password  string
	FirstName *string
}

// Registrar creates accounts.
type Registrar struct {
	store store.Transactor
}

// NewRegistrar builds the production account registrar.
func NewRegistrar(txStore store.Transactor) *Registrar {
	return &Registrar{store: txStore}
}

// Register creates the user.
func (r *Registrar) Register(ctx context.Context, input Registration) (*ent.User, error) {
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	var u *ent.User

	err = r.store.WithinTx(ctx, func(_ *sql.Tx, txClient *ent.Client) error {
		u, err = txClient.User.Create().
			SetEmail(NormalizeEmail(input.Email)).
			SetPasswordDigest(hash).
			SetNillableFirstName(input.FirstName).
			Save(ctx)
		if err != nil {
			return errs.Wrapf(store.ToDomain(err, userUniques()...), "create user")
		}

		return nil
	})
	if err != nil {
		return nil, errs.Wrapf(err, "save user")
	}

	return u, nil
}

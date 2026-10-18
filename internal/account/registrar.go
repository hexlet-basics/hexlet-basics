// Package account owns account-creation business operations.
package account

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/events"
	"hexletbasics/internal/store"
)

// ErrPasswordProcessing identifies passwords bcrypt cannot encode.
var ErrPasswordProcessing = errors.New("password processing failed")

// userUniques names the field behind each unique index of the users table, so a
// taken email is an invalid input naming "email" rather than a bare conflict.
var userUniques = []domain.Unique{
	domain.UniqueOn("index_users_on_email", "email"),
	domain.UniqueOn("index_users_on_LOWER_email", "email"),
}

// Registration is the validated public sign-up input.
type Registration struct {
	Email     string
	Password  string
	FirstName *string
	Locale    string
}

// Registrar atomically creates an account and its UserSignedUp outbox record.
type Registrar struct {
	store     store.Transactor
	publisher events.TxPublisher
	now       func() time.Time
}

// NewRegistrar builds the production account registrar.
func NewRegistrar(txStore store.Transactor, publisher events.TxPublisher, opts ...RegistrarOption) *Registrar {
	registrar := &Registrar{store: txStore, publisher: publisher, now: time.Now}
	for _, opt := range opts {
		opt(registrar)
	}
	return registrar
}

// RegistrarOption adjusts a Registrar as it is built.
type RegistrarOption func(*Registrar)

// WithClock makes the registrar stamp the sign-up fact with now instead of the
// wall clock, so a test can assert on the exact time the fact carries.
func WithClock(now func() time.Time) RegistrarOption {
	return func(r *Registrar) { r.now = now }
}

// Register creates the user and publishes the corresponding domain fact in the
// same PostgreSQL transaction.
func (r *Registrar) Register(ctx context.Context, input Registration) (*ent.User, error) {
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	var u *ent.User
	err = r.store.WithinTx(ctx, func(tx *sql.Tx, txClient *ent.Client) error {
		u, err = txClient.User.Create().
			SetEmail(NormalizeEmail(input.Email)).
			SetPasswordDigest(hash).
			SetNillableFirstName(input.FirstName).
			Save(ctx)
		if err != nil {
			return errs.Wrapf(domain.FromStore(err, userUniques...), "create user")
		}

		if err := r.publisher.Publish(ctx, tx, events.UserSignedUp{
			UserID:     u.ID,
			Email:      u.Email,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Locale:     input.Locale,
			OccurredAt: r.now().UTC(),
		}); err != nil {
			return errs.Wrapf(err, "publish user signed up")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

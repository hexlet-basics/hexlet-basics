// Package accounts owns account-creation business operations.
package accounts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/internal/events"
	"hexletbasics/internal/store"
)

// ErrPasswordProcessing identifies passwords bcrypt cannot encode.
var ErrPasswordProcessing = errors.New("password processing failed")

// Registration is the validated public sign-up input.
type Registration struct {
	Email     string
	Password  string
	FirstName *string
	Locale    string
}

// UserRegistrar is the small account-registration interface shared by handlers
// and tests. The production Registrar hides hashing, transaction management,
// persistence, and domain-event publication behind this seam.
type UserRegistrar interface {
	Register(ctx context.Context, input Registration) (*ent.User, error)
}

// Registrar atomically creates an account and its UserSignedUp outbox record.
type Registrar struct {
	db        *sql.DB
	publisher events.TxPublisher
	now       func() time.Time
}

// NewRegistrar builds the production account registrar.
func NewRegistrar(db *sql.DB, publisher events.TxPublisher) *Registrar {
	return &Registrar{db: db, publisher: publisher, now: time.Now}
}

// Register creates the user and publishes the corresponding domain fact in the
// same PostgreSQL transaction.
func (r *Registrar) Register(ctx context.Context, input Registration) (_ *ent.User, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPasswordProcessing, err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin account registration: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	txClient := store.NewTxClient(tx)
	u, err := txClient.User.Create().
		SetEmail(input.Email).
		SetPasswordDigest(string(hash)).
		SetNillableFirstName(input.FirstName).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := r.publisher.Publish(ctx, tx, events.UserSignedUp{
		UserID:     u.ID,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Locale:     input.Locale,
		OccurredAt: r.now().UTC(),
	}); err != nil {
		return nil, fmt.Errorf("publish user signed up: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit account registration: %w", err)
	}
	return u, nil
}

package account

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/domain"
)

// ErrInvalidCredentials is the one answer for an unknown email, a removed
// account, an account without a password and a wrong password, so a caller
// cannot tell which of them it hit.
var ErrInvalidCredentials = errors.New("invalid credentials")

// NormalizeEmail is legacy `normalizes :email` (strip, downcase): the form an
// email is stored and looked up in. The unique index on lower(email) already
// forbids two rows differing only by case, so folding here loses no account.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ActiveUserByID loads a user who still has an account. Everything that turns a
// token or a session into a user goes through it, so no path forgets the
// removed state (see NotRemoved).
func ActiveUserByID(ctx context.Context, db *ent.Client, id int) (*ent.User, error) {
	u, err := db.User.Query().Where(user.ID(id), NotRemoved()).WithStaffMember().Only(ctx)

	return u, domain.FromStore(err)
}

// ActiveUserByEmail loads the user registered under email, normalized the way
// registration stores it. A removed account has no email, so it never matches.
func ActiveUserByEmail(ctx context.Context, db *ent.Client, email string) (*ent.User, error) {
	u, err := db.User.Query().Where(user.Email(NormalizeEmail(email)), NotRemoved()).WithStaffMember().Only(ctx)

	return u, domain.FromStore(err)
}

// Credentials owns who a person is by what they know: the password digest, how
// it is made and checked, and which email finds which account. It knows nothing
// of cookies or tokens; the handlers turn its answers into sessions.
type Credentials struct {
	db *ent.Client
}

// NewCredentials builds the credentials module over an ent client.
func NewCredentials(db *ent.Client) *Credentials {
	return &Credentials{db: db}
}

// Authenticate returns the active user whose password matches, or
// ErrInvalidCredentials. Legacy kept only an active? user after
// authenticate_by.
func (c *Credentials) Authenticate(ctx context.Context, email, password string) (*ent.User, error) {
	u, err := ActiveUserByEmail(ctx, c.db, email)
	if ent.IsNotFound(err) {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		return nil, err
	}

	if u.PasswordDigest == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordDigest), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

// SetPassword replaces u's password and returns the saved user. Changing the
// digest is also what kills a Password Reset link, whose fingerprint binds it.
// A password bcrypt cannot encode is ErrPasswordProcessing.
//
// It also stamps when the password changed: a session remembers the stamp it
// began under and stops matching when it moves, which signs the user's other
// sessions out (ADR-0017). The stamp is cut to the database's microseconds so the
// saved user and a later load of the row carry the same value.
func (c *Credentials) SetPassword(ctx context.Context, u *ent.User, password string) (*ent.User, error) {
	digest, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	result, err := u.Update().
		SetPasswordDigest(digest).
		SetPasswordChangedAt(time.Now().UTC().Truncate(time.Microsecond)).
		Save(ctx)
	if err != nil {
		return result, errs.Wrapf(err, "save user")
	}

	return result, nil
}

// hashPassword is the one place a password becomes a digest.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Wrapf(errors.Join(ErrPasswordProcessing, err), "hash password")
	}

	return string(hash), nil
}

// Package emailtokens issues and verifies the tokens carried by Magic Link and
// Password Reset emails.
//
// Tokens are stateless, as they were in legacy: a JWT naming the User, the
// purpose, a fifteen-minute expiry and a fingerprint of the one User attribute
// whose change must kill the link — the email for a Magic Link, the password
// digest for a Password Reset. Verification reloads the User and compares the
// fingerprint, so "dies when the email changes" and "dies once the password is
// reset" need no token table.
package emailtokens

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
)

// Lifetime is how long an emailed link stays valid. Legacy used the same
// fifteen minutes for both flows.
const Lifetime = 15 * time.Minute

// Purpose names what a token may be used for. It is signed into the token, so a
// Magic Link can never be replayed as a Password Reset or the other way round.
type Purpose string

const (
	MagicLink     Purpose = "magic_link"
	PasswordReset Purpose = "password_reset"
)

// ErrInvalid covers every refused token — forged, malformed, expired, for the
// other purpose, for a deleted User, or outlived by the attribute it is bound
// to. Callers answer all of them the same way, so they are not told apart.
var ErrInvalid = errors.New("email token is invalid")

type claims struct {
	Purpose     Purpose `json:"pur"`
	Fingerprint string  `json:"fp"`
	jwt.RegisteredClaims
}

// Tokens signs with a secret of its own (config.EmailTokenSecret).
type Tokens struct {
	secret []byte
	now    func() time.Time
}

// Option adjusts a Tokens.
type Option func(*Tokens)

// WithClock replaces the wall clock, so tests can issue a token that has
// already expired.
func WithClock(now func() time.Time) Option {
	return func(t *Tokens) { t.now = now }
}

// New builds a Tokens over the signing secret.
func New(secret string, opts ...Option) *Tokens {
	t := &Tokens{secret: []byte(secret), now: time.Now}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Issue signs a token for the User.
func (t *Tokens) Issue(purpose Purpose, u *ent.User) (string, error) {
	now := t.now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Purpose:     purpose,
		Fingerprint: fingerprint(purpose, u),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(Lifetime)),
		},
	})
	signed, err := token.SignedString(t.secret)
	if err != nil {
		return "", oops.Wrapf(err, "sign %s token", purpose)
	}
	return signed, nil
}

// Verify returns the User a token was issued to, or ErrInvalid.
func (t *Tokens) Verify(ctx context.Context, db *ent.Client, purpose Purpose, raw string) (*ent.User, error) {
	var parsed claims
	_, err := jwt.ParseWithClaims(raw, &parsed, func(*jwt.Token) (any, error) { return t.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(t.now),
	)
	if err != nil || parsed.Purpose != purpose {
		return nil, ErrInvalid
	}
	id, err := strconv.Atoi(parsed.Subject)
	if err != nil {
		return nil, ErrInvalid
	}
	// A removed account takes its links with it. Nulling the email and the
	// digest already breaks most fingerprints, but not a Password Reset issued
	// to a user who never had a password: its bound value was empty before and
	// is empty after.
	u, err := db.User.Query().Where(user.ID(id), accounts.NotRemoved()).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrInvalid
	}
	if err != nil {
		return nil, oops.Wrapf(err, "load user for %s token", purpose)
	}
	if fingerprint(purpose, u) != parsed.Fingerprint {
		return nil, ErrInvalid
	}
	return u, nil
}

// fingerprint hashes the attribute the token is bound to. The hash keeps the
// email and the password digest out of a token that travels through mailboxes.
func fingerprint(purpose Purpose, u *ent.User) string {
	var bound string
	switch purpose {
	case MagicLink:
		if u.Email != nil {
			bound = *u.Email
		}
	case PasswordReset:
		if u.PasswordDigest != nil {
			bound = *u.PasswordDigest
		}
	}
	sum := sha256.Sum256([]byte(string(purpose) + "\x00" + bound))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

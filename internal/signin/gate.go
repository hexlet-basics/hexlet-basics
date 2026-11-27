// Package signin is where every way into an account ends: a password, a Magic
// Link, a Password Reset and a User sign-up (GLOSSARY). Each proves who the
// visitor is in its own way and then takes the same steps, in the same order,
// so none of them can skip one: the sign-in is recorded (not for a sign-up,
// which the registrar records as a UserSignedUp), the Session begins, and the
// visitor's Guest progress becomes the account's.
//
// The package also asks for the two emailed links, and answers the same
// whether or not the email is registered, so neither form tells a visitor
// which emails have accounts.
//
// It is its own package, not part of auth, because jobs imports auth (the
// expired-session purge) and a sign-in enqueues a job.
package signin

import (
	"context"
	"errors"
	"time"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
	"hexletbasics/internal/accountemail"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/progress"
)

// Gate signs visitors in and up.
type Gate struct {
	db          *ent.Client
	credentials *account.Credentials
	registrar   *account.Registrar
	tokens      *emailtoken.Tokens
	sessions    *auth.Authenticator
	progress    *progress.Progress
	events      events.StandalonePublisher
	jobs        *jobs.Enqueuer
}

// New builds the gate.
func New(
	db *ent.Client,
	credentials *account.Credentials,
	registrar *account.Registrar,
	tokens *emailtoken.Tokens,
	sessions *auth.Authenticator,
	tracker *progress.Progress,
	publisher events.StandalonePublisher,
	enqueuer *jobs.Enqueuer,
) *Gate {
	return &Gate{
		db:          db,
		credentials: credentials,
		registrar:   registrar,
		tokens:      tokens,
		sessions:    sessions,
		progress:    tracker,
		events:      publisher,
		jobs:        enqueuer,
	}
}

// Visitor is who is coming in: the locale of the request, and the Guest
// progress its cookie carried (empty when it carried none).
type Visitor struct {
	Locale string
	Guest  progress.GuestProgress
}

// Entered is an account a visitor has just come into, and whether their Guest
// progress is now the account's: only then may the adapter clear the cookie
// that carried it, so a failed merge leaves the progress to retry.
type Entered struct {
	User        *ent.User
	GuestMerged bool
}

// WithPassword signs in the account with this email and password. A wrong pair
// is account.ErrInvalidCredentials.
func (g *Gate) WithPassword(ctx context.Context, email, password string, visitor Visitor) (Entered, error) {
	u, err := g.credentials.Authenticate(ctx, email, password)
	if err != nil {
		return Entered{}, err
	}

	return g.signIn(ctx, u, visitor)
}

// WithMagicLink signs in the account a Magic Link was issued to. The link may
// be followed again while it lives, as in legacy. A dead or forged link is
// emailtoken.ErrInvalid.
func (g *Gate) WithMagicLink(ctx context.Context, token string, visitor Visitor) (Entered, error) {
	u, err := g.tokens.Verify(ctx, g.db, emailtoken.MagicLink, token)
	if err != nil {
		return Entered{}, err
	}

	return g.signIn(ctx, u, visitor)
}

// CheckPasswordReset says whether a Password Reset link still lives, so the
// reset page can refuse a dead one before a password is typed into it.
func (g *Gate) CheckPasswordReset(ctx context.Context, token string) error {
	_, err := g.tokens.Verify(ctx, g.db, emailtoken.PasswordReset, token)

	return err
}

// ResetPassword replaces the password of the account a Password Reset link was
// issued to and signs it in: the visitor has just proved they own the email.
// Changing the digest is what kills the link, whose fingerprint no longer
// matches, so it cannot be used twice. A dead link is emailtoken.ErrInvalid; a
// password that cannot be hashed is account.ErrPasswordProcessing.
func (g *Gate) ResetPassword(ctx context.Context, token, password string, visitor Visitor) (Entered, error) {
	u, err := g.tokens.Verify(ctx, g.db, emailtoken.PasswordReset, token)
	if err != nil {
		return Entered{}, err
	}

	u, err = g.credentials.SetPassword(ctx, u, password)
	if err != nil {
		return Entered{}, err
	}

	return g.signIn(ctx, u, visitor)
}

// SignUp creates an account (User sign-up) and starts its Session. The
// registrar records the sign-up in the same transaction as the account, so no
// sign-in is recorded here. A taken email is an invalid "email" field; a
// password that cannot be hashed is account.ErrPasswordProcessing.
func (g *Gate) SignUp(ctx context.Context, email, password string, firstName *string, visitor Visitor) (Entered, error) {
	u, err := g.registrar.Register(ctx, account.Registration{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		Locale:    visitor.Locale,
	})
	if err != nil {
		return Entered{}, err
	}

	return g.enter(ctx, u, visitor)
}

// RequestMagicLink emails a Magic Link to the account with this email, and
// does nothing for an unknown one. It never creates an account.
func (g *Gate) RequestMagicLink(ctx context.Context, email, locale string) error {
	return g.requestEmail(ctx, emailtoken.MagicLink, email, locale)
}

// RequestPasswordReset emails a Password Reset link to the account with this
// email, and does nothing for an unknown one. Unlike legacy it does not say an
// email is unknown: that would undo the Magic Link form's protection.
func (g *Gate) RequestPasswordReset(ctx context.Context, email, locale string) error {
	return g.requestEmail(ctx, emailtoken.PasswordReset, email, locale)
}

// signIn is what every way of proving who you are ends in for an account that
// already exists: the sign-in is recorded, then the visitor enters.
func (g *Gate) signIn(ctx context.Context, u *ent.User, visitor Visitor) (Entered, error) {
	err := g.events.PublishStandalone(ctx, events.UserSignedIn{
		UserID:          u.ID,
		OccurrenceCount: -1,
		Email:           u.Email,
		Locale:          visitor.Locale,
		OccurredAt:      time.Now().UTC(),
	})
	if err != nil {
		return Entered{}, err
	}

	return g.enter(ctx, u, visitor)
}

// enter is the tail every way in shares, sign-up included: the Session begins
// and the Guest progress is credited to the account. Keeping it in one place is
// what stops sign-in and sign-up from issuing different sessions.
//
// The merge runs in its own transaction rather than one with the session:
// beginning a session writes nothing the merge must be atomic with, and
// sign-up owns its transaction inside the registrar. The merge itself is
// atomic, and GuestMerged is set only after it succeeds.
func (g *Gate) enter(ctx context.Context, u *ent.User, visitor Visitor) (Entered, error) {
	err := g.sessions.Begin(ctx, u)
	if err != nil {
		return Entered{}, err
	}

	if len(visitor.Guest.Entries) == 0 {
		return Entered{User: u, GuestMerged: false}, nil
	}

	err = g.progress.MergeGuest(ctx, u.ID, visitor.Guest, visitor.Locale)
	if err != nil {
		return Entered{}, err
	}

	return Entered{User: u, GuestMerged: true}, nil
}

// requestEmail enqueues the email for a registered address and does nothing
// for an unknown one. The lookup normalizes the email the way registration
// stores it, as legacy's did.
func (g *Gate) requestEmail(ctx context.Context, purpose emailtoken.Purpose, email, locale string) error {
	u, err := account.ActiveUserByEmail(ctx, g.db, email)
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return g.jobs.Enqueue(ctx, jobs.AccountEmailArgs{
		Purpose: purpose,
		UserID:  u.ID,
		Locale:  accountemail.Locale(locale),
	})
}

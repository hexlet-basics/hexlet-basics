// Package session is how a signed-in visitor stays signed in (ADR-0017): an
// opaque token in an httpOnly cookie that names a row in the sessions table.
//
// It owns the cookie, the token's rotation and expiry, and the question "which
// user does this request belong to". It does not decide who may do what; that is
// read from the account on every request.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
)

const (
	// CookieName is the session cookie, and what the contract's cookie security
	// schemes are named for.
	CookieName = "session"

	// MaxLifetime is how long a session lasts from sign-in, whatever the visitor
	// does. It is also the cookie's lifetime, and the cookie never changes after
	// sign-in: the site's pages are rendered by a server that calls this API on
	// the browser's behalf, so a changed cookie would reach the browser only if
	// that server relayed it. A cookie that is stable cannot be left stale.
	MaxLifetime = 31 * 24 * time.Hour

	// IdleTimeout ends a session nobody has used for a week. Idleness is the
	// moment of the last request the session saw, kept inside the session and
	// refreshed at most once an hour (touchEvery), so a read is not a write.
	IdleTimeout = 7 * 24 * time.Hour
	touchEvery  = time.Hour

	userIDKey   = "user_id"
	epochKey    = "password_epoch"
	lastSeenKey = "last_seen"
)

// ErrUnauthenticated means the request carries no session that names a current
// user: none at all, an ended one, a removed account, or a password changed
// since the session began.
var ErrUnauthenticated = errors.New("session: request is not authenticated")

// Authenticator issues, reads and ends sessions. The zero value is not usable; build
// one with New.
type Authenticator struct {
	sessions *scs.SessionManager
	db       *ent.Client
	store    *store
	clock    func() time.Time
}

type options struct {
	secure bool
	clock  func() time.Time
	onErr  func(http.ResponseWriter, *http.Request, error)
}

// Option adjusts a Authenticator.
type Option func(*options)

// WithSecure marks the session cookie Secure. It is on wherever the site is
// served over https.
func WithSecure(secure bool) Option { return func(o *options) { o.secure = secure } }

// WithClock replaces the clock that decides whether a stored session is still
// live. Tests use it to make a session expire without waiting a week.
func WithClock(clock func() time.Time) Option { return func(o *options) { o.clock = clock } }

// WithErrorFunc routes a failure of the session store (the database is down)
// to the application's error responder, so it answers in the API's shape rather
// than scs's plain-text 500.
func WithErrorFunc(fn func(http.ResponseWriter, *http.Request, error)) Option {
	return func(o *options) { o.onErr = fn }
}

// New builds a Authenticator whose sessions live in db.
func New(db *ent.Client, opts ...Option) *Authenticator {
	cfg := options{clock: time.Now}
	for _, opt := range opts {
		opt(&cfg)
	}

	st := &store{db: db, clock: cfg.clock}
	sm := scs.New()
	sm.Store = st
	// scs expires a session MaxLifetime after sign-in and nothing else; idleness
	// is judged by User from the session's own last_seen.
	sm.Lifetime = MaxLifetime
	// A leaked sessions table must not be a list of working cookies.
	sm.HashTokenInStore = true
	sm.Cookie.Name = CookieName
	sm.Cookie.Secure = cfg.secure
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Persist = true
	if cfg.onErr != nil {
		sm.ErrorFunc = cfg.onErr
	}
	return &Authenticator{sessions: sm, db: db, store: st, clock: cfg.clock}
}

// Middleware loads the session a request's cookie names and writes it back, with
// its cookie, when the handler changed it. A request without a session creates
// nothing: no row and no Set-Cookie until someone signs in.
//
// Everything below it may call Begin, End and User; calling them on a request
// that did not pass through it is a programming error.
func (m *Authenticator) Middleware(next http.Handler) http.Handler {
	return m.sessions.LoadAndSave(next)
}

// Begin signs u in on this request: the token is renewed so a token an attacker
// planted before sign-in is worthless (session fixation), and the session
// remembers who u is and which password they signed in with.
func (m *Authenticator) Begin(ctx context.Context, u *ent.User) error {
	if err := m.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("renew session token: %w", err)
	}
	m.remember(ctx, u)
	return nil
}

// Issue starts a session for u outside any request and returns the cookie value
// that names it. It exists for the caller with no request to hang Begin on: a
// test acting as a signed-in user. A request signs in through Begin.
func (m *Authenticator) Issue(ctx context.Context, u *ent.User) (string, error) {
	ctx, err := m.sessions.Load(ctx, "")
	if err != nil {
		return "", fmt.Errorf("load session: %w", err)
	}
	m.remember(ctx, u)
	token, _, err := m.sessions.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("commit session: %w", err)
	}
	return token, nil
}

// remember records on the session who u is and which password they are using.
func (m *Authenticator) remember(ctx context.Context, u *ent.User) {
	m.sessions.Put(ctx, userIDKey, u.ID)
	m.sessions.Put(ctx, epochKey, passwordEpoch(u))
	m.sessions.Put(ctx, lastSeenKey, m.clock().Unix())
}

// End signs the visitor out: the session row is deleted and the cookie expired.
func (m *Authenticator) End(ctx context.Context) error {
	if err := m.sessions.Destroy(ctx); err != nil {
		return fmt.Errorf("destroy session: %w", err)
	}
	return nil
}

// User is the account this request belongs to, or ErrUnauthenticated.
//
// The row decides, not the session: a removed account is signed out at once, and
// so is a session that began under a password that has since changed. Such a
// session is destroyed on the spot rather than left to answer "no" again.
func (m *Authenticator) User(ctx context.Context) (*ent.User, error) {
	id := m.sessions.GetInt(ctx, userIDKey)
	if id == 0 {
		return nil, ErrUnauthenticated
	}

	u, err := account.ActiveUserByID(ctx, m.db, id)
	if ent.IsNotFound(err) {
		return nil, m.endStale(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("load session user: %w", err)
	}
	if m.sessions.GetInt64(ctx, epochKey) != passwordEpoch(u) {
		return nil, m.endStale(ctx)
	}

	now := m.clock()
	idle := now.Sub(time.Unix(m.sessions.GetInt64(ctx, lastSeenKey), 0))
	if idle > IdleTimeout {
		return nil, m.endStale(ctx)
	}
	if idle >= touchEvery {
		m.sessions.Put(ctx, lastSeenKey, now.Unix())
	}
	return u, nil
}

func (m *Authenticator) endStale(ctx context.Context) error {
	if err := m.End(ctx); err != nil {
		return err
	}
	return ErrUnauthenticated
}

// DeleteExpired removes the sessions whose time has passed and reports how many
// it removed. Nothing reads them once they have expired; this is only
// housekeeping, run by a periodic job.
func (m *Authenticator) DeleteExpired(ctx context.Context) (int, error) {
	return m.store.deleteExpired(ctx)
}

// passwordEpoch is the value a session compares to know its user's password has
// not changed: the moment of the last change, zero for a password never changed.
func passwordEpoch(u *ent.User) int64 {
	if u.PasswordChangedAt == nil {
		return 0
	}
	return u.PasswordChangedAt.UnixMicro()
}

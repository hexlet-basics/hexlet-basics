package session_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	entsession "hexletbasics/ent/session"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/session"
	"hexletbasics/internal/testsupport"
)

// clock is a settable time source shared by a Manager and its test.
type clock struct{ now time.Time }

func (c *clock) Now() time.Time          { return c.now }
func (c *clock) advance(d time.Duration) { c.now = c.now.Add(d) }

// app is the smallest server that exercises a session: sign in as a fixture
// user, ask who you are, sign out. Everything goes through the real middleware.
type app struct {
	t        *testing.T
	db       *ent.Client
	sessions *session.Manager
	handler  http.Handler
}

func newApp(t *testing.T, clk *clock, opts ...session.Option) *app {
	t.Helper()
	db := testsupport.NewClient(t)
	sessions := session.New(db, append([]session.Option{session.WithClock(clk.Now)}, opts...)...)
	a := &app{t: t, db: db, sessions: sessions}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /signin/{label}", func(w http.ResponseWriter, r *http.Request) {
		u := db.User.GetX(r.Context(), fixtures.ID(r.PathValue("label")))
		if err := sessions.Begin(r.Context(), u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		u, err := sessions.User(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(strconv.Itoa(u.ID)))
	})
	mux.HandleFunc("POST /signout", func(w http.ResponseWriter, r *http.Request) {
		if err := sessions.End(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	a.handler = sessions.Middleware(mux)
	return a
}

func (a *app) do(method, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	a.t.Helper()
	req := httptest.NewRequestWithContext(a.t.Context(), method, path, nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	a.handler.ServeHTTP(rec, req)
	return rec
}

// signIn signs the fixture user in and returns the session cookie the server
// set.
func (a *app) signIn(label string, with ...*http.Cookie) *http.Cookie {
	a.t.Helper()
	var planted *http.Cookie
	if len(with) > 0 {
		planted = with[0]
	}
	rec := a.do(http.MethodPost, "/signin/"+label, planted)
	require.Equal(a.t, http.StatusOK, rec.Code)
	return cookieNamed(rec.Result(), session.CookieName)
}

func (a *app) me(cookie *http.Cookie) int {
	a.t.Helper()
	rec := a.do(http.MethodGet, "/me", cookie)
	if rec.Code != http.StatusOK {
		return 0
	}
	id, err := strconv.Atoi(rec.Body.String())
	require.NoError(a.t, err)
	return id
}

func (a *app) rows() int { return a.db.Session.Query().CountX(a.t.Context()) }

func cookieNamed(resp *http.Response, name string) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func newClock() *clock { return &clock{now: time.Now()} }

// A visitor who never signs in costs nothing: no row, no cookie.
func TestAnAnonymousRequestCreatesNoSession(t *testing.T) {
	a := newApp(t, newClock())
	before := a.rows()

	rec := a.do(http.MethodGet, "/me", nil)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Empty(t, rec.Result().Cookies())
	assert.Equal(t, before, a.rows())
}

func TestSigningInStartsASessionTheCookieNames(t *testing.T) {
	a := newApp(t, newClock())

	cookie := a.signIn("user_with_password")

	require.NotNil(t, cookie)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	// The cookie's lifetime is the session's maximum and never changes after
	// sign-in, so a page rendered by a server that does not relay cookies cannot
	// leave the browser holding a stale one.
	assert.InDelta(t, int(session.MaxLifetime.Seconds()), cookie.MaxAge, 5)
	assert.Equal(t, fixtures.ID("user_with_password"), a.me(cookie))
}

// A leaked table must not be a list of working cookies: the row is keyed by a
// digest, never by the value a browser sends.
func TestTheStoredSessionIsKeyedByADigestOfTheToken(t *testing.T) {
	a := newApp(t, newClock())

	cookie := a.signIn("user_with_password")

	exists := a.db.Session.Query().Where(entsession.ID(cookie.Value)).ExistX(t.Context())
	assert.False(t, exists, "the token itself is not stored")
	assert.Equal(t, 1, a.rows())
}

// A token planted before sign-in must be worthless after it (session fixation).
func TestSigningInRenewsAPlantedToken(t *testing.T) {
	a := newApp(t, newClock())
	planted := a.signIn("user_blank")

	fresh := a.signIn("user_with_password", planted)

	assert.NotEqual(t, planted.Value, fresh.Value)
	assert.Equal(t, fixtures.ID("user_with_password"), a.me(fresh))
	assert.Zero(t, a.me(planted), "the old token no longer names anyone")
}

func TestSigningOutEndsTheSessionNotJustTheCookie(t *testing.T) {
	a := newApp(t, newClock())
	cookie := a.signIn("user_with_password")

	rec := a.do(http.MethodPost, "/signout", cookie)
	require.Equal(t, http.StatusOK, rec.Code)

	cleared := cookieNamed(rec.Result(), session.CookieName)
	require.NotNil(t, cleared)
	assert.Negative(t, cleared.MaxAge)
	assert.Zero(t, a.me(cookie), "a copy of the cookie opens nothing")
	assert.Zero(t, a.rows())
}

func TestARemovedAccountIsSignedOutAtOnce(t *testing.T) {
	a := newApp(t, newClock())
	cookie := a.signIn("user_with_password")
	a.db.User.UpdateOneID(fixtures.ID("user_with_password")).SetState(accounts.StateRemoved).ExecX(t.Context())

	assert.Zero(t, a.me(cookie))
	assert.Zero(t, a.rows(), "the dead session is not left to answer no again")
}

// Changing the password ends every session that began under the old one, on
// every device, and leaves the one that began under the new one alone.
func TestChangingThePasswordEndsTheOtherSessions(t *testing.T) {
	a := newApp(t, newClock())
	laptop := a.signIn("user_with_password")
	phone := a.signIn("user_with_password")
	u := a.db.User.GetX(t.Context(), fixtures.ID("user_with_password"))

	_, err := accounts.NewCredentials(a.db).SetPassword(t.Context(), u, "a-new-password")
	require.NoError(t, err)

	assert.Zero(t, a.me(laptop))
	assert.Zero(t, a.me(phone))
	assert.Equal(t, fixtures.ID("user_with_password"), a.me(a.signIn("user_with_password")))
}

// A week without a request ends a session.
func TestASessionNobodyUsesEndsAfterTheIdleWeek(t *testing.T) {
	clk := newClock()
	a := newApp(t, clk)
	cookie := a.signIn("user_with_password")

	clk.advance(session.IdleTimeout - time.Hour)
	require.Equal(t, fixtures.ID("user_with_password"), a.me(cookie))

	clk.advance(session.IdleTimeout + time.Hour)
	assert.Zero(t, a.me(cookie))
}

// Use keeps a session alive, but only to the thirty-first day from sign-in.
func TestAnActiveSessionStillEndsAtItsMaximumAge(t *testing.T) {
	clk := newClock()
	a := newApp(t, clk)
	cookie := a.signIn("user_with_password")

	for range 5 {
		clk.advance(6 * 24 * time.Hour)
		require.Equal(t, fixtures.ID("user_with_password"), a.me(cookie), "active, so extended")
	}

	clk.advance(2 * 24 * time.Hour)
	assert.Zero(t, a.me(cookie), "thirty-two days after sign-in")
}

// Reading is not writing: a busy session records that it was seen at most once
// an hour, and its end, which the browser's cookie also carries, never moves.
func TestABusySessionIsWrittenBackAtMostOnceAnHour(t *testing.T) {
	clk := newClock()
	a := newApp(t, clk)
	cookie := a.signIn("user_with_password")
	stored := func() (data []byte, expiry time.Time) {
		row := a.db.Session.Query().OnlyX(t.Context())
		return row.Data, row.Expiry
	}
	signedInData, signedInExpiry := stored()

	clk.advance(30 * time.Minute)
	require.NotZero(t, a.me(cookie))
	data, _ := stored()
	assert.Equal(t, signedInData, data, "within the hour nothing is written")

	clk.advance(90 * time.Minute)
	require.NotZero(t, a.me(cookie))
	data, expiry := stored()
	assert.NotEqual(t, signedInData, data, "after the hour the sighting is recorded")
	assert.Equal(t, signedInExpiry, expiry, "the end of the session does not move")
}

func TestDeleteExpiredRemovesOnlyTheSessionsThatRanOut(t *testing.T) {
	clk := newClock()
	a := newApp(t, clk)
	a.signIn("user_with_password")
	a.signIn("user_blank")
	require.Equal(t, 2, a.rows())

	clk.advance(session.MaxLifetime - time.Hour)
	removed, err := a.sessions.DeleteExpired(t.Context())
	require.NoError(t, err)
	assert.Zero(t, removed, "a session within its lifetime stays, used or not")

	clk.advance(2 * time.Hour)
	removed, err = a.sessions.DeleteExpired(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, removed)
	assert.Zero(t, a.rows())
}

func TestTheCookieIsSecureWhereTheSiteIsHTTPS(t *testing.T) {
	a := newApp(t, newClock(), session.WithSecure(true))

	cookie := a.signIn("user_with_password")

	assert.True(t, cookie.Secure)
}

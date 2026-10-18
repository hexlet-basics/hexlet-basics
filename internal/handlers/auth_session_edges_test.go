package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/auth"
	"hexletbasics/internal/testsupport"
)

// signedIn registers an account and returns the session cookie it was given.
func signedIn(t *testing.T, router http.Handler, email string) *http.Cookie {
	t.Helper()
	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+email+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	cookie := sessionCookie(resp)
	require.NotNil(t, cookie)
	return cookie
}

// A cookie that is not a token this server signed does not open the admin
// area: the answer is 401, as for any protected operation, never a 403 that
// would admit the visitor was recognized.
func TestAGarbledSessionCookieDoesNotOpenTheAdminArea(t *testing.T) {
	router := newAuthRouter(t)
	garbled := &http.Cookie{Name: auth.CookieName, Value: "not-a-token"} //nolint:gosec // G124: a request cookie in a test, the attributes belong to Set-Cookie

	resp := do(t, router, http.MethodGet, "/api/admin/course_categories", "", garbled)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// A signed-in learner who is not staff is refused with 403 rather than 401:
// they are known, they just may not.
func TestASignedInLearnerWhoIsNotStaffIsForbiddenTheAdminArea(t *testing.T) {
	router := newAuthRouter(t)
	cookie := signedIn(t, router, "not-staff@example.com")

	resp := do(t, router, http.MethodGet, "/api/admin/course_categories", "", cookie)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// A browser tells the server where a request came from. A write sent by a page
// on another site is refused before the session is even looked at, whatever
// cookie rides along; the same write from the site itself, or from a client that
// is not a browser, goes through.
func TestACrossSiteWriteIsRefused(t *testing.T) {
	router := newAuthRouter(t)
	cookie := signedIn(t, router, "cross-site@example.com")
	const edit = `{"firstName":"Grace","lastName":null}`

	patch := func(site string) int {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, "/api/account/profile",
			strings.NewReader(edit))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(cookie)
		if site != "" {
			req.Header.Set("Sec-Fetch-Site", site)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}

	assert.Equal(t, http.StatusForbidden, patch("cross-site"))
	assert.Equal(t, http.StatusForbidden, patch("same-site"))
	assert.Equal(t, http.StatusOK, patch("same-origin"))
	assert.Equal(t, http.StatusOK, patch(""), "a client that sends no origin headers is not a browser")
}

// Signing up and signing in are different ways in, and must not be different
// sessions: the cookie that comes back has the same attributes either way.
func TestSignUpAndSignInStartSessionsOfOneShape(t *testing.T) {
	router := newAuthRouter(t)
	const email, password = "one-shape@example.com", "s3cret-pass"

	signUp := sessionCookie(do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+email+`","password":"`+password+`"}`, nil))
	signIn := sessionCookie(do(t, router, http.MethodPost, "/api/session",
		`{"email":"`+email+`","password":"`+password+`"}`, nil))

	require.NotNil(t, signUp)
	require.NotNil(t, signIn)
	assert.NotEqual(t, signUp.Value, signIn.Value, "each way in is its own session")
	assert.Equal(t, signUp.Path, signIn.Path)
	assert.Equal(t, signUp.HttpOnly, signIn.HttpOnly)
	assert.Equal(t, signUp.Secure, signIn.Secure)
	assert.Equal(t, signUp.SameSite, signIn.SameSite)
	assert.InDelta(t, signUp.MaxAge, signIn.MaxAge, 5)
}

// A visitor who is not signed in costs the server nothing: no session row and
// no cookie, however many pages they read.
func TestAnAnonymousReadStartsNoSession(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	router := newAuthRouterWithDB(t, db, tx, transactor)
	before := db.Session.Query().CountX(t.Context())

	resp := do(t, router, http.MethodGet, "/api/me", "", nil)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, resp.Cookies())
	assert.Equal(t, before, db.Session.Query().CountX(t.Context()))
}

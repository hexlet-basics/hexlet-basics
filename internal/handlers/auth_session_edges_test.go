package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// signedIn registers an account and returns the session cookies it was given.
func signedIn(t *testing.T, router http.Handler, email string) (session *http.Cookie, xsrf string) {
	t.Helper()
	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+email+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	session = jwtCookie(resp)
	require.NotNil(t, session)
	token := xsrfCookieFromResponse(resp)
	require.NotNil(t, token)
	return session, token.Value
}

// A cookie that is not a token this server signed does not open the admin
// area: the answer is 401, as for any protected operation, never a 403 that
// would admit the visitor was recognized.
func TestAGarbledSessionCookieDoesNotOpenTheAdminArea(t *testing.T) {
	router := newAuthRouter(t)
	garbled := &http.Cookie{Name: "JWT", Value: "not-a-token"}

	resp := do(t, router, http.MethodGet, "/api/admin/course_categories", "", garbled)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// A signed-in learner who is not staff is refused with 403 rather than 401:
// they are known, they just may not.
func TestASignedInLearnerWhoIsNotStaffIsForbiddenTheAdminArea(t *testing.T) {
	router := newAuthRouter(t)
	session, _ := signedIn(t, router, "not-staff@example.com")

	resp := do(t, router, http.MethodGet, "/api/admin/course_categories", "", session)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// Changing something needs the double-submit token as well as the cookie: a
// page on another site can send the cookie, but cannot read the token.
func TestAnUnsafeRequestNeedsTheSubmittedXsrfToken(t *testing.T) {
	router := newAuthRouter(t)
	session, xsrf := signedIn(t, router, "xsrf-needed@example.com")
	const edit = `{"firstName":"Grace","lastName":null}`

	missing := doWithXSRF(t, router, http.MethodPatch, "/api/account/profile", edit, session, "")
	wrong := doWithXSRF(t, router, http.MethodPatch, "/api/account/profile", edit, session, "not-the-token")
	right := doWithXSRF(t, router, http.MethodPatch, "/api/account/profile", edit, session, xsrf)

	assert.Equal(t, http.StatusUnauthorized, missing.StatusCode)
	assert.Equal(t, http.StatusUnauthorized, wrong.StatusCode)
	assert.Equal(t, http.StatusOK, right.StatusCode)
}

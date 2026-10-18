package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/testsupport"
)

// An email already registered is refused as taken, on the email field, and no
// second account is made.
func TestSignUpRefusesATakenEmail(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	taken := h.DB.User.Query().FirstX(t.Context())
	before := h.DB.User.Query().CountX(t.Context())

	_, err := h.Client.CreateUser(t.Context(), &api.SignUpInput{
		Email:    *taken.Email,
		Password: "correct horse battery staple",
	})
	problem := problemWithStatus(t, err, http.StatusUnprocessableEntity)
	assert.Contains(t, problem.Response.Errors.Value, "email")
	assert.Equal(t, before, h.DB.User.Query().CountX(t.Context()))
}

// A password the hash cannot take (bcrypt stops at 72 bytes) is refused on the
// password field rather than failing the request.
func TestSignUpRefusesAPasswordTheHashCannotTake(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.CreateUser(t.Context(), &api.SignUpInput{
		Email:    "long-password@example.com",
		Password: strings.Repeat("a", 100),
	})
	problem := problemWithStatus(t, err, http.StatusUnprocessableEntity)
	assert.Contains(t, problem.Response.Errors.Value, "password")
	assert.False(t, h.DB.User.Query().Where(user.Email("long-password@example.com")).ExistX(t.Context()))
}

// An operation that needs a signed-in user and is reached without one is a
// programming error in the contract's security wiring, not a user mistake: it
// panics loudly instead of continuing as nobody.
func TestMustAuthenticatedUserPanicsWithoutOne(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() { handlers.MustAuthenticatedUser(context.Background()) })
}

// A session cookie that is not a token this server signed is no session: the
// visitor is anonymous on a public read and refused on a protected one, never
// an error page.
func TestAGarbledSessionCookieIsNoSession(t *testing.T) {
	t.Parallel()
	router := newAuthRouter(t)
	garbled := &http.Cookie{Name: auth.CookieName, Value: "not-a-token"} //nolint:gosec // G124: a request cookie in a test, the attributes belong to Set-Cookie

	me := do(t, router, http.MethodGet, "/api/me", "", garbled)
	require.Equal(t, http.StatusOK, me.StatusCode)
	var current struct {
		User *struct{} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(me.Body).Decode(&current))
	assert.Nil(t, current.User)

	lead := do(t, router, http.MethodPost, "/api/leads",
		`{"contactMethod":"phone","contactValue":"+10000000000","ymClientId":null}`, garbled)
	assert.Equal(t, http.StatusUnauthorized, lead.StatusCode)
}

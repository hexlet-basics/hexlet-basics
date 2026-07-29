package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/testsupport"
)

// newAuthRouter builds the generated ogen server over a transaction-bound ent client,
// proving the auth routes are implemented through the contract seam.
func newAuthRouter(t *testing.T) http.Handler {
	t.Helper()
	db := testsupport.NewClient(t)
	translator := testsupport.NewTranslator(t)
	handler := handlers.NewServer(
		db,
		&config.Config{JWTSecret: "test-secret"},
		&testsupport.RecordingEnqueuer{DB: db},
		translator,
	)
	server, err := api.NewServer(
		handler,
		api.WithErrorHandler(handlers.NewAPIErrorHandler(translator)),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	require.NoError(t, err)
	return translator.Middleware(server)
}

// jwtCookie returns the JWT auth cookie from a response, or nil if absent.
func jwtCookie(resp *http.Response) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == "JWT" {
			return c
		}
	}
	return nil
}

func do(t *testing.T, router http.Handler, method, path, body string, cookie *http.Cookie) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Result()
}

func TestAuthRegisterLoginFlow(t *testing.T) {
	router := newAuthRouter(t)
	const email = "auth-flow@example.com"
	const password = "s3cret-pass"

	// Register sets the cookie and echoes the created user.
	resp := do(t, router, http.MethodPost, "/users",
		`{"firstName":"Ada","email":"`+email+`","password":"`+password+`"}`, nil)
	if resp.StatusCode != http.StatusCreated {
		var failure any
		_ = json.NewDecoder(resp.Body).Decode(&failure)
		t.Fatalf("register returned %d: %#v", resp.StatusCode, failure)
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	cookie := jwtCookie(resp)
	require.NotNil(t, cookie, "register must set the JWT cookie")

	var created struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, email, created.Email)

	// The cookie resolves the current user via /me.
	resp = do(t, router, http.MethodGet, "/me", "", cookie)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var me struct {
		User *struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&me))
	require.NotNil(t, me.User)
	assert.Equal(t, email, me.User.Email)

	// Without the cookie, /me is anonymous, not an error.
	resp = do(t, router, http.MethodGet, "/me", "", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var anon struct {
		User *struct{} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&anon))
	assert.Nil(t, anon.User)

	// Login with the registered credentials succeeds and re-issues the cookie.
	resp = do(t, router, http.MethodPost, "/session",
		`{"email":"`+email+`","password":"`+password+`"}`, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, jwtCookie(resp))

	// A wrong password is rejected with a 422 validation error.
	resp = do(t, router, http.MethodPost, "/session",
		`{"email":"`+email+`","password":"wrong"}`, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAuthLogoutClearsCookie(t *testing.T) {
	router := newAuthRouter(t)

	resp := do(t, router, http.MethodDelete, "/session", "", nil)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Reset writes an expired JWT cookie so the browser drops it.
	cookie := jwtCookie(resp)
	require.NotNil(t, cookie)
	assert.True(t, cookie.MaxAge < 0 || cookie.Value == "",
		"logout must expire or empty the JWT cookie")
}

func TestAuthValidationErrorUsesRequestLocale(t *testing.T) {
	router := newAuthRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/session",
		strings.NewReader(`{"email":"missing@example.com","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "ru-RU")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var body api.ValidationError
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t,
		[]string{"Неверный адрес электронной почты или пароль"},
		body.Errors["password"],
	)
}

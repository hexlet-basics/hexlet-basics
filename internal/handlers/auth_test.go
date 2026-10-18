package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

// newAuthRouter builds the generated ogen server over a transaction-bound ent client,
// proving the auth routes are implemented through the contract seam.
func newAuthRouter(t *testing.T) http.Handler {
	t.Helper()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	return newAuthRouterWithDB(t, db, tx, transactor)
}

func newAuthRouterWithDB(t *testing.T, db *ent.Client, tx *sql.Tx, transactor store.Transactor) http.Handler {
	t.Helper()
	translator := testsupport.NewTranslator(t)
	errorHandler := testsupport.NewAPIErrorHandler(t, translator)
	services := testsupport.NewServices(t, tx, transactor)
	mocks := testsupport.NewMocks()
	handler := handlers.NewServer(handlers.Deps{
		DB:      db,
		Config:  &config.Config{GuestCookieSecret: "test-secret", EmailTokenSecret: "test-email-secret"},
		Starter: services.Starter,
		Jobs:    services.Jobs,
		// The real progress module: the check is a public operation, so these
		// tests reach it while asserting what the contract protects.
		Progress:  progress.New(db, transactor, services.Events, mocks.Runner),
		Sessions:  auth.New(db),
		Registrar: accounts.NewRegistrar(transactor, services.Events),
		Remover:   accounts.NewRemover(transactor),
		Events:    services.Events,
		I18n:      translator,
		Errors:    errorHandler,
	})
	server, err := api.NewServer(
		handler,
		handler.Security(),
		api.WithErrorHandler(errorHandler.Write),
		api.WithMiddleware(handlers.NewAuditMiddleware(transactor, translator, errorHandler)),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	require.NoError(t, err)
	// The same pipeline as production: cross-origin protection, the session, then
	// the optional identity.
	return translator.Middleware(handler.Security().Pipeline(server))
}

// sessionCookie returns the session cookie from a response, or nil if absent.
func sessionCookie(resp *http.Response) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == auth.CookieName {
			return c
		}
	}
	return nil
}

func do(t *testing.T, router http.Handler, method, path, body string, cookie *http.Cookie) *http.Response {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Result()
}

func TestAuthRegisterLoginFlow(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	router := newAuthRouterWithDB(t, db, tx, transactor)
	const email = "auth-flow@example.com"
	const password = "s3cret-pass"

	// Register sets the cookie and echoes the created user.
	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+email+`","password":"`+password+`"}`, nil)
	if resp.StatusCode != http.StatusCreated {
		var failure any
		_ = json.NewDecoder(resp.Body).Decode(&failure)
		t.Fatalf("register returned %d: %#v", resp.StatusCode, failure)
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	cookie := sessionCookie(resp)
	require.NotNil(t, cookie, "register must start a session")
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	assert.Positive(t, cookie.MaxAge, "the session cookie outlives the browser window")

	var created struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, email, created.Email)

	// The cookie resolves the current user via /me.
	resp = do(t, router, http.MethodGet, "/api/me", "", cookie)
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
	resp = do(t, router, http.MethodGet, "/api/me", "", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var anon struct {
		User *struct{} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&anon))
	assert.Nil(t, anon.User)

	// Login with the registered credentials succeeds and re-issues the cookie.
	resp = do(t, router, http.MethodPost, "/api/session",
		`{"email":"`+email+`","password":"`+password+`"}`, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	loginCookie := sessionCookie(resp)
	require.NotNil(t, loginCookie)

	// Login and signup issue equivalent cookies that resolve through /me.
	resp = do(t, router, http.MethodGet, "/api/me", "", loginCookie)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&me))
	require.NotNil(t, me.User)
	assert.Equal(t, email, me.User.Email)

	// A wrong password is rejected with a 422 validation error.
	resp = do(t, router, http.MethodPost, "/api/session",
		`{"email":"`+email+`","password":"wrong"}`, nil)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var wrongPassword api.ValidationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&wrongPassword))

	// Unknown users produce the same response and cannot be enumerated.
	resp = do(t, router, http.MethodPost, "/api/session",
		`{"email":"missing-auth-flow@example.com","password":"wrong"}`, nil)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var missingUser api.ValidationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&missingUser))
	assert.Equal(t, wrongPassword, missingUser)

	// Accounts without a password digest are also rejected as invalid
	// credentials rather than causing a panic or exposing account state.
	resp = do(t, router, http.MethodPost, "/api/session",
		`{"email":"blank-user@example.com","password":"wrong"}`, nil)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var passwordlessUser api.ValidationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&passwordlessUser))
	assert.Equal(t, wrongPassword, passwordlessUser)
}

func TestCurrentUserSurvivesEmailChange(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	router := newAuthRouterWithDB(t, db, tx, transactor)
	const oldEmail = "jwt-old-email@example.com"
	const newEmail = "jwt-new-email@example.com"

	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+oldEmail+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	cookie := sessionCookie(resp)
	require.NotNil(t, cookie)

	var created struct {
		ID int32 `json:"id"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	_, err := db.User.UpdateOneID(int(created.ID)).SetEmail(newEmail).Save(t.Context())
	require.NoError(t, err)

	resp = do(t, router, http.MethodGet, "/api/me", "", cookie)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var me struct {
		User *struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&me))
	require.NotNil(t, me.User)
	assert.Equal(t, newEmail, me.User.Email)
}

func TestAuthLogoutClearsCookie(t *testing.T) {
	router := newAuthRouter(t)

	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"logout-auth-flow@example.com","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	signedInCookie := sessionCookie(resp)
	require.NotNil(t, signedInCookie)

	resp = do(t, router, http.MethodDelete, "/api/session", "", nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode, "signing out needs a session")

	resp = do(t, router, http.MethodDelete, "/api/session", "", signedInCookie)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// The browser is told to drop the cookie...
	cleared := sessionCookie(resp)
	require.NotNil(t, cleared)
	assert.True(t, cleared.MaxAge < 0 || cleared.Value == "",
		"logout must expire or empty the session cookie")

	// ...and the session itself is gone, so the cookie a browser kept (or an
	// attacker copied) opens nothing.
	resp = do(t, router, http.MethodGet, "/api/me", "", signedInCookie)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var me struct {
		User *struct{} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&me))
	assert.Nil(t, me.User)
}

func TestContractSecurityProtectsParityRoutesAndKeepsLessonCheckPublic(t *testing.T) {
	router := newAuthRouter(t)

	resp := do(t, router, http.MethodPost, "/api/leads",
		`{"contactMethod":"phone","contactValue":"+10000000000","ymClientId":null}`, nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp = do(t, router, http.MethodPost, "/api/lessons/1/check",
		`{"code":"puts 1","versionId":1}`, nil)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAdminAuthorizationUsesCurrentDatabaseValue(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	router := newAuthRouterWithDB(t, db, tx, transactor)
	const email = "admin-revocation@example.com"

	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+email+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	sessionCookieValue := sessionCookie(resp)
	require.NotNil(t, sessionCookieValue)

	u, err := db.User.Query().Where(user.Email(email)).Only(t.Context())
	require.NoError(t, err)

	resp = do(t, router, http.MethodGet, "/api/admin/course_categories", "", sessionCookieValue)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	_, err = db.User.UpdateOneID(u.ID).SetAdmin(true).Save(t.Context())
	require.NoError(t, err)
	resp = do(t, router, http.MethodGet, "/api/admin/course_categories", "", sessionCookieValue)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, err = db.User.UpdateOneID(u.ID).SetAdmin(false).Save(t.Context())
	require.NoError(t, err)
	resp = do(t, router, http.MethodGet, "/api/admin/course_categories", "", sessionCookieValue)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestAuthValidationErrorUsesRequestLocale(t *testing.T) {
	router := newAuthRouter(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/session",
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

func TestSignUpRejectsShortPassword(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	router := newAuthRouterWithDB(t, db, tx, transactor)

	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":null,"email":"short-password@example.com","password":"12345"}`, nil)

	// The contract's minLength is enforced by the generated server, which
	// rejects it as a bad request like any other contract constraint.
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	exists, err := db.User.Query().Where(user.Email("short-password@example.com")).Exist(t.Context())
	require.NoError(t, err)
	assert.False(t, exists)
}

// The profile's name rules are contract constraints: the generated server
// refuses them as a bad request, and only a valid name is saved. Raw HTTP, so
// the bodies are exactly what a client that skips the form could send.
func TestUpdateProfileRejectsInvalidNames(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	router := newAuthRouterWithDB(t, db, tx, transactor)
	const email = "profile-validation@example.com"

	resp := do(t, router, http.MethodPost, "/api/users",
		`{"firstName":"Ada","email":"`+email+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	cookie := sessionCookie(resp)
	require.NotNil(t, cookie)

	for _, body := range []string{
		`{"firstName":"` + strings.Repeat("a", 41) + `","lastName":null}`,
		`{"firstName":"Ada","lastName":"Love@lace"}`,
		`{"firstName":"Ada+","lastName":null}`,
	} {
		resp = do(t, router, http.MethodPatch, "/api/account/profile", body, cookie)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, body)
	}
	u := db.User.Query().Where(user.Email(email)).OnlyX(t.Context())
	assert.Equal(t, "Ada", *u.FirstName)
	assert.Nil(t, u.LastName)

	// Forty characters, and a blank last name, are both allowed.
	resp = do(t, router, http.MethodPatch, "/api/account/profile",
		`{"firstName":"`+strings.Repeat("a", 40)+`","lastName":""}`, cookie)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	u = db.User.GetX(t.Context(), u.ID)
	assert.Equal(t, strings.Repeat("a", 40), *u.FirstName)
	assert.Empty(t, *u.LastName)
}

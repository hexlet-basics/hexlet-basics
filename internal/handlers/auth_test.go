package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
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
	db, transactor := testsupport.NewClientWithTransactor(t)
	return newAuthRouterWithDB(t, db, transactor)
}

func newAuthRouterWithDB(t *testing.T, db *ent.Client, transactor store.Transactor) http.Handler {
	t.Helper()
	translator := testsupport.NewTranslator(t)
	errorHandler := testsupport.NewAPIErrorHandler(t, translator)
	enqueuer := &testsupport.RecordingEnqueuer{DB: db}
	handler := handlers.NewServer(
		db,
		&config.Config{JWTSecret: "test-secret"},
		enqueuer,
		enqueuer,
		// The real progress module: the check is a public operation, so these
		// tests reach it while asserting what the contract protects.
		progress.New(db, transactor, &testsupport.RecordingEventPublisher{}, testsupport.NewStubExerciseRunner()),
		testsupport.NewRecordingRegistrar(db),
		&testsupport.RecordingEventPublisher{},
		translator,
		errorHandler,
	)
	server, err := api.NewServer(
		handler,
		handler.AuthHandler(),
		api.WithErrorHandler(errorHandler.Write),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	require.NoError(t, err)
	return translator.Middleware(handler.AuthHandler().Trace(server))
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

// xsrfCookie returns go-pkgz/auth's script-readable double-submit cookie.
func xsrfCookieFromResponse(resp *http.Response) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == "XSRF-TOKEN" {
			return c
		}
	}
	return nil
}

func do(t *testing.T, router http.Handler, method, path, body string, cookie *http.Cookie) *http.Response {
	return doWithXSRF(t, router, method, path, body, cookie, "")
}

func doWithXSRF(
	t *testing.T,
	router http.Handler,
	method, path, body string,
	cookie *http.Cookie,
	xsrf string,
) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	if xsrf != "" {
		req.Header.Set("X-XSRF-TOKEN", xsrf)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Result()
}

func TestAuthRegisterLoginFlow(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	router := newAuthRouterWithDB(t, db, transactor)
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
	xsrf := xsrfCookieFromResponse(resp)
	require.NotNil(t, xsrf, "register must set the XSRF cookie")
	assert.True(t, cookie.HttpOnly)
	assert.False(t, xsrf.HttpOnly)
	assert.NotEmpty(t, xsrf.Value)

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

	userQueries := 0
	db.User.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
			userQueries++
			return next.Query(ctx, query)
		})
	}))

	// Login with the registered credentials succeeds and re-issues the cookie.
	resp = do(t, router, http.MethodPost, "/session",
		`{"email":"`+email+`","password":"`+password+`"}`, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	loginCookie := jwtCookie(resp)
	require.NotNil(t, loginCookie)
	assert.Equal(t, 1, userQueries, "login must load the user only once")

	// Login and signup issue equivalent cookies that resolve through /me.
	resp = do(t, router, http.MethodGet, "/me", "", loginCookie)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&me))
	require.NotNil(t, me.User)
	assert.Equal(t, email, me.User.Email)

	// A wrong password is rejected with a 422 validation error.
	resp = do(t, router, http.MethodPost, "/session",
		`{"email":"`+email+`","password":"wrong"}`, nil)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var wrongPassword api.ValidationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&wrongPassword))

	// Unknown users produce the same response and cannot be enumerated.
	resp = do(t, router, http.MethodPost, "/session",
		`{"email":"missing-auth-flow@example.com","password":"wrong"}`, nil)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var missingUser api.ValidationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&missingUser))
	assert.Equal(t, wrongPassword, missingUser)

	// Accounts without a password digest are also rejected as invalid
	// credentials rather than causing a panic or exposing account state.
	_, err := db.User.Create().SetEmail("passwordless-auth-flow@example.com").Save(t.Context())
	require.NoError(t, err)
	resp = do(t, router, http.MethodPost, "/session",
		`{"email":"passwordless-auth-flow@example.com","password":"wrong"}`, nil)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	var passwordlessUser api.ValidationError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&passwordlessUser))
	assert.Equal(t, wrongPassword, passwordlessUser)
}

func TestCurrentUserSurvivesEmailChange(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	router := newAuthRouterWithDB(t, db, transactor)
	const oldEmail = "jwt-old-email@example.com"
	const newEmail = "jwt-new-email@example.com"

	resp := do(t, router, http.MethodPost, "/users",
		`{"firstName":"Ada","email":"`+oldEmail+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	cookie := jwtCookie(resp)
	require.NotNil(t, cookie)

	var created struct {
		ID int32 `json:"id"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	_, err := db.User.UpdateOneID(int(created.ID)).SetEmail(newEmail).Save(t.Context())
	require.NoError(t, err)

	resp = do(t, router, http.MethodGet, "/me", "", cookie)
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

	resp := do(t, router, http.MethodPost, "/users",
		`{"firstName":"Ada","email":"logout-auth-flow@example.com","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	jwt := jwtCookie(resp)
	xsrf := xsrfCookieFromResponse(resp)
	require.NotNil(t, jwt)
	require.NotNil(t, xsrf)

	resp = do(t, router, http.MethodDelete, "/session", "", jwt)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp = doWithXSRF(t, router, http.MethodDelete, "/session", "", jwt, "wrong")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp = doWithXSRF(t, router, http.MethodDelete, "/session", "", jwt, xsrf.Value)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Reset writes expired JWT and XSRF cookies so the browser drops both.
	cookie := jwtCookie(resp)
	require.NotNil(t, cookie)
	assert.True(t, cookie.MaxAge < 0 || cookie.Value == "",
		"logout must expire or empty the JWT cookie")
	xsrf = xsrfCookieFromResponse(resp)
	require.NotNil(t, xsrf)
	assert.True(t, xsrf.MaxAge < 0 || xsrf.Value == "",
		"logout must expire or empty the XSRF cookie")
}

func TestContractSecurityProtectsParityRoutesAndKeepsLessonCheckPublic(t *testing.T) {
	router := newAuthRouter(t)

	resp := do(t, router, http.MethodPost, "/leads",
		`{"contactMethod":"phone","contactValue":"+10000000000","ymClientId":null}`, nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp = do(t, router, http.MethodPost, "/lessons/1/check",
		`{"code":"puts 1","versionId":1}`, nil)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAdminAuthorizationUsesCurrentDatabaseValue(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)
	router := newAuthRouterWithDB(t, db, transactor)
	const email = "admin-revocation@example.com"

	resp := do(t, router, http.MethodPost, "/users",
		`{"firstName":"Ada","email":"`+email+`","password":"s3cret-pass"}`, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	jwt := jwtCookie(resp)
	require.NotNil(t, jwt)

	u, err := db.User.Query().Where(user.Email(email)).Only(t.Context())
	require.NoError(t, err)

	resp = do(t, router, http.MethodGet, "/admin/course_categories", "", jwt)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	_, err = db.User.UpdateOneID(u.ID).SetAdmin(true).Save(t.Context())
	require.NoError(t, err)
	resp = do(t, router, http.MethodGet, "/admin/course_categories", "", jwt)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	_, err = db.User.UpdateOneID(u.ID).SetAdmin(false).Save(t.Context())
	require.NoError(t, err)
	resp = do(t, router, http.MethodGet, "/admin/course_categories", "", jwt)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
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

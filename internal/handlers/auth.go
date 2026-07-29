package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	authlogger "github.com/go-pkgz/auth/v2/logger"
	authmiddleware "github.com/go-pkgz/auth/v2/middleware"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/config"
	"hexletbasics/internal/events"
	"hexletbasics/internal/ids"
	"hexletbasics/internal/localization"
)

const (
	authTokenTTL  = 24 * time.Hour
	authCookieTTL = 31 * 24 * time.Hour
	authIssuer    = "hexlet-basics"
	authCookie    = "JWT"
	authProvider  = "password"
	xsrfCookie    = "XSRF-TOKEN"
)

var errInvalidCredentials = errors.New("invalid credentials")

// AuthHandler keeps credential verification and go-pkgz/auth JWT handling
// behind the generated ogen interface. Cookie headers are returned through the
// response models declared in TypeSpec, so the public HTTP seam remains the
// generated contract.
type AuthHandler struct {
	db     *ent.Client
	conv   apiconv.Converter
	jwt    *token.Service
	auth   authmiddleware.Authenticator
	i18n   *localization.Translator
	users  accounts.UserRegistrar
	events events.StandalonePublisher
}

// NewAuthHandler builds the auth implementation used by the ogen handlers.
func NewAuthHandler(
	db *ent.Client,
	cfg *config.Config,
	translator *localization.Translator,
	errorHandler *APIErrorHandler,
	registrar accounts.UserRegistrar,
	eventPublisher events.StandalonePublisher,
) *AuthHandler {
	tokenOpts := token.Opts{
		SecretReader: token.SecretFunc(func(string) (string, error) {
			return cfg.JWTSecret, nil
		}),
		Issuer:         authIssuer,
		TokenDuration:  authTokenTTL,
		CookieDuration: authCookieTTL,
		XSRFIgnoreMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
		},
		SameSite: http.SameSiteLaxMode,
	}
	jwtService := token.NewService(tokenOpts)

	return &AuthHandler{
		db:   db,
		conv: &apiconv.ConverterImpl{},
		jwt:  jwtService,
		auth: authmiddleware.Authenticator{
			L:          authlogger.NoOp,
			JWTService: jwtService,
			Providers: []provider.Service{{
				Provider: provider.DirectHandler{ProviderName: authProvider},
			}},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, status int, err error) {
				errorHandler.Write(r.Context(), w, r, withHTTPStatus(status, err))
			},
		},
		i18n:   translator,
		users:  registrar,
		events: eventPublisher,
	}
}

// Trace applies go-pkgz/auth's optional authentication middleware.
func (h *AuthHandler) Trace(next http.Handler) http.Handler {
	return h.auth.Trace(next)
}

// Auth applies go-pkgz/auth's required authentication and XSRF middleware.
func (h *AuthHandler) Auth(next http.Handler) http.Handler {
	return h.auth.Auth(next)
}

// Admin applies go-pkgz/auth's required authentication and admin authorization.
func (h *AuthHandler) Admin(next http.Handler) http.Handler {
	return h.auth.AdminOnly(next)
}

// CreateSession authenticates once and reuses the loaded user for both the JWT
// claims and the response, keeping password login free of internal HTTP
// round-trips and duplicate database queries.
func (h *AuthHandler) CreateSession(ctx context.Context, req *api.SessionInput) (api.CreateSessionRes, error) {
	u, err := h.authenticate(ctx, req.Email, req.Password)
	if errors.Is(err, errInvalidCredentials) {
		return validationError("password", h.i18n.Text(ctx, localization.WrongCredentials)), nil
	}
	if err != nil {
		return nil, err
	}
	if err := h.events.PublishStandalone(ctx, events.UserSignedIn{
		UserID:          u.ID,
		OccurrenceCount: -1,
		Email:           u.Email,
		Locale:          h.i18n.Locale(ctx),
		OccurredAt:      time.Now().UTC(),
	}); err != nil {
		return nil, err
	}

	cookies, err := h.issueCookies(u)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: cookies, Response: h.conv.ToUser(u)}, nil
}

// CreateUser creates an account and issues a go-pkgz/auth JWT cookie.
func (h *AuthHandler) CreateUser(ctx context.Context, req *api.SignUpInput) (api.CreateUserRes, error) {
	var firstName *string
	if value, ok := req.FirstName.Get(); ok {
		firstName = &value
	}

	u, err := h.users.Register(ctx, accounts.Registration{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: firstName,
		Locale:    h.i18n.Locale(ctx),
	})
	if err != nil {
		if errors.Is(err, accounts.ErrPasswordProcessing) {
			return validationError("password", h.i18n.Text(ctx, localization.PasswordProcessingFailed)), nil
		}
		if ent.IsConstraintError(err) {
			return validationError("email", h.i18n.Text(ctx, localization.EmailTaken)), nil
		}
		return nil, err
	}

	cookies, err := h.issueCookies(u)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: cookies, Response: h.conv.ToUser(u)}, nil
}

// DeleteSession returns go-pkgz/auth's expired session cookies.
func (h *AuthHandler) DeleteSession(context.Context) (*api.DeleteSessionNoContent, error) {
	rec := httptest.NewRecorder()
	h.jwt.Reset(rec)
	cookies, err := responseCookies(rec, authCookie, xsrfCookie)
	if err != nil {
		return nil, err
	}
	return &api.DeleteSessionNoContent{SetCookie: cookies}, nil
}

// GetCurrentUser resolves the user populated by go-pkgz/auth's Trace middleware.
// Missing, invalid, or stale sessions are anonymous.
func (h *AuthHandler) GetCurrentUser(ctx context.Context, params api.GetCurrentUserParams) (*api.CurrentUser, error) {
	_ = params
	authReq := httptest.NewRequest(http.MethodGet, "/me", nil).WithContext(ctx)
	authUser, err := token.GetUserInfo(authReq)
	if err != nil {
		return anonymousCurrentUser(), nil
	}

	userID, err := strconv.Atoi(authUser.ID)
	if err != nil {
		return anonymousCurrentUser(), nil
	}

	u, err := h.db.User.Get(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return anonymousCurrentUser(), nil
		}
		return nil, err
	}
	return &api.CurrentUser{User: api.NewNilUser(h.conv.ToUser(u))}, nil
}

func (h *AuthHandler) authenticate(ctx context.Context, email, password string) (*ent.User, error) {
	u, err := h.db.User.Query().Where(user.Email(email)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, errInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if u.PasswordDigest == nil {
		return nil, errInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordDigest), []byte(password)); err != nil {
		return nil, errInvalidCredentials
	}
	return u, nil
}

func (h *AuthHandler) issueCookies(u *ent.User) ([]string, error) {
	rec := httptest.NewRecorder()
	authUser := &token.User{ID: strconv.Itoa(u.ID), Name: *u.Email}
	authUser.SetAdmin(u.Admin != nil && *u.Admin)
	_, err := h.jwt.Set(rec, token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       ids.New(),
			Issuer:   authIssuer,
			Audience: jwt.ClaimStrings{authIssuer},
		},
		User:         authUser,
		AuthProvider: &token.AuthProvider{Name: authProvider},
	})
	if err != nil {
		return nil, err
	}
	return responseCookies(rec, authCookie, xsrfCookie)
}

func responseCookies(rec *httptest.ResponseRecorder, names ...string) ([]string, error) {
	found := make(map[string]string, len(names))
	for _, cookie := range rec.Result().Cookies() {
		for _, name := range names {
			if cookie.Name == name {
				found[name] = cookie.String()
			}
		}
	}

	result := make([]string, 0, len(names))
	for _, name := range names {
		value, ok := found[name]
		if !ok {
			return nil, fmt.Errorf("%s cookie missing from auth response", name)
		}
		result = append(result, value)
	}
	return result, nil
}

func validationError(field, message string) *api.ValidationError {
	return &api.ValidationError{Errors: api.ValidationErrorErrors{field: {message}}}
}

func anonymousCurrentUser() *api.CurrentUser {
	var user api.NilUser
	user.SetToNull()
	return &api.CurrentUser{User: user}
}

// The generated api.Handler seam is implemented by Server; auth remains a
// cohesive module behind four deliberately small forwarding methods.
func (s *Server) CreateSession(ctx context.Context, req *api.SessionInput) (api.CreateSessionRes, error) {
	return s.auth.CreateSession(ctx, req)
}

func (s *Server) CreateUser(ctx context.Context, req *api.SignUpInput) (api.CreateUserRes, error) {
	return s.auth.CreateUser(ctx, req)
}

func (s *Server) DeleteSession(ctx context.Context) (*api.DeleteSessionNoContent, error) {
	return s.auth.DeleteSession(ctx)
}

func (s *Server) GetCurrentUser(ctx context.Context, params api.GetCurrentUserParams) (*api.CurrentUser, error) {
	return s.auth.GetCurrentUser(ctx, params)
}

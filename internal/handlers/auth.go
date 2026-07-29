package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/config"
	"hexletbasics/internal/localization"
)

const (
	authTokenTTL  = 24 * time.Hour
	authCookieTTL = 31 * 24 * time.Hour
	authIssuer    = "hexlet-basics"
	authCookie    = "JWT"
)

var errInvalidCredentials = errors.New("invalid credentials")

// AuthHandler keeps credential verification and go-pkgz/auth JWT handling
// behind the generated ogen interface. Cookie headers are returned through the
// response models declared in TypeSpec, so the public HTTP seam remains the
// generated contract.
type AuthHandler struct {
	db   *ent.Client
	conv apiconv.Converter
	jwt  *token.Service
	i18n *localization.Translator
}

// NewAuthHandler builds the auth implementation used by the ogen handlers.
func NewAuthHandler(db *ent.Client, cfg *config.Config, translator *localization.Translator) *AuthHandler {
	tokenOpts := token.Opts{
		SecretReader: token.SecretFunc(func(string) (string, error) {
			return cfg.JWTSecret, nil
		}),
		Issuer:         authIssuer,
		TokenDuration:  authTokenTTL,
		CookieDuration: authCookieTTL,
		DisableXSRF:    true,
		SameSite:       http.SameSiteLaxMode,
	}

	return &AuthHandler{
		db:   db,
		conv: &apiconv.ConverterImpl{},
		jwt:  token.NewService(tokenOpts),
		i18n: translator,
	}
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

	cookie, err := h.issueCookie(u)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: cookie, Response: h.conv.ToUser(u)}, nil
}

// CreateUser creates an account and issues a go-pkgz/auth JWT cookie.
func (h *AuthHandler) CreateUser(ctx context.Context, req *api.SignUpInput) (api.CreateUserRes, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return validationError("password", h.i18n.Text(ctx, localization.PasswordProcessingFailed)), nil
	}
	var firstName *string
	if value, ok := req.FirstName.Get(); ok {
		firstName = &value
	}

	u, err := h.db.User.Create().
		SetEmail(req.Email).
		SetPasswordDigest(string(hash)).
		SetNillableFirstName(firstName).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return validationError("email", h.i18n.Text(ctx, localization.EmailTaken)), nil
		}
		return nil, err
	}

	cookie, err := h.issueCookie(u)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: cookie, Response: h.conv.ToUser(u)}, nil
}

// DeleteSession returns the expired JWT cookie through ogen's typed response.
func (h *AuthHandler) DeleteSession(context.Context) (*api.DeleteSessionNoContent, error) {
	rec := httptest.NewRecorder()
	h.jwt.Reset(rec)
	cookie, err := responseCookie(rec, authCookie)
	if err != nil {
		return nil, err
	}
	return &api.DeleteSessionNoContent{SetCookie: cookie}, nil
}

// GetCurrentUser resolves the optional Cookie request header modeled by the
// contract. Missing, invalid, or stale sessions are anonymous.
func (h *AuthHandler) GetCurrentUser(ctx context.Context, params api.GetCurrentUserParams) (*api.CurrentUser, error) {
	authReq := httptest.NewRequest(http.MethodGet, "/me", nil).WithContext(ctx)
	if cookie, ok := params.Cookie.Get(); ok {
		authReq.Header.Set("Cookie", cookie)
	}
	claims, _, err := h.jwt.Get(authReq)
	if err != nil || claims.User == nil {
		return anonymousCurrentUser(), nil
	}

	userID, err := strconv.Atoi(claims.User.ID)
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

func (h *AuthHandler) issueCookie(u *ent.User) (string, error) {
	rec := httptest.NewRecorder()
	_, err := h.jwt.Set(rec, token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       strconv.Itoa(u.ID),
			Issuer:   authIssuer,
			Audience: jwt.ClaimStrings{authIssuer},
		},
		User: &token.User{ID: strconv.Itoa(u.ID), Name: *u.Email},
	})
	if err != nil {
		return "", err
	}
	return responseCookie(rec, authCookie)
}

func responseCookie(rec *httptest.ResponseRecorder, name string) (string, error) {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == name {
			return cookie.String(), nil
		}
	}
	return "", fmt.Errorf("%s cookie missing from auth response", name)
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

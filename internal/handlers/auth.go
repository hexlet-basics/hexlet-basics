package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	pkgzauth "github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/config"
)

const (
	authTokenTTL  = 24 * time.Hour
	authCookieTTL = 31 * 24 * time.Hour
	authIssuer    = "hexlet-basics"
	authCookie    = "JWT"
)

// AuthHandler adapts go-pkgz/auth's HTTP-oriented JWT/direct-provider model to
// the generated ogen interface. The adapter captures the provider's cookie and
// returns it as the Set-Cookie response header modeled in TypeSpec; the public
// HTTP seam therefore remains the generated contract.
type AuthHandler struct {
	db         *ent.Client
	conv       apiconv.Converter
	jwt        *token.Service
	directAuth http.Handler
}

// NewAuthHandler builds the auth implementation used by the ogen handlers.
func NewAuthHandler(db *ent.Client, cfg *config.Config) *AuthHandler {
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

	authService := pkgzauth.NewService(pkgzauth.Opts{
		SecretReader:   tokenOpts.SecretReader,
		Issuer:         tokenOpts.Issuer,
		TokenDuration:  tokenOpts.TokenDuration,
		CookieDuration: tokenOpts.CookieDuration,
		DisableXSRF:    tokenOpts.DisableXSRF,
		SameSiteCookie: tokenOpts.SameSite,
	})
	authService.AddDirectProvider("direct", provider.CredCheckerFunc(func(email, password string) (bool, error) {
		u, err := db.User.Query().Where(user.Email(email)).Only(context.Background())
		if ent.IsNotFound(err) || err == nil && u.PasswordDigest == nil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return bcrypt.CompareHashAndPassword([]byte(*u.PasswordDigest), []byte(password)) == nil, nil
	}))
	directAuth, _ := authService.Handlers()

	return &AuthHandler{
		db:         db,
		conv:       &apiconv.ConverterImpl{},
		jwt:        token.NewService(tokenOpts),
		directAuth: directAuth,
	}
}

// CreateSession delegates credential verification and JWT issuance to the
// go-pkgz/auth direct provider, then maps its response to the contract model.
func (h *AuthHandler) CreateSession(ctx context.Context, req *api.SessionInput) (api.CreateSessionRes, error) {
	body := fmt.Sprintf(`{"user":%q,"passwd":%q,"aud":%q}`, req.Email, req.Password, authIssuer)
	authReq := httptest.NewRequest(http.MethodPost, "/auth/direct/login", bytes.NewBufferString(body)).
		WithContext(ctx)
	authReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.directAuth.ServeHTTP(rec, authReq)
	if rec.Code != http.StatusOK {
		return validationError("password", "Wrong email or password"), nil
	}

	u, err := h.db.User.Query().Where(user.Email(req.Email)).Only(ctx)
	if err != nil {
		return nil, err
	}
	cookie, err := responseCookie(rec, authCookie)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: cookie, Response: h.conv.ToUser(u)}, nil
}

// CreateUser creates an account and issues the same go-pkgz/auth JWT cookie
// used by the direct provider.
func (h *AuthHandler) CreateUser(ctx context.Context, req *api.SignUpInput) (api.CreateUserRes, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return validationError("password", "Could not process the password"), nil
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
		return validationError("email", "This email is already taken"), nil
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

	u, err := h.db.User.Query().Where(user.Email(claims.User.Name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return anonymousCurrentUser(), nil
		}
		return nil, err
	}
	return &api.CurrentUser{User: api.NewNilUser(h.conv.ToUser(u))}, nil
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

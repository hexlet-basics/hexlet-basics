package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/config"
)

// AuthHandler implements the email+password auth endpoints (ADR-0003). These
// live outside the generated ogen layer because they must set/clear the JWT
// cookie on the raw http.ResponseWriter, which ogen's typed handlers can't reach
// (the same reason attachments/webhooks are hand-mounted — see router.go). The
// contract paths (`POST /session`, `DELETE /session`, `POST /users`, `GET /me`)
// are mounted ahead of the ogen catch-all, so the generated TS client hits these
// implementations instead of the UnimplementedHandler 501s.
//
// The JWT is minted and read with go-pkgz/auth's token.Service so we reuse its
// cookie handling; the cookie name defaults to `JWT`, matching the frontend
// (lib/auth.ts) and go-pkgz's own default.
type AuthHandler struct {
	db   *ent.Client
	conv apiconv.Converter
	jwt  *token.Service
}

const (
	// Access-token lifetime and cookie lifetime. Short-ish token TTL is the
	// accepted trade for JWT's lack of server-side revocation (ADR-0003).
	authTokenTTL  = 24 * time.Hour
	authCookieTTL = 31 * 24 * time.Hour
	authIssuer    = "hexlet-basics"
)

// NewAuthHandler builds the JWT service from the configured secret.
func NewAuthHandler(db *ent.Client, cfg *config.Config) *AuthHandler {
	jwtService := token.NewService(token.Opts{
		SecretReader: token.SecretFunc(func(string) (string, error) {
			return cfg.JWTSecret, nil
		}),
		Issuer:         authIssuer,
		TokenDuration:  authTokenTTL,
		CookieDuration: authCookieTTL,
		// The SPA sends the cookie via `credentials: include` and does not carry
		// an XSRF header, so XSRF double-submit is disabled here. SameSite=Lax
		// blocks the cross-site POSTs that XSRF protects against for now.
		DisableXSRF: true,
		SameSite:    http.SameSiteLaxMode,
	})
	return &AuthHandler{db: db, conv: &apiconv.ConverterImpl{}, jwt: jwtService}
}

// Login verifies email+password (legacy bcrypt digest) and sets the JWT cookie.
// Contract: createSession — `POST /session` → User | ValidationError.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeValidationError(w, "base", "Invalid request body")
		return
	}

	u, err := h.db.User.Query().Where(user.Email(in.Email)).Only(r.Context())
	// A missing user and a wrong password return the same generic error so the
	// response can't be used to enumerate registered emails.
	if err != nil || u.PasswordDigest == nil ||
		bcrypt.CompareHashAndPassword([]byte(*u.PasswordDigest), []byte(in.Password)) != nil {
		writeValidationError(w, "password", "Wrong email or password")
		return
	}

	if err := h.setSession(w, u); err != nil {
		writeValidationError(w, "base", "Could not start a session")
		return
	}
	apiUser := h.conv.ToUser(u)
	writeJSON(w, http.StatusOK, &apiUser)
}

// Register creates an account (bcrypt-hashing the password) and signs the user
// in. Contract: createUser — `POST /users` → 201 User | ValidationError. Email
// uniqueness is enforced by the baseline DB index; a duplicate surfaces as a
// 422 keyed on `email`, matching the admin CRUD's constraint handling.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		FirstName *string `json:"firstName"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeValidationError(w, "base", "Invalid request body")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		writeValidationError(w, "password", "Could not process the password")
		return
	}

	u, err := h.db.User.Create().
		SetEmail(in.Email).
		SetPasswordDigest(string(hash)).
		SetNillableFirstName(in.FirstName).
		Save(r.Context())
	if err != nil {
		// Most likely the unique-email constraint; surface it on the email field.
		writeValidationError(w, "email", "This email is already taken")
		return
	}

	if err := h.setSession(w, u); err != nil {
		writeValidationError(w, "base", "Could not start a session")
		return
	}
	apiUser := h.conv.ToUser(u)
	writeJSON(w, http.StatusCreated, &apiUser)
}

// Logout clears the JWT cookie. Contract: deleteSession — `DELETE /session` →
// 204.
func (h *AuthHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	h.jwt.Reset(w)
	w.WriteHeader(http.StatusNoContent)
}

// Me resolves the current user from the JWT cookie for SSR. Contract:
// getCurrentUser — `GET /me` → CurrentUser ({ user: User | null }). An absent or
// invalid cookie is not an error: it simply yields `{ user: null }`.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, _, err := h.jwt.Get(r)
	if err != nil || claims.User == nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": nil})
		return
	}
	id, err := strconv.Atoi(claims.User.ID)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": nil})
		return
	}
	u, err := h.db.User.Get(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": nil})
		return
	}
	apiUser := h.conv.ToUser(u)
	writeJSON(w, http.StatusOK, map[string]any{"user": &apiUser})
}

// setSession mints a JWT carrying the user id and writes it as the auth cookie.
func (h *AuthHandler) setSession(w http.ResponseWriter, u *ent.User) error {
	id := strconv.Itoa(u.ID)
	claims := token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			Issuer:    authIssuer,
			Audience:  jwt.ClaimStrings{authIssuer},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(authTokenTTL)),
		},
		User: &token.User{ID: id, Name: userDisplayName(u)},
	}
	_, err := h.jwt.Set(w, claims)
	return err
}

// userDisplayName mirrors the api User `name` fallback (first+last, else email).
func userDisplayName(u *ent.User) string {
	name := ""
	if u.FirstName != nil {
		name = *u.FirstName
	}
	if u.LastName != nil {
		if name != "" {
			name += " "
		}
		name += *u.LastName
	}
	if name == "" && u.Email != nil {
		name = *u.Email
	}
	return name
}

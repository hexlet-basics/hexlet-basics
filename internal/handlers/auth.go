package handlers

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	"hexletbasics/internal/progress"
)

const (
	authTokenTTL  = 24 * time.Hour
	authCookieTTL = 31 * 24 * time.Hour
	authIssuer    = "hexlet-basics"
	authCookie    = "JWT"
	authProvider  = "password"
	xsrfCookie    = "XSRF-TOKEN"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
	errUnauthenticated    = errors.New("request is not authenticated")
	errAdminRequired      = errors.New("administrator access is required")
)

type authenticatedUserContext struct {
	user *ent.User
	jti  string
}

type authenticatedUserContextKey struct{}

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
	errors *APIErrorHandler
	// guests carries the visitor progress a new session inherits, and the codec
	// that reads and clears its cookie.
	guests   progress.Tracker
	guestJar *progress.GuestCodec
	secure   bool
}

// NewAuthHandler builds the auth implementation used by the ogen handlers.
func NewAuthHandler(
	db *ent.Client,
	cfg *config.Config,
	translator *localization.Translator,
	errorHandler *APIErrorHandler,
	registrar accounts.UserRegistrar,
	eventPublisher events.StandalonePublisher,
	tracker progress.Tracker,
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
		i18n:     translator,
		users:    registrar,
		events:   eventPublisher,
		errors:   errorHandler,
		guests:   tracker,
		guestJar: progress.NewGuestCodec(cfg.JWTSecret),
		secure:   strings.HasPrefix(cfg.PublicURL, "https://"),
	}
}

// Trace applies go-pkgz/auth's optional authentication middleware.
func (h *AuthHandler) Trace(next http.Handler) http.Handler {
	return h.auth.Trace(next)
}

// Identify attaches the signed-in user to the request context when a valid
// session cookie is present, and does nothing otherwise.
//
// Contract-declared authentication (ADR-0011) covers operations that REQUIRE a
// session: ogen runs their security handler and rejects the request. The public
// reads are different — they must answer a visitor and a learner alike, with
// the learner getting their progress — and ogen never invokes security for an
// operation that declares none. This middleware supplies that optional
// identity, and never fails a request: an absent or invalid cookie simply
// leaves the context anonymous.
func (h *AuthHandler) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookie)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx, err := h.loadAuthenticatedUser(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// guestProgressContextKey carries the decoded guest cookie.
type guestProgressContextKey struct{}

// CarryGuestProgress decodes the guest cookie into the request context. Like
// Identify it never fails a request: an absent, forged or tampered cookie
// leaves the visitor with no progress, which is the same thing a first visit
// looks like.
func (h *AuthHandler) CarryGuestProgress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(progress.GuestCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		guest, err := h.guestJar.Decode(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), guestProgressContextKey{}, guest),
		))
	})
}

// GuestProgress returns the visitor progress carried by the request cookie.
func GuestProgress(ctx context.Context) (progress.GuestProgress, bool) {
	guest, ok := ctx.Value(guestProgressContextKey{}).(progress.GuestProgress)
	return guest, ok
}

// GuestCookie renders visitor progress as the signed response cookie carrying
// it back. Handlers that advance a guest ask for the cookie rather than for the
// codec, so the secret and the cookie attributes stay in one place.
func (h *AuthHandler) GuestCookie(guest progress.GuestProgress) (string, error) {
	cookie, err := h.guestJar.Cookie(guest, h.secure)
	if err != nil {
		return "", err
	}
	return cookie.String(), nil
}

// mergeGuestProgress credits the visitor's cookie progress to the account that
// just signed in or signed up, and returns the cookie that clears it.
//
// Deviation from #762 worth stating: the merge runs in its own transaction
// rather than the one that issues the session. Session issue writes nothing
// (the JWT is signed in process) and sign-up owns its transaction inside the
// registrar, so one shared transaction would mean restructuring both. The merge
// itself is atomic, and the cookie is cleared only after it succeeds — a failed
// merge leaves the visitor's progress intact to be retried on the next sign-in,
// which is the property the acceptance criterion is protecting.
func (h *AuthHandler) mergeGuestProgress(ctx context.Context, userID int) ([]string, error) {
	guest, ok := GuestProgress(ctx)
	if !ok || len(guest.Entries) == 0 {
		return nil, nil
	}

	if err := h.guests.MergeGuest(ctx, userID, guest, h.i18n.Locale(ctx)); err != nil {
		return nil, err
	}

	cleared := &http.Cookie{
		Name:     progress.GuestCookieName,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	}
	return []string{cleared.String()}, nil
}

// AuthenticatedUser returns the database user loaded by the generated security
// handler. Protected application handlers use this context seam instead of
// reparsing JWT claims or issuing a second user query.
func AuthenticatedUser(ctx context.Context) (*ent.User, bool) {
	value, ok := ctx.Value(authenticatedUserContextKey{}).(*authenticatedUserContext)
	if !ok || value.user == nil {
		return nil, false
	}
	return value.user, true
}

// HandleUserSession implements ogen's generated user-session security seam.
func (h *AuthHandler) HandleUserSession(
	ctx context.Context,
	_ api.OperationName,
	session api.UserSession,
) (context.Context, error) {
	return h.loadAuthenticatedUser(ctx, session.APIKey)
}

// HandleAdminSession implements ogen's generated admin-session security seam.
// The JWT's historical admin claim is intentionally ignored: authorization
// always reads the current database value so revocation takes effect at once.
func (h *AuthHandler) HandleAdminSession(
	ctx context.Context,
	_ api.OperationName,
	session api.AdminSession,
) (context.Context, error) {
	ctx, err := h.loadAuthenticatedUser(ctx, session.APIKey)
	if err != nil {
		return ctx, err
	}
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return ctx, errUnauthenticated
	}
	if u.Admin == nil || !*u.Admin {
		return ctx, withHTTPStatus(http.StatusForbidden, errAdminRequired)
	}
	return ctx, nil
}

// HandleXsrfToken completes the contract's AND security group. The session
// handler stores the verified JWT id in context before ogen invokes this method.
func (h *AuthHandler) HandleXsrfToken(
	ctx context.Context,
	_ api.OperationName,
	xsrf api.XsrfToken,
) (context.Context, error) {
	authenticated, ok := ctx.Value(authenticatedUserContextKey{}).(*authenticatedUserContext)
	if !ok || authenticated.jti == "" {
		return ctx, errUnauthenticated
	}
	if subtle.ConstantTimeCompare([]byte(authenticated.jti), []byte(xsrf.APIKey)) != 1 {
		return ctx, errUnauthenticated
	}
	return ctx, nil
}

// RequireAdmin protects the temporary multipart adapter that ogen cannot
// generate yet. It delegates to the same HandleAdminSession/HandleXsrfToken
// methods ogen invokes for generated operations, so the admin+XSRF policy is
// defined once; only cookie extraction and error writing are manual here.
func (h *AuthHandler) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookie)
		if err != nil {
			h.errors.Write(r.Context(), w, r, withHTTPStatus(http.StatusUnauthorized, errUnauthenticated))
			return
		}

		ctx, err := h.HandleAdminSession(r.Context(), "", api.AdminSession{APIKey: cookie.Value})
		if err == nil {
			_, err = h.HandleXsrfToken(ctx, "", api.XsrfToken{APIKey: r.Header.Get("X-XSRF-TOKEN")})
		}
		if errors.Is(err, errUnauthenticated) {
			err = withHTTPStatus(http.StatusUnauthorized, err)
		}
		if err != nil {
			h.errors.Write(ctx, w, r, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loadAuthenticatedUser verifies a session token and attaches the user it names.
//
// Verification is done here rather than by reading what go-pkgz's Trace
// middleware left in the context. Trace refuses a token on an unsafe method
// without the double-submit XSRF header, which is right for the operations
// whose contract requires a session — they declare the XsrfToken scheme and
// ogen enforces it (ADR-0011) — but wrong for a PUBLIC unsafe operation:
// checking a solution has to answer a guest as well as a learner, so no client
// sends that header, and piggybacking on Trace would silently serve every
// signed-in learner as a guest. The JWT cookie is SameSite=Lax, so a
// cross-site POST carries no session to forge with.
func (h *AuthHandler) loadAuthenticatedUser(ctx context.Context, rawJWT string) (context.Context, error) {
	claims, err := h.jwt.Parse(rawJWT)
	if err != nil || claims.Handshake != nil || claims.User == nil {
		return ctx, errUnauthenticated
	}

	userID, err := strconv.Atoi(claims.User.ID)
	if err != nil {
		return ctx, errUnauthenticated
	}
	u, err := h.db.User.Get(ctx, userID)
	if ent.IsNotFound(err) {
		return ctx, errUnauthenticated
	}
	if err != nil {
		return ctx, withHTTPStatus(http.StatusInternalServerError, fmt.Errorf("load authenticated user: %w", err))
	}

	return context.WithValue(ctx, authenticatedUserContextKey{}, &authenticatedUserContext{
		user: u,
		jti:  claims.ID,
	}), nil
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
	merged, err := h.mergeGuestProgress(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: append(cookies, merged...), Response: h.conv.ToUser(u)}, nil
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
	merged, err := h.mergeGuestProgress(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	return &api.UserHeaders{SetCookie: append(cookies, merged...), Response: h.conv.ToUser(u)}, nil
}

// DeleteSession returns go-pkgz/auth's expired session cookies.
func (h *AuthHandler) DeleteSession(context.Context) (api.DeleteSessionRes, error) {
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

func (s *Server) DeleteSession(ctx context.Context) (api.DeleteSessionRes, error) {
	return s.auth.DeleteSession(ctx)
}

func (s *Server) GetCurrentUser(ctx context.Context, params api.GetCurrentUserParams) (*api.CurrentUser, error) {
	return s.auth.GetCurrentUser(ctx, params)
}

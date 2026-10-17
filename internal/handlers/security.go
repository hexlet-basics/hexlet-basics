package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/staffroles"
)

var (
	errUnauthenticated = errors.New("request is not authenticated")
	errAdminRequired   = errors.New("administrator access is required")
	errCrossOrigin     = errors.New("cross-origin request refused")
)

type authenticatedUserContextKey struct{}

// AuthenticatedUser returns the database user loaded by the generated security
// handler or the Identify middleware. Protected application handlers use this
// context seam instead of asking the session a second time.
func AuthenticatedUser(ctx context.Context) (*ent.User, bool) {
	u, ok := ctx.Value(authenticatedUserContextKey{}).(*ent.User)
	if !ok || u == nil {
		return nil, false
	}
	return u, true
}

// MustAuthenticatedUser is AuthenticatedUser for an operation the contract
// secures: ogen runs the security handler first and answers 401 itself, so a
// request that reaches the handler has a user. Its absence means the operation
// was declared without security while the handler assumes it, a mistake worth a
// loud failure instead of a branch no request can take.
func MustAuthenticatedUser(ctx context.Context) *ent.User {
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		panic("handlers: operation reached its handler without an authenticated user; declare @useAuth in the contract")
	}
	return u
}

// Security is how a request is recognised: it implements ogen's generated
// security seam (ADR-0011) and the middleware that runs around the generated
// server. Both read the user from the session (ADR-0017); neither parses a
// credential of its own.
type Security struct {
	db       *ent.Client
	sessions *sessions.Manager
	policy   *staffroles.Policy
	i18n     *localization.Translator
	errors   *APIErrorHandler
	guests   *guestCookies
}

// newSecurity builds the security seam over the session manager.
func newSecurity(
	db *ent.Client,
	manager *sessions.Manager,
	policy *staffroles.Policy,
	translator *localization.Translator,
	errorHandler *APIErrorHandler,
	guests *guestCookies,
) *Security {
	return &Security{
		db:       db,
		sessions: manager,
		policy:   policy,
		i18n:     translator,
		errors:   errorHandler,
		guests:   guests,
	}
}

// Pipeline wraps the generated server with what every API request passes
// through, outermost first:
//
//   - cross-origin protection, which refuses a browser's cross-site write before
//     any work is done. It judges by the Sec-Fetch-Site and Origin headers the
//     browser sets itself; requests that carry neither (the SSR server, webhooks,
//     tests) are not browsers and pass;
//   - the session, which loads the visitor's session and writes its cookie back;
//   - the caller's address, the signed-in user and the guest progress cookie.
//
// The harness builds exactly this, so a test passes the same gates a request does.
func (s *Security) Pipeline(apiHandler http.Handler) http.Handler {
	identified := WithClientIP(s.identify(s.guests.carry(apiHandler)))

	protection := http.NewCrossOriginProtection()
	protection.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.errors.Write(r.Context(), w, r, withHTTPStatus(http.StatusForbidden, errCrossOrigin))
	}))
	return protection.Handler(s.sessions.Middleware(identified))
}

// identify attaches the signed-in user to the request context when its session
// names one, and does nothing otherwise.
//
// Contract-declared authentication (ADR-0011) covers operations that REQUIRE a
// session: ogen runs their security handler and rejects the request. The
// public reads are different: they must answer a visitor and a learner alike,
// with the learner getting their progress, and ogen never invokes security for
// an operation that declares none. This middleware supplies that optional
// identity, and never fails a request: no session, or one that no longer names a
// current user, leaves the context anonymous.
func (s *Security) identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := s.attach(r.Context())
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandleUserSession implements ogen's generated user-session security seam. The
// cookie value ogen hands over is ignored: the session middleware has already
// resolved it, and a session answers for itself.
func (s *Security) HandleUserSession(
	ctx context.Context,
	_ api.OperationName,
	_ api.UserSession,
) (context.Context, error) {
	return s.attach(ctx)
}

// HandleAdminSession implements ogen's generated admin-session security seam.
// Authorization always reads the current database row, never anything the
// session remembers, so revoking administrator rights takes effect at once.
func (s *Security) HandleAdminSession(
	ctx context.Context,
	operation api.OperationName,
	_ api.AdminSession,
) (context.Context, error) {
	// An MCP tool call carries a Principal in its context (never in a header or
	// a cookie) and takes the branch with the capability check; every other call
	// is a session. The session cookie value of a tool call's in-process request
	// is a placeholder nothing reads: this branch is taken before the session is
	// consulted, and the same request sent over HTTP has no Principal and finds
	// no session under the placeholder.
	if principal, ok := mcpserver.PrincipalFrom(ctx); ok {
		return s.authorizeMCP(ctx, operation, principal)
	}
	ctx, err := s.attach(ctx)
	if err != nil {
		return ctx, err
	}
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return ctx, errUnauthenticated
	}
	err = s.policy.Allow(ctx, u, s.i18n.Locale(ctx), operation)
	if errors.Is(err, staffroles.ErrForbidden) {
		s.recordRefusal(ctx, operation)
		return ctx, withHTTPStatus(http.StatusForbidden, errAdminRequired)
	}
	return ctx, err
}

// recordRefusal leaves an Audit Record with outcome forbidden for a call the
// permission check turned away (ADR-0019, #845). It runs here because a refusal
// never reaches the audit middleware, and it is a single insert on the pooled
// client with no ambient transaction (none exists yet at the security stage), so
// it commits on its own. The record names who, where and which operation and
// nothing else: the request body is not decoded yet, so no value the caller was
// not allowed to read or change can end up in it. A record that cannot be
// written is logged and does not turn a refusal into a different answer.
func (s *Security) recordRefusal(ctx context.Context, operation string) {
	ctx = context.WithoutCancel(ctx)
	entry := auditEntry(ctx, operation, s.i18n, audit.OutcomeForbidden)
	if err := audit.Write(ctx, s.db, entry, nil); err != nil {
		s.errors.logger.ErrorContext(ctx, "recording a refused admin call", "operation", operation, "err", err)
	}
}

// attach puts the session's user in the context.
//
// The session is resolved by the middleware rather than by ogen's security
// handler so that it also serves a PUBLIC unsafe operation: checking a solution
// has to answer a guest as well as a learner, and ogen never runs security for
// an operation that declares none.
func (s *Security) attach(ctx context.Context) (context.Context, error) {
	u, err := s.sessions.User(ctx)
	if errors.Is(err, sessions.ErrUnauthenticated) {
		return ctx, errUnauthenticated
	}
	if err != nil {
		return ctx, withHTTPStatus(http.StatusInternalServerError, fmt.Errorf("load authenticated user: %w", err))
	}
	return context.WithValue(ctx, authenticatedUserContextKey{}, u), nil
}

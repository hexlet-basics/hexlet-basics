package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/permissions"
	"hexletbasics/internal/staffroles"
)

// NewAdminDispatcher is what an MCP admin tool call is served by (ADR-0018
// addendum): the generated api.Server, wrapped as the router wraps it for the
// website except for what does not apply to a call that has no browser. It
// resolves the locale from the request's Accept-Language, which the tool sets
// from its locale argument, and records the caller's address and user agent.
// It does not run the session pipeline (Security.Pipeline), which loads a
// session from the cookie: a tool call has none.
//
// It is the only thing the MCP server is given to reach an operation, so every
// tool call passes the generated security handler (authorizeMCP) and the audit
// middleware installed on apiServer.
func NewAdminDispatcher(apiServer http.Handler, i18n *localization.Translator) http.Handler {
	return i18n.Middleware(WithClientIP(apiServer))
}

// authorizeMCP is the MCP branch of HandleAdminSession: the same checkpoint an
// HTTP admin call passes, for a call whose subject is a Connected App's token
// instead of a session. It judges, from the database on every call:
//
//  1. the admin Scope the user granted the Connected App;
//  2. the operation's Permission and the locale, by Policy.Allow, which is the
//     very call the cookie branch makes;
//  3. the capability Permission matching the operation's kind (mcp.read or
//     mcp.write), by Policy.HasCapability, which the administrator flag does
//     not satisfy (ADR-0017).
//
// A failure of any is the refusal HTTP gives: an Audit Record with outcome
// forbidden (channel mcp, since the tool put its Source in the context) and 403,
// before the body is decoded or a handler runs. Taking a Permission away is
// seen on the next call, because nothing is cached.
func (s *Security) authorizeMCP(ctx context.Context, operation api.OperationName, p mcpserver.Principal) (context.Context, error) {
	u, err := accounts.ActiveUserByID(ctx, s.db, p.UserID)
	if ent.IsNotFound(err) {
		return ctx, errUnauthenticated
	}
	if err != nil {
		return ctx, withHTTPStatus(http.StatusInternalServerError, fmt.Errorf("load authenticated user: %w", err))
	}
	ctx = context.WithValue(ctx, authenticatedUserContextKey{}, u)

	// The locale the generated server resolved must be the one the tool named;
	// anything else (a dispatcher wired without the locale middleware, a header
	// that did not survive) would judge the call in a locale nobody asked for.
	locale := s.i18n.Locale(ctx)
	if locale != p.Locale || !slices.Contains(p.Scopes, oauthserver.ScopeAdmin) {
		s.recordRefusal(ctx, operation)
		return ctx, withHTTPStatus(http.StatusForbidden, errAdminRequired)
	}
	err = s.policy.Allow(ctx, u, locale, operation)
	if err == nil {
		err = s.policy.HasCapability(ctx, u, permissions.MCPCapability(operation))
	}
	if errors.Is(err, staffroles.ErrForbidden) {
		s.recordRefusal(ctx, operation)
		return ctx, withHTTPStatus(http.StatusForbidden, errAdminRequired)
	}
	return ctx, err
}

package mcpserver

import "context"

// Principal is who an admin tool call acts for, as the generated server's
// security handler sees it (ADR-0021 addendum, #854). The tool dispatcher puts
// one in the context of the in-process request it sends through the generated
// server; handlers.Security.HandleAdminSession, the one check the HTTP
// operations also pass, reads it instead of a session and decides from the
// database: the Scope, the operation's Permission and the matching mcp.*
// Permission.
//
// It is a value of the context, never of a header or a cookie, so an HTTP
// client cannot present one: the only code that sets it is in this package, and
// it sets it from a verified OAuth token.
type Principal struct {
	// UserID is the token's subject.
	UserID int
	// Scopes are what the user granted the Connected App.
	Scopes []string
	// Locale is the tool's locale argument, already validated. The check must
	// judge the call in this locale, so the security handler compares it with the
	// locale the request resolved to and refuses on any difference.
	Locale string
}

type principalKey struct{}

// WithPrincipal marks ctx as the context of an MCP tool call.
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

// PrincipalFrom is the Principal WithPrincipal set, if any.
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

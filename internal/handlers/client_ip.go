package handlers

import (
	"context"
	"net/http"

	"github.com/go-pkgz/rest/realip"
)

type clientIPKey struct{}

// requestInfo is what an Audit Record needs from the HTTP request and a
// security handler or MCP tool cannot read from a context otherwise.
type requestInfo struct {
	userAgent string
	requestID string
}

type requestInfoKey struct{}

// WithClientIP puts the submitting client's address into the request context,
// together with the user agent and the caller's request id, which the audit
// trail records for refusals too (a refusal happens in the security handler,
// where only a context is at hand).
// ogen hands operations only a context, and lead attribution needs the IP
// legacy took from ahoy's visit (ADR-0015). Behind the ingress RemoteAddr is
// the proxy, so realip reads the forwarding headers first and skips private
// hops; the value is attribution data for amoCRM, never a security decision.
func WithClientIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestInfoKey{}, requestInfo{
			userAgent: r.UserAgent(),
			requestID: r.Header.Get(requestIDHeader),
		})
		if ip, err := realip.Get(r); err == nil && ip != "" {
			ctx = context.WithValue(ctx, clientIPKey{}, ip)
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ClientIP returns the address WithClientIP resolved, if any.
func ClientIP(ctx context.Context) (string, bool) {
	ip, ok := ctx.Value(clientIPKey{}).(string)
	return ip, ok
}

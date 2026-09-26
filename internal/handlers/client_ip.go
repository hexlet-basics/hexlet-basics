package handlers

import (
	"context"
	"net/http"

	"github.com/go-pkgz/rest/realip"
)

type clientIPKey struct{}

// WithClientIP puts the submitting client's address into the request context.
// ogen hands operations only a context, and lead attribution needs the IP
// legacy took from ahoy's visit (ADR-0015). Behind the ingress RemoteAddr is
// the proxy, so realip reads the forwarding headers first and skips private
// hops; the value is attribution data for amoCRM, never a security decision.
func WithClientIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip, err := realip.Get(r); err == nil && ip != "" {
			r = r.WithContext(context.WithValue(r.Context(), clientIPKey{}, ip))
		}
		next.ServeHTTP(w, r)
	})
}

// ClientIP returns the address WithClientIP resolved, if any.
func ClientIP(ctx context.Context) (string, bool) {
	ip, ok := ctx.Value(clientIPKey{}).(string)
	return ip, ok
}

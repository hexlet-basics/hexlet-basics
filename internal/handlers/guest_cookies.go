package handlers

import (
	"context"
	"net/http"

	"hexletbasics/internal/progress"
)

// guestProgressContextKey carries the decoded guest cookie.
type guestProgressContextKey struct{}

// guestCookies is the visitor-progress cookie as the handlers meet it: read from
// a request, written to a response, cleared once an account has taken the
// progress over. The codec owns the cookie's attributes; this type only decides
// who is asked for it, so the secret and the Secure flag are stated once.
type guestCookies struct {
	codec  *progress.GuestCodec
	secure bool
}

func newGuestCookies(secret string, secure bool) *guestCookies {
	return &guestCookies{codec: progress.NewGuestCodec(secret), secure: secure}
}

// carry decodes the guest cookie into the request context. It never fails a
// request: an absent, forged or tampered cookie leaves the visitor with no
// progress, which is the same thing a first visit looks like.
func (g *guestCookies) carry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(progress.GuestCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		guest, err := g.codec.Decode(cookie.Value)
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

// set renders visitor progress as the signed response cookie carrying it back.
func (g *guestCookies) set(guest progress.GuestProgress) (string, error) {
	cookie, err := g.codec.Cookie(guest, g.secure)
	if err != nil {
		return "", err
	}
	return cookie.String(), nil
}

// cleared is the response cookie that ends the visitor's progress cookie.
func (g *guestCookies) cleared() string {
	return g.codec.ClearedCookie(g.secure).String()
}

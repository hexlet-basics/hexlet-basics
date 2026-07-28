package handlers

import "net/http"

// NewRouter composes the full HTTP surface: the multipart/blob routes that ogen
// cannot generate (ADR-0005, internal/apigen/ogen.yml) plus the generated
// api.Server for everything else.
//
// The generated server is mounted at the catch-all `/`; the explicit method+path
// patterns for the attachment routes are strictly more specific, so Go's
// ServeMux routes them first and only unmatched requests fall through to ogen.
//
// Auth note: there is no auth middleware yet (none of the /admin/* operations
// are protected today). When it lands, wrap the returned http.Handler so it
// covers BOTH the custom routes and the generated server — do not wrap only the
// api.Server, or these uploads would sit outside the protected perimeter.
func NewRouter(apiHandler http.Handler, att *AttachmentHandler, gh *GitHubWebhookHandler, auth *AuthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/attachments", att.Upload)
	mux.HandleFunc("GET /storage/{key}", att.Download)
	// GitHub webhook: verifies its own HMAC signature, so it is safe outside any
	// future admin-auth middleware — but it must stay a build TRIGGER only.
	mux.HandleFunc("POST /webhooks/github", gh.Handle)
	// Auth endpoints (ADR-0003): mounted here rather than in ogen because they
	// set/clear the JWT cookie on the raw ResponseWriter. Their paths mirror the
	// contract, so the generated client reaches these instead of the 501 stubs.
	mux.HandleFunc("POST /session", auth.Login)
	mux.HandleFunc("DELETE /session", auth.Logout)
	mux.HandleFunc("POST /users", auth.Register)
	mux.HandleFunc("GET /me", auth.Me)
	mux.Handle("/", apiHandler)
	return mux
}

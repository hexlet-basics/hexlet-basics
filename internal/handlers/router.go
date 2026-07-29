package handlers

import "net/http"

// NewRouter composes the full HTTP surface: the temporary multipart upload
// adapter required by ogen's missing requestBody.encoding support, the blob
// read route (ADR-0005, internal/apigen/ogen.yml), and the generated api.Server
// for everything else.
//
// The generated server is mounted at the catch-all `/`; the explicit method+path
// patterns for the attachment routes are strictly more specific, so Go's
// ServeMux routes them first and only unmatched requests fall through to ogen.
//
// The inner mux owns transport dispatch. The outer mux applies go-pkgz/auth's
// required middleware to authenticated route families and its optional Trace
// middleware to the public surface. Keeping the generated API and custom
// attachment route behind the same branch prevents uploads from bypassing auth.
func NewRouter(
	apiHandler http.Handler,
	att *AttachmentHandler,
	gh *GitHubWebhookHandler,
	auth *AuthHandler,
) http.Handler {
	transport := http.NewServeMux()
	transport.HandleFunc("POST /admin/attachments", att.Upload)
	transport.HandleFunc("GET /storage/{key}", att.Download)
	// GitHub webhook: verifies its own HMAC signature, so it is safe outside any
	// auth middleware — but it must stay a build TRIGGER only.
	transport.HandleFunc("POST /webhooks/github", gh.Handle)
	transport.Handle("/", apiHandler)

	required := auth.Auth(transport)
	admin := auth.Admin(transport)
	optional := auth.Trace(transport)

	router := http.NewServeMux()
	router.Handle("/admin/", admin)
	router.Handle("/account/", required)
	router.Handle("/ai/lessons/", required)
	router.Handle("/my", required)
	router.Handle("DELETE /session", required)
	router.Handle("/", optional)
	return router
}

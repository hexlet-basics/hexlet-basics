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
// TypeSpec security requirements and ogen's generated SecurityHandler protect
// generated operations. Trace wraps the complete transport only to preserve
// go-pkgz/auth's optional identity extraction and sliding JWT refresh.
//
// The multipart upload remains an exact-route exception while ogen cannot
// generate it; RequireAdmin reuses the same database-backed auth module without
// reintroducing URL-family policy.
func NewRouter(
	apiHandler http.Handler,
	att *AttachmentHandler,
	gh *GitHubWebhookHandler,
	auth *AuthHandler,
) http.Handler {
	transport := http.NewServeMux()
	transport.Handle("POST /admin/attachments", auth.RequireAdmin(http.HandlerFunc(att.Upload)))
	transport.HandleFunc("GET /storage/{key}", att.Download)
	// GitHub webhook: verifies its own HMAC signature, so it is safe outside any
	// auth middleware — but it must stay a build TRIGGER only.
	transport.HandleFunc("POST /webhooks/github", gh.Handle)
	// Identify wraps only the generated operations: it attaches the signed-in
	// user when a session cookie is present, so a public read can answer a
	// learner with their progress and a visitor without it.
	transport.Handle("/", auth.Identify(apiHandler))
	return auth.Trace(transport)
}

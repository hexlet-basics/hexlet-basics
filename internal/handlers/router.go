package handlers

import (
	"net/http"

	"hexletbasics/internal/assetstore"
)

// UploadBodyLimit bounds the whole multipart request. The generated decoder's
// ParseMultipartForm limit is only a memory threshold (larger parts spill to
// disk), so the total size is capped here; the extra megabyte leaves room for
// the multipart envelope while assetstore enforces the exact file-byte limit.
const UploadBodyLimit = assetstore.MaxUploadBytes + 1<<20

// NewRouter composes the full HTTP surface: the blob read route (ADR-0005),
// the GitHub webhook, and the generated api.Server for every contract
// operation. Everything lives under `/api` on the site's own host (ADR-0015);
// the ingress sends the rest of the site to the SSR server.
//
// The generated server is mounted at the catch-all `/` (the contract's paths
// already carry the prefix, and ogen answers anything else with the localized
// problem 404); the explicit method+path patterns are strictly more specific,
// so Go's ServeMux routes them first and only unmatched requests fall through
// to ogen. The upload route is the generated server too, behind a body-size
// cap.
//
// TypeSpec security requirements and ogen's generated SecurityHandler protect
// generated operations. Trace wraps the complete transport only to preserve
// go-pkgz/auth's optional identity extraction and sliding JWT refresh.
func NewRouter(
	apiHandler http.Handler,
	att *AttachmentHandler,
	gh *GitHubWebhookHandler,
	auth *AuthHandler,
) http.Handler {
	// Identify wraps only the generated operations: it attaches the signed-in
	// user when a session cookie is present, so a public read can answer a
	// learner with their progress and a visitor without it.
	// WithClientIP records the caller's address for lead attribution.
	generated := WithClientIP(auth.Identify(auth.CarryGuestProgress(apiHandler)))

	transport := http.NewServeMux()
	transport.Handle("POST /api/admin/attachments", http.MaxBytesHandler(generated, UploadBodyLimit))
	transport.HandleFunc("GET /api/storage/{key}", att.Download)
	// GitHub webhook: verifies its own HMAC signature, so it is safe outside any
	// auth middleware — but it must stay a build TRIGGER only. GitHub still
	// delivers to the legacy `/webhooks/github`; the ingress maps it here.
	transport.HandleFunc("POST /api/webhooks/github", gh.Handle)
	transport.Handle("/", generated)
	return auth.Trace(transport)
}

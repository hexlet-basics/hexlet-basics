package handlers

import (
	"errors"
	"net/http"
	"time"

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
	// The assistant answers as a stream: each chunk the generated encoder copies
	// has to leave the process as it is written, not when the answer ends.
	transport.Handle("POST /api/ai/lessons/{lessonId}/messages", StreamResponses(generated))
	// GitHub webhook: verifies its own HMAC signature, so it is safe outside any
	// auth middleware — but it must stay a build TRIGGER only. GitHub still
	// delivers to the legacy `/webhooks/github`; the ingress maps it here.
	transport.HandleFunc("POST /api/webhooks/github", gh.Handle)
	transport.Handle("/", generated)
	return auth.Trace(transport)
}

// StreamResponses serves a route whose body is produced over time (the
// assistant's answer). Three things stand between a written chunk and the
// learner, and each is lifted here: net/http buffers the body (every Write is
// flushed), nginx buffers proxied responses (X-Accel-Buffering: no), and the
// server's write timeout would cut an answer that takes longer than it (the
// deadline is cleared for this request only).
//
// The controller reaches the connection through every middleware wrapper that
// implements Unwrap; a writer that cannot flush or move its deadline (a test
// recorder) is served unbuffered-in-name-only rather than refused.
func StreamResponses(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := http.NewResponseController(w)
		_ = controller.SetWriteDeadline(time.Time{})
		w.Header().Set("X-Accel-Buffering", "no")
		next.ServeHTTP(&flushingWriter{ResponseWriter: w, controller: controller}, r)
	})
}

// flushingWriter flushes after every Write. A failed flush means the client is
// gone, and is returned so the copy feeding it stops.
type flushingWriter struct {
	http.ResponseWriter
	controller *http.ResponseController
}

func (w *flushingWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	if err != nil {
		return n, err
	}
	if err := w.controller.Flush(); err != nil && !errors.Is(err, http.ErrNotSupported) {
		return n, err
	}
	return n, nil
}

// Unwrap lets response controllers further in reach the connection.
func (w *flushingWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

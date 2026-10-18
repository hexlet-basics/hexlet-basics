package handlers

import (
	"errors"
	"net/http"
	"time"

	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
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
// `generated` is the generated server inside Security.Pipeline: the session, the
// signed-in user and cross-origin protection reach every operation, including the
// two routes that wrap it below. The blob read and the webhook sit outside it:
// the blob is public, and the webhook proves itself by its HMAC signature and
// carries no browser session.
//
// TypeSpec security requirements and ogen's generated SecurityHandler protect
// generated operations; the pipeline only supplies the identity they read.
//
// oauth is the OAuth authorization server (ADR-0022), nil while the feature is
// switched off. It is the one thing mounted outside `/api`: RFC 8414 puts its
// metadata at the root well-known path, and its endpoints sit under `/oauth/`.
// It stays outside the session pipeline: its clients are native or browser
// applications that carry no session cookie, and the pipeline's cross-origin
// protection would refuse a web client's legitimate cross-site token request.
// The consent screen, which does need the session, is a generated operation.
// WithClientIP resolves the caller's address for its registration rate limit.
// mcp, the MCP endpoint (ADR-0021), is mounted outside `/api` the same way, and is
// nil while its own flag is off.
func NewRouter(
	generated http.Handler,
	att *AttachmentHandler,
	gh *GitHubWebhookHandler,
	oauth http.Handler,
	mcp http.Handler,
) http.Handler {
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

	if oauth != nil {
		for _, pattern := range oauthserver.Patterns() {
			transport.Handle(pattern, WithClientIP(oauth))
		}
	}
	// The MCP endpoint and its resource metadata (ADR-0021), nil while switched
	// off. It authenticates a bearer token itself, so it sits beside the
	// generated server and not behind the session pipeline: the website's session
	// cookie is never read on it, and the generated operations accept no bearer token.
	if mcp != nil {
		for _, pattern := range mcpserver.Patterns() {
			transport.Handle(pattern, mcp)
		}
	}

	transport.Handle("/", generated)

	return transport
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
//
// A status written after the body has begun is the generated server reporting
// a failure mid-answer (the model dropped, storing the answer failed). The 200
// is already on the wire, and appending the error document would hand the
// client a short answer with JSON glued to it that reads as complete; aborting
// the handler cuts the connection instead, which the client sees as the broken
// stream it is. The error has been reported by then — the error handler logs
// before it writes.
type flushingWriter struct {
	http.ResponseWriter

	controller *http.ResponseController
	started    bool
}

func (w *flushingWriter) WriteHeader(status int) {
	if w.started {
		panic(http.ErrAbortHandler)
	}

	w.ResponseWriter.WriteHeader(status)
}

func (w *flushingWriter) Write(p []byte) (int, error) {
	w.started = true

	n, err := w.ResponseWriter.Write(p)
	if err != nil {
		return n, errs.Wrapf(err, "write response")
	}

	if err := w.controller.Flush(); err != nil && !errors.Is(err, http.ErrNotSupported) {
		return n, errs.Wrapf(err, "flush response")
	}

	return n, nil
}

// Unwrap lets response controllers further in reach the connection.
func (w *flushingWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

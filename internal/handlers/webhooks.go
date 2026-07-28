package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
)

// maxWebhookBytes bounds a webhook body. GitHub deliveries are well under this;
// the cap stops an oversized body from being buffered whole (also the ceiling for
// the HMAC read, which must see the exact bytes GitHub signed).
const maxWebhookBytes = 5 << 20 // 5 MiB

// coursePrefix is the fixed prefix on every exercises repo name; stripping it
// yields the course slug (the same convention the git fetcher clones by).
const coursePrefix = "exercises-"

// GitHubWebhookHandler receives GitHub `workflow_run` webhooks and triggers a
// course-version build when a course's exercises repo finishes CI successfully on
// its default branch. It lives outside the ogen router (like the attachment
// routes): a signed webhook is not part of the JSON CRUD contract, and it needs
// raw-body access for signature verification that the generated codec would hide.
//
// Building off CI-success (not a raw push) mirrors the legacy invariant that only
// content which passed the repo's pipeline is ever served.
type GitHubWebhookHandler struct {
	db      *ent.Client
	starter VersionBuildStarter
	// secret is the shared HMAC secret GitHub signs deliveries with. An empty
	// secret disables the endpoint (fail closed) — an unauthenticated build
	// trigger must never be reachable.
	secret string
}

// NewGitHubWebhookHandler wires the webhook to its dependencies.
func NewGitHubWebhookHandler(db *ent.Client, starter VersionBuildStarter, secret string) *GitHubWebhookHandler {
	return &GitHubWebhookHandler{db: db, starter: starter, secret: secret}
}

// workflowRunPayload is the slice of the `workflow_run` event this handler acts
// on. GitHub sends far more; only these fields decide whether to build.
type workflowRunPayload struct {
	Action      string `json:"action"`
	WorkflowRun struct {
		Conclusion string `json:"conclusion"`
		HeadBranch string `json:"head_branch"`
	} `json:"workflow_run"`
	Repository struct {
		Name          string `json:"name"`
		DefaultBranch string `json:"default_branch"`
	} `json:"repository"`
}

// Handle processes one delivery. It returns:
//   - 404 if the endpoint is disabled (no secret configured),
//   - 401 on a missing/invalid signature,
//   - 204 for events it deliberately ignores (wrong event, non-success, non-default
//     branch, unrecognized repo, unknown course) — a webhook receiver must ACK
//     these so GitHub does not retry a delivery there is simply nothing to do about,
//   - 202 when a build was enqueued.
func (h *GitHubWebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if h.secret == "" {
		http.NotFound(w, r)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxWebhookBytes))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	if !h.validSignature(r.Header.Get("X-Hub-Signature-256"), body) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	// Only workflow_run events can trigger a build; ACK everything else.
	if r.Header.Get("X-GitHub-Event") != "workflow_run" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var payload workflowRunPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	slug, ok := h.buildTargetSlug(payload)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	built, err := h.triggerBuild(r.Context(), slug)
	if err != nil {
		http.Error(w, "failed to trigger build", http.StatusInternalServerError)
		return
	}
	if !built {
		// No course matches this repo — nothing to build, but a valid delivery.
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// validSignature verifies GitHub's `sha256=<hex>` HMAC over the exact body bytes,
// in constant time.
func (h *GitHubWebhookHandler) validSignature(header string, body []byte) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	want, err := hex.DecodeString(strings.TrimPrefix(header, prefix))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	return hmac.Equal(want, mac.Sum(nil))
}

// buildTargetSlug decides whether a payload should trigger a build and, if so,
// returns the course slug. A build fires only for a completed, successful run on
// the repo's default branch of an `exercises-<slug>` repo — the same gate legacy
// relied on (only CI-built content is served).
func (h *GitHubWebhookHandler) buildTargetSlug(p workflowRunPayload) (string, bool) {
	if p.Action != "completed" || p.WorkflowRun.Conclusion != "success" {
		return "", false
	}
	if p.WorkflowRun.HeadBranch != p.Repository.DefaultBranch {
		return "", false
	}
	if !strings.HasPrefix(p.Repository.Name, coursePrefix) {
		return "", false
	}
	slug := strings.TrimPrefix(p.Repository.Name, coursePrefix)
	if slug == "" {
		return "", false
	}
	return slug, true
}

// triggerBuild creates a `created` version for the course with the given slug and
// enqueues the loader job — the same effect as the admin createVersion action.
// Returns false (no error) when no course matches the slug.
func (h *GitHubWebhookHandler) triggerBuild(ctx context.Context, slug string) (bool, error) {
	c, err := h.db.Course.Query().Where(course.Slug(slug)).Only(ctx)
	if ent.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if _, err := h.starter.Start(ctx, c.ID); err != nil {
		return false, err
	}
	return true, nil
}

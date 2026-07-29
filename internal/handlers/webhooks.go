package handlers

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v89/github"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/localization"
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
	i18n   *localization.Translator
}

// NewGitHubWebhookHandler wires the webhook to its dependencies.
func NewGitHubWebhookHandler(db *ent.Client, starter VersionBuildStarter, secret string, translator *localization.Translator) *GitHubWebhookHandler {
	return &GitHubWebhookHandler{db: db, starter: starter, secret: secret, i18n: translator}
}

// Handle processes one delivery. It returns:
//   - 404 if the endpoint is disabled (no secret configured),
//   - 401 on a missing/invalid signature,
//   - 204 for events it deliberately ignores (wrong event, non-success, non-default
//     branch, unrecognized repo, unknown course) — a webhook receiver must ACK
//     these so GitHub does not retry a delivery there is simply nothing to do about,
//   - 202 when a build was enqueued.
func (h *GitHubWebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h.secret == "" {
		http.Error(w, h.i18n.StatusText(ctx, http.StatusNotFound), http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxWebhookBytes))
	if err != nil {
		http.Error(w, h.i18n.Text(ctx, localization.ReadWebhookBodyFailed), http.StatusBadRequest)
		return
	}

	if err := github.ValidateSignature(r.Header.Get(github.SHA256SignatureHeader), body, []byte(h.secret)); err != nil {
		http.Error(w, h.i18n.Text(ctx, localization.InvalidWebhookSignature), http.StatusUnauthorized)
		return
	}

	// Only workflow_run events can trigger a build; ACK everything else.
	eventType := github.WebHookType(r)
	if eventType != "workflow_run" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsed, err := github.ParseWebHook(eventType, body)
	if err != nil {
		http.Error(w, h.i18n.Text(ctx, localization.InvalidWebhookPayload), http.StatusBadRequest)
		return
	}
	payload, ok := parsed.(*github.WorkflowRunEvent)
	if !ok {
		http.Error(w, h.i18n.Text(ctx, localization.InvalidWebhookPayload), http.StatusBadRequest)
		return
	}

	slug, ok := h.buildTargetSlug(payload)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	built, err := h.triggerBuild(ctx, slug)
	if err != nil {
		http.Error(w, h.i18n.Text(ctx, localization.TriggerBuildFailed), http.StatusInternalServerError)
		return
	}
	if !built {
		// No course matches this repo — nothing to build, but a valid delivery.
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// buildTargetSlug decides whether a payload should trigger a build and, if so,
// returns the course slug. A build fires only for a completed, successful run on
// the repo's default branch of an `exercises-<slug>` repo — the same gate legacy
// relied on (only CI-built content is served).
func (h *GitHubWebhookHandler) buildTargetSlug(event *github.WorkflowRunEvent) (string, bool) {
	run := event.GetWorkflowRun()
	repo := event.GetRepo()
	if event.GetAction() != "completed" || run.GetConclusion() != "success" {
		return "", false
	}
	if run.GetHeadBranch() != repo.GetDefaultBranch() {
		return "", false
	}
	name := repo.GetName()
	if !strings.HasPrefix(name, coursePrefix) {
		return "", false
	}
	slug := strings.TrimPrefix(name, coursePrefix)
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

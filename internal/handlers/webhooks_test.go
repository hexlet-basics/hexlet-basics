package handlers_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hexletbasics/ent"
	"hexletbasics/internal/jobs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/handlers"
	"hexletbasics/internal/testsupport"
)

const webhookSecret = "test-secret"

// newWebhookHandler wires the real build starter over the test transaction, so
// a test reads the version and the loader job the webhook left behind.
func newWebhookHandler(t *testing.T, secret string) (*handlers.GitHubWebhookHandler, *ent.Client, *testsupport.Services) {
	t.Helper()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)

	return handlers.NewGitHubWebhookHandler(services.Starter, secret, testsupport.NewTranslator(t)), db, services
}

func rawWebhookRequest(event, body, signature string) *http.Request {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/webhooks/github", strings.NewReader(body))
	req.Header.Set("X-Github-Event", event)

	if signature != "" {
		req.Header.Set("X-Hub-Signature-256", signature)
	}

	return req
}

// signedRequest builds a POST /api/webhooks/github request for the given event and
// JSON payload, signed with webhookSecret exactly as GitHub would.
func signedRequest(t *testing.T, event string, payload any) *http.Request {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return rawWebhookRequest(event, string(body), sig)
}

// workflowRunSuccess is a minimal successful-CI payload for the given repo.
func workflowRunSuccess(repo string) map[string]any {
	return map[string]any{
		"action":       "completed",
		"workflow_run": map[string]any{"conclusion": "success", "head_branch": "main"},
		"repository":   map[string]any{"name": repo, "default_branch": "main"},
	}
}

func TestGitHubWebhookTriggersBuild(t *testing.T) {
	t.Parallel()
	h, db, services := newWebhookHandler(t, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby")))

	assert.Equal(t, http.StatusAccepted, rec.Code)

	// The loader was queued; look up exactly that version (robust against fixture
	// versions the ruby course already has).
	queued := testsupport.RequireEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})

	version := db.CourseVersion.GetX(context.Background(), queued.VersionID)
	require.NotNil(t, version.State)
	assert.Equal(t, "created", *version.State)
	assert.Equal(t, courseRubyIDA(), version.CourseID)
}

func TestGitHubWebhookRejectsBadSignature(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	req := signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby"))
	req.Header.Set("X-Hub-Signature-256", "sha256=deadbeef")

	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookRejectsMissingOrMalformedSignature(t *testing.T) {
	t.Parallel()

	for _, signature := range []string{"", "not-a-signature", "md5=deadbeef"} {
		t.Run(signature, func(t *testing.T) {
			h, _, services := newWebhookHandler(t, webhookSecret)

			rec := httptest.NewRecorder()
			h.Handle(rec, rawWebhookRequest("workflow_run", `{}`, signature))

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
				VersionID: 0,
			})
		})
	}
}

func TestGitHubWebhookValidatesOfficialSignatureVector(t *testing.T) {
	t.Parallel()

	const (
		secret    = "It's a Secret to Everybody"
		body      = "Hello, World!"
		signature = "sha256=757107ea0eb2509fc211221cce984b8a37570b6d7586c22c46f4379c8b043e17"
	)

	t.Run("valid", func(t *testing.T) {
		h, _, services := newWebhookHandler(t, secret)

		rec := httptest.NewRecorder()
		h.Handle(rec, rawWebhookRequest("ping", body, signature))

		assert.Equal(t, http.StatusNoContent, rec.Code)
		testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
			VersionID: 0,
		})
	})

	t.Run("tampered payload", func(t *testing.T) {
		h, _, services := newWebhookHandler(t, secret)

		rec := httptest.NewRecorder()
		h.Handle(rec, rawWebhookRequest("ping", body+"!", signature))

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
			VersionID: 0,
		})
	})
}

func TestGitHubWebhookDisabledWithoutSecret(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, "")

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby")))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookIgnoresOtherEvent(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "push", map[string]any{"not": "a workflow run"}))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookRejectsMalformedWorkflowRun(t *testing.T) {
	t.Parallel()

	body := `{`
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write([]byte(body))
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	h, _, services := newWebhookHandler(t, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, rawWebhookRequest("workflow_run", body, signature))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookIgnoresIncompleteWorkflowRun(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", map[string]any{"action": "completed"}))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookIgnoresFailedRun(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	payload := workflowRunSuccess("exercises-ruby")
	as[map[string]any](t, payload["workflow_run"])["conclusion"] = "failure"

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", payload))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookIgnoresNonDefaultBranch(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	payload := workflowRunSuccess("exercises-ruby")
	as[map[string]any](t, payload["workflow_run"])["head_branch"] = "feature"

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", payload))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookIgnoresUnrecognizedRepository(t *testing.T) {
	t.Parallel()

	for _, repo := range []string{"ruby", "exercises-"} {
		t.Run(repo, func(t *testing.T) {
			h, _, services := newWebhookHandler(t, webhookSecret)

			rec := httptest.NewRecorder()
			h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess(repo)))

			assert.Equal(t, http.StatusNoContent, rec.Code)
			testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
				VersionID: 0,
			})
		})
	}
}

func TestGitHubWebhookIgnoresUnknownCourse(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-nonexistent-course")))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

func TestGitHubWebhookRejectsOversizedBody(t *testing.T) {
	t.Parallel()
	h, _, services := newWebhookHandler(t, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, rawWebhookRequest("workflow_run", strings.Repeat("x", (5<<20)+1), "sha256=deadbeef"))

	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	testsupport.RequireNotEnqueued(t, services, jobs.ExerciseLoaderArgs{
		VersionID: 0,
	})
}

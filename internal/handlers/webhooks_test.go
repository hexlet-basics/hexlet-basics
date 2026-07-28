package handlers_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

const webhookSecret = "test-secret"

// signedRequest builds a POST /webhooks/github request for the given event and
// JSON payload, signed with webhookSecret exactly as GitHub would.
func signedRequest(t *testing.T, event string, payload any) *http.Request {
	t.Helper()
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", strings.NewReader(string(body)))
	req.Header.Set("X-GitHub-Event", event)
	req.Header.Set("X-Hub-Signature-256", sig)
	return req
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
	db := testsupport.NewClient(t)
	ctx := context.Background()
	enq := &testsupport.RecordingEnqueuer{}
	h := handlers.NewGitHubWebhookHandler(db, enq, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby")))

	assert.Equal(t, http.StatusAccepted, rec.Code)

	// The loader was queued; look up exactly that version (robust against fixture
	// versions the ruby course already has).
	require.Len(t, enq.Inserted, 1)
	args := enq.Inserted[0].(jobs.ExerciseLoaderArgs)

	version := db.CourseVersion.GetX(ctx, args.VersionID)
	require.NotNil(t, version.State)
	assert.Equal(t, "created", *version.State)
	assert.Equal(t, courseRubyIDA, version.LanguageID)
}

func TestGitHubWebhookRejectsBadSignature(t *testing.T) {
	db := testsupport.NewClient(t)
	enq := &testsupport.RecordingEnqueuer{}
	h := handlers.NewGitHubWebhookHandler(db, enq, webhookSecret)

	req := signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby"))
	req.Header.Set("X-Hub-Signature-256", "sha256=deadbeef")

	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Empty(t, enq.Inserted)
}

func TestGitHubWebhookDisabledWithoutSecret(t *testing.T) {
	db := testsupport.NewClient(t)
	enq := &testsupport.RecordingEnqueuer{}
	h := handlers.NewGitHubWebhookHandler(db, enq, "")

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby")))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Empty(t, enq.Inserted)
}

func TestGitHubWebhookIgnoresFailedRun(t *testing.T) {
	db := testsupport.NewClient(t)
	enq := &testsupport.RecordingEnqueuer{}
	h := handlers.NewGitHubWebhookHandler(db, enq, webhookSecret)

	payload := workflowRunSuccess("exercises-ruby")
	payload["workflow_run"].(map[string]any)["conclusion"] = "failure"

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", payload))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, enq.Inserted)
}

func TestGitHubWebhookIgnoresUnknownCourse(t *testing.T) {
	db := testsupport.NewClient(t)
	enq := &testsupport.RecordingEnqueuer{}
	h := handlers.NewGitHubWebhookHandler(db, enq, webhookSecret)

	rec := httptest.NewRecorder()
	h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-nonexistent-course")))

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, enq.Inserted)
}

package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/courseversion"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/testsupport"
)

// A delivery GitHub signed correctly is answered 500 when the build cannot be
// started, so GitHub shows the delivery as failed and it can be redelivered,
// instead of the build being dropped behind a 2xx.
func TestGitHubWebhookAnswersAServerErrorWhenTheBuildCannotStart(t *testing.T) {
	broken := map[string]string{
		"the course cannot be read":        fmt.Sprintf(`ALTER TABLE %q RENAME TO %q`, course.Table, "courses_gone"),
		"the version cannot be written":    fmt.Sprintf(`ALTER TABLE %q RENAME TO %q`, courseversion.Table, "versions_gone"),
		"the course lookup column is gone": fmt.Sprintf(`ALTER TABLE %q RENAME COLUMN %q TO %q`, course.Table, course.FieldSlug, "slug_gone"),
	}
	for name, statement := range broken {
		t.Run(name, func(t *testing.T) {
			db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
			services := testsupport.NewServices(t, tx, transactor)
			translator := testsupport.NewTranslator(t)
			h := handlers.NewGitHubWebhookHandler(db, services.Starter, webhookSecret, translator)
			_, err := tx.ExecContext(t.Context(), statement)
			require.NoError(t, err, fmt.Sprintf("breaking storage with %s", statement))

			rec := httptest.NewRecorder()
			h.Handle(rec, signedRequest(t, "workflow_run", workflowRunSuccess("exercises-ruby")))

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), translator.Text(t.Context(), localization.TriggerBuildFailed))
		})
	}
}

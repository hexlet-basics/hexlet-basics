package audit_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/testsupport"
	"hexletbasics/internal/testsupport/testdb"
)

func TestMain(m *testing.M) {
	os.Exit(testdb.Run(m))
}

// A field marked Sensitive() in the ent schema keeps only its name in the
// record: password material and tokens must never reach the trail, in either
// direction. A neighbouring ordinary field of the same update keeps both values.
func TestASecretFieldIsRecordedByNameOnly(t *testing.T) {
	client := testsupport.NewClient(t)
	ctx, recorder := audit.WithRecorder(t.Context())
	user := client.User.Query().FirstX(t.Context())

	require.NoError(t, client.User.UpdateOneID(user.ID).
		SetPasswordDigest("$2a$10$new-secret-hash").
		SetFirstName("Renamed").
		Exec(ctx))

	changes := recorder.Changes()
	require.Len(t, changes, 1)
	assert.Equal(t, "User", changes[0].Entity)

	byName := map[string]audit.FieldChange{}
	for _, f := range changes[0].Fields {
		byName[f.Name] = f
	}

	secret, ok := byName["password_digest"]
	require.True(t, ok, "a changed secret is still a change")
	assert.True(t, secret.Secret)
	assert.Empty(t, secret.Old)
	assert.Empty(t, secret.New)

	first := byName["first_name"]
	assert.False(t, first.Secret)
	assert.JSONEq(t, `"Renamed"`, string(first.New))
}

func TestTheSecretFieldsAreGeneratedFromTheSchemaMarker(t *testing.T) {
	assert.Contains(t, ent.AuditSecretFields["User"], "password_digest")
	assert.Contains(t, ent.AuditSecretFields["User"], "confirmation_token")
	assert.NotContains(t, ent.AuditSecretFields["User"], "email")
	assert.Empty(t, ent.AuditSecretFields["BlogPost"])
}

func TestWithoutARecorderTheHookDoesNothing(t *testing.T) {
	client := testsupport.NewClient(t)
	user := client.User.Query().FirstX(t.Context())

	require.NoError(t, client.User.UpdateOneID(user.ID).SetFirstName("Quiet").Exec(t.Context()))
}

func TestSourceDefaultsToTheWebsite(t *testing.T) {
	assert.Equal(t, audit.ChannelWeb, audit.SourceFrom(t.Context()).Channel)

	app := 7
	got := audit.SourceFrom(audit.WithSource(t.Context(), audit.Source{
		Channel:        audit.ChannelMCP,
		ConnectedAppID: &app,
		TokenID:        nil,
	}))
	assert.Equal(t, audit.ChannelMCP, got.Channel)
	assert.Equal(t, 7, *got.ConnectedAppID)
}

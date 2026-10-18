package audit_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// A delete keeps what it destroyed: the old value of each field, loaded before
// the row went, with a secret field by name only.
func TestADeleteRecordsTheOldValuesOfTheRow(t *testing.T) {
	client := testsupport.NewClient(t)
	ctx, recorder := audit.WithRecorder(t.Context())
	// A learner with a password and nothing attached, so that the delete goes through.
	user := client.User.GetX(t.Context(), fixture.ID(fixture.UserWithPassword))

	require.NoError(t, client.User.DeleteOneID(user.ID).Exec(ctx))

	changes := recorder.Changes()
	require.Len(t, changes, 1)
	assert.Equal(t, "delete", changes[0].Op)
	assert.Equal(t, strconv.Itoa(user.ID), changes[0].ID)

	byName := map[string]audit.FieldChange{}
	for _, f := range changes[0].Fields {
		byName[f.Name] = f
	}

	assert.JSONEq(t, `"with-password@example.com"`, string(byName["email"].Old))
	assert.JSONEq(t, `"Password"`, string(byName["first_name"].Old))
	assert.Empty(t, byName["email"].New)

	digest, ok := byName["password_digest"]
	require.True(t, ok, "the secret existed, so its name is kept")
	assert.True(t, digest.Secret)
	assert.Empty(t, digest.Old)
	assert.NotContains(t, byName, "created_at", "ent's own timestamps are not recorded")
}

func TestADeleteOfAMissingRowRecordsNothingAndFailsAsBefore(t *testing.T) {
	client := testsupport.NewClient(t)
	ctx, recorder := audit.WithRecorder(t.Context())

	err := client.User.DeleteOneID(987654321).Exec(ctx)

	require.Error(t, err)
	assert.True(t, ent.IsNotFound(err))
	assert.Empty(t, recorder.Changes())
}

package eventhandlers_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/auditrecord"
	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// The consumer writes the record with the time of the read, not of the
// delivery, and for an MCP read carries the Connected App and token.
func TestStaffReadPersonalDataIsRecordedAtTheTimeOfTheRead(t *testing.T) {
	client := testsupport.NewClient(t)
	handler := eventhandlers.StaffReadPersonalData(client)
	app, token, actor := 3, 4, 5
	readAt := time.Now().Add(-2 * time.Hour).Truncate(time.Second)

	require.NoError(t, handler.Handle(t.Context(), &events.StaffReadPersonalData{
		ActorUserID: &actor, Channel: "mcp", ConnectedAppID: &app, TokenID: &token,
		Operation: "AdminGetLead", EntityType: "Lead", EntityIDs: []string{"11"}, Count: 1,
		RequestID: "req-1", OccurredAt: readAt,
	}))

	row := client.AuditRecord.Query().Where(auditrecord.Operation("AdminGetLead")).OnlyX(t.Context())
	assert.True(t, readAt.Equal(row.OccurredAt), "occurred_at is the read's time")
	assert.Equal(t, "mcp", row.Channel)
	assert.Equal(t, app, *row.ConnectedAppID)
	assert.Equal(t, token, *row.TokenID)
	assert.Equal(t, "Lead", *row.EntityType)
	assert.Equal(t, "11", *row.EntityID)
	assert.Equal(t, "succeeded", row.Outcome)
	assert.Equal(t, "staff_read_personal_data_to_audit_v1", handler.HandlerName())
}

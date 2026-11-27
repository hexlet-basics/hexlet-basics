package mcpserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/audit"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/testsupport"
)

// The generated server is built with every middleware in one api.WithMiddleware
// call (a second call replaces the first), so a dispatched call must show the
// effect of each: the audit middleware's record for a write, and the fact of the
// personal-data read middleware for a read of users.
func TestDispatchedCallsPassBothMiddlewares(t *testing.T) {
	mcpserver.AllowForTest(t, "adminListUsers", "adminUpdateBanner")
	o := testsupport.NewOAuth(t)
	o.Grant(editor(), "users.read", "banners.update", "mcp.read", "mcp.write", "mcp.manage")
	session, token := admin(t, o, bob())

	read := call(t, session, "adminListUsers", map[string]any{"locale": "en", "perPage": 2})
	require.False(t, read.IsError, "%+v", read.Content)

	queued := testsupport.EnqueuedJobs[jobs.RecordPersonalDataReadArgs](t, o.Services)
	require.Len(t, queued, 1)
	fact := queued[0].Event
	assert.Equal(t, "mcp", fact.Channel)
	assert.Equal(t, token.ConnectedAppID, *fact.ConnectedAppID)
	assert.Equal(t, token.ID, *fact.TokenID)
	testsupport.WorkEnqueued(t, o.Services, jobs.RecordPersonalDataReadWorker(o.DB))
	records := o.Records("AdminListUsers")
	require.Len(t, records, 1)
	assert.Equal(t, "mcp", records[0].Channel)

	write := call(t, session, "adminUpdateBanner", map[string]any{
		"locale": "en", "id": firstBannerID(t, o), "body": bannerInput("Both middlewares"),
	})
	require.False(t, write.IsError, "%+v", write.Content)

	written := o.Records("AdminUpdateBanner")
	require.Len(t, written, 1)
	assert.Equal(t, audit.OutcomeSucceeded, audit.Outcome(written[0].Outcome))
	assert.Equal(t, "mcp", written[0].Channel)
}

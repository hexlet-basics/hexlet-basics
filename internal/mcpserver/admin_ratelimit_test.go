package mcpserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// Admin tools draw on the same limiter as the catalogue ones: a write spends the
// write budget, a read the personal one, and a limited call never reaches the
// checkpoint (it leaves no Audit Record).
func TestAdminToolsAreRateLimitedByClass(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{
		User: mcpserver.Rates{Personal: 1, Write: 1},
	}))
	o.Grant(editor, "banners.update", "mcp.read", "mcp.write")
	session := connect(t, o, bob, oauthserver.ScopeAdmin)
	read := map[string]any{"locale": "en"}
	write := map[string]any{"locale": "en", "id": firstBannerID(t, o), "body": bannerInput("Once")}

	require.False(t, callTool(t, session, "adminListBanners", read).IsError)
	rejected(t, callTool(t, session, "adminListBanners", read))

	require.False(t, callTool(t, session, "adminUpdateBanner", write).IsError)
	rejected(t, callTool(t, session, "adminUpdateBanner", write))

	assert.Len(t, o.Records("AdminUpdateBanner"), 1, "the refused call never reached the audit middleware")
}

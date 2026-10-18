package mcpserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport"
)

func TestTokenOfARemovedAccountIsNotSignedIn(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor(), "mcp.read")
	session, _ := admin(t, o, bob())
	o.DB.User.UpdateOneID(bob()).SetState("removed").ExecX(t.Context())

	result := call(t, session, "adminListBanners", map[string]any{"locale": "en"})

	require.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "HTTP 401")
}

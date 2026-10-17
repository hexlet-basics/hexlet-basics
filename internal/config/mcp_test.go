package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadLeavesMCPOffByDefault(t *testing.T) {
	setRequiredSecrets(t)

	cfg, err := Load()

	require.NoError(t, err)
	require.False(t, cfg.MCP.Enabled)
}

func TestLoadAcceptsMCPBesideOAuth(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)
	t.Setenv("MCP_ENABLED", "true")

	cfg, err := Load()

	require.NoError(t, err)
	require.True(t, cfg.MCP.Enabled)
}

func TestLoadRejectsMCPWithoutOAuth(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("MCP_ENABLED", "true")

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "MCP_ENABLED requires OAUTH_ENABLED")
}

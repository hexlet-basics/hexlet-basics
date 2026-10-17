package config

import (
	"testing"
	"time"

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

func enableMCP(t *testing.T) {
	t.Helper()
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)
	t.Setenv("MCP_ENABLED", "true")
}

func TestLoadGivesMCPLimitsTheirDefaults(t *testing.T) {
	enableMCP(t)

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, 50, cfg.MCP.MaxPageSize)
	require.Equal(t, 102400, cfg.MCP.MaxResponseBytes)
	require.Equal(t, 30*time.Second, cfg.MCP.ToolTimeout)
	require.Equal(t, time.Minute, cfg.MCP.RateWindow)
	require.Greater(t, cfg.MCP.RateUserCatalogue, cfg.MCP.RateUserPersonal, "personal reads are stricter than catalogue reads")
	require.Greater(t, cfg.MCP.RateUserPersonal, cfg.MCP.RateUserWrite, "writes are the strictest")
	require.Greater(t, cfg.MCP.RateClientCatalogue, cfg.MCP.RateClientPersonal)
	require.Greater(t, cfg.MCP.RateClientPersonal, cfg.MCP.RateClientWrite)
}

func TestLoadRejectsBadMCPLimitsOnlyWhenMCPIsOn(t *testing.T) {
	for env, value := range map[string]string{
		"MCP_MAX_PAGE_SIZE":         "0",
		"MCP_MAX_RESPONSE_BYTES":    "10",
		"MCP_TOOL_TIMEOUT":          "0s",
		"MCP_RATE_WINDOW":           "0s",
		"MCP_RATE_USER_CATALOGUE":   "0",
		"MCP_RATE_USER_PERSONAL":    "-1",
		"MCP_RATE_USER_WRITE":       "0",
		"MCP_RATE_CLIENT_CATALOGUE": "0",
		"MCP_RATE_CLIENT_PERSONAL":  "0",
		"MCP_RATE_CLIENT_WRITE":     "0",
	} {
		t.Run(env, func(t *testing.T) {
			enableMCP(t)
			t.Setenv(env, value)

			cfg, err := Load()

			require.Nil(t, cfg)
			require.ErrorContains(t, err, env)

			t.Setenv("MCP_ENABLED", "false")
			_, err = Load()
			require.NoError(t, err, "an inert endpoint is not validated")
		})
	}
}

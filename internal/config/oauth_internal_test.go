package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Long enough to be accepted; readable on purpose so it is no one's secret.
const validOAuthKey = "test-oauth-key-test-oauth-key-test"

func TestLoadLeavesOAuthOffAndUnconfiguredByDefault(t *testing.T) {
	setRequiredSecrets(t)

	cfg, err := LoadServer()

	require.NoError(t, err)
	require.False(t, cfg.OAuth.Enabled)
}

func TestLoadAcceptsAnEnabledOAuthServer(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)

	cfg, err := LoadServer()

	require.NoError(t, err)
	require.True(t, cfg.OAuth.Enabled)
	require.Equal(t, "https://code-basics.com", cfg.OAuth.Issuer)
}

func TestLoadRejectsAnEnabledOAuthServerWithoutWhatItNeeds(t *testing.T) {
	cases := map[string]struct{ issuer, key, want string }{
		"no issuer":      {"", validOAuthKey, "OAUTH_ISSUER"},
		"issuer path":    {"https://code-basics.com/app", validOAuthKey, "bare origin"},
		"issuer no host": {"https:///", validOAuthKey, "OAUTH_ISSUER"},
		"short key":      {"https://code-basics.com", "short", "OAUTH_ENCRYPTION_KEY must be at least 32"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			setRequiredSecrets(t)
			t.Setenv("OAUTH_ENABLED", "true")
			t.Setenv("OAUTH_ISSUER", tc.issuer)
			t.Setenv("OAUTH_ENCRYPTION_KEY", tc.key)

			cfg, err := LoadServer()

			require.Nil(t, cfg)
			requireProblem(t, err, tc.want)
		})
	}
}

func TestLoadRejectsAnHTTPIssuerInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "http://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)

	cfg, err := LoadServer()

	require.Nil(t, cfg)
	requireProblem(t, err, "OAUTH_ISSUER must be an https origin in production")
}

func TestLoadGivesTheRegistrationLimitADefaultAndRejectsADisabledOne(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)

	cfg, err := LoadServer()
	require.NoError(t, err)
	require.Equal(t, 10, cfg.OAuth.RegistrationLimit)
	require.Equal(t, time.Hour, cfg.OAuth.RegistrationWindow)

	t.Setenv("OAUTH_REGISTRATION_LIMIT", "0")

	_, err = LoadServer()
	requireProblem(t, err, "OAUTH_REGISTRATION_LIMIT")
}

func TestLoadGivesTheTokenLifetimesDefaultsAndRejectsBadOnes(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)

	cfg, err := LoadServer()
	require.NoError(t, err)
	require.Equal(t, time.Hour, cfg.OAuth.AccessTokenTTL)
	require.Equal(t, 30*24*time.Hour, cfg.OAuth.RefreshTokenTTL)

	t.Setenv("OAUTH_ACCESS_TOKEN_TTL", "2h")
	t.Setenv("OAUTH_REFRESH_TOKEN_TTL", "48h")

	cfg, err = LoadServer()
	require.NoError(t, err)
	require.Equal(t, 2*time.Hour, cfg.OAuth.AccessTokenTTL)
	require.Equal(t, 48*time.Hour, cfg.OAuth.RefreshTokenTTL)

	for name, env := range map[string][2]string{
		"zero access":           {"0s", "48h"},
		"refresh not outliving": {"2h", "2h"},
		"refresh shorter":       {"2h", "1h"},
		"access under a minute": {"30s", "48h"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("OAUTH_ACCESS_TOKEN_TTL", env[0])
			t.Setenv("OAUTH_REFRESH_TOKEN_TTL", env[1])

			_, err := LoadServer()
			requireProblem(t, err, "OAUTH_ACCESS_TOKEN_TTL")
		})
	}

	t.Setenv("OAUTH_ACCESS_TOKEN_TTL", "soon")

	_, err = LoadServer()
	require.Error(t, err, "an unparsable duration is a config error")
}

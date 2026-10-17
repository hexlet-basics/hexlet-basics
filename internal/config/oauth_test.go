package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Long enough to be accepted; readable on purpose so it is no one's secret.
const validOAuthKey = "test-oauth-key-test-oauth-key-test"

func TestLoadLeavesOAuthOffAndUnconfiguredByDefault(t *testing.T) {
	setRequiredSecrets(t)

	cfg, err := Load()

	require.NoError(t, err)
	require.False(t, cfg.OAuth.Enabled)
}

func TestLoadAcceptsAnEnabledOAuthServer(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)

	cfg, err := Load()

	require.NoError(t, err)
	require.True(t, cfg.OAuth.Enabled)
	require.Equal(t, "https://code-basics.com", cfg.OAuth.Issuer)
}

func TestLoadRejectsAnEnabledOAuthServerWithoutWhatItNeeds(t *testing.T) {
	cases := map[string]struct{ issuer, key, want string }{
		"no issuer":      {"", validOAuthKey, "OAUTH_ISSUER must be an origin"},
		"issuer path":    {"https://code-basics.com/app", validOAuthKey, "bare origin"},
		"issuer no host": {"https:///", validOAuthKey, "OAUTH_ISSUER must be an origin"},
		"short key":      {"https://code-basics.com", "short", "OAUTH_ENCRYPTION_KEY must be at least 32"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			setRequiredSecrets(t)
			t.Setenv("OAUTH_ENABLED", "true")
			t.Setenv("OAUTH_ISSUER", tc.issuer)
			t.Setenv("OAUTH_ENCRYPTION_KEY", tc.key)

			cfg, err := Load()

			require.Nil(t, cfg)
			require.ErrorContains(t, err, tc.want)
		})
	}
}

func TestLoadRejectsAnHTTPIssuerInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "http://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", validOAuthKey)

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "OAUTH_ISSUER must be an https origin in production")
}

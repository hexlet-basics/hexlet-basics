package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// requireProblem asserts err is a failure whose message or whose problems list
// (the "problems" attribute of an invalid configuration) mentions want.
func requireProblem(t *testing.T, err error, want string) {
	t.Helper()
	require.Error(t, err)
	require.Contains(t, err.Error()+"\n"+strings.Join(Problems(err), "\n"), want)
}

// Every bad setting is reported in one error, each named by its variable, so a
// deployment is fixed in one round instead of one restart per mistake.
func TestLoadServerReportsEverySettingProblemAtOnce(t *testing.T) {
	tests := []struct {
		name  string
		env   map[string]string
		unset []string
		want  []string
	}{
		{
			name:  "missing and malformed together",
			unset: []string{"DATABASE_URL", "GUEST_COOKIE_SECRET"},
			env:   map[string]string{"APP_ENV": "staging", "MCP_ENABLED": "true"},
			want:  []string{"DATABASE_URL", "GUEST_COOKIE_SECRET", "APP_ENV", "MCP_ENABLED"},
		},
		{
			name: "development secrets in production",
			env: map[string]string{
				"APP_ENV":             "production",
				"GUEST_COOKIE_SECRET": developmentGuestCookieSecret,
				"EMAIL_TOKEN_SECRET":  developmentEmailTokenSecret,
			},
			want:  []string{"GUEST_COOKIE_SECRET", "EMAIL_TOKEN_SECRET"},
			unset: nil,
		},
		{
			name:  "oauth enabled with nothing else set",
			env:   map[string]string{"OAUTH_ENABLED": "true"},
			want:  []string{"OAUTH_ISSUER", "OAUTH_ENCRYPTION_KEY"},
			unset: nil,
		},
		{
			name: "oauth issuer with a path and a short key",
			env: map[string]string{
				"OAUTH_ENABLED":        "true",
				"OAUTH_ISSUER":         "https://code-basics.com/oauth",
				"OAUTH_ENCRYPTION_KEY": "short",
			},
			want:  []string{"OAUTH_ISSUER", "OAUTH_ENCRYPTION_KEY"},
			unset: nil,
		},
		{
			name: "oauth lifetimes out of order and mcp limits at zero",
			env: map[string]string{
				"OAUTH_ENABLED":            "true",
				"OAUTH_ISSUER":             "https://code-basics.com",
				"OAUTH_ENCRYPTION_KEY":     validOAuthKey,
				"OAUTH_ACCESS_TOKEN_TTL":   "30s",
				"OAUTH_REFRESH_TOKEN_TTL":  "10s",
				"OAUTH_REGISTRATION_LIMIT": "0",
				"MCP_ENABLED":              "true",
				"MCP_MAX_PAGE_SIZE":        "0",
				"MCP_TOOL_TIMEOUT":         "100ms",
			},
			want: []string{
				"OAUTH_ACCESS_TOKEN_TTL", "OAUTH_REFRESH_TOKEN_TTL", "OAUTH_REGISTRATION_LIMIT",
				"MCP_MAX_PAGE_SIZE", "MCP_TOOL_TIMEOUT",
			},
			unset: nil,
		},
		{
			name: "http oauth issuer in production",
			env: map[string]string{
				"APP_ENV":              "production",
				"OAUTH_ENABLED":        "true",
				"OAUTH_ISSUER":         "http://code-basics.com",
				"OAUTH_ENCRYPTION_KEY": validOAuthKey,
			},
			want:  []string{"OAUTH_ISSUER"},
			unset: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRequiredSecrets(t)

			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			clearEnv(t, tt.unset...)

			cfg, err := LoadServer()

			require.Nil(t, cfg)
			require.Error(t, err)

			for _, name := range tt.want {
				requireProblem(t, err, name)
			}
		})
	}
}

func TestLoadWorkerReportsEverySettingProblemAtOnce(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("EMAIL_TOKEN_SECRET", developmentEmailTokenSecret)
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "")
	t.Setenv("AUDIT_RETENTION_MONTHS", "0")
	clearEnv(t, "SITE_URL")

	cfg, err := LoadWorker()

	require.Nil(t, cfg)

	for _, name := range []string{"SITE_URL", "EMAIL_TOKEN_SECRET", "MAIL_POSTBOX_ACCESS_KEY_ID", "AUDIT_RETENTION_MONTHS"} {
		requireProblem(t, err, name)
	}
}

// A message names a variable and never carries the value of a secret.
func TestLoadErrorNeverContainsASecretValue(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.com")
	t.Setenv("OAUTH_ENCRYPTION_KEY", "short-but-secret")

	_, err := LoadServer()

	requireProblem(t, err, "OAUTH_ENCRYPTION_KEY")
	require.NotContains(t, err.Error(), "short-but-secret")
}

func TestLoadServerHidesSecretsWhenTheConfigIsPrinted(t *testing.T) {
	setRequiredSecrets(t)

	cfg, err := LoadServer()

	require.NoError(t, err)
	require.NotContains(t, fmt.Sprintf("%+v", cfg), "production-secret-from-secret-store")
	require.Equal(t, "production-secret-from-secret-store", cfg.GuestCookieSecret.Reveal())
}

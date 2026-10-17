package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/require"
)

func TestConfigRequiresGuestCookieSecret(t *testing.T) {
	cfg := &Config{}
	err := env.ParseWithOptions(cfg, env.Options{
		// A non-empty map replaces the process environment, making the absence
		// of GUEST_COOKIE_SECRET deterministic even on configured developer machines.
		Environment: map[string]string{"UNRELATED": "value"},
	})

	require.ErrorContains(t, err, `required environment variable "GUEST_COOKIE_SECRET" is not set`)
}

// setRequiredSecrets fills every required secret with a non-development value,
// so each test changes only the one it is about.
func setRequiredSecrets(t *testing.T) {
	t.Helper()
	t.Setenv("GUEST_COOKIE_SECRET", "production-secret-from-secret-store")
	t.Setenv("EMAIL_TOKEN_SECRET", "production-email-secret-from-secret-store")
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "postbox-key-id")
}

func TestLoadRejectsEmptyGuestCookieSecret(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("GUEST_COOKIE_SECRET", "")

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, `environment variable "GUEST_COOKIE_SECRET" should not be empty`)
}

func TestLoadRejectsDevelopmentGuestCookieSecretInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("GUEST_COOKIE_SECRET", developmentGuestCookieSecret)

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "validate production config")
	require.ErrorContains(t, err, "GUEST_COOKIE_SECRET must differ from the public development value")
}

func TestLoadAcceptsExplicitProductionGuestCookieSecret(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("GUEST_COOKIE_SECRET", "production-secret-from-secret-store")

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, "production-secret-from-secret-store", cfg.GuestCookieSecret)
}

func TestLoadAcceptsDevelopmentGuestCookieSecretOutsideProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("SENTRY_ENVIRONMENT", "development")
	t.Setenv("GUEST_COOKIE_SECRET", developmentGuestCookieSecret)

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, developmentGuestCookieSecret, cfg.GuestCookieSecret)
}

func TestLoadRejectsEmptyEmailTokenSecret(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("EMAIL_TOKEN_SECRET", "")

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, `environment variable "EMAIL_TOKEN_SECRET" should not be empty`)
}

func TestLoadRejectsDevelopmentEmailTokenSecretInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("EMAIL_TOKEN_SECRET", developmentEmailTokenSecret)

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "EMAIL_TOKEN_SECRET must differ from the public development value")
}

func TestLoadRequiresPostboxInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "")

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "MAIL_POSTBOX_ACCESS_KEY_ID is required in production")
}

func TestLoadLogsEmailOutsideProductionWithoutPostbox(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "")

	cfg, err := Load()

	require.NoError(t, err)
	require.False(t, cfg.Mail.PostboxEnabled())
	require.Equal(t, "support@hexlet.io", cfg.Mail.From)
}

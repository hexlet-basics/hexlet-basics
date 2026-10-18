package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// setRequiredSecrets fills every required setting of both processes with a
// non-development value, so each test changes only the one it is about.
const (
	developmentGuestCookieSecret = "dev-insecure-guest-cookie-secret-change-me"
	developmentEmailTokenSecret  = "dev-insecure-email-token-secret-change-me"
)

func setRequiredSecrets(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://app@db.example/app")
	t.Setenv("PUBLIC_URL", "https://code-basics.com")
	t.Setenv("SITE_URL", "https://code-basics.com")
	t.Setenv("APP_HOST", "code-basics.com")
	t.Setenv("BLOB_BUCKET_URL", "file://./storage")
	t.Setenv("GUEST_COOKIE_SECRET", "production-secret-from-secret-store")
	t.Setenv("EMAIL_TOKEN_SECRET", "production-email-secret-from-secret-store")
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "postbox-key-id")
}

// clearEnv blanks a variable so a test sees it as not set: t.Setenv("X", "")
// would be "set but empty", which `required` accepts.
func clearEnv(t *testing.T, names ...string) {
	t.Helper()
	for _, name := range names {
		t.Setenv(name, "")
		require.NoError(t, os.Unsetenv(name))
	}
}

func TestLoadServerNamesEveryMissingRequiredSetting(t *testing.T) {
	setRequiredSecrets(t)
	clearEnv(t, "DATABASE_URL", "PUBLIC_URL", "SITE_URL", "APP_HOST", "BLOB_BUCKET_URL", "GUEST_COOKIE_SECRET")

	cfg, err := LoadServer()

	require.Nil(t, cfg)
	for _, name := range []string{"DATABASE_URL", "PUBLIC_URL", "SITE_URL", "APP_HOST", "BLOB_BUCKET_URL", "GUEST_COOKIE_SECRET"} {
		require.ErrorContains(t, err, name)
	}
}

func TestLoadWorkerNamesAMissingDeployCriticalSetting(t *testing.T) {
	setRequiredSecrets(t)
	clearEnv(t, "DATABASE_URL")

	cfg, err := LoadWorker()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "DATABASE_URL")
}

// The worker serves no visitor, so it has no cookie to sign and no public host
// to build feeds on: a deployment must not have to give it those.
func TestLoadWorkerNeedsNoServerOnlySetting(t *testing.T) {
	setRequiredSecrets(t)
	clearEnv(t, "GUEST_COOKIE_SECRET", "APP_HOST")

	cfg, err := LoadWorker()

	require.NoError(t, err)
	require.Equal(t, "postgres://app@db.example/app", cfg.DatabaseURL)
}

func TestLoadRefusesAnUnknownAppEnv(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "staging")

	_, err := LoadServer()
	require.ErrorContains(t, err, "APP_ENV")

	_, err = LoadWorker()
	require.ErrorContains(t, err, "APP_ENV")
}

func TestAppEnvDefaultsToDevelopment(t *testing.T) {
	setRequiredSecrets(t)
	clearEnv(t, "APP_ENV")

	cfg, err := LoadServer()

	require.NoError(t, err)
	require.Equal(t, Development, cfg.AppEnv)
}

func TestLoadRejectsEmptyGuestCookieSecret(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("GUEST_COOKIE_SECRET", "")

	cfg, err := LoadServer()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "GUEST_COOKIE_SECRET is required")
}

func TestLoadRejectsDevelopmentGuestCookieSecretInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("GUEST_COOKIE_SECRET", developmentGuestCookieSecret)

	cfg, err := LoadServer()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "invalid configuration")
	require.ErrorContains(t, err, "GUEST_COOKIE_SECRET must differ from the public development value")
}

func TestLoadAcceptsExplicitProductionGuestCookieSecret(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("GUEST_COOKIE_SECRET", "production-secret-from-secret-store")

	cfg, err := LoadServer()

	require.NoError(t, err)
	require.Equal(t, "production-secret-from-secret-store", cfg.GuestCookieSecret.Reveal())
}

func TestLoadAcceptsDevelopmentGuestCookieSecretOutsideProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "development")
	t.Setenv("GUEST_COOKIE_SECRET", developmentGuestCookieSecret)

	cfg, err := LoadServer()

	require.NoError(t, err)
	require.Equal(t, developmentGuestCookieSecret, cfg.GuestCookieSecret.Reveal())
}

func TestLoadRejectsEmptyEmailTokenSecret(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("EMAIL_TOKEN_SECRET", "")

	cfg, err := LoadServer()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "EMAIL_TOKEN_SECRET is required")
}

func TestLoadRejectsDevelopmentEmailTokenSecretInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("EMAIL_TOKEN_SECRET", developmentEmailTokenSecret)

	cfg, err := LoadServer()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "EMAIL_TOKEN_SECRET must differ from the public development value")
}

func TestLoadWorkerRequiresPostboxInProduction(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "")

	cfg, err := LoadWorker()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "MAIL_POSTBOX_ACCESS_KEY_ID is required in production")
}

func TestLoadLogsEmailOutsideProductionWithoutPostbox(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("MAIL_POSTBOX_ACCESS_KEY_ID", "")

	cfg, err := LoadWorker()

	require.NoError(t, err)
	require.False(t, cfg.Mail.PostboxEnabled())
	require.Equal(t, "support@hexlet.io", cfg.Mail.From)
}

func TestAuditRetentionDefaultsToTwoYears(t *testing.T) {
	setRequiredSecrets(t)

	cfg, err := LoadWorker()

	require.NoError(t, err)
	require.Equal(t, 24, cfg.AuditRetentionMonths)
}

func TestAuditRetentionCanBeShortenedButNotRemoved(t *testing.T) {
	setRequiredSecrets(t)
	t.Setenv("AUDIT_RETENTION_MONTHS", "6")
	cfg, err := LoadWorker()
	require.NoError(t, err)
	require.Equal(t, 6, cfg.AuditRetentionMonths)

	t.Setenv("AUDIT_RETENTION_MONTHS", "0")
	_, err = LoadWorker()
	require.ErrorContains(t, err, "AUDIT_RETENTION_MONTHS")
}

package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/require"
)

func TestConfigRequiresJWTSecret(t *testing.T) {
	cfg := &Config{}
	err := env.ParseWithOptions(cfg, env.Options{
		// A non-empty map replaces the process environment, making the absence
		// of JWT_SECRET deterministic even on configured developer machines.
		Environment: map[string]string{"UNRELATED": "value"},
	})

	require.ErrorContains(t, err, `required environment variable "JWT_SECRET" is not set`)
}

func TestLoadRejectsEmptyJWTSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, `environment variable "JWT_SECRET" should not be empty`)
}

func TestLoadRejectsDevelopmentJWTSecretInProduction(t *testing.T) {
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", developmentJWTSecret)

	cfg, err := Load()

	require.Nil(t, cfg)
	require.ErrorContains(t, err, "validate production config")
	require.ErrorContains(t, err, "JWT_SECRET must differ from the public development value")
}

func TestLoadAcceptsExplicitProductionJWTSecret(t *testing.T) {
	t.Setenv("SENTRY_ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", "production-secret-from-secret-store")

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, "production-secret-from-secret-store", cfg.JWTSecret)
}

func TestLoadAcceptsDevelopmentJWTSecretOutsideProduction(t *testing.T) {
	t.Setenv("SENTRY_ENVIRONMENT", "development")
	t.Setenv("JWT_SECRET", developmentJWTSecret)

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, developmentJWTSecret, cfg.JWTSecret)
}

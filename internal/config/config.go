// Package config loads runtime configuration from the environment.
package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/samber/oops"
)

const developmentJWTSecret = "dev-insecure-jwt-secret-change-me"

// Config holds the runtime configuration for the server. Fields are populated
// from environment variables (12-factor); defaults keep local dev zero-config.
type Config struct {
	Addr        string `env:"ADDR" envDefault:":3001"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@127.0.0.1:54330/code_basics_development"`
	// BlobBucketURL selects the gocloud.dev/blob backend (ADR-0005): local disk
	// (fileblob) in dev, s3blob in prod, chosen by URL scheme. The dev default
	// writes under ./storage (created on first open) so uploads survive restarts.
	BlobBucketURL string `env:"BLOB_BUCKET_URL" envDefault:"file://./storage?create_dir=true"`
	// AppHost is the public site host used to build canonical page URLs (e.g. a
	// blog post's `url`), mirroring legacy AppHost.canonical. HTTPS is assumed in
	// prod; the default matches the legacy default.
	AppHost string `env:"APP_HOST" envDefault:"code-basics.com"`
	// PublicURL is this server's own public origin, used to build absolute asset
	// URLs it serves itself (the `/storage/{key}` blob read path). Separate from
	// AppHost because the API may sit on a different host than the site.
	PublicURL string `env:"PUBLIC_URL" envDefault:"http://localhost:3001"`
	// CourseRepoBaseURL is the GitHub org base under which each course's exercises
	// repo lives: the loader clones `<base>/exercises-<slug>.git`. The default
	// matches the legacy hard-coded convention (apiconv.repositoryURL).
	CourseRepoBaseURL string `env:"COURSE_REPO_BASE_URL" envDefault:"https://github.com/hexlet-basics"`
	// GitHubToken authenticates clones (private repos / higher rate limits). The
	// exercises repos are public, so it is optional; when set it is used as the
	// HTTP basic-auth password (x-access-token).
	GitHubToken string `env:"GITHUB_TOKEN"`
	// GitHubWebhookSecret is the shared secret GitHub signs webhook deliveries
	// with (HMAC-SHA256). The webhook route rejects any delivery whose signature
	// doesn't verify; an empty secret disables the webhook (fail closed).
	GitHubWebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	// amoCRM receives lead_created events through a Watermill -> River handler.
	AmoCRMBaseURL   string `env:"AMOCRM_BASE_URL"`
	AmoCRMAuthToken string `env:"AMOCRM_AUTH_TOKEN"`
	YMCounter       string `env:"YM_COUNTER"`
	// SentryDSN enables exception delivery. An empty DSN intentionally creates a
	// disabled client so local development and tests stay off-process.
	SentryDSN         string `env:"SENTRY_DSN"`
	SentryEnvironment string `env:"SENTRY_ENVIRONMENT" envDefault:"development"`
	ReleaseVersion    string `env:"HEXLET_BASICS_RELEASE_VERSION"`
	// JWTSecret signs the session JWT stored in the auth cookie (ADR-0003). The
	// committed .env.example supplies a development value, while deployments
	// must provide a non-empty secret explicitly.
	JWTSecret string `env:"JWT_SECRET,required,notEmpty"`
}

// Load reads configuration from environment variables, applying defaults.
// A local .env is loaded first if present; it only fills variables that are not
// already set in the real environment (godotenv never overrides), so exported
// vars and 12-factor deploys keep precedence over the dev-convenience file.
func Load() (*Config, error) {
	// Ignore a missing file: .env is a local convenience, absent in prod/CI.
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, oops.Wrapf(err, "parse config from environment")
	}
	if err := validateProduction(cfg); err != nil {
		return nil, oops.Wrapf(err, "validate production config")
	}
	return cfg, nil
}

// validateProduction keeps environment-specific safety policy separate from
// env decoding. The development secret is committed for local convenience and
// therefore must never be accepted as a session-signing key in production.
func validateProduction(cfg *Config) error {
	if cfg.SentryEnvironment == "production" && cfg.JWTSecret == developmentJWTSecret {
		return oops.Errorf("JWT_SECRET must differ from the public development value")
	}
	return nil
}

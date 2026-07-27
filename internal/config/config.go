// Package config loads runtime configuration from the environment.
package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/samber/oops"
)

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
	return cfg, nil
}

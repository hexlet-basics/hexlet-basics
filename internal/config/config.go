// Package config loads runtime configuration from the environment.
package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/samber/oops"
)

// Config holds the runtime configuration for the server. Fields are populated
// from environment variables (12-factor); defaults keep local dev zero-config.
type Config struct {
	Addr        string `env:"ADDR" envDefault:":3001"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@127.0.0.1:54330/code_basics_development"`
}

// Load reads configuration from environment variables, applying defaults.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, oops.Wrapf(err, "parse config from environment")
	}
	return cfg, nil
}

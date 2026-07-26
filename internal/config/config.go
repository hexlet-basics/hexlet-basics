// Package config loads runtime configuration from the environment.
package config

import (
	"strings"

	"github.com/samber/oops"
	"github.com/spf13/viper"
)

// Config holds the runtime configuration for the server.
type Config struct {
	Addr        string
	DatabaseURL string
}

// Load reads configuration from environment variables, applying defaults.
func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("addr", ":3001")
	v.SetDefault("database_url", "postgres://postgres:postgres@127.0.0.1:54330/code_basics_development")

	cfg := &Config{
		Addr:        v.GetString("addr"),
		DatabaseURL: v.GetString("database_url"),
	}

	if cfg.DatabaseURL == "" {
		return nil, oops.Errorf("database_url is required")
	}

	return cfg, nil
}

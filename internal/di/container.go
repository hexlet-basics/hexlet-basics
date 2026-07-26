// Package di wires the application's dependency-injection container.
package di

import (
	"log/slog"

	"github.com/samber/do/v2"

	"hexletbasics/internal/config"
	"hexletbasics/internal/logging"
)

// New builds the DI container and registers the platform services (configuration
// and logger) that the rest of the application depends on.
func New() do.Injector {
	injector := do.New()

	do.Provide(injector, func(do.Injector) (*config.Config, error) {
		return config.Load()
	})

	do.Provide(injector, func(do.Injector) (*slog.Logger, error) {
		return logging.New(slog.LevelInfo), nil
	})

	return injector
}

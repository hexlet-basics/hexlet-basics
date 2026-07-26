// Package di wires the application's dependency-injection container.
package di

import (
	"log/slog"

	"github.com/samber/do/v2"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/logging"
	"hexletbasics/internal/store"
)

// New builds the DI container and registers the application's services.
// Providers resolve their own dependencies from the injector, so the plain
// constructors (store.NewClient, handlers.NewServer, api.NewServer) stay
// injector-agnostic and remain usable directly in tests.
func New() do.Injector {
	injector := do.New()

	do.Provide(injector, func(do.Injector) (*config.Config, error) {
		return config.Load()
	})

	do.Provide(injector, func(do.Injector) (*slog.Logger, error) {
		return logging.New(slog.LevelInfo), nil
	})

	do.Provide(injector, func(i do.Injector) (*ent.Client, error) {
		return store.NewClient(do.MustInvoke[*config.Config](i).DatabaseURL)
	})

	do.Provide(injector, func(i do.Injector) (*handlers.Server, error) {
		return handlers.NewServer(do.MustInvoke[*ent.Client](i)), nil
	})

	do.Provide(injector, func(i do.Injector) (*api.Server, error) {
		return api.NewServer(do.MustInvoke[*handlers.Server](i))
	})

	return injector
}

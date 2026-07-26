// Package di wires the application's dependency-injection container.
package di

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/samber/do/v2"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/logging"
	"hexletbasics/internal/store"
)

// HTTP server timeouts. ReadHeaderTimeout caps slow-header (Slowloris) clients;
// the rest bound a misbehaving peer without cutting normal API traffic.
const (
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
)

// New builds the DI container and registers the application's services.
// Providers resolve their own dependencies from the injector, so the plain
// constructors (store.NewClient, handlers.NewServer, api.NewServer) stay
// injector-agnostic and remain usable directly in tests.
//
// It returns the concrete *do.RootScope (not the do.Injector interface) so the
// caller can drive the injector's lifecycle helpers (graceful shutdown).
func New() *do.RootScope {
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

	// *http.Server natively satisfies do's ShutdownerWithContextAndError (its
	// Shutdown(ctx) drains in-flight requests), so injector.Shutdown gracefully
	// stops it — no bespoke teardown wiring needed.
	do.Provide(injector, func(i do.Injector) (*http.Server, error) {
		cfg := do.MustInvoke[*config.Config](i)
		apiServer := do.MustInvoke[*api.Server](i)

		// Dev CORS so the Vite frontend (any localhost port) can call the API.
		handler := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:*", "http://127.0.0.1:*"},
		}).Handler(apiServer)

		return &http.Server{
			Addr:              cfg.Addr,
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		}, nil
	})

	return injector
}

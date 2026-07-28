// Package di wires the application's dependency-injection container.
package di

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/rs/cors"
	"github.com/samber/do/v2"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
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
		return handlers.NewServer(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*config.Config](i),
			do.MustInvoke[*river.Client[pgx.Tx]](i),
		), nil
	})

	// Blob bucket for uploaded assets (ADR-0005). Closed explicitly in main.go on
	// shutdown, like the ent client and pgx pool (blob.Bucket has Close, not a do
	// Shutdowner).
	do.Provide(injector, func(i do.Injector) (*blob.Bucket, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return store.NewBucket(context.Background(), cfg.BlobBucketURL)
	})

	// pgx pool + river client back the background-job queue (ADR-0004). The pool
	// is separate from ent's database/sql handle because riverpgxv5 needs a
	// native *pgxpool.Pool.
	do.Provide(injector, func(i do.Injector) (*pgxpool.Pool, error) {
		return store.NewPool(context.Background(), do.MustInvoke[*config.Config](i).DatabaseURL)
	})

	// The exercise loader (course-version builds) needs the ent client for writes,
	// the blob bucket for lesson theory images, and a git fetcher for the repo.
	do.Provide(injector, func(i do.Injector) (courseloader.Fetcher, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return courseloader.NewGitFetcher(cfg.CourseRepoBaseURL, cfg.GitHubToken), nil
	})

	do.Provide(injector, func(i do.Injector) (*courseloader.Loader, error) {
		return courseloader.NewLoader(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*blob.Bucket](i),
			do.MustInvoke[courseloader.Fetcher](i),
			do.MustInvoke[*config.Config](i),
		), nil
	})

	do.Provide(injector, func(i do.Injector) (*river.Client[pgx.Tx], error) {
		return jobs.NewClient(
			do.MustInvoke[*pgxpool.Pool](i),
			do.MustInvoke[*courseloader.Loader](i),
		)
	})

	do.Provide(injector, func(i do.Injector) (*api.Server, error) {
		// WithErrorHandler installs the central ent-error -> HTTP-status mapping
		// (404/409), so handlers return raw ent errors instead of typed DTOs.
		return api.NewServer(
			do.MustInvoke[*handlers.Server](i),
			api.WithErrorHandler(handlers.APIErrorHandler),
		)
	})

	// *http.Server natively satisfies do's ShutdownerWithContextAndError (its
	// Shutdown(ctx) drains in-flight requests), so injector.Shutdown gracefully
	// stops it — no bespoke teardown wiring needed.
	do.Provide(injector, func(i do.Injector) (*http.Server, error) {
		cfg := do.MustInvoke[*config.Config](i)
		apiServer := do.MustInvoke[*api.Server](i)

		// Compose the generated server with the multipart/blob routes ogen can't
		// generate (ADR-0005), then wrap the whole router in CORS so both surfaces
		// are covered (and a future auth middleware should wrap here too).
		att := handlers.NewAttachmentHandler(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*blob.Bucket](i),
		)
		gh := handlers.NewGitHubWebhookHandler(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*river.Client[pgx.Tx]](i),
			cfg.GitHubWebhookSecret,
		)
		auth := handlers.NewAuthHandler(do.MustInvoke[*ent.Client](i), cfg)
		router := handlers.NewRouter(apiServer, att, gh, auth)

		// Dev CORS so the Vite frontend (any localhost port) can call the API.
		// AllowCredentials is required for the auth cookie to round-trip
		// cross-origin (the SPA fetches with `credentials: include`).
		handler := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
			AllowCredentials: true,
		}).Handler(router)

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

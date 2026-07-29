// Package di wires the application's dependency-injection container.
package di

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

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
	"hexletbasics/internal/versionbuilds"
)

// HTTP server timeouts. ReadHeaderTimeout caps slow-header (Slowloris) clients;
// the rest bound a misbehaving peer without cutting normal API traffic.
const (
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second

	routerServiceName      = "router"
	httpHandlerServiceName = "http-handler"
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

	do.Provide(injector, func(i do.Injector) (*sql.DB, error) {
		return store.NewDB(do.MustInvoke[*config.Config](i).DatabaseURL)
	})

	do.Provide(injector, func(i do.Injector) (*ent.Client, error) {
		return store.NewClient(do.MustInvoke[*sql.DB](i)), nil
	})

	do.Provide(injector, func(i do.Injector) (*handlers.Server, error) {
		return handlers.NewServer(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*config.Config](i),
			do.MustInvoke[*versionbuilds.Starter](i),
		), nil
	})

	// Blob bucket for uploaded assets (ADR-0005). Closed explicitly in main.go on
	// shutdown, like the ent client and pgx pool (blob.Bucket has Close, not a do
	// Shutdowner).
	do.Provide(injector, func(i do.Injector) (*blob.Bucket, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return store.NewBucket(context.Background(), cfg.BlobBucketURL)
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

	do.Provide(injector, func(i do.Injector) (*river.Client[*sql.Tx], error) {
		return jobs.NewClient(
			do.MustInvoke[*sql.DB](i),
			do.MustInvoke[*courseloader.Loader](i),
		)
	})

	do.Provide(injector, func(i do.Injector) (*versionbuilds.Starter, error) {
		return versionbuilds.NewStarter(
			do.MustInvoke[*sql.DB](i),
			do.MustInvoke[*river.Client[*sql.Tx]](i),
		), nil
	})

	do.Provide(injector, func(i do.Injector) (*api.Server, error) {
		// WithErrorHandler installs the central ent-error -> HTTP-status mapping
		// (404/409), so handlers return raw ent errors instead of typed DTOs.
		return api.NewServer(
			do.MustInvoke[*handlers.Server](i),
			api.WithErrorHandler(handlers.APIErrorHandler),
		)
	})

	do.Provide(injector, func(i do.Injector) (*handlers.AttachmentHandler, error) {
		return handlers.NewAttachmentHandler(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*blob.Bucket](i),
		), nil
	})

	do.Provide(injector, func(i do.Injector) (*handlers.GitHubWebhookHandler, error) {
		return handlers.NewGitHubWebhookHandler(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*versionbuilds.Starter](i),
			do.MustInvoke[*config.Config](i).GitHubWebhookSecret,
		), nil
	})

	// Both the raw router and the middleware-wrapped application handler have
	// the same http.Handler interface. Named bindings keep those two seams
	// explicit without adding wrapper types whose only purpose would be DI.
	do.ProvideNamed(injector, routerServiceName, func(i do.Injector) (http.Handler, error) {
		return handlers.NewRouter(
			do.MustInvoke[*api.Server](i),
			do.MustInvoke[*handlers.AttachmentHandler](i),
			do.MustInvoke[*handlers.GitHubWebhookHandler](i),
		), nil
	})

	do.ProvideNamed(injector, httpHandlerServiceName, func(i do.Injector) (http.Handler, error) {
		// Dev CORS lets the Vite frontend (on any localhost port) call both the
		// generated API and the hand-mounted routes. AllowCredentials is needed
		// for the auth cookie to make the cross-origin round trip.
		return cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
			AllowCredentials: true,
		}).Handler(do.MustInvokeNamed[http.Handler](i, routerServiceName)), nil
	})

	// *http.Server natively satisfies do's ShutdownerWithContextAndError (its
	// Shutdown(ctx) drains in-flight requests), so injector.Shutdown gracefully
	// stops it — no bespoke teardown wiring needed.
	do.Provide(injector, func(i do.Injector) (*http.Server, error) {
		cfg := do.MustInvoke[*config.Config](i)

		return &http.Server{
			Addr:              cfg.Addr,
			Handler:           do.MustInvokeNamed[http.Handler](i, httpHandlerServiceName),
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		}, nil
	})

	return injector
}

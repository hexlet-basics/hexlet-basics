// Package di wires the application's dependency-injection container.
package di

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/rs/cors"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/amocrm"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/config"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/events"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/logging"
	"hexletbasics/internal/store"
	"hexletbasics/internal/telemetry"
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

// NewServer builds the synchronous application graph. Its River client can
// enqueue jobs but has no queues or workers, so this process cannot execute
// background work.
func NewServer() *do.RootScope {
	return newContainer(false)
}

// NewWorker builds the asynchronous application graph. It owns the Watermill
// subscribers and the River workers but exposes no HTTP server.
func NewWorker() *do.RootScope {
	return newContainer(true)
}

// newContainer registers the process-specific graph while keeping construction
// providers local to this package. cmd/server and cmd/worker remain the sole
// lifecycle owners for the objects they resolve.
func newContainer(worker bool) *do.RootScope {
	injector := do.New()

	do.Provide(injector, func(do.Injector) (*config.Config, error) {
		return config.Load()
	})

	do.Provide(injector, func(do.Injector) (*slog.Logger, error) {
		return logging.New(slog.LevelInfo), nil
	})

	do.Provide(injector, func(i do.Injector) (*sentry.Client, error) {
		return telemetry.NewSentryClient(do.MustInvoke[*config.Config](i))
	})

	do.Provide(injector, func(i do.Injector) (*otelconf.SDK, error) {
		return telemetry.NewOpenTelemetrySDK(
			context.Background(),
			do.MustInvoke[*slog.Logger](i),
		)
	})

	do.Provide(injector, func(i do.Injector) (*handlers.APIErrorHandler, error) {
		return handlers.NewAPIErrorHandler(
			do.MustInvoke[*localization.Translator](i),
			do.MustInvoke[*slog.Logger](i),
			do.MustInvoke[*sentry.Client](i),
		), nil
	})

	do.Provide(injector, func(do.Injector) (*localization.Translator, error) {
		return localization.New()
	})

	do.Provide(injector, func(i do.Injector) (*sql.DB, error) {
		return store.NewDB(do.MustInvoke[*config.Config](i).DatabaseURL)
	})

	do.Provide(injector, func(i do.Injector) (*ent.Client, error) {
		return store.NewClient(do.MustInvoke[*sql.DB](i)), nil
	})

	do.Provide(injector, func(i do.Injector) (*store.Store, error) {
		return store.New(do.MustInvoke[*sql.DB](i)), nil
	})

	if !worker {
		do.Provide(injector, func(i do.Injector) (*events.Publisher, error) {
			return events.NewPublisher(
				do.MustInvoke[*store.Store](i),
				do.MustInvoke[*slog.Logger](i),
			), nil
		})

		do.Provide(injector, func(i do.Injector) (*accounts.Registrar, error) {
			return accounts.NewRegistrar(
				do.MustInvoke[*store.Store](i),
				do.MustInvoke[*events.Publisher](i),
			), nil
		})
	}

	// Blob bucket for uploaded assets (ADR-0005). Closed explicitly in main.go on
	// shutdown, like the ent client and pgx pool (blob.Bucket has Close, not a do
	// Shutdowner).
	do.Provide(injector, func(i do.Injector) (*blob.Bucket, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return store.NewBucket(context.Background(), cfg.BlobBucketURL)
	})

	// Asset storage owns MIME policy, blob writes, attachment persistence,
	// compensation, and public URLs for every upload path.
	do.Provide(injector, func(i do.Injector) (*assetstore.Store, error) {
		return assetstore.New(
			do.MustInvoke[*ent.Client](i),
			do.MustInvoke[*blob.Bucket](i),
			do.MustInvoke[*config.Config](i).PublicURL,
		), nil
	})

	if worker {
		// The exercise loader (course-version builds) uses the shared asset store
		// for lesson theory images and a git fetcher for the source repository.
		do.Provide(injector, func(i do.Injector) (courseloader.Fetcher, error) {
			cfg := do.MustInvoke[*config.Config](i)
			return courseloader.NewGitFetcher(cfg.CourseRepoBaseURL, cfg.GitHubToken), nil
		})

		do.Provide(injector, func(i do.Injector) (*courseloader.Loader, error) {
			return courseloader.NewLoader(
				do.MustInvoke[*ent.Client](i),
				do.MustInvoke[*store.Store](i),
				do.MustInvoke[*assetstore.Store](i),
				do.MustInvoke[courseloader.Fetcher](i),
			), nil
		})

		do.Provide(injector, func(i do.Injector) (*amocrm.Client, error) {
			cfg := do.MustInvoke[*config.Config](i)
			return amocrm.NewClient(cfg.AmoCRMBaseURL, cfg.AmoCRMAuthToken, cfg.YMCounter), nil
		})
	}

	do.Provide(injector, func(i do.Injector) (*river.Client[*sql.Tx], error) {
		if worker {
			return jobs.NewWorkerClient(
				do.MustInvoke[*sql.DB](i),
				do.MustInvoke[*courseloader.Loader](i),
				do.MustInvoke[*amocrm.Client](i),
				do.MustInvoke[*slog.Logger](i),
				jobs.NewErrorHandler(do.MustInvoke[*sentry.Client](i)),
			)
		}
		return jobs.NewInsertOnlyClient(
			do.MustInvoke[*sql.DB](i),
			do.MustInvoke[*slog.Logger](i),
		)
	})

	if worker {
		do.Provide(injector, func(i do.Injector) (*events.Runtime, error) {
			return events.NewRuntime(
				do.MustInvoke[*sql.DB](i),
				do.MustInvoke[*slog.Logger](i),
				eventhandlers.LeadCreated(do.MustInvoke[*river.Client[*sql.Tx]](i)),
			)
		})
	} else {
		do.Provide(injector, func(i do.Injector) (*versionbuilds.Starter, error) {
			return versionbuilds.NewStarter(
				do.MustInvoke[*store.Store](i),
				do.MustInvoke[*river.Client[*sql.Tx]](i),
			), nil
		})

		do.Provide(injector, func(i do.Injector) (*handlers.Server, error) {
			return handlers.NewServer(
				do.MustInvoke[*ent.Client](i),
				do.MustInvoke[*config.Config](i),
				do.MustInvoke[*versionbuilds.Starter](i),
				do.MustInvoke[*accounts.Registrar](i),
				do.MustInvoke[*events.Publisher](i),
				do.MustInvoke[*localization.Translator](i),
				do.MustInvoke[*handlers.APIErrorHandler](i),
			), nil
		})

		do.Provide(injector, func(i do.Injector) (*api.Server, error) {
			// WithErrorHandler installs the central ent-error -> HTTP-status
			// mapping (404/409), so handlers return raw ent errors instead of
			// typed DTOs.
			return api.NewServer(
				do.MustInvoke[*handlers.Server](i),
				api.WithErrorHandler(do.MustInvoke[*handlers.APIErrorHandler](i).Write),
				api.WithTracerProvider(do.MustInvoke[*otelconf.SDK](i).TracerProvider()),
				api.WithNotFound(handlers.NewNotFoundHandler(do.MustInvoke[*localization.Translator](i))),
				api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(do.MustInvoke[*localization.Translator](i))),
			)
		})

		do.Provide(injector, func(i do.Injector) (*handlers.AttachmentHandler, error) {
			return handlers.NewAttachmentHandler(
				do.MustInvoke[*assetstore.Store](i),
				do.MustInvoke[*localization.Translator](i),
				do.MustInvoke[*handlers.APIErrorHandler](i),
			), nil
		})

		do.Provide(injector, func(i do.Injector) (*handlers.GitHubWebhookHandler, error) {
			return handlers.NewGitHubWebhookHandler(
				do.MustInvoke[*ent.Client](i),
				do.MustInvoke[*versionbuilds.Starter](i),
				do.MustInvoke[*config.Config](i).GitHubWebhookSecret,
				do.MustInvoke[*localization.Translator](i),
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
				do.MustInvoke[*handlers.Server](i).AuthHandler(),
			), nil
		})

		do.ProvideNamed(injector, httpHandlerServiceName, func(i do.Injector) (http.Handler, error) {
			// Dev CORS lets the Vite frontend (on any localhost port) call both
			// the generated API and the hand-mounted routes.
			localized := do.MustInvoke[*localization.Translator](i).Middleware(
				do.MustInvokeNamed[http.Handler](i, routerServiceName),
			)
			return cors.New(cors.Options{
				AllowedOrigins: []string{"http://localhost:*", "http://127.0.0.1:*"},
				AllowedMethods: []string{
					http.MethodGet,
					http.MethodHead,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
					http.MethodOptions,
				},
				AllowedHeaders:   []string{"Accept", "Content-Type", "X-Requested-With", "X-XSRF-TOKEN"},
				AllowCredentials: true,
			}).Handler(localized), nil
		})

		// The process lifecycle coordinator starts and gracefully stops this
		// server. Keeping the provider on the vendor type avoids coupling DI to
		// supervision.
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
	}

	return injector
}

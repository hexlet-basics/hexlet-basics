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

// ServerDependencies is the process-level root of the synchronous application
// graph. cmd/server owns the lifecycle of every resource exposed here.
type ServerDependencies struct {
	Logger           *slog.Logger   `do:""`
	SentryClient     *sentry.Client `do:""`
	HTTPServer       *http.Server   `do:""`
	Database         *ent.Client    `do:""`
	Bucket           *blob.Bucket   `do:""`
	OpenTelemetrySDK *otelconf.SDK  `do:""`
}

// WorkerDependencies is the process-level root of the asynchronous application
// graph. cmd/worker owns the lifecycle of every resource exposed here.
type WorkerDependencies struct {
	Logger           *slog.Logger           `do:""`
	SentryClient     *sentry.Client         `do:""`
	Database         *ent.Client            `do:""`
	RiverClient      *river.Client[*sql.Tx] `do:""`
	EventRuntime     *events.Runtime        `do:""`
	Bucket           *blob.Bucket           `do:""`
	OpenTelemetrySDK *otelconf.SDK          `do:""`
}

// BuildServer resolves the complete synchronous application graph. Its River
// client can enqueue jobs but has no queues or workers, so this process cannot
// execute background work.
func BuildServer() (ServerDependencies, error) {
	return do.InvokeStruct[ServerDependencies](newContainer(false))
}

// BuildWorker resolves the complete asynchronous application graph. It owns
// the Watermill subscribers and River workers but exposes no HTTP server.
func BuildWorker() (WorkerDependencies, error) {
	return do.InvokeStruct[WorkerDependencies](newContainer(true))
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
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return telemetry.NewSentryClient(cfg)
	})

	do.Provide(injector, func(i do.Injector) (*otelconf.SDK, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return telemetry.NewOpenTelemetrySDK(context.Background(), logger)
	})

	do.Provide(injector, func(i do.Injector) (*handlers.APIErrorHandler, error) {
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		sentryClient, err := do.Invoke[*sentry.Client](i)
		if err != nil {
			return nil, err
		}
		return handlers.NewAPIErrorHandler(translator, logger, sentryClient), nil
	})

	do.Provide(injector, func(do.Injector) (*localization.Translator, error) {
		return localization.New()
	})

	do.Provide(injector, func(i do.Injector) (*sql.DB, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return store.NewDB(cfg.DatabaseURL)
	})

	do.Provide(injector, func(i do.Injector) (*ent.Client, error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		return store.NewClient(db), nil
	})

	do.Provide(injector, func(i do.Injector) (*store.Store, error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		return store.New(db), nil
	})

	if !worker {
		do.Provide(injector, func(i do.Injector) (*events.Publisher, error) {
			db, err := do.Invoke[*store.Store](i)
			if err != nil {
				return nil, err
			}
			logger, err := do.Invoke[*slog.Logger](i)
			if err != nil {
				return nil, err
			}
			return events.NewPublisher(db, logger), nil
		})

		do.Provide(injector, func(i do.Injector) (*accounts.Registrar, error) {
			db, err := do.Invoke[*store.Store](i)
			if err != nil {
				return nil, err
			}
			publisher, err := do.Invoke[*events.Publisher](i)
			if err != nil {
				return nil, err
			}
			return accounts.NewRegistrar(db, publisher), nil
		})
	}

	// Blob bucket for uploaded assets (ADR-0005). Closed explicitly in main.go on
	// shutdown, like the ent client and pgx pool (blob.Bucket has Close, not a do
	// Shutdowner).
	do.Provide(injector, func(i do.Injector) (*blob.Bucket, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return store.NewBucket(context.Background(), cfg.BlobBucketURL)
	})

	// Asset storage owns MIME policy, blob writes, attachment persistence,
	// compensation, and public URLs for every upload path.
	do.Provide(injector, func(i do.Injector) (*assetstore.Store, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		bucket, err := do.Invoke[*blob.Bucket](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return assetstore.New(db, bucket, cfg.PublicURL), nil
	})

	if worker {
		// The exercise loader (course-version builds) uses the shared asset store
		// for lesson theory images and a git fetcher for the source repository.
		do.Provide(injector, func(i do.Injector) (courseloader.Fetcher, error) {
			cfg, err := do.Invoke[*config.Config](i)
			if err != nil {
				return nil, err
			}
			return courseloader.NewGitFetcher(cfg.CourseRepoBaseURL, cfg.GitHubToken), nil
		})

		do.Provide(injector, func(i do.Injector) (*courseloader.Loader, error) {
			db, err := do.Invoke[*ent.Client](i)
			if err != nil {
				return nil, err
			}
			store, err := do.Invoke[*store.Store](i)
			if err != nil {
				return nil, err
			}
			assets, err := do.Invoke[*assetstore.Store](i)
			if err != nil {
				return nil, err
			}
			fetcher, err := do.Invoke[courseloader.Fetcher](i)
			if err != nil {
				return nil, err
			}
			return courseloader.NewLoader(db, store, assets, fetcher), nil
		})

		do.Provide(injector, func(i do.Injector) (*amocrm.Client, error) {
			cfg, err := do.Invoke[*config.Config](i)
			if err != nil {
				return nil, err
			}
			return amocrm.NewClient(cfg.AmoCRMBaseURL, cfg.AmoCRMAuthToken, cfg.YMCounter), nil
		})
	}

	do.Provide(injector, func(i do.Injector) (*river.Client[*sql.Tx], error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		otelSDK, err := do.Invoke[*otelconf.SDK](i)
		if err != nil {
			return nil, err
		}

		if worker {
			loader, err := do.Invoke[*courseloader.Loader](i)
			if err != nil {
				return nil, err
			}
			amoCRMClient, err := do.Invoke[*amocrm.Client](i)
			if err != nil {
				return nil, err
			}
			sentryClient, err := do.Invoke[*sentry.Client](i)
			if err != nil {
				return nil, err
			}
			return jobs.NewWorkerClient(
				db,
				loader,
				amoCRMClient,
				logger,
				jobs.NewErrorHandler(sentryClient),
				otelSDK.TracerProvider(),
				otelSDK.MeterProvider(),
			)
		}
		return jobs.NewInsertOnlyClient(
			db,
			logger,
			otelSDK.TracerProvider(),
			otelSDK.MeterProvider(),
		)
	})

	if worker {
		do.Provide(injector, func(i do.Injector) (*events.Runtime, error) {
			db, err := do.Invoke[*sql.DB](i)
			if err != nil {
				return nil, err
			}
			logger, err := do.Invoke[*slog.Logger](i)
			if err != nil {
				return nil, err
			}
			riverClient, err := do.Invoke[*river.Client[*sql.Tx]](i)
			if err != nil {
				return nil, err
			}
			return events.NewRuntime(db, logger, eventhandlers.LeadCreated(riverClient))
		})
	} else {
		do.Provide(injector, func(i do.Injector) (*versionbuilds.Starter, error) {
			db, err := do.Invoke[*store.Store](i)
			if err != nil {
				return nil, err
			}
			riverClient, err := do.Invoke[*river.Client[*sql.Tx]](i)
			if err != nil {
				return nil, err
			}
			return versionbuilds.NewStarter(db, riverClient), nil
		})

		do.Provide(injector, func(i do.Injector) (*handlers.Server, error) {
			db, err := do.Invoke[*ent.Client](i)
			if err != nil {
				return nil, err
			}
			cfg, err := do.Invoke[*config.Config](i)
			if err != nil {
				return nil, err
			}
			starter, err := do.Invoke[*versionbuilds.Starter](i)
			if err != nil {
				return nil, err
			}
			registrar, err := do.Invoke[*accounts.Registrar](i)
			if err != nil {
				return nil, err
			}
			publisher, err := do.Invoke[*events.Publisher](i)
			if err != nil {
				return nil, err
			}
			translator, err := do.Invoke[*localization.Translator](i)
			if err != nil {
				return nil, err
			}
			errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
			if err != nil {
				return nil, err
			}
			return handlers.NewServer(
				db,
				cfg,
				starter,
				registrar,
				publisher,
				translator,
				errorHandler,
			), nil
		})

		do.Provide(injector, func(i do.Injector) (*api.Server, error) {
			handler, err := do.Invoke[*handlers.Server](i)
			if err != nil {
				return nil, err
			}
			errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
			if err != nil {
				return nil, err
			}
			otelSDK, err := do.Invoke[*otelconf.SDK](i)
			if err != nil {
				return nil, err
			}
			translator, err := do.Invoke[*localization.Translator](i)
			if err != nil {
				return nil, err
			}
			// WithErrorHandler installs the central ent-error -> HTTP-status
			// mapping (404/409), so handlers return raw ent errors instead of
			// typed DTOs.
			return api.NewServer(
				handler,
				api.WithErrorHandler(errorHandler.Write),
				api.WithTracerProvider(otelSDK.TracerProvider()),
				api.WithNotFound(handlers.NewNotFoundHandler(translator)),
				api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
			)
		})

		do.Provide(injector, func(i do.Injector) (*handlers.AttachmentHandler, error) {
			assets, err := do.Invoke[*assetstore.Store](i)
			if err != nil {
				return nil, err
			}
			translator, err := do.Invoke[*localization.Translator](i)
			if err != nil {
				return nil, err
			}
			errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
			if err != nil {
				return nil, err
			}
			return handlers.NewAttachmentHandler(assets, translator, errorHandler), nil
		})

		do.Provide(injector, func(i do.Injector) (*handlers.GitHubWebhookHandler, error) {
			db, err := do.Invoke[*ent.Client](i)
			if err != nil {
				return nil, err
			}
			starter, err := do.Invoke[*versionbuilds.Starter](i)
			if err != nil {
				return nil, err
			}
			cfg, err := do.Invoke[*config.Config](i)
			if err != nil {
				return nil, err
			}
			translator, err := do.Invoke[*localization.Translator](i)
			if err != nil {
				return nil, err
			}
			return handlers.NewGitHubWebhookHandler(db, starter, cfg.GitHubWebhookSecret, translator), nil
		})

		// Both the raw router and the middleware-wrapped application handler have
		// the same http.Handler interface. Named bindings keep those two seams
		// explicit without adding wrapper types whose only purpose would be DI.
		do.ProvideNamed(injector, routerServiceName, func(i do.Injector) (http.Handler, error) {
			apiServer, err := do.Invoke[*api.Server](i)
			if err != nil {
				return nil, err
			}
			attachmentHandler, err := do.Invoke[*handlers.AttachmentHandler](i)
			if err != nil {
				return nil, err
			}
			webhookHandler, err := do.Invoke[*handlers.GitHubWebhookHandler](i)
			if err != nil {
				return nil, err
			}
			handler, err := do.Invoke[*handlers.Server](i)
			if err != nil {
				return nil, err
			}
			return handlers.NewRouter(
				apiServer,
				attachmentHandler,
				webhookHandler,
				handler.AuthHandler(),
			), nil
		})

		do.ProvideNamed(injector, httpHandlerServiceName, func(i do.Injector) (http.Handler, error) {
			translator, err := do.Invoke[*localization.Translator](i)
			if err != nil {
				return nil, err
			}
			router, err := do.InvokeNamed[http.Handler](i, routerServiceName)
			if err != nil {
				return nil, err
			}
			// Dev CORS lets the Vite frontend (on any localhost port) call both
			// the generated API and the hand-mounted routes.
			localized := translator.Middleware(router)
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
			cfg, err := do.Invoke[*config.Config](i)
			if err != nil {
				return nil, err
			}
			handler, err := do.InvokeNamed[http.Handler](i, httpHandlerServiceName)
			if err != nil {
				return nil, err
			}

			return &http.Server{
				Addr:              cfg.Addr,
				Handler:           handler,
				ReadHeaderTimeout: readHeaderTimeout,
				ReadTimeout:       readTimeout,
				WriteTimeout:      writeTimeout,
				IdleTimeout:       idleTimeout,
			}, nil
		})
	}

	return injector
}

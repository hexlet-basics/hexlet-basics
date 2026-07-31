package di

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/rs/cors"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/config"
	"hexletbasics/internal/events"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/lessonreviews"
	"hexletbasics/internal/localization"
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

// serverPackage contains the synchronous graph: HTTP delivery, event
// publishing, and insert-only River access.
var serverPackage = do.Package(
	do.Lazy[*handlers.APIErrorHandler](func(i do.Injector) (*handlers.APIErrorHandler, error) {
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return handlers.NewAPIErrorHandler(translator, logger), nil
	}),
	do.Lazy[*events.Publisher](func(i do.Injector) (*events.Publisher, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return events.NewPublisher(db, logger), nil
	}),
	do.Lazy[*accounts.Registrar](func(i do.Injector) (*accounts.Registrar, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return accounts.NewRegistrar(db, publisher), nil
	}),
	do.Lazy[*river.Client[*sql.Tx]](func(i do.Injector) (*river.Client[*sql.Tx], error) {
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
		return jobs.NewInsertOnlyClient(
			db,
			logger,
			otelSDK.TracerProvider(),
			otelSDK.MeterProvider(),
		)
	}),
	do.Lazy[*versionbuilds.Starter](func(i do.Injector) (*versionbuilds.Starter, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		riverClient, err := do.Invoke[*river.Client[*sql.Tx]](i)
		if err != nil {
			return nil, err
		}
		return versionbuilds.NewStarter(db, riverClient), nil
	}),
	do.Lazy[*lessonreviews.Enqueuer](func(i do.Injector) (*lessonreviews.Enqueuer, error) {
		riverClient, err := do.Invoke[*river.Client[*sql.Tx]](i)
		if err != nil {
			return nil, err
		}
		return lessonreviews.NewEnqueuer(riverClient), nil
	}),
	do.Lazy[*handlers.Server](func(i do.Injector) (*handlers.Server, error) {
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
		reviews, err := do.Invoke[*lessonreviews.Enqueuer](i)
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
			reviews,
			registrar,
			publisher,
			translator,
			errorHandler,
		), nil
	}),
	do.Lazy[*api.Server](func(i do.Injector) (*api.Server, error) {
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
			handler.AuthHandler(),
			api.WithErrorHandler(errorHandler.Write),
			api.WithTracerProvider(otelSDK.TracerProvider()),
			api.WithNotFound(handlers.NewNotFoundHandler(translator)),
			api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
		)
	}),
	do.Lazy[*handlers.AttachmentHandler](func(i do.Injector) (*handlers.AttachmentHandler, error) {
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
	}),
	do.Lazy[*handlers.GitHubWebhookHandler](func(i do.Injector) (*handlers.GitHubWebhookHandler, error) {
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
	}),
	// Both the raw router and the middleware-wrapped application handler have
	// the same http.Handler interface. Named bindings keep those two seams
	// explicit without adding wrapper types whose only purpose would be DI.
	do.LazyNamed[http.Handler](routerServiceName, func(i do.Injector) (http.Handler, error) {
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
	}),
	do.LazyNamed[http.Handler](httpHandlerServiceName, func(i do.Injector) (http.Handler, error) {
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		router, err := do.InvokeNamed[http.Handler](i, routerServiceName)
		if err != nil {
			return nil, err
		}
		sentryClient, err := do.Invoke[*sentry.Client](i)
		if err != nil {
			return nil, err
		}
		// Dev CORS lets the Vite frontend (on any localhost port) call both
		// the generated API and the hand-mounted routes.
		localized := translator.Middleware(router)
		corsHandler := cors.New(cors.Options{
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
		}).Handler(localized)
		return telemetry.NewSentryHTTPHandler(sentryClient, corsHandler), nil
	}),
	// The process lifecycle coordinator starts and gracefully stops this
	// server. Keeping the provider on the vendor type avoids coupling DI to
	// supervision.
	do.Lazy[*http.Server](func(i do.Injector) (*http.Server, error) {
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
	}),
)

package di

import (
	"database/sql"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"

	"hexletbasics/ent"
	"hexletbasics/internal/amocrm"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/config"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
)

// workerPackage contains the asynchronous graph: course loading, external
// integrations, River workers, and Watermill subscribers.
var workerPackage = do.Package(
	// The exercise loader uses the shared asset store for lesson theory images
	// and a git fetcher for the source repository.
	do.Lazy[courseloader.Fetcher](func(i do.Injector) (courseloader.Fetcher, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return courseloader.NewGitFetcher(cfg.CourseRepoBaseURL, cfg.GitHubToken), nil
	}),
	do.Lazy[*courseloader.Loader](func(i do.Injector) (*courseloader.Loader, error) {
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
	}),
	do.Lazy[*amocrm.Client](func(i do.Injector) (*amocrm.Client, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return amocrm.NewClient(cfg.AmoCRMBaseURL, cfg.AmoCRMAuthToken, cfg.YMCounter), nil
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
	}),
	do.Lazy[*events.Runtime](func(i do.Injector) (*events.Runtime, error) {
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
	}),
)

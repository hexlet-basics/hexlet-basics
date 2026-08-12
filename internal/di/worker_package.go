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
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/config"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/lessonreviews"
	"hexletbasics/internal/progress"
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
		completion, err := do.Invoke[*progress.Progress](i)
		if err != nil {
			return nil, err
		}
		return courseloader.NewLoader(db, store, assets, fetcher, completion), nil
	}),
	// Promoting a version re-evaluates completion, so the worker graph carries
	// the progress module too. It also needs a publisher of its own: the worker
	// otherwise only consumes facts, and promotion is the one place it raises
	// them.
	do.Lazy[*events.Publisher](func(i do.Injector) (*events.Publisher, error) {
		txStore, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return events.NewPublisher(txStore, logger), nil
	}),
	do.Lazy[*progress.Progress](func(i do.Injector) (*progress.Progress, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		txStore, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		// No exercise runner in this graph: the worker reaches the progress
		// module only through promotion, which re-evaluates completion and never
		// runs a submission. Wiring a real runner here would suggest otherwise.
		return progress.New(db, txStore, publisher, progress.UnavailableRunner{}), nil
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
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		entClient, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		// Without OpenAI credentials the reviewer stays unregistered: enqueued
		// review jobs wait in the queue instead of failing against a dead client.
		var reviewer jobs.LessonReviewer
		if cfg.OpenAIAccessToken != "" {
			reviewer = lessonreviews.NewReviewer(
				entClient,
				assistant.NewOpenAI(cfg.OpenAIAccessToken, cfg.OpenAIModel),
			)
		}
		return jobs.NewWorkerClient(
			db,
			loader,
			amoCRMClient,
			reviewer,
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

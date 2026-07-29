package di

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"
	"gocloud.dev/blob"

	"hexletbasics/ent"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/config"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/logging"
	"hexletbasics/internal/store"
	"hexletbasics/internal/telemetry"
)

// commonPackage contains infrastructure shared by the HTTP and worker
// processes. Process-specific packages build on this graph.
var commonPackage = do.Package(
	do.Lazy[*config.Config](func(do.Injector) (*config.Config, error) {
		return config.Load()
	}),
	do.Lazy[*slog.Logger](func(do.Injector) (*slog.Logger, error) {
		return logging.New(slog.LevelInfo), nil
	}),
	do.Lazy[*sentry.Client](func(i do.Injector) (*sentry.Client, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return telemetry.NewSentryClient(cfg)
	}),
	do.Lazy[*otelconf.SDK](func(i do.Injector) (*otelconf.SDK, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return telemetry.NewOpenTelemetrySDK(context.Background(), logger)
	}),
	do.Lazy[*localization.Translator](func(do.Injector) (*localization.Translator, error) {
		return localization.New()
	}),
	do.Lazy[*sql.DB](func(i do.Injector) (*sql.DB, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return store.NewDB(cfg.DatabaseURL)
	}),
	do.Lazy[*ent.Client](func(i do.Injector) (*ent.Client, error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		return store.NewClient(db), nil
	}),
	do.Lazy[*store.Store](func(i do.Injector) (*store.Store, error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		return store.New(db), nil
	}),
	// Blob buckets are closed explicitly by the process lifecycle coordinator
	// because blob.Bucket has Close, not a do Shutdowner.
	do.Lazy[*blob.Bucket](func(i do.Injector) (*blob.Bucket, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return store.NewBucket(context.Background(), cfg.BlobBucketURL)
	}),
	// Asset storage owns MIME policy, blob writes, attachment persistence,
	// compensation, and public URLs for every upload path.
	do.Lazy[*assetstore.Store](func(i do.Injector) (*assetstore.Store, error) {
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
	}),
)

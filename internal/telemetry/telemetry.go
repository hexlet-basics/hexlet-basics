// Package telemetry initializes the native Sentry and OpenTelemetry SDKs.
//
// It deliberately returns vendor types: ogen, River, DI, and tests use the
// extension points those libraries already expose instead of learning an
// application-specific observability interface.
package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/contrib/otelconf"
	"go.opentelemetry.io/otel"

	"hexletbasics/internal/config"
)

// NewSentryClient builds a process-wide client. sentry-go treats an empty DSN
// as disabled; an explicitly malformed DSN is returned as a startup error.
func NewSentryClient(dsn string, environment config.AppEnv, release string) (*sentry.Client, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: string(environment),
		Release:     release,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize Sentry: %w", err)
	}
	return client, nil
}

// NewOpenTelemetrySDK builds the declaratively configured OpenTelemetry SDK and
// installs all of its providers globally. Without OTEL_CONFIG_FILE, otelconf
// returns noop providers, preserving zero-config local development.
func NewOpenTelemetrySDK(
	ctx context.Context,
	logger *slog.Logger,
) (*otelconf.SDK, error) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.ErrorContext(context.WithoutCancel(ctx), "OpenTelemetry error", "err", err)
	}))
	sdk, err := otelconf.NewSDK(otelconf.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("initialize OpenTelemetry SDK: %w", err)
	}

	otel.SetTracerProvider(sdk.TracerProvider())
	otel.SetMeterProvider(sdk.MeterProvider())
	otel.SetLoggerProvider(sdk.LoggerProvider())
	otel.SetTextMapPropagator(sdk.Propagator())
	return &sdk, nil
}

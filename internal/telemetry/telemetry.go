// Package telemetry initializes the native Sentry and OpenTelemetry SDKs.
//
// It deliberately returns vendor types: ogen, River, DI, and tests use the
// extension points those libraries already expose instead of learning an
// application-specific observability interface.
package telemetry

import (
	"context"
	"errors"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/samber/oops"
	"go.opentelemetry.io/contrib/otelconf"
	"go.opentelemetry.io/otel"

	"hexletbasics/internal/config"
)

// SentryOption adjusts the client options NewSentryClient builds; tests use it
// to swap the transport.
type SentryOption func(*sentry.ClientOptions)

// WithSentryTransport replaces the network transport.
func WithSentryTransport(transport sentry.Transport) SentryOption {
	return func(o *sentry.ClientOptions) { o.Transport = transport }
}

// NewSentryClient builds a process-wide client. sentry-go treats an empty DSN
// as disabled; an explicitly malformed DSN is returned as a startup error.
//
// BeforeSend copies oops attributes into every event (ADR-0018), so the HTTP
// and River paths, which both report through this client, carry them.
func NewSentryClient(dsn string, environment config.AppEnv, release string, opts ...SentryOption) (*sentry.Client, error) {
	options := sentry.ClientOptions{
		Dsn:         dsn,
		Environment: string(environment),
		Release:     release,
		BeforeSend:  attachOopsAttributes,
	}
	for _, opt := range opts {
		opt(&options)
	}
	client, err := sentry.NewClient(options)
	if err != nil {
		return nil, errs.Wrapf(err, "initialize Sentry")
	}
	return client, nil
}

// attachOopsAttributes exposes what sentry-go cannot see: it reads the stack of
// an oops error by reflection but not its attributes. The event gets the
// package as the oops.domain tag and every With(...) pair, from each wrapped
// layer, in the "oops" context.
func attachOopsAttributes(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	if hint == nil || hint.OriginalException == nil {
		return event
	}
	oopsErr, ok := errors.AsType[oops.OopsError](hint.OriginalException)
	if !ok {
		return event
	}
	if domain := oopsErr.Domain(); domain != "" {
		if event.Tags == nil {
			event.Tags = map[string]string{}
		}
		event.Tags["oops.domain"] = domain
	}
	if attrs := oopsErr.Context(); len(attrs) > 0 {
		if event.Contexts == nil {
			event.Contexts = map[string]sentry.Context{}
		}
		event.Contexts["oops"] = attrs
	}
	return event
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
		return nil, errs.Wrapf(err, "initialize OpenTelemetry SDK")
	}

	otel.SetTracerProvider(sdk.TracerProvider())
	otel.SetMeterProvider(sdk.MeterProvider())
	otel.SetLoggerProvider(sdk.LoggerProvider())
	otel.SetTextMapPropagator(sdk.Propagator())
	return &sdk, nil
}

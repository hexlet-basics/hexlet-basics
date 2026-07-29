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
	"net/url"
	"os"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"

	"hexletbasics/internal/config"
)

const (
	defaultServiceName = "hexlet-basics-api"
	defaultSampleRatio = 0.01
)

// NewSentryClient builds a process-wide client. sentry-go treats an empty DSN
// as disabled; an explicitly malformed DSN is returned as a startup error.
func NewSentryClient(cfg *config.Config) (*sentry.Client, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         cfg.SentryDSN,
		Environment: cfg.SentryEnvironment,
		Release:     cfg.ReleaseVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize Sentry: %w", err)
	}
	return client, nil
}

// NewTracerProvider builds the provider used explicitly by ogen and installs it
// globally for other instrumentation. Without an endpoint it returns a
// non-recording SDK provider, preserving zero-config local development.
func NewTracerProvider(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
) (*sdktrace.TracerProvider, error) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.ErrorContext(context.Background(), "OpenTelemetry error", "err", err)
	}))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if cfg.OTLPEndpoint == "" && cfg.OTLPTracesEndpoint == "" {
		provider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.NeverSample()))
		otel.SetTracerProvider(provider)
		return provider, nil
	}
	if err := validateEndpoint("OTEL_EXPORTER_OTLP_ENDPOINT", cfg.OTLPEndpoint); err != nil {
		return nil, err
	}
	if err := validateEndpoint("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", cfg.OTLPTracesEndpoint); err != nil {
		return nil, err
	}

	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize OTLP HTTP trace exporter: %w", err)
	}

	res, err := telemetryResource(cfg)
	if err != nil {
		return nil, fmt.Errorf("build OpenTelemetry resource: %w", err)
	}
	options := []sdktrace.TracerProviderOption{
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	}
	// The SDK parses OTEL_TRACES_SAMPLER* itself. Only supply the application
	// default when the operator did not choose a standard sampler explicitly.
	if _, configured := os.LookupEnv("OTEL_TRACES_SAMPLER"); !configured {
		options = append(options, sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(defaultSampleRatio)),
		))
	}

	provider := sdktrace.NewTracerProvider(options...)
	otel.SetTracerProvider(provider)
	return provider, nil
}

func telemetryResource(cfg *config.Config) (*resource.Resource, error) {
	attrs := make([]attribute.KeyValue, 0, 2)
	if _, configured := os.LookupEnv("OTEL_SERVICE_NAME"); !configured {
		attrs = append(attrs, semconv.ServiceName(defaultServiceName))
	}
	if cfg.ReleaseVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ReleaseVersion))
	}
	if len(attrs) == 0 {
		return resource.Default(), nil
	}
	return resource.Merge(resource.Default(), resource.NewSchemaless(attrs...))
}

func validateEndpoint(name, value string) error {
	if value == "" {
		return nil
	}
	endpoint, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", name, err)
	}
	if endpoint.Scheme != "http" && endpoint.Scheme != "https" {
		return fmt.Errorf("parse %s: endpoint scheme must be http or https", name)
	}
	if endpoint.Host == "" {
		return fmt.Errorf("parse %s: endpoint host is required", name)
	}
	return nil
}

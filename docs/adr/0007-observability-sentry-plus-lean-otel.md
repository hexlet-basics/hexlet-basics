# Observability: Sentry for errors, OTLP tracing, Prometheus deferred

Error tracking stays on **Sentry** (self-hosted, as in legacy) via the
`getsentry/sentry-go` SDK. Logs use `slog` + `lmittmann/tint`. OpenTelemetry
**traces** use ogen's native spans. The OpenTelemetry SDK is built by the
official contrib `otelconf` package from the declarative file selected by
`OTEL_CONFIG_FILE`; without that variable all providers are noop, so local
development remains zero-config and non-recording.

Production defaults to parent-based 1% sampling, matching the order of
magnitude used by the legacy application. Exporter, resource, propagator, and
sampler policy live together in the declarative configuration instead of an
application-specific environment parser.

OTel traces/metrics do NOT replace Sentry's exception aggregation; the two are
complementary, so both are kept.

## Consequences

- `getsentry/sentry-go`, the OTel SDK, and contrib `otelconf` are runtime
  dependencies.
- ogen's generated HTTP spans are exported without an additional HTTP tracing
  middleware.
- The application installs tracer, meter, and logger providers from one SDK;
  metrics and OTel logs remain noop until their pipelines are added to the
  declarative file.
- Unexpected HTTP failures and River job failures are captured by Sentry at
  their native error-handler seams; expected HTTP 4xx responses are not.
- Prometheus and an OTel metrics exporter remain deferred until a metrics
  receiver and operational requirements are confirmed.

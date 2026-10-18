# Observability: Sentry for errors, OTLP tracing, Prometheus deferred

> **Extended by [ADR-0026](0026-logging-conventions.md).** The logging part (slog plus
> tint) stays; ADR-0026 adds JSON off a terminal, `trace_id` correlation,
> `LOG_LEVEL` and redaction. Error reporting to Sentry also follows
> [ADR-0018](0018-errors-via-oops.md) (`oops` attributes) and
> [ADR-0027](0027-worker-conventions-river-and-watermill.md) (jobs report on the final attempt).

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
- The official `sentryhttp` middleware wraps the complete HTTP application. It
  gives every request an isolated hub, captures panics, and repanics so
  `net/http` retains its normal failure semantics.
- The application installs tracer, meter, and logger providers from one SDK;
  metrics and OTel logs remain noop until their pipelines are added to the
  declarative file.
- Unexpected returned HTTP failures are captured through the request-local hub
  at ogen's error-handler seams. River failures stay at River's native error
  seam, and expected HTTP 4xx responses are not captured.
- Prometheus and an OTel metrics exporter remain deferred until a metrics
  receiver and operational requirements are confirmed.

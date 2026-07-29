# Observability: Sentry for errors, OTLP tracing, Prometheus deferred

Error tracking stays on **Sentry** (self-hosted, as in legacy) via the
`getsentry/sentry-go` SDK. Logs use `slog` + `lmittmann/tint`. OpenTelemetry
**traces** use ogen's native spans and the OTLP/HTTP protobuf exporter. The
exporter is enabled only when a standard `OTEL_EXPORTER_OTLP_*_ENDPOINT` is
configured; local development remains zero-config and non-recording.

Production defaults to parent-based 1% sampling, matching the order of
magnitude used by the legacy application. Operators can override it through
the standard `OTEL_TRACES_SAMPLER*` variables.

OTel traces/metrics do NOT replace Sentry's exception aggregation; the two are
complementary, so both are kept.

## Consequences

- `getsentry/sentry-go`, the OTel SDK, and the OTLP/HTTP trace exporter are
  runtime dependencies.
- ogen's generated HTTP spans are exported without an additional HTTP tracing
  middleware.
- Unexpected HTTP failures and River job failures are captured by Sentry at
  their native error-handler seams; expected HTTP 4xx responses are not.
- Prometheus and an OTel metrics exporter remain deferred until a metrics
  receiver and operational requirements are confirmed.

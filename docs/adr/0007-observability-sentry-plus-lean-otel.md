# Observability: Sentry for errors, lean OpenTelemetry, Prometheus deferred

Error tracking stays on **Sentry** (self-hosted, as in legacy) via the
`getsentry/sentry-go` SDK — added because it was missing from the reference lib
set. Logs use `slog` + `lmittmann/tint`. OpenTelemetry **traces** are adopted
cheaply since `ogen` emits them natively, but the heavier OTLP exporter +
Prometheus metrics wiring is **deferred** until a metrics/trace backend receiver
is confirmed — we don't ship exporters into a void.

OTel traces/metrics do NOT replace Sentry's exception aggregation; the two are
complementary, so both are kept.

## Consequences

- `getsentry/sentry-go` is added to the dependency set.
- The full OTel exporter/Prometheus block from the reference list is not part of
  the initial set; revisit when a collector/Grafana/Prometheus receiver exists.

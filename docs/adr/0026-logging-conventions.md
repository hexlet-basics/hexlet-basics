---
status: accepted
---

# Logging: JSON off a terminal, `trace_id` correlation, `LOG_LEVEL`, redaction

Part of the backend clean-up (spec #905, ticket #920, decisions in #889).
Extends [ADR-0007](0007-observability-sentry-otel.md), which keeps `slog` plus
`tint` as the logger and decides sampling; it does not replace it.

**Problem.** Log lines had no link to the trace of the request that wrote them
(only `report()` added ids, by hand), the level was a constant, output was text
even in a container, and nothing stopped a secret or an email from landing in a
log line.

**Decision.**

- **Format.** `internal/logging` writes tint text when stderr is a terminal (or
  `FORCE_COLOR` says a human reads the file) and `slog.JSONHandler` otherwise.
  `NO_COLOR` only turns colour off. There is no `LOG_FORMAT` variable.
- **Context.** A call with a `ctx` in scope uses the `*Context` variant;
  `sloglint` (`context: scope`) enforces it. Calls with no context (`main`
  before the signal context exists) stay plain.
- **Trace correlation.** One wrapper handler adds `trace_id` and `span_id` when
  `SpanContext().IsValid()`. It does not test `IsRecording()`: production keeps
  1% of traces (ADR-0007), and that test would leave 99% of lines without the key
  that finds them. `trace_id` is the only correlation key; there is no request
  id of our own, and `X-Request-ID` stays an audit-only field from the proxy.
  `report()` no longer adds the ids to its log call (the Sentry tags stay).
  Exporting `slog` to OTLP stays off until a log pipeline exists.
- **Attributes.** `snake_case` keys, `err` for an error (never `error`), static
  messages with the variable parts as attributes, and one call style
  (key-value pairs, no mixing with `slog.Attr`), all enforced by `sloglint`.
- **`report()`** logs the error, whose `oops` attributes the handler expands,
  plus `status` and `method`/`path`. A 4xx is neither logged nor sent.
- **Level.** `LOG_LEVEL` (`debug|info|warn|error`, default `info`) in
  `config.Common`; any other value stops the start.
- **Sensitive data.** No passwords, tokens, `Authorization`/`Cookie` values or
  emails in logs; identity is `user_id`. A `ReplaceAttr` safety net replaces the
  value of any attribute whose key contains `password`, `passwd`, `secret`,
  `token`, `authorization`, `cookie`, `email`, `api_key` or `apikey`. It is a net
  for a slip, not the standard: the standard is review. The development mailer
  (`mailer.Log`, used only where no Postbox key is set, which production
  refuses) prints the recipient and the message so a developer can follow the
  link; it is never active in production.

**Rejected.** The samber slog libraries, checked against their repositories:
`slog-otel` (v0.1.0, an internal helper) gates on `IsRecording()`;
`slog-formatter` masks by keeping a value prefix, and its full-redaction
`FormatByKey` is as long as our `ReplaceAttr`; `slog-sentry` is a second path
into Sentry beside `report()` and `BeforeSend`, so one failure would arrive
twice. `slog-multi` is considered only if a second sink appears.

## Consequences

- Container logs are one JSON object per line with `trace_id`; a Sentry event
  and its log lines share the id.
- A key containing one of the fragments is redacted whatever its value, so an
  innocent attribute (a count of tokens) must be named without the fragment.
- `sloglint` joins the lint set; a new `Info`/`Error` call with a `ctx` in
  scope fails lint until it takes the `*Context` form.

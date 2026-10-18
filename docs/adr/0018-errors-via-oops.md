---
status: accepted
---

# Errors are built with samber/oops, one template per package

Every error in production code is created and wrapped with
[`samber/oops`](https://github.com/samber/oops); `fmt.Errorf` is banned by
`forbidigo`. Each package declares `var errs = oops.In("<package>")` and builds
errors from it, with the variable data in `.With(...)` and a static message.

Before this, the two styles lived side by side with no rule: 94 `oops.Wrapf` and
`oops.Errorf` calls against 166 `fmt.Errorf`, split by package (`courseloader`
and `assistant` on `oops`, `progress` and `exerciserunner` on `fmt.Errorf`), and
69 of the `oops` calls interpolated data into the message. `oops` was paid for
(stack trace, attributes) but its main benefit, data that travels with the error,
was not used.

**The rules**

- Only `In`, `With`, `Wrapf` and `Errorf`. No `Code`, because nothing maps a code
  to a status here (that is what the `ErrX` values are for). No `Public`:
  responses are localised through the translator, and a single English string
  would fight it. No `Tags`: it duplicates `In`. No `Hint`.
- A message is static; the data goes in `.With("lesson_id", id)`, so errors group
  in Sentry and Loki and the data stays searchable.
- Wrap once, at the package boundary: the exported function that knows the
  operation wraps. Helpers inside a package return the error as is, and when
  there is nothing to add the function returns `err`. Every `Wrap` adds another
  stack trace ([samber/oops#81](https://github.com/samber/oops/issues/81)).
- Log or return, never both. A function that returns an error does not also log
  it: the one place that handles it (the HTTP error handler, the River error
  handler) logs and reports it once, with the attributes the error carries.
- Messages are lowercase and static, with no trailing punctuation, so they
  compose when wrapped (`load lesson: ...`).
- Look into an error chain with `errors.AsType[T](err)` (Go 1.26), not
  `errors.As` with a declared target variable.
- Domain errors (`ErrNotFound`, `ErrLastAdmin`) stay plain `errors.New` values
  checked with `errors.Is`, and `errorStatus` in `handlers/errors.go` keeps
  mapping them.
- `forbidigo` rejects `fmt.Errorf` in production code and allows it in tests, the
  way `.Create()` is rejected in tests: the rule survives people and agents who
  do not remember it. The rule is active now with an exclusion list in
  `.golangci.yml` holding every package not yet migrated; migrating a package
  removes its entry, and the list never grows. `internal/books` is the migrated
  template (`var errs = oops.In("books")`, static messages, `.With`, one wrap at
  the exported function).
- `oops.Recover` guards the one goroutine that does real work without a safety
  net, the assistant stream. HTTP (`sentryhttp` with `Repanic`), Watermill
  (`middleware.Recoverer`) and River (its job executor) already recover.

**What reaches the tools**

`telemetry.NewSentryClient` installs that `BeforeSend` (the `oops.domain` tag and
an `oops` context holding the attributes). It reads the error from
`EventHint.OriginalException`, which `Hub.CaptureException` sets and
`Client.CaptureException` does not, so a caller using the client directly (the
River handler) passes it in the hint.

`slog` receives every attribute through `OopsError.LogValue`. `sentry-go` finds
`OopsError.StackFrames()` by reflection, so the stack reaches Sentry, but the
attributes do not. One `BeforeSend` in `telemetry.NewSentryClient` (about 40
lines) copies them into the event; both error paths, HTTP (`handlers/errors.go`)
and River (`jobs/errors.go`), go through that client. A plain `fmt.Errorf` would
not carry a stack at all: `sentry-go` then reports the stack of the place that
called `CaptureException`.

`sentry-go` shows a duplicate frame and a line one too low for errors that expose
`[]runtime.Frame`, because it resolves the PCs a second time
([getsentry/sentry-go#1461](https://github.com/getsentry/sentry-go/issues/1461),
fix in [#1462](https://github.com/getsentry/sentry-go/pull/1462)). We fix that
upstream rather than rebuild frames ourselves; until the release the duplicate
frame stays.

## Considered Options

- **Keep both styles.** Rejected: it is the state that produced the split, and it
  costs the data that travels with an error.
- **`cockroachdb/errors` or `go-errors/errors`.** Their Sentry support is more
  mature, but both are weaker on attributes and neither gives `slog` a structured
  value; `cockroachdb/errors` also pulls in a large dependency tree. `oops` is
  already a direct dependency used in 17 files.
- **`sentry-go` alone, with `fmt.Errorf`.** Rejected: without an error that
  carries its stack, Sentry shows only where the error was reported, not where
  it was created.

## Consequences

- One mechanical rewrite of the 166 `fmt.Errorf` calls, and of the 69 messages
  that interpolate data, to the template above.
- The Sentry `BeforeSend` adapter and the `forbidigo` rule land with or before
  that rewrite.
- `samber/oops` has no Sentry integration of its own (its README asks for one);
  a helper for it may be contributed upstream later.

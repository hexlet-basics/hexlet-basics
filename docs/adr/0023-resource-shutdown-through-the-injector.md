---
status: accepted
---

# Release resources through injector.Shutdown

`cmd/server` and `cmd/worker` each carried a copy of `closeResources` that
closed the bucket, the database, OpenTelemetry and Sentry by hand, in an order
only a comment explained. A new resource meant editing both copies and both
tests. The injector (ADR-0010 builds the graph with `samber/do` v2) already
knows what it built and what depends on what, so teardown moves there.

## Decision

- After the actors have stopped, each process calls
  `dependencies.Shutdown(ctx)`, which runs `injector.ShutdownWithContext`.
  Both `closeResources` functions are gone.
- Stopping the actors stays staged in the errgroup coordinators of
  ADR-0010: HTTP drain on the server; Watermill, then River on the worker.
  The injector shuts down only after them, because it does not know that
  events enqueue jobs and so must stop first.
- Resources that have no do shutdown method get a small closer service
  (`internal/di/closers.go`): `databaseCloser` (ent client), `bucketCloser`
  (blob bucket), `sentryCloser` (flush). The resource itself stays registered
  under its vendor type, so no other code sees a wrapper. `*otelconf.SDK` and
  `*http.Server` already have `Shutdown(ctx) error`, which do calls itself.
- do shuts a service down only after its dependents. The database and bucket
  closers resolve the Sentry closer in their providers, so Sentry is flushed
  last and sees errors raised while closing.
- A failed build releases whatever was built before returning the error.
- `ShutdownOnSignals` is not used: the signal context belongs to the
  coordinators, and a second signal listener would race the staged stop.

## Consequences

- Adding a resource means registering it and, if it has no shutdown method, a
  closer; no change to `cmd/*`.
- Shutdown errors come back as one report naming every resource that failed.
- Health and readiness endpoints, which depend on the same lifecycle, extend
  this ADR when they land.

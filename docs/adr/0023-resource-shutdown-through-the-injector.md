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

## Health and readiness

- `internal/health` serves `/healthz` and `/readyz` on a plain `http.ServeMux`
  wrapped around the ogen handler (outside the Sentry and locale middleware).
  They are not in the TypeSpec contract, not in the OpenAPI document and not
  under `/api`.
- `/healthz` is 200 whenever the process serves. It checks nothing: a database
  outage must not make the kubelet restart a pod that cannot fix it.
- `/readyz` runs `injector.HealthCheckWithContext` under a 2 s timeout and is
  503 on any failure. Only the database is checked: `databaseCloser`
  implements `HealthCheck(ctx)` (one cheap ent query). The bucket and other
  services stay out, so a storage blip does not drain the whole fleet.
- `/readyz` is also 503 from the moment the process stops. `application.shutdown`
  calls `Probes.Drain()` before `http.Server.Shutdown`, so the pod leaves
  rotation while in-flight requests finish; `/healthz` stays 200 during the
  drain.
- Tests drive the built handler with `httptest` and inject the database fault
  with an ent interceptor, never DDL.

## Ops note (chart change, handed off)

`k8s/app-chart` still probes with `tcpSocket`. For the `api` and `api-check`
deployments switch to `httpGet` on the API port: `startupProbe` and
`livenessProbe` on `/healthz`, `readinessProbe` on `/readyz`, each with
`timeoutSeconds: 3` (above the 2 s readiness check). The chart is not edited
in this change.

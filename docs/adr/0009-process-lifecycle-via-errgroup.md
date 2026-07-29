# ADR-0009: Supervise the process lifecycle with errgroup

**Status:** Accepted

## Context

The server has three long-lived runtime components: River workers, the
Watermill router, and the HTTP server. They have readiness dependencies and
must stop in a deliberate order. `samber/do` constructs these dependencies but
does not model blocking actors, readiness, or propagation of runtime failures.

Hand-written goroutines and an error channel in `main` duplicated part of a
supervisor. They also left a partial-startup hole: if Watermill failed after
River started, `os.Exit` bypassed River shutdown and resource cleanup.

## Decision

`cmd/server` owns one package-local lifecycle coordinator based on
`golang.org/x/sync/errgroup`.

Startup is ordered by readiness:

1. start River;
2. run Watermill and wait for all subscribers to become ready;
3. start HTTP.

The first signal, startup error, runtime error, or unexpected actor completion
cancels the group. Shutdown then proceeds in the reverse dependency order:

1. drain HTTP, using `Close` as a fallback if graceful shutdown fails;
2. close Watermill;
3. stop River;
4. close the blob bucket and database, then flush telemetry.

The River and Watermill actors receive a process context that is cancelled only
after their graceful stop methods run. The context derived by `errgroup` is
used for supervision, not passed directly to actors, because its immediate
cancellation would race the required staged shutdown.

`samber/do` remains the dependency-construction mechanism. It is not a second
lifecycle owner, and production shutdown does not call `injector.Shutdown`.

## Alternatives

`oklog/run` provides execute/interrupt actor pairs, but it would not remove the
application-specific readiness gates or staged shutdown coordinator. Using it
here would add another abstraction without owning the hard parts of this
lifecycle.

Using only `samber/do` lifecycle hooks would hide ordering inside container
teardown and still would not supervise blocking actor failures.

## Consequences

- Partial startup follows the same cleanup path as steady-state failures.
- A signal is an expected termination and exits successfully; runtime failures
  remain non-zero exits.
- Unexpected clean completion of a long-lived actor is normalized to an error.
- Adding another runtime actor requires an explicit readiness and shutdown
  relationship in the coordinator rather than an implicit goroutine in `main`.

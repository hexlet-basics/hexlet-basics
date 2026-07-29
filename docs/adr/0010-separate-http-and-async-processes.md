# ADR-0010: Separate HTTP and asynchronous process runtimes

**Status:** Accepted

## Context

The Go server previously supervised HTTP, River workers, and Watermill
subscribers in one process. That made HTTP readiness and availability depend on
background consumers, pulled worker-only dependencies into the HTTP graph, and
coupled River concurrency to the number of HTTP replicas.

Both transports are PostgreSQL-backed and durable. Business writes can publish
Watermill events and enqueue River jobs without either consumer runtime being
available.

## Decision

Run two independent Go processes from the same codebase:

- `cmd/server` serves HTTP, publishes domain events, and uses an insert-only
  River client for transactional job enqueueing.
- `cmd/worker` runs Watermill subscribers and River workers.

The DI package exposes separate `NewServer` and `NewWorker` graphs. The server
graph has no River queues or worker registry and cannot execute jobs.

Each process owns its lifecycle through a package-local `errgroup`
coordinator. The worker starts River before Watermill so event handlers can
enqueue immediately, then stops Watermill before draining River so no new jobs
arrive during shutdown.

Watermill and River retain distinct semantics: Watermill routes facts and River
executes commands. Co-locating their consumers is only a deployment choice; a
Watermill handler still delegates durable, retryable work by enqueueing a River
job.

## Consequences

- Background startup, failure, and deployment no longer affect HTTP.
- HTTP replicas do not multiply River worker concurrency.
- Events and jobs accumulate durably while the worker is unavailable.
- Production and local development must launch both Go binaries.
- Watermill subscribers may move to a third process later without changing
  publishers or event contracts if they need independent scaling.

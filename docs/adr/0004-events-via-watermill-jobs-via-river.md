# Domain events via watermill, background jobs via river

Async work is split across two Postgres-backed mechanisms with an explicit
boundary:

- **watermill** owns **domain events** — "something happened" facts (e.g.
  `course_started`, `lesson_finished`, `solution_checked`, `user_signed_up`)
  published once and fanned out to N independent subscribers (amoCRM sync,
  Scenario Triggers keyed on event name, analytics, cache invalidation, email
  triggers). Replaces the legacy `app/events` + `app/handlers`.
- **river** owns **jobs** — "do this" commands: exercise-loader course-version
  builds, SMS/email sending, AI lesson review, finish-enrollments,
  reap-stuck-builds, scheduled/cron work, retries with backoff. Replaces the
  legacy Solid Queue `app/jobs`.

**Boundary rule:** a watermill subscriber that needs heavy, durable, retryable
work does NOT do it inline — it enqueues a river job. Events =
notification/routing; jobs = execution.

## Considered Options

- **Both (chosen)** — decoupled fan-out for events plus a real job queue.
- **river only (rejected)** — a pure job queue expresses a multi-consumer event
  poorly: the publisher would enqueue one job per handler and become coupled to
  its consumers, and there is no natural home for event-name-keyed Scenario
  Triggers.

## Consequences

- Two Postgres-backed async systems to operate and monitor; the boundary rule
  above keeps their responsibilities from bleeding together.

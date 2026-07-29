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

## Implementation

Domain events use Watermill's CQRS components over the PostgreSQL SQL Pub/Sub.
All event types share the durable `domain_events` topic. Each handler has a
stable, independently persisted consumer-group name; changing that name would
replay the stream for that handler.

Business state and its event are inserted through the same `database/sql`
transaction. Atlas creates the Watermill message and offset tables, and runtime
schema initialization is disabled. Handlers publish domain facts through their
own business module rather than calling downstream integrations or maintaining
a subscriber list.

Delivery is at least once. Subscribers must therefore be idempotent. Watermill
recovers panics, retries a handler three times with short exponential backoff,
then Nacks the message for durable SQL redelivery. A stuck consumer group does
not block independent groups.

Schema version 1 defines all legacy event names and data shapes:
`user_signed_up`, `user_signed_in`, `book_requested`, `course_started`,
`course_finished`, `lesson_started`, `lesson_finished`, `solution_checked`,
`email_confirmed`, and `lead_created`. Occurrence time is added to each Go
contract because Rails Event Store previously supplied that fact in its event
record rather than in `data`.

The runtime observes every event and logs only its event id, stable name, and
schema version, never its payload. The legacy `AmocrmHandler` is represented by
the `lead_created_to_amocrm_v1` Watermill consumer. It enqueues an
arguments-unique River job; the job performs the amoCRM call so Watermill
remains routing-only and River owns external-call retries. The amoCRM
`source_uid` also includes the lead id to make redelivery idempotent.

Only events whose corresponding write use case exists in the Go application
are emitted today. New write paths must publish the matching contract in the
same database transaction instead of calling integrations from HTTP handlers.

Version 1 retains domain-event rows indefinitely for replay and diagnosis.
Retention or PII erasure must be designed explicitly before adding a cleanup
job; deleting rows without considering every consumer offset can lose events.

## Considered Options

- **Both (chosen)** — decoupled fan-out for events plus a real job queue.
- **river only (rejected)** — a pure job queue expresses a multi-consumer event
  poorly: the publisher would enqueue one job per handler and become coupled to
  its consumers, and there is no natural home for event-name-keyed Scenario
  Triggers.

## Consequences

- Two Postgres-backed async systems to operate and monitor; the boundary rule
  above keeps their responsibilities from bleeding together.
- Event handlers can receive a message more than once and must make external
  effects idempotent.
- Event names, schema versions, and consumer-group names become persistent
  compatibility identifiers.

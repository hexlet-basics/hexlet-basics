---
status: accepted
---

# A domain event is a River job for each consumer that reacts to it

Supersedes the Watermill half of [ADR-0004](0004-events-via-watermill-jobs-via-river.md)
(River for jobs stays) and the Watermill paragraph of
[ADR-0027](0027-worker-conventions-river-and-watermill.md). Amends
[ADR-0009](0009-process-lifecycle-via-errgroup.md),
[ADR-0010](0010-separate-http-and-async-processes.md) and
[ADR-0023](0023-resource-shutdown-through-the-injector.md): the worker runs one actor, River.

**Problem.** Domain events reached their consumers through a Watermill outbox
(watermill-sql over PostgreSQL), and River ran the work. Two consumers used it:
`lead_created`, whose handler only inserted the amoCRM River job, and
`staff_read_personal_data`, whose handler wrote an Audit Record. Every other
event (nine of eleven) had no consumer but a log line, and stayed in the topic
table forever. For this the stack carried a second queue with its own retry
middleware, poison table, Sentry hook, consumer offsets, router lifecycle in the
worker and three tables, alongside River, which already does each of these.
Sphericon dropped the same outbox for the same reasons (its ADR 0101).

## Decision

- **Publishing a fact inserts the River job of each consumer that reacts to it,**
  in the caller's transaction (`events.Publisher.Publish`). The job exists exactly
  when the state change commits; an insert that fails rolls it back.
- **The consumers are named once, at composition.** `jobs.EventConsumers()` lists
  them; each is an `events.Consumer` that, given a fact, returns its job or
  declines. `di` hands them to the publisher of the server process. A producer
  raises a fact without knowing who reacts, and a new consumer touches no
  producer: this answers ADR-0004's objection to "River only" (the publisher would
  be coupled to its consumers).
- **A fact exists only with a consumer.** The nine events nobody consumed are
  removed with their publications (sign-up, sign-in, book request, course and
  lesson start and finish, solution check, email confirmation). One is added back
  together with the consumer that needs it.
- **Every job carries the whole fact** (`AmoCRMLeadArgs.Event`,
  `RecordPersonalDataReadArgs.Event`), so a consumer needs nothing the request had.
- **Failures are River's** (ADR-0027): attempts named per kind, `JobCancel` for what
  retrying cannot fix, Sentry on the final attempt, `discarded` rows as the
  dead-letter store. The poison table, its replay and the Watermill Sentry hook
  are gone.
- **The worker runs River alone.** It is ready once River has started; on SIGTERM
  it drains River, then the injector releases the resources. Only the server
  process publishes; the worker graph has no publisher.
- **A migration drops** `watermill_domain_events`, its offsets and its poison
  table. The Go stack has not served production, so no row is carried over.

## Considered options

- **Keep a plain append-only `domain_events` journal next to the jobs.** It would
  have kept ADR-0004's "retained indefinitely for replay" for the facts without a
  consumer. Rejected: nothing reads them, and a fact with no reader is a promise
  with no owner; the rows of the state (enrollments, lesson progress, leads) are
  the durable record.
- **Enqueue the job directly from the producer** (`leadcapture` inserting the
  amoCRM job). Simpler for two consumers, but it couples a business module to its
  integrations, which ADR-0004 rejected for good reason.
- **Keep Watermill.** Ordering and a cursor a new consumer could read back from:
  no consumer uses either, and the cost is a second queue to operate.

## Consequences

- One async system to operate and monitor. The amoCRM and audit work keeps its
  retries and history in `river_job`.
- A consumer added later sees only the facts published after it is deployed.
- Tests assert on the jobs a fact queued (`testsupport.EnqueuedJobs`) and run a
  job with `testsupport.WorkEnqueued` (River's `rivertest`), not through an outbox.

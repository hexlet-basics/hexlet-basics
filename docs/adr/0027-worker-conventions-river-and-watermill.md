---
status: accepted
---

# Worker conventions: River and Watermill

Supersedes only the "River usage rules" paragraph of
[ADR-0004](0004-events-via-watermill-jobs-via-river.md). The
split itself (Watermill for domain events, River for jobs, the boundary rule,
the event catalogue, the poison table) stays as ADR-0004 decided. This ADR adds
what ADR-0004 left unstated: who retries what, where a failure ends up, and
what a worker may assume about being run twice.

## Decision

**Insert policy (carried over).** Each job kind owns its policy in its Args'
`InsertOpts` (queue, attempts, uniqueness), its `Timeout`, and a `jobs.Module`
constructor. The worker composes the modules it can run and sets
`FetchOnlyKnownKinds`, so a job with no module in this worker waits for one that
has it. Queues isolate classes of work (`mail`, `builds`, `llm`, default). The
HTTP process schedules through `jobs.Enqueuer`; a Watermill handler has no
business write and inserts directly, relying on the kind's uniqueness.

**Idempotent by construction.** At-least-once delivery is the rule for both
systems, so every job and every subscriber must be safe to run twice, and
achieves it in its design (a claim in SQL, a state check, uniqueness over the
args), not with a "seen" table. An exception is allowed only when named and
commented where it lives. The two today:

- `ExerciseLoader` (`MaxAttempts: 1`): a build fetches, uploads and promotes, and
  records its own failure on the version row; the queue never retries it, an
  admin or the next webhook re-triggers it.
- `StaffReadPersonalData` consumer: a redelivered fact can leave a second audit
  line for one read, the safe direction for a trail (ADR-0020).

**Attempts are named.** Every kind sets `MaxAttempts` through a constant whose
comment gives the reason (`accountEmailMaxAttempts`, `llmMaxAttempts`,
`amoCRMLeadMaxAttempts`, `exerciseLoaderMaxAttempts`). River's default of 25
means weeks of retries and is never inherited by accident.

**Errors that retrying cannot fix are cancelled.** A worker returns
`river.JobCancel` for `domain.ErrNotFound`, `ErrConflict` and `ErrInvalid`
(`jobs.cancelUnfixable`): the row is gone or the data refuses the change, and
the next attempt gets the same answer. Anything else keeps River's backoff.

**Sentry sees the end, slog sees every attempt.** The River error handler logs
each failed attempt with the job's kind, queue and attempt numbers, and sends to
Sentry only when `Attempt >= MaxAttempts`. A retry is routine; the final
failure is the incident. A panic follows the same rule.

**Discarded jobs are the dead-letter store.** A job that used all its attempts
becomes `discarded`; the worker client sets `DiscardedJobRetentionPeriod`
explicitly (`jobs.DiscardedRetention`, 30 days) so the window is ours and not a
library default. An operator reads the error on the row and replays it with
River's own retry; there is no custom replay tooling.

**Watermill.** A handler is retried three times with short exponential backoff,
then the message goes to `domain_events_poison` and is acked (ADR-0004, kept).
A poisoned message is also captured in Sentry, between the poison queue and the
retry middleware so only an exhausted message is reported, so growth of the
poison table is noticed. Replay stays manual.

**The transaction adapter stays at the seams.** `database/sql` transactions and
the River and Watermill adapters that join them belong to `store`, to
`jobs.Enqueuer` (`Enqueue`, `EnqueueTx`) and to `events.Publisher`. A business
module asks to run in a transaction, to enqueue, or to publish; it does not
build a driver or a Watermill publisher. Modules that still receive a `*sql.Tx`
from `WithinTx` do so only to pass it to `EnqueueTx` or `Publish`.

## Consequences

- A flaky dependency no longer pages on every retry; the final failure does.
- A cancelled job keeps its error on the row, so a hopeless job is visible
  without having spent its attempts.
- Removing the `*sql.Tx` parameter from the progress, courseloader and other
  module signatures, so only the seams see it, is not done here and is not
  enforced by lint yet.

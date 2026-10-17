---
status: accepted
---

# Tests run in a rolled-back transaction on real PostgreSQL, over fixtures, with state asserted and no interaction checks

The test rules in `AGENTS.md` were reviewed against common Go practice and kept,
with three additions (parallelism for pure packages, `goleak`, `synctest`). Each
kept rule has a reason; a rule without one would have lost to practice.

**Kept, with the reason**

- **Real PostgreSQL, one rolled-back transaction per test.** Testcontainers start
  PostgreSQL 18, apply the Atlas migrations and load `fixtures/`; a test gets a
  transaction (`testsupport.NewClient*`, `NewHarness*`) that is rolled back. No
  database mocks. Code that needs commit semantics goes through the savepoint
  transactor in `testsupport`, or carries a `//nolint:forbidigo // why`. The
  savepoint adapter is the price of running production transaction code inside
  the rollback.
- **Fixtures only.** A test reads its data from `fixtures/` by generated
  `fixtures.Label`; `forbidigo` rejects `.Create()` in `_test.go`. Practice
  prefers per-test factories, but a labelled row names the case it exists for and
  a typo fails the build. The known cost is a growing fixture set; every row
  carries a comment saying which case it is for.
- **Assert on what an operation left behind** (rows, `PublishedEvents`,
  `EnqueuedJobs`), never on calls, their count or arguments. Mocks (moq,
  generated) exist only at seams to the outside world: today the LLM and the
  Docker runner, tomorrow any new external boundary, without a rule change.
- **`testsupport.Today` and the `WithClock` options** stay for tests that touch
  the database: PostgreSQL runs on real time, so `testing/synctest` cannot fake
  the `NOW()` in a query.

**Added**

- **`t.Parallel()` is not the default for DB-backed tests.** Transactions do not
  share state, but two of them updating the same fixture row block each other
  until one rolls back, and `Today` is package-level. It is adopted only where a
  measurement shows a clear gain. Packages with no database use `t.Parallel()`
  from now on.
- **`goleak.VerifyTestMain` in the packages that spawn goroutines** (the
  assistant, the exercise runner, the event publisher), not repo-wide, where it
  would be noise and would collide with the testcontainers startup.
- **`testing/synctest`** is allowed for in-memory timing and concurrency logic
  (timeouts, retries, debounce).
- **Style.** Table-driven subtests for pure functions with three or more cases;
  scenario tests for handlers. `require` for preconditions, `assert` for the
  checks after them. This is a review standard, not a linter rule: a table
  forced onto a handler test is a worse test.

## Considered Options

- **Mock the database or use SQLite.** Rejected: the schema, constraints,
  `jsonb`, and the outbox/queue tables are PostgreSQL behaviour that a stand-in
  does not reproduce.
- **Per-test factories instead of fixtures.** Rejected for the reasons above;
  revisit if the fixture set stops being reviewable.
- **`t.Parallel()` everywhere.** Deferred, not rejected: it needs the
  measurement above.

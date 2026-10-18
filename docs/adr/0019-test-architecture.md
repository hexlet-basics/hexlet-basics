---
status: accepted
---

# Tests run in a rolled-back transaction on real PostgreSQL, over fixtures, with state asserted and no interaction checks

The test rules in `AGENTS.md` were reviewed against common Go practice and kept,
with three additions (parallelism for pure packages, `goleak`, `synctest`). Each
kept rule has a reason; a rule without one would have lost to practice.

**Kept, with the reason**

- **Real PostgreSQL, one rolled-back transaction per test.** Each package's
  database is cloned by [pgtestdb](https://github.com/peterldowns/pgtestdb) on the
  development PostgreSQL 18 daemon (a service container in CI) from a template
  holding the Atlas migrations and `fixtures/`; a test gets a transaction
  (`testsupport.NewClient*`, `NewHarness*`) on a handle from
  [pgcleaner](https://github.com/mokevnin/pgcleaner)'s `Transaction` strategy
  (`testdb.DB(t)`), so it is rolled back. No database
  mocks. Code that needs commit semantics goes through the savepoint transactor
  in `testsupport`, or carries a `//nolint:forbidigo // why`. The savepoint
  adapter is the price of running production transaction code inside the
  rollback.
- **A test whose act is the commit gets a database of its own.** `testdb.Clone(t)`
  is `pgtestdb.New` over the same template: two connections racing for a row, or
  which connection sees an uncommitted write, run against a fresh clone that is
  dropped with the test, rather than committing into the package's database and
  cleaning up after themselves.
- **The server is shared, the databases are not.** A template is built once per
  content hash of `migrations/` and `fixtures/` and reused by every package, and a
  clone costs milliseconds where a container per package cost seconds. Clients
  keep the server's time zone: containers were always UTC and hid a partition
  function whose month bounds overlapped across a DST change on a server in
  another zone, so a zone-dependent migration must fail on a developer's machine
  rather than in production.
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
  would be noise.
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
- **A container per package (testcontainers-go).** Replaced: it needed Docker for
  every `go test`, and started and migrated a server per package. The development
  daemon already exists, and `CREATE DATABASE … TEMPLATE` is what Rails does for
  parallel tests.
- **A cloned database per test, for every test.** Deferred, not rejected, in
  favour of the rollback. A measurement on a loaded machine (load average ~30)
  put one clone at 130–265 ms for `CREATE DATABASE` alone and 46–236 ms for the
  drop, on a 15 MB template, and there are over a thousand tests. It would remove
  the savepoint adapter, so revisit it with a measurement on an idle machine;
  `testdb.Clone` is already the per-test path where it is needed.

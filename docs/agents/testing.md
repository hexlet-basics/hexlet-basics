# Writing tests

The rules for tests in this repository (moved out of `AGENTS.md`, which stays under
the 32 KiB that Codex reads). Commands to run them are in `AGENTS.md`, "Test".

- **A test never creates the data it arranges: it reads it from `fixtures/`.**
  Every row there is named for the case it exists for, with a label:
  `{{ id "user_removed" }}` in the YAML (primary keys and references alike), and
  `fixture.ID(fixture.UserRemoved)` (`fixture.ID32` for the contract's int32)
  in a test: every label has a generated `fixture.Label` variable (the
  `snake_case` label in `PascalCase`), and `ID` takes a `Label`, not a string, so
  a typo or a removed row fails the build. `NewHarnessAs(t,
  fixture.UserLearnerWithAGap)` signs a request in as that user. Add the row,
  then run `mise run gen:fixtures`; `fixture.LabelOf(name)` is only for a label
  read at run time (a request path). Ids count up in the order rows are written
  (the application reads "oldest first" off the id, so write older rows first);
  a label no row defines panics in the YAML, and `LabelOf` of one panics in
  `ID`. A case a test needs and no row has becomes a new row, with a comment
  above it saying what it is for, never an
  `.Create()` in the test: `forbidigo` rejects `.Create()`, `.CreateBulk()`,
  `sql.Open` and `store.NewDB` in `_test.go`, and a test whose creation is the
  act under test (a sign-up, a concurrent request that wins a race, a database
  constraint) says so with a `//nolint:forbidigo // why`. Flipping one flag on a
  fixture row (`admin=false`, `ClearCurrentVersionID`) or deleting a row is fine;
  a state two tests share is a row of its own. Do not assert how many rows a
  table holds (a new fixture breaks that) or a raw id; count the rows you start
  from, and assert on business facts (slug, order, label). A test that depends on
  today's date does not read the wall clock: the harness's clock is
  `testsupport.Today` (`assistant.WithClock`, `courseloader.WithClock`,
  `account.WithClock`) and the fixtures are written relative to it.
  The fixtures are edited by hand; `mise run db:fixtures` and the e2e database load
  them through `cmd/fixtures`, the same loader the test template database uses.
- A test owns no connection and cleans nothing up: it gets a transaction that is
  rolled back (`testsupport.NewClient*`, `NewHarness*`), the way a Rails test does.
  Jobs are worked in it with `rivertest.Worker` (`jobs.<Kind>Worker`), and what
  was queued or published is read back through the same transaction
  (`rivertest.RequireInsertedTx`, `testsupport.PublishedEvents`). Nothing runs a
  live River client or Watermill router against the shared database.
- Tests use `stretchr/testify` (`require` for preconditions, `assert` for the
  checks after them; `testifylint` enforces the idioms); each handler test gets an
  explicit `sql.Tx` rolled back during cleanup, matching Rails transactional
  tests, and builds a `Server` through `testsupport`. Why these rules: ADR-0019.
- **Parallel and leaks** (ADR-0019): a package with no database uses
  `t.Parallel()`; `handlers` runs its tests in parallel, which is safe only
  because no two tests touch the same fixture row (two transactions updating one
  row block each other, so a test that needs a row to itself gets its own
  fixture row). Adopt `t.Parallel()` in another DB-backed package only after
  that is true and measured. `goleak.VerifyTestMain` guards the
  goroutine-spawning packages (assistant, exercise runner, events), and
  `testing/synctest` is for in-memory timing, never for a query's `NOW()`.
  Table-driven subtests for pure functions with three or more cases, scenario
  tests for handlers.
  The harness runs the real event publisher and job enqueuers over the test
  transaction, so a test asserts on what an operation left behind
  (`testsupport.PublishedEvents`, `testsupport.EnqueuedJobs`, DB rows), never on
  what was called, how often or with which arguments: no call-verification
  tests. Only the LLM and the Docker runner are mocks (`testsupport.NewMocks`,
  moq): set a mock's `…Func` for what the boundary answers, and read the prompt
  back through the response.

## End-to-end specs

A spec only says what a person does and sees: it takes `test` from
`e2e/support/test.ts`, gets the `site` fixture and calls steps on it
(`site.session.signInAs(learner)`, `site.lesson.open(path)`), never `page`, a
selector or a wait. Three things hold that line: that `test` hands a scenario
no `page` (destructuring it fails `tsc -b`), `Site` keeps its browser in a
`#private` field (`site.page` fails `tsc -b`), and oxlint rejects importing
`@playwright/test` under `e2e/tests`. Steps live in a file per concept under
`e2e/support/` (`session.ts`, `catalog.ts`, `lesson.ts`, `admin.ts`,
`locale.ts`), and a scenario that needs a step the site does not have adds it
there.

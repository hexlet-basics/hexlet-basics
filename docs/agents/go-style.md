# Go style: errors, use-case modules, deviations

The rules for errors and use-case modules in the Go code (moved out of `AGENTS.md`, which stays
under the 32 KiB that Codex reads). The general Go style is in `AGENTS.md`, "Go Style".

## Errors (ADR-0018, ADR-0028)

- Every error in production code is built with `samber/oops`; `fmt.Errorf` is
  banned by `forbidigo` (allowed in tests). One `var errs = oops.In("<package>")`
  per package; only `In`, `With`, `Wrapf` and `Errorf` (no `Code`, `Public`,
  `Tags`, `Hint`).
- A message is static, lowercase, without trailing punctuation; the variable data
  goes in `.With("lesson_id", id)`, so errors group in Sentry and carry their data
  (one `BeforeSend` in the Sentry client copies the `oops` attributes).
- Wrap once, at the package boundary (the exported function that knows the
  operation); helpers return the error as is.
- Log or return, never both. The HTTP error handler and the River error handler
  are the one place that logs and reports.
- Look into a chain with `errors.AsType[T](err)`, not `errors.As` with a declared
  target. `oops.Recover` guards a goroutine that has no safety net (the
  assistant stream).
- Domain errors are plain `errors.New` values named `ErrX` (`domain.ErrNotFound`,
  `ErrConflict`, `ErrInvalid`, `ErrForbidden`, and a module's own, such as
  `account.ErrLastAdmin`, which wraps one of them), checked with `errors.Is`.
  `store.ToDomain` is the only translator of ent's errors (a module names the
  field behind a unique index with `store.UniqueOn`); a handler never reads an
  ent error (`forbidigo` rejects `ent.Is*`).
- HTTP: every non-2xx is RFC 9457 Problem Details, a 422 adds `errors:
  {field: [messages]}`. A handler returns the success variant and reports a
  failure through the `error` return; `errorStatus` is the one table from a
  domain error to a status and `detailFor` the one table to a localized `detail`
  (never `err.Error()`). A handler validates nothing: the contract, module rules
  returning `ErrInvalid`, and database constraints do.

## Use-case modules

The business rules behind an operation live in a package named for the concept
(`internal/<concept>/`), never in a handler body, so the HTTP handlers and any
other adapter (MCP) call the same code. There is no umbrella `usecases` or
`services` package. The rules below apply to every such module.

- **Package names** (ADR-0025): singular, named for the concept, and never equal
  to an `ent/<x>` package, so no file needs an import alias (`blog`,
  `permission`, `auth`, `catalog`, `curriculum`, `leadcapture`, `coursereview`,
  `sitebanner`, `account`, `book`, ...). `events`, `jobs` and `handlers` stay
  plural; a mass noun (`progress`, `dashboard`) stays as it is. `mise run
  lint:packages` enforces it. Name types for their role (`Editor`, `Directory`,
  `Registrar`), not `Service`, `Manager` or `Handler`.
- **Constructor:** `New(...)` returns the concrete type, with the required
  dependencies as positional arguments and functional options only for optional
  settings.
- **Interfaces:** a module exports none, and `handlers.Server` holds the concrete
  types. An adapter declares an interface only for a test substitute or a second
  real implementation: the LLM, the exercise runner, the mailer, the course
  fetcher (`testsupport` mocks them with moq), and the MCP server's `Verifier` and
  `Catalogue`. Where the dependency direction forbids naming the concrete type,
  take functions (`mcpserver.Deps`). Never add an interface with one
  implementation and no substitute.
- **Commands and results:** an operation takes plain Go values or a command
  struct named for the operation, and never an `internal/api` type (ADR-0001);
  adapters convert in both directions.
- **Pages:** a module's list takes a `domain.Page` and answers rows and a total;
  `store.Paged` runs the count and the window over one scoped, ordered query.
  The handler wraps them with `listFrom`.
- **Resources with no rules** (a plain create, read, update, delete of one
  entity) need no module: annotate the ent schema with `AdminCRUD` and the
  handler calls the generated `AdminList`/`AdminGet`/`AdminCreate`/`AdminUpdate`/
  `AdminDelete`, which is the single implementation (they take the contract's
  input type, because the generated `SetInput` does). The moment an operation
  has a rule of its own (a creator, a scope, a guard, a job to schedule), it
  moves to a module.
- **State moves** (Publication State, Course Readiness) are one shared rule: a
  `schema.StateMove` field gets a generated `Move<Field>` (ADR-0029), the only
  writer of the column; the handler passes it a `transition`.
- **Enforced:** `forbidigo` rejects an ent query builder (`Query`, `Create`,
  `Update*`, `Delete*`, `Get`) in `internal/handlers`, so a rule cannot slip
  back into a handler; the generated `Admin*` helpers pass. Keeping writes and
  reads in separate modules is review guidance only.

## Where we deviate from the vendor Go skills

The vendor Go skills (samber's set) are installed unedited. Where they disagree
with us, the rule below wins, for the reason in the linked ADR.

- **No `is`/`has`/`can` prefix on unexported bool fields and variables.** It is
  not standard Go, and ent and ogen output ignores it. Name a bool for what it
  says (`enabled`, `closed`). Recorded in [ADR-0025](../adr/0025-singular-package-names-and-adapter-declared-interfaces.md).
- **`ShutdownOnSignals` is not used.** The errgroup coordinators in `cmd/server`
  and `cmd/worker` own the signals and stop the actors in a fixed order; the
  injector shuts resources down after them. See [ADR-0023](../adr/0023-resource-shutdown-through-the-injector.md)
  and [ADR-0010](../adr/0010-separate-http-and-async-processes.md).
- **`events`, `jobs` and `handlers` stay plural**, against the singular-package
  rule, because `event`, `job` and `handler` are the commonest local variable
  names. See [ADR-0025](../adr/0025-singular-package-names-and-adapter-declared-interfaces.md).

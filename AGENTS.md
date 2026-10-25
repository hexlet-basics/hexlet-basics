# AGENTS Guide

This file is for coding agents working in `hexlet-basics`, the Go rewrite of the
original Rails app (retired and removed from the repository).
Prefer these instructions over generic Go or React defaults.

## What This Repo Is (Read First)

`hexlet-basics` is the Go replacement of the Rails app. It is one stack: a Go API
(`cmd/`, `internal/`, `ent/`) whose HTTP contract is generated from a TypeSpec
definition (`api-spec/`), plus a TanStack Start (React 19 + TypeScript) frontend
in `src/`.

The plan of record is **hard cutover at parity** (no long side-by-side run),
keeping backward-compat only on bcrypt passwords and URL routes. See
`docs/adr/` (0001–0022) for the binding architecture decisions. **Read the
relevant ADR before changing how a subsystem works** — they are the source of truth for why each library was
chosen, and several decisions (contract-first, hard cutover) constrain how you
should implement.

## Architecture: The Contract-First Pipeline

This is the single most important thing to understand (ADR-0001). The HTTP
contract is the source of truth, and both the Go server and the TS client are
**generated** from it. The flow is one-directional:

```
api-spec/*.tsp  ──tsp──▶  api-spec/dist/openapi.yaml  ──┬──ogen──▶  internal/api/   (Go server: routing, validation, types)
  (you edit)                    (generated)             └──hey-api▶  src/client/     (TS client + TanStack Query hooks)
```

- You hand-write **TypeSpec** in `api-spec/*.tsp` (`main.tsp`, `public.tsp`,
  `auth.tsp`, `admin.tsp`, `lesson-player.tsp`, `models.tsp`, `common.tsp`).
- `mise run gen` runs the whole pipeline: `gen:spec` (TypeSpec → OpenAPI) →
  `gen:api` (`go generate ./internal/apigen`, ogen → `internal/api/`) →
  `gen:client` (openapi-ts → `src/client/`).
- **Never hand-edit `internal/api/oas_*_gen.go`, `api-spec/dist/`, or the
  generated TS client.** Change the `.tsp`, run `mise run gen`, commit the outputs.
- Adding a contract operation adds a method to the generated `api.Handler`
  interface. Until you implement it, the embedded `api.UnimplementedHandler`
  keeps it compiling as "not implemented" — so the build never breaks on a
  contract addition; you fill in the handler after.
- ogen generates multipart uploads, but not OpenAPI `requestBody.encoding`
  ([ogen#1159](https://github.com/ogen-go/ogen/issues/1159)). TypeSpec emits
  that block for `HttpPart<File>`, so declare file parts as `HttpPart<bytes>`
  (same `format: binary` schema, no encoding block) and ogen generates the
  operation with an `ht.MultipartFile` field.
- **`@hey-api/openapi-ts` must stay on the `next` channel** while `typescript`
  is 7.x: the stable releases call the TypeScript compiler API (`ts.SyntaxKind`),
  which the native TS 7 package no longer exposes, so `gen:client` dies with
  `Cannot read properties of undefined`. The next-channel builds dropped that
  peer dependency. A "latest" bump therefore breaks codegen — the pin lives in
  both `package.json` and the `minimumReleaseAgeExclude` list in
  `pnpm-workspace.yaml`, and both must move together.

## Architecture: Handlers, DI, and Data

- **Handlers** live in `internal/handlers/`. `Server` (in `courses.go`) embeds
  `api.UnimplementedHandler` and implements contract methods. It holds the
  concrete module types from `handlers.Deps` (ADR-0025), not interfaces. A
  handler only translates the contract: the rules live in modules (see
  "Use-case modules").
- **ent → api conversion** is in `internal/apiconv/` (`apiconv.go` hand-written,
  `apiconv.gen.go` generated). Keep DB entities out of the API surface; convert
  through this layer. API models mirror the serializers, not `schema.rb`.
- **Push ordering/filtering into SQL** via ent query builders (see the ordered
  `LandingPage.Query()` in the `landingpage` module), not into Go post-processing.
- **DI** is `samber/do` v2, wired in `internal/di/`. Providers resolve their own
  deps from the injector, so the plain constructors (`store.NewClient`,
  `handlers.NewServer`, `api.NewServer`) stay injector-agnostic and are usable
  directly in tests. `di.NewServer` builds the HTTP/producer graph;
  `di.NewWorker` builds the Watermill/River consumer graph. `cmd/server` and
  `cmd/worker` each own a local `errgroup` coordinator (ADR-0010); `do` only
  constructs dependencies and, after the actors stop, tears them down
  (`dependencies.Shutdown`, ADR-0023). Add a new service to the appropriate
  process graph instead of threading it through constructors. A resource with no
  `Shutdown(ctx) error` of its own gets a small closer service in
  `internal/di/closers.go`.
- **ent ORM** lives in `ent/`; schema is `ent/schema/*.go`. Regenerate the
  client with `mise run gen:ent` (`go generate ./ent`) after editing a schema.
- Other `internal/` packages: `config/` (env config), `logging/` (`slog`),
  `store/` (ent client, `Transactor`, ambient transaction), `domain/` (domain
  errors, `FromStore`, `Page`), `health/` (`/healthz`, `/readyz`), `jobs/`
  (River), `events/` (Watermill), `ident/`, `apiconv/`, `fixture/`,
  `testsupport/`.
- **Configuration** (ADR-0024): `APP_ENV` (`development` default, `production`,
  `test`) drives environment policy. `config.LoadServer()` and `LoadWorker()`
  each parse and validate only their process's settings (`config.Common` plus
  their own), with `go-playground/validator` tags next to the `env` tags; the
  deploy-critical settings are required, with no code default. Secrets are
  `config.Secret` (prints `[REDACTED]`; `Reveal()` reads it). A constructor
  takes the section or the scalars it reads, never `*config.Server` /
  `*config.Worker` outside `di`, and `os.Getenv` stays in `internal/config`
  (`mise run lint:config`). `.env` is read only by `config.LoadDotenv()` in
  `cmd/*/main`, and only under `APP_ENV` unset or `development`.
- **Logging** (ADR-0026, extends ADR-0007): `slog`; JSON off a terminal, tint on
  one; `LOG_LEVEL` (default `info`). Use the `*Context` variant when a `ctx` is
  in scope, `snake_case` keys, `err` for an error, a static message with the
  variable parts as attributes (`sloglint` enforces it). `trace_id` is the only
  correlation key. Never log passwords, tokens, cookies or emails (identity is
  `user_id`); the key-name `ReplaceAttr` is a net, not the standard.
- **Lifecycle and health** (ADR-0023): the server drains HTTP; the worker stops
  Watermill, then River; then the injector shuts the resources down. `/healthz`
  is liveness only and checks nothing; `/readyz` checks the database alone and
  turns 503 at SIGTERM, before the drain. Both are outside the contract and
  `/api`; the worker serves them on `HEALTH_ADDR`. Chart probes are an ops
  handoff (the chart still uses `tcpSocket`).
- **Workers** (ADR-0027): every job and subscriber is idempotent by construction
  (a named, commented exception otherwise: `ExerciseLoader`); `MaxAttempts` is a
  named constant with a reason, never River's default; a domain error retrying
  cannot fix returns `river.JobCancel`; Sentry gets a job error only on its
  final attempt; `discarded` rows (retention `jobs.DiscardedRetention`) are the
  dead-letter store, replayed with River's own retry; a poisoned Watermill
  message goes to `domain_events_poison` and to Sentry.
- **Transactions, a known gap** (ADR-0027, #924): the tx adapter belongs to
  `store`, `jobs.Enqueuer` and `events.Publisher`, but business modules still
  receive a `*sql.Tx` from `WithinTx` (to pass to `EnqueueTx`/`Publish`), and no
  lint enforces the confinement. Do not add new modules that take a `*sql.Tx`
  for anything else; removing the parameter is not done yet.

## Database & Migrations (atlas owns the schema)

- **atlas** owns the schema going forward (versioned SQL in `migrations/`,
  config in `atlas.hcl`). The baseline is the retired Rails schema. This is *not*
  ent auto-migrate and *not* Rails migrations.
- New schema change: edit `ent/schema`, then `mise run db:migrate-new <name>` to
  scaffold a migration to hand-author, then `atlas migrate hash`.
- **No database triggers of our own**, in migrations or in tests (river's are the only ones). A rule
  lives in Go, where it can be read and tested; a test injects a fault with an ent hook or
  interceptor, never DDL or a trigger. `mise run lint:triggers` rejects a migration that creates one.
- Local PostgreSQL 18 runs as the `db` mise daemon (`mise.toml`, pitchfork) on
  port **54330** (via `mise run db:up`) to avoid clashing with a `5432` DB. It
  is a standing service: leave it running. Tests use the same server but never
  its development database: each test package gets a database of its own,
  cloned from a migrated, fixture-loaded template (see Test).

## Setup & Dev Commands (run from repo root)

- Tasks live in `mise/tasks.toml` (`mise tasks` lists them) and the dev stack is
  the `[daemons]` table in `mise.toml`; there is no Makefile at the root.
- `mise run setup` — the one task a fresh checkout needs: toolchain + git hooks,
  seed `.env`, then the separate `install` task (dependencies) and `db:fixtures`
  (db daemon, migrations, `fixtures/`) and `dev`: after a clone, `mise setup`
  alone brings the whole project up with a populated site.
- `mise run install` — only the dependencies: `go mod download` + `pnpm install`
  (root package.json is shared by the frontend and the api-spec tooling).
- `mise run db:up` — start the Postgres daemon; idempotently creates the
  development database.
- `mise run dev` — start the API, async worker (both air live-reload), and Vite
  as daemons; `mise run dev:down` stops them (never the database).
  `mise daemons logs <api|worker|web|db>` shows output. The api and web ports
  are dynamic (base `3002` / `5173`, next free when taken; see `mise env | grep PORT`),
  and `mise.toml` pins `DATABASE_URL` to the development DB because the postgres
  preset would export its own.
- `mise run dev:api` / `dev:worker` / `dev:web` — one process in the foreground.
- `mise run dev:spec` — watch TypeSpec and re-emit OpenAPI on change.
- Go tooling that is cheap to build is pinned by `tool` directives in `go.mod`
  and invoked as `go tool <name>` (air, ent, ogen, goverter, goi18n) — add new
  ones with `go get -tool`, never by hand. `mise.toml` pins the rest: the go
  toolchain, atlas as the Apache-2.0
  `atlas-community` build (ADR-0014; its public CLI module is frozen at a
  2024-11 0.13.x snapshot and will not install — its go.mod has replace
  directives; `mise run lint:atlas-pin` keeps `mise.toml` and the `Dockerfile` on
  one version, and the nightly `lint:atlas-retention` fails before the pinned
  binary ages out; the ent/atlas exit plan, sqlc, is in ADR-0014), and golangci-lint (a tool directive would drag its whole
  analyzer tree into go.mod to rebuild what mise fetches as a binary).
  Git hooks are hk (`hk.pkl`, pinned in `mise.toml` with pkl), installed
  by `mise run setup` (and mise's postinstall hook) via `hk install --mise`.
  Use `pnpm`, never `npm`/`npx` directly for JS deps.

## Code search

`ast-grep` is installed (pinned in `mise.toml`). For any search that depends on
syntax or code structure (calls, declarations, a method on a type, JSX props),
default to `ast-grep --lang <go|tsx|typescript> -p '<pattern>'` and adjust
`--lang` per file type. Use `rg`/grep only for plain-text searches (strings,
comments, config values). Skip generated files (`internal/api/oas_*_gen.go`,
`src/client/`, `ent/` outside `schema/`) unless that is the point of the search.

## Codegen Commands

- `mise run gen` — full contract pipeline (spec → ogen server → hey-api client).
- `mise run gen:spec` / `mise run gen:api` / `mise run gen:client` — individual stages.
- `mise run gen:ent` — regenerate the ent ORM from `ent/schema`.
- `mise run gen:all` — `gen:ent` + `gen` + `gen:locales` + `gen:amocrm` + `gen:mocks` + `gen:fixtures`.
- `mise run gen:mocks` — regenerate the moq mocks of the external boundaries (LLM, exercise runner) in `internal/testsupport/mock_*_gen.go` (directives in `internal/testsupport/gen.go`). `mise run gen:fixtures` — regenerate `internal/fixture/labels_gen.go`, one `fixtures.Label` variable per label in `fixtures/*.yml` (directive in `internal/fixture/gen.go`, generator in `cmd/fixturelabels`). `mise run tidy` — `go mod tidy`.
- `mise run gen:check` — runs `gen:all` and fails if the tree changes; the CI `gen-check` job runs it, so generated output must be committed with its source.
- After changing a `.tsp`, ent schema, or deps, run the matching `gen` target
  and commit the regenerated outputs alongside the source change.

## Lint & Typecheck

- `mise run lint` — Go (`gofmt` check, `go vet`, `golangci-lint`: see below) +
  frontend (`oxlint`, `oxfmt --check`, `tsc -b`) + ADR number uniqueness
  (`lint:adr`) + `lint:packages` (internal package names are singular, ADR-0025)
  + `lint:config` (ADR-0024) + `lint:atlas-pin` (Atlas version in `mise.toml`
  and `Dockerfile` agree, ADR-0014) + `lint:triggers` + `lint:deps` (`go mod tidy
  -diff`) + `lint:typos` + `lint:workflows` (zizmor, jactionlint)
  + `lint:shell` (shellcheck) + `lint:unused` (knip) + `lint:i18n` (generated
  types current, no hardcoded strings outside `src/components/pages`, no key missing
  from a locale catalog) + `lint:helm` (helm lint, render, kubeconform) +
  `lint:migrations` (shipped migrations immutable, new ones squawk-clean; an edit
  needs a `Migration-Edit: <why>` line in its commit) + `lint:contract` (vacuum) +
  `lint:agents` (agnix); `lint:contract-breaking` (oasdiff) is manual. CI runs the
  same task. The pre-push hook runs the Go, frontend and ADR checks; the
  pre-commit hook runs gitleaks, oxlint, oxfmt, golangci-lint, typos,
  shellcheck, zizmor, jactionlint, terraform and helm on the staged files.
  Pushes need them clean.
- **golangci-lint** (`.golangci.yml`) enables, beyond the standard set,
  `gosec` (tests exempt from G115, G304, G306, G301, G101 only; a production
  finding is fixed or carries a `//nolint:gosec // why`), `modernize`,
  `testifylint`, `errname`, `sloglint`, `errorlint`, `contextcheck`, `noctx`,
  `bodyclose`, `forbidigo` and the rest listed there. The rule is: a linter is
  enabled in the change that clears its findings, with no baseline and no
  exclusion list that hides debt. `forbidigo` bans raw SQL, `.Create()` in tests,
  ent query builders and `ent.Is*` in handlers, and `fmt.Errorf` in production
  code (no exclusions left). `nolintlint` requires a reason on every `nolint`.
- `mise run lint:size` — the client bundle budgets (size-limit config in
  `package.json`): builds the web app, then fails if the entry chunk or the
  total client JavaScript grows past its gzip limit. CI runs it as the
  `bundle-size` job; it is not part of `lint` because it needs a build.
- `mise run lint:fix` — `gofmt -w`, `golangci-lint --fix`, `oxlint --fix`,
  `oxfmt` and `typos --write-changes`.
- Go vet/build only: `go vet ./...`, `go build ./...`.

## Test

- **Heavy tasks are serialized.** Several sessions and worktrees share one machine, so the
  leaf tasks that saturate it (`go test`, `go vet`, `golangci-lint`, `tsc -b`, vitest, builds,
  e2e, `gen:check`, the pre-commit golangci-lint) run through
  [`heavy`](https://github.com/mokevnin/heavy) (`"github:mokevnin/heavy"` in `mise.toml`), which
  holds one of a few machine-wide slots for the command's duration and prints
  `all N slots of pool 'default' are busy, waiting for one` when it queues. Use the mise tasks;
  when you must run one of these commands by hand, prefix it with `heavy --` and never start
  them in parallel yourself. `HEAVY_SLOTS` sets how many run at once, `heavy --status` shows who
  holds them; a nested call (`HEAVY_HELD`) and `CI` run at once. Wrap leaf tasks only, never a
  parent such as `lint`, or it would hold a slot for children that wait on it.

- `mise run test` — the Go suite (`go test -race ./...`, plus `-shuffle=on` when `CI` is set: a
  shuffle disables go's result cache, so locally unchanged packages show `(cached)`). In CI tests run in a random
  order and `go` prints the seed, so a test must not depend on another's side
  effects; replay an order with `mise run test -- -shuffle=<seed> ./pkg`.
- `mise run test:docker` — the tests that need a real Docker daemon (build tag
  `docker`, `internal/exerciserunner`: the container lifecycle). `mise run test` leaves them
  out; CI runs them as a separate step. A test of that kind goes in a `//go:build docker` file.
- `mise run test:cover` — the Go suite (with `-race`, docker-tagged tests included) over `cmd/` and `internal/` with
  coverage, failing under the total in `.testcoverage.yml` (generated code and
  `testsupport` excluded). `pnpm test:coverage` does the same for the frontend
  with the thresholds in `vite.config.ts`. Raise a floor when coverage rises;
  never lower it.
- **Nightly, not per PR.** `.github/workflows/nightly.yml` runs both coverage
  gates, `mise run lint:vulns` (govulncheck) and `mise run lint:atlas-retention`
  (the pinned Atlas binary is still downloadable) once a night and on demand
  (`workflow_dispatch`). A new advisory can fail `lint:vulns` with no change
  on our side, so it is deliberately not part of `lint`. Put a check there
  when its result does not depend on the diff under review.
- DB-dependent packages clone their database with
  [pgtestdb](https://github.com/peterldowns/pgtestdb) on the `db` daemon
  (start it with `mise run db:up`; a missing server fails fast with that hint). A template database holds the
  Atlas migrations and the committed `fixtures/`, is built once per content hash
  of both and shared by every package, and each package (`testdb.Run` in its
  `TestMain`) clones it in milliseconds and drops the clone at exit. No Docker and
  no test DB preparation. CI points `TEST_DB_HOST`/`TEST_DB_PORT` at its
  PostgreSQL service. The clients keep the server's time zone on purpose: a
  migration or function that depends on the session zone (month arithmetic on a
  `timestamptz` across a DST change) must fail here, on a developer's machine,
  not in production. Do such arithmetic on UTC timestamps without a zone.
  A test that needs real commits (the commit is the act: two connections, a race
  for a row) takes a database of its own with `testdb.Clone(t)`, dropped when the
  test ends, instead of committing into the package's database.
- `mise run db:prepare` prepares the development DB without replacing its data.
- `mise run db:fixtures` migrates the development DB and loads `fixtures/` into it
  (replaces the fixture tables' data; the old `rails db:fixtures:load`). A fresh
  development DB shows an empty site until this has run.
- `mise run test:e2e` — the Playwright suite in `e2e/`: Chromium against the
  built Go API and SSR frontend on a throw-away `code_basics_e2e_test` database
  (migrations + `fixtures/`, `test:e2e:db`). Needs the db daemon (`mise run
  db:up`) and `pnpm exec playwright install chromium` once; never part of
  `mise run test`. It asserts on what pages show, mocks nothing, and does not run
  lesson code (that needs an environment the stack lacks).
  A spec only says what a person does and sees, through the `site` fixture of
  `e2e/support/test.ts`; see `docs/agents/testing.md`, "End-to-end specs".
- Single package / test: `go test ./internal/handlers/`,
  `go test ./internal/handlers/ -run TestListCourses -v`.
- **Writing a test: read `docs/agents/testing.md` first.** Its rules are enforced
  by `forbidigo` and review: a test reads the data it arranges from `fixtures/`
  (never `.Create()`), runs in a rolled-back transaction, asserts on what an
  operation left behind and not on calls, and mocks only the LLM and the Docker
  runner (ADR-0019).

## Build

- `mise run build` — `build:api` (`bin/server`) + `build:worker` (`bin/worker`) +
  `build:web` (`pnpm build`). `mise run clean` removes `bin/` and `dist/`.

## Frontend (`src/`)

- TanStack Start (SSR, ADR-0008) + React 19 + TypeScript (`strict`) + Vite.
  Routing is file-based under `src/routes/`; `routeTree.gen.ts` is **generated**
  — don't hand-edit it. UI uses Mantine; data fetching uses TanStack Query hooks
  generated by hey-api into `src/client/`.
- **Lint policy (oxlint, `oxlint.config.ts`): everything is on; only contradiction
  cases are off.** Every rule of every stable category (correctness, suspicious,
  pedantic, perf, style, restriction) across all plugins, type-aware rules
  included. A rule is switched off only in that file, with a one-line reason
  beside it, and only when it directly contradicts oxfmt or another enabled rule
  (not because its fix is large). Prefer an option over an off switch: the size
  rules (`max-lines`, `max-lines-per-function`, `max-statements`, `complexity`)
  run with a threshold, and `no-void` allows `void` as a statement. Two
  exceptions are not contradictions and say so beside them: the `react-perf`
  rules (nothing here uses `React.memo`) and `oxc/no-optional-chaining`. No
  baseline, no blanket ignore: fix a finding, or show the contradiction.
- `any` is banned (oxlint `typescript/no-explicit-any` plus the type-aware `no-unsafe-*` rules): use `unknown` and narrow, or a precise type.
- Call the API through the **generated client / Query hooks**, not hand-written
  `fetch`. If an endpoint is missing, add it in TypeSpec and regenerate.
- State: `zustand`; validation: `zod`; i18n: `i18next` + `react-i18next`
  (`src/locales/`, config via `i18next-cli`). Animation: `motion`.
- **Library docs, ask the source before guessing an API.** TanStack (Router,
  Start, Query, Form, Table): `pnpm exec tanstack search-docs "<query>"
  [--library router]`, then `pnpm exec tanstack doc <library> <path>` for the
  page. Mantine: the `mantine` MCP server in `.mcp.json` (`search_docs`,
  `get_item_props`). The TanStack CLI has no `mcp` command anymore.

## Generated / Do-Not-Edit Files (root)

Change the source, run the generator, commit the output. Never hand-edit or
`git restore` a generated file — regenerate it from current source.

- `internal/testsupport/mock_*_gen.go` (moq mocks) ← the production interfaces via `mise run gen:mocks`.
- `internal/fixture/labels_gen.go` (fixture labels) ← `fixtures/*.yml` via `mise run gen:fixtures`.
- `api-spec/dist/**` (OpenAPI) ← `api-spec/*.tsp` via `mise run gen:spec`.
- `internal/api/oas_*_gen.go` (ogen server) ← OpenAPI via `mise run gen:api`.
- `internal/amocrm/generated/**` (ogen client) ← the pinned amoCRM OpenAPI via
  `mise run gen:amocrm` (`go:generate` in `internal/amocrm/gen.go`).
- `src/client/**` (hey-api TS client + Query hooks) ← OpenAPI via `mise run gen:client`.
- `ent/**` (except `ent/schema/`) ← `ent/schema` via `mise run gen:ent`.
- `internal/apiconv/apiconv.gen.go` — generated converter.
- `src/routeTree.gen.ts` — TanStack Router route tree.
- `CHANGELOG.md` and `version` in `package.json` — release-please; never bump
  by hand.
- Sources you DO edit: `api-spec/*.tsp`, `ent/schema/*.go`, Go handlers,
  `src/locales` string values, `migrations/` (hand-authored, then `atlas migrate hash`).

## Go Style

- `gofmt`-clean, passes `go vet` and `golangci-lint` (config governs `cmd`,
  `internal`, `ent/schema`; generated files are skipped).
- Where our rules and established Go practice disagree, practice wins; a rule of
  ours survives only with a reason recorded in an ADR (see "Where we deviate
  from the vendor Go skills").
- Match the surrounding package: doc comments on exported types/functions
  explaining *why* (the existing files are heavily commented on rationale:
  mirror that density on non-obvious logic, e.g. why an ordering has a specific
  tie-breaker), constructors return concrete types.
- Return early on errors; don't swallow them. API types stay at the boundary:
  convert via `apiconv`. The business rules behind an operation live in a
  use-case module (below), not in the handler.

### Errors (ADR-0018, ADR-0028)

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
  `domain.FromStore` is the only translator of ent's errors (a module names the
  field behind a unique index with `domain.UniqueOn`); a handler never reads an
  ent error (`forbidigo` rejects `ent.Is*`).
- HTTP: every non-2xx is RFC 9457 Problem Details, a 422 adds `errors:
  {field: [messages]}`. A handler returns the success variant and reports a
  failure through the `error` return; `errorStatus` is the one table from a
  domain error to a status and `detailFor` the one table to a localized `detail`
  (never `err.Error()`). A handler validates nothing: the contract, module rules
  returning `ErrInvalid`, and database constraints do.

### Use-case modules

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
  `domain.Paged` runs the count and the window over one scoped, ordered query.
  The handler wraps them with `listFrom`.
- **Resources with no rules** (a plain create, read, update, delete of one
  entity) need no module: annotate the ent schema with `AdminCRUD` and the
  handler calls the generated `AdminList`/`AdminGet`/`AdminCreate`/`AdminUpdate`/
  `AdminDelete`, which is the single implementation (they take the contract's
  input type, because the generated `SetInput` does). The moment an operation
  has a rule of its own (a creator, a scope, a guard, a job to schedule), it
  moves to a module.
- **Enforced:** `forbidigo` rejects an ent query builder (`Query`, `Create`,
  `Update*`, `Delete*`, `Get`) in `internal/handlers`, so a rule cannot slip
  back into a handler; the generated `Admin*` helpers pass. Keeping writes and
  reads in separate modules is review guidance only.

### Where we deviate from the vendor Go skills

The vendor Go skills (samber's set) are installed unedited. Where they disagree
with us, the rule below wins, for the reason in the linked ADR.

- **No `is`/`has`/`can` prefix on unexported bool fields and variables.** It is
  not standard Go, and ent and ogen output ignores it. Name a bool for what it
  says (`enabled`, `closed`). Recorded in [ADR-0025](docs/adr/0025-singular-package-names-and-adapter-declared-interfaces.md).
- **`ShutdownOnSignals` is not used.** The errgroup coordinators in `cmd/server`
  and `cmd/worker` own the signals and stop the actors in a fixed order; the
  injector shuts resources down after them. See [ADR-0023](docs/adr/0023-resource-shutdown-through-the-injector.md)
  and [ADR-0010](docs/adr/0010-separate-http-and-async-processes.md).
- **`events`, `jobs` and `handlers` stay plural**, against the singular-package
  rule, because `event`, `job` and `handler` are the commonest local variable
  names. See [ADR-0025](docs/adr/0025-singular-package-names-and-adapter-declared-interfaces.md).

## Conventional Commits & Releases

- Every commit follows [Conventional Commits](https://www.conventionalcommits.org/)
  (`<type>(<scope>?): <description>`); `commitlint` (commit-msg hook) enforces
  types. Commit messages are in English; the body explains *why*, not *what*.
- Releases are cut by `release-please` (config `release-please-config.json`,
  same type→bump rules as before). It maintains a release PR bumping
  `package.json` and `CHANGELOG.md`; merging it tags `vX.Y.Z`.
- Type → bump (pre-1.0, so bumps shift down one level): `feat` → patch,
  `fix`/`perf`/`revert` → patch, `chore`/`build` → patch (via
  `extra-release-commit-types`), breaking (`type!:` / `BREAKING CHANGE:`) →
  minor. `refactor`/`test`/`docs`/`ci`/`style` → no release. Changelog sections:
  Features, Bug Fixes, Performance Improvements, Reverts, Miscellaneous
  (`chore`), Build System (`build`).

## Agent skills

### Issue tracker

Issues live in GitHub Issues on `hexlet-basics/hexlet-basics`, driven by the `gh`
CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical triage roles, each label string equal to its name. See
`docs/agents/triage-labels.md`.

### Domain docs

Single-context: `GLOSSARY.md` + `docs/adr/` at the repo root. See
`docs/agents/domain.md`.

### Git workflow

The main checkout stays on `main` (never switch branches there); trivial changes
are committed on `main`, anything bigger is done in a worktree. See
`docs/agents/git-workflow.md`.

## High-Value Agent Workflow

- To change the API: edit TypeSpec → `mise run gen` → implement/adjust the Go
  handler and, if needed, the frontend client usage → run the narrowest Go test.
- After Go changes: `go build ./...`, then `go test ./...` (or a single package);
  DB-dependent packages clone their own database automatically, from the `db`
  daemon: run `mise run db:up` first when calling `go test` directly.
- After frontend changes: run `tsc -b` (`pnpm check`).
- Before finishing substantial work: `mise run lint` and the most relevant tests.
- Prefer minimal diffs; don't start broad refactors or "normalize" unrelated
  files. When unsure whether a file is generated, find its generator before
  editing (`grep` `mise/tasks.toml` / `go:generate` directives).

## Ops

- The Go stack deploys with `k8s/app-chart` (Helm, `k8s/Makefile`) from the
  images built by `Dockerfile` (server, worker, atlas) and `Dockerfile.web`
  (SSR); `release-please.yml` publishes them to GHCR on release. The sops secrets
  file is `k8s/secrets.yaml` (`SECRETS=` overrides it in `k8s/Makefile`) and
  `terraform/` reads it (`mise run lint:infra` checks the Terraform). If you
  touch any of it, call out the change clearly in handoff notes. Never commit secrets.
- Handed-off chart changes (the chart is not edited by the Go work): HTTP probes
  on `/healthz` and `/readyz` for the API and the worker (`HEALTH_ADDR`), and
  `APP_ENV=production` plus the required variables of ADR-0024 on both workloads.

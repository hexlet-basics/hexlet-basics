# AGENTS Guide

This file is for coding agents working at the **root** of `hexlet-basics`. The
root is the Go rewrite; the original Rails app lives in `legacy/` with its own
`AGENTS.md`.
Prefer these instructions over generic Go or React defaults.

## What This Repo Is (Read First)

`hexlet-basics` is mid-migration from Rails to Go. Two apps coexist in one repo:

- **Root = the new stack.** A Go API (`cmd/`, `internal/`, `ent/`) whose HTTP
  contract is generated from a TypeSpec definition (`api-spec/`), plus a
  TanStack Start (React 19 + TypeScript) frontend in `src/`.
- **`legacy/` = the Rails app** being replaced (Rails + Inertia + React). It has
  its own Makefile, package.json, and agent guide. When your work is inside
  `legacy/`, follow `legacy/AGENTS.md`, not this file.

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
  `api.UnimplementedHandler` and implements contract methods; it holds an
  `*ent.Client` and an `apiconv.Converter`. Handler methods run ent queries and
  return generated `api.*` types.
- **ent → api conversion** is in `internal/apiconv/` (`apiconv.go` hand-written,
  `apiconv.gen.go` generated). Keep DB entities out of the API surface; convert
  through this layer. API models mirror the serializers, not `schema.rb`.
- **Push ordering/filtering into SQL** via ent query builders (see the ordered
  `LandingPage.Query()` in `courses.go`), not into Go post-processing.
- **DI** is `samber/do` v2, wired in `internal/di/container.go`. Providers
  resolve their own deps from the injector, so the plain constructors
  (`store.NewClient`, `handlers.NewServer`, `api.NewServer`) stay
  injector-agnostic and are usable directly in tests. `di.NewServer` builds the
  HTTP/producer graph; `di.NewWorker` builds the Watermill/River consumer graph.
  `cmd/server` and `cmd/worker` each own a local `errgroup` coordinator; `do`
  only constructs dependencies. Add a new service to the appropriate process
  graph instead of threading it through constructors.
- **ent ORM** lives in `ent/`; schema is `ent/schema/*.go`. Regenerate the
  client with `mise run gen:ent` (`go generate ./ent`) after editing a schema.
- Other `internal/` packages: `config/` (env config via `caarlos0/env`),
  `logging/` (`slog` + `tint`), `store/` (ent client construction),
  `ids/`, `apiconv/`, `testsupport/`.

## Database & Migrations (atlas owns the schema)

- **atlas** owns the schema going forward (versioned SQL in `migrations/`,
  config in `atlas.hcl`). The baseline is the retired Rails schema. This is *not*
  ent auto-migrate and *not* Rails migrations.
- New schema change: edit `ent/schema`, then `mise run db:migrate-new <name>` to
  scaffold a migration to hand-author, then `atlas migrate hash`.
- Local PostgreSQL 18 runs as the `db` mise daemon (`mise.toml`, pitchfork) on
  port **54330** (via `mise run db:up`) to avoid clashing with a `5432` DB. It
  is a standing service: leave it running. It is only for development; tests own
  disposable PostgreSQL 18 containers.

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
  directives), and golangci-lint (a tool directive would drag its whole
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
- `mise run gen:mocks` — regenerate the moq mocks of the external boundaries (LLM, exercise runner) in `internal/testsupport/mock_*_gen.go` (directives in `internal/testsupport/gen.go`). `mise run gen:fixtures` — regenerate `internal/fixtures/labels_gen.go`, one `fixtures.Label` variable per label in `fixtures/*.yml` (directive in `internal/fixtures/gen.go`, generator in `cmd/fixturelabels`). `mise run tidy` — `go mod tidy`.
- `mise run gen:check` — runs `gen:all` and fails if the tree changes; the CI `gen-check` job runs it, so generated output must be committed with its source.
- After changing a `.tsp`, ent schema, or deps, run the matching `gen` target
  and commit the regenerated outputs alongside the source change.

## Lint & Typecheck

- `mise run lint` — Go (`gofmt` check, `go vet`, `golangci-lint` with
  contextcheck/noctx/errorlint/bodyclose among others) + frontend (`oxlint`,
  `oxfmt --check`, `tsc -b`) + ADR number uniqueness + `lint:deps` (`go mod tidy
  -diff`) + `lint:typos` + `lint:workflows` (zizmor, jactionlint)
  + `lint:shell` (shellcheck) + `lint:unused` (knip) + `lint:i18n` (generated
  types current, no hardcoded strings outside `src/components/pages`, no key missing
  from a locale catalog) + `lint:helm` (helm lint, render, kubeconform) +
  `lint:infra` (terraform fmt and validate for `legacy/terraform`). CI runs the
  same task. The pre-push hook runs the Go, frontend and ADR checks; the
  pre-commit hook runs gitleaks, oxlint, oxfmt, golangci-lint, typos,
  shellcheck, zizmor, jactionlint, terraform and helm on the staged files.
  Pushes need them clean.
- `mise run lint:size` — the client bundle budgets (size-limit config in
  `package.json`): builds the web app, then fails if the entry chunk or the
  total client JavaScript grows past its gzip limit. CI runs it as the
  `bundle-size` job; it is not part of `lint` because it needs a build.
- `mise run lint:fix` — `gofmt -w`, `golangci-lint --fix`, `oxlint --fix`,
  `oxfmt` and `typos --write-changes`.
- Go vet/build only: `go vet ./...`, `go build ./...`.

## Test

- `mise run test` — the Go suite (`go test -shuffle=on ./...`). Tests run in a random
  order and `go` prints the seed, so a test must not depend on another's side
  effects; replay an order with `mise run test -- -shuffle=<seed> ./pkg`.
- `mise run test:cover` — the Go suite over `cmd/` and `internal/` with
  coverage, failing under the total in `.testcoverage.yml` (generated code and
  `testsupport` excluded). `pnpm test:coverage` does the same for the frontend
  with the thresholds in `vite.config.ts`. Raise a floor when coverage rises;
  never lower it.
- **Nightly, not per PR.** `.github/workflows/nightly.yml` runs both coverage
  gates and `mise run lint:vulns` (govulncheck) once a night and on demand
  (`workflow_dispatch`). A new advisory can fail `lint:vulns` with no change
  on our side, so it is deliberately not part of `lint`. Put a check there
  when its result does not depend on the diff under review.
- DB-dependent packages use testcontainers-go to start isolated PostgreSQL 18
  containers, apply the Atlas migrations, and load the committed `fixtures/`
  snapshot. Docker is required; no separate test DB preparation is needed.
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
- Single package / test: `go test ./internal/handlers/`,
  `go test ./internal/handlers/ -run TestListCourses -v`.
- **A test never creates the data it arranges: it reads it from `fixtures/`.**
  Every row there is named for the case it exists for, with a label:
  `{{ id "user_removed" }}` in the YAML (primary keys and references alike), and
  `fixtures.ID(fixtures.UserRemoved)` (`fixtures.ID32` for the contract's int32)
  in a test: every label has a generated `fixtures.Label` variable (the
  `snake_case` label in `PascalCase`), and `ID` takes a `Label`, not a string, so
  a typo or a removed row fails the build. `NewHarnessAs(t,
  fixtures.UserLearnerWithAGap)` signs a request in as that user. Add the row,
  then run `mise run gen:fixtures`; `fixtures.LabelOf(name)` is only for a label
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
  `accounts.WithClock`) and the fixtures are written relative to it.
  The fixtures were exported once from the legacy Rails set and have been edited
  by hand since; `mise run fixtures:import` skips every file that exists, and
  `mise run db:fixtures` and the e2e database load them through `cmd/fixtures`,
  the same loader the test containers use.
- A test owns no connection and cleans nothing up: it gets a transaction that is
  rolled back (`testsupport.NewClient*`, `NewHarness*`), the way a Rails test does.
  Jobs are worked in it with `rivertest.Worker` (`jobs.<Kind>Worker`), and what
  was queued or published is read back through the same transaction
  (`rivertest.RequireInsertedTx`, `testsupport.PublishedEvents`). Nothing runs a
  live River client or Watermill router against the shared database.
- Tests use `stretchr/testify` (`require`/`assert`); each handler test gets an
  explicit `sql.Tx` rolled back during cleanup, matching Rails transactional
  tests, and builds a `Server` through `testsupport`.
  The harness runs the real event publisher and job enqueuers over the test
  transaction, so a test asserts on what an operation left behind
  (`testsupport.PublishedEvents`, `testsupport.EnqueuedJobs`, DB rows), never on
  what was called, how often or with which arguments: no call-verification
  tests. Only the LLM and the Docker runner are mocks (`testsupport.NewMocks`,
  moq): set a mock's `…Func` for what the boundary answers, and read the prompt
  back through the response.

## Build

- `mise run build` — `build:api` (`bin/server`) + `build:worker` (`bin/worker`) +
  `build:web` (`pnpm build`). `mise run clean` removes `bin/` and `dist/`.

## Frontend (`src/`)

- TanStack Start (SSR, ADR-0008) + React 19 + TypeScript (`strict`) + Vite.
  Routing is file-based under `src/routes/`; `routeTree.gen.ts` is **generated**
  — don't hand-edit it. UI uses Mantine; data fetching uses TanStack Query hooks
  generated by hey-api into `src/client/`.
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
- `internal/fixtures/labels_gen.go` (fixture labels) ← `fixtures/*.yml` via `mise run gen:fixtures`.
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
- Match the surrounding package: doc comments on exported types/functions
  explaining *why* (the existing files are heavily commented on rationale —
  mirror that density on non-obvious logic, e.g. why an ordering has a specific
  tie-breaker), constructors return concrete types, errors wrapped with context.
- Return early on errors; don't swallow them. Keep DB access in handlers/ent and
  API types at the boundary — convert via `apiconv`. The business rules behind an operation
  live in a use-case module (next section), not in the handler.

### Use-case modules

The business rules behind an operation live in a package named for the concept
(`internal/<concept>/`), never in a handler body, so the HTTP handlers and any
other adapter (MCP) call the same code. There is no umbrella `usecases` or
`services` package. The rules below apply to every such module.

- **Package:** one concept, plural, lowercase (`accounts`, `blogposts`). Plural
  because ent generates a singular package per entity (`ent/blogpost`), and a
  file that needs both must not alias one; a mass noun (`progress`, `dashboard`)
  stays as it is.
- **Constructor:** `New(...)` returns the concrete type.
- **Interfaces:** the caller declares the small interface it needs (handlers
  and the MCP adapter each declare their own; see `handlers/collaborators.go`);
  the module exports none. The exception is a seam to the outside with a real
  second adapter, such as the LLM, the exercise runner, the mailer, or the
  course fetcher: the module that owns the boundary declares it, and `testsupport`
  mocks it.
- **Commands and results:** an operation takes a plain command named
  `<Verb><Noun>` (`UpdateUser`) and never an `internal/api` type (ADR-0001);
  adapters convert in both directions.
- **Errors:** domain errors named `ErrX` (`domain.ErrNotFound`,
  `domain.ErrConflict`, `accounts.ErrLastAdmin`, which is a conflict); the module
  maps ent errors to them with `domain.FromStore`, which keeps the cause in the
  chain, and each adapter maps them to its wire format (the HTTP one is
  `errorStatus`).
- **Separation:** writes and reads are separate modules or files (admin writes
  apart from the public read model); within a module the command, the result
  and the implementation each have their own file.
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
  back into a handler; the generated `Admin*` helpers pass.

## Conventional Commits & Releases

- Every commit follows [Conventional Commits](https://www.conventionalcommits.org/)
  (`<type>(<scope>?): <description>`); `commitlint` (commit-msg hook) enforces
  types. Commit messages are in English; the body explains *why*, not *what*.
- Releases are cut by `release-please` (config `release-please-config.json`,
  same type→bump rules as legacy). It maintains a release PR bumping
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

- Confirm you're working at the root (Go) vs `legacy/` (Rails) and follow the
  right guide.
- To change the API: edit TypeSpec → `mise run gen` → implement/adjust the Go
  handler and, if needed, the frontend client usage → run the narrowest Go test.
- After Go changes: `go build ./...`, then `go test ./...` (or a single package);
  DB-dependent packages prepare their own testcontainers automatically.
- After frontend changes: run `tsc -b` (`pnpm check`).
- Before finishing substantial work: `mise run lint` and the most relevant tests.
- Prefer minimal diffs; don't start broad refactors or "normalize" unrelated
  files. When unsure whether a file is generated, find its generator before
  editing (`grep` `mise/tasks.toml` / `go:generate` directives).

## Ops

- The Go stack deploys with `k8s/app-chart` (Helm, `k8s/Makefile`) from the
  images built by `Dockerfile` (server, worker, atlas) and `Dockerfile.web`
  (SSR); `release-please.yml` publishes them to GHCR on release. The sops
  secrets file and `terraform/` still live under `legacy/k8s` and
  `legacy/terraform`. If you touch any of it, call out the change clearly in
  handoff notes. Never commit secrets.

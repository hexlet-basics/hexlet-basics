# AGENTS Guide

This file (surfaced to Claude Code via the `CLAUDE.md` symlink) is for coding
agents working at the **root** of `hexlet-basics`. The root is the Go rewrite;
the original Rails app lives in `legacy/` with its own `AGENTS.md`/`CLAUDE.md`.
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
`docs/STACK.md` for the dependency plan and `docs/adr/0001`–`0008` for the
binding architecture decisions. **Read the relevant ADR before changing how a
subsystem works** — they are the source of truth for why each library was
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
- `make gen` runs the whole pipeline: `gen-spec` (TypeSpec → OpenAPI) →
  `gen-api` (`go generate ./internal/apigen`, ogen → `internal/api/`) →
  `gen-client` (openapi-ts → `src/client/`).
- **Never hand-edit `internal/api/oas_*_gen.go`, `api-spec/dist/`, or the
  generated TS client.** Change the `.tsp`, run `make gen`, commit the outputs.
- Adding a contract operation adds a method to the generated `api.Handler`
  interface. Until you implement it, the embedded `api.UnimplementedHandler`
  keeps it compiling as "not implemented" — so the build never breaks on a
  contract addition; you fill in the handler after.
- ogen 1.24 generates multipart uploads, but not OpenAPI `requestBody.encoding`.
  TypeSpec currently emits `encoding.file.contentType: "*/*"` for
  `HttpPart<File>`, so that upload operation is skipped in `ogen.yml` and
  temporarily handled outside the generated layer.
- **`@hey-api/openapi-ts` must stay on the `next` channel** while `typescript`
  is 7.x: the stable releases call the TypeScript compiler API (`ts.SyntaxKind`),
  which the native TS 7 package no longer exposes, so `gen-client` dies with
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
  client with `make gen-ent` (`go generate ./ent`) after editing a schema.
- Other `internal/` packages: `config/` (env config via `caarlos0/env`),
  `logging/` (`slog` + `tint`), `store/` (ent client construction), `slugs/`,
  `ids/`, `apiconv/`, `testsupport/`.

## Database & Migrations (atlas owns the schema)

- **atlas** owns the schema going forward (versioned SQL in `migrations/`,
  config in `atlas.hcl`). The baseline is the retired Rails schema. This is *not*
  ent auto-migrate and *not* Rails migrations.
- New schema change: edit `ent/schema`, then `make migrate-new NAME=...` to
  scaffold a migration to hand-author, then `atlas migrate hash`.
- Local PostgreSQL 18 runs in Docker on port **54330** (via
  `make services-start`) to avoid clashing with a `5432` DB. It is only for
  development; tests own disposable PostgreSQL 18 containers.

## Setup & Dev Commands (run from repo root)

- `make setup` — `go mod download` + `pnpm install` (root package.json is shared
  by the frontend and the api-spec tooling).
- `make services-start` / `make services-stop` — local Postgres in Docker;
  startup idempotently creates the development database.
- `make dev` — run API, async worker (both air live-reload), and Vite together.
- `make dev-api` — Go API only. `make dev-worker` — async worker only.
  `make dev-web` — Vite only.
- `make dev-spec` — watch TypeSpec and re-emit OpenAPI on change.
- Go-based tooling is pinned by `tool` directives in `go.mod` and invoked as
  `go tool <name>` (air, golangci-lint, ent, ogen, goverter, goi18n) — add new
  ones with `go get -tool`, never by hand. `mise.toml` pins only what cannot be
  a Go module: the go toolchain itself, kiota (.NET), and atlas (its public CLI
  module is frozen at a 2024-11 0.13.x snapshot that will not even install).
  Use `pnpm`, never `npm`/`npx` directly for JS deps.

## Codegen Commands

- `make gen` — full contract pipeline (spec → ogen server → hey-api client).
- `make gen-spec` / `make gen-api` / `make gen-client` — individual stages.
- `make gen-ent` — regenerate the ent ORM from `ent/schema`.
- `make gen-all` — `gen-ent` + `gen`. `make tidy` — `go mod tidy`.
- After changing a `.tsp`, ent schema, or deps, run the matching `gen` target
  and commit the regenerated outputs alongside the source change.

## Lint & Typecheck

- `make lint` — Go (`gofmt` check, `go vet`, `golangci-lint`) + frontend
  (`tsc -b`). The pre-push hook runs this; pushes need it clean.
- `make lint-fix` — `gofmt -w`, `golangci-lint --fix`.
- Go vet/build only: `go vet ./...`, `go build ./...`.

## Test

- `make test` — the Go suite (`go test ./...`).
- DB-dependent packages use testcontainers-go to start isolated PostgreSQL 18
  containers, apply the Atlas migrations, and load the committed `fixtures/`
  snapshot. Docker is required; no separate test DB preparation is needed.
- `make dev-prepare` prepares the development DB without replacing its data.
- Single package / test: `go test ./internal/handlers/`,
  `go test ./internal/handlers/ -run TestListCourses -v`.
- **Fixtures** in `fixtures/` are the starting data (converted from the legacy
  Rails set via `make fixtures-import`, which borrows the Rails runtime). Their
  ids are legacy crc32 values, so **assert on business facts (slug, order),
  never on raw ids** — see `courses_test.go` for the pattern.
- Tests use `stretchr/testify` (`require`/`assert`); each handler test gets an
  explicit `sql.Tx` rolled back during cleanup, matching Rails transactional
  tests, and builds a `Server` through `testsupport`.

## Build

- `make build` — `build-api` (`bin/server`) + `build-worker` (`bin/worker`) +
  `build-web` (`pnpm build`). `make clean` removes `bin/` and `dist/`.

## Frontend (`src/`)

- TanStack Start (SSR, ADR-0008) + React 19 + TypeScript (`strict`) + Vite.
  Routing is file-based under `src/routes/`; `routeTree.gen.ts` is **generated**
  — don't hand-edit it. UI uses Mantine; data fetching uses TanStack Query hooks
  generated by hey-api into `src/client/`.
- Call the API through the **generated client / Query hooks**, not hand-written
  `fetch`. If an endpoint is missing, add it in TypeSpec and regenerate.
- State: `zustand`; validation: `zod`; i18n: `i18next` + `react-i18next`
  (`src/locales/`, config via `i18next-cli`). Animation: `motion`.

## Generated / Do-Not-Edit Files (root)

Change the source, run the generator, commit the output. Never hand-edit or
`git restore` a generated file — regenerate it from current source.

- `api-spec/dist/**` (OpenAPI) ← `api-spec/*.tsp` via `make gen-spec`.
- `internal/api/oas_*_gen.go` (ogen server) ← OpenAPI via `make gen-api`.
- `internal/amocrm/generated/**` (Kiota client) ← the pinned amoCRM OpenAPI via
  `make gen-amocrm`.
- `src/client/**` (hey-api TS client + Query hooks) ← OpenAPI via `make gen-client`.
- `ent/**` (except `ent/schema/`) ← `ent/schema` via `make gen-ent`.
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
  API types at the boundary — convert via `apiconv`.

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

Single-context: `CONTEXT.md` + `docs/adr/` at the repo root. See
`docs/agents/domain.md`.

## High-Value Agent Workflow

- Confirm you're working at the root (Go) vs `legacy/` (Rails) and follow the
  right guide.
- To change the API: edit TypeSpec → `make gen` → implement/adjust the Go
  handler and, if needed, the frontend client usage → run the narrowest Go test.
- After Go changes: `go build ./...`, then `go test ./...` (or a single package);
  DB-dependent packages prepare their own testcontainers automatically.
- After frontend changes: run `tsc -b` (`pnpm check`).
- Before finishing substantial work: `make lint` and the most relevant tests.
- Prefer minimal diffs; don't start broad refactors or "normalize" unrelated
  files. When unsure whether a file is generated, find its generator before
  editing (`grep` the Makefile / `go:generate` directives).

## Ops

- Infra (`k8s/`, `terraform/`) currently lives under `legacy/`. If you touch it,
  call out the change clearly in handoff notes. Never commit secrets.

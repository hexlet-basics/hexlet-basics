[![check](https://github.com/hexlet-basics/hexlet-basics/actions/workflows/check.yml/badge.svg?event=push)](https://github.com/hexlet-basics/hexlet-basics/actions/workflows/check.yml)

# hexlet-basics

This repository is mid-migration from Rails to Go. The **root** is the new stack
— a Go API whose HTTP contract is generated from TypeSpec, plus a TanStack Start
(React 19 + TypeScript) frontend. The original **Rails app lives in `legacy/`**
with its own README and Makefile.

- Go API: `cmd/`, `internal/`, `ent/`
- API contract (source of truth): `api-spec/*.tsp` → `api-spec/dist/openapi.yaml`
- Frontend: `src/`
- Architecture decisions: [`docs/STACK.md`](./docs/STACK.md), [`docs/adr/`](./docs/adr)

## Setup

### Requirements

- `docker` (Go tests start disposable PostgreSQL containers)
- [`mise`](https://mise.jdx.dev) — provisions the pinned toolchain from
  `mise.toml` (`go`, `node`, `pnpm`, `golangci-lint`, `atlas`, `hk`, …), runs the
  project tasks (`mise tasks`) and the dev daemons; `mise run setup` runs
  `mise install` for you. Go codegen/dev tools (`air`, `ent`, `ogen`, …) come
  from `go.mod` `tool` directives — nothing to install.

### Run

1. Bootstrap everything from scratch with one task: toolchain, git hooks, `.env`,
   then `install` (dependencies) and `db:fixtures` (Postgres daemon, migrations,
   starting data from `fixtures/`) as separate tasks it composes:

   ```bash
   mise run setup
   ```

   PostgreSQL 18 is published on port `54330` so it does not conflict with a
   local database on `5432`. This service is only for development; Go tests
   create their own disposable databases.

   Dependencies alone: `mise run install` (Go modules + frontend/api-spec
   packages). The database steps are also available on their own: `mise run db:up` starts
   Postgres, `mise run db:prepare` creates and migrates the schema without
   touching data, and `mise run db:fixtures` loads the starting data (it replaces
   what the fixture tables hold, like `rails db:fixtures:load` did; without it the
   site shows nothing). Running `setup` again reloads the fixtures.

   Under mise, `DATABASE_URL` points at the local development DB and the API/Vite
   ports are dynamic (`mise env | grep PORT`): API from `3002`, Vite from `5173`,
   the next free port when one is taken.

3. Run the stack (Go API, asynchronous worker, and Vite frontend):

   ```bash
   mise run dev
   ```

   `mise run dev` starts them as background daemons (`mise daemons ls`,
   `mise daemons logs api`, `mise run dev:down` to stop). Run one process alone in
   the foreground with `mise run dev:api`, `mise run dev:worker`, or
   `mise run dev:web`.

## Code Generation

The HTTP contract is the source of truth (see
[ADR-0001](./docs/adr/0001-contract-first-pipeline-with-ogen.md)). Edit the
TypeSpec in `api-spec/*.tsp`, then regenerate — never hand-edit generated files.

```bash
mise run gen       # TypeSpec -> OpenAPI -> ogen Go server + hey-api TS client
mise run gen:ent   # regenerate the ent ORM from ent/schema
mise run gen:all   # gen-ent + gen
mise run dev:spec  # watch TypeSpec and re-emit OpenAPI on change
```

## Lint & Test

```bash
mise run lint          # gofmt + go vet + golangci-lint, and tsc
mise run lint:fix      # auto-fix Go formatting and lint findings

mise run test          # go test ./... with disposable PostgreSQL 18 containers
mise run db:prepare   # create/migrate dev DB without replacing development data
```

`mise run test` requires Docker. Each DB-dependent Go package starts an independent
PostgreSQL 18 through testcontainers-go, applies the [atlas](https://atlasgo.io)
migrations from `migrations/`, and loads `fixtures/` before its tests. Scaffold
a schema change with `mise run db:migrate-new <name>`.

## Build

```bash
mise run build  # builds bin/server, bin/worker, and the frontend bundle
```

## Releases

Commits follow [Conventional Commits](https://www.conventionalcommits.org/)
(enforced by `commitlint`). `release-please` maintains a release PR that bumps
`package.json` and `CHANGELOG.md`; merging it tags `vX.Y.Z`.

---

<!-- [![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics).
-->

## License

This project is licensed under the GNU General Public License v3.0.
See the [LICENSE](./LICENSE) file for details.

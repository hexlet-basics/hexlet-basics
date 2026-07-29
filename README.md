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

- `make`, `docker`
- [`mise`](https://mise.jdx.dev) — provisions the pinned toolchain from
  `mise.toml` (`go`, `golangci-lint`, `atlas`). Run `mise install`.
- `node` = 26.5.0 and [`pnpm`](https://pnpm.io)
- [`air`](https://github.com/air-verse/air) for API live-reload:
  `go install github.com/air-verse/air@latest`

### Run

1. Install dependencies (Go modules + frontend/api-spec packages):

   ```bash
   make setup
   ```

2. Start local Postgres (Docker):

   ```bash
   make services-start
   ```

   PostgreSQL 18 is published on port `54330` so it does not conflict with a
   local database on `5432`. This service is only for development; Go tests
   create their own disposable databases.

   Prepare the development schema:

   ```bash
   make dev-prepare
   ```

   The server reads `DATABASE_URL` / `ADDR` from the environment but defaults
   to the local development DB and `:3001`.

3. Run the stack (Go API on `:3001`, asynchronous worker, and Vite frontend):

   ```bash
   make dev
   ```

   Run one process alone with `make dev-api`, `make dev-worker`, or
   `make dev-web`.

## Code Generation

The HTTP contract is the source of truth (see
[ADR-0001](./docs/adr/0001-contract-first-pipeline-with-ogen.md)). Edit the
TypeSpec in `api-spec/*.tsp`, then regenerate — never hand-edit generated files.

```bash
make gen       # TypeSpec -> OpenAPI -> ogen Go server + hey-api TS client
make gen-ent   # regenerate the ent ORM from ent/schema
make gen-all   # gen-ent + gen
make dev-spec  # watch TypeSpec and re-emit OpenAPI on change
```

## Lint & Test

```bash
make lint          # gofmt + go vet + golangci-lint, and tsc
make lint-fix      # auto-fix Go formatting and lint findings

make test          # go test ./... with disposable PostgreSQL 18 containers
make dev-prepare   # create/migrate dev DB without replacing development data
```

`make test` requires Docker. Each DB-dependent Go package starts an independent
PostgreSQL 18 through testcontainers-go, applies the [atlas](https://atlasgo.io)
migrations from `migrations/`, and loads `fixtures/` before its tests. Scaffold
a schema change with `make migrate-new NAME=...`.

## Build

```bash
make build  # builds bin/server, bin/worker, and the frontend bundle
```

## Releases

Commits follow [Conventional Commits](https://www.conventionalcommits.org/)
(enforced by `commitlint`). `release-please` maintains a release PR that bumps
`package.json` and `CHANGELOG.md`; merging it tags `vX.Y.Z`.

---

<!-- [![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics).

See most active contributors on [hexlet-friends](https://friends.hexlet.io/). -->

## License

This project is licensed under the GNU General Public License v3.0.
See the [LICENSE](./LICENSE) file for details.

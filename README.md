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
  `mise.toml` (`go`, `golangci-lint`, `atlas`, `testfixtures`). Run `mise install`.
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

   Postgres is published on port `54330` so it does not conflict with a local
   database on `5432`. The server reads `DATABASE_URL` / `ADDR` from the
   environment but defaults to this local DB and `:3001`, so no config is needed
   for local dev.

3. Run the stack (Go API on `:3001` with live-reload + Vite frontend):

   ```bash
   make dev
   ```

   Run one side alone with `make dev-api` or `make dev-web`.

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
make lint          # gofmt + go vet + golangci-lint, and biome + tsc
make lint-fix      # auto-fix Go formatting and frontend lint

make test-prepare  # apply atlas migrations + load fixtures/ into the test DB
make test          # go test ./...
```

Run `make test-prepare` once before `make test` (CI does the same). The schema
is owned by [atlas](https://atlasgo.io) migrations in `migrations/`; scaffold a
change with `make migrate-new NAME=...`. Fixtures in `fixtures/` are the test
starting data.

## Build

```bash
make build  # go build -o bin/server ./cmd/server, and pnpm build
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

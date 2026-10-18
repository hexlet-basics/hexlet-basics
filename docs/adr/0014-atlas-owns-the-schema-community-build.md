# Atlas owns the schema, ent is decoupled from migrations (Community build)

The database schema is owned by **Atlas versioned migrations** (`migrations/`,
`atlas.hcl`), not by ent auto-migrate and not by Rails. The baseline migration
is the retired Rails schema captured once; every later change is hand-authored
SQL scaffolded with `atlas migrate new`. ent schemas in `ent/schema` only
describe the tables the Go code reads and writes, so an ent schema over an
existing legacy table needs no migration at all. Tests and development apply
the same migration directory, so the test schema cannot drift from production.

The CLI is the **Atlas Community build** (Apache-2.0, built only from the OSS
repository), installed by mise as `atlas-community`. The default Atlas binary is
distributed under the proprietary Atlas MSA and unlocks Pro features after
`atlas login`; we use none of them, and every command the workflow needs
(`migrate new`, `migrate hash`, `migrate apply`, `migrate diff`) is in the
Community build.

## Considered Options

- **ent auto-migrate** — rejected: the schema already exists (legacy data is
  migrated as-is, ADR-0002), and auto-migrate would own DDL we must control by
  hand.
- **goose / golang-migrate** — rejected: plain runners of hand-written SQL with
  no `migrate diff` planning from a desired schema.
- **Default (MSA) Atlas binary** — rejected: licence terms to accept for no
  feature we use.

## Consequences

- Community lacks `migrate lint/down/checkpoint/rebase`, `schema plan/lint`,
  and support for views, triggers, functions, sequences, and RLS. Needing any of
  these means revisiting this decision, not quietly switching binaries.
- Ariga deletes release binaries older than ~6 months and supports only the two
  latest minors, so the `atlas-community` pin in `mise.toml` must be bumped at
  least that often or `mise install` breaks for new machines and CI.
- The Atlas CLI cannot be a `go tool`: its public `cmd/atlas` module is frozen at
  a 0.13.x snapshot. Tests use the `ariga.io/atlas` Go library (Apache-2.0)
  instead.
- The pin lives in two places, `mise.toml` and the Dockerfile's
  `arigaio/atlas:<v>-community` image; `mise run lint:atlas-pin` (part of `lint`)
  fails when they disagree. The nightly workflow runs `mise run
  lint:atlas-retention`, which fails once the pinned binary is no longer
  downloadable from Ariga's release server, so the bump happens before a new
  machine or CI run breaks.

## Exit plan

The data-access stack has one vendor behind it twice over: ent and Atlas are both
maintained by Ariga. The exit is **sqlc** over the same Atlas-owned migrations
(the schema stays; only the query layer and the generated client change). A
working slice is on the `prototype/sqlc-slice` branch and is the evidence that
the move is feasible; it is not a plan to move.

Triggers, any one of which starts the exit:

1. No ent release for 12 months.
2. A Go toolchain bump we need breaks `go generate ./ent` and there is no tagged
   fix.
3. Atlas Community drops `migrate new`, `migrate hash`, `migrate apply` or
   `migrate diff`.

The first response to trigger 1 or 2 is not the exit: pin ent to its `master`
branch (a pseudo-version in `go.mod`), which usually carries the fix before a
tag. Only if that fails, or on trigger 3 with no other migration runner that
covers the same four commands, do we move to sqlc.

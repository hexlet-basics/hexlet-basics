# Data access: ent + atlas (community) versus alternatives

Resolves wayfinder ticket [#885](https://github.com/hexlet-basics/hexlet-basics/issues/885), part of map [#880](https://github.com/hexlet-basics/hexlet-basics/issues/880).
Criterion: codegen, off-the-shelf, best practice. Date of research: 2026-10-17.

Evidence tags: **[verified]** I ran it or counted it in this tree or read it from a
primary API response (GitHub API, repo file). **[read]** taken from an official
doc page fetched through a summarising tool, so wording is second-hand.
**[unverified]** general knowledge not re-checked against a primary source today.

## Short answer

Keep ent + atlas. I found no blocker. Two real risks, neither forces a switch:
(1) ent is 0.x, releases rarely, and is staffed by the same vendor as Atlas;
(2) Atlas Community is a deliberately cut-down build, pinned by hand with a
retention cliff. Details and an explicit list of what would change the verdict
are at the end.

## 1. Maintenance and licensing

### ent

| Fact | Value | Source |
|---|---|---|
| Licence | Apache-2.0, not archived | [verified] `gh api repos/ent/ent` |
| Last release | v0.14.6, 2026-03-23 | [verified] `gh api repos/ent/ent/releases` |
| Release cadence | v0.12.5 2023-11, v0.14.0 2024-07, v0.14.6 2026-03; tags v0.14.1-v0.14.5 exist between | [verified] releases + tags API |
| Master activity | last push 2026-10-14; 10 PRs merged since the last release; authors since then: masseelch (5), a8m (2), giautm, two outside contributors | [verified] commits and search API |
| Open issues / PRs | 629 open issues (GitHub count, includes PRs), 127 open PRs | [verified] repos API, search API |
| Who maintains it | README: "developed and maintained by the Atlas team" (Ariga), v1 roadmap in [ent#46](https://github.com/ent/ent/issues/46) | [verified] [README](https://github.com/ent/ent/blob/master/README.md) |

Reading: alive but slow, and a single vendor owns both ent and Atlas. Unreleased
master already carries a Go 1.27 fix
([ent commit 8037bc8, issue #4516](https://github.com/ent/ent/issues/4516),
closed): `entc/load` failed on Go 1.27 with x/tools older than v0.50. Our tree
runs `go 1.27.2` with `golang.org/x/tools v0.52.0` in `go.mod` (pulled by other
tools), so MVS already supplies the fix. **[verified]** `go generate ./ent` on
go1.27.2 here regenerates with zero diff. The lesson: a new Go release can break
ent codegen until Ariga ships, and we survive on a transitive-dependency
coincidence.

### Atlas and the community pin

- The repo `ariga/atlas` is Apache-2.0, not archived, last push 2026-10-11;
  releases v1.0.0 2025-12-24, v1.1.0 2026-02-05, v1.2.0 2026-04-10, v1.3.0
  2026-08-02 (roughly every 2 to 4 months). **[verified]** repos and releases API.
- Two CLI distributions. Default: "free to download and use under the terms of
  the Atlas MSA", `atlas login` unlocks Pro features, 30-day Pro trial.
  Community: "a minimal build distributed under the Apache 2.0 license", built
  from the GitHub repo. **[read]** [atlasgo.io/community-edition](https://atlasgo.io/community-edition).
  That page does not say when the split happened; I did not find a primary
  dated announcement, so the history of the split is unknown to me.
- Missing from Community, per the same page: `migrate checkpoint/down/rebase/edit/rm`,
  `schema plan/validate/lint/test/push`, migration lint and pre-migration
  checks, drift detection, views, triggers, functions, sequences, partitions,
  RLS, several drivers, and the integrations (K8s operator, Terraform, GitHub
  Actions, IDE plugins). What we use (`migrate new/hash/apply/diff`) is in both.
  **[read]**
- ADR-0014 and the earlier review add: binaries older than about 6 months are
  deleted and only the two latest minors are supported, so the pin must move
  every ~6 months; `cmd/atlas` on the Go proxy is frozen at 0.13.x so the CLI
  cannot be a `go tool`. The review verified the `v0.13.1` end of
  `go list -m -versions` ([review section 4](2026-09-24-go-stack-review.md));
  the retention rule is the maintainer's statement in
  [ariga/atlas#3296](https://github.com/ariga/atlas/issues/3296) (cited there,
  not re-read by me).
- Our footprint on Atlas: 10 migration files (1,829 lines of SQL), **[verified]**
  `wc -l migrations/*.sql`; `atlas.hcl` has one env; tests use the
  `ariga.io/atlas` v1.3.0 Go library (Apache-2.0), **[verified]** `go.mod`. The
  Go SDK listed as "missing" in Community refers to an integration product, not
  this library; the library works for us today.

The pin is a maintenance chore, not a licensing risk: the OSS repo is Apache-2.0
and the binary we run is built only from it. The risk is vendor-side: if Ariga
moved features out of the OSS tree, our fallback is plain SQL migrations (see
section 4), and because our migrations are hand-authored SQL the exit cost is
low (section 6).

## 2. Comparison on the criterion

| Option | Codegen | Off-the-shelf | Fit here | Verdict |
|---|---|---|---|---|
| **ent** (current) | Full: client, builders, predicates, typed edges, from Go schema; custom templates | Yes, mature | Already wired, 32 schemas | keep |
| **sqlc** | Generates Go from hand-written SQL; type-safe via the SQL analyser | Yes (MIT, release v1.31.1 2026-04-22, active pushes) **[verified]** API | Every query hand-written; no eager-load graph; dynamic filters only via `sqlc.narg`/COALESCE tricks **[read]** | strongest alternative, more hand code |
| **pgx + query builder** (squirrel, goqu, Jet) | None for squirrel and goqu (stringly columns); Jet generates a typed DSL from the live DB | squirrel last release 2023-03 **[verified]**; Jet v2.16.1 2026-10-06, Apache-2.0 **[verified]** | Jet is closest; no hooks, no templates, no graph loading **[unverified]** | no gain |
| **Bun** | None; struct tags and reflection **[read]** | Active (v1.3.0 2026-10-03, BSD-2-Clause) **[verified]** | Relations via `Relation("Posts")` tags **[read]**; column names are strings | worse type safety |
| **GORM** | Optional/none; reflection by default **[unverified]** (my attempt to read the GORM generics page returned 404) | Active (v1.31.2 2026-06-25, MIT) **[verified]** | `Preload` strings, runtime errors | worse type safety |
| **Plain SQL migrations** (goose v3.29.0 2026-10-11; golang-migrate v4.20.1 2026-09-09) | n/a | Both active **[verified]** release API; both report licence `NOASSERTION` in the GitHub API, i.e. check before adopting | Runners of hand SQL, no `migrate diff` | fallback for Atlas only |

Not benchmarked: I did not measure query performance for any option. All of
them issue one SQL statement per query; ent's cost is allocation in the scanner,
and the `LandingPage` ordering is a single ORDER BY in SQL (section 3). The
ticket's performance question is therefore answered structurally, not
empirically.

## 3. How the repo actually uses ent

Counts from `rg` over `internal/` and `cmd/`, excluding generated `internal/api`
and `internal/amocrm/generated`, production code only unless noted. **[verified]**
(`rg` patterns are approximate: they count call-site text, not types.)

| Pattern | Count |
|---|---|
| Files importing `hexletbasics/ent` | 71 production, 41 test files |
| `.With<Edge>(` eager loading | 63 hits in 31 files |
| `ent.Asc/Desc(` ordering | 27 in 17 files; `.Order(` 5 in 5 files |
| `.Where(` | 29 in 19 files |
| Custom `sql.Selector` or `sql.Expr` | 3 in 3 files (e.g. `landingpages.go:95`, `curriculum.go:49` order by `s.C(FieldID)`) |
| `GroupBy`/`Aggregate` | `blogposts/reader.go:179`, `progress/completion.go:127,169` |
| `OnConflict` (upsert feature) | 9 in 7 files |
| `CreateBulk` | 5 in 4 files |
| `.Create()` in production | 28 in 19 files |
| Hooks, interceptors, privacy | 0 (no `Hooks(`, `Policy(`, `Intercept`, `privacy` anywhere in `ent/schema` or `internal`) |
| Row locking (`ForUpdate`) | 0 |
| Raw SQL (`ExecContext/QueryContext/QueryRowContext`) | 0 in production; 2 files, both `internal/testsupport` (savepoints, outbox reads) |
| Transactions | one `Store.WithinTx` (`internal/store/store.go:66`) |

Shape of use: query builders with typed predicates and eager loading, SQL-side
ordering with `NULLS LAST` and tie-breakers
(`internal/landingpages/landingpages.go`, ~lines 79-100), a few aggregations,
upserts. The advanced ent features that create lock-in (hooks, privacy,
interceptors, GraphQL/entgql, edge-field magic) are unused. Lock-in is the
generated builder API itself plus two custom templates.

Codegen footprint **[verified]**: `ent/` generated code is 151,660 lines of Go;
schemas are 1,971 lines in 32 files; two hand-written templates,
`ent/template/adminput.tmpl` (166 lines) and `admincrud.tmpl` (121 lines),
generate `SetInput` methods that map contract input types (`api.*Input`) onto
builder setters, used by 6+ packages (`handlers/courses_admin.go`,
`blog_posts.go`, `staffroles`, ...). This is the piece with no ready-made
equivalent in sqlc/Bun/GORM: with sqlc we would write each admin update query,
or a template of our own.

forbidigo rules in `.golangci.yml` **[verified]**: production code may not call
`sql.DB/Tx/Conn/Stmt` Exec/Query methods ("use ent query builders instead of raw
SQL"); tests may not call `ent.*Client.Create/CreateBulk` (arrange from
`fixtures/`), nor `sql.Open`/`store.NewDB`. The linter rules are written against
ent types; a switch rewrites them.

## 4. Type safety and fit with the ogen/goverter pipeline

- ent: column names, predicates, edges and enums are generated Go identifiers, so
  a renamed field fails to compile. **[verified]** by use: compile errors are the
  normal feedback in this repo (e.g. `ByCourseField(course.FieldOrder, ...)`).
- goverter converts `*ent.X` to `api.X` in `internal/apiconv/apiconv.go` (128
  `goverter` annotation mentions) **[verified]**. goverter works on any Go
  structs, so the pipeline survives a data-layer change; but the converter
  interfaces name ent types, so every one is rewritten, and eager-loaded
  `Edges.X` access would become explicit joined structs.
- sqlc would give flat row structs per query, so loading Course + current Version
  becomes a hand-written join and a mapping step that ent's `WithCourse(...)`
  does for free. Type safety of sqlc is good (checked against the schema at
  generate time, by its design **[read]** from the docs); the cost is
  query volume, not safety.
- The contract pipeline (TypeSpec to ogen to hey-api) never touches the data
  layer; handlers only convert at the boundary (AGENTS.md), so a data-access
  switch does not touch the contract.

## 5. Transaction-based testing

The repo's test model is: one `sql.Tx` per test, rolled back in `t.Cleanup`;
`store.NewTxClient(tx)` binds ent to it; production code that opens a
transaction uses `Store.WithinTx`, and tests substitute a `savepointTransactor`
(`internal/testsupport/transactor.go`) that turns nested "transactions" into
SAVEPOINTs. **[verified]** by reading `internal/testsupport/testsupport.go`,
`transactor.go`, `store.go`.

What this says about alternatives:
- It works because ent's SQL driver accepts any `ExecQuerier` (`entsql.Conn{ExecQuerier: tx}`),
  with a small `txBoundDriver` wrapper so ent does not try to commit.
  That is ~10 lines of glue, not an ent feature we cannot replace.
- sqlc has the same seam by design: generated `DBTX` interface and
  `WithTx(tx)` **[read]** ([sqlc transactions doc](https://docs.sqlc.dev/en/latest/howto/transactions.html)),
  so the rollback-per-test harness ports almost unchanged.
- Bun and GORM have transaction handles too **[unverified]**; I did not check
  their savepoint behaviour.
- The weight of testing lives in `fixtures/` plus label codegen, which is
  independent of the ORM.

## 6. Migration risk and size of a switch (sqlc as the realistic target)

Estimated from the counts above; **estimate, not measured**.
- Rewrite roughly 71 production files and 41 test files that import ent, plus
  all `apiconv` converters, the two `adminput`/`admincrud` templates, and the
  forbidigo rules.
- About 150 `With*/Order/Where/GroupBy/OnConflict` call sites become hand-written
  SQL plus sqlc config; eager loading becomes joins or N queries.
- The 32 schemas in `ent/schema` would be replaced by SQL as source of truth.
  That is actually closer to our ADR-0014 state (Atlas owns the schema; ent
  schema over existing tables), but means the schema is described three ways
  (SQL, Go) during any partial migration.
- Because the migrations are plain SQL already, the Atlas-to-goose exit is cheap
  (ADR-0014 lists it as rejected for lacking `migrate diff`, not for being
  hard). Atlas Community is only needed for `diff`; we author by hand.
- Order of magnitude: weeks of full-time work, no user-visible gain.
  Bias going in was "keep"; the numbers do not argue against it.

## 7. Blockers looked for

| Candidate blocker | Finding |
|---|---|
| ent abandoned or relicensed | No: Apache-2.0, pushed 4 days ago, 10 merged PRs since last release |
| ent breaks on current Go | Go 1.27 broke `entc/load` once; fixed on master, covered in our tree by x/tools 0.52; regeneration is clean here |
| ent 0.x instability | Last tagged release 7 months old; 0.14.x line since 2024-07; the 0.14.x line has been stable since 2024-07, but v1 is unscheduled |
| Atlas licensing trap | Community is Apache-2.0 and what we run; default binary under MSA is avoided by ADR-0014 |
| Atlas pin rot | Real, mechanical: bump about every 6 months or `mise install` fails on new machines |
| Features we need but Community lacks | None today: views, triggers, functions, sequences, RLS and `migrate lint` unused. ADR-0014 says "revisit" if needed |
| ent cannot do a query we need | No raw SQL in production; `sql.Selector` escape hatch covers the three odd cases |

## Recommendation

**Clear keep** of ent + atlas-community, with these low-cost hardening items
(none needs a decision ticket):
1. Keep Go tool upgrades gated by `go generate ./ent` (the `gen:check` job
   already does this); note that ent codegen can lag new Go releases.
2. Add the 6-month Atlas pin bump to a recurring reminder or renovate rule
   (ADR-0014 already states the requirement; I found no automation).
3. Keep hooks, privacy and interceptors out of ent usage to preserve the
   cheap-exit property (they are unused now).

## Open questions

- Is a bus-factor concern about Ariga owning both ent and Atlas enough for the
  owner to want an exit plan written down, or is the section 6 estimate enough?
- When and why did Ariga split the Community build? I found no dated primary
  source; matters only if the owner wants to model the vendor's direction.
- Not measured: query cost of ent versus sqlc on `LandingPage`-style queries.
  Cheap to benchmark if the owner wants numbers, but I expect it to be dominated
  by Postgres.
- goose and golang-migrate show licence `NOASSERTION` in the GitHub API; verify
  the actual LICENSE text if the migration fallback is ever exercised.
- Bun, GORM, squirrel and Jet rows rest partly on unverified knowledge; do a
  proper read of them only if the keep verdict is challenged.

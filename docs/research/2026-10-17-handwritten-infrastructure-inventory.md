# Hand-written infrastructure in `internal/` and `cmd/`: inventory

**Question** (issue #886, map #880). Which hand-written infrastructure (not
business logic) could an established library or generator replace? For each:
what it does, size, the off-the-shelf option, cost of switching.

**Date.** 2026-10-17. **Status.** Inventory only. No recommendations and no code
changed; "cost" is a statement of what would have to move, not advice.

## How this was checked

- **Ran:** `wc -l` over non-test, non-generated Go files (excluding `*_gen.go`,
  `*.gen.go`, `internal/api/`, `*/generated/`); `gh api repos/<owner>/<repo>`
  for licence, `archived` and `pushed_at` of every library named below (queried
  2026-10-17); `rg`/`grep` for usages (health endpoints, `do.Shutdown`,
  `go-pkgz`, `jwt`); `ls` of module-cache sources (`$GOMODCACHE`) to confirm
  that a library ships a given package or method.
- **Read, not run:** every hand-written file cited below. No test, build or
  benchmark was run; nothing was prototyped, so every "cost" is an estimate
  from reading call sites, not a measured migration.
- **Not checked:** behaviour of any alternative under our constraints (for
  example whether `scs` pgx stores can join a rolled-back test transaction was
  taken from the repo's own comment at `internal/sessions/store.go:10-17`, not
  re-tested). Library maintenance is judged only by `pushed_at`, `archived` and
  licence; releases and issue health were not read.
- Sizes are non-blank, non-generated physical lines (`wc -l`), comments
  included; the codebase is heavily commented (AGENTS.md "Go Style"), so code-only
  counts are smaller.
- Whole non-generated Go under `cmd/` + `internal/` is about 17,150 lines
  (same `wc -l` filter), so the infrastructure below is a minority of the code.

## Summary table

Library facts from `gh api repos/...` on 2026-10-17: licence / archived /
last push.

| # | Item | Where (lines) | What it does | Off-the-shelf option (primary repo) | Switch cost |
|---|---|---|---|---|---|
| 1 | DI wiring | `internal/di/*.go` 885 (`server_package.go` 481, `worker_package.go` 224, `common_package.go` 119, `container.go` 61) | `samber/do` providers, one per service, each resolving its deps by hand | already `samber/do` (MIT, active). Alternatives: `uber-go/fx` (MIT, pushed 2025-12-27), `uber-go/dig` (MIT, 2025-05-13), `google/wire` (Apache-2.0, **archived**) , or plain constructor calls | medium to high, see note 1 |
| 2 | API error mapping | `internal/handlers/errors.go` 194 | status mapping, RFC 9457 body, 5xx log + Sentry, 404/405 localization | no library for the ogen seams; pieces exist (`ogenerrors`, `sentry-go`, `slog-sentry`) | low value, moderate effort |
| 3 | Session store glue | `internal/sessions/store.go` 81, `session.go` 215 | scs `Store` over ent; cookie policy; rotation, idle timeout | `alexedwards/scs` `postgresstore` / `pgxstore` (MIT) | blocked by transaction-rollback test design |
| 4 | apiconv hand-written part | `internal/apiconv/apiconv.go` 845 (384 code lines, 122 `goverter:` directives, 62 funcs) | goverter interface plus bridges for ogen `Opt*`/`Nil*` types | already `jmattheis/goverter` (MIT, active); alternatives are reflection mappers | small; mostly inherent |
| 5 | Config loading | `internal/config/config.go` 185 | `caarlos0/env` struct, `godotenv`, production safety checks | already `caarlos0/env` (MIT); alternatives `sethvargo/go-envconfig`, `knadh/koanf`, `spf13/viper` | very low value |
| 6 | Logging | `internal/logging/logging.go` 32 | tint handler plus TTY/`NO_COLOR` detection | `lmittmann/tint` (MIT) README documents `go-isatty` for this | trivial |
| 7 | Process lifecycle | `cmd/server/runtime.go` 102 + `main.go` 95; `cmd/worker/runtime.go` 165 + `main.go` 95 | errgroup supervisors, signal handling, ordered start/stop, `closeResources` | `do` has `ShutdownOnSignals`, `Shutdown`, `HealthCheck` (unused here); `oklog/run`, `sourcegraph/conc`, `thejerf/suture` | medium; ADR-0009/0010 constraints |
| 8 | Health/readiness endpoint | none found | (absent) | `alexliesenfeld/health`, `hellofresh/health-go`; `do.HealthCheck` | new feature, not a replacement |
| 9 | Event glue (Watermill) | `internal/events/*.go` 422 (`runtime.go` 165, `events.go` 169, `publisher.go` 88) | outbox publisher per tx, router, retry/poison middleware, observer | already Watermill (MIT) + `watermill-sql`; its `components/forwarder` is an outbox forwarder; River as the only queue | high if dropping Watermill (not prototyped) |
| 10 | Job glue (River) | `internal/jobs/*.go` 714 (`jobs.go` 127, `enqueuer.go` 55, `errors.go` 66) | queues, `Module` registration, tx-bound `Enqueuer` | already River (**MPL-2.0**, active) | low; mostly thin |
| 11 | Exercise runner | `internal/exerciserunner/runner.go` 452 | Docker run of untrusted code: limits, tar copy, capped output, timeout exit code | `moby/moby` client already used; no library for this logic | no off-the-shelf equivalent found |
| 12 | Mailer | `internal/mailer/mailer.go` 120 | SES v2 client at Postbox, log adapter | `aws-sdk-go-v2` sesv2 already; `wneessen/go-mail` (MIT) for SMTP | low value |
| 13 | Asset store | `internal/assetstore/assetstore.go` 211 | MIME policy, blob write plus DB row, compensation, URLs | already `gocloud.dev/blob` | no equivalent for the policy part |
| 14 | HTTP router and middleware | `internal/handlers/router.go` 111, `security.go` 168, `client_ip.go` 30, `guest_cookies.go` 65, `localization.go` `Middleware` ~15 | ServeMux dispatch, SSE flush wrapper, realip, auth pipeline | `go-pkgz/rest` already used for realip; ogen generated router | low |
| 15 | Tx adapter | `internal/store/store.go` 102 | one `sql.Tx` bound to ent and River/Watermill | ent `Tx` and River pgx drivers cover parts; this bridges both on `database/sql` | high, structural |
| 16 | Tokens and small helpers | `emailtokens` 142, `paging` 30, `inputconv` 39, `ids` 9, `htmltext` 27 | JWT email tokens, paging window, `Opt*` to pointer | `golang-jwt/jwt` (MIT) already used; the rest are tiny | not worth listing further |

## Per-item notes

### 1. DI (`internal/di`, 885 lines)

- Every service is a `do.Lazy[T]` closure that calls `do.Invoke` per dependency
  and then the plain constructor, e.g. `internal/di/common_package.go:30-123`
  and `internal/di/server_package.go:140-481`. The repetition is the bulk of the
  885 lines.
- `handlerServices` (`server_package.go`, struct near the top) uses `do:""`
  tags with `do.InvokeStruct` (also `container.go:44-51`). The file comment says
  `handlers.Deps` holds interfaces "which the container cannot tell apart by type
  alone".
- Named `http.Handler` bindings (`server_package.go:276-319` region) exist to
  separate the router from the Sentry-wrapped handler.
- Alternatives (all verified to exist, licence MIT unless noted): `uber-go/fx`
  and `uber-go/dig` are reflection-based and add their own lifecycle;
  `google/wire` is Apache-2.0 and **archived** (`gh api`: `archived: true`,
  last push 2025-08-22), so it is not an active option. Manual wiring is the
  non-library option.
- Cost: AGENTS.md ("Architecture: Handlers, DI, and Data") makes `do` the
  documented convention; `di.NewServer`/`di.NewWorker` are called from `cmd/`.
  Any change touches all four `di` files and the process `main`s. Not
  prototyped.
- `do` v2.1.0 ships `ShutdownOnSignals`, `Shutdown*` and `HealthCheck*`
  (`root_scope.go:164-321` in the module cache). `rg 'do\.Shutdown|ShutdownOnSignals|HealthCheck'`
  over the repo found no use, so these features of the library already in the
  tree are not used (see item 7).

### 2. `handlers/errors.go` (194 lines)

- `errorStatus` (`errors.go:110`) maps `statusError`, ogen errors, `context.Canceled`
  to 499, ent not-found and constraint errors, and domain errors to HTTP
  statuses. `Write` (`:90`) writes `application/problem+json` through `jx`.
  `report` (`:137`) logs 5xx with trace ids and captures to Sentry. `NewNotFoundHandler`
  (`:177`) and `NewMethodNotAllowedHandler` (`:185`) localize 404/405.
- No library implements ogen's `NewError` / `WithErrorHandler` seams (both are
  ogen-specific, per the file comment `errors.go:16-27`). Reusable parts:
  `getsentry/sentry-go` (MIT) and `samber/slog-sentry` (MIT, 54 stars) for
  forwarding log records to Sentry; `report` currently does both logging and
  `hub.CaptureException` by hand. `samber/oops` is already used.
- Cost: whichever parts are replaced must keep the generated
  `api.ProblemDetails` shape. The central mapping itself would remain.

### 3. Sessions (`internal/sessions`, 296 lines)

- `store.go` implements scs's `Store` (`FindCtx`, `CommitCtx`, `DeleteCtx`
  at `:29`, `:45`, `:57`) over ent, plus `deleteExpired` (`:64`) and three
  context-free methods the interface requires. `session.go` configures scs
  (`HashTokenInStore`, cookie policy, `session.go:97`) and adds lifetime,
  idle-timeout and password-epoch logic.
- The repo's own reasons for not using scs's `postgresstore`
  (`store.go:10-17`): it takes a `*sql.DB` so it cannot join a test's rolled-back
  transaction, and it starts a cleanup goroutine that ADR-0010 forbids. Both
  `postgresstore` and `pgxstore` exist as modules in the cache
  (`alexedwards/scs/postgresstore@v0.0.0-20251002...`,
  `pgxstore@...`). I did not test whether either can be given a transaction.
- `alexedwards/scs`: MIT, last push 2025-11-20. `gorilla/sessions` (BSD-3, last
  push 2024-08-20) is the other well-known option; it is not a drop-in for
  server-side hashed tokens.
- Cost: the store file is 81 lines; the larger part (`session.go`) is policy
  (ADR-0017), not plumbing.

### 4. `apiconv` hand-written part (845 lines)

- Already generated by goverter (`apiconv.go:26` `go:generate`, `:30`
  `goverter:converter`); the generated output is `apiconv.gen.go` (833
  lines, excluded from the counts above). The hand-written file is mostly
  `goverter:` directives (122 lines) and 62 small bridge functions
  (`grep -c`), e.g. the `Nil*FromPtr` family listed at `apiconv.go:34-56`.
- Many bridges exist only because ogen represents optionals as `Opt*`/`Nil*`
  wrappers (package comment, `apiconv.go:1-6`).
- Alternative: reflection mappers (not evaluated). Changing the TypeSpec
  optionality rules or ogen options would shrink the bridges; I did not
  check which ogen generator option would do so.
- Cost: low to moderate; this layer is already the "generator" answer.

### 5. Config (185 lines)

- `Load` (`config.go:154`) does `godotenv.Load`, `env.Parse` and
  `validateProduction` (`:171`). The remaining ~150 lines are the struct with
  `env:` tags (`config.go:20-140`).
- Already `caarlos0/env` v11 (MIT, pushed 2026-10-07) and `joho/godotenv`
  (MIT). Alternatives that bundle files/env/validation: `knadh/koanf`
  (MIT), `spf13/viper` (MIT), `sethvargo/go-envconfig` (Apache-2.0).
- Cost: the struct is the contract with Helm values (`k8s/`, not read for
  this note). Unchecked.

### 6. Logging (32 lines)

- `logging.New` builds a tint text handler (`logging.go:12`);
  `colorEnabled` (`:23`) implements `NO_COLOR`, `FORCE_COLOR` and a TTY check
  with `os.Stat`. tint's README (v1.2.1, line 102-108) shows the same TTY check
  with `mattn/go-isatty` (MIT, pushed 2026-07-23). `FORCE_COLOR` is the part the
  README example lacks.
- Cost: trivial.

### 7. Lifecycle (`cmd/*`, 457 lines)

- Each process has a hand-written `errgroup` supervisor: `cmd/server/runtime.go:36`
  (`run`), `:72` (`shutdown`); `cmd/worker/runtime.go:41`, `:118`, with a
  `runtimeState` of atomics. The worker's ordering (River, then Watermill, reverse
  on stop; "jobsStartSettled" race guard) is documented inline
  (`cmd/worker/runtime.go:41-70`).
- `closeResources` in both `main.go` files (`cmd/server/main.go:75`,
  `cmd/worker/main.go:75`) closes the bucket, DB, OTel and Sentry in a fixed
  order by hand, although `do` already supports `Shutdown` interfaces
  (`di_lifecycle.go:62-125`); the comment at `common_package.go:71-72` says
  `blob.Bucket` has `Close` and not a do Shutdowner.
- Options: `oklog/run` (Apache-2.0, pushed 2025-06-27), `sourcegraph/conc`
  (MIT, 2026-09-01), `thejerf/suture` (MIT, 2026-07-22). I did not evaluate
  whether any of them expresses the readiness gating (wait for Watermill
  `Running()`) that the code needs.
- Cost: ADR-0009 (superseded) and ADR-0010 choose errgroup explicitly; a switch
  needs an ADR and keeps the readiness logic regardless.

### 8. Health endpoint (absent)

- `rg -i 'healthz|/health|readyz'` over non-generated Go found nothing.
  Listed because the issue names "health/lifecycle"; Kubernetes probes
  configuration in `k8s/` was not read. Libraries exist (`alexliesenfeld/health`
  MIT, pushed 2026-06-15; `hellofresh/health-go` Apache-2.0, 2026-04-17), as does
  `do.HealthCheck`.

### 9. Events (`internal/events`, 422 lines)

- `Publisher.Publish` (`publisher.go:49`) builds a Watermill event bus per call
  on the caller's `sql.Tx` (see `busConfig` comment, `:73`). `Runtime`
  (`runtime.go:33`) builds the router, SQL subscribers (`newSubscriber`,
  `:134`), a poison queue and Retry/Recoverer middleware
  (`deliveryMiddleware`, `:104`). Events are a typed enum in `events.go`.
- The libraries are already Watermill (MIT, pushed 2026-08-25). Its
  `components/forwarder` package exists in v1.5.3 (module cache), which is the
  library's own outbox forwarder; I did not read whether it fits a
  Postgres-table-as-topic setup like this one.
- ADR-0004 (`docs/adr/0004-events-via-watermill-jobs-via-river.md`) keeps events
  and jobs as separate mechanisms; I did not check it for a River-only
  alternative. Dropping Watermill would also remove `runtime.go`'s router and the worker's second
  actor (item 7).
- Cost: high; events are part of the ADR-0004 design.

### 10. Jobs (`internal/jobs`, 714 lines)

- Most of the 714 lines are per-kind job Args and workers (business); the
  infrastructure is `jobs.go` (queue map `:38`, `Module`/`module` `:59`,
  `NewWorkerClient` `:88`), `enqueuer.go` (tx-bound batch insert via
  `InsertManyTx`) and `errors.go`.
- River: MPL-2.0 (file-level copyleft; recorded here as a licence fact), last
  push 2026-10-17, plus `otelriver` from `rivercontrib` already used.
- Cost: low; these files wrap River rather than replace it.

### 11. Exercise runner (452 lines)

- `Docker.run` (`runner.go:255`) creates a container with `hostConfig`
  limits (`:369`), copies the submission in as a tar (`tarFile`, `:398`),
  detects timeout through coreutils exit status 124 (`:41`), and caps output
  (`cappedWriter`). Governed by ADR-0013 and
  `docs/research/2026-08-04-untrusted-code-runner-in-go.md`, which I did not
  re-read for this note.
- The Docker client is already the official `moby/moby` module. I found no
  library that implements "run learner code with these limits" as a unit;
  `testcontainers-go` (MIT) is in `go.mod` but oriented to tests. Not evaluated
  beyond that.

### 12. Mailer (120 lines)

- `mailer.go:38` picks the Postbox SES v2 client or a log adapter. It is
  already a thin wrapper over `aws-sdk-go-v2/service/sesv2`. Postbox's SMTP
  option and `wneessen/go-mail` (MIT, pushed 2026-10-13) would be a different
  transport; not checked against Postbox docs.

### 13. Asset store (211 lines)

- `Put` (`assetstore.go:96`) enforces `MaxUploadBytes` 15 MiB (`:25`),
  MIME policy, writes the blob, persists the attachment row and runs
  `compensate` (`:203`) on failure. `gocloud.dev/blob` (ADR-0005) already
  provides the storage seam; the hand-written part is the policy and
  compensation, with no library equivalent found.

### 14. HTTP router and middleware (about 390 lines)

- `NewRouter` (`router.go`) uses `http.ServeMux` for four hand-routed paths and
  delegates the rest to ogen; `StreamResponses` plus `flushingWriter` (~45
  lines) disable the write deadline and flush for AI streaming. `WithClientIP`
  already uses `go-pkgz/rest/realip` (`client_ip.go:7`; MIT, 63 stars,
  pushed 2026-08-19). `localization.Translator.Middleware`
  (`localization.go:64`) reads `Accept-Language` into the context.
- Whether ogen can generate the streaming and webhook routes was not checked.

### 15. Transaction adapter (`internal/store`, 102 lines)

- `Store.WithinTx` (`store.go:66`) gives one callback both `*sql.Tx` and an ent
  client bound to it via `txBoundDriver` (`:60`, `:90`), so River/Watermill
  inserts and ent writes commit together; tests depend on it (AGENTS.md "Test").
  It is the structural reason for items 3, 9 and 10 being shaped as they are.
  No library found that bridges these three; not searched exhaustively.

### 16. Small helpers

- `emailtokens` (142 lines) uses `golang-jwt/jwt` v5 (MIT, active) directly
  (`emailtokens.go:78-100`); `paging`, `inputconv`, `ids`, `htmltext` total
  about 140 lines. Not analysed further.

## Not in scope or not measured

- `internal/handlers` business handlers (3,276 lines), `progress`,
  `courseloader`, `testsupport` (1,205 lines of test harness) and `fixtures`
  (252 lines around `testfixtures`, MIT) were not inventoried item by item.
- `legacy/`, `src/`, k8s and ops are outside the map (#880).
- No claim here about performance, bundle size or runtime behaviour was
  measured.

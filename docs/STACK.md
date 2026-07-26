# Go backend — dependency plan of record

Finalized via a grilling session on 2026-07-26, starting from a reference
`go.mod` the user pasted and trimming to what has a real consumer at parity.
Architecture decisions live in [`docs/adr/`](./adr) (0001–0008). Frontend
rollout is planned in [`docs/FRONTEND_PLAN.md`](./FRONTEND_PLAN.md).

**Cutover:** hard cutover from Rails at parity (no side-by-side). Backward-compat
kept on **bcrypt passwords** and **URL routes** only (ADR-0002).

## Dependencies to use (near-term)

| Area | Libraries | Notes |
|---|---|---|
| API / HTTP | `ogen` (+ `go-faster/jx`, `go-faster/errors`), `rs/cors` | ADR-0001. Replaces `oapi-codegen`+`kin-openapi`+`echo`. Generates validation in-code. |
| DB | `entgo.io/ent` + `ariga.io/atlas`, driver `jackc/pgx/v5` | ent replaces hand-written pgx SQL. |
| Jobs | `river` (+ `riverpgxv5`, `rivertype`, `otelriver`) | ADR-0004. Replaces Solid Queue. |
| Events | `watermill` (+ `watermill-sql`) | ADR-0004. Domain-event fan-out. |
| Auth | `go-pkgz/auth`, `golang-jwt`, `go-crypt`, **`go-webauthn/webauthn`** | ADR-0003. webauthn ADDED (was missing). JWT sessions. |
| Email | `aws-sdk-go-v2` (+ `config`, `credentials`) + `service/ses` | ADR-0006. Talks to **Yandex Postbox** SES-compat API (not AWS). |
| Assets | **`gocloud.dev/blob`** | ADR-0005. ADDED. s3blob (prod) / fileblob (dev) + ent attachment table. |
| i18n | `nicksnyder/go-i18n/v2` | Backend-emitted strings only. Content localization = ent data; UI i18n = React i18next. |
| Config / DI / utils | `caarlos0/env`, `samber/do`, `samber/lo`, `samber/oops`, `oklog/ulid`, `gosimple/slug` | `caarlos0/env` over `spf13/viper`: config is env-only (12-factor), no files/remote/flags — struct-tag parsing fits, viper's weight didn't. |
| Logs / errors | `slog` + `lmittmann/tint`, **`getsentry/sentry-go`** | ADR-0007. sentry-go ADDED (OTel ≠ error tracking). |
| Tracing (lean) | `go.opentelemetry.io/otel` + `sdk` + `trace` | ogen emits OTel natively. |
| Tests | `stretchr/testify`, `go-txdb`, `testfixtures` | Fixtures-based, like legacy. |

## Added beyond the reference list

- `go-webauthn/webauthn` — passkeys (ADR-0003).
- `gocloud.dev/blob` — ActiveStorage-like asset storage (ADR-0005).
- `getsentry/sentry-go` — error tracking parity (ADR-0007).

## Dropped (no consumer / superseded)

`osteele/liquid` · `standard-webhooks` · `tink-crypto/tink-go` · `doyensec/safeurl`
· `wneessen/go-mail` (SES API instead of SMTP) · `emersion/go-msgauth` (Postbox
does DKIM) · `robbiet480/go.sns` (Postbox monitoring ≠ AWS SNS).

## Deferred (add when there is a concrete need)

- `preslavrachev/gomjml`, `k3a/html2text` — nicer emails, post-parity.
- Full OTel exporters (`otlptrace*`, `otlpmetric*`), `prometheus/client_golang`,
  `exporters/prometheus`, `contrib/instrumentation/runtime` — until a
  metrics/trace receiver exists.

## No library needed

- **Content build & solution runner** — river jobs shelling to the `docker` CLI
  via `os/exec` (mirrors legacy `Open3`). No Docker SDK, no git client, no runner
  microservice: course content ships as a Docker image (`ghcr.io/.../exercises-*`).
- **SMS** — SMSC.ru HTTP API via `net/http`. Legacy is a stub; low priority.

## Rollout

`go.mod` advances **per phase** — unimported deps get removed by `go mod tidy`,
so libraries land with the code that uses them.

0. **Foundation** — caarlos0/env, slog+tint, samber/do+lo+oops, ulid, slug, lean otel.
1. **Persistence** — ent + atlas; rewrite `internal/handlers/courses.go`; testify + go-txdb + testfixtures.
2. **API layer** — ogen + rs/cors; drop echo/oapi-codegen/kin-openapi.
3. **Auth** — go-crypt + go-pkgz/auth + go-webauthn + golang-jwt.
4. **Jobs & events** — river + watermill.
5. **Email** — service/ses → Postbox.
6. **Assets & i18n** — gocloud.dev/blob; go-i18n.

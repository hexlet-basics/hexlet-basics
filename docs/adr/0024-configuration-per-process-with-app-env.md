---
status: accepted
---

# Configuration: one root per process, `APP_ENV` policy, no process-wide config in constructors

Part of the backend clean-up (spec #905, ticket #918). Amends the configuration
parts of [ADR-0010](0010-separate-http-and-async-processes.md) only where they say both
processes parse one shared struct.

(Number 0023 is reserved for the lifecycle ADR of the same spec.)

**Problem.** `config.Load` parsed every setting in both processes, so the worker
needed the guest cookie secret it never reads; deploy-critical values
(`DATABASE_URL`, hosts, bucket) had development defaults that a forgotten
variable would silently put in production; the production policy hung on a
Sentry label (`SENTRY_ENVIRONMENT`); and `*config.Config` was passed to
constructors, so a component could read anything.

**Decision.**

- `APP_ENV` is `development` (default), `production` or `test`; anything else
  stops the start. It, not a Sentry label, drives environment policy: in
  production a secret equal to the public development value is refused, and the
  worker requires a Postbox key. Sentry's environment is `APP_ENV`;
  `SENTRY_ENVIRONMENT` is no longer read.
- Two roots, `config.LoadServer()` and `config.LoadWorker()`, each embedding
  `config.Common` (database, bucket, site and public URLs, email token secret,
  OpenAI, Sentry) and adding only its own settings. Server only: `ADDR`,
  `APP_HOST`, `GUEST_COOKIE_SECRET`, OAuth, MCP, exercise runner, webhook secret,
  book key. Worker only: GitHub token, amoCRM, mail, audit retention.
- Required, with no code default: `DATABASE_URL`, `PUBLIC_URL`, `SITE_URL`,
  `BLOB_BUCKET_URL`, `EMAIL_TOKEN_SECRET`, plus `APP_HOST` and
  `GUEST_COOKIE_SECRET` for the server. A missing one fails the start naming
  every missing variable at once (`caarlos0/env` aggregates them).
- Only `internal/di` and `cmd` hold a root. A constructor takes the section it
  reads or scalars (`handlers.Settings`, `telemetry.NewSentryClient(dsn, env,
  release)`). `mise run lint:config` rejects `*config.Server` / `*config.Worker`
  outside `di`, `config` and tests, and `os.Getenv` / `LookupEnv` / `Environ`
  outside `internal/config` (and `testsupport`).
- `.env` is no longer read by `Load*`: `config.LoadDotenv()` is called from
  `cmd/*/main` and reads it only when `APP_ENV` is unset or `development`, so
  production and tests never read a dotenv file. mise's `[env]` supplies the
  dev stack's `DATABASE_URL`, `PUBLIC_URL`, `SITE_URL`, `APP_HOST` and
  `BLOB_BUCKET_URL`, so an existing checkout keeps working.

**Not decided here.** Validator tags, a redacting `Secret` type and the
`origin`/`notdevsecret` validators are separate tickets of the spec; the
hand-written `validate*` functions stay until then.

**Ops note (the chart is outside this rewrite).** Before deploying, give both
workloads `APP_ENV=production` (replacing `SENTRY_ENVIRONMENT=production`), and
make sure `DATABASE_URL`, `PUBLIC_URL`, `SITE_URL`, `BLOB_BUCKET_URL` and
`EMAIL_TOKEN_SECRET` are set on both; the server additionally needs `APP_HOST`
and `GUEST_COOKIE_SECRET`. The worker no longer needs `GUEST_COOKIE_SECRET`,
`APP_HOST` or the server-only settings. A missing variable now fails the pod at
start instead of running on a development value.

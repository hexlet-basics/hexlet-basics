# Design

## Context

See proposal.md for motivation; the specs under `specs/` state the behaviour.

- `internal/handlers/auth.go` already owns sign-in: `CreateSession` publishes
  `UserSignedIn`, `issueCookies` sets the go-pkgz/auth JWT cookie, and
  `mergeGuestProgress` credits guest progress and clears its cookie.
- River runs in the worker process (ADR-0010); `internal/jobs` shows the
  pattern, including `UniqueOpts` in `ExerciseLoaderArgs`.
- `internal/localization` resolves the request locale and holds go-i18n
  bundles; its fallback is currently `en`.
- `golang-jwt/jwt/v5` is already a direct dependency; `aws-sdk-go-v2` is only
  indirect.
- Links built today use `config.AppHost` (HTTPS assumed).

## Goals / Non-Goals

**Goals:**
- One mail path that both flows, and later ones (book request, phone-less
  notifications), reuse.
- Sign-in through a link shares one code path with password sign-in.

**Non-Goals:**
- Changing the locale set or the default fallback of the rest of the site
  (the `en` decision). Account emails pick `ru` on their own when the request
  locale is not `ru`/`es`.
- A text/plain part, MJML templates, bounce handling (ADR-0006 defers them).

## Decisions

**Email tokens are JWTs signed with a dedicated secret.** Each token carries
the purpose (`magic_link` / `password_reset`), the user id, an expiry of
fifteen minutes, and a fingerprint: a hash of the email for a Magic Link, a hash
of the current `password_digest` for a Password Reset. Verifying reloads the
User and compares the fingerprint, which gives "dies on email change" and "dies
after one reset" without storage. `golang-jwt/jwt/v5` handles signing and
expiry, so no HMAC or encoding is hand-written; the purpose claim stops a
Magic Link from being replayed as a reset token. The secret is
`EMAIL_TOKEN_SECRET`, separate from `JWT_SECRET` so either can rotate alone.
*Alternative:* a token table — rejected: a migration and a sweeper for no
behaviour the specs need. *Alternative:* reusing the guest-cookie HMAC codec —
rejected: it is hand-written and cookie-shaped.

**A `Mailer` port with three adapters.** `Send(ctx, Message)` where `Message`
is recipient, subject, HTML body. Adapters: Postbox via `sesv2` with the
endpoint overridden (production), log via `slog` (development, selected when
Postbox credentials are absent), and an in-memory recorder (tests). The choice
is made in DI, not in the job.

**Rendering happens in the worker; the job carries intent, not HTML.** The API
enqueues `AccountEmailArgs{Kind, UserID, Locale}`; the worker loads the User,
issues the token, renders `html/template` with go-i18n strings and sends.
Issuing the token at send time means a retried job never mails an expired link,
and the job table never holds a usable credential.
*Alternative:* render in the API — rejected for exactly that reason.

**The cooldown is River uniqueness.** `UniqueOpts{ByArgs: true, ByPeriod:
time.Minute}` over `{Kind, UserID}` (locale excluded from uniqueness) makes a
second insert within the minute a no-op. Unknown emails never reach the queue,
so they need no cooldown. *Alternative:* a cache or table — none exists and
River already provides it.

**One sign-in helper.** Extract from `CreateSession` a
`signIn(ctx, user) (*api.UserHeaders, error)` that publishes `UserSignedIn`,
issues cookies and merges guest progress; `CreateSession`, `ConsumeMagicLink`
and `UpdatePassword` all return through it.

**Email lookup is normalised.** Trim and lower-case before the lookup, as
legacy did; the contract's `format: email` still rejects garbage first.

**Contract.** `updatePassword` returns `User | NotFoundError | ValidationError |
ApiError` (with the cookie header, like `createSession`); `@minLength(6)` on
`ResetPasswordInput.password` and `SignUpInput.password`. ogen validates length,
so the handlers need no length check.

**Frontend.** Four routes under `src/routes/{-$locale}/`: `magic_links/new`,
`magic_links/$token` (calls `consumeMagicLink` on load, redirects home on
success or to `new` with a message on 404), `remind_password/new`,
`password/$token/edit` (calls `checkPasswordResetToken` in the loader). All
through the generated Query hooks and Mantine forms, like `session/new`.

## Risks / Trade-offs

- [A Magic Link can be reused within fifteen minutes, so a leaked link is a
  session] → same as legacy; the lifetime is short and the link dies on email
  change.
- [Anti-enumeration costs the typo hint on the reminder form] → the form tells
  the visitor to check the address if no email arrives.
- [Sign-up rejects five-character passwords that the current contract accepts]
  → only affects new accounts; existing bcrypt hashes are untouched.
- [Postbox sandbox limits or DKIM not verified at cutover] → covered in the
  migration plan; the log adapter makes the flows testable without it.

## Migration Plan

1. Before cutover: verify the sender domain in Postbox, create a service
   account key, and set `EMAIL_TOKEN_SECRET`, Postbox credentials and
   `MAIL_FROM` in the cluster secrets (`legacy/k8s`).
2. Links emailed by Rails in the last fifteen minutes before cutover stop
   working; no action.
3. Rollback: the Rails app is untouched and still sends its own email.

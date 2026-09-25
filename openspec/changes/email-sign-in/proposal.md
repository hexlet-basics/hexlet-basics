# Proposal

## Why

The Go stack cannot send email at all, so the two email-based ways back into an
account — Magic Link and Password Reset — are unimplemented, and they block the
hard cutover (ADR-0002, `docs/PARITY.md` blocker #3). Five contract operations
(`createMagicLink`, `consumeMagicLink`, `createPasswordReminder`,
`checkPasswordResetToken`, `updatePassword`) wait on this one foundation.

## What Changes

- Transactional email: the API process enqueues, the worker process renders and
  delivers through Yandex Postbox (ADR-0006). Emails go out in the request's
  locale (`ru` by default, or `es`). Development logs each email instead of
  sending it, and tests capture it in memory.
- A second request for the same email within one minute sends nothing new.
- **Magic Link**: requesting one always answers the same way; an existing User
  gets a fifteen-minute link bound to their email; following it signs them in.
- **Password Reset**: requesting one always answers the same way — **BREAKING**
  against legacy, which told the visitor the email was unknown. The link lives
  fifteen minutes and dies when the password changes. Setting a new password
  now signs the User in (legacy left them signed out).
- Signing in through either link does what a password sign-in does: session
  cookie, guest progress merged, sign-in recorded.
- Passwords must be at least six characters wherever one is set (sign-up and
  reset); the contract currently accepts one.
- `updatePassword` returns the signed-in `User` instead of no content.
- Email links keep the legacy URLs (`/magic_links/:token`,
  `/password/:token/edit`), and the frontend gains the four pages behind them.
- Links emailed by the Rails app stop working at cutover; they live only
  fifteen minutes, so nobody is locked out for longer than that.

Out of scope: dropping the `en` locale (an open decision in `docs/PARITY.md`),
phone sign-in, passkeys, and post-parity email polish (text part, MJML).

## Capabilities

### New Capabilities

- `account-emails`: how account emails are addressed, localized, throttled and
  delivered.
- `magic-link`: passwordless sign-in through an emailed link.
- `password-reset`: replacing a forgotten password through an emailed link, and
  the password rule shared with sign-up.

### Modified Capabilities

None — `openspec/specs/` has no capabilities yet.

## Impact

- Contract: `api-spec/auth.tsp` (`updatePassword` response, `@minLength(6)` on
  `ResetPasswordInput` and `SignUpInput`); regenerated ogen server and hey-api
  client.
- Go: a new mailer package, a new River job, email-token signing, the five
  handlers in `internal/handlers`, new config (`EMAIL_TOKEN_SECRET`, Postbox
  credentials, sender address), DI wiring in both process graphs.
- Dependency: `aws-sdk-go-v2/service/sesv2` promoted to a direct dependency.
- Frontend: four new routes under `src/routes/{-$locale}/`, plus `ru`/`es`
  strings.
- Locales: email subjects and bodies in `internal/localization/locales`.
- Docs: `docs/PARITY.md` blocker #3 closed.

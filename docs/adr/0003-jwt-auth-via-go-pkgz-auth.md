# Authentication via go-pkgz/auth (JWT), adapting to its model

Authentication is built on **go-pkgz/auth** rather than assembled from
primitives. We adapt our design to its JWT-cookie token model to reuse the large
amount of ready-made logic it provides (OAuth flows, middleware, token
handling). The legacy app used opaque DB-backed sessions; we deliberately move
to JWT here.

## How each legacy sign-in method maps

- **Email + password** — go-pkgz/auth `direct` provider with a `CredCheckerFunc`
  that verifies existing bcrypt hashes via `go-crypt/crypt` (password compat).
- **Google / GitHub / Facebook** — built-in OAuth2 providers.
- **Passkey (WebAuthn)** — verified with `go-webauthn/webauthn` (added to the
  stack; it was missing from the reference lib set), then a go-pkgz/auth token is
  issued for the verified user.
- **Magic link, password reset, email confirmation** — short-lived signed tokens
  via `golang-jwt`.

## Consequences

- **No instant server-side revocation.** JWTs cannot be killed once issued.
  "Sign out everywhere" and admin bans are handled by short token TTL + refresh,
  or a custom `Validator` blocklist. This is the accepted price of reusing the
  library's logic.
- `go-webauthn/webauthn` is added to the dependency set for passkeys.

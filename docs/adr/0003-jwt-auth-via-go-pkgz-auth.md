# Authentication via go-pkgz/auth JWT tokens

Authentication uses **go-pkgz/auth's `token.Service`** rather than assembling
JWT signing, parsing, and cookie handling from primitives. The application owns
credential verification and adapts the token service's HTTP-oriented methods to
the generated ogen contract. The legacy app used opaque DB-backed sessions; we
deliberately move to JWT here.

## How each legacy sign-in method maps

- **Email + password** — the auth module loads the user through ent, verifies
  the existing bcrypt hash, and passes that same user to `token.Service` for JWT
  issuance. The library's `direct` provider is not used because its HTTP-only
  interface would add an internal HTTP round-trip at the ogen handler seam.
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

# Authentication via go-pkgz/auth JWT tokens

Authentication uses **go-pkgz/auth's `token.Service`** rather than assembling
JWT signing, parsing, and cookie handling from primitives. The application owns
credential verification and adapts the token service's HTTP-oriented methods to
the generated ogen contract. The legacy app used opaque DB-backed sessions; we
deliberately move to JWT here.

Cookie-authenticated requests use go-pkgz/auth's built-in double-submit XSRF
protection. `token.Service.Set` writes both the httpOnly `JWT` cookie and the
script-readable `XSRF-TOKEN` cookie. The generated Axios client uses its native
XSRF support to copy that value into `X-XSRF-TOKEN`; go-pkgz/auth's
`middleware.Authenticator` validates it before protected handlers run. Safe
`GET`, `HEAD`, and `OPTIONS` requests are exempt. `SameSite=Lax` remains
defense-in-depth rather than the primary CSRF control.

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
- XSRF issuance and validation stay inside go-pkgz/auth. Application code only
  adapts its multiple `Set-Cookie` headers to the generated contract and mounts
  the library's optional/required auth middleware at the router seam.
- `go-webauthn/webauthn` is added to the dependency set for passkeys.

# Authentication via go-pkgz/auth JWT tokens

**Status:** Accepted, amended by [ADR-0011](0011-contract-declared-authentication.md)

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

- **No instant cryptographic session revocation.** JWTs cannot be killed once
  issued. "Sign out everywhere" still needs a short token TTL, refresh policy,
  or a custom `Validator` blocklist. Administrative authorization is no longer
  accepted from JWT claims: ADR-0011 requires a current database check on every
  admin operation, so removing admin rights takes effect immediately.
- XSRF issuance and validation stay inside go-pkgz/auth. Application code only
  adapts its multiple `Set-Cookie` headers to the generated contract. The
  required policy is declared in TypeSpec and enforced through ogen's generated
  security seam; go-pkgz/auth's optional trace middleware remains responsible
  for sliding cookie refresh.
- `go-webauthn/webauthn` is added to the dependency set for passkeys.

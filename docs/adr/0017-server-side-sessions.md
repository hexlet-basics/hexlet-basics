---
status: accepted
---

# Server-side sessions instead of JWT cookies

Supersedes [ADR-0003](0003-jwt-auth-via-go-pkgz-auth.md) and amends
[ADR-0011](0011-contract-declared-authentication.md).

A signed-in visitor holds an opaque session token in an httpOnly cookie, and the
session row lives in PostgreSQL. We stop issuing JWTs for sessions and drop
`go-pkgz/auth`.

ADR-0003 chose a stateless JWT. The statelessness was never used: every request
already loads the user (`accounts.ActiveUserByID`) so that a removed account is
signed out at once, and the token carried nothing but a user id and the `jti`
that doubled as the XSRF token. We paid for a library whose HTTP-only API forced
`httptest.NewRecorder` into production code to read back its Set-Cookie headers,
and got no revocation in return. Nothing has been deployed, so there is no live
session to migrate.

**The session**

- [`alexedwards/scs`](https://github.com/alexedwards/scs) owns the cookie, the
  token rotation, the idle and absolute timeouts and the middleware. The token is
  stored hashed (`HashTokenInStore`), so a leaked table logs nobody in.
- The store is an ent entity (`user_sessions`; the Rails baseline's own `sessions`
  table is a different shape and is left alone), not scs's `postgresstore`. The store
  has to run inside the test's rolled-back transaction, which a `*sql.DB` store
  cannot, and expired rows are removed by a River periodic job in the worker
  rather than by a goroutine the store starts itself (ADR-0010).
- A session lives at most 31 days from sign-in and ends after 7 days without a
  request. The cookie is set once, at sign-in, for the full 31 days and never
  changes: pages are rendered by the SSR server, which calls this API for the
  browser and does not relay the API's `Set-Cookie`, so a cookie that moved with
  use would be refreshed for the SSR server and never for the browser. Idleness
  is the time of the last request, kept inside the session and recorded at most
  once an hour, so a read is not a write (scs's own idle timeout rewrites the
  session on every request and moves the cookie). The token is renewed on sign-in, sign-up and password change (session
  fixation) and destroyed on sign-out and account removal.
- A session remembers the user and the user's `password_changed_at` when it
  started. Changing or resetting the password makes every other session of that
  user stop matching, so it signs them out without scanning the table.

**One writer for the session cookie.** The scs middleware writes it. The
operations that used to return the auth cookies in `Set-Cookie` response headers
no longer do; `Set-Cookie` stays on the contract only where a cookie that is not
the session travels (the guest progress cookie). The cookie scheme in the
contract is named for the session cookie rather than `JWT`.

**CSRF.** The double-submit XSRF token (ADR-0003, ADR-0011) goes away with the
`jti` it was built on. Cross-origin writes are refused by the standard library's
`http.CrossOriginProtection`, which judges a request by the `Sec-Fetch-Site` and
`Origin` headers browsers set themselves, on top of the `SameSite=Lax` cookie.
Requests without either header (the SSR server, the GitHub webhook, tests) are
not browsers and pass. `XsrfToken` leaves the contract and with it every
`@useAuth([…, XsrfToken])`.

**What stays from ADR-0011.** Protected operations still declare their security
in TypeSpec and ogen's generated handler guards them; the `Identify` middleware
still supplies optional identity to public operations. Both now read the user
from the session instead of parsing a token. Administrator rights are still read
from the current database row, never from the session.

**Magic Link and Password Reset** are unchanged: they are short-lived signed
tokens (`emailtokens`) that end in a sign-in, not sessions.

## Consequences

- A session can be ended from the server, and a password change ends the others.
- A busy session is written back at most once an hour, not on every request. A
  session that went idle stays in the table until its 31 days are over, when the
  periodic job removes it; nothing reads it in between.
- `CrossOriginProtection` compares a request's `Origin` with its `Host` when the
  browser sent no `Sec-Fetch-Site`, so the ingress must pass `Host` through
  unchanged.
- `go-pkgz/auth` leaves `go.mod`; `golang-jwt` stays for `emailtokens`.

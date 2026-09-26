# Hard cutover from Rails to Go, with password and route compatibility

The Go backend replaces the Rails app in a single cutover once it reaches
parity — Rails and Go do NOT run side by side in production (no strangler
proxy). A broken state at HEAD is acceptable until parity is reached.

What "parity" covers — parity with production, not with the legacy codebase —
and when the English locale goes are pinned down in ADR-0015.

Two backward-compatibility constraints survive the cutover:

- **Passwords** — existing bcrypt password hashes must keep working, so users
  sign in with their current passwords without a reset (verified via
  `golang.org/x/crypto/bcrypt`).
- **URL routes** — existing public URLs must not break (bookmarks, SEO). This is
  why the API keeps legacy paths such as `/languages` even though the domain
  concept is `Course`.

## Consequences

- We migrate the existing Postgres data as-is; no session-cookie/CSRF interop
  with Rails is needed (sessions/CSRF are greenfield on the Go side).
- The URL surface mirrors legacy Rails routes rather than a clean REST redesign.

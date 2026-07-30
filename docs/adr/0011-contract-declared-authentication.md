# ADR-0011: Declare authentication in the HTTP contract

**Status:** Accepted

## Context

Protected route families were selected by handwritten URL-prefix checks in the
HTTP router. Authentication therefore lived outside the contract-first
pipeline: a newly added operation could be exposed by omission, and generated
clients could not see its cookie, XSRF, 401, or 403 requirements.

The JWT also contained an `admin` value. Because go-pkgz/auth can refresh a
valid token, accepting that claim for authorization could preserve privileges
after an administrator had been demoted in the database.

## Decision

TypeSpec is the source of truth for protected operations. It declares:

- `UserSession`, an API key in the httpOnly `JWT` cookie;
- `AdminSession`, the same cookie with current-administrator semantics;
- `XsrfToken`, an API key in the `X-XSRF-TOKEN` header.

Read-only protected operations require their session scheme. Unsafe operations
require the session and XSRF schemes as one AND security group. Protected
operations explicitly document `401` and, for administrator operations, `403`
Problem Details responses.

ogen's generated `SecurityHandler` runs before application handlers. Its
implementation verifies the JWT, loads the referenced user from PostgreSQL,
and stores that `*ent.User` in `context.Context`. Administrator authorization
checks the current nullable `users.admin` value from that loaded row and never
trusts an authorization claim in the JWT. New tokens therefore omit the admin
claim.

go-pkgz/auth's optional `Trace` middleware remains outside the generated server
to support the library's sliding cookie refresh and anonymous `/me` behavior;
it does not choose which routes require authentication.

The multipart attachment upload remains a temporary manual adapter because ogen
cannot generate its OpenAPI encoding. Only the exact
`POST /admin/attachments` route is wrapped, using the same JWT, database-admin,
XSRF, context, and Problem Details implementation as generated operations.

## Consequences

- Adding a protected operation requires an explicit TypeSpec security
  declaration, visible in OpenAPI and generated clients.
- Missing or invalid sessions and XSRF tokens are rejected before the business
  handler. Current database state can revoke admin access immediately without
  revoking the JWT itself.
- The router no longer contains policy-bearing path-prefix lists.
- Application handlers can reuse the authenticated database user from context
  instead of parsing claims or querying it again.
- Explicit 401/403 responses make generated operation result types unions; test
  helpers unwrap successful variants while contract tests exercise the error
  variants directly.

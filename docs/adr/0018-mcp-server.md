---
status: accepted
---

# A self-contained MCP server at `/mcp`

Amends [ADR-0015](0015-cutover-scope-and-api-prefix.md).

AI assistants should be able to read the catalogue, a learner's own progress and,
for staff, manage content, without a scraper or a copy of the admin UI. The
server exposes a remote MCP endpoint, `/mcp`, from `cmd/server`. It is
self-contained: its own tools and its own sign-in, not an extension of the
separate Hexlet platform MCP, because Hexlet Basics is a different database and
a different stack, and mixing the two tool sets would only confuse an agent.

**Transport and access.** Streamable HTTP in stateless mode through the
official `modelcontextprotocol/go-sdk` (v1.8.0 at the time of writing): one
request, one response, no sessions and no long-lived streams, so its load is
that of any other API call. The handler is mounted in `handlers.NewRouter`,
outside the `/api` prefix (ADR-0015 is amended to allow this; the ingress also
carries `/oauth/*` and `/.well-known/oauth-*`, ADR-0019). The whole endpoint
requires an OAuth bearer token: there is no anonymous tier, because MCP clients
only begin to sign in after a 401, so a half-open endpoint would never tell a
client it can. The public catalogue is already open on the site, in the sitemap
and in the feeds. The endpoint is behind a configuration flag and off by
default.

Token verification uses the SDK's bearer-token middleware and its protected
resource metadata handler. The verifier looks the token up (ADR-0019), checks
that its audience is the canonical `/mcp` URL, and refuses a token whose
subject is a client rather than a user.

**Tools.**

- *Catalogue*: courses, categories, lesson content, blog and books, read-only.
- *Learner*: the signed-in user's own courses and Enrollment, Lesson Progress
  and check results, read-only. There is no parameter naming whose data to
  read; the subject is always the token's. Not offered: a learner's solution
  code, assistant chat history, email and the rest of the profile.
- *Admin*: every admin operation, as a tool named after its operationId, which
  also keeps the names from colliding with the Hexlet platform MCP. They are
  registered from the contract through an allowlist of operationIds kept in
  code, so an operation added to the contract is not exposed to agents until it
  is listed. Each takes the locale as a required argument; it is never
  inferred. Operations that delete carry the `destructiveHint` annotation.
  Attachment upload takes base64 within the website's size limit.

**One checkpoint.** The permission check does not live in the handlers; it ran
in the generated security handler (ADR-0011). A tool that called a handler
directly would skip it. Every tool call therefore goes through the same
operation-layer check as an HTTP call. For an admin tool that is the
operation's Permission, the matching `mcp.*` Permission and the locale, read
from the database on each call (ADR-0017). Catalogue and learner tools are
checked against the token's Scope and its user subject only.

**Results and limits.** A tool returns structured content whose schema comes
from the contract, plus a short text. Defaults, all configurable: pages of at
most 50 records, a response of at most 100 KB, a 30-second timeout per call,
and a per-user and per-client rate limit that is stricter for writes and for
reads of personal data.

**Order of delivery**, each a pull request that can be checked and reverted:
the Permission model, then the audit trail, then the OAuth server and consent
screen, then catalogue and learner tools, then admin tools, then a check with
real Claude Code and Cursor.

## Considered Options

- **Generate every tool from the OpenAPI document.** Nothing to keep in sync,
  but open by default: each new operation becomes an agent capability the
  moment it is added.
- **Hand-written admin tools.** Friendlier tool design, but 59 operations that
  drift from the contract.
- **A separate MCP process.** Rejected: stateless calls cost what any API call
  costs, and ADR-0010 already pays for two runtimes.
- **Anonymous access to the public tools.** Rejected for the sign-in reason
  above; a separate unauthenticated endpoint can be added if a need appears.

## Consequences

- The agent is as powerful as the user: an administrator's token can reach
  every admin operation, including users and roles. That is deliberate, bounded
  by `mcp.manage`, the self-protection rules and the audit trail (ADR-0017).
- Tool schemas are part of the public surface: renaming an operationId renames
  a tool.
- Not verified yet: the SDK's exact stateless and middleware API against our
  router, and how real clients behave (ADR-0019).

## Addendum: the tracer slice (#849)

What the first implementation settled, and where it stops short of the text above.

- **SDK API, verified.** `modelcontextprotocol/go-sdk` v1.8.0: stateless mode is
  `mcp.NewStreamableHTTPHandler(..., &StreamableHTTPOptions{Stateless: true})`; we
  also set `JSONResponse` because nothing streams. `auth.RequireBearerToken` wraps
  it with a `TokenVerifier`, and `auth.ProtectedResourceMetadataHandler` serves
  the metadata. Both work in `handlers.NewRouter` unchanged. Stateless mode
  answers GET and DELETE with 405.
- **Switch.** `MCP_ENABLED` (default off); enabling it without `OAUTH_ENABLED` is a
  configuration error. The package is `internal/mcpserver`; it does not import
  `internal/handlers`, which imports it to mount `/mcp`.
- **Verification.** The verifier takes `oauthserver.Server.VerifyAccessToken` and
  additionally requires `Token.Audience` to contain `Resource()` and the subject
  to be a user; anything else is a 401 with a `WWW-Authenticate` challenge whose
  `resource_metadata` is the path-inserted URL. Metadata is served at both
  `/.well-known/oauth-protected-resource` and `.../mcp`. A failure of our own is a
  500 whose detail goes to the log, not the body. The library requires a token
  expiry that `Token` does not carry, so `AllowMissingExpiration` is set: the
  lookup already refuses an expired token.
- **Scope check per tool.** The endpoint-wide middleware option would answer 403
  to a connection lacking one Scope for every tool, so each catalogue tool checks
  `catalog:read` itself and answers a tool error.
- **Output schema from the contract.** `internal/mcpserver/contract.go` reads the
  operation's 200 schema from the embedded OpenAPI document, inlines `$ref` and
  rewrites `nullable`. Structured content is an object, so the contract's array
  sits under one key (`{"courses": [...]}`); a test validates real output against
  the advertised schema.
- **Tool name.** `listCourses`, the operationId, as admin tools will be.
- **Not done in this slice.** Paging and the 100 KB / 30 s / rate limits, the
  locale as a tool argument (the request's `Accept-Language` decides, as for the
  HTTP operation; a test pins that it reaches the tool), and the `mcp.*` Permission check, which belong to the tickets
  for the other tools.
- **Ingress not yet routed.** The Helm ingress (`k8s/app-chart/templates/ingress.yaml`)
  still sends only `/api/` to Go. `/mcp`, `/.well-known/oauth-protected-resource*`
  (and `/oauth/*` from ADR-0019) must be added before either flag is switched on.
  The Vite development proxy already carries them.

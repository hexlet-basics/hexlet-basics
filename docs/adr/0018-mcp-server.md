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

## Addendum: catalogue tools, response limits and rate limiting (#851)

- **Catalogue tools.** Seven read-only tools, named after the operationId they
  mirror, each needing `catalog:read` and calling the same handler method as the
  HTTP operation (`mcpserver.Catalogue` is satisfied by `handlers.Server`):
  `listCourses`, `listPublicCourseCategories`, `getPublicCourseCategory`,
  `getCourse` (a Course with its Lessons; its slug is a landing page's, as
  `listCourses` returns it), `getCourseLesson` (Lesson content with its Lesson
  Tests; `courseSlug` is the Course's own), `listBlogPosts` and `getBlogPost`.
  The output schema is the operation's 200 schema from the contract plus one
  optional `truncation` property; the two operations that return a bare array
  (`listCourses`, `listPublicCourseCategories`) are wrapped under `courses` /
  `categories` with `total`, `page` and `perPage`.
- **Visitor rules.** A bearer token is not a session, so the handlers answer as
  they do a signed-out visitor: no Enrollment, progress at the beginning,
  drafts and archived posts not found. A Lesson's `preparedCode` (Starter Code),
  `originalCode` (Reference Solution) and `testCode` are in the answer because
  the HTTP operation gives them to every visitor; the tools add no rule of their
  own. A test compares a tool's result with the HTTP response for a visitor.
- **Not exposed.** `getBook` (the book page's state is the caller's own book
  request, so for a visitor it is always empty, and the PDF itself is behind a
  sign-in and a request); `getNextBlogPost` (infinite-scroll plumbing, covered
  by `listBlogPosts`); the Yandex feeds and the sitemap (machine formats for
  other consumers); everything that writes (`likeBlogPost`, `createLead`...).
  There is no public read of a book's content in the contract, so "books" has
  no tool until there is.
- **Page size: clamp.** A `perPage` above `MCP_MAX_PAGE_SIZE` (50) is clamped,
  not refused, and the text of the result says so; the schema states the limit
  in the argument's description, not as a `maximum`, because the limit is
  configurable. Lists the contract does not page are paged in the tool layer
  (`page`, `perPage`); `listBlogPosts` passes the clamped size to the handler,
  which pages in SQL. Lists nested in a record (a Course's Lessons, a
  category's courses) are not paged and are bounded by the response limit.
- **Response size: cut at an element, with a hint.** The structured content plus
  the text is measured against `MCP_MAX_RESPONSE_BYTES` (100 KB). Over it, the
  largest array (then the largest long string) is shortened, repeatedly, on the
  decoded value, so the result is always valid JSON, and `truncation` says what
  was cut. For a list it carries `suggestedPerPage` and `suggestedPage` that
  continue the reading (the first call may repeat a few records already seen).
  A result that cannot be shortened to the limit is a tool error naming the
  limit.
- **Timeout.** `MCP_TOOL_TIMEOUT` (30 s) becomes the deadline of the context the
  handler runs with, which is what cancels its queries; the agent gets a tool
  error.
- **Rate limit.** Per user and per Connected App (the token's OAuth client id,
  already in `Token.ClientID`; the `Token` is unchanged), per Class of tool:
  `catalogue`, `personal` (reads of the user's own data) and `write`, each with
  its own budget per `MCP_RATE_WINDOW`. Defaults per window, user / app:
  catalogue 60 / 600, personal 30 / 300, write 10 / 100. A call spends one of
  both budgets, or neither. A tool declares its Class when it is registered
  (`addTool`), so #852 and #854 only pick one; only the catalogue class is
  exercised by a tool today, the others are unit tested. The Scope check comes
  first, so a call refused for its Scope spends nothing.
- **Over the limit: a tool error, not HTTP 429.** The Class is known only from
  the tool name inside the JSON-RPC body, and a 429 on the endpoint would also
  stop `initialize` and `tools/list`, which cost nothing, and an SDK client
  takes it for a broken transport rather than a message to the agent. The error
  text and `_meta.retryAfterSeconds` name the wait.
- **One limiter.** The counter of `oauthserver`'s registration limit moved to
  `internal/ratelimit` (a fixed window per string key, `Allow` over several
  charges atomically) and the OAuth server keeps only the address bucketing. Its
  tests are kept, minus the cap test, which moved with the counter. It is in
  memory, per replica: with several replicas a caller gets its budget once per
  replica, a brake on a runaway agent and not an exact quota.
- **Configuration.** `MCP_MAX_PAGE_SIZE`, `MCP_MAX_RESPONSE_BYTES`,
  `MCP_TOOL_TIMEOUT`, `MCP_RATE_WINDOW` and `MCP_RATE_{USER,CLIENT}_{CATALOGUE,
  PERSONAL,WRITE}`, all with defaults; a value below its minimum is a
  configuration error when `MCP_ENABLED` is on.

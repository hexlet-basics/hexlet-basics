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

## Addendum: admin tools and the one checkpoint (#854)

Settled by the first admin-tool slice (banners: `adminListBanners`,
`adminUpdateBanner`), which every later admin ticket inherits. It meets ADR-0017
and ADR-0018 as written: one checkpoint, identical to the HTTP check. One schema
change came with it (the `token_id` column type, below).

**How a tool call reaches an operation.** It does not. A tool call becomes an
HTTP request synthesized from the arguments and is served, in process, by the
generated `api.Server` that serves the website (`handlers.NewAdminDispatcher`:
that server behind the locale middleware and `WithClientIP`, without `Identify`
and `Trace`, which read the session cookie). So the security handler, body
decoding and validation, the audit middleware (transaction, Audit Record,
`failed`) and the central error mapping are the code an HTTP call runs, not a
copy of it. `internal/mcpserver` is given only that `http.Handler`
(`Deps.API`); it holds no `api.Handler`, no ent client and does not import
`internal/handlers`. A test (`TestThePackageCannotReachAnOperationExceptThroughTheGeneratedServer`)
fails if it starts to, and `TestNoRegisteredAdminToolRunsWithoutTheCheck` calls
every registered admin tool as a caller whose role holds nothing, with arguments
generated from the tool's own schema, and requires a refusal, no `succeeded`
Audit Record and no changed row.

**The principal.** The tool puts an `mcpserver.Principal` (user, the token's
Scopes, the locale argument, a per-call nonce) in the request context.
`AuthHandler.HandleAdminSession`, the one gate of every admin operation, takes
the MCP branch (`authorizeMCP`) when the context holds one. A context value
cannot arrive over HTTP, so no header or cookie can claim to be a Principal
(`TestBearerTokenIsNotAcceptedOnAnyAPIOperation` and
`TestSessionPlaceholderOpensNothingOverHTTP` pin it). The branch loads the user
from the database and judges, on every call and in this order: the `admin`
Scope; `Policy.Allow(user, locale, operation)`, the call the cookie branch makes;
`Policy.HasCapability` for `permissions.MCPCapability(operation)`, which is
`mcp.write` for an admin write and `mcp.read` for the rest, from the contract's
declared kind (never the administrator flag, ADR-0017). A refusal is the HTTP
refusal: an Audit Record `forbidden` (channel `mcp`) and a 403, before the body
is decoded. Nothing is cached, so removing a grant applies on the next call.
The tool does not check the Scope or the capability itself: a second place that
refuses would refuse without recording.

**Two things the generated code forces, and how they are kept harmless.** The
generated server will not call a security handler without a `JWT` cookie, and a
write also needs the `X-XSRF-TOKEN` header. The synthesized request therefore
carries a placeholder cookie (`mcp-in-process`) and the Principal's nonce as the
XSRF header, and `authorizeMCP` stores the nonce as the session id, so the
second scheme passes by the same comparison the website's requests pass. Outside
a tool call the placeholder is a JWT that does not parse (401). This is the
single mechanism that departs from "no cookie on MCP"; it opens nothing.

**Locale.** `locale` is a required argument of every admin tool, validated
against `en` and `ru` (the locales staff work in; the website's matcher would
otherwise fall back to English for a language it does not know, judging a call
in a locale nobody named). It is sent as the request's `Accept-Language`, which
is what the HTTP check reads, and the Principal carries it too: `authorizeMCP`
refuses when the locale the generated server resolved differs. The connection's
own `Accept-Language` is not read (a test sends `en` there and `ru` as the
argument). The request body is nested under `body` in the input, not merged with
the parameters, because a body field can collide with a call argument (a banner
has a `locale` of its own).

**The registry.** `adminOperations` in `internal/mcpserver/admin.go` is the
allowlist of operationIds, in code. Building the server fails for a name that is
not in the contract or does not require the `AdminSession` scheme (which would
skip the checkpoint). Tool name is the operationId; the input schema is the
operation's path and query parameters (references to `components/parameters`
resolved) plus `locale` and `body`; the output schema is the first 2xx JSON
response, wrapped under `result` only when it is not an object. A write
carries `destructiveHint: false`, a delete `true`; GET is read-only. A non-2xx
answer is a tool error whose text is `HTTP <status>: <problem document>`, byte
for byte what an HTTP caller receives (`TestHTTPAndMCPBehaveTheSameForTheSameOperation`).
Adding a tool is adding its operationId; #855 lists the rest. Operations that
declare `x-unit-of-work: standalone` (the upload) need nothing more: the
middleware already handles them.

**Audit Record, Connected App and token.** The dispatcher sets
`audit.WithSource` before dispatch, so `succeeded`, `forbidden` and `failed`
records alike carry channel `mcp`. The Connected App is the id of the
`oauth_clients` row behind the token's `client_id` (`oauthserver.Token.ConnectedAppID`,
filled by `VerifyAccessToken`, which also refuses a token whose client row is
gone). The token is the access token's own id (`oauth_access_tokens.id`, not the
bearer, which is that id encrypted). That id is a random string and the column
was reserved as `bigint`, so migration `audit_records_token_id_text` changes
`audit_records.token_id` to text and `audit.Source.TokenID` to `*string`
(nothing had ever written it). The request's user agent, `X-Forwarded-For`,
`X-Real-IP` and `X-Request-ID` are copied from the `/mcp` request to the
synthesized one, so a record keeps the client's agent and address as an HTTP
record does.

**Test seam.** `testsupport.NewOAuth` now builds the generated server with the
production audit middleware and hands `NewAdminDispatcher` of it to the MCP
server, and gains `Grant`, `RevokeGrant` and `Records` (documented in
`internal/testsupport/mcp.go`). A tool test is: grant the role what the case
needs, `Connect` with the Scope, call the tool, assert on the result, the rows
and `Records`.

**Not done here, for the next tickets.** `mcp.manage` (ADR-0017: operations on
users, staff and roles) is not applied by `MCPCapability`; the operations that
need it must not be added to the allowlist until it is (#855 decides how it is
declared in the contract, and the rule must go in `MCPCapability`, the one
place). Upload (#856) takes its bytes as base64, which the synthesizer does not
build yet (the body here is JSON). The page-size cap of #851 is not applied to
admin lists, which keep the contract's own `perPage` maximum of 100.

**Limits (after #851).** Admin tools share the limiter, timeout and response size
limit of `addTool`, but are registered through their own handler, because their
Scope is judged by the security handler (so a refusal is recorded) and not by
`admit`. `Server.spend` is the rate-limit half of `admit`, used by both. An admin
write is `ClassWrite`; every admin read is `ClassPersonal` (stricter than the
public catalogue: it can show staff-only and personal data), whatever its
operation, so no list of "personal" operations is kept in code. A call refused for
the limit never reaches the checkpoint and leaves no Audit Record.

**One `WithMiddleware` call.** ogen's `api.WithMiddleware` replaces the previous
option, so the audit middleware and the personal-data read middleware are passed
in a single call (`di`, `testsupport.NewHarness` and `testsupport.NewOAuth`). A
dispatcher is never given a server of its own, so a call by tool shows the effect
of both (`TestDispatchedCallsPassBothMiddlewares`: a write leaves its Audit
Record, a read of users publishes its read fact, both with channel `mcp`, the
Connected App and the token). The read fact's `TokenID` is a string for the same
reason as the column.

## Addendum: the learner tools (#852)

- **Tools.** `getMyDashboard` is the operationId of the HTTP operation and returns
  its contract schema (the handler's code, called with the token's user).
  `getMyCourseProgress` (argument `courseSlug`) and `getMyLessonProgress`
  (`courseSlug`, `lessonSlug`) are named wrappers: the contract has no single
  operation that tells a Lesson that is Started from one that is Not started
  (`LessonProgressItem` carries only `finished`), so their output schemas are written
  in `internal/mcpserver/learner.go`, not read from the contract, and the data
  comes from `progress.Tracker.LearnerCourse`, which builds on `CourseStates` so
  Completion and the Next Lesson are the numbers the website shows. All need
  `progress:read` and are read-only.
- **The subject is the token's.** No argument names a user, and arguments are
  decoded strictly: one the schema does not declare (a `userId`) is a tool error,
  not ignored. A course or lesson slug may name anything; the answer is still the
  token's user's rows (Not started for what they never opened).
- **Never offered.** Solution code, assistant chat, email and the rest of the
  profile appear in no schema or result; a test walks every schema and result for
  those names.
- **"Result of the last check" is a deviation from the ticket's wording.** The
  platform stores neither the submitted code nor the check output (a failed
  attempt only leaves the Lesson started), so there is no last-check record to
  expose. `getMyLessonProgress` reports `passedCheck`: whether a check of this
  Lesson has passed (Finished). Exposing output would need a new table and, since
  output can echo the submission, a decision about solution code; it is left out.
- **Rate limiting.** The tools register through `addTool` in the `personal` Class
  of #851 (30 calls per user per window by default), behind the Scope check.
- **Guest progress** does not exist over MCP: a bearer token always has a user.

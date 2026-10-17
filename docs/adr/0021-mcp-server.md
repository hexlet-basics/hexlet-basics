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
carries `/oauth/*` and `/.well-known/oauth-*`, ADR-0022). The whole endpoint
requires an OAuth bearer token: there is no anonymous tier, because MCP clients
only begin to sign in after a 401, so a half-open endpoint would never tell a
client it can. The public catalogue is already open on the site, in the sitemap
and in the feeds. The endpoint is behind a configuration flag and off by
default.

Token verification uses the SDK's bearer-token middleware and its protected
resource metadata handler. The verifier looks the token up (ADR-0022), checks
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
from the database on each call (ADR-0020). Catalogue and learner tools are
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
  by `mcp.manage`, the self-protection rules and the audit trail (ADR-0020).
- Tool schemas are part of the public surface: renaming an operationId renames
  a tool.
- Not verified yet: the SDK's exact stateless and middleware API against our
  router, and how real clients behave (ADR-0022).

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
  (and `/oauth/*` from ADR-0022) must be added before either flag is switched on.
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
`adminUpdateBanner`), which every later admin ticket inherits. It meets ADR-0020
and ADR-0021 as written: one checkpoint, identical to the HTTP check. One schema
change came with it (the `token_id` column type, below).

**How a tool call reaches an operation.** It does not. A tool call becomes an
HTTP request synthesized from the arguments and is served, in process, by the
generated `api.Server` that serves the website (`handlers.NewAdminDispatcher`:
that server behind the locale middleware and `WithClientIP`, without the session
pipeline, `Security.Pipeline`, which loads a session from the cookie). So the security handler, body
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
Scopes, the locale argument) in the request context.
`Security.HandleAdminSession`, the one gate of every admin operation, takes
the MCP branch (`authorizeMCP`) when the context holds one, before it looks at
the session. A context value cannot arrive over HTTP, so no header or cookie can
claim to be a Principal
(`TestBearerTokenIsNotAcceptedOnAnyAPIOperation` and
`TestSessionPlaceholderOpensNothingOverHTTP` pin it). The branch loads the user
from the database and judges, on every call and in this order: the `admin`
Scope; `Policy.Allow(user, locale, operation)`, the call the session branch makes;
`Policy.HasCapability` for `permissions.MCPCapability(operation)`, which is
`mcp.write` for an admin write and `mcp.read` for the rest, from the contract's
declared kind (never the administrator flag, ADR-0020). A refusal is the HTTP
refusal: an Audit Record `forbidden` (channel `mcp`) and a 403, before the body
is decoded. Nothing is cached, so removing a grant applies on the next call.
The tool does not check the Scope or the capability itself: a second place that
refuses would refuse without recording.

**What the generated code forces, and how it is kept harmless.** The generated
server will not call a security handler without the `session` cookie. The
synthesized request therefore carries a placeholder cookie (`mcp-in-process`).
The branch never reads it: it is taken before the session is consulted, and the
user comes from the Principal and the database. Sent over HTTP, where the context
holds no Principal, the same value names no session, so the session lookup finds
nothing and answers 401. This is the single mechanism that departs from "no
cookie on MCP"; it opens nothing. (Rebased onto ADR-0017's server-side sessions:
the design first needed a second placeholder, a per-call nonce sent as the XSRF
header the contract then asked of a write and stored as the JWT id so the
comparison passed. The XSRF token left the contract with ADR-0017, which
protects a browser's cross-site write by `Sec-Fetch-Site` and `Origin`, headers
an in-process request does not carry, so the nonce is gone and the Principal
holds only the user, the Scopes and the locale. Nothing else changed: the order
and the content of the checks are those above.)

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

**Not done here, for the next tickets.** `mcp.manage` (ADR-0020: operations on
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

## Addendum: the attachment upload tool (#856)

`adminUploadAttachment` is the first tool whose operation takes a multipart
request. It is not a second upload path: the tool call is still an HTTP request
served by the generated server (#854), so the security handler (`any-staff`,
the `admin` Scope, `mcp.write`), validation, `assetstore.Put` (size, raster
allowlist, sniffed type), the `x-unit-of-work: standalone` handling and the Audit
Record (written after the call in its own transaction, channel `mcp`, Connected
App and token) are the website's.

**Input.** The registry reads the operation's request body from the contract: a
`multipart/form-data` body of exactly one binary part is a "file operation"
(anything else multipart is refused when the server is built). Its `body`
argument is `{fileName, contentBase64, contentType?}`; `buildRequest` builds the
multipart body with the part named as the contract names it. `contentType` is
advisory, as on the website, where the stored type comes from the bytes. The
output is the contract's `Attachment` (id, url, filename, contentType, byteSize).
An upload is `ClassWrite` (it is an admin write), like every other write.

**Limits.** The website's limit, `assetstore.MaxUploadBytes` (15 MiB), is the
tool's: the length of the base64 text is checked before it is decoded, so an
oversized argument costs no decoded copy, and the tool answers "too large" and
stores nothing. Invalid base64, a missing name and unknown fields are tool
errors. The multipart envelope cap of the website's route (`UploadBodyLimit`) is
not needed in process: the tool builds the envelope itself, from at most
15 MiB. Departure to note: the SDK's default `/mcp` request limit of 4 MiB would
refuse the base64 of a 15 MiB file (20 MiB), so the limit of the `/mcp` endpoint
is raised to the base64 length of the largest upload plus 1 MiB. It applies to
every tool, but only to callers with a valid token (the bearer check runs before
the body is read), and a request still has to fit the rate limit.

## Addendum: every admin operation as a tool, and `mcp.manage` (#855)

Meets ADR-0020 and ADR-0021 as written; no departure. It closes the "not done
here" of #854.

**The allowlist is the whole admin surface.** `adminOperations` lists every
operation that requires the `AdminSession` scheme (86 tools), minus
`excludedAdminOperations`, a map from operationId to the reason it is not a tool.
A test (`TestEveryAdminOperationIsAToolOrExcludedWithAReason`) enumerates the
contract's admin operations and requires each to be in exactly one of the two, so
an operation added later fails the build until someone decides. Excluded:
only `adminGetMyCapabilities` (it feeds the admin menu; an agent learns what it
may do from the answers of the tools, which name the refusal). The audit trail
(`adminListAuditRecords`, #857) is a tool. The state transitions (`publish`, `unpublish`, `archive`,
`adminMarkCourse...`) are tools of their own, as they are operations of their
own; the `update` tools carry no state.

**`mcp.manage` is declared in the contract.** A TypeSpec extension,
`x-requires-capability: "mcp.manage"`, on the operation; `permissions.Parse`
reads it (and rejects a name not listed under `x-capabilities`), and
`permissions.MCPCapabilities(operation)` returns every capability the operation
needs over MCP: `mcp.write` or `mcp.read` by its kind, plus the declared one.
It replaces `MCPCapability`, and `authorizeMCP` checks each in turn through
`Policy.HasCapability`, so an administrator without an explicit grant is refused
like anyone else. There is no list in Go.

**What manages who has access.** Every operation whose Permission is in
`users.*`, `staff_members.*` or `staff_roles.*`, 22 operations, reads included:
the user CRUD and search, the management users, the Staff Members, the Staff
Roles, a role's grants and the Permission catalog. Writes are the obvious case
(they create administrators, move people between roles and rewrite grants). Reads
are included because the user and staff reads are personal data (ADR-0020's
reads of personal data), and the roles and the catalog are the map an agent needs
to plan a privilege escalation; keeping the reads under `mcp.read` would let a
role meant for content work list every learner. So `mcp.manage` covers the
access-management operations whatever their kind, on top of the operation's
Permission and `mcp.read`/`mcp.write`. Not included: the leads, messages and
lesson progress (personal data, but about learning and sales, not access), the
learner-facing reads, and the content operations.

**Self-protection holds, by construction.** The #844 rules are in the handlers,
and the tool dispatches through the generated server with the Principal's user as
the authenticated user, so the refusals are the handlers': own Staff Member
record (403), own account (403), the last administrator (409), the
administrator flag by a non-administrator (403) and the grants of the caller's
own role by a non-administrator (403). Each is tested through a tool, with the
Audit Record `failed`, the channel `mcp` and the row unchanged.

**Hints.** `readOnlyHint` for a GET. `destructiveHint` only for an HTTP DELETE:
archive and unpublish keep the record and can be undone, and a `PUT` that
replaces a set (a role's grants, a post's related courses) is reversible by
another call, so it is `destructiveHint: false`, which says that the data is not
lost. `idempotentHint` for a PUT and a DELETE. A description is the contract's
text, or, for the operations that have none, the operationId read as words,
followed by what the call needs (its Permission, and `mcp.manage`).

**Schemas.** The derivation of #854 (path and query parameters, a nested `body`,
`$ref`s inlined, `nullable` as a type list or `anyOf` with null) held for every
operation; no change was needed. A test builds every tool and validates an
example generated from its input schema against it, and checks the output schema
is valid. The admin tool list is about 225 KB of JSON (every input and output
schema is inlined), and a test fails if it passes 300 KB; trimming it (for
instance by advertising shared output shapes once) is left for a ticket of its
own. The test harness (`testsupport.NewOAuth`) now wires the roles writer, the
blog post remover and the related-courses replacer, which the new tools reach.

## Addendum: a cheap, discoverable tool list (#878)

Loading and searching tools lazily is a client feature (Claude Code and the
Claude API defer tools; the MCP proposals for server-side search of `tools/list`
were closed unmerged), so the server cooperates instead of inventing a
mechanism: it says early what it has, keeps the list small, stable and
per-caller, and offers one lookup for a client that cannot defer. It does not
offer a way to run a tool by name.

**Measured.** The admin tools (87) were 228,757 bytes of JSON (about 57K
tokens): 166 KB of output schemas, 38 KB of input schemas, 8 KB of descriptions.
They are now 79,086 bytes (about 20K tokens), 65% less: 25 KB, 33 KB and 5 KB.
The whole list for a caller who sees everything (98 tools: admin, catalogue,
learner, `findTools`) is 88 KB, and a test fails past 100 KB; the admin tools
alone are held to 90 KB. The budget is a decision to make on purpose.

**What was cut.**
- *Output schemas* are advertised as a shape (`shape.go`): the result's own
  fields by type, one level (`items: array`, `id: integer`), without prose,
  enums, formats or `required`. A client reads values and does not validate
  them against ours, and every reduction accepts more than the contract's
  schema, so a real result still validates (the existing test that checks
  results against the advertised schema holds). The root is always a plain
  object: a root `anyOf`/`oneOf` is dropped, because Claude Code flattens it and
  loses the enforcement. The `truncation` property is its name and type only; what
  the hint contains is told once, in the instructions.
- *Input schemas* keep their constraints (an agent that sees `minLength` or an
  enum writes a valid call first time) and lose titles, examples, vendor
  extensions, number widths and the long part of every description (a property
  keeps its first sentence, at most 100 characters). The generated server still
  validates against the full contract, so nothing is loosened where it counts.
- *Descriptions* are the contract's first sentence (at most 140 characters)
  plus what the call needs; the rest of the contract's text said how errors
  surface, which a refusal says when it matters. The hand-written catalogue and
  learner descriptions were shortened the same way. A tool's `title` is gone
  where it only repeated the name.
- *Hints*: every tool sets all four annotations (`hints` in `tool.go`), since the
  protocol's defaults (destructive, open-world) are the unsafe ones. Nothing here
  reaches beyond the site's own data, so `openWorldHint` is false everywhere, and
  reads are `destructiveHint: false`, `idempotentHint: true`.

**Instructions.** About 1.9 KB (the limit is 4 KB, where Claude Code cuts): the
domains (catalogue, learner progress, content administration, access management,
audit), what each needs, when to call `findTools`, the rule that a state change is
its own tool, and the shape of a cut result.

**Names are not renamed.** The ticket asks for resource prefixes so that one
search matches a group. A tool's name is its operationId (this ADR, above) and the
names are a public surface that clients cache, so they stay: every admin name
starts with `admin` and carries its resource (`adminListBanners`,
`adminPublishBanner`), so a search for "banner" matches the group, which a test
checks. Order is deterministic: by name, the same on every call.

**The list is per caller; the call is judged all the same.** `tools/list` goes
through a receiving middleware that shows a connection a tool only when its token
holds the tool's Scope and, for an admin tool, the user's Staff Role holds every
capability Permission the tool needs (`mcp.read` or `mcp.write` by kind,
`mcp.manage`, `audit.read`). The capabilities come from `Deps.Capabilities`
(`Security.HeldCapabilities`, which asks `Policy.HasCapability`, the check
`authorizeMCP` makes), read on each listing; nothing is cached, so a grant taken
away shows at the next list. The operation's own Permission and locale are not
used to hide a tool (they depend on the call). If the capabilities cannot be
read the list is narrowed by Scope only: showing a tool costs one refusal and
hiding is not authorization. A test calls a hidden tool by name (no admin Scope,
no `mcp.write`) and gets the same 403 as before, with no succeeded Audit Record;
`TestNoRegisteredAdminToolRunsWithoutTheCheck` now takes its tools from the
contract, not from the filtered list, so it still calls every one. List results
are `cacheScope: private` (`SetCacheable`), since they differ by caller; the TTL
is left at zero so a change of grants is seen.

**The discovery tool: `findTools`.** Measured, the list is still about 20K tokens
for a client that cannot defer tools, so the one lookup the ticket allows was added.
It is read-only and looks at a registry of the tools and what each needs; it holds
no handler of any other tool, so there is nothing in it that could run one. Arguments:
`query` (words that must all match a name or description), `detail` (`name`,
`description`, the default, with the Permission each needs, or `schema` with the
full input schema) and `limit` (20, at most 50). It shows only what the caller's
`tools/list` shows, needs no Scope and spends the catalogue rate budget. There is
no `call_tool(name, args)` and no `execute(code)`: either would bypass or flatten
the per-operation checkpoint and its Audit Record. A test pins that no tool is named
or takes an argument like one, and that `findTools` asked to run `adminCreateBanner`
refuses the argument and leaves no row and no Audit Record. The name is `findTools`
and not `search`, which the content search of #877 takes.

## Addendum: catalogue search, and the `search` and `fetch` tools (#877)

Meets the ADR as written; no departure. (The ticket calls this record "ADR-0018",
its number when the ticket was written.)

**Why these two tools and not a tool-finding tool.** Searching and loading tools
lazily is a client feature (Claude Code and the Claude API defer tools, and the
MCP proposals for a server-side search of `tools/list` were closed unmerged), so
the server cooperates with it and invents nothing: scope-filtered `tools/list`
stays as it is, and a generic `call_tool(name, args)` or `execute(code)` is not
built, since it would bypass the per-operation checkpoint and audit (ADR-0020).
For content, `search` and `fetch` is the de facto convention of ChatGPT
connectors and deep research, and it is what a client built to it looks for by
name.

**The operation.** `searchCatalogue` (`GET /search`) is a public contract
operation: `q` (1 to 200 characters), an optional `type` (`course`, `lesson`,
`post`) and the usual `page` and `perPage`, in the request locale. It answers a
`Page` of `SearchResult`: `type`, a typed `id`, `title`, a plain-text `snippet`
and the page's `path` with the locale prefix and no origin. The HTTP handler and
the MCP tool run the same method. The business rules are in `internal/search`
(`SearchCatalogue` command, `Hit` result), per the use-case module rules.

**What is searched, and what a visitor may see.** Only text the authors wrote:
a landing page's name, header, description and outcomes; a lesson's name,
description, theory and instructions; a post's name, description and body. The
records are those the public reads already scope, by the same predicates
(`landingpages.Listed`, `blogposts.Published`, the course's current Version): a
listed landing page of a completed Course, a lesson of that Course's current
build translated into the locale, a published post of the locale. A draft or
unlisted page, a retired build, an unpublished post and another locale are not
found, each tested by making the record readable and seeing it appear.
A lesson is found only while its course has a listed landing page in the
locale, so a lesson never outlives its course on the list.

**No learner-written text, by construction.** The module reads three tables
(landing pages, lesson translations, blog posts) and the two it needs for slugs.
It never touches reviews, assistant messages, leads, users or lesson reviews; a
test records the statements a search sends and fails if one names such a table,
and another searches for the words of a fixture review and finds nothing. The
**prompt-injection note** the research asked for: result text comes from authors
and not from visitors, which is why it can be shown to a model as the site's own
copy; but it is still data, not instruction, and a client is expected to treat it
so. Keeping every user-generated field out of `search` results (and out of
`fetch`: a post's author name is a Staff Member's public byline) is therefore a
rule, and adding a field that carries visitor text to either needs this note
revisited.

**PostgreSQL full-text search, no new service.** One `IMMUTABLE` SQL function,
`catalogue_search_vector(config, title, summary, body)`, builds the document
(title weight A, summary B, body C), and a migration creates a GIN index on it
for each table and for `russian` and `english`. The query writes the
configuration as a literal chosen from a fixed map (`ru` and `en`; any other
locale falls back to `simple` and an unindexed scan, which is fine while the
catalogue has no content in it), so the planner can match the index expression,
and reads the text with `websearch_to_tsquery`: quotes, `or` and `-word` work and
no input is a syntax error. The word `Hello` in a Russian post and in an English
one are the same query under two configurations, which is what the locale scope
is for. An `EXPLAIN` test records the statements a search sends and shows each
language's index of each table in the plan. The fixtures are a handful of rows,
for which PostgreSQL rightly reads the table, so the test turns sequential scans
off for its transaction, and for the lessons also drops the locale index that
would otherwise be the cheaper way in (in that transaction only): what it proves
is that the query's expression is the one the index holds, not that the planner
picks it at that size. Results are ranked by `ts_rank` (best first, then kind,
then id, so a page is stable). Each kind is queried for its best
`page * perPage` rows and the lists are merged; the window is limited to the
first 1000 hits, because a search is for finding a page and the lists are for
paging, and an unbounded page number would be an unbounded query. Snippets
(`ts_headline`) are made in a second query, only for the hits of the window,
since producing one reads the whole text of the record.

**Typed ids.** `course:<landing page slug>` (the slug `getCourse` takes, so per
locale), `lesson:<course slug>/<lesson slug>` (the course's own slug, as
`getCourseLesson` and the lesson page's URL have it) and `post:<slug>`. They are
made of what the page's URL carries, so they are stable as long as the page is,
cannot be guessed to reach anything a URL does not, and hold no database id.
Because a landing page slug depends on the locale, an id of a course or lesson
is read in the locale of the connection it came from; a `fetch` of a course in
another locale is "no record has the id", and the error says so.

**The tools.** `search(query, type?, page?, perPage?)` answers `{results: [{id,
title, url, type, snippet}], total, page, perPage}`; `fetch(id)` answers `{id,
title, text, url, metadata}`. Both return that JSON as `structuredContent` and,
identically, as one text block (a flag on the tool result makes the text the JSON
of the data, written after any cut, so the two cannot differ). Both need
`catalog:read`, spend the catalogue rate budget, carry `readOnlyHint`, and carry
`_meta["anthropic/alwaysLoad"]: true` so a client that defers tools still loads
them first. The `url` is absolute, from the configured site URL (`SITE_URL`,
`Deps.SiteURL`; the Issuer when it is empty), never from the request, with the
locale prefix of the call (`Deps.Locale` is the translator's own, so a prefix
cannot disagree with the content). A result and the record it names carry the
same URL, tested for all three types. `page` and `perPage` follow the other
catalogue tools (a size above the limit is clamped).

**`fetch` is the catalogue operations, as one document.** It calls `getCourse`,
`getCourseLesson` or `getBlogPost` through the `Catalogue` and writes what they
return as Markdown-ish text: a course's landing copy, lessons (each with the id
that fetches it), modules and questions; a lesson's theory, exercise, definitions
and tips; and, as #851 decided for `getCourseLesson`, its starter code,
reference solution and tests (the HTTP read gives them to every visitor, and
fetch adds no rule of its own); a post's text with the markup dropped and its
paragraphs kept, and the courses it promotes. `metadata` holds the kind, slugs,
position, reading time and the like. The text is held to the response limit of
#851: the document is in the structured content and in the text block, so each
gets half of `MCP_MAX_RESPONSE_BYTES`, the longest string is cut with the usual
`truncation` hint, and the hint says the full text is on the website, at `url`.
An id that is malformed, or that names nothing a visitor may read (missing,
draft, archived, in another locale), is a tool error that says which; a landing
page that was archived and moved says where.

**A finding, not changed here.** `getCourse` and `getCourseLesson` answer a
request for a landing page that is a draft, or for a lesson of a course that is
not completed, whatever its state: the website does not link to them, but the
URL works (`coursepages` checks the slug and the locale, and an archived page's
redirect, not the state). The HTTP behaviour is not this ticket's to change, so
`fetch` guards against it instead: a course is fetched only when its landing
page is published and its Course completed, a lesson only when its Course is
completed, so `fetch` is never more permissive than `search`. Whether the HTTP
operations should answer 404 for those is a question for the owner of the
public course pages.

**Not done here.** The site has no search box; the operation is in the contract
and the client for whoever builds one. Books have no tool (see #851).

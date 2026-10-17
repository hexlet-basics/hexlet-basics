---
status: accepted
---

# An in-process OAuth 2.1 authorization server for MCP clients

Amends [ADR-0001](0001-contract-first-pipeline-with-ogen.md),
[ADR-0017](0017-server-side-sessions.md) and
[ADR-0015](0015-cutover-scope-and-api-prefix.md).

(Written against ADR-0003's JWT cookie and rebased onto ADR-0017's server-side
sessions, which replaced it; the reasoning below holds for the session cookie.)

Remote MCP clients (Claude Code, Cursor) cannot use the session cookie, which
a browser holds and a client has no way to obtain; the MCP authorization specification (revision 2026-07-28) expects the
server to sit behind an OAuth 2.1 authorization server. We build one as a
self-contained package, `internal/oauthserver`, inside `cmd/server`, on the
`zitadel/oidc` `op` library. It owns its tables and configuration and reaches
the rest of the application through one narrow seam, so it can be lifted out
later; it does not import `internal/handlers`.

**Shape.** The package is given two things by the application: who is signed in
(the existing session) and what the consent screen shows. Everything else, from
client registration to token storage, is inside it. It runs in the server
process, so ADR-0010 is unchanged: no third runtime, no second database.

**Protocol surface.**

- RFC 8414 metadata, advertising `registration_endpoint`, S256 as the only PKCE
  method, and `authorization_response_iss_parameter_supported`. Protected
  resource metadata (RFC 9728) is the MCP endpoint's side (ADR-0019).
- PKCE S256 is mandatory for every client. Redirect URIs match exactly, except
  that a loopback redirect may use any port (RFC 8252), because native clients
  pick a free port per run.
- The `resource` parameter (RFC 8707) is honoured and bound to the token's
  audience, which is the canonical `/mcp` URL. The authorization response
  carries `iss` (RFC 9207).
- Refresh tokens rotate. Presenting an already-used refresh token revokes the
  whole family, access tokens included.
- Clients register either by Client ID Metadata Document or by Dynamic Client
  Registration; we support both. Registration is anonymous, so it is
  constrained: `authorization_code` and `refresh_token` grants only, `https` or
  loopback redirect URIs only, a fixed list of requestable scopes, no
  caller-chosen audience, and a rate limit. Fetching a metadata document is
  guarded against SSRF and cached.
- Tokens are opaque, and a database leak alone must yield no usable token. The
  library mints an access token by encrypting its id and subject with the
  provider's key, so the database holds ids, not secrets; refresh tokens are
  chosen by our storage and are kept only as hashes. The MCP endpoint verifies
  a token by looking its id up, so revocation is immediate. That key is now a
  secret to protect and rotate. Access tokens live one hour, refresh tokens
  thirty days.

**Endpoint layout.** The issuer is the bare canonical origin from configuration
(never derived from the `Host` header), with no path, so discovery stays at the
root `/.well-known/oauth-authorization-server`. The library mounts `/auth`,
`/keys`, `/revoke`, `/userinfo` and `/end_session` at the root by default;
every authorization-server endpoint is instead mounted under `/oauth/*`, with
the library defaults overridden, because the root belongs to the SSR server.
The consent screen is an SSR route outside that prefix (`/connect`), and the
list of a user's Connected Apps lives at `/account/connected_apps`. The
issuer is baked into every registration and token, so it is the one value here
that is expensive to change.

**Ingress.** ADR-0015 sent only `/api/*` to Go. The ingress now also sends
`/mcp`, `/oauth/*`, `/.well-known/oauth-*` and `/.well-known/openid-configuration`
to Go. Clients that miss the RFC 8414 document fall back to the OpenID one,
which the library serves, so both paths answer with the same document. None of
`/mcp`, `/oauth` and `/connect` exists as a legacy URL or an SSR page.

**What stays outside the contract.** ADR-0001 generates every Go operation from
TypeSpec. The endpoints an RFC defines (`/oauth/authorize`, `/oauth/token`,
`/oauth/register`, the discovery documents) use form-encoded bodies, redirects
and fixed paths that TypeSpec does not express, so they are hand-written, and
that is deliberate: do not move them into the contract. What our own frontend
calls is still contract-first (ADR-0011): the data behind the consent screen,
approving or denying it, listing Connected Apps and revoking one are TypeSpec
operations guarded by the session scheme (the two writes also sit behind the
pipeline's cross-origin protection, ADR-0017).

**A second credential.** The session cookie (ADR-0017) is not the only credential
anymore: an OAuth bearer token is accepted, but only on `/mcp`, never on an
`/api` operation. A token whose subject is a client rather than a user is
rejected there. The consent screen can only be completed with the cookie
session, on the canonical host; a visitor who signed in on another of our
domains signs in once more there.

**Scopes.** A Scope records what a user consented to give a Connected App, and
is deliberately coarse: `catalog:read`, `progress:read`, `admin`. What a call
may actually do comes from the user's Permissions at the moment of the call
(ADR-0018); a Scope can only narrow that, never widen it.

## Considered Options

- **Ory Hydra v26.2.0 as a separate process.** Spiked with its login and
  consent delegation. Its login/consent split fits an app that already owns
  sign-in, but it leaves gaps we would fill ourselves anyway:
  - the discovery document has no `registration_endpoint`, no `iss` flag and no
    Client ID Metadata Documents (open upstream issue);
  - `resource` is ignored, and a dynamically registered client without
    `audience` set at registration fails every refresh with "audience has not
    been whitelisted";
  - registration is open to anyone and accepted a foreign audience, any scope,
    a plain-`http` non-loopback redirect, a wildcard redirect, and the
    `client_credentials` and `implicit` grants;
  - PKCE is not required unless `oauth2.pkce.enforced_for_public_clients` is
    set.

  The fixes are a registration gateway and our own metadata in front of it,
  plus a third runtime with its own database, which is what ADR-0010 avoided.
  Rejected.
- **Keycloak, Zitadel, Authentik, Dex.** Read from vendor documentation and
  repositories on 2026-10-16, not tested. Keycloak has the most MCP
  documentation but its Client ID Metadata and resource-indicator support is
  experimental, it cannot delegate login to our own screens without a plugin or
  brokering, and it is the heaviest to run. Zitadel is AGPL-3.0 and keeps its
  own users. Authentik and Dex lacked the needed registration and metadata
  support.
- **fosite.** The standalone repository has had no release since v0.49.0
  (December 2024); development continues inside the Hydra repository. It is a
  toolkit, so the storage and endpoints we must write are the same. Importing
  the copy inside Hydra was not attempted: its `go.mod` uses `replace`
  directives, which a consumer does not inherit.
- **zitadel/oidc v3.51.16 (chosen).** Active, Apache-2.0, and spiked through its
  example server: the authorization-code flow with S256-only PKCE, refresh
  rotation, introspection and revocation all work without the `openid` scope.

## Consequences

- The code that is ours: the `Storage` implementation on ent, scope admission,
  refresh-family revocation, `iss`, audience binding, both registration paths,
  the metadata document and token verification. This is more than the Hydra
  route would have been, in exchange for no extra process.
- Seen in the spike and still to solve: the library silently drops a custom
  scope such as `catalog:read` unless the client admits it; it issues an ID
  token even without `openid`, which MCP clients ignore; its example storage
  does not detect refresh-token reuse.
- Not verified yet, to be checked against the real clients before the feature
  is switched on: whether Claude Code and Cursor choose Client ID Metadata or
  Dynamic Registration, how they handle loopback redirects, and whether they
  check `iss`. A wildcard redirect registered by a client was accepted by
  Hydra's registration; whether it matched at authorization was not tested.

## Addendum: the tracer slice (#836)

What the first implementation settled, and where it stops short of the text above.

- **Switch and configuration.** `OAUTH_ENABLED` (default off), `OAUTH_ISSUER`
  (a bare origin; `https` required in production) and `OAUTH_ENCRYPTION_KEY`
  (at least 32 characters). The last two are required only when the flag is on,
  so the default configuration needs none of them.
- **Seam.** The package is given `oauthserver.Application`, which names the
  signed-in user. The consent screen's three operations are TypeSpec operations
  under `/api/oauth/consents/{requestId}` (read; `approve` and `deny` are
  also behind the cross-origin protection); the handlers call `Pending`, `Approve` and `Deny`. A test
  runs `go list -deps` and fails if the package depends on `internal/handlers`.
- **ID-token signing key.** An Ed25519 key derived from the encryption secret,
  not a stored or per-process RSA key: every replica signs and verifies alike
  and there is no second secret in the database. Nothing authorizes with an ID
  token.
- **Access tokens** are the library's encrypted compact JWE of `id:subject`. The
  table holds ids, refresh tokens and authorization codes are kept as SHA-256
  hashes, and `VerifyAccessToken` decrypts and then looks the id up.
- **PKCE.** S256 and a challenge are enforced in `CreateAuthRequest`; the
  library alone accepts a missing challenge and `plain`.
- **Closed library surface.** `/oauth/introspect` (every client is public) and
  `/oauth/device_authorization` answer 404, and the metadata document is the
  library's narrowed to `code`, `authorization_code` and `refresh_token`, S256
  and public clients.
- **A client row lists the custom scopes it may request.** That is the minimum
  the library needs to keep `catalog:read` instead of dropping it; the policy
  for what a registration may list is still scope admission.
- **Not yet done**, each for its own ticket: `registration_endpoint` (not
  advertised, since it would answer 404), the `iss` authorization-response
  parameter, the `resource` parameter and audience binding (the audience is the
  client for now), and refresh-token reuse detection. Loopback redirects already
  match on any port, because the library does that for native clients.

## Addendum: scopes, resource, `iss` and loopback redirects (#841)

- **Canonical MCP URL.** One value, `OAUTH_ISSUER` plus `/mcp`
  (`Server.Resource()`), built once from configuration and never from the `Host`
  header; there is no second setting to drift from the issuer. It is the
  audience of every access token.
- **Scope admission.** The Scopes `catalog:read`, `progress:read` and `admin`,
  plus `openid` and `offline_access`, are admitted; a custom Scope is admitted
  only when the client row lists it. Anything else (an unknown name, `profile`,
  `email`, a Scope the client was not registered for) is refused with
  `invalid_scope` redirected to the client, with the supported list in the
  description. We chose a refusal over the library's silent drop because a client
  that believed it held `admin` would fail later in a confusing way; the cost is
  that a client sending generic `profile email` scopes must stop. The check is
  an `AuthorizeValidator` over the library's own validation
  (`internal/oauthserver/policy.go`). The flow does not need `openid`; without
  it the library still issues an ID token, which clients ignore.
- **`resource` (RFC 8707).** Every value of `resource` at the authorization and
  token endpoints must equal the canonical URL, otherwise `invalid_target` (a
  redirect at the authorization endpoint, a 400 at the token endpoint, which does
  not spend the code). An absent `resource` is accepted and bound to the same
  URL, since this server has a single resource. The audience is not stored per
  token: `VerifyAccessToken` returns `Token.Audience` as the canonical URL, which
  a resource server must find itself in. Deleting the row still revokes at once.
  The ID token's `aud` carries the URL as well as the client id.
- **`iss` (RFC 9207).** Every redirect to the client from the authorization
  endpoint, its callback and the consent denial carries `iss`, errors included;
  the hop to the consent screen does not. The metadata document, built here for
  both well-known paths, declares `authorization_response_iss_parameter_supported`.
- **Loopback redirects.** The library, not our code, matches them: for a native
  client, once both the registered and the requested URI are loopback hosts, only
  path and query are compared, so any port matches, and so do another loopback
  address or name (`localhost`, `[::1]`) and the other scheme. That reaches only
  the user's own machine. A non-loopback redirect must equal a registered one
  exactly, and a plain-`http` non-loopback redirect is refused even if
  registered. Tests pin all of this.

## Addendum: dynamic client registration (#840)

- **Endpoint.** `POST /oauth/register` (RFC 7591, JSON in and out, `201` with `Cache-Control: no-store`), advertised as `registration_endpoint` in both metadata documents. It is ours entirely (`internal/oauthserver/registration.go`); the library has none. A client row is public by construction, so a registration can only ever produce a public client with the `authorization_code` and `refresh_token` grants and the `code` response type.
- **Refusals** (`400`, RFC 7591 `invalid_redirect_uri` or `invalid_client_metadata`, nothing stored): a redirect URI that is not `https` or `http` on a loopback host (`localhost`, `127.0.0.0/8`, `::1`), is relative, has a custom scheme, a fragment, userinfo (the `http://127.0.0.1@evil` lookalike) or a wildcard (literal or `%2A`), or more than ten redirect URIs; any grant other than the two above, or `refresh_token` without `authorization_code`; any response type but `code`; a `token_endpoint_auth_method` other than `none`; any scope outside `catalog:read`, `progress:read`, `admin` and `offline_access` (`openid` and `profile` included: a client row cannot use them); a body that is not one JSON object, is not `application/json` (`415`) or exceeds 16 KiB (`413`). The request is judged by an allow-list of fields: `client_secret`, `audience`, `resource`, `client_id`, `jwks`, `software_statement`, `skip_consent` and any other unknown field are refused, never ignored, because a field ignored today could gain meaning by accident. Harmless descriptive metadata (`client_uri`, `logo_uri`, `tos_uri`, `policy_uri`, `contacts`, `software_id`, `software_version`, `application_type`) is accepted and dropped.
- **Default scopes.** A registration that names no scope may later request `catalog:read` and `progress:read`: read-only, never `admin`, which a client must ask for by name (the user still consents, and their Permissions still decide, ADR-0018). The client row stores the custom Scopes only; `offline_access` is admitted for every client.
- **Name.** `client_name` is stored as plain text, trimmed, control characters removed, at most 100 characters, `Unnamed application` when absent. It is what the consent screen shows, so it is untrusted text, and escaping stays the screen's job.
- **Rate limit.** Nothing in the repository limited requests, so the package has the smallest limiter that does the job: a fixed window per client address, in memory, `OAUTH_REGISTRATION_LIMIT` attempts (default 10) per `OAUTH_REGISTRATION_WINDOW` (default `1h`), a `429` with `Retry-After` beyond it. Every attempt spends it, refused ones included, so probing the validation is not free. The address comes from the existing client-IP middleware (`handlers.WithClientIP`, now also wrapped around the OAuth routes) and reaches the package through `Deps.ClientAddress`, never from a header read in the package; an unresolved address shares one bucket, and an IPv6 address counts as its /64. Memory is bounded: past 100,000 addresses in one window new addresses wait for the window to pass. Being in memory, the budget is per replica. There is no separate cap on the clients one address may own: ten an hour from one address is the cap, and an unused row costs a few hundred bytes; pruning clients that never completed a flow is left for when the table shows it needs it. The client-IP middleware trusts the forwarding headers as `realip` does (ADR-0015 chose it for attribution); an attacker who can put a public address first in `X-Forwarded-For` on a request that reaches the pod directly would pick their own bucket, which is the ingress's job to prevent by overwriting the header.
- **Supersedes** the "not yet done" note on `registration_endpoint` in the #836 addendum.

## Addendum: ingress and configuration (#842)

- **Routing.** The Helm ingress sends `/mcp`, the four well-known documents (`oauth-authorization-server`, `openid-configuration`, `oauth-protected-resource` and its `/mcp` form) and the `/oauth/` prefix to the `api` service, ahead of the `/` catch-all (Yandex ALB takes the first matching path), on every `full` host only; hosts without `full` (code-basics.ru) keep `/` alone. The main Ingress sets `use-regex`, so its Exact paths are RE2 regexes and the dots of `.well-known` are escaped. `/connect` and `/account/connected_apps` stay on the web service. `mise run lint:helm` asserts the per-host path list.
- **The ingress does not follow the flags.** `OAUTH_ENABLED` and `MCP_ENABLED` default to `"false"` in `goEnv`, and the paths are routed either way: with the flags off Go mounts nothing and answers its ordinary 404 (a test pins every path), so the same release serves both states and turning the feature on is a configuration change, not a routing one.
- **Configuration.** `OAUTH_ISSUER` (the canonical origin, `https://code-basics.com`), `OAUTH_REGISTRATION_LIMIT`, `OAUTH_REGISTRATION_WINDOW` and `AUDIT_RETENTION_MONTHS` are plain `goEnv` values shared by the server and the worker (the partition maintenance job runs in the worker). `OAUTH_ENCRYPTION_KEY` is a secret read from `secrets.code_basics.env` like the others; it must be added to the sops file (`legacy/k8s/secrets.yaml`) before the OAuth flag goes on.
- **`X-Forwarded-For`.** The registration limit trusts the address `realip` takes from this header, so the proxy must set it from the real peer and never pass a client's value through. The chart cannot configure that on the ALB and it was not verified here: check it on the live ALB before enabling OAuth.

## Addendum: Client ID Metadata Documents (#847)

- **Flow.** A `client_id` that is an `https` URL is resolved at the authorization endpoint (`internal/oauthserver/clientdocument.go`), before the library looks the client up: the document is fetched the first time the id is seen, judged, and stored as an ordinary row of `oauth_clients` (the URL is its `client_id`), so consent, tokens and refresh work unchanged. It is accepted only if its own `client_id` equals the URL character for character; the request's redirect URI must then be listed in it, which the library's usual redirect check enforces. The token endpoint never fetches: a flow already under way, and a refresh token already issued, do not depend on a remote host being up.
- **Same constraints as registration.** The document goes through the same allow-list and validation as a registration body (`parseMetadata`, shared with `register`): `https` or loopback redirects, no wildcards, `authorization_code`/`refresh_token` and `code` only, `token_endpoint_auth_method` `none`, the same scope list, and any unknown field (a secret, an audience, `jwks_uri`) is refused. The one extra field is `client_id`. A stricter reading than some clients' documents may assume; a client that fails here is told which field.
- **The URL.** `https` only, a path other than `/`, no credentials, no fragment, no dot segments, the default port only, at most 2048 characters (a query string is allowed, as the draft allows it). A URL that fails is refused before any request leaves the server.
- **SSRF guard.** The destination is judged where the socket connects: the dialer's `Control` hook sees the resolved IP for every candidate address, so a name that rebinds between a check and the connection (DNS rebinding) cannot pass, and a literal address, a name resolving to loopback and an IPv4-mapped IPv6 form are all caught the same way. Only a global unicast address outside the blocked ranges passes: loopback, RFC 1918 private, carrier-grade NAT, link-local (including the 169.254.169.254 metadata service), unique-local, multicast, unspecified, reserved and documentation ranges, and the IPv6 forms that embed an IPv4 address (IPv4-mapped, NAT64, 6to4, Teredo). No proxy is used, redirects are not followed (a 3xx is a refusal, not a hop to re-validate), the answer must be a `200` `application/json` of at most 16 KiB, and the whole fetch has five seconds.
- **Caching.** The cached copy lives in the client row (`metadata_expires_at`, null for any other client). Its lifetime comes from `s-maxage`, `max-age` or `Expires`, 1 hour without any of them, clamped between 5 minutes and 24 hours; `no-store` and `no-cache` mean the 5-minute floor, because the row has to exist to run the flow. **Refresh:** the next authorization request after expiry fetches the document again and replaces the row's name, redirect URIs and scopes. If that fails, the authorization request is refused and the row is kept, so an application whose document can no longer be verified starts no new grant, while grants it already holds go on until the user revokes them. Failures are remembered in memory for 5 minutes per URL (bounded at 10,000), so an unreachable or hostile URL is not fetched again on every request.
- **Rate limit.** The fetch is triggered by an anonymous request, so each fetch that would actually leave the server spends one attempt of a per-address limiter, a second instance of the registration limiter with the same limit and window (a `429` with `Retry-After` beyond it); an id served from the row or from the failure memory costs nothing. It has the same limits as the registration limiter: in memory, per replica, and the address only comes from `Deps.ClientAddress`.
- **Test seam.** `oauthserver.New` takes options; `WithMetadataHTTPClient` replaces the fetching client and exists for tests, since a test server on loopback is refused by design. Nothing in the application wiring passes it, and a test through the production client shows that loopback, private, link-local, metadata-service and IPv4-mapped targets are refused.
- **Consent.** The consent screen already shows the redirect host (with its port); it now also shows a visible warning when that host is `localhost`, a loopback address, or `::1`, since an approval there lands on a program running on the user's own machine.
- **Not done.** A document's `client_name` is untrusted text like a registered one; the screen shows the redirect host beside it, but not the host of the `client_id` URL. Deleting a client row of this kind lets the next authorization request recreate it from the document, which is intended.

## Addendum: refresh rotation, replay detection and revocation (#839)

- **Rotation.** Every refresh returns a new access and refresh token; the presented refresh token is marked used (`used_at`) and its paired access token is deleted. A used refresh token is kept, not deleted, until its own expiry: the row is the tombstone that lets a replay be told from an unknown token.
- **Family.** Every token descends from one authorization and carries its `family_id` (a random id made at the code exchange, inherited by every rotation; migration `20261019090000_oauth_token_families.sql` backfills existing rows as one family each). Presenting a used refresh token revokes the whole family: every access token and every refresh token, used or live, with the family id is deleted, the freshly issued ones included. The thief and the victim are both locked out and the user signs in again, which is what the MCP authorization spec asks for. Other sign-ins of the same user, other clients and other users are separate families and are untouched.
- **Strict reuse detection, no grace window.** Two exchanges of one refresh token (a client retry after a lost response) race on one conditional `UPDATE ... WHERE used_at IS NULL`, so the database lets exactly one claim the token; the other is treated as a replay and revokes the family, the winner's child included. We chose this over a short grace window (accept the second presentation for a few seconds and return the same child) because a grace window needs the previous child's secret to be recoverable, and the database holds only hashes; keeping them recoverable, or issuing two live children, would weaken the property that a stolen token is detected. The cost is that a client that retries a refresh after a lost response is signed out and must reconnect; MCP clients reconnect through the browser, and such a loss is rare. Two parallel exchanges therefore never leave two live children (tested on a real connection pool).
- **Revocation.** `/oauth/revoke` (RFC 7009) ends an access token at once (the row is deleted, so `VerifyAccessToken` refuses it; the client may still refresh) and ends the whole family when given a refresh token, used ones included, since a user withdrawing a connection means every token of it. A token already gone answers `200`; only the client a token was issued to may revoke it.
- **Lifetimes.** Access tokens live one hour (`OAUTH_ACCESS_TOKEN_TTL`) and refresh tokens thirty days (`OAUTH_REFRESH_TOKEN_TTL`, `720h`); each rotation grants a fresh full refresh lifetime. With OAuth enabled the access lifetime must be at least a minute and the refresh lifetime longer than it, otherwise configuration fails to load.
- **Cleanup.** The periodic worker job `purge_oauth_tokens` (hourly, also at start) deletes expired access and refresh tokens, used tombstones among them, and authorization requests that expired uncompleted. Expiry is also checked on every use, so the sweep only bounds table growth.
- **Supersedes** "refresh-token reuse detection" in the #836 "not yet done" note.

## Addendum: Connected Apps (#848)

- **What the list is.** A Connected App is one OAuth client the signed-in user holds a live grant to: an unused, unexpired refresh token, or an unexpired access token (a grant made without `offline_access`, which lives out its hour). A used refresh token is the replay tombstone of #839, not a grant, so it never keeps an app listed on its own. Several sign-ins of one user to one client are one app. The list shows the client's display name as stored (untrusted text, escaped by the page), the Scopes of its newest live grant, when it was connected, and, for a client-document app, the host of its `client_id` (the one fact its author does not pick through a name; absent for other clients, whose ids are opaque). "Connected at" is the approval time (`auth_time`, which every rotation inherits) of the oldest grant still on file, so it can move forward once an old family's tombstones expire.
- **Contract.** `GET /api/oauth/connected_apps` (session) and `DELETE /api/oauth/connected_apps/{id}` (session; the revoke is also behind the cross-origin protection), handled by `ListConnectedApps` and `RevokeConnectedApp`, which reach the package through the `OAuthConsents` seam (`ConnectedApps`, `Disconnect`); the package still does not import handlers. The handle is the client row's integer id, not the `client_id`, because a metadata-document `client_id` is a URL that does not belong in a path segment. A guest gets 401; with the OAuth server off the list is empty (no app can exist) and a revoke is 404.
- **Revoking.** `Disconnect` runs the same primitive as `TerminateSession` (#839): in one transaction every access token and every refresh token of that user and client, in every family and used tombstones included, is deleted. The effect is immediate: the access token fails `VerifyAccessToken` at its next use, and a refresh with any old refresh token is `invalid_grant`. Other users' grants to the same client and the same user's other clients are untouched. An id that is not one of the caller's listed apps (another user's, unknown, already revoked or lapsed) is a 404, one answer for all of them. The revoke leaves the client row alone, so the app may ask for consent again.
- **Lapse after thirty days: what "last used" means.** There is no last-used column and no sliding window on access. A grant lives as long as its newest refresh token, which is thirty days from its issue, and every refresh issues a new one. So an app stays connected while it **refreshes** at least once in thirty days, and an app that has not refreshed for thirty days has no live refresh token: it cannot obtain access, and it drops off the list on its own. "Used" therefore means a refresh, **not an MCP call**: an access token is good for an hour, so a client that is calling tools is refreshing hourly anyway, but an app holding a still-valid access token and never refreshing lapses on schedule regardless of calls it made. The page does not show a "last used" time for the same reason, since the only honest value would be the last refresh. Pinned with the test helper's clock: 30 days minus a minute after the last refresh it is listed, two minutes later it is not and its refresh token is refused; an app refreshed every five days is still listed at day 35.

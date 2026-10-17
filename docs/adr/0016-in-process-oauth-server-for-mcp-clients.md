---
status: accepted
---

# An in-process OAuth 2.1 authorization server for MCP clients

Amends [ADR-0001](0001-contract-first-pipeline-with-ogen.md),
[ADR-0003](0003-jwt-auth-via-go-pkgz-auth.md) and
[ADR-0015](0015-cutover-scope-and-api-prefix.md).

Remote MCP clients (Claude Code, Cursor) cannot use the JWT cookie with its XSRF
header; the MCP authorization specification (revision 2026-07-28) expects the
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
  resource metadata (RFC 9728) is the MCP endpoint's side (ADR-0018).
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
operations guarded by the session and XSRF schemes.

**A second credential.** ADR-0003's cookie JWT is not the only credential
anymore: an OAuth bearer token is accepted, but only on `/mcp`, never on an
`/api` operation. A token whose subject is a client rather than a user is
rejected there. The consent screen can only be completed with the cookie
session, on the canonical host; a visitor who signed in on another of our
domains signs in once more there.

**Scopes.** A Scope records what a user consented to give a Connected App, and
is deliberately coarse: `catalog:read`, `progress:read`, `admin`. What a call
may actually do comes from the user's Permissions at the moment of the call
(ADR-0017); a Scope can only narrow that, never widen it.

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

import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";

using Http;
using Rest;
using OpenAPI;

namespace HexletBasics;

// ---------------------------------------------------------------------------
// OAuth consent screen (ADR-0019).
//
// The protocol endpoints an RFC defines (`/oauth/authorize`, `/oauth/token`,
// the discovery documents) are hand-written outside this contract on purpose:
// they use form bodies, redirects and fixed paths TypeSpec does not express.
// What our own frontend calls is contract-first like everything else: the data
// behind the consent screen, and approving or denying it. All three are guarded
// by the session cookie (the two writes also sit behind the pipeline's
// cross-origin protection, ADR-0017), because the consent screen can only be completed by the signed-in user on the canonical
// host.
// ---------------------------------------------------------------------------

/** What the consent screen shows for one authorization request. */
model OAuthConsent {
  /** Name of the application asking for access. */
  clientName: string;

  /** Host the user is sent back to after deciding. */
  redirectHost: string;

  /** What the application asks for: OAuth scopes (`catalog:read`, `progress:read`, `admin`, `offline_access`). */
  scopes: string[];

  /** The signed-in account that would grant it. */
  accountName: string;
}

/** Where the browser goes next after the user has decided. */
model OAuthDecision {
  /** Absolute URL: the authorization endpoint's callback after approving, the application's redirect URI (with `error=access_denied`) after denying. */
  redirectUrl: string;
}

@route("/oauth/consents/{requestId}")
@tag("OAuth")
interface OAuthConsents {
  /** The request the browser arrived at the consent screen with. 404 when it is unknown, expired or already decided. */
  @get
  @useAuth(UserSession)
  @operationId("getOAuthConsent")
  get(@path requestId: string): OAuthConsent | UnauthorizedError | NotFoundError | ApiError;

  /** Grant the request as the signed-in user. */
  @post
  @route("/approve")
  @useAuth(UserSession)
  @operationId("approveOAuthConsent")
  approve(@path requestId: string): OAuthDecision | UnauthorizedError | NotFoundError | ApiError;

  /** Refuse the request; the application is told `access_denied`. */
  @post
  @route("/deny")
  @useAuth(UserSession)
  @operationId("denyOAuthConsent")
  deny(@path requestId: string): OAuthDecision | UnauthorizedError | NotFoundError | ApiError;
}

// ---------------------------------------------------------------------------
// Connected Apps (ADR-0019, #848): the OAuth clients the signed-in user has
// approved, and withdrawing one. Both are session-guarded and show only the
// caller's own apps; the revoke needs the XSRF header too.
// ---------------------------------------------------------------------------

/** One application the signed-in user has connected. */
model ConnectedApp {
  /** Handle to revoke it by (not the OAuth client_id, which can be a URL). */
  id: int32;

  /** Display name the application registered with; untrusted text. */
  name: string;

  /** What it may do: OAuth scopes (`catalog:read`, `progress:read`, `admin`, `offline_access`). */
  scopes: string[];

  /** When the user first connected it (the oldest grant still on file). */
  connectedAt: utcDateTime;

  /** Host of the application's client_id when that is a metadata document URL; absent otherwise. */
  clientHost?: string;
}

/** The signed-in user's Connected Apps, newest first. */
model ConnectedApps {
  items: ConnectedApp[];
}

@route("/oauth/connected_apps")
@tag("OAuth")
interface OAuthConnectedApps {
  /** The apps the signed-in user has connected and that can still act for them. An app unused for thirty days (no refresh) has lapsed and is not listed. */
  @get
  @useAuth(UserSession)
  @operationId("listConnectedApps")
  list(): ConnectedApps | UnauthorizedError | ApiError;

  /** Revoke one app at once: every access and refresh token the user granted it ends. 404 when it is not one of the caller's apps. */
  @delete
  @route("/{id}")
  @useAuth([UserSession, XsrfToken])
  @operationId("revokeConnectedApp")
  revoke(@path id: int32): NoContent | UnauthorizedError | NotFoundError | ApiError;
}

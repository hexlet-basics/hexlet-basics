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

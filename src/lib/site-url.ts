import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequest } from "@tanstack/react-start/server";

// The site's public origin, for the absolute URLs a page's head carries
// (canonical links), as legacy's `*_url` helpers built them.
//
// On the server SITE_URL wins: behind the ingress the SSR process sees the
// internal request (TLS terminated at the load balancer), so the request's own
// origin would say http. Without it — development, where Vite serves the page
// on the origin the browser uses — the request origin is the public one. In the
// browser the page's own origin is the public one by definition. Server-only
// imports are stripped from the client bundle by the isomorphic split.
const getSiteOrigin = createIsomorphicFn()
  .server(() => process.env.SITE_URL ?? new URL(getRequest().url).origin)
  .client(() => window.location.origin);

// The absolute URL of a path on this site, e.g. a route match's pathname.
export function siteUrl(pathname: string): string {
  return new URL(pathname, getSiteOrigin()).toString();
}

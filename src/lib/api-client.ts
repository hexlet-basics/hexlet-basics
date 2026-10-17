import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequest, getRequestHeader } from "@tanstack/react-start/server";
import { client } from "@/client/client.gen";
import { apiErrorOf } from "@/lib/api-error";
import { localeFromPathname } from "@/lib/locale-path";

// Read the incoming SSR request's cookie header. Server-only (the client impl
// is a no-op), so the server import is stripped from the browser bundle — the
// browser sends the auth cookies itself through Axios `withCredentials`.
const getRequestCookie = createIsomorphicFn()
  .server(() => getRequestHeader("cookie"))
  .client(() => undefined);

// Resolve the active URL locale at request time. Reading inside the interceptor
// is request-safe during SSR and avoids mutating the generated client singleton
// when concurrent renders use different locales.
const getRequestLocale = createIsomorphicFn()
  .server(() => localeFromPathname(new URL(getRequest().url).pathname))
  .client(() => localeFromPathname(window.location.pathname));

// The generated hey-api client is a singleton, and every contract path starts
// with `/api` (ADR-0015). In the browser the API is the site's own origin — the
// ingress routes `/api/*` to Go in production and the Vite proxy does it in
// development. hey-api puts `baseUrl` into every query key, so it must be the
// same constant during SSR and in the browser: a server-only value would make
// the dehydrated keys differ from the browser's, the hydrated cache would never
// match, and each page would refetch and fail hydration. A bare relative URL
// cannot be it, because the fetch client builds a `Request` before any
// interceptor runs and Node cannot parse a relative one. So the base is a
// placeholder origin that the request interceptor below always replaces: with
// the internal API_URL on the server, with the site's own origin (a relative
// URL) in the browser.
const PLACEHOLDER_ORIGIN = "http://api.invalid";

client.setConfig({ baseUrl: PLACEHOLDER_ORIGIN });

// Forward the SSR request's cookie so authenticated pages render server-side
// with the user's session (ADR-0008). Registered once at module scope — never
// per-request — and safe on the singleton because the cookie is read inside the
// request-scoped isomorphic fn (Start's async request context), not by mutating
// shared client state.
client.interceptors.request.use(async (request) => {
  const { pathname, search } = new URL(request.url);
  // The dead branch is stripped per build target, so `process` never leaks
  // into the client bundle.
  const target = import.meta.env.SSR
    ? `${(process.env.API_URL ?? "http://localhost:3001").replace(/\/$/, "")}${pathname}${search}`
    : `${pathname}${search}`;
  const headers = new Headers(request.headers);
  const cookie = getRequestCookie();
  if (cookie) headers.set("cookie", cookie);
  headers.set("accept-language", getRequestLocale());
  // A new Request is the way to change the URL; the body is read out because a
  // streamed one cannot be handed to a new Request without `duplex`.
  const body = request.body ? await request.arrayBuffer() : undefined;
  return new Request(target, {
    method: request.method,
    headers,
    body,
    signal: request.signal,
    redirect: request.redirect,
    credentials: request.credentials,
  });
});

// A non-2xx answer reaches callers as one ApiError (status, parsed body), so
// no route imports the HTTP library.
client.interceptors.error.use(apiErrorOf);

export { client };

import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequest, getRequestHeader } from "@tanstack/react-start/server";
import { client } from "@/client/client.gen";
import { localeOf } from "@/lib/locale-path";

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
  .server(() => localeOf(new URL(getRequest().url).pathname))
  .client(() => localeOf(window.location.pathname));

// The generated hey-api client is a singleton, and every contract path starts
// with `/api` (ADR-0015). In the browser the API is the site's own origin — the
// ingress routes `/api/*` to Go in production and the Vite proxy does it in
// development — so the base URL is empty. The config stays empty during SSR
// too: hey-api puts `baseURL` into every query key, so a server-only value would
// make the dehydrated keys differ from the browser's, the hydrated cache would
// never match, and each page would refetch and fail hydration. The server
// reaches Go over the internal network (API_URL) through the request
// interceptor below instead.
client.setConfig({ baseURL: "" });

// Forward the SSR request's cookie so authenticated pages render server-side
// with the user's session (ADR-0008). Registered once at module scope — never
// per-request — and safe on the singleton because the cookie is read inside the
// request-scoped isomorphic fn (Start's async request context), not by mutating
// shared client state.
client.instance.interceptors.request.use((request) => {
  // The dead branch is stripped per build target, so `process` never leaks
  // into the client bundle.
  if (import.meta.env.SSR) request.baseURL = process.env.API_URL ?? "http://localhost:3001";
  const cookie = getRequestCookie();
  if (cookie) request.headers.set("cookie", cookie);
  request.headers.set("accept-language", getRequestLocale());
  return request;
});

export { client };

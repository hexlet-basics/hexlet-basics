import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequest, getRequestHeader } from "@tanstack/react-start/server";
import { client } from "@/client/client.gen";
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

// The generated hey-api client is a singleton. During SSR the Node process
// reaches Go over the internal network (API_URL); in the browser it uses the
// public origin (VITE_API_URL). The dead branch is stripped per build target,
// so `process` never leaks into the client bundle.
client.setConfig({
  baseURL: import.meta.env.SSR
    ? (process.env.API_URL ?? "http://localhost:3001")
    : (import.meta.env.VITE_API_URL ?? "http://localhost:3001"),
  // Axios owns the browser-side XSRF protocol. Its defaults match
  // go-pkgz/auth; withXSRFToken opts the dev cross-origin API into the same
  // built-in behavior used for same-origin requests.
  withCredentials: true,
  withXSRFToken: true,
  xsrfCookieName: "XSRF-TOKEN",
  xsrfHeaderName: "X-XSRF-TOKEN",
});

// Forward the SSR request's cookie so authenticated pages render server-side
// with the user's session (ADR-0008). Registered once at module scope — never
// per-request — and safe on the singleton because the cookie is read inside the
// request-scoped isomorphic fn (Start's async request context), not by mutating
// shared client state.
client.instance.interceptors.request.use((request) => {
  const cookie = getRequestCookie();
  if (cookie) request.headers.set("cookie", cookie);
  request.headers.set("accept-language", getRequestLocale());
  return request;
});

export { client };

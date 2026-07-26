import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";
import { client } from "@/client/client.gen";

// Read the incoming SSR request's cookie header. Server-only (the client impl
// is a no-op), so the server import is stripped from the browser bundle — the
// browser sends the JWT cookie itself via `credentials: "include"`.
const getRequestCookie = createIsomorphicFn()
  .server(() => getRequestHeader("cookie"))
  .client(() => undefined);

// The generated hey-api client is a singleton. During SSR the Node process
// reaches Go over the internal network (API_URL); in the browser it uses the
// public origin (VITE_API_URL). The dead branch is stripped per build target,
// so `process` never leaks into the client bundle.
client.setConfig({
  baseUrl: import.meta.env.SSR
    ? (process.env.API_URL ?? "http://localhost:3001")
    : (import.meta.env.VITE_API_URL ?? "http://localhost:3001"),
  // Send the httpOnly JWT cookie to the Go API cross-origin (Go must reply with
  // CORS allow-credentials). Harmless on the Node SSR fetch.
  credentials: "include",
});

// Forward the SSR request's cookie so authenticated pages render server-side
// with the user's session (ADR-0008). Registered once at module scope — never
// per-request — and safe on the singleton because the cookie is read inside the
// request-scoped isomorphic fn (Start's async request context), not by mutating
// shared client state.
client.interceptors.request.use((request) => {
  const cookie = getRequestCookie();
  if (cookie) request.headers.set("cookie", cookie);
  return request;
});

export { client };

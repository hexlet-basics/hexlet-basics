import { redirect } from "@tanstack/react-router";
import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";

// Stub user shape until the Go auth contract lands (ADR-0003). Kept minimal on
// purpose — the real User model will come from TypeSpec, not be hand-invented.
export interface AuthUser {
  id: string;
}

// go-pkgz/auth issues the JWT in a cookie named `JWT` by default (ADR-0003).
const JWT_COOKIE = "JWT";

function readCookie(
  header: string | undefined,
  name: string,
): string | undefined {
  if (!header) return undefined;
  for (const pair of header.split(";")) {
    const [key, ...rest] = pair.trim().split("=");
    if (key === name) return rest.join("=");
  }
  return undefined;
}

// Resolve the current user from the httpOnly JWT cookie during SSR. STUB: the
// auth backend and a `/me` operation don't exist yet, so this returns null even
// when a token is present. When they land, verify the JWT / call `GET /me` here
// and map the result to AuthUser. Client-side always null — the cookie is
// httpOnly, so the browser relies on the SSR-resolved value.
export const getCurrentUser = createIsomorphicFn()
  .server((): AuthUser | null => {
    const token = readCookie(getRequestHeader("cookie"), JWT_COOKIE);
    void token; // TODO(auth): verify JWT / call GET /me, map to AuthUser.
    return null;
  })
  .client((): AuthUser | null => null);

// beforeLoad guard for protected routes (wired in Wave 3). Uses an untyped href
// because the `/login` route doesn't exist yet.
export function requireAuth(user: AuthUser | null, redirectHref: string): void {
  if (!user) {
    throw redirect({
      href: `/login?redirect=${encodeURIComponent(redirectHref)}`,
    });
  }
}

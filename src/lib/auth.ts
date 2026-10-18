import type { QueryClient } from "@tanstack/react-query";
import { redirect } from "@tanstack/react-router";
import { getCurrentUserOptions } from "@/client/@tanstack/react-query.gen";
import type { User } from "@/client/types.gen";

// The authenticated principal, sourced from the generated `User` model (the auth
// Backend + `GET /me` landed in ADR-0003). No hand-invented shape — the contract
// Owns it.
export type AuthUser = User;

// Resolve the current user for the router context. Runs in the root beforeLoad on
// Both server and client; the shared QueryClient caches the result so the render
// Tree and later navigations reuse it instead of re-fetching. `GET /me` returns
// 200 with `{ user: null }` for anonymous requests (never 401), so this never
// Throws and every page — public or guarded — can rely on it.
export async function resolveCurrentUser(queryClient: QueryClient): Promise<AuthUser | null> {
  // Static: `onIdentityChanged` (lib/identity) re-seeds this query when it changes.
  const { user } = await queryClient.query({ ...getCurrentUserOptions(), staleTime: "static" });
  return user;
}

// BeforeLoad guard for authenticated routes. Redirects anonymous callers to the
// Login page, preserving the intended destination.
export function requireAuth(user: AuthUser | null, redirectHref: string): void {
  if (!user) {
    throw redirect({
      href: `/session/new?redirect=${encodeURIComponent(redirectHref)}`,
    });
  }
}

// BeforeLoad guard for the admin area. Gates on `canAccessAdmin` (the serializer
// Flag legacy used for the admin namespace), not the raw `admin` column, so the
// Same rule as the Rails app decides who sees the back office. Anonymous callers
// Go to login; signed-in users without access are bounced to the home page rather
// Than looped back through login.
export function requireAdmin(user: AuthUser | null, redirectHref: string): void {
  requireAuth(user, redirectHref);
  if (!user?.canAccessAdmin) {
    throw redirect({ to: "/{-$locale}" });
  }
}

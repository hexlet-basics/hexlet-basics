import type { QueryClient } from "@tanstack/react-query";
import { redirect } from "@tanstack/react-router";
import { getCurrentUserOptions } from "@/client/@tanstack/react-query.gen";
import type { User } from "@/client/types.gen";

// The authenticated principal, sourced from the generated `User` model (the auth
// backend + `GET /me` landed in ADR-0003). No hand-invented shape — the contract
// owns it.
export type AuthUser = User;

// Resolve the current user for the router context. Runs in the root beforeLoad on
// both server and client; the shared QueryClient caches the result so the render
// tree and later navigations reuse it instead of re-fetching. `GET /me` returns
// 200 with `{ user: null }` for anonymous requests (never 401), so this never
// throws and every page — public or guarded — can rely on it.
export async function resolveCurrentUser(queryClient: QueryClient): Promise<AuthUser | null> {
  const { user } = await queryClient.ensureQueryData(getCurrentUserOptions());
  return user;
}

// beforeLoad guard for authenticated routes. Redirects anonymous callers to the
// login page, preserving the intended destination.
export function requireAuth(user: AuthUser | null, redirectHref: string): void {
  if (!user) {
    throw redirect({
      href: `/session/new?redirect=${encodeURIComponent(redirectHref)}`,
    });
  }
}

// beforeLoad guard for the admin area. Gates on `canAccessAdmin` (the serializer
// flag legacy used for the admin namespace), not the raw `admin` column, so the
// same rule as the Rails app decides who sees the back office. Anonymous callers
// go to login; signed-in users without access are bounced to the home page rather
// than looped back through login.
export function requireAdmin(user: AuthUser | null, redirectHref: string): void {
  requireAuth(user, redirectHref);
  if (!user?.canAccessAdmin) {
    throw redirect({ to: "/{-$locale}" });
  }
}

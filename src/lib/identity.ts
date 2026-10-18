import type { QueryClient } from "@tanstack/react-query";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import type { User } from "@/client/types.gen";

// What a visitor's identity changes is not one query but everything the site
// Has read for them: the course and lesson payloads carry their progress, the
// Profile is theirs, and a page cached as a guest says what a guest sees. So
// When who is visiting changes (sign-in, sign-up, Magic Link, password reset,
// Sign-out) the whole cache goes, and the router reloads what is on screen.
//
// `user` is the account the flow already holds (the response of the call that
// Signed them in); it is put back as the current user so the page does not ask
// For it again. Without it the next read asks the API, which answers a guest
// With no user. `router` is null for a flow that is itself a navigation (the
// Magic Link loader): the redirect it throws loads every route afresh.
export async function onIdentityChanged(
  queryClient: QueryClient,
  router: { invalidate: () => Promise<void> } | null,
  user?: User | null,
): Promise<void> {
  queryClient.clear();
  if (user) {
    queryClient.setQueryData(getCurrentUserQueryKey(), { user });
  }
  await router?.invalidate();
}

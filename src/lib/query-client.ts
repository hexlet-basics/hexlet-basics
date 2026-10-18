import type { DefaultOptions } from "@tanstack/react-query";

// Shared by the app's QueryClient and the test harness, so tests fetch the way
// The app does. With the library's staleTime of 0, data a route loader has just
// Read is already stale when the page's useQuery mounts on it, and the page
// Fetches it a second time. A minute keeps the loader's read for the render it
// Was made for (TanStack Query's SSR guidance); an invalidation still refetches
// Active queries at once, whatever their staleTime.
export const queryDefaults: DefaultOptions = {
  queries: { staleTime: 60_000 },
};

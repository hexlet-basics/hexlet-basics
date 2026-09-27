import type { DefaultOptions } from "@tanstack/react-query";

// Shared by the app's QueryClient and the test harness, so tests fetch the way
// the app does. With the library's staleTime of 0, data a route loader has just
// read is already stale when the page's useQuery mounts on it, and the page
// fetches it a second time. A minute keeps the loader's read for the render it
// was made for (TanStack Query's SSR guidance); an invalidation still refetches
// active queries at once, whatever their staleTime.
export const queryDefaults: DefaultOptions = {
  queries: { staleTime: 60_000 },
};

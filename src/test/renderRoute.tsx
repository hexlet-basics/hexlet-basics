import { QueryClient } from "@tanstack/react-query";
import {
  type AnyRoute,
  createMemoryHistory,
  createRootRouteWithContext,
  createRoute,
  createRouter,
  RouterProvider,
} from "@tanstack/react-router";
import type { ReactNode } from "react";
import { createI18n } from "@/lib/i18n";
import { renderWithProviders } from "./renderWithProviders";

// Mounts a real file route — its loader, its component, its staticData — at a
// real URL, under a synthetic root.
//
// The root is synthetic because the application's own `__root.tsx` renders
// `<html>`, `<head>` and Start's `<Scripts>`, which cannot be nested inside a
// test container. Everything below the root is the real thing: the route module
// under test is passed in, not reimplemented, so what a test drives is what a
// visitor's URL resolves to.
//
// One QueryClient serves both the router context and the provider, so a loader's
// prefetch is what the hooks under it read — the same arrangement
// setupRouterSsrQueryIntegration gives the app.
export async function renderRoute(
  route: AnyRoute,
  {
    path,
    initialPath,
    wrap = (element) => element,
  }: {
    // The route's own path pattern, e.g. "/{-$locale}/languages/$slug".
    path: string;
    // The URL to open, which must match that pattern.
    initialPath: string;
    // Wraps the router, for a page that needs a sized container to render into.
    wrap?: (element: ReactNode) => ReactNode;
  },
) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  const rootRoute = createRootRouteWithContext<{ queryClient: QueryClient }>()();
  // The file route carries its parent from the generated tree, so it is rebuilt
  // here from its own options with the synthetic root as its parent.
  const mounted = createRoute({
    ...(route.options as object),
    getParentRoute: () => rootRoute,
    path,
  } as never);

  const router = createRouter({
    routeTree: rootRoute.addChildren([mounted as never]),
    context: { queryClient, i18n: createI18n(), user: null } as never,
    history: createMemoryHistory({ initialEntries: [initialPath] }),
  });

  const screen = await renderWithProviders(
    wrap(<RouterProvider router={router as never} />),
    queryClient,
  );

  return { screen, router, queryClient };
}

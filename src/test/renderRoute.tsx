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
import { NotFoundPage, ServerErrorPage } from "@/components/ErrorPage";
import type { AuthUser } from "@/lib/auth";
import { createI18n, type Locale } from "@/lib/i18n";
import { localeFromPathname } from "@/lib/locale-path";
import { queryDefaults } from "@/lib/query-client";
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
    user = null,
    locale,
    wrap = (element) => element,
    also = [],
  }: {
    // The route's own path pattern, e.g. "/{-$locale}/languages/$slug".
    path: string;
    // The URL to open, which must match that pattern.
    initialPath: string;
    // Who is visiting. The application resolves this in the root route's
    // beforeLoad, which the synthetic root does not run, so a test that cares
    // whether the visitor is a guest says so here. The default is a guest.
    user?: AuthUser | null;
    // The language the page renders in. The application sets it in the
    // `{-$locale}` layout's beforeLoad, which the synthetic root does not run;
    // the default is the unprefixed en.
    locale?: Locale;
    // Wraps the router, for a page that needs a sized container to render into.
    wrap?: (element: ReactNode) => ReactNode;
    // Further real routes mounted beside the one under test, each at its own
    // path pattern, for a test that follows the page into another one.
    also?: { route: AnyRoute; path: string }[];
  },
) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { ...queryDefaults.queries, retry: false },
      mutations: { retry: false },
    },
  });

  const rootRoute = createRootRouteWithContext<{ queryClient: QueryClient }>()();
  // The file route carries its parent from the generated tree, so it is rebuilt
  // here from its own options with the synthetic root as its parent.
  const mount = (fileRoute: AnyRoute, pattern: string) =>
    createRoute({
      ...(fileRoute.options as object),
      getParentRoute: () => rootRoute,
      path: pattern,
    } as never);
  const mounted = [mount(route, path), ...also.map((extra) => mount(extra.route, extra.path))];

  // The locale layout switches the router's i18n to the URL's locale before a
  // route's head runs; the synthetic root has no such layout, so it is done here.
  const i18n = createI18n();
  await i18n.changeLanguage(locale ?? localeFromPathname(initialPath));

  const router = createRouter({
    routeTree: rootRoute.addChildren(mounted as never[]),
    context: { queryClient, i18n, user } as never,
    history: createMemoryHistory({ initialEntries: [initialPath] }),
    // The application's error pages, as getRouter sets them.
    defaultNotFoundComponent: NotFoundPage,
    defaultErrorComponent: ServerErrorPage,
  });

  const screen = await renderWithProviders(
    wrap(<RouterProvider router={router as never} />),
    queryClient,
    i18n,
  );

  return { screen, router, queryClient };
}

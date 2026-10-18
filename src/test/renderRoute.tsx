import { QueryClient } from "@tanstack/react-query";
import {
  RouterProvider,
  createMemoryHistory,
  createRootRouteWithContext,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import type { AnyRoute } from "@tanstack/react-router";
import type { ReactNode } from "react";
import { NotFoundPage, ServerErrorPage } from "@/components/ErrorPage";
import { Route as localeLayoutRoute } from "@/routes/{-$locale}";
import { Route as adminLayoutRoute } from "@/routes/{-$locale}/admin";
import type { AuthUser } from "@/lib/auth";
import { createI18n } from "@/lib/i18n";
import type { Locale } from "@/lib/i18n";
import { localeOf } from "@/lib/locale-path";
import { queryDefaults } from "@/lib/query-client";
import { renderWithProviders } from "./renderWithProviders";

// The application's file routes carry their parents from the generated tree, so
// each is rebuilt here from its own options under the parent the test gives it.
const make = (fileRoute: AnyRoute, parent: AnyRoute, pattern: string) =>
  createRoute({
    ...(fileRoute.options as object),
    getParentRoute: () => parent,
    path: pattern,
  } as never) as AnyRoute;

const under = (pattern: string, prefix: string) =>
  pattern === prefix || pattern.startsWith(`${prefix}/`);

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
  // The application's own layouts, outermost first: the `{-$locale}` layout (site
  // chrome, locale check) and, under it, the admin layout (its guard). A page is
  // mounted as a child of the deepest layout its path is under, so the chrome and
  // guards run as they do for a visitor; any other page hangs off the root. The
  // file routes carry their parents from the generated tree, so each is rebuilt
  // here from its own options.
  const LOCALE = "/{-$locale}";
  const ADMIN = `${LOCALE}/admin`;
  const layout = make(localeLayoutRoute, rootRoute, LOCALE);
  const adminLayout = make(adminLayoutRoute, layout, ADMIN.slice(LOCALE.length));
  const top: AnyRoute[] = [];
  const inLocale: AnyRoute[] = [adminLayout];
  const inAdmin: AnyRoute[] = [];
  const mount = (fileRoute: AnyRoute, pattern: string) => {
    if (under(pattern, ADMIN)) {
      inAdmin.push(make(fileRoute, adminLayout, pattern.slice(ADMIN.length) || "/"));
    } else if (under(pattern, LOCALE)) {
      inLocale.push(make(fileRoute, layout, pattern.slice(LOCALE.length) || "/"));
    } else {
      top.push(make(fileRoute, rootRoute, pattern));
    }
  };
  mount(route, path);
  for (const extra of also) {
    mount(extra.route, extra.path);
  }
  adminLayout.addChildren(inAdmin as never[]);
  layout.addChildren(inLocale as never[]);

  // The locale layout switches the router's i18n to the URL's locale before a
  // route's head runs; the synthetic root has no such layout, so it is done here.
  const i18n = createI18n();
  await i18n.changeLanguage(locale ?? localeOf(initialPath));

  const router = createRouter({
    routeTree: rootRoute.addChildren([...top, layout] as never[]),
    context: { queryClient, i18n, user } as never,
    history: createMemoryHistory({ initialEntries: [initialPath] }),
    // Links preload on hover, as in getRouter.
    defaultPreload: "intent",
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

import {
  RouterProvider,
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import type { ReactNode } from "react";
import { vi } from "vitest";
import { renderWithProviders } from "./renderWithProviders";

interface RenderWithRouterOptions {
  // Memory-history entry the router starts on. Default carries a `locale` param
  // So components calling useParams({ strict: false }) see the app's shape.
  initialPath?: string;
  // Route pattern the component mounts under. Must line up with initialPath.
  path?: string;
}

// Mounts `ui` inside a minimal in-memory TanStack router so components that call
// UseNavigate/useParams/Link render without the real route tree. Navigation is
// NOT performed — `router.navigate` is a spy the test asserts against.
export async function renderWithRouter(
  ui: ReactNode,
  {
    initialPath = "/admin",
    path = "/admin",
    navigate: navigation = "spy",
  }: RenderWithRouterOptions = {},
) {
  const rootRoute = createRootRoute();
  const route = createRoute({
    getParentRoute: () => rootRoute,
    path,
    component: () => ui,
  });
  const router = createRouter({
    routeTree: rootRoute.addChildren([route]),
    history: createMemoryHistory({ initialEntries: [initialPath] }),
  });

  // Router.navigate is generic over the route tree, which a Mock cannot carry,
  // So the mock takes any arguments and is assigned onto the router as it is.
  const navigate = vi.fn<(...args: unknown[]) => Promise<void>>().mockResolvedValue(undefined);
  router.navigate = navigate;

  const screen = await renderWithProviders(<RouterProvider router={router} />);
  return { screen, navigate };
}

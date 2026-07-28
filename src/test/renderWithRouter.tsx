import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from "@tanstack/react-router";
import type { ReactNode } from "react";
import { type Mock, vi } from "vitest";
import { renderWithProviders } from "./renderWithProviders";

type RenderWithRouterOptions = {
  // Memory-history entry the router starts on. Default carries a `locale` param
  // so components calling useParams({ strict: false }) see the app's shape.
  initialPath?: string;
  // Route pattern the component mounts under. Must line up with initialPath.
  path?: string;
};

// Mounts `ui` inside a minimal in-memory TanStack router so components that call
// useNavigate/useParams/Link render without the real route tree. Navigation is
// NOT performed — `router.navigate` is a spy the test asserts against.
export async function renderWithRouter(
  ui: ReactNode,
  { initialPath = "/admin", path = "/admin" }: RenderWithRouterOptions = {},
) {
  const rootRoute = createRootRoute();
  const route = createRoute({
    getParentRoute: () => rootRoute,
    path,
    component: () => <>{ui}</>,
  });
  const router = createRouter({
    routeTree: rootRoute.addChildren([route]),
    history: createMemoryHistory({ initialEntries: [initialPath] }),
  });

  const navigate: Mock = vi.fn().mockResolvedValue(undefined);
  router.navigate = navigate as unknown as typeof router.navigate;

  const screen = await renderWithProviders(<RouterProvider router={router} />);
  return { screen, navigate };
}

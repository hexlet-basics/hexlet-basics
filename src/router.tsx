import {
  createRootRoute,
  createRoute,
  createRouter,
  Outlet,
} from "@tanstack/react-router";
import CatalogIndex from "@/pages/languages/index";
import CourseShow from "@/pages/languages/show";

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

// Catalog is served at both `/` and `/languages` (legacy URL kept for compat).
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: CatalogIndex,
});

const languagesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/languages",
  component: CatalogIndex,
});

const courseRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/languages/$slug",
  component: CourseShow,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  languagesRoute,
  courseRoute,
]);

export const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

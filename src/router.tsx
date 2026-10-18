import { QueryClient } from "@tanstack/react-query";
import { createRouter } from "@tanstack/react-router";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import { createI18n } from "@/lib/i18n";
import { queryDefaults } from "@/lib/query-client";
// Side-effect import: configures the hey-api client (baseUrl, credentials,
// SSR cookie forwarding) once at module scope.
import "@/lib/api-client";
import { NotFoundPage, ServerErrorPage } from "@/components/ErrorPage";
import { routeTree } from "./routeTree.gen";

// Start calls getRouter() per request, so each SSR render gets a fresh
// QueryClient. setupRouterSsrQueryIntegration wires dehydration/hydration and
// Wraps the tree in QueryClientProvider.
export function getRouter() {
  const queryClient = new QueryClient({ defaultOptions: queryDefaults });
  // Request-scoped i18n instance; the `{-$locale}` layout resolves the URL
  // Locale in beforeLoad and calls changeLanguage on it.
  const i18n = createI18n();

  const router = createRouter({
    routeTree,
    // `user` is resolved server-side in the root beforeLoad; null is the base.
    context: { queryClient, i18n, user: null },
    defaultPreload: "intent",
    // Legacy errors#show for every route: a match that finds nothing renders
    // The 404 page and a load that throws the 500 page, in place of the route
    // And so inside the site layout. The server answers with the same status.
    defaultNotFoundComponent: NotFoundPage,
    defaultErrorComponent: ServerErrorPage,
    scrollRestoration: true,
    defaultViewTransition: true,
    // On hydration the client does not rerun beforeLoad — it adopts the
    // Server's results — so the `{-$locale}` layout never calls changeLanguage
    // On the browser's own i18n instance, and every page repaints in the `en`
    // Default. The resolved language travels with the dehydrated router and
    // Is applied before the first client render.
    dehydrate: () => ({ locale: i18n.language }),
    hydrate: async ({ locale }) => {
      await i18n.changeLanguage(locale);
    },
  });

  setupRouterSsrQueryIntegration({ router, queryClient });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}

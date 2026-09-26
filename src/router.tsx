import { QueryClient } from "@tanstack/react-query";
import { createRouter } from "@tanstack/react-router";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import { createI18n } from "@/lib/i18n";
// Side-effect import: configures the hey-api client (baseUrl, credentials,
// SSR cookie forwarding) once at module scope.
import "@/lib/api-client";
import { routeTree } from "./routeTree.gen";

// Start calls getRouter() per request, so each SSR render gets a fresh
// QueryClient. setupRouterSsrQueryIntegration wires dehydration/hydration and
// wraps the tree in QueryClientProvider.
export function getRouter() {
  const queryClient = new QueryClient();
  // Request-scoped i18n instance; the `{-$locale}` layout resolves the URL
  // locale in beforeLoad and calls changeLanguage on it.
  const i18n = createI18n();

  const router = createRouter({
    routeTree,
    // `user` is resolved server-side in the root beforeLoad; null is the base.
    context: { queryClient, i18n, user: null },
    defaultPreload: "intent",
    scrollRestoration: true,
    // On hydration the client does not rerun beforeLoad — it adopts the
    // server's results — so the `{-$locale}` layout never calls changeLanguage
    // on the browser's own i18n instance, and every page repaints in the `en`
    // default. The resolved language travels with the dehydrated router and
    // is applied before the first client render.
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

import { QueryClient } from "@tanstack/react-query";
import { createRouter } from "@tanstack/react-router";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import { client } from "@/client/client.gen";
import { routeTree } from "./routeTree.gen";

// The generated hey-api client is a singleton; point it at the Go API. During
// SSR the Node process reaches Go over the internal network (API_URL); in the
// browser it uses the public origin (VITE_API_URL). The dead branch is stripped
// per build target, so `process` never leaks into the client bundle.
client.setConfig({
  baseUrl: import.meta.env.SSR
    ? (process.env.API_URL ?? "http://localhost:3001")
    : (import.meta.env.VITE_API_URL ?? "http://localhost:3001"),
});

// Start calls getRouter() per request, so each SSR render gets a fresh
// QueryClient. setupRouterSsrQueryIntegration wires dehydration/hydration and
// wraps the tree in QueryClientProvider.
export function getRouter() {
  const queryClient = new QueryClient();

  const router = createRouter({
    routeTree,
    context: { queryClient },
    defaultPreload: "intent",
    scrollRestoration: true,
  });

  setupRouterSsrQueryIntegration({ router, queryClient });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}

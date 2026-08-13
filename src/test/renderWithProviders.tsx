import { MantineProvider } from "@mantine/core";
import { ModalsProvider } from "@mantine/modals";
import { Notifications } from "@mantine/notifications";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { I18nextProvider } from "react-i18next";
import { render } from "vitest-browser-react";
import { createI18n } from "@/lib/i18n";

// Mirrors the provider stack from __root.tsx so components render the same way
// they do in the app: Mantine + Modals + Notifications, a request-scoped i18n
// instance (so `t(($) => $.x)` resolves), and a fresh QueryClient per render with
// retries off so failed requests surface immediately.
//
// Components that use the router (useNavigate/useParams/Link) need
// renderWithRouter instead.
function makeQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
}

// `queryClient` is passed in by renderRoute, so a route loader and the hooks
// under it share one cache, as they do in the app.
export function renderWithProviders(ui: ReactNode, queryClient = makeQueryClient()) {
  const i18n = createI18n();

  function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <I18nextProvider i18n={i18n}>
          <MantineProvider>
            <ModalsProvider>
              <Notifications position="top-right" />
              {children}
            </ModalsProvider>
          </MantineProvider>
        </I18nextProvider>
      </QueryClientProvider>
    );
  }

  return render(ui, { wrapper: Wrapper });
}

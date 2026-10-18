import { MantineProvider } from "@mantine/core";
import { ModalsProvider } from "@mantine/modals";
import { Notifications } from "@mantine/notifications";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { CookiesProvider } from "react-cookie";
import { I18nextProvider } from "react-i18next";
import { render } from "vitest-browser-react";
import { createI18n } from "@/lib/i18n";
import { queryDefaults } from "@/lib/query-client";

// Mirrors the provider stack from __root.tsx so components render the same way
// They do in the app: Mantine + Modals + Notifications, a request-scoped i18n
// Instance (so `t(($) => $.x)` resolves), and a fresh QueryClient per render with
// Retries off so failed requests surface immediately.
//
// Components that use the router (useNavigate/useParams/Link) need
// RenderWithRouter instead.
function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { ...queryDefaults.queries, retry: false },
      mutations: { retry: false },
    },
  });
}

// `queryClient` and `i18n` are passed in by renderRoute, so a route loader and
// The hooks under it share one cache and one language, as they do in the app.
export function renderWithProviders(
  ui: ReactNode,
  queryClient = makeQueryClient(),
  i18n = createI18n(),
) {
  function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <I18nextProvider i18n={i18n}>
          <CookiesProvider>
            {/* Env="test": no transitions and portals rendered in place, so nothing
              of a test (a modal still closing when it ends) outlives its tree. */}
            <MantineProvider env="test">
              <ModalsProvider>
                <Notifications position="top-right" />
                {children}
              </ModalsProvider>
            </MantineProvider>
          </CookiesProvider>
        </I18nextProvider>
      </QueryClientProvider>
    );
  }

  return render(ui, { wrapper: Wrapper });
}

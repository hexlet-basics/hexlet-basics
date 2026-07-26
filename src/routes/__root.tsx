import "@mantine/core/styles.css";
import {
  ColorSchemeScript,
  MantineProvider,
  mantineHtmlProps,
} from "@mantine/core";
import type { QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts,
} from "@tanstack/react-router";
import type { i18n as I18n } from "i18next";
import type { ReactNode } from "react";
import { I18nextProvider } from "react-i18next";
import { type AuthUser, getCurrentUser } from "@/lib/auth";

interface RouterContext {
  queryClient: QueryClient;
  i18n: I18n;
  user: AuthUser | null;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  // Resolve the current user from the httpOnly JWT cookie server-side, so
  // authenticated pages render on the server (ADR-0008). STUB: null until the
  // auth backend lands.
  beforeLoad: () => ({ user: getCurrentUser() }),
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "Hexlet Basics" },
    ],
  }),
  component: RootComponent,
});

function RootComponent() {
  // The `{-$locale}` layout has already run changeLanguage in beforeLoad, so
  // i18n.language is the resolved locale by the time the document renders.
  const { i18n } = Route.useRouteContext();

  return (
    <I18nextProvider i18n={i18n}>
      <RootDocument lang={i18n.language}>
        <Outlet />
      </RootDocument>
    </I18nextProvider>
  );
}

function RootDocument({
  children,
  lang,
}: {
  children: ReactNode;
  lang: string;
}) {
  return (
    <html lang={lang} {...mantineHtmlProps}>
      <head>
        <ColorSchemeScript />
        <HeadContent />
      </head>
      <body>
        <MantineProvider>{children}</MantineProvider>
        <Scripts />
      </body>
    </html>
  );
}

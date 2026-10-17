import "@mantine/core/styles.css";
import "@mantine/code-highlight/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/tiptap/styles.css";
import "@mantine/notifications/styles.css";
import "@/view-transitions.css";
import { ColorSchemeScript, MantineProvider, mantineHtmlProps } from "@mantine/core";
import { ModalsProvider } from "@mantine/modals";
import { Notifications } from "@mantine/notifications";
import type { QueryClient } from "@tanstack/react-query";
import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";
import { createRootRouteWithContext, HeadContent, Outlet, Scripts } from "@tanstack/react-router";
import type { i18n as I18n } from "i18next";
import { type ReactNode, useEffect } from "react";
import { Cookies, CookiesProvider } from "react-cookie";
import { I18nextProvider } from "react-i18next";
import { NotFoundPage } from "@/components/ErrorPage";
import ApplicationLayout from "@/components/layout/ApplicationLayout";
import { type AuthUser, resolveCurrentUser } from "@/lib/auth";
import { recordFirstVisit } from "@/lib/first-visit";
import { localeFromPathname } from "@/lib/locale-path";
import { SITE_NAME } from "@/lib/seo-head";
import { theme } from "@/theme";

interface RouterContext {
  queryClient: QueryClient;
  i18n: I18n;
  user: AuthUser | null;
}

// The incoming request's Cookie header, for the cookie provider below. Server
// only: in the browser the provider reads document.cookie itself.
const getRequestCookies = createIsomorphicFn()
  .server(() => getRequestHeader("cookie"))
  .client(() => undefined);

export const Route = createRootRouteWithContext<RouterContext>()({
  // Resolve the current user from the httpOnly session cookie server-side, so
  // authenticated pages (and the admin guard) render on the server (ADR-0008).
  // Cached in the shared QueryClient, so this is one fetch per request.
  //
  // The language is also set here from the URL prefix: a miss under `/ru/...`
  // is rendered without the `{-$locale}` layout's beforeLoad having run, and
  // its 404 would otherwise come out in the `en` default.
  beforeLoad: async ({ context, location }) => {
    await context.i18n.changeLanguage(localeFromPathname(location.pathname));
    return {
      user: await resolveCurrentUser(context.queryClient),
      // A string, so it travels with the dehydrated router and the browser's
      // first render reads the same cookies the server rendered with.
      cookieHeader: getRequestCookies(),
    };
  },
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: SITE_NAME },
    ],
  }),
  // A path the locale layout refuses (`/unknownpage`, `/xx/...`, `/en/...`)
  // never mounts that layout, so its 404 brings the site chrome itself. Every
  // other miss renders the router's default 404 inside the layout.
  notFoundComponent: () => (
    <ApplicationLayout>
      <NotFoundPage />
    </ApplicationLayout>
  ),
  component: RootComponent,
});

function RootComponent() {
  // The `{-$locale}` layout has already run changeLanguage in beforeLoad, so
  // i18n.language is the resolved locale by the time the document renders.
  const { i18n, cookieHeader } = Route.useRouteContext();

  // Effects run only in the browser, so the first visit is recorded from the
  // page the visitor actually landed on — once, whichever page that is.
  useEffect(() => {
    recordFirstVisit();
  }, []);

  return (
    <I18nextProvider i18n={i18n}>
      <CookiesProvider cookies={new Cookies(cookieHeader)}>
        <RootDocument lang={i18n.language}>
          <Outlet />
        </RootDocument>
      </CookiesProvider>
    </I18nextProvider>
  );
}

function RootDocument({ children, lang }: { children: ReactNode; lang: string }) {
  return (
    <html lang={lang} {...mantineHtmlProps}>
      <head>
        <ColorSchemeScript />
        <HeadContent />
      </head>
      <body>
        <MantineProvider theme={theme}>
          <Notifications />
          <ModalsProvider>{children}</ModalsProvider>
        </MantineProvider>
        <Scripts />
      </body>
    </html>
  );
}

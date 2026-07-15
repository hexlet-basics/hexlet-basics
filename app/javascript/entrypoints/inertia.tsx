import { createInertiaApp } from "@inertiajs/react";
import * as Sentry from "@sentry/react";
import { configureRoutes } from "@/init.ts";
import RootLayout from "@/layouts/RootLayout.tsx";
import I18nAppProvider from "@/lib/I18nAppProvider.tsx";

// Browser error monitoring. Runs before the app boots so RootLayout's
// Sentry.ErrorBoundary/withProfiler and the global window handlers are live.
// dsn is only set on production builds (VITE_SENTRY_DSN build secret); without
// it Sentry.init is a no-op, so dev/SSR are unaffected. Source maps are uploaded
// by sentryVitePlugin (vite.config.ts) with debug ids under the same release.
if (!import.meta.env.SSR && import.meta.env.VITE_SENTRY_DSN) {
  Sentry.init({
    dsn: import.meta.env.VITE_SENTRY_DSN,
    release: import.meta.env.VITE_RELEASE_VERSION,
    environment: import.meta.env.VITE_NODE_ENV,
  });
}

createInertiaApp({
  defaults: {
    visitOptions: (href, options) => {
      return { viewTransition: true };
    },
  },
  layout: (_component, page) => {
    configureRoutes(page.props.suffix);
    return RootLayout;
  },
  pages: "../pages",
  withApp(app, { ssr }) {
    return <I18nAppProvider app={app} ssr={ssr} />;
  },
});

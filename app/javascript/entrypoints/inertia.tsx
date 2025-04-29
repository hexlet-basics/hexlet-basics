import Root from "@/components/Root.tsx";
import type { ResolvedComponent, RootProps } from "@/types";
import { createInertiaApp } from "@inertiajs/react";
import * as Sentry from "@sentry/react";
import { createRoot, hydrateRoot } from "react-dom/client";
import "@/init.ts";
import configure from "@/configure";

if (import.meta.env.DEV) {
  localStorage.debug = "app:*";
}

Sentry.init({
  debug: import.meta.env.DEV,
  dsn: import.meta.env.VITE_SENTRY_DSN,
  ignoreErrors: [
    "Failed to fetch dynamically imported module",
    "Unable to preload CSS",
  ],
});

createInertiaApp({
  // Set default page title
  // see https://inertia-rails.netlify.app/guide/title-and-meta
  //
  // title: title => title ? `${title} - App` : 'App',

  // Disable progress bar
  //
  // see https://inertia-rails.netlify.app/guide/progress-indicators
  // progress: false,

  resolve: async (name) => {
    const pages = import.meta.glob<ResolvedComponent>("../pages/**/*.tsx");
    const pageFn = pages[`../pages/${name}.tsx`];
    if (!pageFn) {
      console.error(`Missing Inertia page component: '${name}.tsx'`);
    }
    const page = await pageFn();
    page.default.layout ??= (page) => <Root>{page}</Root>;

    return page;
  },

  setup({ el, App, props }) {
    if (el) {
      const typedProps = props as RootProps;
      const { locale, suffix } = typedProps.initialPage.props;
      configure(locale, suffix)
      const vdomFn = () => {
        // return <Root locale={locale} suffix={suffix}><App {...props} /></Root>;
        return <App {...props} />
      };
      if (import.meta.env.MODE !== "development") {
        hydrateRoot(el, vdomFn(), {
          // Callback called when an error is thrown and not caught by an ErrorBoundary.
          onUncaughtError: Sentry.reactErrorHandler((error, errorInfo) => {
            console.warn("Uncaught error", error, errorInfo.componentStack);
          }),
          // Callback called when React catches an error in an ErrorBoundary.
          onCaughtError: Sentry.reactErrorHandler(),
          // Callback called when React automatically recovers from errors.
          onRecoverableError: Sentry.reactErrorHandler(),
        });
      } else {
        const root = createRoot(el);
        root.render(vdomFn());
      }
    } else {
      console.error(
        "Missing root element.\n\n" +
        "If you see this error, it probably means you load Inertia.js on non-Inertia pages.\n" +
        'Consider moving <%= vite_typescript_tag "inertia" %> to the Inertia-specific layout instead.',
      );
    }
  },
});

import Root from "@/components/Root.tsx";
import { createInertiaApp } from "@inertiajs/react";
import * as Sentry from "@sentry/react";
import type { ReactNode } from "react";
import { createRoot, hydrateRoot } from "react-dom/client";

window.addEventListener("vite:preloadError", (event) => {
  window.location.reload(); // for example, refresh the page
});

if (import.meta.env.DEV) {
  localStorage.debug = "app:*";
}

// Temporary type definition, until @inertiajs/react provides one
type ResolvedComponent = {
  default: ReactNode;
  layout?: (page: ReactNode) => ReactNode;
};

Sentry.init({
  dsn: import.meta.env.VITE_SENTRY_DSN,
  ignoreErrors: [],
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

  resolve: (name) => {
    const pages = import.meta.glob<ResolvedComponent>("../pages/**/*.tsx");
    const pageFn = pages[`../pages/${name}.tsx`];
    if (!pageFn) {
      console.error(`Missing Inertia page component: '${name}.tsx'`);
    }
    const page = pageFn();

    // To use a default layout, import the Layout component
    // and use the following line.
    // see https://inertia-rails.netlify.app/guide/pages#default-layouts
    //
    // page.default.layout ||= (page) => createElement(Layout, null, page)

    return page;
  },

  setup({ el, App, props }) {
    if (el) {
      const vdomFn = () => {
        return (
          <Root {...props}>
            <App {...props} />
          </Root>
        );
      };
      if (import.meta.env.MODE === "production") {
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

import Root from "@/components/Root.tsx";
import type { ResolvedComponent } from "@/types";
import { createInertiaApp } from "@inertiajs/react";
import createServer from "@inertiajs/react/server";
import * as Sentry from "@sentry/node";
import ReactDOMServer from "react-dom/server";

Sentry.init({
  dsn: import.meta.env.VITE_SENTRY_DSN,
});

createServer((page) =>
  createInertiaApp({
    page,
    render: (...args) => {
      try {
        return ReactDOMServer.renderToString(...args);
      } catch (error) {
        Sentry.setContext("page", {
          url: page.url,
          component: page.component,
        });

        Sentry.captureException(error);

        throw error;
      }
    },
    resolve: (name) => {
      // const pages = import.meta.glob("../pages/**/*.jsx", { eager: true });
      const pages = import.meta.glob<ResolvedComponent>("../pages/**/*.tsx", {
        eager: true,
      });
      const page = pages[`../pages/${name}.tsx`];

      page.default.layout ??= (page) => <Root>{page}</Root>;

      return page;
    },
    setup: ({ App, props }) => {
      const vdom = <App {...props} />;

      return vdom;
    },
  }),
);

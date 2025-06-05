import Root from "@/components/Root.tsx";
import type { ResolvedComponent, RootProps } from "@/types";
import { createInertiaApp } from "@inertiajs/react";
import createServer from "@inertiajs/react/server";
import * as Sentry from "@sentry/node";
import ReactDOMServer from "react-dom/server";
import "@/init.ts";
import configure from "@/lib/configure";

Sentry.init({
  // debug: import.meta.env.DEV,
  dsn: import.meta.env.VITE_SENTRY_DSN,
});

createServer((page) => {
  // NOTE: используется для просмотра последних попыток рендера перед падением контейнера ssr на проде (дебаг)
  console.log(page.url);
  const used = process.memoryUsage().rss / 1024 / 1024;
  console.log(`Memory usage: ${Math.round(used)} MiB`);

  return createInertiaApp({
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
      const typedProps = props as RootProps;
      const { locale, suffix } = typedProps.initialPage.props;

      configure(locale, suffix)
      const vdom = <App {...props} />;
      return vdom;
    },
  });
}, { cluster: true });

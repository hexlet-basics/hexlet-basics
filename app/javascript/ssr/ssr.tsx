import Root from "@/components/Root.tsx";
import * as Routes from "@/routes.js";
import type { ResolvedComponent, RootProps } from "@/types";
import { createInertiaApp } from "@inertiajs/react";
import createServer from "@inertiajs/react/server";
import * as Sentry from "@sentry/node";
import dayjs from "dayjs";
import i18next from "i18next";
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
      const typedProps = props as RootProps;
      const { locale, suffix } = typedProps.initialPage.props;

      i18next.changeLanguage(locale);
      dayjs.locale(locale);
      Routes.configure({ default_url_options: { suffix } });

      const vdom = <App {...props} />;

      return vdom;
    },
  }),
);

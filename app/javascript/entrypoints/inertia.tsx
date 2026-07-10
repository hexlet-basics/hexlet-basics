import { createInertiaApp } from "@inertiajs/react";
import { configureRoutes } from "@/init.ts";
import RootLayout from "@/layouts/RootLayout.tsx";
import I18nAppProvider from "@/lib/I18nAppProvider.tsx";

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

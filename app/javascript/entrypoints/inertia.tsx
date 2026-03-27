import { createInertiaApp } from "@inertiajs/react";
import RootLayout from "@/layouts/RootLayout.tsx";
import I18nAppProvider from "@/lib/I18nAppProvider.tsx";

import "@/init.ts";

createInertiaApp({
  defaults: {
    visitOptions: (href, options) => {
      return { viewTransition: true };
    },
  },
  layout: (_component, page) => {
    return RootLayout;
  },
  pages: "../pages",
  withApp(app, { ssr }) {
    return <I18nAppProvider app={app} ssr={ssr} />;
  },
});

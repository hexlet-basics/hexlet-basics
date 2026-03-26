import { createInertiaApp } from "@inertiajs/react";
import configure from "@/configure";
import RootLayout from "@/layouts/RootLayout.tsx";

// import "@/init.ts";

createInertiaApp({
  defaults: {
    visitOptions: (href, options) => {
      return { viewTransition: true };
    },
  },
  layout: (_component, page) => {
    configure(page.props.locale, page.props.suffix);

    return RootLayout;
  },
  pages: "../pages",
});

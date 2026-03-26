import { createInertiaApp } from "@inertiajs/react";
import dayjs from "dayjs";
import i18next from "i18next";
import configure from "@/configure";
import RootLayout from "@/layouts/RootLayout.tsx";
import * as Routes from "@/routes.js";

configure();

createInertiaApp({
  defaults: {
    visitOptions: (href, options) => {
      return { viewTransition: true };
    },
  },
  layout: (_component, page) => {
    i18next.changeLanguage(page.props.locale);
    Routes.configure({
      default_url_options: {
        suffix: page.props.suffix,
        protocol: "https",
        host: import.meta.env.VITE_APP_HOST,
      },
    });
    dayjs.locale(page.props.locale);

    return RootLayout;
  },
  pages: "../pages",
});

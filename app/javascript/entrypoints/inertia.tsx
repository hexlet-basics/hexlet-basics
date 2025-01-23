import * as Sentry from "@sentry/react";
import { createInertiaApp } from "@inertiajs/react";
import { PrimeReactProvider } from "primereact/api";
import type { ReactNode } from "react";
import { createRoot, hydrateRoot } from "react-dom/client";

import * as Routes from "@/routes.js";
import pt from "../primereact.ts";

import "react-bootstrap";

import type { SharedProps } from "@/types";
import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import locales from "../locales.json";

export const resources = locales;
export const defaultNS = "web";

i18next.use(initReactI18next);
i18next.init({
  resources,
  defaultNS,
  ns: Object.keys(resources.ru),
  // lng: locale,
  interpolation: {
    prefix: "%{",
    suffix: "}",
    escapeValue: false, // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
  },
});

// Temporary type definition, until @inertiajs/react provides one
type ResolvedComponent = {
  default: ReactNode;
  layout?: (page: ReactNode) => ReactNode;
};

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
    const RootComponent = () => {
      // const { locale, suffix } = usePage<SharedProps>().props;
      const { locale, suffix } = props.initialPage
        .props as unknown as SharedProps;

      // useEffect?
      i18next.changeLanguage(locale);
      Routes.configure({ default_url_options: { suffix } });

      return (
        <PrimeReactProvider value={{ pt, unstyled: true }}>
          <App {...props} />
        </PrimeReactProvider>
      );
    };
    if (el) {
      if (import.meta.env.MODE === "production") {
        hydrateRoot(el, <RootComponent />, {
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
        root.render(<RootComponent />);
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

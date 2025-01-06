import { createInertiaApp } from "@inertiajs/react";
import createServer from "@inertiajs/react/server";
import React from "react";
import ReactDOMServer from "react-dom/server";

import "react-bootstrap";
import "@/i18n";

createServer((page) =>
  createInertiaApp({
    page,
    render: ReactDOMServer.renderToString,
    resolve: (name) => {
      // const pages = import.meta.glob("../pages/**/*.jsx", { eager: true });
      const pages = import.meta.glob("../pages/**/*.tsx", {
        eager: true,
      });
      return pages[`../pages/${name}.tsx`];
    },
    setup: ({ App, props }) => <App {...props} />,
  }),
);

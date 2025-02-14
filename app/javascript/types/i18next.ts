// import the original type declarations
import "i18next";
import type locales from "../locales.json";

const defaultNS = "web" as const;

declare module "i18next" {
  interface i18n {
    language: "en" | "ru";
  }
  // Extend CustomTypeOptions
  interface CustomTypeOptions {
    // custom namespace type, if you changed it
    defaultNS: typeof defaultNS;
    resources: (typeof locales)["ru"];
  }
}

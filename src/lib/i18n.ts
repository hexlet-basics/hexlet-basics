import type { i18n as I18n } from "i18next";
import { createInstance } from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/locales/en/translation";
import es from "@/locales/es/translation";
import ru from "@/locales/ru/translation";

export type Locale = "en" | "ru" | "es";

// `en` is the default and is served unprefixed; `ru`/`es` live behind a path
// prefix (legacy `scope "(:suffix)", suffix: /es|ru/`, ADR-0002 / ADR-0008).
export const DEFAULT_LOCALE: Locale = "en";
// Locales that appear as a URL prefix (everything except the unprefixed default).
export const PREFIXED_LOCALES = ["ru", "es"] as const;

export function isLocale(value: string): value is Locale {
  return value === "en" || value === "ru" || value === "es";
}

// One i18next instance per request (created in getRouter), never a module
// singleton: on the SSR server a shared instance would let concurrent `/ru` and
// `/es` renders clobber each other's language. Resources are bundled and init is
// synchronous, so strings are available on first server paint.
export function createI18n(): I18n {
  const instance = createInstance();

  instance.use(initReactI18next).init({
    lng: DEFAULT_LOCALE,
    fallbackLng: DEFAULT_LOCALE,
    defaultNS: "translation",
    ns: ["translation"],
    resources: { en, es, ru },
    interpolation: { escapeValue: false },
    react: { useSuspense: false },
  });

  return instance;
}

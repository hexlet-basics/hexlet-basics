import { createInstance, type i18n as I18n, type Resource } from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/locales/en/translation.ts";
import es from "@/locales/es/translation.ts";
import ru from "@/locales/ru/translation.ts";
import type { Locale } from "@/types";

const resources = {
  en,
  es,
  ru,
} as const satisfies Resource;

let clientI18n: I18n | null = null;

function initI18n(instance: I18n, locale: Locale, resourceStore: Resource) {
  instance.use(initReactI18next);
  void instance.init({
    lng: locale,
    fallbackLng: "ru",
    defaultNS: "translation",
    ns: ["translation"],
    debug: import.meta.env.DEV,
    interpolation: {
      escapeValue: false,
    },
    react: {
      useSuspense: false,
    },
    resources: resourceStore,
  });

  return instance;
}

export function buildInitialI18nStore(locale: Locale): Resource {
  return {
    [locale]: resources[locale],
  };
}

export function createServerI18n(locale: Locale) {
  return initI18n(createInstance(), locale, buildInitialI18nStore(locale));
}

export function getClientI18n(locale: Locale) {
  if (clientI18n) {
    return clientI18n;
  }

  clientI18n = initI18n(
    createInstance(),
    locale,
    buildInitialI18nStore(locale),
  );

  return clientI18n;
}

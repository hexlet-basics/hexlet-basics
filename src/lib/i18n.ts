import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/locales/en/translation";
import es from "@/locales/es/translation";
import ru from "@/locales/ru/translation";

export type Locale = "ru" | "en" | "es";

// Locale resources are copied verbatim from the legacy app
// (legacy/app/javascript/locales), generated there by i18next-cli.
i18n.use(initReactI18next).init({
  lng: "ru",
  fallbackLng: "ru",
  defaultNS: "translation",
  ns: ["translation"],
  resources: { en, es, ru },
  interpolation: { escapeValue: false },
  react: { useSuspense: false },
});

export default i18n;

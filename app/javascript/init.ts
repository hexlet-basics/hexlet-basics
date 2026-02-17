import "dayjs/locale/ru";

import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import { gon } from "./lib/gon.ts";
import en from "./locales/en/translation.ts";
// import es from './locales/es/translation.ts';
import ru from "./locales/ru/translation.ts";

i18next.use(initReactI18next).init({
  lng: gon.locale,
  debug: import.meta.env.DEV,
  defaultNS: "translation",
  ns: ["translation"],
  interpolation: {
    prefix: "%{",
    suffix: "}",
    escapeValue: false,
  },
  resources: {
    ru,
    en,
    // es,
  },
});

export default i18next;

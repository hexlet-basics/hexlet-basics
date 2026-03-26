import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import * as Routes from "@/routes.js";
import type { Locale } from "@/types";
import en from "./locales/en/translation.ts";
// import es from './locales/es/translation.ts';
import ru from "./locales/ru/translation.ts";

export default function configure(locale: Locale, suffix: string | null) {
  i18next.use(initReactI18next).init({
    lng: locale,
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
  // if (i18next.language !== locale) {
  //   i18next.changeLanguage(locale);
  // }

  dayjs.extend(localizedFormat);
  dayjs.locale(locale);

  dayjs.extend(relativeTime);
  dayjs.extend(duration);

  Routes.configure({
    default_url_options: {
      suffix,
      protocol: "https",
      host: import.meta.env.VITE_APP_HOST,
    },
  });
}

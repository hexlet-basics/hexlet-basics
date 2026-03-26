import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import { gon } from "./lib/gon.ts";
import en from "./locales/en/translation.ts";
// import es from './locales/es/translation.ts';
import ru from "./locales/ru/translation.ts";

export default function configure() {
  i18next.use(initReactI18next).init({
    lng: gon.locale || "en",
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

  dayjs.extend(relativeTime);
  dayjs.extend(duration);
}

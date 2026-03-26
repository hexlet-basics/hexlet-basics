import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import i18next from "i18next";
import * as Routes from "@/routes.js";
import type { Locale } from "@/types";

export default function configure(locale: Locale, suffix: string | null) {
  if (i18next.language !== locale) {
    i18next.changeLanguage(locale);
  }

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

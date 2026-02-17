import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import i18next from "i18next";
import type { Locale } from "@/types";
// import { configure as configureRoutes } from "../routes.js";

export default function configure(locale: Locale, host?: string) {
  if (i18next.language !== locale) {
    i18next.changeLanguage(locale);
  }

  dayjs.extend(localizedFormat);
  dayjs.locale(locale);

  dayjs.extend(relativeTime);
  dayjs.extend(duration);

  // if (host) {
  //   configureRoutes({
  //     default_url_options: {
  //       host: host,
  //       protocol: "https",
  //     },
  //   });
  // }
}

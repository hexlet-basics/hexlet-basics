import * as Routes from "@/routes.js";
import i18next from "i18next";
import { Locale } from "../types";
// import { gon } from "@/lib/gon";
import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import relativeTime from "dayjs/plugin/relativeTime";


export default function configure(locale: Locale, suffix: string | null) {
  i18next.changeLanguage(locale);
  dayjs.locale(locale);

  dayjs.extend(relativeTime);
  dayjs.extend(duration);
  Routes.configure({ default_url_options: { suffix: suffix, protocol: "https", host: import.meta.env.VITE_APP_HOST } });
}

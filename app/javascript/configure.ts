import * as Routes from "@/routes.js";
import { dayjs } from "@/lib/utils";
import i18next from "i18next";
import { Locale } from "./types";
// import { gon } from "@/lib/gon";

export default function configure(locale: Locale, suffix: string | null) {
  i18next.changeLanguage(locale);
  dayjs.locale(locale);
  Routes.configure({ default_url_options: { suffix: suffix, protocol: "https", host: import.meta.env.VITE_APP_HOST } });
}

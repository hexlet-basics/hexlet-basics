import { useTranslation } from "react-i18next";
import { dateTime } from "@/lib/time-ago";

// Formats an ISO timestamp for a table cell in the app's selected language (the
// One `t` translates in), not the browser's, so the back office reads in one
// Language throughout.
export function useDateTime() {
  const { i18n } = useTranslation();
  return (iso: string) => dateTime(iso, i18n.language);
}

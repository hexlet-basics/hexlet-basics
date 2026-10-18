import dayjs from "dayjs";
import "dayjs/locale/es";
import "dayjs/locale/ru";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);
dayjs.extend(localizedFormat);

// "3 days ago" in the page's language, as legacy rendered `dayjs().to(date)`.
// The locale is set on the instance, never globally: the SSR server renders
// Several locales concurrently and a global dayjs locale would be shared.
export function timeAgo(iso: string, locale: string): string {
  return dayjs(iso).locale(locale).fromNow();
}

// The long localized date ("July 26, 2026", "26 июля 2026 г."), as legacy
// Rendered `dayjs(date).format("LL")`. Per-instance locale, for the same reason.
export function longDate(iso: string, locale: string): string {
  return dayjs(iso).locale(locale).format("LL");
}

// Date and time in the short localized form ("07/26/2026 3:30 PM", "26.07.2026
// 15:30"): what a table cell of the back office shows. The language is the app's,
// Never the browser's, so a page reads the same wherever it is opened.
export function dateTime(iso: string, locale: string): string {
  return dayjs(iso).locale(locale).format("L LT");
}

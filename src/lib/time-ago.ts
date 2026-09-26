import dayjs from "dayjs";
import "dayjs/locale/es";
import "dayjs/locale/ru";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);
dayjs.extend(localizedFormat);

// "3 days ago" in the page's language, as legacy rendered `dayjs().to(date)`.
// The locale is set on the instance, never globally: the SSR server renders
// several locales concurrently and a global dayjs locale would be shared.
export function timeAgo(iso: string, locale: string): string {
  return dayjs(iso).locale(locale).fromNow();
}

// The long localized date ("July 26, 2026", "26 июля 2026 г."), as legacy
// rendered `dayjs(date).format("LL")`. Per-instance locale, for the same reason.
export function longDate(iso: string, locale: string): string {
  return dayjs(iso).locale(locale).format("LL");
}

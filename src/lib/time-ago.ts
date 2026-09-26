import dayjs from "dayjs";
import "dayjs/locale/es";
import "dayjs/locale/ru";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);

// "3 days ago" in the page's language, as legacy rendered `dayjs().to(date)`.
// The locale is set on the instance, never globally: the SSR server renders
// several locales concurrently and a global dayjs locale would be shared.
export function timeAgo(iso: string, locale: string): string {
  return dayjs(iso).locale(locale).fromNow();
}

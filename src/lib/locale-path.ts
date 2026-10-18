import { DEFAULT_LOCALE, PREFIXED_LOCALES } from "@/lib/i18n";
import type { Locale } from "@/lib/i18n";

// The one place that knows how a locale shows in a URL path (ADR-0002 /
// ADR-0008): the default locale is unprefixed, every other contract locale is
// The first segment. The Go side of the same rule is localization.PathPrefix;
// Both are checked against the same table (see locale-path.test.ts).

// The prefixed locale a path starts with, if any. `/ruby` is not `/ru`.
function prefixOf(pathname: string): (typeof PREFIXED_LOCALES)[number] | undefined {
  const firstSegment = pathname.split("/").find(Boolean);
  return PREFIXED_LOCALES.find((locale) => locale === firstSegment);
}

// The backend has no locale-prefixed routes of its own, so API requests carry
// The frontend route locale explicitly through Accept-Language.
export function localeOf(pathname: string): Locale {
  return prefixOf(pathname) ?? DEFAULT_LOCALE;
}

// The path without its locale prefix, always starting with `/`.
export function stripLocale(pathname: string): string {
  const prefix = prefixOf(pathname);
  if (!prefix) {
    return pathname;
  }
  return pathname.slice(pathname.indexOf(prefix) + prefix.length) || "/";
}

// The same page in another locale (the language switcher's link target).
export function localize(pathname: string, locale: Locale): string {
  const base = stripLocale(pathname);
  if (locale === DEFAULT_LOCALE) {
    return base;
  }
  return `/${locale}${base === "/" ? "" : base}`;
}

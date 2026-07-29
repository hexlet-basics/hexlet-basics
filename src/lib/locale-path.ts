import { DEFAULT_LOCALE, type Locale, PREFIXED_LOCALES } from "@/lib/i18n";

// The backend has no locale-prefixed routes of its own, so API requests carry
// the frontend route locale explicitly through Accept-Language.
export function localeFromPathname(pathname: string): Locale {
  const firstSegment = pathname.split("/").find(Boolean);
  const prefixed = PREFIXED_LOCALES.find((locale) => locale === firstSegment);
  return prefixed ?? DEFAULT_LOCALE;
}

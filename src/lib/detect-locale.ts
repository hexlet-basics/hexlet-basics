import { createIsomorphicFn } from "@tanstack/react-start";
import { getCookie, getRequestHeader } from "@tanstack/react-start/server";
import { isLocale, type Locale } from "@/lib/i18n";

// First app-supported locale in Accept-Language (region stripped), mirroring
// legacy `locale_from_header`.
function localeFromAcceptLanguage(header: string | undefined): Locale | undefined {
  if (!header) return undefined;
  for (const part of header.split(",")) {
    const tag = part.split(";")[0]?.trim().split("-")[0];
    if (tag && isLocale(tag)) return tag;
  }
  return undefined;
}

// Locale to redirect the unprefixed root to, or undefined to stay on `en`.
// Parity with legacy `prepare_locale_settings`: a locale the visitor switched
// to (the `locale` cookie switchLocale sets, legacy `session[:locale]`) wins,
// so choosing English keeps a Russian-preferring browser on `/`; without one,
// a Russian-preferring browser goes to `/ru`. Runs server-side only (the
// client impl is a no-op, so the server-only imports are stripped from the
// browser bundle).
//
// Still stubbed: country-by-IP via getRequestIP → "RU", skip-redirect-for-bots
// (SEO), and legacy's remembering of the locale on every non-root visit.
export const detectRootLocale = createIsomorphicFn()
  .server((): Locale | undefined => {
    const remembered = getCookie("locale");
    if (remembered && isLocale(remembered)) {
      return remembered === "en" ? undefined : remembered;
    }
    const fromHeader = localeFromAcceptLanguage(getRequestHeader("accept-language"));
    return fromHeader === "ru" ? "ru" : undefined;
  })
  .client((): Locale | undefined => undefined);

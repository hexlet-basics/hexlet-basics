import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";
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
// Parity with legacy `prepare_locale_settings`: send a Russian-preferring
// browser to `/ru`. Runs server-side only (the client impl is a no-op, so the
// server-only header import is stripped from the browser bundle).
//
// Stubbed until their backends land: remembered `session[:locale]` (auth),
// country-by-IP via getRequestIP → "RU", and skip-redirect-for-bots (SEO).
export const detectRootLocale = createIsomorphicFn()
  .server((): Locale | undefined => {
    const fromHeader = localeFromAcceptLanguage(getRequestHeader("accept-language"));
    return fromHeader === "ru" ? "ru" : undefined;
  })
  .client((): Locale | undefined => undefined);

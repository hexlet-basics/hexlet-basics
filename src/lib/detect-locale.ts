import { createIsomorphicFn } from "@tanstack/react-start";
import { getCookie, getRequestHeader, setCookie } from "@tanstack/react-start/server";
import { isLocale, type Locale } from "@/lib/i18n";

// First app-supported locale in Accept-Language (region stripped), mirroring
// legacy `locale_from_header`.
function localeFromAcceptLanguage(header: string | undefined): Locale | undefined {
  if (!header) {
    return undefined;
  }
  for (const part of header.split(",")) {
    const tag = part.split(";")[0]?.trim().split("-")[0];
    if (tag && isLocale(tag)) {
      return tag;
    }
  }
  return undefined;
}

// The cookie the SSR server keeps the visitor's locale in (ADR-0002). Only the
// root `/` reads it, and following a language link changes no other state.
const LOCALE_COOKIE = "locale";
const LOCALE_COOKIE_MAX_AGE = 365 * 24 * 60 * 60;

// Remember the locale of the page being served, so the root `/` sends the
// visitor back to it. Runs server-side only (the client impl is a no-op, so the
// server-only imports are stripped from the browser bundle).
export const rememberLocale = createIsomorphicFn()
  .server((locale: Locale): void => {
    setCookie(LOCALE_COOKIE, locale, {
      path: "/",
      maxAge: LOCALE_COOKIE_MAX_AGE,
      httpOnly: true,
      sameSite: "lax",
    });
  })
  .client((): void => {});

// Locale to redirect the unprefixed root to, or undefined to stay on `en`.
// The remembered locale wins, then a Russian-preferring Accept-Language goes to
// `/ru`, else English. A visit arriving from a page of this site is a link the
// visitor followed (the English entry of the language menu is `/`), so it is
// not sent back to the remembered locale. Server-side only, as above.
export const detectRootLocale = createIsomorphicFn()
  .server((): Locale | undefined => {
    const host = getRequestHeader("host");
    const referer = getRequestHeader("referer");
    if (host && referer && URL.canParse(referer) && new URL(referer).host === host) {
      return undefined;
    }
    const remembered = getCookie(LOCALE_COOKIE);
    if (remembered && isLocale(remembered)) {
      return remembered === "en" ? undefined : remembered;
    }
    const fromHeader = localeFromAcceptLanguage(getRequestHeader("accept-language"));
    return fromHeader === "ru" ? "ru" : undefined;
  })
  .client((): Locale | undefined => undefined);

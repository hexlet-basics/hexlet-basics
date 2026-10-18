import { describe, expect, it } from "vitest";
import { PREFIXED_LOCALES } from "@/lib/i18n";
import { localeOf, localize, stripLocale } from "@/lib/locale-path";

// Cross-reference: the same prefix rule is tested in Go by TestPathPrefix in
// internal/localization/paths_test.go. Keep the (locale, prefix) table below
// identical on both sides: en carries no prefix, every other locale carries
// `/<locale>`.
const PREFIX_CASES = [
  ["en", ""],
  ["ru", "/ru"],
  ["es", "/es"],
] as const;

describe("PREFIXED_LOCALES", () => {
  it("is every contract locale except the unprefixed default", () => {
    expect([...PREFIXED_LOCALES].toSorted()).toStrictEqual(["es", "ru"]);
  });
});

describe("localeOf", () => {
  it.each([
    ["/", "en"],
    ["/languages/javascript", "en"],
    ["/ru", "ru"],
    ["/ru/admin/courses", "ru"],
    ["/es/", "es"],
    ["/es/session/new", "es"],
    ["/ruby", "en"],
  ])("maps %s to %s", (pathname, expected) => {
    expect(localeOf(pathname)).toBe(expected);
  });
});

describe("stripLocale", () => {
  it.each([
    ["/", "/"],
    ["/ru", "/"],
    ["/ru/", "/"],
    ["/es/languages/php", "/languages/php"],
    ["/languages/php", "/languages/php"],
    ["/ruby", "/ruby"],
  ])("turns %s into %s", (pathname, expected) => {
    expect(stripLocale(pathname)).toBe(expected);
  });
});

describe("localize", () => {
  it.each([
    ["/", "en", "/"],
    ["/", "ru", "/ru"],
    ["/ru", "en", "/"],
    ["/ru/languages/php", "es", "/es/languages/php"],
    ["/es/languages/php", "en", "/languages/php"],
    ["/languages/php", "ru", "/ru/languages/php"],
  ] as const)("localizes %s to %s as %s", (pathname, locale, expected) => {
    expect(localize(pathname, locale)).toBe(expected);
  });

  it.each(PREFIX_CASES)("gives %s the prefix %j", (locale, prefix) => {
    expect(localize("/x", locale)).toBe(`${prefix}/x`);
  });
});

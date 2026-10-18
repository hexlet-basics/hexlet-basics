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

describe("the unprefixed default locale", () => {
  it("is every contract locale except the unprefixed default", () => {
    expect.hasAssertions();
    expect([...PREFIXED_LOCALES].toSorted()).toStrictEqual(["es", "ru"]);
  }, 60_000);
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
  ])(
    "maps %s to %s",
    (pathname, expected) => {
      expect.hasAssertions();
      expect(localeOf(pathname)).toBe(expected);
    },
    60_000,
  );
});

describe("stripLocale", () => {
  it.each([
    ["/", "/"],
    ["/ru", "/"],
    ["/ru/", "/"],
    ["/es/languages/php", "/languages/php"],
    ["/languages/php", "/languages/php"],
    ["/ruby", "/ruby"],
  ])(
    "turns %s into %s",
    (pathname, expected) => {
      expect.hasAssertions();
      expect(stripLocale(pathname)).toBe(expected);
    },
    60_000,
  );
});

describe("localize", () => {
  it.each([
    ["/", "en", "/"],
    ["/", "ru", "/ru"],
    ["/ru", "en", "/"],
    ["/ru/languages/php", "es", "/es/languages/php"],
    ["/es/languages/php", "en", "/languages/php"],
    ["/languages/php", "ru", "/ru/languages/php"],
  ] as const)(
    "localizes %s to %s as %s",
    (pathname, locale, expected) => {
      expect.hasAssertions();
      expect(localize(pathname, locale)).toBe(expected);
    },
    60_000,
  );

  it.each(PREFIX_CASES)(
    "gives %s the prefix %j",
    (locale, prefix) => {
      expect.hasAssertions();
      expect(localize("/x", locale)).toBe(`${prefix}/x`);
    },
    60_000,
  );
});

import { describe, expect, it } from "vitest";
import { localeFromPathname } from "@/lib/locale-path";

describe("localeFromPathname", () => {
  it.each([
    ["/", "en"],
    ["/languages/javascript", "en"],
    ["/ru", "ru"],
    ["/ru/admin/courses", "ru"],
    ["/es/", "es"],
    ["/es/session/new", "es"],
    ["/ruby", "en"],
  ])("maps %s to %s", (pathname, expected) => {
    expect(localeFromPathname(pathname)).toBe(expected);
  });
});

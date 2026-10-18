import { describe, expect, it } from "vitest";
import { pageCount } from "@/lib/page-count";

describe("pageCount", () => {
  it.each([
    { total: 0, perPage: 25, expected: 1 },
    { total: 1, perPage: 25, expected: 1 },
    { total: 25, perPage: 25, expected: 1 },
    { total: 26, perPage: 25, expected: 2 },
    { total: 30, perPage: 25, expected: 2 },
    { total: 100, perPage: 25, expected: 4 },
  ])(
    "$total rows at $perPage a page make $expected pages",
    ({ total, perPage, expected }) => {
      expect.hasAssertions();
      expect(pageCount(total, perPage)).toBe(expected);
    },
    60_000,
  );

  it("is one page for an empty list whatever the page size", () => {
    expect.hasAssertions();
    expect(pageCount(0, 0)).toBe(1);
  }, 60_000);
});

import { describe, expect, it } from "vitest";
import { pageCount } from "@/lib/page-count";

describe("pageCount", () => {
  it.each([
    { expected: 1, perPage: 25, total: 0 },
    { expected: 1, perPage: 25, total: 1 },
    { expected: 1, perPage: 25, total: 25 },
    { expected: 2, perPage: 25, total: 26 },
    { expected: 2, perPage: 25, total: 30 },
    { expected: 4, perPage: 25, total: 100 },
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

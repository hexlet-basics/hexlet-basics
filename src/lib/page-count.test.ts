import { describe, expect, it } from "vitest";
import { pageCount } from "@/lib/page-count";

describe("pageCount", () => {
  it.each([
    [0, 25, 1],
    [1, 25, 1],
    [25, 25, 1],
    [26, 25, 2],
    [30, 25, 2],
    [100, 25, 4],
  ])(
    "%i rows at %i a page make %i pages",
    (total, perPage, expected) => {
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

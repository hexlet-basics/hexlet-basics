import { describe, expect, it } from "vitest";
import { dateTime } from "@/lib/time-ago";

// The day depends on the machine's time zone, so the assertions pin the shape of
// the format, which is what follows the language.
describe("dateTime", () => {
  const iso = "2026-07-26T12:30:00Z";

  it("writes the English date month first", () => {
    expect.hasAssertions();
    expect(dateTime(iso, "en")).toMatch(/^\d{2}\/\d{2}\/2026 \d{1,2}:\d{2} [AP]M$/u);
  }, 60_000);

  it("writes the Russian date day first, on a 24-hour clock", () => {
    expect.hasAssertions();
    expect(dateTime(iso, "ru")).toMatch(/^\d{2}\.\d{2}\.2026 \d{1,2}:\d{2}$/u);
  }, 60_000);

  it("writes the Spanish date day first", () => {
    expect.hasAssertions();
    expect(dateTime(iso, "es")).toMatch(/^\d{1,2}\/\d{1,2}\/2026 \d{1,2}:\d{2}$/u);
  }, 60_000);
});

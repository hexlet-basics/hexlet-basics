import { describe, expect, it } from "vitest";
import { client } from "@/lib/api-client";

describe("the API client auth transport", () => {
  it("sends no CSRF token of its own", () => {
    expect.hasAssertions();
    expect(client.getConfig()).not.toHaveProperty("xsrfHeaderName");
  }, 60_000);
});

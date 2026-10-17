import { describe, expect, it } from "vitest";
import { client } from "@/lib/api-client";

describe("API client auth transport", () => {
  it("sends no CSRF token of its own", () => {
    expect(client.getConfig()).not.toHaveProperty("xsrfHeaderName");
  });
});

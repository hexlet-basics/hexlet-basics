import { describe, expect, it } from "vitest";
import { client } from "@/lib/api-client";

describe("API client auth transport", () => {
  it("calls the same-origin /api and sends no CSRF token of its own", () => {
    expect(client.getConfig()).toMatchObject({ baseURL: "" });
    expect(client.getConfig()).not.toHaveProperty("xsrfHeaderName");
  });
});

import { describe, expect, it } from "vitest";
import { client } from "@/lib/api-client";

describe("API client auth transport", () => {
  it("calls the same-origin /api and delegates XSRF handling to Axios", () => {
    expect(client.getConfig()).toMatchObject({
      baseURL: "",
      xsrfCookieName: "XSRF-TOKEN",
      xsrfHeaderName: "X-XSRF-TOKEN",
    });
  });
});

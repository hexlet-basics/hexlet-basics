import { describe, expect, it } from "vitest";
import { client } from "@/lib/api-client";

describe("API client auth transport", () => {
  it("delegates credential and XSRF handling to Axios", () => {
    expect(client.getConfig()).toMatchObject({
      withCredentials: true,
      withXSRFToken: true,
      xsrfCookieName: "XSRF-TOKEN",
      xsrfHeaderName: "X-XSRF-TOKEN",
    });
  });
});

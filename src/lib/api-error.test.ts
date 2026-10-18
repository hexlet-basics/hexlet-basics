import { isNotFound } from "@tanstack/react-router";
import { describe, expect, it } from "vitest";
import { ApiError, apiErrorOf, isApiError, rethrowMissingAsNotFound } from "@/lib/api-error";

describe("apiError", () => {
  it("turns a non-2xx response into an error carrying its status and parsed body", () => {
    expect.hasAssertions();
    const response = new Response(null, { status: 422 });
    const error = apiErrorOf({ errors: { email: ["taken"] } }, response);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({ status: 422, body: { errors: { email: ["taken"] } } });
  }, 60_000);

  it("leaves a failure with no response, a network error, as it was", () => {
    expect.hasAssertions();
    const failure = new TypeError("Failed to fetch");

    expect(apiErrorOf(failure, undefined)).toBe(failure);
  }, 60_000);

  it("recognises an ApiError, optionally by status", () => {
    expect.hasAssertions();
    const error = new ApiError(404, {});

    expect(isApiError(error)).toBe(true);
    expect(isApiError(error, 404)).toBe(true);
    expect(isApiError(error, 422)).toBe(false);
    expect(isApiError(new Error("x"))).toBe(false);
  }, 60_000);

  it("rethrows a 404 as the router's notFound and any other failure untouched", () => {
    expect.hasAssertions();
    let thrown: unknown;
    try {
      rethrowMissingAsNotFound(new ApiError(404, {}));
    } catch (error) {
      thrown = error;
    }
    expect(isNotFound(thrown)).toBe(true);

    const boom = new ApiError(500, {});
    expect(() => rethrowMissingAsNotFound(boom)).toThrow(boom);
  }, 60_000);
});

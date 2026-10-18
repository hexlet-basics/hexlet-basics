import { isNotFound } from "@tanstack/react-router";
import { describe, expect, it } from "vitest";
import { ApiError, apiErrorOf, isApiError, rethrowMissingAsNotFound } from "@/lib/api-error";
import { HTTP_NOT_FOUND, HTTP_SERVER_ERROR, HTTP_UNPROCESSABLE } from "@/lib/http-status";

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
    const error = new ApiError(HTTP_NOT_FOUND, {});

    expect(isApiError(error)).toBe(true);
    expect(isApiError(error, HTTP_NOT_FOUND)).toBe(true);
    expect(isApiError(error, HTTP_UNPROCESSABLE)).toBe(false);
    expect(isApiError(new Error("x"))).toBe(false);
  }, 60_000);

  it("rethrows a 404 as the router's notFound and any other failure untouched", () => {
    expect.hasAssertions();
    let thrown: unknown = null;
    try {
      rethrowMissingAsNotFound(new ApiError(HTTP_NOT_FOUND, {}));
    } catch (error) {
      thrown = error;
    }
    expect(isNotFound(thrown)).toBe(true);

    const boom = new ApiError(HTTP_SERVER_ERROR, {});
    expect(() => rethrowMissingAsNotFound(boom)).toThrow(boom);
  }, 60_000);
});

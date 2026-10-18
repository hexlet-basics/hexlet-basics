import { notFound } from "@tanstack/react-router";
import { HTTP_NOT_FOUND } from "@/lib/http-status";

// What the generated client throws for a non-2xx answer: the status and the
// parsed body (JSON when the server sent it, text otherwise). Routes and forms
// depend on this class, never on the HTTP library under the client.
export class ApiError extends Error {
  readonly status: number;
  readonly body: unknown;

  constructor(status: number, body: unknown) {
    super(`API request failed with status ${status}`);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }
}

// The client's error interceptor: only a failure that came with a response
// becomes an ApiError. A network error or an abort has no status and passes
// through as it is.
export function apiErrorOf(error: unknown, response: Response | undefined): unknown {
  if (!response || response.ok) {
    return error;
  }
  return new ApiError(response.status, error);
}

export function isApiError(error: unknown, status?: number): error is ApiError {
  return error instanceof ApiError && (status === undefined || error.status === status);
}

// For a loader's catch: a 404 from the API is a missing page, so it becomes the
// router's `notFound()`; anything else is rethrown as it was.
export function rethrowMissingAsNotFound(error: unknown): never {
  if (isApiError(error, HTTP_NOT_FOUND)) {
    throw notFound();
  }
  throw error;
}

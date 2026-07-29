import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

namespace HexletBasics;

// ---------------------------------------------------------------------------
// Shared primitives, envelopes and errors.
//
// Naming rule (see main.tsp): the domain concept is `Course`. The legacy Rails
// app called it `Language`; we keep `Language` ONLY in URL paths for backward
// compatibility (ADR-0002), but every type/operation/schema uses `Course`.
//
// These types are consumed by BOTH generators — ogen (Go server validators)
// and hey-api (TS client + zod). Keep response shapes uniform.
// ---------------------------------------------------------------------------

/** UI locale. Mirrors the frontend `Locale` union and the legacy locale set. */
enum Locale {
  ru: "ru",
  en: "en",
  es: "es",
}

/**
 * A page of results. Generic envelope reused by every admin list so the CRUD
 * engine (TanStack Table) can read pagination uniformly.
 */
@friendlyName("{name}Page", T)
model Page<T> {
  items: T[];

  /** Total rows across all pages (before pagination). */
  total: int32;

  /** 1-based page number this response represents. */
  page: int32;

  /** Page size used for this response. */
  perPage: int32;
}

/** Standard pagination + sorting query params for every admin list. */
model ListQuery {
  @query
  @minValue(1)
  page?: int32;

  @query
  @minValue(1)
  @maxValue(100)
  perPage?: int32;

  /** Sort field (ransack `sf`). */
  @query sortField?: string;

  /** Sort order (ransack `so`). */
  @query sortOrder?: "asc" | "desc";
}

/**
 * Field-level validation errors, keyed by field name (each value is the list
 * of messages for that field). Returned when a write fails validation —
 * including constraints the schema cannot express, like uniqueness.
 */
@error
model ValidationError {
  @statusCode statusCode: 422;
  errors: Record<string[]>;
}

/** A resource was not found. */
@error
model NotFoundError {
  @statusCode statusCode: 404;
  message: string;
}

/** The request is not authenticated (no/invalid session cookie). */
@error
model UnauthorizedError {
  @statusCode statusCode: 401;
  message: string;
}

/** The caller is authenticated but lacks permission for this resource. */
@error
model ForbiddenError {
  @statusCode statusCode: 403;
  message: string;
}

/**
 * RFC 9457 problem details returned by the transport when request decoding,
 * authorization, persistence, or an unexpected server failure prevents an
 * operation from producing its declared domain response.
 */
model ProblemDetails {
  /** URI identifying the problem type. `about:blank` means the HTTP status. */
  type: string;

  /** Short, localized summary for this problem type. */
  title: string;

  /** HTTP status code generated for this occurrence. */
  status: int32;

  /** Optional occurrence-specific explanation that is safe for clients. */
  detail?: string;

  /** Optional URI identifying this particular occurrence. */
  instance?: string;
}

/**
 * Default error response shared by every operation.
 *
 * `@error` emits an OpenAPI `default` response, keeping central transport
 * failures typed without enumerating every status on every operation.
 */
@error
model ApiError {
  @header contentType: "application/problem+json";
  @body body: ProblemDetails;
}

/** Empty 204 response for deletes and other content-less successes. */
model NoContent {
  @statusCode statusCode: 204;
}

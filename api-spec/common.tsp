import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

namespace HexletBasics;

// ---------------------------------------------------------------------------
// Authentication schemes.
//
// UserSession and AdminSession deliberately read the same session cookie
// (ADR-0017). Distinct scheme names give ogen two generated SecurityHandler
// methods, so the contract — rather than Go operation-name lists — selects
// whether a current database admin check is required. There is no CSRF token
// scheme: cross-origin writes are refused by the server from the headers a
// browser sets itself.
// ---------------------------------------------------------------------------

model UserSession {
  type: AuthType.apiKey;
  in: ApiKeyLocation.cookie;
  name: "session";
}

model AdminSession {
  type: AuthType.apiKey;
  in: ApiKeyLocation.cookie;
  name: "session";
}

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

/** UI locale (the legacy locale set). The frontend `Locale` type is generated from this. */
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

/**
 * Standard pagination query params for every list. The default page size is
 * declared here once (the legacy pagy default); the generated server applies
 * it and every response echoes the size used in `perPage`.
 */
model ListQuery {
  @query
  @minValue(1)
  page?: int32;

  @query
  @minValue(1)
  @maxValue(100)
  perPage?: int32 = 20;
}

/**
 * A list an editor looks things up in: the page window plus an optional text
 * filter. `q` is a case-insensitive substring match over the resource's
 * human-readable fields (title, name, slug, text), applied in the query, so it
 * composes with the window. Searching is reading: the operation keeps its
 * Permission and its personal-data marker.
 */
model SearchListQuery {
  ...ListQuery;

  /** Text to find, matched case-insensitively against the names, slugs and text of the rows. */
  @query
  q?: string;
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
  @header contentType: "application/problem+json";
  @body body: ProblemDetails;
}

/** The caller is authenticated but lacks permission for this resource. */
@error
model ForbiddenError {
  @statusCode statusCode: 403;
  @header contentType: "application/problem+json";
  @body body: ProblemDetails;
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

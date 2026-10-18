---
status: accepted
---

# HTTP errors: Problem Details for every non-2xx, failures through the `error` return

Part of the backend clean-up (spec #905, ticket #917). Extends
[ADR-0018](0018-errors-via-oops.md) to the wire: how a domain error becomes an
HTTP answer. Supersedes the typed error models of the contract
(`NotFoundError`, `UnauthorizedError`, `ForbiddenError`, `ValidationError`).

**Problem.** An operation could fail in four ways that looked alike to a caller
and unlike to the code: a typed model in the contract (`ValidationError`,
`NotFoundError`, with the two others), a Problem Details document from the
central handler, a handler-local `refusal` carrying its own status and message,
and an ent error interpreted in the handler. A client had to decode a different
body per status, handlers returned `(success | error model, nil)` pairs, the
audit middleware had to recognise the validation model to record a failed
call, and every handler repeated how to word a not-found.

**Decision.**

- **One wire shape.** Every non-2xx response is RFC 9457 Problem Details
  (`application/problem+json`). The contract keeps one shared `ApiError` default
  response; the four typed models are removed, so an operation declares only its
  success.
- **A 422 adds an `errors` extension:** `errors: {field: [messages]}`, optional
  on `ProblemDetails` and present on a 422 only. It is what a form puts next to
  its fields.
- **Handlers return success only.** A failure is the `error` return. The one
  `errorStatus` table maps `domain.ErrNotFound`, `ErrConflict`, `ErrInvalid` and
  `ErrForbidden` to 404, 409, 422 and 403; security failures keep ogen's own
  401. `ent.Is*` is not used in a handler or in `errorStatus`: ent's errors reach
  the table already translated by `domain.FromStore`.
- **Unique violations.** A module names the field behind a unique index with
  `domain.UniqueOn`; a violation of it is `ErrInvalid` naming that field (422),
  any other constraint violation is `ErrConflict` (409). The account editor and
  the registrar declare the email index; other tables stay 409 until a module
  declares its fields.
- **Wording lives in the adapter.** `detailFor(err)` in `handlers/errors.go` is
  the table from an error to its localized `detail`: the self-protection
  refusals, the last administrator, a refused state move (named in the
  vocabulary the `transition.Illegal` carries) and the not-founds that explain
  themselves. The domain does not import localization. `refusal`,
  `withHTTPStatus`, `statusError` and `validationError` are deleted. A 422 whose
  field and message the handler knows is `invalid(field, message)`, an
  `ErrInvalid`.
- **Validation in three layers.** The contract (ogen, a 400 as generated), module
  rules returning `ErrInvalid`, and database constraints through `FromStore`;
  a handler validates nothing.
- **Audit.** A failed admin call is recognised by its error return, so the
  middleware no longer inspects a response type. A failure leaves an Audit
  Record with outcome `failed` and no changes, as before.
- **Unchanged.** 5xx are logged and sent to Sentry; a cancelled request is a
  silent 499; `detail` is never taken from `err.Error()`. A declared 429 (the
  assistant's daily quota) stays an operation response that is Problem Details.

## Considered Options

- **Keep the typed models for the three common statuses.** Rejected: the
  generated client then has a union per operation and every handler a second
  way to fail.
- **Put the localized message in the domain error.** Rejected: it makes every
  module depend on localization and fixes the wording at the place that cannot
  know the caller's language.
- **Make every unique violation a 422.** Rejected: the table cannot tell a
  field the person typed from a key the system derived; a module that knows
  declares it.

## Consequences

- The generated Go client returns a failure as an error
  (`*api.ProblemDetailsStatusCode`), so `testsupport.Client` no longer unwraps
  response unions, and the TS client reads `errors` from the Problem Details
  body (`zProblemDetails`).
- A new failure that needs its own wording is a row in `details`, not a new
  handler type.
- Frontend call sites that read the removed models move to the Problem Details
  fields (`src/lib/server-validation.ts`).

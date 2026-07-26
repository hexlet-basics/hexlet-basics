# Contract-first API pipeline with ogen

The API contract is authored once in TypeSpec (`api-spec/main.tsp`), emitted to
OpenAPI (`api-spec/dist/openapi.yaml`), and code-generated into a Go `net/http`
server via **ogen** and a TypeScript client via **hey-api**. One contract, no
client/server drift; generated artifacts are committed.

## Considered Options

- **ogen** (chosen) — generates value validation in code (enums, min/max,
  patterns, formats, required via `oas_validators_gen.go`), typed
  `Optional`/`Nullable` wrappers, reflection-free JSON via `go-faster/jx`, a
  static radix router, and native OpenTelemetry. Fits the generation-first rule
  (least hand-written code) and needs no runtime validation dependency.
- **oapi-codegen + Echo + kin-openapi** (rejected) — generates types only and
  offloads schema validation to a runtime middleware; no response validation,
  enums unchecked. More hand-wiring for less coverage.

## Consequences

- ogen is stricter about the spec: some OpenAPI constructs need the TypeSpec
  source adjusted rather than worked around in Go.
- `oapi-codegen`, `kin-openapi`, and `echo` leave the stack; `rs/cors` provides
  CORS on `net/http`.

import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";

// Shared building blocks first, then domain read models, then the operation
// surfaces. Each file reopens `namespace HexletBasics`.
import "./common.tsp";
import "./models.tsp";
import "./public.tsp";
import "./auth.tsp";
import "./admin.tsp";
import "./lesson-player.tsp";

using Http;
using Rest;
using OpenAPI;

/**
 * Hexlet Basics API.
 *
 * Contract-first (ADR-0001): this TypeSpec compiles to `dist/openapi.yaml`
 * (`pnpm spec`), which drives BOTH the Go server (ogen validators) and the TS
 * frontend client + zod (hey-api, `pnpm generate`). Never hand-edit `dist/`.
 *
 * Naming: the domain concept is `Course`. The legacy Rails app called it
 * `Language`; we keep `Language` ONLY in URL paths (e.g. `/languages`) for
 * backward compatibility (ADR-0002), but every type/operation/schema uses
 * `Course`. The full model/enum set lives in `models.tsp`.
 */
@service(#{ title: "Hexlet Basics API" })
namespace HexletBasics;

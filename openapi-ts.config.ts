import { defineConfig } from "@hey-api/openapi-ts";

// Generate the typed TS client + TanStack Query hooks + Zod validators from the
// Same OpenAPI the Go server is built from. Regenerate with `pnpm generate`
// After the contract changes.
//
// The `zod` plugin emits a schema per request/response body (e.g.
// `zBannerInput`) so TanStack Form validates against the SAME contract the CRUD
// Engine renders from (ADR-0008, FRONTEND_PLAN Wave 1).
export default defineConfig({
  input: "api-spec/dist/openapi.yaml",
  output: "src/client",
  plugins: ["@hey-api/client-fetch", "@tanstack/react-query", "zod"],
});

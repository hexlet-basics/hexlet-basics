import { defineConfig } from "@hey-api/openapi-ts";

// Generate the typed TS client + TanStack Query hooks from the same OpenAPI the
// Go server is built from. Regenerate with `pnpm generate` after the contract
// changes.
export default defineConfig({
  input: "api-spec/dist/openapi.yaml",
  output: "src/client",
  plugins: ["@hey-api/client-fetch", "@tanstack/react-query"],
});

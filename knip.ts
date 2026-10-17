import type { KnipConfig } from "knip";

// Unused files, exports and dependencies of the frontend and the TypeSpec
// tooling. legacy/ is the Rails app being retired and has its own package.json.
const config: KnipConfig = {
  ignore: ["legacy/**", "src/client/**", "src/types/**"],
  ignoreDependencies: [
    // Run by name from the SSR image (Dockerfile.web), never imported.
    "srvx",
    // The frontend state library the plan names (docs/FRONTEND_PLAN.md), not
    // used by any screen yet.
    "zustand",
    // Run as `pnpm exec tanstack search-docs`, never imported (see AGENTS.md).
    "@tanstack/cli",
    // TypeSpec packages are loaded by the compiler from api-spec/*.tsp.
    "@typespec/.*",
    // Grammar sources that `mise run lint:syntax` compiles with the tree-sitter CLI.
    // (The ones the lesson editor loads as wasm assets are imported, so knip sees them.)
    "tree-sitter-(css|embedded-template|go|html|javascript|json|python|ruby|typescript)",
    "@tree-sitter-grammars/tree-sitter-(toml|yaml)",
  ],
  ignoreExportsUsedInFile: true,
  entry: ["i18next.config.ts", "src/routes/**/*.{ts,tsx}"],
  // The end-to-end suite keeps its Playwright config beside its tests.
  playwright: { config: ["e2e/playwright.config.ts"], entry: ["e2e/tests/**/*.spec.ts"] },
};

export default config;

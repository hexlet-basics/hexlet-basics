import type { KnipConfig } from "knip";

// Unused files, exports and dependencies of the frontend and the TypeSpec
// Tooling. legacy/ is the Rails app being retired and has its own package.json.
const config: KnipConfig = {
  entry: ["i18next.config.ts", "src/routes/**/*.{ts,tsx}"],
  ignore: ["legacy/**", "src/client/**", "src/types/**"],
  ignoreDependencies: [
    // Run by name from the SSR image (Dockerfile.web), never imported.
    "srvx",
    // The frontend state library the plan names (docs/FRONTEND_PLAN.md), not
    // Used by any screen yet.
    "zustand",
    // Run as `pnpm exec tanstack search-docs`, never imported (see AGENTS.md).
    "@tanstack/cli",
    // TypeSpec packages are loaded by the compiler from api-spec/*.tsp.
    "@typespec/.*",
  ],
  ignoreExportsUsedInFile: true,
  // The end-to-end suite keeps its Playwright config beside its tests.
  playwright: { config: ["e2e/playwright.config.ts"], entry: ["e2e/tests/**/*.spec.ts"] },
};

export default config;

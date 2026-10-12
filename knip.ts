import type { KnipConfig } from "knip";

// Unused files, exports and dependencies of the frontend and the TypeSpec
// tooling. legacy/ is the Rails app being retired and has its own package.json.
const config: KnipConfig = {
  ignore: ["legacy/**", "src/client/**", "src/types/**"],
  ignoreDependencies: [
    // Run by name from the SSR image (Dockerfile.web) and the git hooks (hk.pkl),
    // never imported.
    "srvx",
    "@commitlint/cli",
    // The frontend state library the plan names (docs/FRONTEND_PLAN.md), not
    // used by any screen yet.
    "zustand",
    // TypeSpec packages are loaded by the compiler from api-spec/*.tsp.
    "@typespec/.*",
  ],
  ignoreExportsUsedInFile: true,
  entry: ["i18next.config.ts", "src/routes/**/*.{ts,tsx}"],
};

export default config;

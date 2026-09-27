import { defineConfig } from "oxlint";

export default defineConfig({
  plugins: [
    "eslint",
    "typescript",
    "unicorn",
    "oxc",
    "react",
    "import",
    "jsx-a11y",
    "vitest",
    "promise",
  ],
  options: {
    // Type-aware rules (no-floating-promises & co.) run through oxlint-tsgolint.
    typeAware: true,
  },
  rules: {
    // An unhandled promise is a silently swallowed failure; mark deliberate
    // fire-and-forget calls (navigate, invalidateQueries) with `void`.
    "typescript/no-floating-promises": "error",
    "typescript/no-misused-promises": "error",
    // A union grown by the contract must fail the switch, not fall into default.
    "typescript/switch-exhaustiveness-check": "error",
    "typescript/no-deprecated": "error",
  },
  ignorePatterns: [
    "legacy/**/*",
    "src/client/**",
    "src/routeTree.gen.ts",
    "src/types/resources.d.ts",
  ],
});

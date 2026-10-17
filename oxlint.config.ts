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
    // TanStack Router signals redirects and 404s by throwing plain objects.
    "typescript/only-throw-error": [
      "error",
      {
        allow: [
          {
            from: "package",
            package: "@tanstack/router-core",
            name: ["Redirect", "NotFoundError"],
          },
        ],
      },
    ],
  },
  overrides: [
    {
      // A spec says what a person does and sees through the `site` fixture. It
      // takes `test` from the support file, whose type hands it no `page`;
      // Playwright's own `test` would bring `page`, selectors and waits back.
      files: ["e2e/tests/**"],
      rules: {
        "no-restricted-imports": [
          "error",
          {
            paths: [
              {
                name: "@playwright/test",
                message: "import test from ../support/test and drive the site through its steps",
              },
            ],
          },
        ],
      },
    },
  ],
  ignorePatterns: [
    "legacy/**/*",
    "src/client/**",
    "src/routeTree.gen.ts",
    "src/types/resources.d.ts",
  ],
});

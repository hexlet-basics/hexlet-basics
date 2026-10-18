import { defineConfig } from "oxlint";

// Policy: every rule of every stable category is on. A rule is switched off
// only below, with the reason next to it (it conflicts with oxfmt or with
// another enabled rule, with a convention of this repo or of TanStack, or it is
// a taste rule whose fix churns hundreds of lines for no safety). There is no
// baseline and no blanket ignore: fix a finding, or justify turning its rule off.
export default defineConfig({
  plugins: [
    "eslint",
    "typescript",
    "unicorn",
    "oxc",
    "react",
    "import",
    "jsx-a11y",
    "promise",
    "react-perf",
    "node",
  ],
  categories: {
    correctness: "error",
    suspicious: "error",
    pedantic: "error",
    perf: "error",
    style: "error",
    restriction: "error",
  },
  options: {
    // Type-aware rules (no-floating-promises & co.) run through oxlint-tsgolint.
    typeAware: true,
  },
  rules: {
    "eslint/prefer-arrow-callback": ["error", { allowNamedFunctions: true }],
    "node/no-process-env": [
      "error",
      {
        allowedVariables: [
          "API_URL",
          "SITE_URL",
          "CI",
          "E2E_API_PORT",
          "E2E_WEB_PORT",
          "E2E_DB_HOST",
          "E2E_DB_PORT",
          "E2E_DB_NAME",
          "E2E_DB_PASSWORD",
        ],
      },
    ],
    "eslint/no-bitwise": ["error", { allow: ["|"] }],
    "eslint/no-magic-numbers": [
      "error",
      {
        ignore: [0, 1, -1, 2, 60_000],
        ignoreArrayIndexes: true,
        ignoreDefaultValues: true,
        ignoreEnums: true,
        ignoreNumericLiteralTypes: true,
        ignoreTypeIndexes: true,
      },
    ],
    "eslint/func-style": ["error", "declaration", { allowArrowFunctions: true }],
    "eslint/no-use-before-define": ["error", { functions: false }],
    "eslint/no-duplicate-imports": ["error", { allowSeparateTypeImports: true }],
    "unicorn/no-useless-undefined": [
      "error",
      { checkArguments: false, checkArrowFunctionBody: false },
    ],
    "unicorn/text-encoding-identifier-case": ["error", { withDash: true }],
    "eslint/capitalized-comments": ["error", "always", { ignoreConsecutiveComments: true }],
    "eslint/sort-imports": ["error", { ignoreDeclarationSort: true }],
    "eslint/one-var": ["error", "never"],
    // Rules switched off, each with its reason.
    // Pending: the rule flags every route component that is not exported next to Route.
    "react/only-export-components": "off",
    // The automatic JSX runtime needs no React import.
    "react/react-in-jsx-scope": "off",
    // User-facing strings are policed by lint:i18n (pages are exempt there on purpose).
    "react/jsx-no-literals": "off",
    // Demands deep readonly on generated client and library types we cannot change.
    "typescript/prefer-readonly-parameter-types": "off",
    // Taste rule, would reorder about 1000 object literals that are grouped by meaning.
    "eslint/sort-keys": "off",
    // Bans conventional short names (t, i, e, _) with no safety gain.
    "eslint/id-length": "off",
    // Mantine layouts nest deeply by design, splitting only adds indirection.
    "react/jsx-max-depth": "off",
    // Conditional rendering with a ternary is idiomatic React.
    "eslint/no-ternary": "off",
    // Taste rule, undefined is the idiomatic absent value next to optional props.
    "eslint/no-undefined": "off",
    // The repo marks deliberate fire-and-forget promises with void (no-floating-promises).
    "eslint/no-void": "off",
    // Restriction rule banning modern syntax the whole app is built on.
    "oxc/no-async-await": "off",
    // Restriction rule banning modern syntax the whole app is built on.
    "oxc/no-optional-chaining": "off",
    // Restriction rule banning modern syntax the whole app is built on.
    "oxc/no-rest-spread-properties": "off",
    // The API contract and DOM and React APIs use null.
    "unicorn/no-null": "off",
    // Return types are inferred and tsc checks callers.
    "typescript/explicit-function-return-type": "off",
    // Return types are inferred and tsc checks callers.
    "typescript/explicit-module-boundary-types": "off",
    // Contradicts no-default-export and TanStack route modules (named Route export).
    "import/no-named-export": "off",
    // Taste rule, conflicts with exporting at the declaration.
    "import/group-exports": "off",
    // Contradicts no-default-export.
    "import/prefer-default-export": "off",
    // Taste rule, conflicts with exporting at the declaration.
    "import/exports-last": "off",
    // Taste rule, page components legitimately import many modules.
    "import/max-dependencies": "off",
    // Taste rule, TypeScript members are public by default.
    "typescript/explicit-member-accessibility": "off",
    // The legal pages hold English prose with hundreds of quotes, escaping them hides the text.
    "react/no-unescaped-entities": "off",
    // CamelCase and PascalCase module names and TanStack $param routes are the repo convention.
    "unicorn/filename-case": "off",
    // Contradicts require-await and return-await, returning a promise directly is fine.
    "typescript/promise-function-async": "off",
    // Spreading Mantine and TanStack Form props is the library idiom.
    "react/jsx-props-no-spreading": "off",
    // Component modules default-export by convention and lazy imports rely on it.
    "import/no-default-export": "off",
    // Small colocated helper components are intended.
    "react/no-multi-comp": "off",
    // Taste threshold, page components are long declarative JSX.
    "eslint/max-lines-per-function": "off",
    // Taste threshold.
    "eslint/max-statements": "off",
    // Taste threshold.
    "eslint/max-lines": "off",
    // Taste threshold.
    "eslint/complexity": "off",
    // Micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-new-function-as-prop": "off",
    // Micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-new-object-as-prop": "off",
    // Micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-new-array-as-prop": "off",
    // Micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-jsx-as-prop": "off",
    // Forbids idiomatic truthiness checks on optional strings and objects.
    "typescript/strict-boolean-expressions": "off",
    // Forbids onClick={() => doIt()} shorthand used throughout.
    "typescript/no-confusing-void-expression": "off",
    // Side-effect imports (css, polyfills, setup) are the point.
    "import/no-unassigned-import": "off",
    // Mantine components take style and className by design.
    "react/forbid-component-props": "off",
    // False positive on render callbacks such as columns and actionsColumn edit and extra, which are not component definitions.
    "react/no-unstable-nested-components": "off",
    // Contradicts itself on forwardRef wrappers and createLink components and has no safety value.
    "react/function-component-definition": "off",
    // ForwardRef wrappers are named by the const they are assigned to.
    "react/display-name": "off",
    // Generated admin resources and e2e specs reach siblings with ../, the @/ alias covers src only.
    "import/no-relative-parent-imports": "off",
    // Monaco and the i18n catalogs are namespace modules.
    "import/no-namespace": "off",
    // StarterKit is the package default export and also a named one, both are the same value.
    "import/no-named-as-default": "off",
    // False positive, the matched function is a fetch interceptor not an Express handler.
    "oxc/no-async-endpoint-handlers": "off",
    // Keys of skeleton placeholders and of positional lists (footer groups, repeated audit changes) have no better identity.
    "react/no-array-index-key": "off",
    // Assertions sit at untyped boundaries (request.json(), Object.keys, JSON.parse, test router doubles) where a guard would only repeat zod validation.
    "typescript/no-unsafe-type-assertion": "off",
    // DangerouslySetInnerHTML carries the JSON-LD script of the blog article, escaped on purpose.
    "react/no-danger": "off",
    // Handlers come from records and hooks (actions[name].run), naming them handleX adds nothing.
    "react/jsx-handler-names": "off",
    // A lazily initialised per-mount constant has no setter to destructure.
    "react/hook-use-state": "off",
    // The default allows JSX only in .jsx files; this repo writes it in .tsx.
    "react/jsx-filename-extension": ["error", { extensions: [".tsx"] }],
    // An unhandled promise is a silently swallowed failure; mark deliberate
    // fire-and-forget calls (navigate, invalidateQueries) with `void`.
    "typescript/no-floating-promises": "error",
    "typescript/no-misused-promises": "error",
    // A union grown by the contract must fail the switch, not fall into default.
    "typescript/switch-exhaustiveness-check": "error",
    "typescript/no-deprecated": "error",
    // `any` switches the checker off for everything it touches; use `unknown`
    // and narrow, or a precise type.
    "typescript/no-explicit-any": "error",
    "typescript/no-unsafe-assignment": "error",
    "typescript/no-unsafe-member-access": "error",
    "typescript/no-unsafe-call": "error",
    "typescript/no-unsafe-return": "error",
    "typescript/no-unsafe-argument": "error",
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
      // These run in Node, not in the browser. A build script cannot use the async
      // fs API without top-level await, which node/no-top-level-await bans, so it
      // reads and writes synchronously.
      files: ["scripts/**", "vite.config.ts", "e2e/playwright.config.ts"],
      rules: {
        // Node's own modules are what a Node script is made of.
        "import/no-nodejs-modules": "off",
        // The contradiction with node/no-top-level-await is described above.
        "node/no-sync": "off",
      },
    },
    {
      // The vitest rules describe vitest specs; the Playwright e2e specs, scripts
      // and application modules are not vitest code, so the plugin is loaded for
      // the unit-test files only.
      files: ["**/*.test.ts", "**/*.test.tsx"],
      plugins: ["vitest"],
      rules: {
        // The stricter toBe(true) is the assertion, and prefer-strict-boolean-matchers requires it.
        "vitest/prefer-to-be-truthy": "off",
        // The stricter toBe(false) is the assertion, and prefer-strict-boolean-matchers requires it.
        "vitest/prefer-to-be-falsy": "off",
        // The valid-title rule rejects the function-valued describe title this rule demands.
        "vitest/prefer-describe-function-title": "off",
        // Contradicts prefer-importing-vitest-globals, explicit imports are used.
        "vitest/no-importing-vitest-globals": "off",
      },
    },
    {
      // A spec says what a person does and sees through the `site` fixture. It
      // takes `test` from the support file, whose type hands it no `page`;
      // Playwright's own `test` would bring `page`, selectors and waits back.
      files: ["e2e/support/**"],
      rules: {
        // Playwright fixtures hand the test body a `use` callback, which is not React's `use` hook.
        "react-hooks/rules-of-hooks": "off",
      },
    },
    {
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
    // Written by `msw init`, not by us.
    "public/mockServiceWorker.js",
  ],
});

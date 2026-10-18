import { defineConfig } from "oxlint";

// Policy: every rule of every stable category is on. A rule is switched off
// only below, with the reason next to it (it conflicts with oxfmt or with
// another enabled rule, with a convention of this repo or of TanStack, or it is
// a taste rule whose fix churns hundreds of lines for no safety). There is no
// baseline and no blanket ignore: fix a finding, or justify turning its rule off.
export default defineConfig({
  categories: {
    correctness: "error",
    pedantic: "error",
    perf: "error",
    restriction: "error",
    style: "error",
    suspicious: "error",
  },
  ignorePatterns: [
    "legacy/**/*",
    "src/client/**",
    "src/routeTree.gen.ts",
    "src/types/resources.d.ts",
    // Written by `msw init`, not by us.
    "public/mockServiceWorker.js",
  ],
  options: {
    // Type-aware rules (no-floating-promises & co.) run through oxlint-tsgolint.
    typeAware: true,
  },
  overrides: [
    {
      // A route module exports its `Route`, a router object and not a component, next to the
      // component it renders. oxlint 1.88 flags that component and ignores allowExportNames
      // for it; the router's own fast-refresh handling covers route files. Re-check on upgrade.
      files: ["src/routes/**"],
      rules: {
        "react/only-export-components": "off",
      },
    },
    {
      // The legal pages and the locale catalogs are text, and the tests list one case after another:
      // their length is their content, not a function to split.
      files: ["src/components/pages/**", "src/locales/**", "src/test/**"],
      rules: {
        "eslint/max-lines": "off",
        "eslint/max-lines-per-function": "off",
        "eslint/max-statements": "off",
      },
    },
    {
      // The legal pages are the text itself, one component per language, so the
      // prose has to sit in the JSX as written and its quotes are not escaped.
      files: ["src/components/pages/**"],
      rules: {
        "react/jsx-no-literals": "off",
        "react/no-unescaped-entities": "off",
      },
    },
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
        // Contradicts prefer-importing-vitest-globals, explicit imports are used.
        "vitest/no-importing-vitest-globals": "off",
        // The valid-title rule rejects the function-valued describe title this rule demands.
        "vitest/prefer-describe-function-title": "off",
        // The stricter toBe(false) is the assertion, and prefer-strict-boolean-matchers requires it.
        "vitest/prefer-to-be-falsy": "off",
        // The stricter toBe(true) is the assertion, and prefer-strict-boolean-matchers requires it.
        "vitest/prefer-to-be-truthy": "off",
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
                message: "import test from ../support/test and drive the site through its steps",
                name: "@playwright/test",
              },
            ],
          },
        ],
      },
    },
  ],
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
  rules: {
    "eslint/capitalized-comments": ["error", "always", { ignoreConsecutiveComments: true }],
    "eslint/complexity": ["error", { max: 25 }],
    "eslint/func-style": ["error", "declaration", { allowArrowFunctions: true }],
    "eslint/id-length": [
      "error",
      { checkGeneric: false, exceptions: ["$", "t", "_"], min: 2, properties: "never" },
    ],
    "eslint/max-lines": ["error", { max: 500, skipBlankLines: true, skipComments: true }],
    "eslint/max-lines-per-function": [
      "error",
      { max: 200, skipBlankLines: true, skipComments: true },
    ],
    "eslint/max-statements": ["error", { max: 30 }],
    "eslint/no-bitwise": ["error", { allow: ["|"] }],
    "eslint/no-duplicate-imports": ["error", { allowSeparateTypeImports: true }],
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
    // Contradicts unicorn/prefer-ternary, which asks for a ternary where an if/else assigns one value.
    "eslint/no-ternary": "off",
    // Contradicts unicorn/no-typeof-undefined: a comparison with `undefined` is how an absent value is tested.
    "eslint/no-undefined": "off",
    "eslint/no-use-before-define": ["error", { functions: false }],
    // The repo marks deliberate fire-and-forget promises with `void` as a statement (no-floating-promises).
    "eslint/no-void": ["error", { allowAsStatement: true }],
    "eslint/one-var": ["error", "never"],
    "eslint/prefer-arrow-callback": ["error", { allowNamedFunctions: true }],
    "eslint/sort-imports": ["error", { ignoreDeclarationSort: true }],
    "eslint/sort-keys": ["error", "asc", { allowLineSeparatedGroups: true }],
    // Taste rule, conflicts with exporting at the declaration.
    "import/exports-last": "off",
    // Taste rule, conflicts with exporting at the declaration.
    "import/group-exports": "off",
    // Taste rule, page components legitimately import many modules.
    "import/max-dependencies": "off",
    // Component modules default-export by convention and lazy imports rely on it.
    "import/no-default-export": "off",
    // StarterKit is the package default export and also a named one, both are the same value.
    "import/no-named-as-default": "off",
    // Contradicts no-default-export and TanStack route modules (named Route export).
    "import/no-named-export": "off",
    // Monaco and the i18n catalogs are namespace modules.
    "import/no-namespace": "off",
    // Generated admin resources and e2e specs reach siblings with ../, the @/ alias covers src only.
    "import/no-relative-parent-imports": "off",
    // Side-effect imports (css, polyfills, setup) are the point.
    "import/no-unassigned-import": "off",
    // Contradicts no-default-export.
    "import/prefer-default-export": "off",
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
    // Contradicts promise/prefer-await-to-then, which asks for async/await over promise chains.
    "oxc/no-async-await": "off",
    // Bans `?.`, the operator that typescript/strict-boolean-expressions code reads nullable values through.
    "oxc/no-optional-chaining": "off",
    // Contradicts eslint/prefer-object-spread, which asks for spread where Object.assign copies.
    "oxc/no-rest-spread-properties": "off",
    // The four react-perf rules flag a new identity per render, which costs only where a
    // memoised child compares props; nothing here uses React.memo.
    "react-perf/jsx-no-jsx-as-prop": "off",
    "react-perf/jsx-no-new-array-as-prop": "off",
    "react-perf/jsx-no-new-function-as-prop": "off",
    "react-perf/jsx-no-new-object-as-prop": "off",
    "react/function-component-definition": [
      "error",
      {
        namedComponents: ["function-declaration", "function-expression", "arrow-function"],
        unnamedComponents: ["arrow-function", "function-expression"],
      },
    ],
    // The default allows JSX only in .jsx files; this repo writes it in .tsx.
    "react/jsx-filename-extension": ["error", { extensions: [".tsx"] }],
    "react/jsx-max-depth": ["error", { max: 6 }],
    "react/jsx-no-literals": [
      "error",
      {
        allowedStrings: [
          "✕",
          "support@hexlet.io",
          "8 800 100 22 47",
          "+7 495 085 21 62",
          "ООО «Хекслет Рус»",
          "ОГРН 1217300010476",
          "Failed to load courses",
        ],
      },
    ],
    "react/jsx-props-no-spreading": ["error", { custom: "ignore", exceptions: ["html"] }],
    "react/no-multi-comp": ["error", { ignoreStateless: true }],
    "react/no-unstable-nested-components": [
      "error",
      { propNamePattern: "{render*,edit,extra,cell}" },
    ],
    "react/only-export-components": "error",
    // Rules switched off, each with its reason.
    // The automatic JSX runtime needs no React import.
    "react/react-in-jsx-scope": "off",
    "typescript/explicit-function-return-type": [
      "error",
      { allowExpressions: true, allowedNames: ["getRouter"] },
    ],
    // The one function left without a return type is getRouter: the router's Register
    // type is read from it, so naming the type there would make it refer to itself.
    "typescript/explicit-module-boundary-types": ["error", { allowedNames: ["getRouter"] }],
    "typescript/no-deprecated": "error",
    // `any` switches the checker off for everything it touches; use `unknown`
    // and narrow, or a precise type.
    "typescript/no-explicit-any": "error",
    // An unhandled promise is a silently swallowed failure; mark deliberate
    // fire-and-forget calls (navigate, invalidateQueries) with `void`.
    "typescript/no-floating-promises": "error",
    "typescript/no-misused-promises": "error",
    "typescript/no-unsafe-argument": "error",
    "typescript/no-unsafe-assignment": "error",
    "typescript/no-unsafe-call": "error",
    "typescript/no-unsafe-member-access": "error",
    "typescript/no-unsafe-return": "error",
    // TanStack Router signals redirects and 404s by throwing plain objects.
    "typescript/only-throw-error": [
      "error",
      {
        allow: [
          {
            from: "package",
            name: ["Redirect", "NotFoundError"],
            package: "@tanstack/router-core",
          },
        ],
      },
    ],
    // Demands deep readonly on generated client and library types we cannot change.
    "typescript/prefer-readonly-parameter-types": "off",
    // Contradicts require-await and return-await: a function that only returns a promise
    // would have to be async, then has no await to satisfy require-await, and an await
    // added for it is what return-await forbids.
    "typescript/promise-function-async": "off",
    "typescript/strict-boolean-expressions": [
      "error",
      { allowNullableBoolean: true, allowNullableString: true },
    ],
    // A union grown by the contract must fail the switch, not fall into default.
    "typescript/switch-exhaustiveness-check": "error",
    // CamelCase and PascalCase module names and TanStack $param routes are the repo convention.
    "unicorn/filename-case": "off",
    // The API contract and DOM and React APIs use null.
    "unicorn/no-null": "off",
    "unicorn/no-useless-undefined": [
      "error",
      { checkArguments: false, checkArrowFunctionBody: false },
    ],
    "unicorn/text-encoding-identifier-case": ["error", { withDash: true }],
  },
});

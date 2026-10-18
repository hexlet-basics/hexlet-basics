import { defineConfig } from "oxlint";

// Policy: every rule of every stable category is on. A rule is switched off
// Only below, with the reason next to it (it conflicts with oxfmt or with
// Another enabled rule, with a convention of this repo or of TanStack, or it is
// A taste rule whose fix churns hundreds of lines for no safety). There is no
// Baseline and no blanket ignore: fix a finding, or justify turning its rule off.
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
    // Rules switched off, each with its reason.
    // The automatic JSX runtime needs no React import.
    "react/react-in-jsx-scope": "off",
    // User-facing strings are policed by lint:i18n (pages are exempt there on purpose).
    "react/jsx-no-literals": "off",
    // Demands deep readonly on generated client and library types we cannot change.
    "typescript/prefer-readonly-parameter-types": "off",
    // Taste rule, would reorder about 1000 object literals that are grouped by meaning.
    "eslint/sort-keys": "off",
    // Conflicts with oxfmt, which owns import ordering.
    "eslint/sort-imports": "off",
    // Bans conventional short names (t, i, e, _) with no safety gain.
    "eslint/id-length": "off",
    // Contradicts init-declarations and one-declaration-per-line style.
    "eslint/one-var": "off",
    // Taste rule, flags every literal including HTTP codes, sizes and Mantine props.
    "eslint/no-magic-numbers": "off",
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
    // Expect.assertions(n) bookkeeping adds noise without safety.
    "vitest/prefer-expect-assertions": "off",
    // The global vitest timeout applies.
    "vitest/require-test-timeout": "off",
    // Flat test files are the repo style.
    "vitest/require-top-level-describe": "off",
    // Contradicts prefer-importing-vitest-globals, explicit imports are used.
    "vitest/no-importing-vitest-globals": "off",
    // Route modules and components are hoisted function declarations referenced from Route, and it contradicts function-component-definition.
    "eslint/func-style": "off",
    // TanStack route modules declare Route first and the hoisted component functions below it.
    "eslint/no-use-before-define": "off",
    // The legal pages hold English prose with hundreds of quotes, escaping them hides the text.
    "react/no-unescaped-entities": "off",
    // TanStack route files export Route next to their component by design.
    "react/only-export-components": "off",
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
    // Taste threshold.
    "eslint/max-params": "off",
    // Taste threshold, one scenario may assert several facts.
    "vitest/max-expects": "off",
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
    // The u flag changes escape semantics, existing patterns are ASCII.
    "eslint/require-unicode-regexp": "off",
    // Window is deliberate, it documents browser-only code behind SSR guards.
    "unicorn/prefer-global-this": "off",
    // Side-effect imports (css, polyfills, setup) are the point.
    "import/no-unassigned-import": "off",
    // E2e specs are Playwright *.spec.ts files, unit tests are *.test.ts.
    "vitest/consistent-test-filename": "off",
    // BeforeEach and afterEach are the way to reset shared state.
    "vitest/no-hooks": "off",
    // False positive on module-level code outside tests.
    "vitest/require-hook": "off",
    // Mantine components take style and className by design.
    "react/forbid-component-props": "off",
    // Env is read at the edge (playwright config, SSR site url).
    "node/no-process-env": "off",
    // Sequential awaits are intended where order matters.
    "eslint/no-await-in-loop": "off",
    // Tests assert that javascript: urls are neutralised.
    "eslint/no-script-url": "off",
    // Its fix drops arguments that zod .catch(undefined), mutate(undefined) and createIsomorphicFn callbacks need to type-check.
    "unicorn/no-useless-undefined": "off",
    // The HTML charset meta must read utf-8, the fix writes utf8.
    "unicorn/text-encoding-identifier-case": "off",
    // Describe titles name classes and components (ApiError), lowercasing misnames them.
    "vitest/prefer-lowercase-title": "off",
    // Type-only imports stay a separate statement because import/consistent-type-specifier-style requires it (import/no-duplicates still catches real duplicates).
    "eslint/no-duplicate-imports": "off",
    // False positive on render callbacks such as columns and actionsColumn edit and extra, which are not component definitions.
    "react/no-unstable-nested-components": "off",
    // Contradicts itself on forwardRef wrappers and createLink components and has no safety value.
    "react/function-component-definition": "off",
    // ForwardRef wrappers are named by the const they are assigned to.
    "react/display-name": "off",
    // Contradicts no-importing-vitest-globals, vitest specs import explicitly and e2e specs use Playwright test.
    "vitest/prefer-importing-vitest-globals": "off",
    // Generated admin resources and e2e specs reach siblings with ../, the @/ alias covers src only.
    "import/no-relative-parent-imports": "off",
    // A let assigned in a callback and read after it cannot be initialised with a sound placeholder.
    "eslint/init-declarations": "off",
    // ToBe(true) is the stricter assertion.
    "vitest/prefer-to-be-truthy": "off",
    // ToBe(false) is the stricter assertion.
    "vitest/prefer-to-be-falsy": "off",
    // False positive on zod .catch and TanStack Router loaders return chains.
    "promise/prefer-await-to-then": "off",
    // False positive on zod .catch, which is not a promise.
    "unicorn/prefer-top-level-await": "off",
    // Tests build deferred promises on purpose.
    "promise/avoid-new": "off",
    // Build scripts read and write files synchronously.
    "node/no-sync": "off",
    // Scripts and configs run in Node.
    "import/no-nodejs-modules": "off",
    // Node:path namespace import is the Node idiom in scripts.
    "unicorn/import-style": "off",
    // Spread in map is clearer, no hot path here.
    "oxc/no-map-spread": "off",
    // Taste rule, the groups are tiny.
    "eslint/prefer-named-capture-group": "off",
    // Monaco and the i18n catalogs are namespace modules.
    "import/no-namespace": "off",
    // Switch-exhaustiveness-check already proves every case, a default would hide a new one.
    "eslint/default-case": "off",
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
    // Describes name constants and classes too, and valid-title rejects a non-string title for those.
    "vitest/prefer-describe-function-title": "off",
    // Handlers come from records and hooks (actions[name].run), naming them handleX adds nothing.
    "react/jsx-handler-names": "off",
    // False positive, a default parameter does not replace null, which the argument may be.
    "unicorn/prefer-default-parameters": "off",
    // Monaco key chords (KeyMod.CtrlCmd | KeyCode.Enter) are bit flags by API.
    "eslint/no-bitwise": "off",
    // Number.parseFloat accepts the "50%" strings the Splitter reports, Number would give NaN.
    "unicorn/prefer-number-coercion": "off",
    // Contradicts React effects that return a cleanup on one path, and loaders and switches that end without a value.
    "typescript/consistent-return": "off",
    // A lazily initialised per-mount constant has no setter to destructure.
    "react/hook-use-state": "off",
    // ?? and ternaries only narrow nullable DOM and query results before the assertion that follows.
    "vitest/no-conditional-in-test": "off",
    // The default allows JSX only in .jsx files; this repo writes it in .tsx.
    "react/jsx-filename-extension": ["error", { extensions: [".tsx"] }],
    // An unhandled promise is a silently swallowed failure; mark deliberate
    // Fire-and-forget calls (navigate, invalidateQueries) with `void`.
    "typescript/no-floating-promises": "error",
    "typescript/no-misused-promises": "error",
    // A union grown by the contract must fail the switch, not fall into default.
    "typescript/switch-exhaustiveness-check": "error",
    "typescript/no-deprecated": "error",
    // `any` switches the checker off for everything it touches; use `unknown`
    // And narrow, or a precise type.
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
      // A spec says what a person does and sees through the `site` fixture. It
      // Takes `test` from the support file, whose type hands it no `page`;
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

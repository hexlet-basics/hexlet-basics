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
    // the automatic JSX runtime needs no React import.
    "react/react-in-jsx-scope": "off",
    // user-facing strings are policed by lint:i18n (pages are exempt there on purpose).
    "react/jsx-no-literals": "off",
    // demands deep readonly on generated client and library types we cannot change.
    "typescript/prefer-readonly-parameter-types": "off",
    // taste rule, would reorder about 1000 object literals that are grouped by meaning.
    "eslint/sort-keys": "off",
    // conflicts with oxfmt, which owns import ordering.
    "eslint/sort-imports": "off",
    // taste rule, would churn about 970 comments.
    "eslint/capitalized-comments": "off",
    // bans conventional short names (t, i, e, _) with no safety gain.
    "eslint/id-length": "off",
    // contradicts init-declarations and one-declaration-per-line style.
    "eslint/one-var": "off",
    // taste rule, flags every literal including HTTP codes, sizes and Mantine props.
    "eslint/no-magic-numbers": "off",
    // Mantine layouts nest deeply by design, splitting only adds indirection.
    "react/jsx-max-depth": "off",
    // conditional rendering with a ternary is idiomatic React.
    "eslint/no-ternary": "off",
    // taste rule, undefined is the idiomatic absent value next to optional props.
    "eslint/no-undefined": "off",
    // the repo marks deliberate fire-and-forget promises with void (no-floating-promises).
    "eslint/no-void": "off",
    // restriction rule banning modern syntax the whole app is built on.
    "oxc/no-async-await": "off",
    // restriction rule banning modern syntax the whole app is built on.
    "oxc/no-optional-chaining": "off",
    // restriction rule banning modern syntax the whole app is built on.
    "oxc/no-rest-spread-properties": "off",
    // the API contract and DOM and React APIs use null.
    "unicorn/no-null": "off",
    // return types are inferred and tsc checks callers.
    "typescript/explicit-function-return-type": "off",
    // return types are inferred and tsc checks callers.
    "typescript/explicit-module-boundary-types": "off",
    // contradicts no-default-export and TanStack route modules (named Route export).
    "import/no-named-export": "off",
    // taste rule, conflicts with exporting at the declaration.
    "import/group-exports": "off",
    // contradicts no-default-export.
    "import/prefer-default-export": "off",
    // taste rule, conflicts with exporting at the declaration.
    "import/exports-last": "off",
    // taste rule, page components legitimately import many modules.
    "import/max-dependencies": "off",
    // taste rule, TypeScript members are public by default.
    "typescript/explicit-member-accessibility": "off",
    // expect.assertions(n) bookkeeping adds noise without safety.
    "vitest/prefer-expect-assertions": "off",
    // the global vitest timeout applies.
    "vitest/require-test-timeout": "off",
    // flat test files are the repo style.
    "vitest/require-top-level-describe": "off",
    // contradicts prefer-importing-vitest-globals, explicit imports are used.
    "vitest/no-importing-vitest-globals": "off",
    // route modules and components are hoisted function declarations referenced from Route, and it contradicts function-component-definition.
    "eslint/func-style": "off",
    // TanStack route modules declare Route first and the hoisted component functions below it.
    "eslint/no-use-before-define": "off",
    // the legal pages hold English prose with hundreds of quotes, escaping them hides the text.
    "react/no-unescaped-entities": "off",
    // TanStack route files export Route next to their component by design.
    "react/only-export-components": "off",
    // camelCase and PascalCase module names and TanStack $param routes are the repo convention.
    "unicorn/filename-case": "off",
    // contradicts require-await and return-await, returning a promise directly is fine.
    "typescript/promise-function-async": "off",
    // spreading Mantine and TanStack Form props is the library idiom.
    "react/jsx-props-no-spreading": "off",
    // component modules default-export by convention and lazy imports rely on it.
    "import/no-default-export": "off",
    // small colocated helper components are intended.
    "react/no-multi-comp": "off",
    // taste threshold, page components are long declarative JSX.
    "eslint/max-lines-per-function": "off",
    // taste threshold.
    "eslint/max-statements": "off",
    // taste threshold.
    "eslint/max-lines": "off",
    // taste threshold.
    "eslint/complexity": "off",
    // taste threshold.
    "eslint/max-params": "off",
    // taste threshold, one scenario may assert several facts.
    "vitest/max-expects": "off",
    // micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-new-function-as-prop": "off",
    // micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-new-object-as-prop": "off",
    // micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-new-array-as-prop": "off",
    // micro-optimisation, no memoised children depend on identity.
    "react-perf/jsx-no-jsx-as-prop": "off",
    // forbids idiomatic truthiness checks on optional strings and objects.
    "typescript/strict-boolean-expressions": "off",
    // forbids onClick={() => doIt()} shorthand used throughout.
    "typescript/no-confusing-void-expression": "off",
    // the u flag changes escape semantics, existing patterns are ASCII.
    "eslint/require-unicode-regexp": "off",
    // window is deliberate, it documents browser-only code behind SSR guards.
    "unicorn/prefer-global-this": "off",
    // side-effect imports (css, polyfills, setup) are the point.
    "import/no-unassigned-import": "off",
    // e2e specs are Playwright *.spec.ts files, unit tests are *.test.ts.
    "vitest/consistent-test-filename": "off",
    // beforeEach and afterEach are the way to reset shared state.
    "vitest/no-hooks": "off",
    // false positive on module-level code outside tests.
    "vitest/require-hook": "off",
    // Mantine components take style and className by design.
    "react/forbid-component-props": "off",
    // env is read at the edge (playwright config, SSR site url).
    "node/no-process-env": "off",
    // sequential awaits are intended where order matters.
    "eslint/no-await-in-loop": "off",
    // tests assert that javascript: urls are neutralised.
    "eslint/no-script-url": "off",
    // its fix drops arguments that zod .catch(undefined), mutate(undefined) and createIsomorphicFn callbacks need to type-check.
    "unicorn/no-useless-undefined": "off",
    // the HTML charset meta must read utf-8, the fix writes utf8.
    "unicorn/text-encoding-identifier-case": "off",
    // describe titles name classes and components (ApiError), lowercasing misnames them.
    "vitest/prefer-lowercase-title": "off",
    // type-only imports stay a separate statement because import/consistent-type-specifier-style requires it (import/no-duplicates still catches real duplicates).
    "eslint/no-duplicate-imports": "off",
    // false positive on render callbacks such as columns and actionsColumn edit and extra, which are not component definitions.
    "react/no-unstable-nested-components": "off",
    // contradicts itself on forwardRef wrappers and createLink components and has no safety value.
    "react/function-component-definition": "off",
    // forwardRef wrappers are named by the const they are assigned to.
    "react/display-name": "off",
    // contradicts no-importing-vitest-globals, vitest specs import explicitly and e2e specs use Playwright test.
    "vitest/prefer-importing-vitest-globals": "off",
    // generated admin resources and e2e specs reach siblings with ../, the @/ alias covers src only.
    "import/no-relative-parent-imports": "off",
    // a let assigned in a callback and read after it cannot be initialised with a sound placeholder.
    "eslint/init-declarations": "off",
    // toBe(true) is the stricter assertion.
    "vitest/prefer-to-be-truthy": "off",
    // toBe(false) is the stricter assertion.
    "vitest/prefer-to-be-falsy": "off",
    // false positive on zod .catch and TanStack Router loaders return chains.
    "promise/prefer-await-to-then": "off",
    // false positive on zod .catch, which is not a promise.
    "unicorn/prefer-top-level-await": "off",
    // tests build deferred promises on purpose.
    "promise/avoid-new": "off",
    // build scripts read and write files synchronously.
    "node/no-sync": "off",
    // scripts and configs run in Node.
    "import/no-nodejs-modules": "off",
    // node:path namespace import is the Node idiom in scripts.
    "unicorn/import-style": "off",
    // spread in map is clearer, no hot path here.
    "oxc/no-map-spread": "off",
    // taste rule, the groups are tiny.
    "eslint/prefer-named-capture-group": "off",
    // monaco and the i18n catalogs are namespace modules.
    "import/no-namespace": "off",
    // switch-exhaustiveness-check already proves every case, a default would hide a new one.
    "eslint/default-case": "off",
    // StarterKit is the package default export and also a named one, both are the same value.
    "import/no-named-as-default": "off",
    // false positive, the matched function is a fetch interceptor not an Express handler.
    "oxc/no-async-endpoint-handlers": "off",
    // keys of skeleton placeholders and of positional lists (footer groups, repeated audit changes) have no better identity.
    "react/no-array-index-key": "off",
    // assertions sit at untyped boundaries (request.json(), Object.keys, JSON.parse, test router doubles) where a guard would only repeat zod validation.
    "typescript/no-unsafe-type-assertion": "off",
    // dangerouslySetInnerHTML carries the JSON-LD script of the blog article, escaped on purpose.
    "react/no-danger": "off",
    // describes name constants and classes too, and valid-title rejects a non-string title for those.
    "vitest/prefer-describe-function-title": "off",
    // handlers come from records and hooks (actions[name].run), naming them handleX adds nothing.
    "react/jsx-handler-names": "off",
    // false positive, a default parameter does not replace null, which the argument may be.
    "unicorn/prefer-default-parameters": "off",
    // monaco key chords (KeyMod.CtrlCmd | KeyCode.Enter) are bit flags by API.
    "eslint/no-bitwise": "off",
    // Number.parseFloat accepts the "50%" strings the Splitter reports, Number would give NaN.
    "unicorn/prefer-number-coercion": "off",
    // contradicts React effects that return a cleanup on one path, and loaders and switches that end without a value.
    "typescript/consistent-return": "off",
    // a lazily initialised per-mount constant has no setter to destructure.
    "react/hook-use-state": "off",
    // ?? and ternaries only narrow nullable DOM and query results before the assertion that follows.
    "vitest/no-conditional-in-test": "off",
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

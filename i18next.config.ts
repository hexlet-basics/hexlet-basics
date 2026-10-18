import { defineConfig } from "i18next-cli";

// I18next-cli config for the Go-rewrite frontend. Its primary job here is
// Regenerating the type definitions (`src/types/resources.d.ts` +
// `src/types/i18next.d.ts`) from the Russian source of truth after editing the
// Locale files — run `pnpm exec i18next-cli types`. The `extract` block mirrors
// The legacy setup so key extraction from source stays available, but the
// Hand-maintained locale files under `src/locales` are never overwritten by it
// (they are ignored as inputs).
const ignore = [
  "src/locales/**",
  "src/types/**",
  "src/client/**",
  "src/routeTree.gen.ts",
  "src/**/*.d.ts",
];

export default defineConfig({
  extract: {
    defaultNS: "translation",
    ignore,
    input: ["src/**/*.tsx", "src/**/*.ts"],
    mergeNamespaces: true,
    output: "src/locales/{{language}}/translation.ts",
    outputFormat: "ts",
    primaryLanguage: "ru",
    removeUnusedKeys: false,
  },
  lint: {
    ignore: [
      ...ignore,
      // Legal and about pages are written as one component per language
      // (PrivacyRu/PrivacyEn, TosRu/TosEn, ...), their text is the page, not a
      // Catalog entry. The hardcoded-string rule does not apply to them.
      "src/components/pages/**",
    ],
  },
  locales: ["ru", "es", "en"],
  types: {
    enableSelector: true,
    input: ["src/locales/ru/translation.ts"],
    output: "src/types/i18next.d.ts",
    resourcesFile: "src/types/resources.d.ts",
  },
});

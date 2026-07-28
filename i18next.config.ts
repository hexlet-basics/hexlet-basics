import { defineConfig } from "i18next-cli";

// i18next-cli config for the Go-rewrite frontend. Its primary job here is
// regenerating the type definitions (`src/types/resources.d.ts` +
// `src/types/i18next.d.ts`) from the Russian source of truth after editing the
// locale files — run `pnpm exec i18next-cli types`. The `extract` block mirrors
// the legacy setup so key extraction from source stays available, but the
// hand-maintained locale files under `src/locales` are never overwritten by it
// (they are ignored as inputs).
const ignore = [
  "src/locales/**",
  "src/types/**",
  "src/client/**",
  "src/routeTree.gen.ts",
  "src/**/*.d.ts",
];

export default defineConfig({
  locales: ["ru", "es", "en"],
  extract: {
    ignore,
    input: ["src/**/*.tsx", "src/**/*.ts"],
    output: "src/locales/{{language}}/translation.ts",
    outputFormat: "ts",
    mergeNamespaces: true,
    removeUnusedKeys: false,
    defaultNS: "translation",
    primaryLanguage: "ru",
  },
  lint: { ignore },
  types: {
    input: ["src/locales/ru/translation.ts"],
    output: "src/types/i18next.d.ts",
    resourcesFile: "src/types/resources.d.ts",
    enableSelector: true,
  },
});

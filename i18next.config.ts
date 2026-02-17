import { defineConfig } from "i18next-cli";

const ignore = [
  "app/javascript/**/*.ru.tsx",
  "app/javascript/**/*.kz.tsx",
  "app/javascript/**/*.es.tsx",
  "app/javascript/**/*.en.tsx",
  "app/javascript/locales/**/*.ts",
  "app/javascript/types/**/*.ts",
  "app/javascript/generated/**/*.ts",
  "app/javascript/**/*.d.ts",
];

export default defineConfig({
  locales: ["ru", "es", "en"],
  extract: {
    ignore: ignore,
    input: ["app/javascript/**/*.tsx", "app/javascript/**/*.ts"],
    output: "app/javascript/locales/{{language}}/translation.ts",
    outputFormat: "ts",
    mergeNamespaces: true,
    removeUnusedKeys: false,
    defaultNS: "translation",
    primaryLanguage: "ru",
  },

  lint: {
    ignore: ignore,
  },
  types: {
    input: ["app/javascript/locales/ru/translation.ts"],
    output: "app/javascript/types/i18next.d.ts",
    resourcesFile: "app/javascript/types/resources.d.ts",
    enableSelector: true, // Enable type-safe key selection
  },
});

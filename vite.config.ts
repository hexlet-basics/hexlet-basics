import { fileURLToPath, URL } from "node:url";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";
// `vitest/config` extends vite's defineConfig with the `test` field.
import { defineConfig } from "vitest/config";

// https://vite.dev/config/
export default defineConfig({
  plugins: [tanstackStart(), viteReact()],
  // The browser reaches the Go API on the site's own origin under `/api`
  // (ADR-0015); in development Vite stands in for the production ingress.
  server: {
    proxy: {
      "/api": "http://localhost:3001",
    },
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  test: {
    // Long enough for the tests that wait on the editor; see `editorLoad` in
    // src/test/lesson-player.test.tsx for why that wait is what it is.
    testTimeout: 60_000,
    setupFiles: ["./src/test/setup.tsx"],
    browser: {
      enabled: true,
      provider: playwright(),
      instances: [{ browser: "chromium" }],
      headless: true,
    },
  },
});

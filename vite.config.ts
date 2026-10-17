import { fileURLToPath, URL } from "node:url";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";
// `vitest/config` extends vite's defineConfig with the `test` field.
import { defineConfig } from "vitest/config";

// https://vite.dev/config/
export default defineConfig({
  // `compiler`: React Compiler via oxc (`oxc-transform-react`), not the Babel
  // plugin — it memoizes components, so the code carries no manual useMemo.
  plugins: [tanstackStart(), viteReact({ compiler: true })],
  // The browser reaches the Go API on the site's own origin under `/api`
  // (ADR-0015); in development Vite stands in for the production ingress.
  server: {
    // The OAuth server (ADR-0020) is the one thing outside `/api`: its
    // endpoints under `/oauth` and the two metadata documents go to Go too, so
    // the consent flow works on the origin the browser is on. The well-known
    // keys are exact patterns: `/.well-known` as a whole belongs to the SSR server.
    proxy: {
      "/api": process.env.API_URL ?? "http://localhost:3001",
      "/oauth": process.env.API_URL ?? "http://localhost:3001",
      "^/\\.well-known/(oauth-authorization-server|openid-configuration)$":
        process.env.API_URL ?? "http://localhost:3001",
      // The MCP endpoint (ADR-0019) and its protected resource metadata.
      "^/mcp$": process.env.API_URL ?? "http://localhost:3001",
      "^/\\.well-known/oauth-protected-resource(/mcp)?$":
        process.env.API_URL ?? "http://localhost:3001",
    },
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  test: {
    // The default include also matches e2e/tests/*.spec.ts, the Playwright suite.
    // Vitest would load those in the browser, where Vite fails to bundle
    // @playwright/test and mounts its error overlay over every other test.
    include: ["src/**/*.test.{ts,tsx}"],
    // Long enough for the tests that wait on the editor; see `editorLoad` in
    // src/test/lesson-player.test.tsx for why that wait is what it is.
    testTimeout: 60_000,
    setupFiles: ["./src/test/setup.tsx"],
    coverage: {
      provider: "v8",
      include: ["src/**/*.{ts,tsx}"],
      // Generated code, the test harness itself, and the SSR entry points that a
      // test cannot mount under a synthetic root (the Playwright suite covers them).
      exclude: ["src/client/**", "src/routeTree.gen.ts", "src/test/**", "src/**/*.test.*"],
      // The floor the suite holds; raise it when coverage rises, never lower it.
      thresholds: { lines: 90, statements: 90, functions: 85, branches: 70 },
    },
    browser: {
      enabled: true,
      provider: playwright(),
      instances: [{ browser: "chromium" }],
      headless: true,
    },
  },
});

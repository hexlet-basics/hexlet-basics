import { fileURLToPath } from "node:url";
import { defineConfig, devices } from "@playwright/test";

// End-to-end suite: a real browser against the real stack, with nothing mocked.
// The Go API and the SSR web server start below against a throw-away database
// that `scripts/e2e_db.sh` (mise run test:e2e:db) migrated and loaded with the
// committed fixtures. `mise run test:e2e` does all of it.
const root = fileURLToPath(new URL("..", import.meta.url));

const apiPort = process.env.E2E_API_PORT ?? "3101";
const webPort = process.env.E2E_WEB_PORT ?? "4101";
const siteURL = `http://localhost:${webPort}`;

const dbHost = process.env.E2E_DB_HOST ?? "127.0.0.1";
const dbPort = process.env.E2E_DB_PORT ?? "54330";
const dbName = process.env.E2E_DB_NAME ?? "code_basics_e2e_test";
const dbPassword = process.env.E2E_DB_PASSWORD ?? "postgres";

export default defineConfig({
  forbidOnly: Boolean(process.env.CI),
  // Every journey signs up its own user or uses a fixture user read-only, so
  // the files run in parallel. A failure is a failure: no retries to hide a flake.
  fullyParallel: true,
  outputDir: "./.results",
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  reporter: process.env.CI
    ? [["github"], ["html", { open: "never", outputFolder: "./.report" }]]
    : "list",
  retries: 0,
  testDir: "./tests",
  use: {
    baseURL: siteURL,
    // The fixtures hold English, Russian and Spanish; the English site is the default.
    locale: "en-US",
    trace: "retain-on-failure",
  },
  webServer: [
    {
      // The Go API, built by the test:e2e task.
      command: "bin/server",
      cwd: root,
      env: {
        ADDR: `:${apiPort}`,
        // "test" also keeps the server from reading a developer's .env.
        APP_ENV: "test",
        APP_HOST: "localhost",
        BLOB_BUCKET_URL: `file://${root}e2e/.storage?create_dir=true`,
        DATABASE_URL: `postgres://postgres:${dbPassword}@${dbHost}:${dbPort}/${dbName}?sslmode=disable`,
        EMAIL_TOKEN_SECRET: "e2e-email-token-secret",
        GUEST_COOKIE_SECRET: "e2e-guest-cookie-secret",
        PUBLIC_URL: siteURL,
        SITE_URL: siteURL,
      },
      reuseExistingServer: false,
      timeout: 60_000,
      url: `http://127.0.0.1:${apiPort}/api/languages`,
    },
    {
      // The SSR frontend: the production build (e2e/.build, kept apart from
      // dist/ so a dev build cannot swap files under a running server). Preview
      // proxies /api to the Go server on the site's own origin, as the ingress
      // does in production (ADR-0015); API_URL is where SSR reaches Go.
      // Run through node, not `pnpm exec`: pnpm starts vite in a process group of
      // its own, which Playwright's shutdown signal never reaches, so the run
      // would hang after the last test with the server still listening.
      command: `node node_modules/vite/bin/vite.js preview --outDir e2e/.build --port ${webPort} --strictPort`,
      cwd: root,
      env: { API_URL: `http://127.0.0.1:${apiPort}` },
      reuseExistingServer: false,
      timeout: 60_000,
      url: siteURL,
    },
  ],
});

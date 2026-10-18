// Global setup for Vitest Browser Mode tests, loaded once per test file via
// `test.setupFiles` in vite.config.ts.
//
// - Mantine + dates + notifications styles so components render like production.
// - A fixed client baseUrl so the generated client emits absolute URLs MSW can
//   match; the real `@/lib/api-client` (with its SSR cookie interceptor) is never
//   imported in tests.
// - MSW worker lifecycle: start once, reset handlers between tests. No explicit
//   stop: msw 3's `worker.stop()` never settles in the browser runner (the hook
//   times out), and the page is discarded after each test file anyway.
import "@mantine/core/styles.css";
import "@mantine/code-highlight/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/tiptap/styles.css";
import "@mantine/notifications/styles.css";
import { notifications } from "@mantine/notifications";
import { afterEach, beforeAll } from "vitest";
// The production client setup (base URL, request and error interceptors), so a
// test goes through the same transport as the application.
import "@/lib/api-client";
import { worker } from "./msw";

beforeAll(() => worker.start({ onUnhandledFrame: "bypass", quiet: true }));
afterEach(() => {
  worker.resetHandlers();
  // Notifications render into a portal outside the rendered tree, so unmounting a
  // component leaves them on screen, where they sit over the next test's clicks.
  notifications.clean();
});

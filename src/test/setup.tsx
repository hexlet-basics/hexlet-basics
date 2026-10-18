// Global setup for Vitest Browser Mode tests, loaded once per test file via
// `test.setupFiles` in vite.config.ts.
//
// - Mantine + dates + notifications styles so components render like production.
// - A fixed client baseUrl so the generated client emits absolute URLs MSW can
//   Match; the real `@/lib/api-client` (with its SSR cookie interceptor) is never
//   Imported in tests.
// - MSW worker lifecycle: start once, reset handlers between tests. No explicit
//   Stop: msw 3's `worker.stop()` never settles in the browser runner (the hook
//   Times out), and the page is discarded after each test file anyway.
import "@mantine/core/styles.css";
import "@mantine/code-highlight/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/tiptap/styles.css";
import "@mantine/notifications/styles.css";
import { notifications } from "@mantine/notifications";
import { afterEach, beforeAll } from "vitest";
// The production client setup (base URL, request and error interceptors), so a
// Test goes through the same transport as the application.
import "@/lib/api-client";
import { worker } from "./msw";

beforeAll(() => worker.start({ onUnhandledFrame: "bypass", quiet: true }));
afterEach(() => {
  worker.resetHandlers();
  // Notifications render into a portal outside the rendered tree, so unmounting a
  // Component leaves them on screen, where they sit over the next test's clicks.
  notifications.clean();
});

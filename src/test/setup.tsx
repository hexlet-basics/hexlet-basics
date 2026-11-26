// Global setup for Vitest Browser Mode tests, loaded once per test file via
// `test.setupFiles` in vite.config.ts.
//
// - Mantine + dates + notifications styles so components render like production.
// - A fixed client baseUrl so the generated client emits absolute URLs MSW can
//   match; the real `@/lib/api-client` (with its SSR cookie interceptor) is never
//   imported in tests.
// - MSW worker lifecycle: start once, reset handlers between tests, and fail a test
//   whose page made an API request no handler answered: it would reach Vite's proxy
//   and the absent backend, and the page would render an error the test never meant.
//   No explicit stop: msw 3's `worker.stop()` never settles in the browser runner
//   (the hook times out), and the page is discarded after each test file anyway.
import "@mantine/core/styles.css";
import "@mantine/code-highlight/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/tiptap/styles.css";
import "@mantine/notifications/styles.css";
import { notifications } from "@mantine/notifications";
import { afterEach, beforeAll, expect } from "vitest";
import { cleanup } from "vitest-browser-react";
// The production client setup (base URL, request and error interceptors), so a
// test goes through the same transport as the application.
import "@/lib/api-client";
import { worker } from "./msw";

const unanswered: string[] = [];

// The API request a network frame carries; the page's own modules, styles and fonts are not one.
function apiRequest(data: unknown): Request | undefined {
  if (typeof data !== "object" || data === null || !("request" in data)) {
    return undefined;
  }
  const { request } = data;
  if (!(request instanceof Request) || !new URL(request.url).pathname.startsWith("/api/")) {
    return undefined;
  }
  return request;
}

beforeAll(() =>
  worker.start({
    onUnhandledFrame: ({ frame }) => {
      const request = apiRequest(frame.data);
      if (request !== undefined) {
        const url = new URL(request.url);
        unanswered.push(`${request.method} ${url.pathname}${url.search}`);
      }
    },
    quiet: true,
  }),
);
afterEach(async () => {
  // The vitest-browser-react cleanup runs in the next test's beforeEach; until then this test's page
  // would go on fetching (a refetch, a focus) against the handlers reset below, and its
  // requests would count against the next test. So the page goes first.
  await cleanup();
  worker.resetHandlers();
  // Notifications render into a portal outside the rendered tree, so unmounting a
  // component leaves them on screen, where they sit over the next test's clicks.
  notifications.clean();
  expect(unanswered.splice(0), "API requests no MSW handler answered").toStrictEqual([]);
});

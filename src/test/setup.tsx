// Global setup for Vitest Browser Mode tests, loaded once per test file via
// `test.setupFiles` in vite.config.ts.
//
// - Mantine + dates + notifications styles so components render like production.
// - A fixed client baseUrl so the generated client emits absolute URLs MSW can
//   match; the real `@/lib/api-client` (with its SSR cookie interceptor) is never
//   imported in tests.
// - MSW worker lifecycle: start once, reset handlers between tests, stop at end.
import "@mantine/core/styles.css";
import "@mantine/code-highlight/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/tiptap/styles.css";
import "@mantine/notifications/styles.css";
import { notifications } from "@mantine/notifications";
import { afterAll, afterEach, beforeAll } from "vitest";
import { client } from "@/client/client.gen";
import { worker } from "./msw";

client.setConfig({ baseURL: "http://localhost" });

beforeAll(() => worker.start({ onUnhandledRequest: "bypass", quiet: true }));
afterEach(() => worker.resetHandlers());
// Notifications render into a portal outside the rendered tree, so unmounting a
// component leaves them on screen, where they sit over the next test's clicks.
afterEach(() => notifications.clean());
afterAll(() => worker.stop());

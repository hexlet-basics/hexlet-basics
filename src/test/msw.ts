import { HttpResponse, http } from "msw";
import { setupWorker } from "msw/browser";
import type {
  ActiveBanner,
  CourseOption,
  CurrentUser,
  LessonAssistantChat,
  MyCapabilities,
} from "@/client/types.gen";

// What the API answers a guest on an empty site for the requests every page makes on its own
// (the layout's session, banner, catalogue menu; the admin shell's capabilities and course
// picker; the lesson page's assistant history). A test overrides any of them with
// `worker.use(...)`.
const ambient = [
  http.get("*/api/me", () => HttpResponse.json<CurrentUser>({ user: null })),
  http.get("*/api/banner", () => HttpResponse.json<ActiveBanner>({ banner: null })),
  http.get("*/api/languages", () => HttpResponse.json([])),
  http.get("*/api/language_categories", () => HttpResponse.json([])),
  http.get("*/api/admin/capabilities", () => HttpResponse.json<MyCapabilities>({ held: [] })),
  http.get("*/api/admin/courses/options", () => HttpResponse.json<CourseOption[]>([])),
  http.get("*/api/ai/lessons/:lessonId/messages", () =>
    HttpResponse.json<LessonAssistantChat>({ messages: [], quotaExceeded: false }),
  ),
];

// Shared MSW worker for Vitest Browser Mode. It starts with the ambient handlers only; each test
// declares what else it needs via `worker.use(http.post(...))`, and the setup file resets the
// handlers between tests. MSW intercepts the real fetch the generated hey-api client makes, so
// tests never touch the client's transport by hand.
export const worker = setupWorker(...ambient);

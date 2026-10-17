import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type {
  Course,
  CourseCatalogItem,
  DashboardCourse,
  MyDashboard,
  User,
} from "@/client/types.gen";
import { Route as myRoute } from "@/routes/{-$locale}/my";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const user: User = {
  id: 1,
  firstName: "Dora",
  lastName: null,
  name: "Dora",
  email: "dora@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: 0,
  createdAt: "2026-09-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

function catalogItem(id: number, slug: string, name: string): CourseCatalogItem {
  const course: Course = {
    id,
    slug,
    name,
    learnAs: "first_language",
    readiness: "completed",
    categoryId: null,
    currentVersionId: null,
    currentVersion: null,
    createdAt: "2026-01-01T00:00:00Z",
    enrollmentsCount: 0,
    lessonsCount: 10,
    ratingCount: 0,
    ratingValue: 0,
    repositoryUrl: null,
    hexletProgramLandingPage: null,
    coverListVariant: null,
    coverThumbVariant: null,
  };
  return {
    id,
    slug,
    header: name,
    name,
    locale: "en",
    enrollmentsCount: 0,
    duration: 12,
    coverUrl: null,
    course,
  };
}

// A dashboard row: the enrollment, where the learner stands, and the card to draw.
function dashboardCourse(
  id: number,
  card: CourseCatalogItem,
  state: "started" | "finished",
): DashboardCourse {
  return {
    enrollment: { id, userId: 1, courseId: card.id, state },
    progress: {
      state,
      completion: state === "finished" ? 100 : 40,
      nextLesson: state === "started" ? { id: 5, slug: "variables", name: "Variables" } : null,
    },
    card,
  };
}

function openMy(dashboard: MyDashboard, visitor: User | null = user) {
  worker.use(http.get("*/api/my", () => HttpResponse.json(dashboard)));
  return renderRoute(myRoute, { path: "/{-$locale}/my", initialPath: "/my", user: visitor });
}

test("a guest is sent to sign in first", async () => {
  const { router } = await openMy({ started: [], finished: [] }, null);

  await expect.poll(() => router.state.location.pathname).toBe("/session/new");
});

test("the dashboard shows courses in progress with where to resume, and finished ones", async () => {
  await openMy({
    started: [dashboardCourse(1, catalogItem(1, "javascript", "JavaScript"), "started")],
    finished: [dashboardCourse(2, catalogItem(2, "python", "Python"), "finished")],
  });

  await expect.element(page.getByRole("heading", { name: "JavaScript" })).toBeVisible();
  await expect.element(page.getByText("Variables →")).toBeVisible();
  await expect.element(page.getByText("Python")).toBeVisible();
});

test("a learner with no courses sees an empty state under both headings", async () => {
  await openMy({ started: [], finished: [] });

  await expect.element(page.getByText("Empty").first()).toBeVisible();
  expect(page.getByText("Empty").elements()).toHaveLength(2);
});

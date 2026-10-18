import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
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
import type { RenderedRoute } from "@/test/renderRoute";

const FULL_COMPLETION = 100;
const STARTED_COMPLETION = 40;

const user: User = {
  admin: false,
  assistantMessagesCount: 0,
  canAccessAdmin: false,
  createdAt: "2026-09-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "dora@example.com",
  firstName: "Dora",
  id: 1,
  lastName: null,
  name: "Dora",
  type: "user",
};

function catalogItem(id: number, slug: string, name: string): CourseCatalogItem {
  const course: Course = {
    categoryId: null,
    coverListVariant: null,
    coverThumbVariant: null,
    createdAt: "2026-01-01T00:00:00Z",
    currentVersion: null,
    currentVersionId: null,
    enrollmentsCount: 0,
    hexletProgramLandingPage: null,
    id,
    learnAs: "first_language",
    lessonsCount: 10,
    name,
    ratingCount: 0,
    ratingValue: 0,
    readiness: "completed",
    repositoryUrl: null,
    slug,
  };
  return {
    course,
    coverUrl: null,
    duration: 12,
    enrollmentsCount: 0,
    header: name,
    id,
    locale: "en",
    name,
    slug,
  };
}

// A dashboard row: the enrollment, where the learner stands, and the card to draw.
function dashboardCourse(
  id: number,
  card: CourseCatalogItem,
  state: "started" | "finished",
): DashboardCourse {
  return {
    card,
    enrollment: { courseId: card.id, id, state, userId: 1 },
    progress: {
      completion: state === "finished" ? FULL_COMPLETION : STARTED_COMPLETION,
      nextLesson: state === "started" ? { id: 5, name: "Variables", slug: "variables" } : null,
      state,
    },
  };
}

function openMy(dashboard: MyDashboard, visitor: User | null = user): Promise<RenderedRoute> {
  worker.use(http.get("*/api/my", () => HttpResponse.json(dashboard)));
  return renderRoute(myRoute, { initialPath: "/my", path: "/{-$locale}/my", user: visitor });
}

describe("my", () => {
  it("a guest is sent to sign in first", async () => {
    expect.hasAssertions();
    const { router } = await openMy({ finished: [], started: [] }, null);

    await expect.poll(() => router.state.location.pathname).toBe("/session/new");
  }, 60_000);

  it("the dashboard shows courses in progress with where to resume, and finished ones", async () => {
    expect.hasAssertions();
    await openMy({
      finished: [dashboardCourse(2, catalogItem(2, "python", "Python"), "finished")],
      started: [dashboardCourse(1, catalogItem(1, "javascript", "JavaScript"), "started")],
    });

    await expect.element(page.getByRole("heading", { name: "JavaScript" })).toBeVisible();
    await expect.element(page.getByText("Variables →")).toBeVisible();
    await expect.element(page.getByText("Python")).toBeVisible();
  }, 60_000);

  it("a learner with no courses sees an empty state under both headings", async () => {
    expect.hasAssertions();
    await openMy({ finished: [], started: [] });

    await expect.element(page.getByText("Empty").first()).toBeVisible();
    expect(page.getByText("Empty").elements()).toHaveLength(2);
  }, 60_000);
});

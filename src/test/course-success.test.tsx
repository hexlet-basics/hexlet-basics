import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type {
  Course,
  CourseLandingPage,
  CourseProgress,
  CourseView,
  EnrollmentState,
  User,
} from "@/client/types.gen";
import { Route as courseRoute } from "@/routes/{-$locale}/languages/$slug/index";
import { Route as successRoute } from "@/routes/{-$locale}/languages/$slug/success";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const FULL_COMPLETION = 100;

// The course completion page, driven through its real route with the course
// read faked at the HTTP boundary: who sees it, and where everyone else goes.

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

const course: Course = {
  categoryId: null,
  coverListVariant: null,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  currentVersion: null,
  currentVersionId: 99,
  enrollmentsCount: 1,
  hexletProgramLandingPage: null,
  id: 1,
  learnAs: null,
  lessonsCount: 1,
  name: "javascript",
  ratingCount: 0,
  ratingValue: 0,
  readiness: "completed",
  repositoryUrl: null,
  slug: "javascript",
};

const landingPage: CourseLandingPage = {
  courseId: 1,
  courseSlug: "javascript",
  createdAt: "2026-01-01T00:00:00Z",
  description: "Learn JavaScript in the browser",
  duration: 10,
  enrollmentsCount: 1,
  footer: true,
  footerName: null,
  header: "JavaScript",
  id: 10,
  landingPageToRedirectId: null,
  listed: true,
  main: true,
  metaDescription: "Learn JavaScript",
  metaTitle: "JavaScript course",
  name: "JavaScript",
  order: null,
  outcomesDescription: null,
  outcomesHeader: null,
  outcomesImage: null,
  slug: "javascript",
  state: "published",
  usedInDescription: null,
  usedInHeader: null,
};

// The course read for a learner whose Enrollment is in `state`, or who has none.
function courseView(state: EnrollmentState | null): CourseView {
  const progress: CourseProgress = {
    completion: state === "finished" ? FULL_COMPLETION : 0,
    nextLesson:
      state === "finished" ? null : { id: 1001, name: "Hello, World!", slug: "hello-world" },
    state,
  };
  return {
    course,
    enrollment: state ? { courseId: course.id, id: 5, state, userId: user.id } : null,
    landingPage,
    lessons: [
      {
        description: null,
        finished: state === "finished",
        id: 1001,
        name: "Hello, World!",
        position: 1,
        slug: "hello-world",
      },
    ],
    modules: [],
    movedTo: null,
    progress,
    qnaItems: [],
  };
}

function serveCourse(view: CourseView) {
  worker.use(http.get("*/api/languages/javascript", () => HttpResponse.json(view)));
}

function openSuccess(visitor: User | null) {
  return renderRoute(successRoute, {
    initialPath: "/ru/languages/javascript/success",
    path: "/{-$locale}/languages/$slug/success",
    user: visitor,
  });
}

describe("course success", () => {
  it("a learner who finished the course is congratulated and offered the lead form", async () => {
    expect.hasAssertions();
    serveCourse(courseView("finished"));

    await openSuccess(user);

    await expect
      .element(page.getByRole("heading", { name: "Поздравляем, вы завершили курс «JavaScript»!" }))
      .toBeVisible();
    const career = page.getByRole("link", { name: "Карьера" });
    await expect.element(career).toHaveAttribute("target", "_blank");
    await expect
      .element(career)
      .toHaveAttribute("href", expect.stringContaining("ru.hexlet.io/courses_for_beginners"));
    await expect
      .element(page.getByRole("link", { name: "Навыки и Инструменты" }))
      .toHaveAttribute("href", expect.stringContaining("ru.hexlet.io/courses_for_programmers"));
    await expect.element(page.getByRole("button", { name: "Отправить" })).toBeVisible();
  }, 60_000);

  it("a learner with lessons left is sent to the course page with the warning", async () => {
    expect.hasAssertions();
    serveCourse(courseView("started"));

    const { router } = await openSuccess(user);

    await expect.poll(() => router.state.location.pathname).toBe("/ru/languages/javascript");
    expect({ ...(router.state.location.search as Record<string, unknown>) }).toStrictEqual({
      unfinished: true,
    });
    expect(page.getByRole("button", { name: "Отправить" }).query()).toBeNull();
  }, 60_000);

  it("a signed-in visitor who never enrolled is sent to the course page too", async () => {
    expect.hasAssertions();
    serveCourse(courseView(null));

    const { router } = await openSuccess(user);

    await expect.poll(() => router.state.location.pathname).toBe("/ru/languages/javascript");
    expect({ ...(router.state.location.search as Record<string, unknown>) }).toStrictEqual({
      unfinished: true,
    });
  }, 60_000);

  it("a guest is sent to sign in first", async () => {
    expect.hasAssertions();
    serveCourse(courseView(null));

    const { router } = await openSuccess(null);

    await expect.poll(() => router.state.location.pathname).toBe("/session/new");
  }, 60_000);

  it("the course page carries the legacy warning", async () => {
    expect.hasAssertions();
    serveCourse(courseView("started"));
    await renderRoute(courseRoute, {
      initialPath: "/ru/languages/javascript?unfinished=true",
      path: "/{-$locale}/languages/$slug",
    });

    await expect
      .element(
        page.getByText(
          "В этом курсе ещё есть незавершённые уроки. Пройдите их, чтобы завершить курс!",
        ),
      )
      .toBeVisible();
  }, 60_000);
});

import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type {
  Course,
  CourseLandingPage,
  CourseView,
  EnrollmentState,
  User,
} from "@/client/types.gen";
import { Route as courseRoute } from "@/routes/{-$locale}/languages/$slug/index";
import { Route as successRoute } from "@/routes/{-$locale}/languages/$slug/success";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The course completion page, driven through its real route with the course
// read faked at the HTTP boundary: who sees it, and where everyone else goes.

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

const course: Course = {
  id: 1,
  slug: "javascript",
  name: "javascript",
  learnAs: null,
  readiness: "completed",
  categoryId: null,
  currentVersionId: 99,
  currentVersion: null,
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 1,
  lessonsCount: 1,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: null,
  coverThumbVariant: null,
};

const landingPage: CourseLandingPage = {
  id: 10,
  courseId: 1,
  courseSlug: "javascript",
  createdAt: "2026-01-01T00:00:00Z",
  slug: "javascript",
  name: "JavaScript",
  main: true,
  listed: true,
  state: "published",
  order: null,
  footer: true,
  footerName: null,
  landingPageToRedirectId: null,
  metaTitle: "JavaScript course",
  metaDescription: "Learn JavaScript",
  header: "JavaScript",
  description: "Learn JavaScript in the browser",
  usedInHeader: null,
  usedInDescription: null,
  outcomesHeader: null,
  outcomesDescription: null,
  outcomesImage: null,
  duration: 10,
  enrollmentsCount: 1,
};

// The course read for a learner whose Enrollment is in `state`, or who has none.
function courseView(state: EnrollmentState | null): CourseView {
  const progress = {
    state,
    completion: state === "finished" ? 100 : 0,
    nextLessonSlug: state === "finished" ? null : "hello-world",
    furthestFinishedPosition: state === "finished" ? 1 : 0,
    lessons: [
      { slug: "hello-world", position: 1, finished: state === "finished", available: true },
    ],
  };
  return {
    course,
    landingPage,
    lessons: [{ id: 1001, name: "Hello, World!", description: null, slug: "hello-world" }],
    modules: [],
    qnaItems: [],
    enrollment: state
      ? {
          id: 5,
          userId: user.id,
          courseId: course.id,
          state,
          completion: progress.completion,
          nextLessonName: null,
          progress,
        }
      : null,
    progress,
  };
}

function serveCourse(view: CourseView) {
  worker.use(http.get("*/api/languages/javascript", () => HttpResponse.json(view)));
}

function openSuccess(visitor: User | null) {
  return renderRoute(successRoute, {
    path: "/{-$locale}/languages/$slug/success",
    initialPath: "/ru/languages/javascript/success",
    user: visitor,
  });
}

test("a learner who finished the course is congratulated and offered the lead form", async () => {
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
});

test("a learner with lessons left is sent to the course page with the warning", async () => {
  serveCourse(courseView("started"));

  const { router } = await openSuccess(user);

  await expect.poll(() => router.state.location.pathname).toBe("/ru/languages/javascript");
  expect(router.state.location.search).toEqual({ unfinished: true });
  expect(page.getByRole("button", { name: "Отправить" }).query()).toBeNull();
});

test("a signed-in visitor who never enrolled is sent to the course page too", async () => {
  serveCourse(courseView(null));

  const { router } = await openSuccess(user);

  await expect.poll(() => router.state.location.pathname).toBe("/ru/languages/javascript");
  expect(router.state.location.search).toEqual({ unfinished: true });
});

test("a guest is sent to sign in first", async () => {
  serveCourse(courseView(null));

  const { router } = await openSuccess(null);

  await expect.poll(() => router.state.location.pathname).toBe("/session/new");
});

test("the course page carries the legacy warning", async () => {
  serveCourse(courseView("started"));
  await renderRoute(courseRoute, {
    path: "/{-$locale}/languages/$slug",
    initialPath: "/ru/languages/javascript?unfinished=true",
  });

  await expect
    .element(
      page.getByText(
        "В этом курсе ещё есть незавершённые уроки. Пройдите их, чтобы завершить курс!",
      ),
    )
    .toBeVisible();
});

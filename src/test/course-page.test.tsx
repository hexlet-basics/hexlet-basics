import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { Course, CourseLandingPage, CourseProgress, CourseView } from "@/client/types.gen";
import type { AuthUser } from "@/lib/auth";
import { Route as courseRoute } from "@/routes/{-$locale}/languages/$slug/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The Course page, driven through its real route with the API faked at the HTTP
// boundary. What is asserted is what a visitor sees — the landing copy, the
// marks in the lesson list, the one button — and the requests pressing it made.

const course: Course = {
  id: 1,
  slug: "javascript",
  name: "javascript",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: 99,
  currentVersion: null,
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 1234,
  lessonsCount: 3,
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
  slug: "javascript-ru",
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
  header: "JavaScript course",
  description: "Learn JavaScript in the browser",
  usedInHeader: "Where JavaScript is used",
  usedInDescription: "Browsers, servers and everything in between",
  outcomesHeader: null,
  outcomesDescription: null,
  outcomesImage: null,
  duration: 10,
  enrollmentsCount: 0,
};

// Three payloads the server can send, each a fact about where the visitor
// stands. The page is expected to render them as given.

// Has finished the first lesson: the second is next, the third still locked.
const midCourse: CourseProgress = {
  state: "started",
  completion: 33,
  nextLessonSlug: "variables",
  furthestFinishedPosition: 1,
  lessons: [
    { slug: "hello-world", position: 1, finished: true, available: true },
    { slug: "variables", position: 2, finished: false, available: true },
    { slug: "strings", position: 3, finished: false, available: false },
  ],
};

// Has finished everything: there is no Next Lesson.
const finished: CourseProgress = {
  state: "finished",
  completion: 100,
  nextLessonSlug: null,
  furthestFinishedPosition: 3,
  lessons: [
    { slug: "hello-world", position: 1, finished: true, available: true },
    { slug: "variables", position: 2, finished: true, available: true },
    { slug: "strings", position: 3, finished: true, available: true },
  ],
};

// Has never been here: the first lesson open, the rest locked.
const firstVisit: CourseProgress = {
  state: null,
  completion: 0,
  nextLessonSlug: "hello-world",
  furthestFinishedPosition: 0,
  lessons: [
    { slug: "hello-world", position: 1, finished: false, available: true },
    { slug: "variables", position: 2, finished: false, available: false },
    { slug: "strings", position: 3, finished: false, available: false },
  ],
};

function courseView(progress: CourseProgress): CourseView {
  return {
    course,
    landingPage,
    lessons: [
      { id: 1001, name: "Hello, World!", description: null, slug: "hello-world" },
      { id: 1002, name: "Variables", description: null, slug: "variables" },
      { id: 1003, name: "Strings", description: null, slug: "strings" },
    ],
    enrollment: null,
    progress,
  };
}

const learner: AuthUser = {
  id: 7,
  firstName: "Ada",
  lastName: null,
  name: "Ada",
  email: "ada@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: 0,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

// The real route at a real URL. Start commands are recorded from before the
// first render, together with where the visitor was when each arrived — so a
// test can tell "nothing started on load" and "started, then navigated" apart
// from their opposites.
async function renderCourse(user: AuthUser | null = null, initialPath = "/languages/javascript") {
  const starts: { id: string; pathname: string }[] = [];
  let pathname = () => initialPath;
  worker.use(
    http.post("*/lessons/:id/start", ({ params }) => {
      starts.push({ id: String(params.id), pathname: pathname() });
      return HttpResponse.json(midCourse);
    }),
  );

  const rendered = await renderRoute(courseRoute, {
    path: "/{-$locale}/languages/$slug",
    initialPath,
    user,
  });
  pathname = () => rendered.router.state.location.pathname;

  return { ...rendered, starts };
}

// The lesson entries in the list, in the order a visitor reads them. The marks
// are labelled icons, so they are asserted by accessible name, not here.
function lessonNames() {
  return page
    .getByRole("link")
    .elements()
    .filter((el) => el.getAttribute("href")?.includes("/lessons/"))
    .map((el) => el.textContent?.trim());
}

test("a learner mid-course sees the landing copy, the marked lessons, and Continue", async () => {
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(midCourse))));

  const { router, starts } = await renderCourse(learner);

  await expect.element(page.getByRole("heading", { name: "JavaScript course" })).toBeVisible();
  await expect.element(page.getByText("Learn JavaScript in the browser")).toBeVisible();
  await expect
    .element(page.getByRole("heading", { name: "Where JavaScript is used" }))
    .toBeVisible();

  // Every current lesson, in course order.
  expect(lessonNames()).toEqual(["Hello, World!", "Variables", "Strings"]);
  await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
  await expect.element(page.getByRole("link", { name: "Locked Strings" })).toBeVisible();
  expect(page.getByLabelText("Finished").elements()).toHaveLength(1);
  expect(page.getByLabelText("Locked").elements()).toHaveLength(1);
  // Locked stays reachable: theory is public.
  await expect
    .element(page.getByRole("link", { name: "Locked Strings" }))
    .toHaveAttribute("href", "/languages/javascript/lessons/strings");

  // Completion as the server reported it.
  await expect.element(page.getByText("Completed 33%")).toBeVisible();

  // Nothing started by merely loading the page.
  expect(starts).toEqual([]);

  // Continue enters the Next Lesson — started first, then navigated to.
  await page.getByRole("button", { name: "Continue Learning" }).click();
  await expect
    .poll(() => router.state.location.pathname)
    .toBe("/languages/javascript/lessons/variables");
  expect(starts).toEqual([{ id: "1002", pathname: "/languages/javascript" }]);
});

test("a finished learner is told so and offered no lesson", async () => {
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(finished))));

  await renderCourse(learner);

  await expect.element(page.getByText("Course finished!")).toBeVisible();
  await expect.element(page.getByText("Completed 100%")).toBeVisible();
  expect(page.getByLabelText("Finished").elements()).toHaveLength(3);
  expect(page.getByRole("button", { name: "Continue Learning" }).elements()).toHaveLength(0);
  expect(page.getByRole("button", { name: "Try It" }).elements()).toHaveLength(0);
});

test("a guest carrying progress sees the same page and the same action", async () => {
  // The same payload a signed-in learner gets: the guest's position is the
  // cookie, and the page has no reason to know the difference.
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(midCourse))));

  const { router, starts } = await renderCourse(null);

  await expect.element(page.getByText("Completed 33%")).toBeVisible();
  expect(lessonNames()).toEqual(["Hello, World!", "Variables", "Strings"]);
  await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
  await expect.element(page.getByRole("link", { name: "Locked Strings" })).toBeVisible();
  expect(page.getByLabelText("Finished").elements()).toHaveLength(1);
  expect(page.getByLabelText("Locked").elements()).toHaveLength(1);

  await page.getByRole("button", { name: "Continue Learning" }).click();
  await expect
    .poll(() => router.state.location.pathname)
    .toBe("/languages/javascript/lessons/variables");
  expect(starts).toEqual([{ id: "1002", pathname: "/languages/javascript" }]);
});

test("a first-time visitor sees the first lesson open and the rest locked", async () => {
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(firstVisit))));

  const { router, starts } = await renderCourse(null);

  await expect.element(page.getByText("Completed 0%")).toBeVisible();
  expect(page.getByLabelText("Finished").elements()).toHaveLength(0);
  expect(page.getByLabelText("Locked").elements()).toHaveLength(2);
  expect(lessonNames()).toEqual(["Hello, World!", "Variables", "Strings"]);
  await expect.element(page.getByRole("link", { name: "Locked Variables" })).toBeVisible();
  await expect.element(page.getByRole("link", { name: "Locked Strings" })).toBeVisible();

  await page.getByRole("button", { name: "Try It" }).click();
  await expect
    .poll(() => router.state.location.pathname)
    .toBe("/languages/javascript/lessons/hello-world");
  expect(starts).toEqual([{ id: "1001", pathname: "/languages/javascript" }]);
});

test("keeps the URL locale when entering the course", async () => {
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(firstVisit))));

  const { router } = await renderCourse(null, "/ru/languages/javascript");

  await page.getByRole("button").first().click();
  await expect
    .poll(() => router.state.location.pathname)
    .toBe("/ru/languages/javascript/lessons/hello-world");
});

test("stays on the course page when the start command fails", async () => {
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(firstVisit))));

  const { router } = await renderCourse(null);
  // Registered after the recorder, so it is the handler that answers.
  worker.use(http.post("*/lessons/:id/start", () => HttpResponse.error()));

  await page.getByRole("button", { name: "Try It" }).click();
  await expect
    .element(page.getByText("There was a network problem", { exact: false }))
    .toBeVisible();
  expect(router.state.location.pathname).toBe("/languages/javascript");
});

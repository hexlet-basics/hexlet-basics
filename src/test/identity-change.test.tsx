import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { Course, CourseView, User } from "@/client/types.gen";
import { Route as courseRoute } from "@/routes/{-$locale}/languages/$slug/index";
import { Route as lessonRoute } from "@/routes/{-$locale}/languages/$slug/lessons/$lessonSlug";
import { Route as consumeRoute } from "@/routes/{-$locale}/magic_links/$token";
import { Route as resetRoute } from "@/routes/{-$locale}/password/$token/edit";
import { Route as signInRoute } from "@/routes/{-$locale}/session/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// What a page cached for one visitor shows the next. The course page carries the
// visitor's progress, so a copy read as a guest must not outlive the sign-in (or
// the sign-out) that changes who is looking. Each test caches the page as the
// previous visitor, runs one flow, and opens the page again.

const ada: User = {
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
  lessonsCount: 2,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: null,
  coverThumbVariant: null,
};

// The same course as seen by whoever is asking: Ada has finished a third of it,
// a guest has not started.
function courseFor(visitor: "ada" | "guest"): CourseView {
  const lessons = [
    {
      id: 1001,
      name: "Hello, World!",
      description: null,
      slug: "hello-world",
      position: 1,
      finished: visitor === "ada",
    },
    {
      id: 1002,
      name: "Variables",
      description: null,
      slug: "variables",
      position: 2,
      finished: false,
    },
  ];
  return {
    course,
    landingPage: null,
    movedTo: null,
    lessons,
    modules: [],
    qnaItems: [],
    enrollment:
      visitor === "ada" ? { id: 1, userId: ada.id, courseId: course.id, state: "started" } : null,
    progress:
      visitor === "ada"
        ? {
            state: "started",
            completion: 33,
            nextLesson: { id: 1002, slug: "variables", name: "Variables" },
          }
        : {
            state: null,
            completion: 0,
            nextLesson: { id: 1001, slug: "hello-world", name: "Hello, World!" },
          },
  };
}

// The API as a server that remembers who is signed in.
function serverWith(signedIn: boolean) {
  const state = { signedIn };
  worker.use(
    http.get("*/languages/javascript", () =>
      HttpResponse.json(courseFor(state.signedIn ? "ada" : "guest")),
    ),
    http.get("*/me", () => HttpResponse.json({ user: state.signedIn ? ada : null })),
  );
  return state;
}

// The course page, read and on screen as `visitor` sees it, beside the route a
// flow runs on.
async function courseOpenedAs(
  visitor: "ada" | "guest",
  also: { route: never; path: string }[] = [],
) {
  const rendered = await renderRoute(courseRoute, {
    path: "/{-$locale}/languages/$slug",
    initialPath: "/languages/javascript",
    user: visitor === "ada" ? ada : null,
    also,
  });
  const onScreen =
    visitor === "ada"
      ? page.getByText("Completed 33%")
      : page.getByRole("button", { name: "Try It" });
  await expect.element(onScreen).toBeVisible();
  return rendered;
}

test("a course page cached as a guest shows the learner's progress after signing in", async () => {
  const server = serverWith(false);
  worker.use(
    http.post("*/session", () => {
      server.signedIn = true;
      return HttpResponse.json(ada);
    }),
  );
  const { router } = await courseOpenedAs("guest", [
    { route: signInRoute as never, path: "/{-$locale}/session/new" },
  ]);

  await router.navigate({ href: "/session/new" });
  await page.getByLabelText("Email", { exact: false }).fill("ada@example.com");
  await page.getByLabelText("Password", { exact: false }).fill("correct-horse-battery");
  await page.getByRole("button", { name: "Sign In" }).click();
  await expect.poll(() => router.state.location.pathname).toBe("/");

  await router.navigate({ href: "/languages/javascript" });
  await expect.element(page.getByText("Completed 33%")).toBeVisible();
});

test("a course page cached as a guest shows the learner's progress after a Magic Link", async () => {
  const server = serverWith(false);
  worker.use(
    http.get("*/magic_links/good", () => {
      server.signedIn = true;
      return HttpResponse.json(ada);
    }),
  );
  const { router } = await courseOpenedAs("guest", [
    { route: consumeRoute as never, path: "/{-$locale}/magic_links/$token" },
  ]);

  await router.navigate({ href: "/magic_links/good" });
  await expect.poll(() => router.state.location.pathname).toBe("/");

  await router.navigate({ href: "/languages/javascript" });
  await expect.element(page.getByText("Completed 33%")).toBeVisible();
});

test("a course page cached as a guest shows the learner's progress after a password reset", async () => {
  const server = serverWith(false);
  worker.use(
    http.get("*/password/good/edit", () => new HttpResponse(null, { status: 204 })),
    http.patch("*/password/good", () => {
      server.signedIn = true;
      return HttpResponse.json(ada);
    }),
  );
  const { router } = await courseOpenedAs("guest", [
    { route: resetRoute as never, path: "/{-$locale}/password/$token/edit" },
  ]);

  await router.navigate({ href: "/password/good/edit" });
  await page.getByLabelText("Password", { exact: false }).fill("new-password");
  await page.getByRole("button", { name: "Save password" }).click();
  await expect.poll(() => router.state.location.pathname).toBe("/");

  await router.navigate({ href: "/languages/javascript" });
  await expect.element(page.getByText("Completed 33%")).toBeVisible();
});

test("a course page cached as a learner shows a guest's page after signing out", async () => {
  const server = serverWith(true);
  worker.use(
    http.delete("*/session", () => {
      server.signedIn = false;
      return new HttpResponse(null, { status: 204 });
    }),
  );
  const { router } = await courseOpenedAs("ada");

  // The test viewport is a phone's: the account menu is behind the burger.
  await page.getByRole("button", { name: "Toggle Main Menu Visibility" }).click();
  await page.getByRole("button", { name: "Ada" }).click();
  await page.getByRole("menuitem", { name: "Sign out" }).click();
  await expect.poll(() => router.state.location.pathname).toBe("/");

  await router.navigate({ href: "/languages/javascript" });
  await expect.element(page.getByRole("button", { name: "Try It" })).toBeVisible();
  expect(page.getByText("Completed 33%").elements()).toHaveLength(0);
});

test("pointing at a lesson preloads it and starts nothing", async () => {
  serverWith(true);
  const starts: string[] = [];
  let preloads = 0;
  worker.use(
    http.post("*/lessons/:id/start", ({ params }) => {
      starts.push(String(params.id));
      return new HttpResponse(null, { status: 204 });
    }),
    http.get("*/languages/javascript/lessons/variables", () => {
      preloads += 1;
      return HttpResponse.json({ status: 500 }, { status: 500 });
    }),
  );
  await courseOpenedAs("ada", [
    { route: lessonRoute as never, path: "/{-$locale}/languages/$slug/lessons/$lessonSlug" },
  ]);

  await page.getByRole("link", { name: "Variables" }).hover();

  await expect.poll(() => preloads).toBeGreaterThan(0);
  expect(starts).toEqual([]);
});

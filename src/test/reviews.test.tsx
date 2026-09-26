import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { Course, Review, ReviewPage, User } from "@/client/types.gen";
import { Route as indexRoute } from "@/routes/{-$locale}/reviews/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The public reviews page, driven through its real route with the API faked at
// the HTTP boundary. Asserted is what a visitor sees: the reviews, the page it
// asked for, and the pager's links.

const user: User = {
  id: 1,
  firstName: "Alice",
  lastName: null,
  name: "Alice",
  email: "alice@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: null,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

const course: Course = {
  id: 1,
  slug: "ruby",
  name: "ruby",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: null,
  currentVersion: null,
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 0,
  lessonsCount: 3,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: null,
  coverThumbVariant: null,
};

function review(overrides: Partial<Review>): Review {
  return {
    id: 1,
    user,
    course,
    userId: user.id,
    courseId: course.id,
    body: "Clear lessons and quick practice",
    firstName: "Jane",
    lastName: "Smith",
    fullName: "Jane Smith",
    state: "published",
    locale: "en",
    pinned: false,
    createdAt: "2026-07-29T10:00:00Z",
    ...overrides,
  };
}

test("the reviews page lists the reviews the API returns, with a pager", async () => {
  const list: ReviewPage = {
    items: [review({}), review({ id: 2, fullName: "John Doe", body: "Great course" })],
    total: 45,
    page: 1,
    perPage: 20,
  };
  worker.use(http.get("*/api/reviews", () => HttpResponse.json(list)));

  const { router } = await renderRoute(indexRoute, {
    path: "/{-$locale}/reviews/",
    initialPath: "/reviews",
  });

  await expect.element(page.getByRole("heading", { name: "Reviews", level: 1 })).toBeVisible();
  await expect.element(page.getByText("Jane Smith")).toBeVisible();
  await expect.element(page.getByText("Great course")).toBeVisible();
  await expect.element(page.getByText("July 29, 2026").first()).toBeVisible();

  // 45 reviews at 20 a page are three pages, each a link to its `?page=`.
  await expect
    .element(page.getByRole("link", { name: "2", exact: true }))
    .toHaveAttribute("href", "/reviews?page=2");
  await expect
    .element(page.getByRole("link", { name: "3", exact: true }))
    .toHaveAttribute("href", "/reviews?page=3");

  // The route's head: legacy title behind the site name, and its Twitter card.
  const meta = router.state.matches.at(-1)?.meta;
  expect(meta).toContainEqual({ title: "CodeBasics | Reviews" });
  expect(meta).toContainEqual({ name: "twitter:site", content: "@hexlet_io" });
});

test("a page in the URL is the page asked of the API", async () => {
  let asked: string | null = null;
  worker.use(
    http.get("*/api/reviews", ({ request }) => {
      asked = new URL(request.url).searchParams.get("page");
      const list: ReviewPage = {
        items: [review({ fullName: "Ivan Petrov" })],
        total: 21,
        page: 2,
        perPage: 20,
      };
      return HttpResponse.json(list);
    }),
  );

  await renderRoute(indexRoute, {
    path: "/{-$locale}/reviews/",
    initialPath: "/reviews?page=2",
  });

  await expect.element(page.getByText("Ivan Petrov")).toBeVisible();
  expect(asked).toBe("2");
  // Back to page one is the bare URL.
  await expect
    .element(page.getByRole("link", { name: "1", exact: true }))
    .toHaveAttribute("href", "/reviews");
});

test("the pager keeps the URL's locale prefix", async () => {
  const list: ReviewPage = { items: [review({})], total: 21, page: 1, perPage: 20 };
  worker.use(http.get("*/api/reviews", () => HttpResponse.json(list)));

  await renderRoute(indexRoute, {
    path: "/{-$locale}/reviews/",
    initialPath: "/ru/reviews",
  });

  await expect
    .element(page.getByRole("link", { name: "2", exact: true }))
    .toHaveAttribute("href", "/ru/reviews?page=2");
});

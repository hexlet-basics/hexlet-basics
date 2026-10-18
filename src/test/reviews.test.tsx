import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { Course, Review, ReviewPage, User } from "@/client/types.gen";
import { Route as indexRoute } from "@/routes/{-$locale}/reviews/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The public reviews page, driven through its real route with the API faked at
// the HTTP boundary. Asserted is what a visitor sees: the reviews, the page it
// asked for, and the pager's links.

const user: User = {
  admin: false,
  assistantMessagesCount: null,
  canAccessAdmin: false,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "alice@example.com",
  firstName: "Alice",
  id: 1,
  lastName: null,
  name: "Alice",
  type: "user",
};

const course: Course = {
  categoryId: null,
  coverListVariant: null,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  currentVersion: null,
  currentVersionId: null,
  enrollmentsCount: 0,
  hexletProgramLandingPage: null,
  id: 1,
  learnAs: null,
  lessonsCount: 3,
  name: "ruby",
  ratingCount: 0,
  ratingValue: 0,
  readiness: null,
  repositoryUrl: null,
  slug: "ruby",
};

function review(overrides: Partial<Review>): Review {
  return {
    body: "Clear lessons and quick practice",
    course,
    courseId: course.id,
    createdAt: "2026-07-29T10:00:00Z",
    firstName: "Jane",
    fullName: "Jane Smith",
    id: 1,
    lastName: "Smith",
    locale: "en",
    pinned: false,
    state: "published",
    user,
    userId: user.id,
    ...overrides,
  };
}

describe("reviews", () => {
  it("the reviews page lists the reviews the API returns, with a pager", async () => {
    expect.hasAssertions();
    const list: ReviewPage = {
      items: [review({}), review({ body: "Great course", fullName: "John Doe", id: 2 })],
      page: 1,
      perPage: 20,
      total: 45,
    };
    worker.use(http.get("*/api/reviews", () => HttpResponse.json(list)));

    const { router } = await renderRoute(indexRoute, {
      initialPath: "/reviews",
      path: "/{-$locale}/reviews/",
    });

    await expect.element(page.getByRole("heading", { level: 1, name: "Reviews" })).toBeVisible();
    await expect.element(page.getByText("Jane Smith")).toBeVisible();
    await expect.element(page.getByText("Great course")).toBeVisible();
    await expect.element(page.getByText("July 29, 2026").first()).toBeVisible();

    // 45 reviews at 20 a page are three pages, each a link to its `?page=`.
    await expect
      .element(page.getByRole("link", { exact: true, name: "2" }))
      .toHaveAttribute("href", "/reviews?page=2");
    await expect
      .element(page.getByRole("link", { exact: true, name: "3" }))
      .toHaveAttribute("href", "/reviews?page=3");

    // The route's head: legacy title behind the site name, and its Twitter card.
    const meta = router.state.matches.at(-1)?.meta;
    expect(meta).toContainEqual({ title: "CodeBasics | Reviews" });
    expect(meta).toContainEqual({ content: "@hexlet_io", name: "twitter:site" });
    // The canonical is the list itself, whatever page the visitor is on.
    expect(router.state.matches.at(-1)?.links).toContainEqual({
      href: `${globalThis.location.origin}/reviews`,
      rel: "canonical",
    });
  }, 60_000);

  it("a page in the URL is the page asked of the API", async () => {
    expect.hasAssertions();
    let asked: string | null = null;
    worker.use(
      http.get("*/api/reviews", ({ request }) => {
        asked = new URL(request.url).searchParams.get("page");
        const list: ReviewPage = {
          items: [review({ fullName: "Ivan Petrov" })],
          page: 2,
          perPage: 20,
          total: 21,
        };
        return HttpResponse.json(list);
      }),
    );

    await renderRoute(indexRoute, {
      initialPath: "/reviews?page=2",
      path: "/{-$locale}/reviews/",
    });

    await expect.element(page.getByText("Ivan Petrov")).toBeVisible();
    expect(asked).toBe("2");
    // Back to page one is the bare URL.
    await expect
      .element(page.getByRole("link", { exact: true, name: "1" }))
      .toHaveAttribute("href", "/reviews");
  }, 60_000);

  it("the pager keeps the URL's locale prefix", async () => {
    expect.hasAssertions();
    const list: ReviewPage = { items: [review({})], page: 1, perPage: 20, total: 21 };
    worker.use(http.get("*/api/reviews", () => HttpResponse.json(list)));

    await renderRoute(indexRoute, {
      initialPath: "/ru/reviews",
      path: "/{-$locale}/reviews/",
    });

    await expect
      .element(page.getByRole("link", { exact: true, name: "2" }))
      .toHaveAttribute("href", "/ru/reviews?page=2");
  }, 60_000);
});

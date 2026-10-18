import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type {
  Course,
  CourseCatalogItem,
  CourseCategory,
  CourseCategoryView,
} from "@/client/types.gen";
import { Route as showRoute } from "@/routes/{-$locale}/language_categories/$slug";
import { Route as indexRoute } from "@/routes/{-$locale}/language_categories/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The public course categories, driven through their real routes with the API
// faked at the HTTP boundary. Asserted is what a visitor sees: the categories,
// a category's courses and questions
// and the page's head.

const programming: CourseCategory = {
  createdAt: "2026-01-01T00:00:00Z",
  description: "Learn to program from scratch",
  header: "Programming",
  id: 1,
  locale: "en",
  name: "Programming",
  slug: "programming",
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

const rubyLanding: CourseCatalogItem = {
  course,
  coverUrl: null,
  duration: 3,
  enrollmentsCount: 42,
  header: "Ruby course",
  id: 100,
  locale: "en",
  name: "Ruby",
  slug: "ruby-en",
};

const view: CourseCategoryView = {
  category: programming,
  landingPages: [rubyLanding],
  qnaItems: [{ answer: "No, **none**.", id: 1, question: "Is experience required?" }],
};

function serveCategory(): void {
  worker.use(http.get("*/api/language_categories/programming", () => HttpResponse.json(view)));
}

describe("course categories", () => {
  it("the index links every category the API returns", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/api/language_categories", () => HttpResponse.json([programming])));

    const { router } = await renderRoute(indexRoute, {
      initialPath: "/language_categories/",
      path: "/{-$locale}/language_categories/",
    });

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Course categories" }))
      .toBeVisible();
    // The footer lists the same categories, after the page's own card.
    const category = page.getByRole("link", { name: /Programming/u }).first();
    await expect.element(category).toHaveAttribute("href", "/language_categories/programming");

    const meta = router.state.matches.at(-1)?.meta;
    expect(meta).toContainEqual({ title: "CodeBasics | Course categories" });
    expect(meta).toContainEqual({ content: "@hexlet_io", name: "twitter:site" });
    expect(router.state.matches.at(-1)?.links).toContainEqual({
      href: `${globalThis.location.origin}/language_categories`,
      rel: "canonical",
    });
  }, 60_000);

  it("a category shows its courses and questions, with the legacy head", async () => {
    expect.hasAssertions();
    serveCategory();

    const { router } = await renderRoute(showRoute, {
      initialPath: "/language_categories/programming",
      path: "/{-$locale}/language_categories/$slug",
    });

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Programming" }))
      .toBeVisible();
    await expect.element(page.getByText("Learn to program from scratch")).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Ruby" })).toBeVisible();
    await expect.element(page.getByText("Is experience required?")).toBeVisible();
    await expect.element(page.getByText("none", { exact: true })).toBeVisible();

    const match = router.state.matches.at(-1);
    expect(match?.meta).toContainEqual({
      title: "CodeBasics | Courses in the Programming category",
    });
    const description: unknown = expect.stringContaining("in the Programming category");
    expect(match?.meta).toContainEqual({
      content: description,
      name: "description",
    });
    expect(match?.links).toContainEqual({
      href: `${globalThis.location.origin}/language_categories/programming`,
      rel: "canonical",
    });
  }, 60_000);

  it("a category that is not in this locale is not found", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/api/language_categories/programming", () =>
        HttpResponse.json({ status: 404, title: "Not Found" }, { status: 404 }),
      ),
    );

    await renderRoute(showRoute, {
      initialPath: "/language_categories/programming",
      path: "/{-$locale}/language_categories/$slug",
    });

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);
});

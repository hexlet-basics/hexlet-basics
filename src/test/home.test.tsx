import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { BlogPost, Course, CourseCatalogItem } from "@/client/types.gen";
import { Route as homeRoute } from "@/routes/{-$locale}/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The home page through its real route, inside the real site layout: the hero,
// the course grid a visitor picks a course from, the latest posts and the
// registration form.

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

const courses = [catalogItem(1, "javascript", "JavaScript"), catalogItem(2, "python", "Python")];

const blogPost: BlogPost = {
  id: 10,
  creator: { id: 1, firstName: "Alice", lastName: null, name: "Alice" },
  name: "Hello world",
  slug: "hello-world",
  description: "The very first post",
  state: "published",
  locale: "en",
  richBodyHtml: "<p>Hello</p>",
  readingTime: 0,
  likesCount: 2,
  relatedCourseIds: [],
  coverUrl: null,
  createdAt: "2026-01-01T00:00:00Z",
};

function serveHome() {
  worker.use(
    http.get("*/languages", () => HttpResponse.json(courses)),
    http.get("*/api/blog_posts", () =>
      HttpResponse.json({ items: [blogPost], total: 1, page: 1, perPage: 3 }),
    ),
  );
}

describe("home", () => {
  it("greets a visitor with the hero and the course grid, each course linking to its page", async () => {
    expect.hasAssertions();
    serveHome();

    await renderRoute(homeRoute, { path: "/{-$locale}/", initialPath: "/" });

    const heading = page.getByRole("heading", { level: 1 });
    await expect.element(heading).toBeVisible();

    const main = page.getByRole("main");
    const javascript = main.getByRole("link", { name: /JavaScript/ }).first();
    await expect.element(javascript).toHaveAttribute("href", "/languages/javascript");
    const python = main.getByRole("link", { name: /Python/ }).first();
    await expect.element(python).toHaveAttribute("href", "/languages/python");
  }, 60_000);

  it("shows the latest blog posts and the registration form to a visitor", async () => {
    expect.hasAssertions();
    serveHome();

    await renderRoute(homeRoute, { path: "/{-$locale}/", initialPath: "/" });

    const post = page.getByRole("link", { name: /Hello world/ });
    await expect.element(post.first()).toHaveAttribute("href", "/blog_posts/hello-world");
    const email = page.getByRole("textbox", { name: /Email/ });
    await expect.element(email.first()).toBeVisible();
  }, 60_000);

  it("answers a page that cannot be loaded with the server error page", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/languages", () => HttpResponse.json({ status: 500 }, { status: 500 })));

    await renderRoute(homeRoute, { path: "/{-$locale}/", initialPath: "/" });

    const failure = page.getByRole("heading", { name: "Server Error" });
    await expect.element(failure).toBeVisible();
  }, 60_000);
});

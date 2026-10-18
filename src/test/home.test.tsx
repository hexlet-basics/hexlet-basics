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

const courses = [catalogItem(1, "javascript", "JavaScript"), catalogItem(2, "python", "Python")];

const blogPost: BlogPost = {
  coverUrl: null,
  createdAt: "2026-01-01T00:00:00Z",
  creator: { firstName: "Alice", id: 1, lastName: null, name: "Alice" },
  description: "The very first post",
  id: 10,
  likesCount: 2,
  locale: "en",
  name: "Hello world",
  readingTime: 0,
  relatedCourseIds: [],
  richBodyHtml: "<p>Hello</p>",
  slug: "hello-world",
  state: "published",
};

function serveHome() {
  worker.use(
    http.get("*/languages", () => HttpResponse.json(courses)),
    http.get("*/api/blog_posts", () =>
      HttpResponse.json({ items: [blogPost], page: 1, perPage: 3, total: 1 }),
    ),
  );
}

describe("home", () => {
  it("greets a visitor with the hero and the course grid, each course linking to its page", async () => {
    expect.hasAssertions();
    serveHome();

    await renderRoute(homeRoute, { initialPath: "/", path: "/{-$locale}/" });

    const heading = page.getByRole("heading", { level: 1 });
    await expect.element(heading).toBeVisible();

    const main = page.getByRole("main");
    const javascript = main.getByRole("link", { name: /JavaScript/u }).first();
    await expect.element(javascript).toHaveAttribute("href", "/languages/javascript");
    const python = main.getByRole("link", { name: /Python/u }).first();
    await expect.element(python).toHaveAttribute("href", "/languages/python");
  }, 60_000);

  it("shows the latest blog posts and the registration form to a visitor", async () => {
    expect.hasAssertions();
    serveHome();

    await renderRoute(homeRoute, { initialPath: "/", path: "/{-$locale}/" });

    const post = page.getByRole("link", { name: /Hello world/u });
    await expect.element(post.first()).toHaveAttribute("href", "/blog_posts/hello-world");
    const email = page.getByRole("textbox", { name: /Email/u });
    await expect.element(email.first()).toBeVisible();
  }, 60_000);

  it("answers a page that cannot be loaded with the server error page", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/languages", () => HttpResponse.json({ status: 500 }, { status: 500 })));

    await renderRoute(homeRoute, { initialPath: "/", path: "/{-$locale}/" });

    const failure = page.getByRole("heading", { name: "Server Error" });
    await expect.element(failure).toBeVisible();
  }, 60_000);
});

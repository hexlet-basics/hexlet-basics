import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { BlogPost, BlogPostState } from "@/client/types.gen";
import { Route as listRoute } from "@/routes/{-$locale}/admin/blog_posts/index";
import { Route as editRoute } from "@/routes/{-$locale}/admin/blog_posts/$id";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The Publication State of a blog post is changed by action buttons, one per
// operation, never by the form. Through the real routes under the real admin
// guard: the buttons offered follow the current state, a click calls that
// operation and the post shown is the one it answered with.

function post(state: BlogPostState): BlogPost {
  return {
    id: 7,
    creator: { id: 1, firstName: "Alice", lastName: null, name: "Alice" },
    name: "Hello",
    slug: "hello",
    description: null,
    state,
    locale: "en",
    richBodyHtml: "<p>x</p>",
    readingTime: 0,
    likesCount: 0,
    relatedCourseIds: [],
    coverUrl: null,
    createdAt: "2026-01-01T00:00:00Z",
  };
}

function openPost() {
  return renderRoute(editRoute, {
    path: "/{-$locale}/admin/blog_posts/$id",
    initialPath: "/admin/blog_posts/7",
    user: admin,
  });
}

const noCourses = http.get("*/admin/courses", () =>
  HttpResponse.json({ items: [], total: 0, page: 1, perPage: 100 }),
);

describe("admin blog post state", () => {
  it.each<[BlogPostState, string[], string[]]>([
    ["draft", ["Publish", "Archive"], ["Unpublish"]],
    ["published", ["Unpublish", "Archive"], ["Publish"]],
    ["archived", ["Publish"], ["Unpublish", "Archive"]],
  ])(
    "a %s post offers only the transitions legal for it",
    async (state, offered, hidden) => {
      expect.hasAssertions();
      worker.use(
        noCourses,
        http.get("*/admin/blog_posts/7", () => HttpResponse.json(post(state))),
      );

      await openPost();

      await Promise.all([
        ...offered.map(async (name) => {
          await expect.element(page.getByRole("button", { name, exact: true })).toBeVisible();
        }),
        ...hidden.map(async (name) => {
          await expect
            .element(page.getByRole("button", { name, exact: true }))
            .not.toBeInTheDocument();
        }),
      ]);
    },
    60_000,
  );

  it("the form has no state field", async () => {
    expect.hasAssertions();
    worker.use(
      noCourses,
      http.get("*/admin/blog_posts/7", () => HttpResponse.json(post("draft"))),
    );

    await openPost();

    await expect.element(page.getByLabelText("Name *")).toBeVisible();
    await expect.element(page.getByLabelText("State")).not.toBeInTheDocument();
  }, 60_000);

  it("publishing calls the publish operation and shows the new state", async () => {
    expect.hasAssertions();
    let state: BlogPostState = "draft";
    let published = false;
    worker.use(
      noCourses,
      http.get("*/admin/blog_posts/7", () => HttpResponse.json(post(state))),
      http.post("*/admin/blog_posts/7/publish", () => {
        published = true;
        state = "published";
        return HttpResponse.json(post(state));
      }),
    );

    await openPost();
    await page.getByRole("button", { name: "Publish", exact: true }).click();

    await expect.element(page.getByText("State changed")).toBeVisible();
    expect(published).toBe(true);
    await expect
      .element(page.getByRole("button", { name: "Unpublish", exact: true }))
      .toBeVisible();
    await expect
      .element(page.getByRole("button", { name: "Publish", exact: true }))
      .not.toBeInTheDocument();
  }, 60_000);

  it("archiving and unpublishing call their own operations", async () => {
    expect.hasAssertions();
    const called: string[] = [];
    let state: BlogPostState = "published";
    worker.use(
      noCourses,
      http.get("*/admin/blog_posts/7", () => HttpResponse.json(post(state))),
      http.post("*/admin/blog_posts/7/unpublish", () => {
        called.push("unpublish");
        state = "draft";
        return HttpResponse.json(post(state));
      }),
      http.post("*/admin/blog_posts/7/archive", () => {
        called.push("archive");
        state = "archived";
        return HttpResponse.json(post(state));
      }),
    );

    await openPost();
    await page.getByRole("button", { name: "Unpublish", exact: true }).click();
    await expect.element(page.getByRole("button", { name: "Publish", exact: true })).toBeVisible();
    await page.getByRole("button", { name: "Archive", exact: true }).click();
    await expect
      .element(page.getByRole("button", { name: "Unpublish", exact: true }))
      .not.toBeInTheDocument();

    expect(called).toStrictEqual(["unpublish", "archive"]);
  }, 60_000);

  it("a refused transition tells the admin and keeps the state", async () => {
    expect.hasAssertions();
    worker.use(
      noCourses,
      http.get("*/admin/blog_posts/7", () => HttpResponse.json(post("draft"))),
      http.post("*/admin/blog_posts/7/publish", () =>
        HttpResponse.json({ type: "about:blank", title: "Conflict", status: 409 }, { status: 409 }),
      ),
    );

    await openPost();
    await page.getByRole("button", { name: "Publish", exact: true }).click();

    await expect.element(page.getByText("Could not change the state")).toBeVisible();
    await expect.element(page.getByRole("button", { name: "Publish", exact: true })).toBeVisible();
  }, 60_000);

  it("the list shows each post's state", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/blog_posts", () =>
        HttpResponse.json({ items: [post("archived")], total: 1, page: 1, perPage: 25 }),
      ),
    );

    await renderRoute(listRoute, {
      path: "/{-$locale}/admin/blog_posts/",
      initialPath: "/admin/blog_posts",
      user: admin,
    });

    await expect.element(page.getByRole("cell", { name: "Archived" })).toBeVisible();
  }, 60_000);
});

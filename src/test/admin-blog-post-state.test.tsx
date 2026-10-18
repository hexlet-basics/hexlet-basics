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
    coverUrl: null,
    createdAt: "2026-01-01T00:00:00Z",
    creator: { firstName: "Alice", id: 1, lastName: null, name: "Alice" },
    description: null,
    id: 7,
    likesCount: 0,
    locale: "en",
    name: "Hello",
    readingTime: 0,
    relatedCourseIds: [],
    richBodyHtml: "<p>x</p>",
    slug: "hello",
    state,
  };
}

function openPost() {
  return renderRoute(editRoute, {
    initialPath: "/admin/blog_posts/7",
    path: "/{-$locale}/admin/blog_posts/$id",
    user: admin,
  });
}

const noCourses = http.get("*/admin/courses", () =>
  HttpResponse.json({ items: [], page: 1, perPage: 100, total: 0 }),
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
          await expect.element(page.getByRole("button", { exact: true, name })).toBeVisible();
        }),
        ...hidden.map(async (name) => {
          await expect
            .element(page.getByRole("button", { exact: true, name }))
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
    await page.getByRole("button", { exact: true, name: "Publish" }).click();

    await expect.element(page.getByText("State changed")).toBeVisible();
    expect(published).toBe(true);
    await expect
      .element(page.getByRole("button", { exact: true, name: "Unpublish" }))
      .toBeVisible();
    await expect
      .element(page.getByRole("button", { exact: true, name: "Publish" }))
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
    await page.getByRole("button", { exact: true, name: "Unpublish" }).click();
    await expect.element(page.getByRole("button", { exact: true, name: "Publish" })).toBeVisible();
    await page.getByRole("button", { exact: true, name: "Archive" }).click();
    await expect
      .element(page.getByRole("button", { exact: true, name: "Unpublish" }))
      .not.toBeInTheDocument();

    expect(called).toStrictEqual(["unpublish", "archive"]);
  }, 60_000);

  it("a refused transition tells the admin and keeps the state", async () => {
    expect.hasAssertions();
    worker.use(
      noCourses,
      http.get("*/admin/blog_posts/7", () => HttpResponse.json(post("draft"))),
      http.post("*/admin/blog_posts/7/publish", () =>
        HttpResponse.json({ status: 409, title: "Conflict", type: "about:blank" }, { status: 409 }),
      ),
    );

    await openPost();
    await page.getByRole("button", { exact: true, name: "Publish" }).click();

    await expect.element(page.getByText("Could not change the state")).toBeVisible();
    await expect.element(page.getByRole("button", { exact: true, name: "Publish" })).toBeVisible();
  }, 60_000);

  it("the list shows each post's state", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/blog_posts", () =>
        HttpResponse.json({ items: [post("archived")], page: 1, perPage: 25, total: 1 }),
      ),
    );

    await renderRoute(listRoute, {
      initialPath: "/admin/blog_posts",
      path: "/{-$locale}/admin/blog_posts/",
      user: admin,
    });

    await expect.element(page.getByRole("cell", { name: "Archived" })).toBeVisible();
  }, 60_000);
});

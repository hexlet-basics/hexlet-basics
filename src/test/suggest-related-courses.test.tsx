import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { BlogPost, CoursePage, User } from "@/client/types.gen";
import { Route as editRoute } from "@/routes/{-$locale}/admin/blog_posts/$id";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The admin's AI pick of a post's related courses (legacy `related_courses`
// member action): the button only enqueues the job, so what the admin sees is
// the confirmation, and what the API sees is the post it was asked about.

const author: User = {
  id: 1,
  firstName: "Alice",
  lastName: null,
  name: "Alice",
  email: "alice@example.com",
  admin: true,
  canAccessAdmin: true,
  assistantMessagesCount: null,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

const post: BlogPost = {
  id: 10,
  creator: author,
  name: "Hello world",
  slug: "hello-world",
  description: "The very first post",
  state: "published",
  locale: "ru",
  url: "https://code-basics.com/ru/blog_posts/hello-world",
  richBodyHtml: "<p>Hello world</p>",
  readingTime: 0,
  likesCount: 0,
  relatedCourseItemsCount: 0,
  relatedCourseIds: [],
  coverThumbVariant: null,
  coverListVariant: null,
  coverMainVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
};

const courses: CoursePage = { items: [], total: 0, page: 1, perPage: 100 };

test("enqueues the related-courses suggestion for the post", async () => {
  let suggestedFor: string | undefined;
  worker.use(
    http.get("*/admin/blog_posts/10", () => HttpResponse.json(post)),
    http.get("*/admin/courses", () => HttpResponse.json(courses)),
    http.post("*/admin/blog_posts/:id/suggest_related_courses", ({ params }) => {
      suggestedFor = String(params.id);
      return new HttpResponse(null, { status: 204 });
    }),
  );

  await renderRoute(editRoute, {
    path: "/{-$locale}/admin/blog_posts/$id",
    initialPath: "/admin/blog_posts/10",
  });

  await page.getByRole("button", { name: "Suggest courses" }).click();

  await expect
    .element(page.getByText("Course suggestion started — reload the page shortly"))
    .toBeVisible();
  expect(suggestedFor).toBe("10");
});

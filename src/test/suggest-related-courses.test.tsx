import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { BlogPost, CourseOption, User } from "@/client/types.gen";
import { Route as editRoute } from "@/routes/{-$locale}/admin/blog_posts/$id";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The admin's AI pick of a post's related courses (legacy `related_courses`
// member action): the button only enqueues the job, so what the admin sees is
// the confirmation, and what the API sees is the post it was asked about.

const author: User = {
  admin: true,
  assistantMessagesCount: null,
  canAccessAdmin: true,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "alice@example.com",
  firstName: "Alice",
  id: 1,
  lastName: null,
  name: "Alice",
  type: "user",
};

const post: BlogPost = {
  coverUrl: null,
  createdAt: "2026-01-01T00:00:00Z",
  creator: author,
  description: "The very first post",
  id: 10,
  likesCount: 0,
  locale: "ru",
  name: "Hello world",
  readingTime: 0,
  relatedCourseIds: [],
  richBodyHtml: "<p>Hello world</p>",
  slug: "hello-world",
  state: "published",
};

const courses: CourseOption[] = [];

describe("suggest related courses", () => {
  it("enqueues the related-courses suggestion for the post", async () => {
    expect.hasAssertions();
    let suggestedFor: string | null = null;
    worker.use(
      http.get("*/admin/blog_posts/10", () => HttpResponse.json(post)),
      http.get("*/admin/courses/options", () => HttpResponse.json(courses)),
      http.post("*/admin/blog_posts/:id/suggest_related_courses", ({ params }) => {
        suggestedFor = String(params.id);
        return new HttpResponse(null, { status: 204 });
      }),
    );

    await renderRoute(editRoute, {
      initialPath: "/admin/blog_posts/10",
      path: "/{-$locale}/admin/blog_posts/$id",
      user: author,
    });

    const suggest = page.getByRole("button", { name: "Suggest courses" });
    await suggest.click();

    await expect
      .element(page.getByText("Course suggestion started — reload the page shortly"))
      .toBeVisible();
    expect(suggestedFor).toBe("10");
  }, 60_000);
});

import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page, userEvent } from "vitest/browser";
import type { BlogPost, CourseOption } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/blog_posts/$id";
import { admin, eva } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const COURSE_ID = 150;

const post: BlogPost = {
  coverUrl: null,
  createdAt: "2026-01-01T00:00:00Z",
  creator: eva,
  description: null,
  id: 3,
  likesCount: 0,
  locale: "en",
  name: "Hello",
  readingTime: 1,
  relatedCourseIds: [],
  richBodyHtml: "<p>Hi</p>",
  slug: "hello",
  state: "draft",
};

// The picker of a post's related courses lists every course the catalog has, not
// a first page of them: a course past the hundredth is still one an admin can pick.
describe("admin blog post courses", () => {
  it("the related-courses picker offers every course, however many there are", async () => {
    expect.hasAssertions();
    const options: CourseOption[] = Array.from({ length: 150 }, (_, index) => ({
      id: index + 1,
      name: `Course ${index + 1}`,
    }));
    let saved: unknown = null;
    worker.use(
      http.get("*/admin/blog_posts/3", () => HttpResponse.json(post)),
      http.get("*/admin/courses/options", () => HttpResponse.json(options)),
      http.post("*/admin/blog_posts/3/related_courses", async ({ request }) => {
        saved = await request.json();
        return HttpResponse.json(post);
      }),
    );

    await renderRoute(Route, {
      initialPath: "/admin/blog_posts/3",
      path: "/{-$locale}/admin/blog_posts/$id",
      user: admin,
    });

    await page.getByPlaceholder("Search…").click();
    await page.getByRole("option", { name: "Course 150" }).click();
    // The dropdown stays open over the button until it is dismissed.
    await userEvent.keyboard("{Escape}");
    await page.getByRole("button", { name: "Save" }).last().click();

    await expect.element(page.getByText("Record updated")).toBeVisible();
    expect(saved).toStrictEqual({ courseIds: [COURSE_ID] });
  }, 60_000);
});

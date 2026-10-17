import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page, userEvent } from "vitest/browser";
import type { BlogPost, CourseOption } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/blog_posts/$id";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const post: BlogPost = {
  id: 3,
  creator: { id: 1, name: "Ada" } as BlogPost["creator"],
  name: "Hello",
  slug: "hello",
  description: null,
  state: "draft",
  locale: "en",
  richBodyHtml: "<p>Hi</p>",
  readingTime: 1,
  likesCount: 0,
  relatedCourseItemsCount: 0,
  relatedCourseIds: [],
  coverThumbVariant: null,
  coverListVariant: null,
  coverMainVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
};

// The picker of a post's related courses lists every course the catalog has, not
// a first page of them: a course past the hundredth is still one an admin can pick.
test("the related-courses picker offers every course, however many there are", async () => {
  const options: CourseOption[] = Array.from({ length: 150 }, (_, index) => ({
    id: index + 1,
    name: `Course ${index + 1}`,
  }));
  let saved: unknown;
  worker.use(
    http.get("*/admin/blog_posts/3", () => HttpResponse.json(post)),
    http.get("*/admin/courses/options", () => HttpResponse.json(options)),
    http.post("*/admin/blog_posts/3/related_courses", async ({ request }) => {
      saved = await request.json();
      return HttpResponse.json(post);
    }),
  );

  await renderRoute(Route, {
    path: "/{-$locale}/admin/blog_posts/$id",
    initialPath: "/admin/blog_posts/3",
    user: admin,
  });

  await page.getByPlaceholder("Search…").click();
  await page.getByRole("option", { name: "Course 150" }).click();
  // The dropdown stays open over the button until it is dismissed.
  await userEvent.keyboard("{Escape}");
  await page.getByRole("button", { name: "Save" }).last().click();

  await expect.element(page.getByText("Record updated")).toBeVisible();
  expect(saved).toEqual({ courseIds: [150] });
});

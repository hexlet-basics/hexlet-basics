import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { CourseCategoryInput } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/course_categories/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import { zCourseCategoryInput } from "@/client/zod.gen";

// Positive full-cycle use-case: an admin fills the create form, it POSTs the
// right body to the API, and the app confirms + navigates back to the list.
describe("create course category", () => {
  it("creates a course category", async () => {
    expect.hasAssertions();
    let received: CourseCategoryInput | null = null;
    worker.use(
      http.post("*/admin/course_categories", async ({ request }) => {
        received = zCourseCategoryInput.parse(await request.json());
        return HttpResponse.json({
          createdAt: "2026-01-01T00:00:00Z",
          description: null,
          header: "Frontend courses",
          id: 1,
          locale: "en",
          name: "Frontend",
          slug: "frontend",
        });
      }),
    );

    const { router } = await renderRoute(Route, {
      initialPath: "/admin/course_categories/new",
      path: "/{-$locale}/admin/course_categories/new",
      user: admin,
    });

    await page.getByLabelText("Name *").fill("Frontend");
    await page.getByLabelText("Header *").fill("Frontend courses");
    await page.getByLabelText("Slug *").fill("frontend");
    await page.getByRole("button", { name: "Create" }).click();

    await expect.element(page.getByText("Record created")).toBeVisible();
    expect(received).toMatchObject({
      header: "Frontend courses",
      name: "Frontend",
      slug: "frontend",
    });
    await expect.poll(() => router.state.location.pathname).toBe("/admin/course_categories");
  }, 60_000);
});

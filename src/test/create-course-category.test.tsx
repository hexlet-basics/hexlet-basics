import { HttpResponse, http } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { CourseCategoryInput } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/course_categories/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// Positive full-cycle use-case: an admin fills the create form, it POSTs the
// Right body to the API, and the app confirms + navigates back to the list.
test("creates a course category", async () => {
  let received: CourseCategoryInput | undefined;
  worker.use(
    http.post("*/admin/course_categories", async ({ request }) => {
      received = (await request.json()) as CourseCategoryInput;
      return HttpResponse.json({
        id: 1,
        name: "Frontend",
        header: "Frontend courses",
        slug: "frontend",
        description: null,
        locale: "en",
        createdAt: "2026-01-01T00:00:00Z",
      });
    }),
  );

  const { router } = await renderRoute(Route, {
    path: "/{-$locale}/admin/course_categories/new",
    initialPath: "/admin/course_categories/new",
    user: admin,
  });

  await page.getByLabelText("Name *").fill("Frontend");
  await page.getByLabelText("Header *").fill("Frontend courses");
  await page.getByLabelText("Slug *").fill("frontend");
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Record created")).toBeVisible();
  expect(received).toMatchObject({
    name: "Frontend",
    header: "Frontend courses",
    slug: "frontend",
  });
  await expect.poll(() => router.state.location.pathname).toBe("/admin/course_categories");
});

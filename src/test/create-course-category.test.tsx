import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { CourseCategoryInput } from "@/client/types.gen";
// Route components live under the `{-$locale}` directory, whose braces break
// Vitest's transform when a *test* lives there — so the flow tests sit in
// src/test/ and import the page components across the alias.
import { NewCourseCategory } from "@/routes/{-$locale}/admin/course_categories/new";
import { worker } from "@/test/msw";
import { renderWithRouter } from "@/test/renderWithRouter";

// Positive full-cycle use-case: an admin fills the create form, it POSTs the
// right body to the API, and the app confirms + navigates back to the list.
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

  const { navigate } = await renderWithRouter(<NewCourseCategory />);

  await page.getByLabelText("Name").fill("Frontend");
  await page.getByLabelText("Header").fill("Frontend courses");
  await page.getByLabelText("Slug").fill("frontend");
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Record created")).toBeVisible();
  expect(received).toMatchObject({
    name: "Frontend",
    header: "Frontend courses",
    slug: "frontend",
  });
  expect(navigate).toHaveBeenCalledWith(
    expect.objectContaining({ to: "/{-$locale}/admin/course_categories" }),
  );
});

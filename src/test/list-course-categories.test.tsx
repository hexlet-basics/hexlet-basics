import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { CourseCategoriesList } from "@/routes/{-$locale}/admin/course_categories/index";
import { worker } from "@/test/msw";
import { renderWithRouter } from "@/test/renderWithRouter";

// Positive full-cycle use-case for the read half of the CRUD engine: the list
// hook fetches a server-paginated page and CrudList renders its rows and pager.
//
// This exists because Table v9 only exposes an API when its feature is
// registered — a CrudList missing `rowSortingFeature`/`rowPaginationFeature`
// still type-checks and then renders nothing (or throws in a header). One list
// covers all 14: every admin list is the same CrudList with other columns.
test("renders a server-paginated list through the CRUD engine", async () => {
  worker.use(
    http.get("*/admin/course_categories", () =>
      HttpResponse.json({
        items: [
          {
            id: 1,
            name: "Frontend",
            header: "Frontend courses",
            slug: "frontend",
            description: null,
            locale: "en",
            createdAt: "2026-01-01T00:00:00Z",
          },
          {
            id: 2,
            name: "Backend",
            header: "Backend courses",
            slug: "backend",
            description: null,
            locale: "ru",
            createdAt: "2026-01-02T00:00:00Z",
          },
        ],
        // More than one page of 25, so the pager renders — that is what proves
        // the pagination feature is live rather than merely typed.
        total: 30,
      }),
    ),
  );

  await renderWithRouter(<CourseCategoriesList />);

  // Cell role with `exact`, because getByText matches case-insensitively and
  // each row carries the same word in both its name and its slug.
  await expect.element(page.getByRole("cell", { name: "Frontend", exact: true })).toBeVisible();
  await expect.element(page.getByRole("cell", { name: "backend", exact: true })).toBeVisible();
  await expect.element(page.getByRole("button", { name: "2" })).toBeVisible();
});

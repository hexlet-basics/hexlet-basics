import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import { Route } from "@/routes/{-$locale}/admin/course_categories/index";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// Positive full-cycle use-case for the read half of the CRUD engine: the list
// hook fetches a server-paginated page and CrudList renders its rows and pager.
//
// This exists because Table v9 only exposes an API when its feature is
// registered — a CrudList missing `rowSortingFeature`/`rowPaginationFeature`
// still type-checks and then renders nothing (or throws in a header). One list
// covers all 14: every admin list is the same CrudList with other columns.
describe("list course categories", () => {
  it("renders a server-paginated list through the CRUD engine", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/course_categories", () =>
        HttpResponse.json({
          items: [
            {
              createdAt: "2026-01-01T00:00:00Z",
              description: null,
              header: "Frontend courses",
              id: 1,
              locale: "en",
              name: "Frontend",
              slug: "frontend",
            },
            {
              createdAt: "2026-01-02T00:00:00Z",
              description: null,
              header: "Backend courses",
              id: 2,
              locale: "ru",
              name: "Backend",
              slug: "backend",
            },
          ],
          // More than one page of 25, so the pager renders — that is what proves
          // the pagination feature is live rather than merely typed.
          total: 30,
        }),
      ),
    );

    await renderRoute(Route, {
      initialPath: "/admin/course_categories",
      path: "/{-$locale}/admin/course_categories/",
      user: admin,
    });

    // Cell role with `exact`, because getByText matches case-insensitively and
    // each row carries the same word in both its name and its slug.
    await expect.element(page.getByRole("cell", { exact: true, name: "Frontend" })).toBeVisible();
    await expect.element(page.getByRole("cell", { exact: true, name: "backend" })).toBeVisible();
    await expect.element(page.getByRole("button", { name: "2" })).toBeVisible();
  }, 60_000);
});

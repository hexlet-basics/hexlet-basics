import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { Banner, UserCrud } from "@/client/types.gen";
import { Route as bannersRoute } from "@/routes/{-$locale}/admin/banners/index";
import { Route as managementUsersRoute } from "@/routes/{-$locale}/admin/management/users/index";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const ALL_ROWS = 30;

// What the admin lists do around their table, shown on the banners list: the
// pager, the confirmed delete and the search box. Every list repeats the same
// wiring on a mantine-datatable table, so one real route stands for them.

const banner = (id: number, body: string): Banner => ({
  background: "dark",
  body,
  createdAt: "2026-01-01T00:00:00Z",
  finishesAt: null,
  id,
  locale: "en",
  startsAt: null,
  state: "published",
  url: null,
});

// The page of banners the API answers a list request with: one match for the
// "sale" search, thirty rows otherwise.
function searchedBanners(query: URLSearchParams): {
  items: Banner[];
  page: number;
  perPage: number;
  total: number;
} {
  const found = query.get("q") === "sale";
  return {
    items: [banner(1, found ? "Spring sale" : "Anything")],
    page: Number(query.get("page")),
    perPage: 25,
    total: found ? 1 : ALL_ROWS,
  };
}

const openBanners = (): ReturnType<typeof renderRoute> =>
  renderRoute(bannersRoute, {
    initialPath: "/admin/banners",
    path: "/{-$locale}/admin/banners/",
    user: admin,
  });

describe("admin list", () => {
  it("asks the server for the next page when the pager moves", async () => {
    expect.hasAssertions();
    const requested: URLSearchParams[] = [];
    worker.use(
      http.get("*/admin/banners", ({ request }) => {
        const query = new URL(request.url).searchParams;
        requested.push(query);
        const pageNumber = Number(query.get("page"));
        return HttpResponse.json({
          items: [banner(pageNumber, `Banner on page ${pageNumber}`)],
          page: pageNumber,
          perPage: 25,
          total: ALL_ROWS,
        });
      }),
    );

    await openBanners();

    await expect.element(page.getByRole("cell", { name: "Banner on page 1" })).toBeVisible();
    expect(requested[0]?.get("page")).toBe("1");

    await page.getByRole("button", { name: "2" }).click();

    await expect.element(page.getByRole("cell", { name: "Banner on page 2" })).toBeVisible();
    expect(requested.at(-1)?.get("page")).toBe("2");
  }, 60_000);

  it("asks before deleting a row, then removes it and refreshes the list", async () => {
    expect.hasAssertions();
    let rows = [banner(1, "Spring sale"), banner(2, "Autumn sale")];
    let listed = 0;
    worker.use(
      http.get("*/admin/banners", () => {
        listed += 1;
        return HttpResponse.json({ items: rows, page: 1, perPage: 25, total: rows.length });
      }),
      http.delete("*/admin/banners/1", () => {
        rows = rows.filter((row) => row.id !== 1);
        return new HttpResponse(null, { status: 204 });
      }),
    );

    await openBanners();
    await expect.element(page.getByRole("cell", { name: "Spring sale" })).toBeVisible();

    await page.getByRole("button", { name: "Delete" }).first().click();
    // The confirmation names what is about to go.
    await expect.element(page.getByRole("dialog").getByText("Spring sale")).toBeVisible();
    await page.getByRole("dialog").getByRole("button", { name: "Delete" }).click();

    await expect.element(page.getByText("Record deleted")).toBeVisible();
    await expect.element(page.getByRole("cell", { name: "Spring sale" })).not.toBeInTheDocument();
    await expect.element(page.getByRole("cell", { name: "Autumn sale" })).toBeVisible();
    expect(listed).toBe(2);
  }, 60_000);

  it("sends the typed term as q and starts again from the first page", async () => {
    expect.hasAssertions();
    const requested: URLSearchParams[] = [];
    worker.use(
      http.get("*/admin/banners", ({ request }) => {
        const query = new URL(request.url).searchParams;
        requested.push(query);
        return HttpResponse.json(searchedBanners(query));
      }),
    );

    await openBanners();
    await expect.element(page.getByRole("cell", { name: "Anything" })).toBeVisible();
    expect(requested[0]?.has("q")).toBe(false);
    await page.getByRole("button", { name: "2" }).click();
    await expect.poll(() => requested.at(-1)?.get("page")).toBe("2");

    await page.getByPlaceholder("Search…").fill("sale");

    await expect.element(page.getByRole("cell", { name: "Spring sale" })).toBeVisible();
    expect(requested.at(-1)?.get("q")).toBe("sale");
    expect(requested.at(-1)?.get("page")).toBe("1");
  }, 60_000);

  it("shows only the edit link on a list whose rows cannot be deleted", async () => {
    expect.hasAssertions();
    const learner: UserCrud = {
      admin: false,
      email: "learner@example.com",
      firstName: "Lea",
      id: 1,
      lastName: "Learner",
    };
    worker.use(
      http.get("*/admin/management/users", () =>
        HttpResponse.json({ items: [learner], page: 1, perPage: 25, total: 1 }),
      ),
    );

    await renderRoute(managementUsersRoute, {
      initialPath: "/admin/management/users",
      path: "/{-$locale}/admin/management/users/",
      user: admin,
    });

    await expect.element(page.getByRole("link", { name: "Edit" })).toBeVisible();
    await expect.element(page.getByRole("button", { name: "Delete" })).not.toBeInTheDocument();
  }, 60_000);
});

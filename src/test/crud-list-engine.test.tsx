import { HttpResponse, http } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import {
  adminDeleteBannerMutation,
  adminListBannersOptions,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { Banner } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { worker } from "@/test/msw";
import { renderWithRouter } from "@/test/renderWithRouter";

// The list engine on its own, mounted the way a resource's index route uses it:
// The hook is handed the generated list options and delete mutation, and the
// Resource contributes columns only.

const banner = (id: number, body: string): Banner => ({
  id,
  state: "published",
  locale: "en",
  background: "dark",
  body,
  url: null,
  startsAt: null,
  finishesAt: null,
  createdAt: "2026-01-01T00:00:00Z",
});

function Banners({
  deletable = true,
  searchable = false,
}: {
  deletable?: boolean;
  searchable?: boolean;
}) {
  const list = useCrudList({
    list: adminListBannersOptions,
    listKey: adminListBannersQueryKey,
    remove: deletable ? adminDeleteBannerMutation : undefined,
    searchable,
  });

  const columns: CrudColumnDef<Banner>[] = [
    { accessorKey: "body", header: "Body" },
    list.actionsColumn({
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/banners/$id" params={{ id: String(row.id) }} />
      ),
      describe: (row) => row.body,
    }),
  ];

  return <CrudList columns={columns} {...list.props} />;
}

test("the first request leaves the page size to the server, the pager asks for the next page", async () => {
  const requested: URLSearchParams[] = [];
  worker.use(
    http.get("*/admin/banners", ({ request }) => {
      const query = new URL(request.url).searchParams;
      requested.push(query);
      const pageNumber = Number(query.get("page"));
      return HttpResponse.json({
        items: [banner(pageNumber, `Banner on page ${pageNumber}`)],
        total: 30,
        page: pageNumber,
        perPage: 25,
      });
    }),
  );

  await renderWithRouter(<Banners />, { navigate: "real" });

  await expect.element(page.getByRole("cell", { name: "Banner on page 1" })).toBeVisible();
  expect(requested[0]?.get("page")).toBe("1");
  expect(requested[0]?.has("perPage")).toBe(false);

  await page.getByRole("button", { name: "2" }).click();

  await expect.element(page.getByRole("cell", { name: "Banner on page 2" })).toBeVisible();
  expect(requested.at(-1)?.get("page")).toBe("2");
  expect(requested.at(-1)?.has("perPage")).toBe(false);
});

test("deleting a row asks first, then removes it and refreshes the list", async () => {
  let rows = [banner(1, "Spring sale"), banner(2, "Autumn sale")];
  let listed = 0;
  worker.use(
    http.get("*/admin/banners", () => {
      listed += 1;
      return HttpResponse.json({ items: rows, total: rows.length, page: 1, perPage: 25 });
    }),
    http.delete("*/admin/banners/1", () => {
      rows = rows.filter((row) => row.id !== 1);
      return new HttpResponse(null, { status: 204 });
    }),
  );

  await renderWithRouter(<Banners />, { navigate: "real" });
  await expect.element(page.getByRole("cell", { name: "Spring sale" })).toBeVisible();

  await page.getByRole("button", { name: "Delete" }).first().click();
  // The confirmation names what is about to go.
  await expect.element(page.getByRole("dialog").getByText("Spring sale")).toBeVisible();
  await page.getByRole("dialog").getByRole("button", { name: "Delete" }).click();

  await expect.element(page.getByText("Record deleted")).toBeVisible();
  await expect.element(page.getByRole("cell", { name: "Spring sale" })).not.toBeInTheDocument();
  await expect.element(page.getByRole("cell", { name: "Autumn sale" })).toBeVisible();
  expect(listed).toBe(2);
});

test("a resource without a delete mutation shows only the edit link", async () => {
  worker.use(
    http.get("*/admin/banners", () =>
      HttpResponse.json({ items: [banner(1, "Spring sale")], total: 1, page: 1, perPage: 25 }),
    ),
  );

  await renderWithRouter(<Banners deletable={false} />);

  await expect.element(page.getByRole("link", { name: "Edit" })).toBeVisible();
  await expect.element(page.getByRole("button", { name: "Delete" })).not.toBeInTheDocument();
});

test("a searchable list sends the typed term as q and starts again from the first page", async () => {
  const requested: URLSearchParams[] = [];
  worker.use(
    http.get("*/admin/banners", ({ request }) => {
      const query = new URL(request.url).searchParams;
      requested.push(query);
      const found = query.get("q") === "sale";
      return HttpResponse.json({
        items: [banner(1, found ? "Spring sale" : "Anything")],
        total: found ? 1 : 30,
        page: Number(query.get("page")),
        perPage: 25,
      });
    }),
  );

  await renderWithRouter(<Banners searchable />, { navigate: "real" });
  await expect.element(page.getByRole("cell", { name: "Anything" })).toBeVisible();
  expect(requested[0]?.has("q")).toBe(false);
  await page.getByRole("button", { name: "2" }).click();
  await expect.poll(() => requested.at(-1)?.get("page")).toBe("2");

  await page.getByPlaceholder("Search…").fill("sale");

  await expect.element(page.getByRole("cell", { name: "Spring sale" })).toBeVisible();
  expect(requested.at(-1)?.get("q")).toBe("sale");
  expect(requested.at(-1)?.get("page")).toBe("1");
});

test("a list that is not searchable shows no search box", async () => {
  worker.use(
    http.get("*/admin/banners", () =>
      HttpResponse.json({ items: [banner(1, "Spring sale")], total: 1, page: 1, perPage: 25 }),
    ),
  );

  await renderWithRouter(<Banners />, { navigate: "real" });

  await expect.element(page.getByRole("cell", { name: "Spring sale" })).toBeVisible();
  await expect.element(page.getByPlaceholder("Search…")).not.toBeInTheDocument();
});

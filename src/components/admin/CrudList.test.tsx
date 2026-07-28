import type { ColumnDef } from "@tanstack/react-table";
import { expect, test, vi } from "vitest";
import { page } from "vitest/browser";
import { renderWithProviders } from "@/test/renderWithProviders";
import { CrudList } from "./CrudList";

interface Row {
  id: number;
  name: string;
}

const columns: ColumnDef<Row>[] = [{ accessorKey: "name", header: "Name" }];

const base = {
  columns,
  pagination: { pageIndex: 0, pageSize: 25 },
  sorting: [],
  onPaginationChange: vi.fn(),
  onSortingChange: vi.fn(),
};

test("renders a row per item", async () => {
  renderWithProviders(
    <CrudList
      {...base}
      data={[
        { id: 1, name: "Ruby" },
        { id: 2, name: "Go" },
      ]}
      total={2}
    />,
  );

  await expect.element(page.getByText("Ruby")).toBeVisible();
  await expect.element(page.getByText("Go")).toBeVisible();
});

test("shows the empty state when there are no rows", async () => {
  renderWithProviders(<CrudList {...base} data={[]} total={0} />);

  await expect.element(page.getByText("No records")).toBeVisible();
});

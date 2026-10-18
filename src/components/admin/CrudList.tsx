import { EmptyState, Group, Pagination, Skeleton, Stack, Table, TextInput } from "@mantine/core";
import { IconDatabaseOff } from "@tabler/icons-react";
import { flexRender, rowPaginationFeature, tableFeatures, useTable } from "@tanstack/react-table";
import type { ColumnDef, OnChangeFn, PaginationState, RowData } from "@tanstack/react-table";
import type { ReactElement, ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { pageCount } from "@/lib/page-count";

// Table v9 no longer bundles every feature — each one is registered here, and an
// unregistered feature simply has no API. The admin lists paginate and nothing
// else (the server orders; no column sorts), so that is all we take.
// Deliberately absent is the `paginatedRowModel` slot: the server already
// returns one page, so a client-side row model would re-slice that page (wrong
// results, not just wasted work). The core row model is automatic in v9.
const features = tableFeatures({ rowPaginationFeature });

// V9 threads the feature set through every table type, so the column type is
// exported pre-applied: a resource declares `CrudColumnDef<Course>[]` and stays
// unaware of which features the shared engine registers.
//
// The `TValue` generic is left at its `unknown` default, which is what plain
// literal columns (accessorKey + a custom cell / display columns) produce — no
// `any` needed since our cells read `row.original`, never a typed `getValue()`.
export type CrudColumnDef<T extends RowData> = ColumnDef<typeof features, T>;

// Generic admin list, the read half of the CRUD engine (Wave 1). It is purely
// presentational over TanStack Table in *manual* mode: the owning route runs the
// generated list hook and feeds server-paginated rows plus the table state back
// in, so filtering/ordering stay in SQL (never client-side over one page). Per
// resource you supply only `columns` (including an optional actions column) —
// everything below is shared.
export interface CrudListProps<T extends RowData> {
  columns: CrudColumnDef<T>[];
  data?: T[];
  // Total rows across all pages (from the `XxxPage` envelope), used to derive the
  // page count for the pager — the table itself only ever holds one page.
  total: number;
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  isLoading?: boolean;
  // A search box above the table, for a list the server filters on `q`.
  search?: { value: string; onChange: (value: string) => void };
}

// The placeholder rows shown while a page loads; they have no identity of their own.
// The rows of a list with no data yet: one stable reference, not a new array per render.
const NO_ROWS: never[] = [];

const SKELETON_ROWS = ["first", "second", "third", "fourth", "fifth"];

export function CrudList<T extends RowData>({
  columns,
  data = NO_ROWS,
  total,
  pagination,
  onPaginationChange,
  isLoading,
  search,
}: CrudListProps<T>): ReactElement {
  const { t } = useTranslation();
  const pages = pageCount(total, pagination.pageSize);

  const table = useTable({
    columns,
    data,
    features,
    manualPagination: true,
    onPaginationChange,
    pageCount: pages,
    state: { pagination },
  });

  // Three states of the body, as early returns instead of a nested ternary.
  function renderBody(): ReactNode {
    if (isLoading) {
      return SKELETON_ROWS.map((row) => (
        <Table.Tr key={row}>
          {table.getAllLeafColumns().map((column) => (
            <Table.Td key={column.id}>
              <Skeleton h={16} />
            </Table.Td>
          ))}
        </Table.Tr>
      ));
    }
    if (data.length === 0) {
      return (
        <Table.Tr>
          <Table.Td colSpan={columns.length}>
            <EmptyState
              icon={<IconDatabaseOff />}
              title={t(($) => $.admin.crud.empty)}
              size="sm"
              py="xl"
            />
          </Table.Td>
        </Table.Tr>
      );
    }
    return table.getRowModel().rows.map((row) => (
      <Table.Tr key={row.id}>
        {/* `getAllCells` and not `getVisibleCells`: the latter belongs to
            `columnVisibilityFeature`, and no admin list hides columns —
            taking the core accessor keeps that feature unregistered. */}
        {row.getAllCells().map((cell) => (
          <Table.Td key={cell.id}>
            {flexRender(cell.column.columnDef.cell, cell.getContext())}
          </Table.Td>
        ))}
      </Table.Tr>
    ));
  }

  return (
    <Stack>
      {search && (
        <TextInput
          placeholder={t(($) => $.admin.crud.search)}
          value={search.value}
          onChange={(event) => {
            search.onChange(event.currentTarget.value);
          }}
          maw={360}
        />
      )}
      <Table.ScrollContainer minWidth={600}>
        <Table striped highlightOnHover withTableBorder>
          <Table.Thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <Table.Tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <Table.Th key={header.id}>
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </Table.Th>
                ))}
              </Table.Tr>
            ))}
          </Table.Thead>
          <Table.Tbody>{renderBody()}</Table.Tbody>
        </Table>
      </Table.ScrollContainer>

      {pages > 1 && (
        <Group justify="flex-end">
          <Pagination
            total={pages}
            value={pagination.pageIndex + 1}
            onChange={(page) => {
              onPaginationChange({ ...pagination, pageIndex: page - 1 });
            }}
          />
        </Group>
      )}
    </Stack>
  );
}

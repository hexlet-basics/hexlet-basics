import {
  Center,
  Group,
  Loader,
  Pagination,
  Stack,
  Table,
  Text,
} from "@mantine/core";
import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
  type OnChangeFn,
  type PaginationState,
  type SortingState,
  useReactTable,
} from "@tanstack/react-table";
import { useTranslation } from "react-i18next";

// Generic admin list, the read half of the CRUD engine (Wave 1). It is purely
// presentational over TanStack Table in *manual* mode: the owning route runs the
// generated list hook and feeds server-paginated rows plus the table state back
// in, so filtering/ordering stay in SQL (never client-side over one page). Per
// resource you supply only `columns` (including an optional actions column) —
// everything below is shared.
export interface CrudListProps<T> {
  // `ColumnDef<T>` defaults the cell-value type to `unknown`, which is what plain
  // literal columns (accessorKey + a custom cell / display columns) produce — no
  // `any` needed since our cells read `row.original`, never a typed `getValue()`.
  columns: ColumnDef<T>[];
  data: T[];
  // Total rows across all pages (from the `XxxPage` envelope), used to derive the
  // page count for the pager — the table itself only ever holds one page.
  total: number;
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  sorting: SortingState;
  onSortingChange: OnChangeFn<SortingState>;
  isLoading?: boolean;
}

export function CrudList<T>({
  columns,
  data,
  total,
  pagination,
  onPaginationChange,
  sorting,
  onSortingChange,
  isLoading,
}: CrudListProps<T>) {
  const { t } = useTranslation();
  const pageCount = Math.max(1, Math.ceil(total / pagination.pageSize));

  const table = useReactTable({
    data,
    columns,
    state: { pagination, sorting },
    manualPagination: true,
    manualSorting: true,
    pageCount,
    onPaginationChange,
    onSortingChange,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Stack>
      <Table.ScrollContainer minWidth={600}>
        <Table striped highlightOnHover withTableBorder>
          <Table.Thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <Table.Tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  const canSort = header.column.getCanSort();
                  const sorted = header.column.getIsSorted();
                  return (
                    <Table.Th
                      key={header.id}
                      onClick={
                        canSort
                          ? header.column.getToggleSortingHandler()
                          : undefined
                      }
                      style={{ cursor: canSort ? "pointer" : undefined }}
                    >
                      <Group gap={4} wrap="nowrap">
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                        {sorted === "asc" && "▲"}
                        {sorted === "desc" && "▼"}
                      </Group>
                    </Table.Th>
                  );
                })}
              </Table.Tr>
            ))}
          </Table.Thead>
          <Table.Tbody>
            {isLoading ? (
              <Table.Tr>
                <Table.Td colSpan={columns.length}>
                  <Center py="xl">
                    <Loader />
                  </Center>
                </Table.Td>
              </Table.Tr>
            ) : data.length === 0 ? (
              <Table.Tr>
                <Table.Td colSpan={columns.length}>
                  <Text c="dimmed" ta="center" py="xl">
                    {t(($) => $.admin.crud.empty)}
                  </Text>
                </Table.Td>
              </Table.Tr>
            ) : (
              table.getRowModel().rows.map((row) => (
                <Table.Tr key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <Table.Td key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </Table.Td>
                  ))}
                </Table.Tr>
              ))
            )}
          </Table.Tbody>
        </Table>
      </Table.ScrollContainer>

      {pageCount > 1 && (
        <Group justify="flex-end">
          <Pagination
            total={pageCount}
            value={pagination.pageIndex + 1}
            onChange={(page) =>
              onPaginationChange({ ...pagination, pageIndex: page - 1 })
            }
          />
        </Group>
      )}
    </Stack>
  );
}

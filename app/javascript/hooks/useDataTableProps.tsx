import { router } from "@inertiajs/react";
import type { DataTableSortStatus } from "mantine-datatable";
import { useMemo, useState } from "react";
import { getCurrentUrl } from "@/lib/utils";
import type { Grid } from "@/types";

function cleanObject<T extends Record<string, unknown>>(obj: T): Partial<T> {
  return Object.fromEntries(
    Object.entries(obj).filter(([, v]) => v != null && v !== ""),
  ) as Partial<T>;
}

export default function useDataTableProps<
  T,
  Fields extends Record<string, string | number | undefined>,
>(grid: Grid & { fields: Fields }) {
  const url = getCurrentUrl();

  const [filterValues, setFilterValues] = useState<Fields>(grid.fields ?? {});

  const updateGrid = (extra: Record<string, unknown>) => {
    const params = {
      ...grid,
      fields: cleanObject(filterValues),
      ...extra,
    };
    router.get(url!, params, { preserveScroll: true });
  };

  const onFilterChange = <K extends keyof Fields>(
    key: K,
    value: Fields[K] | null,
  ) => {
    const updated = { ...filterValues, [key]: value ?? undefined };
    setFilterValues(updated);
    updateGrid({ page: 1, fields: cleanObject(updated) });
  };

  const filters: {
    values: Fields;
    getOnChange: <K extends keyof Fields>(
      key: K,
    ) => (val: string | null) => void;
    // biome-ignore lint/correctness/useExhaustiveDependencies: -
  } = useMemo(() => {
    return {
      values: filterValues,
      getOnChange: <K extends keyof Fields>(key: K) => {
        return (val: string | null) => {
          onFilterChange(key, val as Fields[K]);
        };
      },
    };
  }, [filterValues]);

  const sortStatus: DataTableSortStatus<T> = {
    columnAccessor: grid.sf,
    direction: grid.so,
  };

  const handleSortStatusChange = (status: DataTableSortStatus<T>) => {
    updateGrid({
      sf: status.columnAccessor as string,
      so: status.direction,
      page: 1,
    });
  };

  const gridProps = {
    sortStatus,
    onSortStatusChange: handleSortStatusChange,
    page: grid.page,
    onPageChange: (page: number) => updateGrid({ page }),
    recordsPerPage: grid.per,
    totalRecords: grid.tr,
  };

  return { gridProps, filters };
}

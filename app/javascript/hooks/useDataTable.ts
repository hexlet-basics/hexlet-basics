import { url } from "@/lib/utils";
import type { Grid } from "@/types/serializers";
import { router } from "@inertiajs/react";
import type { DataTableStateEvent } from "primereact/datatable";

const useDataTable = () => {
  const handleDataTable = (
    e: Partial<DataTableStateEvent> & { page?: number },
  ) => {
    const newGrid: Partial<Grid> = {};
    if (e.sortField) {
      newGrid.sf = e.sortField;
    }
    if (e.sortOrder) {
      newGrid.so = e.sortOrder;
    }
    if (e.page && e.page > 0) {
      newGrid.page = e.page + 1;
    }
    if (e.filters) {
      newGrid.fields = {};
      Object.entries(e.filters)
        .filter(([_, filter]) => filter.value !== null && filter.value !== "")
        .map(([field, filter]) => {
          newGrid.fields[`${field}_cont`] = filter.value;
        });
    }

    router.get(url(), newGrid);
  };
  return handleDataTable;
};

export default useDataTable;

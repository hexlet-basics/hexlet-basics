import type { Grid } from '@/types';
import type { DataTableSortStatus } from 'mantine-datatable';
import { router } from '@inertiajs/react';
import { url } from '@/lib/utils';

export default function useDataTableProps<T>( grid: Grid) {
  const sortStatus: DataTableSortStatus<T> = {
    columnAccessor: grid.sf,
    direction: grid.so,
  };

  const handleSortStatusChange = (status: DataTableSortStatus<T>) => {
    router.get(url(), {
      ...grid,
      sf: status.columnAccessor as string,
      so: status.direction,
      page: 1, // сбрасываем на первую страницу при смене сортировки
    });
  };

  const handlePageChange = (page: number) => {
    router.get(url(), {
      ...grid,
      page,
    });
  };

  return {
    sortStatus,
    onSortStatusChange: handleSortStatusChange,
    page: grid.page,
    onPageChange: handlePageChange,
    recordsPerPage: grid.per,
    totalRecords: grid.tr,
  };
}


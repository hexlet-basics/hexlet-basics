import type { PropsWithChildren } from "react";

import { useTranslation } from "react-i18next";

import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, User } from "@/types/serializers";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import useDataTable from "@/hooks/useDataTable";

type Props = PropsWithChildren & {
  admins: User[];
  grid: Grid;
};

export default function Index({ admins, grid }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  const handleDataTable = useDataTable();

  return (
    <AdminLayout header={t("admin.home.index.dashboard")}>
      <DataTable
        // header={header}
        value={admins}
        filterDisplay="row"
        lazy
        totalRecords={grid.tr}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        filters={fieldsToFilters(grid.fields)}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        globalFilterFields={["email"]}
      >
        <Column field="id" header="id" />
        <Column field="name" header="name" sortable />
        <Column sortable filter filterPlaceholder="Search" field="email" header="email" />
      </DataTable>
    </AdminLayout>
  );
}

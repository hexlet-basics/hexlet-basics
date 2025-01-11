import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, User } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";

type Props = PropsWithChildren & {
  users: User[];
  grid: Grid;
};

export default function Index({ grid, users }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  const handleDataTable = useDataTable(grid);

  const actionBodyTemplate = (data: User) => {
    return (
      <Link
        className="link-body-emphasis"
        href={Routes.edit_admin_management_user_path(data.id)}
      >
        <i className="bi bi-pencil-fill" />
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.management.users.index.users")}>
      <DataTable
        lazy
        paginator
        totalRecords={grid.tr}
        // header={header}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        filters={fieldsToFilters(grid.filters)}
        value={users}
        // globalFilterFields={["email"]}
      >
        <Column field="id" header="id" />
        <Column field="admin" header="isAdmin" />
        <Column field="name" header="name" sortable />
        <Column sortable field="email" header="email" />
        <Column
          field="created_at"
          header="created_at"
          sortable
          body={DTDateTemplate}
        />
        <Column header="actions" body={actionBodyTemplate} />
      </DataTable>
    </AdminLayout>
  );
}

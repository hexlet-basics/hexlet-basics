import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, User } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  users: User[];
  grid: Grid;
};

export default function Index({ grid, users }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

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
    <AdminLayout header={t("admin.management.users.index.header")}>
      <Menu />
      <DataTable
        lazy
        paginator
        filterDisplay="row"
        totalRecords={grid.tr}
        rows={grid.per}
        first={grid.first}
        sortField={grid.sf}
        sortOrder={grid.so}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        filters={fieldsToFilters(grid.fields)}
        // globalFilterFields={["email"]}
        value={users}
      >
        <Column field="id" header="id" />
        <Column field="admin" header="Admin?" />
        <Column field="assistant_messages_count" header="Messages" />
        <Column field="name" header="name" sortable />
        <Column sortable filter field="email" header="email" />
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

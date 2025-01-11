import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, Review } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";

type Props = PropsWithChildren & {
  reviews: Review[];
  grid: Grid;
};

export default function Index({ grid, reviews }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: Review) => {
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
    <AdminLayout header={t("admin.reviews.index.header")}>
      <DataTable
        lazy
        paginator
        totalRecords={grid.tr}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        filters={fieldsToFilters(grid.fields)}
        value={reviews}
      >
        <Column field="id" header="id" />
        <Column field="first_name" header="First Name" sortable />
        <Column field="last_name" header="Last Name" sortable />
        <Column field="state" header="state" sortable />
        {/* <Column field="body" header="body" /> */}
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

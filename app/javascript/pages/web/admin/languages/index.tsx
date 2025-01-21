import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, OriginalLanguage } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  originalCourses: OriginalLanguage[];
  grid: Grid;
};

export default function Index({ grid, originalCourses }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: OriginalLanguage) => {
    console.log(data)
    return (
      <Link
        className="link-body-emphasis"
        href={Routes.edit_admin_language_path(data.id)}
      >
        <i className="bi bi-pencil-fill" />
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.languages.index.header")}>
      <Menu />
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
        value={originalCourses}
      >
        <Column field="id" header="id" />
        <Column field="name" header="name" sortable />
        <Column field="state" header="state" sortable />
        <Column field="slug" header="slug" />
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


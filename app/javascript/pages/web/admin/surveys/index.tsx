import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, Survey, } from "@/types";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Menu } from "./shared/menu";

type Props = {
  surveys: Survey[];
  grid: Grid;
};

export default function Index({ grid, surveys }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: Survey) => {
    return (
      <Link
        className="link-body-emphasis"
        href={Routes.edit_admin_survey_path(data.id)}
      >
        <i className="bi bi-pencil-fill" />
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.surveys.index.header")}>
      <Menu />
      <DataTable
        lazy
        first={grid.first}
        paginator
        totalRecords={grid.tr}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        filters={fieldsToFilters(grid.fields)}
        value={surveys}
      >
        <Column field="id" header="id" />
        <Column field="parent_survey_id" header="Parent Survey" />
        <Column field="parent_survey_item_value" header="Parent Item" />
        <Column field="slug" header="slug" />
        <Column field="state" header="state" sortable />
        <Column field="question" header="question" />
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


import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, Language } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  courses: Language[];
  grid: Grid;
};

export default function Index({ grid, courses }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: Language) => {
    return (
      <>
        {/* {data.main_landing_slug && ( */}
        {/*   <a */}
        {/*     className="link-body-emphasis me-2" */}
        {/*     target="_blank" */}
        {/*     rel="noreferrer" */}
        {/*     href={Routes.language_path(data.main_landing_slug!)} */}
        {/*   > */}
        {/*     <i className="bi bi-box-arrow-up-right" /> */}
        {/*   </a> */}
        {/* )} */}
        <Link
          className="link-body-emphasis"
          href={Routes.edit_admin_language_path(data.id)}
        >
          <i className="bi bi-pencil-fill" />
        </Link>
      </>
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
        value={courses}
      >
        <Column field="id" header="id" />
        <Column field="current_version_id" header="Version Id" />
        {/* <Column field="current_version.result" header="Version Id" /> */}
        <Column field="slug" header="slug" />
        <Column field="progress" header="Progress" sortable />
        <Column field="learn_as" header="Learn As" sortable />
        <Column field="order" header="order" />
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

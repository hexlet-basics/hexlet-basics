import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, LanguageLandingPage } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Menu } from "./shared/menu";

type Props = PropsWithChildren & {
  landingPages: LanguageLandingPage[];
  grid: Grid;
};

export default function Index({ grid, landingPages }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: LanguageLandingPage) => {
    return (
      <>
        <a
          className="link-body-emphasis me-2"
          target="_blank"
          rel="noreferrer"
          href={Routes.language_path(data.slug!)}
        >
          <i className="bi bi-box-arrow-up-right" />
        </a>
        <Link
          className="link-body-emphasis"
          href={Routes.edit_admin_language_landing_page_path(data.id)}
        >
          <i className="bi bi-pencil-fill" />
        </Link>
      </>
    );
  };

  const languageTemplate = (data: LanguageLandingPage) => {
    return (
      <Link href={Routes.edit_admin_language_path(data.language.id)}>
        {data.language.slug}
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.language_landing_pages.index.header")}>
      <Menu />
      <DataTable
        lazy
        paginator
        filterDisplay="row"
        totalRecords={grid.tr}
        rows={grid.per}
        sortField={grid.sf}
        sortOrder={grid.so}
        filters={fieldsToFilters(grid.fields)}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        value={landingPages}
      >
        <Column field="id" header="id" />
        <Column field="state" header="state" />
        <Column field="order" header="order" />
        <Column
          sortable
          filter
          field="language_slug"
          body={languageTemplate}
          header="language"
        />
        <Column field="slug" header="slug" />
        <Column field="header" header="header" />
        {/* <Column field="order" header="order" /> */}
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

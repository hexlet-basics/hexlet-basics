import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type {
  Grid,
  LanguageCategory,
  LanguageCategoryCrud,
} from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Menu } from "./shared/menu";
import useConfirmation from "@/hooks/useConfirmation";

type Props = PropsWithChildren & {
  categories: LanguageCategoryCrud[];
  grid: Grid;
};

export default function Index({ grid, categories }: Props) {
  const { t } = useTranslation();
  const confirmDeleting = useConfirmation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: LanguageCategory) => {
    return (
      <>
        <a
          className="link-body-emphasis me-2"
          target="_blank"
          rel="noreferrer"
          href={Routes.language_category_path(data.slug!)}
        >
          <i className="bi bi-box-arrow-up-right" />
        </a>
        <Link
          className="link-body-emphasis me-2"
          href={Routes.edit_admin_language_category_path(data.id)}
        >
          <i className="bi bi-pencil-fill" />
        </Link>
        <Link
          onClick={confirmDeleting}
          className="btn btn-link link-body-emphasis p-0 m-0"
          method="delete"
          href={Routes.admin_language_category_path(data.id)}
        >
          {<i className="bi bi-file-x" />}
        </Link>
      </>
    );
  };

  return (
    <AdminLayout header={t("admin.language_categories.index.header")}>
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
        filters={fieldsToFilters(grid.fields)}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        value={categories}
      >
        <Column field="id" header="id" />
        <Column field="slug" header="slug" />
        <Column field="name" header="name" />
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

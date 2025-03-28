import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, LanguageLessonMemberMessage } from "@/types/serializers";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
// import { Menu } from "./shared/menu";

type Props = {
  messages: LanguageLessonMemberMessage[];
  grid: Grid;
};

export default function Index({ grid, messages }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  return (
    <AdminLayout header={t("admin.messages.index.header")}>
      {/* <Menu /> */}
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
        value={messages}
      >
        <Column field="id" header="id" />
        <Column field="role" header="role" />
        <Column field="body" header="body" />
        {/* <Column field="body" header="body" /> */}
        <Column
          field="created_at"
          header="created_at"
          sortable
          body={DTDateTemplate}
        />
        {/* <Column header="actions" body={actionBodyTemplate} /> */}
      </DataTable>
    </AdminLayout>
  );
}

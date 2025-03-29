import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, LanguageLessonMemberMessage } from "@/types/serializers";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Dialog } from "primereact/dialog";
import { useState } from "react";
// import { Menu } from "./shared/menu";

type Props = {
  messages: LanguageLessonMemberMessage[];
  grid: Grid;
};

export function MessageBodyTemplate(message: LanguageLessonMemberMessage) {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <a
        type="button"
        className="link-body-emphasis"
        // biome-ignore lint/a11y/useValidAnchor: <explanation>
        onClick={() => setVisible(true)}
      >
        {message.body?.slice(1, 50)}
      </a>
      <Dialog
        visible={visible}
        modal
        onHide={() => {
          if (!visible) return;
          setVisible(false);
        }}
      >
        {message.body}
      </Dialog>
    </>
  );
}

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
        <Column field="body" header="body" body={MessageBodyTemplate} />
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

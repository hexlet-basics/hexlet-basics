import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, Lead, Review } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { useState } from "react";
import { Dialog } from "primereact/dialog";

type Props = {
  leads: Lead[];
  grid: Grid;
};

export function SurveyAnswersDataTemplate(lead: Lead) {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <a
        type="button"
        className="link-body-emphasis"
        onClick={() => setVisible(true)}
      >
        {"answers_data"}
      </a>
      <Dialog
        style={{ width: "50vw" }}
        visible={visible}
        modal
        onHide={() => {
          if (!visible) return;
          setVisible(false);
        }}
      >
        <pre>{JSON.stringify(lead.survey_answers_data, null, 2)}</pre>
      </Dialog>
    </>
  );
}

export default function Index({ grid, leads }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  // const actionBodyTemplate = (data: Review) => {
  //   return (
  //     <Link
  //       className="link-body-emphasis"
  //       href={Routes.edit_admin_review_path(data.id)}
  //     >
  //       <i className="bi bi-pencil-fill" />
  //     </Link>
  //   );
  // };

  return (
    <AdminLayout header={t("admin.leads.index.header")}>
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
        value={leads}
      >
        <Column field="id" header="id" />
        <Column field="email" header="email" />
        <Column field="user_id" header="user_id" />
        <Column field="full_name" header="full_name" />
        <Column field="phone" header="phone" />
        <Column field="telegram" header="telegram" />
        <Column field="whatsapp" header="whatsapp" />
        <Column field="answers" header="answers" body={SurveyAnswersDataTemplate} />
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


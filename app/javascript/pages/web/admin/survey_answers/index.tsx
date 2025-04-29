import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, Survey, SurveyAnswer, } from "@/types";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";

type Props = {
  surveyAnswers: SurveyAnswer[];
  grid: Grid;
};

export default function Index({ grid, surveyAnswers }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  // const actionBodyTemplate = (data: Survey) => {
  //   return (
  //     <Link
  //       className="link-body-emphasis"
  //       href={Routes.edit_admin_survey_path(data.id)}
  //     >
  //       <i className="bi bi-pencil-fill" />
  //     </Link>
  //   );
  // };

  return (
    <AdminLayout header={t("admin.survey_answers.index.header")}>
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
        value={surveyAnswers}
      >
        <Column field="id" header="id" />
        <Column field="user_id" header="user_id" />
        <Column field="survey_slug" header="survey_slug" />
        <Column field="survey_item_value" header="survey_item_value" />
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



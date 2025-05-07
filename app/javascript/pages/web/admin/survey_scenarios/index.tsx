import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, Survey, SurveyScenarioCrud, } from "@/types";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Menu } from "./shared/menu";

type Props = {
  surveyScenarios: SurveyScenarioCrud[];
  grid: Grid;
};

export default function Index({ grid, surveyScenarios }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: Survey) => {
    return (
      <Link
        className="link-body-emphasis"
        href={Routes.edit_admin_survey_scenario_path(data.id)}
      >
        <i className="bi bi-pencil-fill" />
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.survey_scenarios.index.header")}>
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
        value={surveyScenarios}
      >
        <Column field="id" header="id" />
        {/* <Column field="state" header="state" sortable /> */}
        {/* <Column field="event_name" header="event_name" sortable /> */}
        {/* <Column field="event_treshold_count" header="event treshold count" /> */}
        <Column field="name" header="name" />
        <Column field="survey_item_value" header="Parent Item" />
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


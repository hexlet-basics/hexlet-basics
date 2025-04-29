import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { DTDateTemplate } from "@/components/dtTemplates";
import useDataTable from "@/hooks/useDataTable";
import { fieldsToFilters } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Grid, User, UserSurveyPivot } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";

type Props = PropsWithChildren & {
  userSurveyAnswers: UserSurveyPivot[];
  grid: Grid;
};

export default function Index({ grid, userSurveyAnswers }: Props) {
  const { t } = useTranslation();

  const handleDataTable = useDataTable();

  const actionBodyTemplate = (data: User) => {
    return (
      <Link
        className="link-body-emphasis"
        href={Routes.edit_admin_management_user_path(data.id)}
      >
        <i className="bi bi-pencil-fill" />
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.management.users.index.header")}>
      <DataTable
        lazy
        paginator
        filterDisplay="row"
        totalRecords={grid.tr}
        rows={grid.per}
        first={grid.first}
        sortField={grid.sf}
        sortOrder={grid.so}
        onSort={handleDataTable}
        onFilter={handleDataTable}
        onPage={handleDataTable}
        filters={fieldsToFilters(grid.fields)}
        // globalFilterFields={["email"]}
        value={userSurveyAnswers}
      >
        <Column field="id" header="id" />
        <Column field="user_id" header="user_id" />
        <Column field="goal" header="goal" />
        <Column field="coding_experience" header="coding_experience" />
        <Column field="study_plan" header="study_plan" />
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


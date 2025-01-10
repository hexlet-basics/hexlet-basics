import type { PropsWithChildren } from "react";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import { handleOnSort } from "@/lib/utils";
import AdminLayout from "@/pages/layouts/AdminLayout";
import type { Q, User } from "@/types/serializers";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { Button } from "react-bootstrap";
import { Link } from "@inertiajs/react";
import { DTDateTemplate } from "@/components/dtTemplates";

type Props = PropsWithChildren & {
  users: User[];
  q: Q;
};

export default function Index({ q, users }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  const actionBodyTemplate = (data: User) => {
    return (
      <Link className="link-body-emphasis" href={Routes.edit_admin_management_user_path(data.id)}>
        <i className="bi bi-pencil-fill" />
      </Link>
    );
  };

  return (
    <AdminLayout header={t("admin.management.users.index.users")}>
      <DataTable
        // header={header}
        sortField={q.sf}
        sortOrder={q.so}
        value={users}
        onSort={handleOnSort}
        globalFilterFields={["email"]}
      >
        <Column field="id" header="id" />
        <Column field="admin" header="isAdmin" />
        <Column field="name" header="name" sortable />
        <Column sortable field="email" header="email" />
        <Column field="created_at" header="created_at" sortable body={DTDateTemplate} />
        <Column header="actions" body={actionBodyTemplate} />
      </DataTable>
    </AdminLayout>
  );
}

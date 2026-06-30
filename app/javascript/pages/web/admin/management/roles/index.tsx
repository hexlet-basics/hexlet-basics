import { Link } from "@inertiajs/react";
import { IconEdit, IconLock, IconTrash } from "@tabler/icons-react";
import { DataTable } from "mantine-datatable";
import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { StaffRole } from "@/types";
import { Menu } from "./shared/menu";

type Props = {
  roles: StaffRole[];
  totalPermissionsCount: number;
};

export default function Index({ roles, totalPermissionsCount }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.management.roles.index.header)}>
      <Menu />
      <DataTable
        records={roles}
        columns={[
          { accessor: "id" },
          { accessor: "name" },
          { accessor: "description" },
          {
            accessor: "permissions_count",
            render: (r) => `${r.permissions_count} / ${totalPermissionsCount}`,
          },
          {
            accessor: "edit",
            title: "",
            render: (r) => (
              <Link href={Routes.edit_admin_management_role_path(r.id)}>
                <IconEdit size={14} />
              </Link>
            ),
          },
          {
            accessor: "permissions",
            title: "",
            render: (r) => (
              <Link href={Routes.admin_management_role_permission_path(r.id)}>
                <IconLock size={14} />
              </Link>
            ),
          },
          {
            accessor: "destroy",
            title: "",
            render: (r) => (
              <Link
                href={Routes.admin_management_role_path(r.id)}
                method="delete"
                as="button"
              >
                <IconTrash size={14} />
              </Link>
            ),
          },
        ]}
      />
    </AdminLayout>
  );
}

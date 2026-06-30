import { Link } from "@inertiajs/react";
import { IconEdit, IconTrash } from "@tabler/icons-react";
import { DataTable } from "mantine-datatable";
import { useTranslation } from "react-i18next";
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Grid, StaffMember } from "@/types";
import { Menu } from "./shared/menu";

type Props = {
  staffMembers: StaffMember[];
  grid: Grid;
};

export default function Index({ staffMembers, grid }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<StaffMember, typeof grid.fields>(
    grid,
  );

  return (
    <AdminLayout
      header={t(($) => $.admin.management.staff_members.index.header)}
    >
      <Menu />
      <DataTable
        records={staffMembers}
        columns={[
          { accessor: "id" },
          { accessor: "user.name", title: "User", render: (r) => r.user.name },
          { accessor: "role.name", title: "Role", render: (r) => r.role.name },
          {
            accessor: "allowed_locales",
            render: (r) => r.allowed_locales.join(", "),
          },
          {
            accessor: "edit",
            title: "",
            render: (r) => (
              <Link href={Routes.edit_admin_management_staff_member_path(r.id)}>
                <IconEdit size={14} />
              </Link>
            ),
          },
          {
            accessor: "destroy",
            title: "",
            render: (r) => (
              <Link
                href={Routes.admin_management_staff_member_path(r.id)}
                method="delete"
                as="button"
              >
                <IconTrash size={14} />
              </Link>
            ),
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}

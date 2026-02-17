import { TextInput } from "@mantine/core";
import { IconEdit } from "@tabler/icons-react";
import dayjs from "dayjs";
import { DataTable } from "mantine-datatable";
import { useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { Grid, User } from "@/types/serializers";
import { Menu } from "./shared/menu";

type Props = {
  users: User[];
  grid: Grid;
};

export default function Index({ grid, users }: Props) {
  const { t } = useTranslation();
  const { gridProps, filters } = useDataTableProps<User, typeof grid.fields>(
    grid,
  );

  const renderActions = (item: User) => (
    <AppAnchor href={Routes.edit_admin_management_user_path(item.id)}>
      <IconEdit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t(($) => $.admin.management.users.index.header)}>
      <Menu />
      <DataTable
        records={users}
        columns={[
          { accessor: "id" },
          { accessor: "admin", title: "Admin?" },
          { accessor: "assistant_messages_count", title: "Messages" },
          { accessor: "name", sortable: true },
          {
            accessor: "email",
            sortable: true,
            // filter: (
            //   <TextInput
            //     value={filters.values['fields[email_cont]']}
            //     onChange={filters.getOnChange('fields[email_cont]')}
            //     placeholder={t('admin.management.users.index.search_by_email')}
            //   />
            // ),
          },
          {
            accessor: "created_at",
            sortable: true,
            render: (r) => dayjs(r.created_at).format("LL"),
          },
          { accessor: "actions", title: "actions", render: renderActions },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}

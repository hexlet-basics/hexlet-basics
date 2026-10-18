import { Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminListManagementUsersOptions,
  adminListManagementUsersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListManagementUsersQuery } from "@/client/zod.gen";
import type { UserCrud } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";

// Management users (legacy `/admin/management/users`) — list + edit only,
// no create/delete (parity with legacy: accounts are self-registered).
export const Route = createFileRoute("/{-$locale}/admin/management/users/")({
  validateSearch: zAdminListManagementUsersQuery,
  component: ManagementUsersList,
});

function ManagementUsersList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListManagementUsersOptions,
    searchable: true,
    listKey: adminListManagementUsersQueryKey,
  });

  const columns: CrudColumnDef<UserCrud>[] = [
    {
      accessorKey: "email",
      header: t(($) => $.models.attributes.user.email),
    },
    {
      accessorKey: "firstName",
      header: t(($) => $.models.attributes.user.first_name),
    },
    {
      accessorKey: "lastName",
      header: t(($) => $.models.attributes.user.last_name),
    },
    {
      id: "admin",
      header: t(($) => $.models.attributes.user.admin),
      cell: ({ row }) => (row.original.admin ? "✓" : ""),
    },
    list.actionsColumn({
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/management/users/$id" params={{ id: String(row.id) }} />
      ),
      describe: (row) => row.email ?? "",
    }),
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.managementUsers)}</Title>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteUserMutation,
  adminListUsersOptions,
  adminListUsersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListUsersQuery } from "@/client/zod.gen";
import type { UserCrud } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { userColumns } from "@/components/admin/generated/user";
import { useGeneratedColumns } from "@/components/admin/useGeneratedFields";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

// Users admin (legacy `/admin/api/users` — the `api/` segment was a Rails
// namespacing artifact and is dropped from the UI path). The search box is the
// list's own `q` filter, so a search pages and totals like the list does.
export const Route = createFileRoute("/{-$locale}/admin/users/")({
  validateSearch: zAdminListUsersQuery,
  component: UsersList,
});

function UsersList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListUsersOptions,
    listKey: adminListUsersQueryKey,
    remove: adminDeleteUserMutation,
    searchable: true,
  });

  const columns: CrudColumnDef<UserCrud>[] = [
    ...useGeneratedColumns<UserCrud>(userColumns, ["email", "firstName", "lastName"]),
    {
      id: "admin",
      header: t(($) => $.models.attributes.user.admin),
      cell: ({ row }) => (row.original.admin ? "✓" : ""),
    },
    list.actionsColumn({
      edit: (row) => <EditLink to="/{-$locale}/admin/users/$id" params={{ id: String(row.id) }} />,
      describe: (row) => row.email ?? "",
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.users)}</Title>
        <ButtonLink to="/{-$locale}/admin/users/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

import { Group, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListManagementUsersOptions } from "@/client/@tanstack/react-query.gen";
import type { UserCrud } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";
import { ButtonLink } from "@/components/RouterLink";

// Management users (legacy `/admin/management/users`) — list + edit only,
// no create/delete (parity with legacy: accounts are self-registered).
export const Route = createFileRoute("/{-$locale}/admin/management/users/")({
  component: ManagementUsersList,
});

function ManagementUsersList() {
  const { t } = useTranslation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListManagementUsersOptions({
      query: pager.query,
    }),
  );

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
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/management/users/$id"
            params={{ id: String(row.original.id) }}
            size="xs"
            variant="light"
          >
            {t(($) => $.admin.crud.edit)}
          </ButtonLink>
        </Group>
      ),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.managementUsers)}</Title>

      <CrudList
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        {...pager.props(data?.perPage)}
        isLoading={isLoading}
      />
    </Stack>
  );
}

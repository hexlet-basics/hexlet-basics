import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteRoleMutation,
  adminListRolesOptions,
  adminListRolesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListRolesQuery } from "@/client/zod.gen";
import type { StaffRole } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/")({
  validateSearch: zAdminListRolesQuery,
  component: RolesList,
});

function RolesList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListRolesOptions,
    listKey: adminListRolesQueryKey,
    remove: adminDeleteRoleMutation,
  });

  const columns: CrudColumnDef<StaffRole>[] = [
    {
      accessorKey: "name",
      header: t(($) => $.models.attributes.base.name),
    },
    {
      accessorKey: "description",
      header: t(($) => $.models.attributes.base.description),
    },
    {
      accessorKey: "permissionsCount",
      header: t(($) => $.models.attributes.staff_role.permissions_count),
    },
    list.actionsColumn({
      describe: (row) => row.name,
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/management/roles/$id" params={{ id: String(row.id) }} />
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.roles)}</Title>
        <ButtonLink to="/{-$locale}/admin/management/roles/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

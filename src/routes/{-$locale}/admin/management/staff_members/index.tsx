import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteStaffMemberMutation,
  adminListStaffMembersOptions,
  adminListStaffMembersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListStaffMembersQuery } from "@/client/zod.gen";
import type { StaffMember } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/")({
  component: StaffMembersList,
  validateSearch: zAdminListStaffMembersQuery,
});

function StaffMembersList(): ReactElement {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListStaffMembersOptions,
    listKey: adminListStaffMembersQueryKey,
    remove: adminDeleteStaffMemberMutation,
  });

  const columns: CrudColumnDef<StaffMember>[] = [
    {
      cell: ({ row }) => row.original.user.email ?? String(row.original.userId),
      header: t(($) => $.models.attributes.staff_member.user_id),
      id: "user",
    },
    {
      cell: ({ row }) => row.original.role.name,
      header: t(($) => $.models.attributes.staff_member.role_id),
      id: "role",
    },
    {
      cell: ({ row }) => row.original.allowedLocales.join(", "),
      header: t(($) => $.models.attributes.staff_member.allowed_locales),
      id: "allowedLocales",
    },
    list.actionsColumn({
      describe: (row) => row.user.email ?? "",
      edit: (row) => (
        <EditLink
          to="/{-$locale}/admin/management/staff_members/$id"
          params={{ id: String(row.id) }}
        />
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.staffMembers)}</Title>
        <ButtonLink to="/{-$locale}/admin/management/staff_members/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

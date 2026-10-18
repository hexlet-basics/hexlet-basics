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

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/")({
  validateSearch: zAdminListStaffMembersQuery,
  component: StaffMembersList,
});

function StaffMembersList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListStaffMembersOptions,
    listKey: adminListStaffMembersQueryKey,
    remove: adminDeleteStaffMemberMutation,
  });

  const columns: CrudColumnDef<StaffMember>[] = [
    {
      id: "user",
      header: t(($) => $.models.attributes.staff_member.user_id),
      cell: ({ row }) => row.original.user.email ?? String(row.original.userId),
    },
    {
      id: "role",
      header: t(($) => $.models.attributes.staff_member.role_id),
      cell: ({ row }) => row.original.role.name,
    },
    {
      id: "allowedLocales",
      header: t(($) => $.models.attributes.staff_member.allowed_locales),
      cell: ({ row }) => row.original.allowedLocales.join(", "),
    },
    list.actionsColumn({
      edit: (row) => (
        <EditLink
          to="/{-$locale}/admin/management/staff_members/$id"
          params={{ id: String(row.id) }}
        />
      ),
      describe: (row) => row.user.email ?? "",
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

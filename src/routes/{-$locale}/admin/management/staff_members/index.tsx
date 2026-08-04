import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteStaffMemberMutation,
  adminListStaffMembersOptions,
  adminListStaffMembersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { StaffMember } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/")({
  component: StaffMembersList,
});

function StaffMembersList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListStaffMembersOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteStaffMemberMutation(),
    invalidate: [adminListStaffMembersQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: CrudColumnDef<StaffMember>[] = [
    {
      id: "user",
      header: t(($) => $.models.attributes.staff_member.user_id),
      cell: ({ row }) => row.original.user.email ?? String(row.original.userId),
      enableSorting: false,
    },
    {
      id: "role",
      header: t(($) => $.models.attributes.staff_member.role_id),
      cell: ({ row }) => row.original.role.name,
      enableSorting: false,
    },
    {
      id: "allowedLocales",
      header: t(($) => $.models.attributes.staff_member.allowed_locales),
      cell: ({ row }) => row.original.allowedLocales.join(", "),
      enableSorting: false,
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/management/staff_members/$id"
            params={{ id: String(row.original.id) }}
            size="xs"
            variant="light"
          >
            {t(($) => $.admin.crud.edit)}
          </ButtonLink>
          <Tooltip label={t(($) => $.admin.crud.delete)}>
            <ActionIcon
              color="red"
              variant="light"
              aria-label={t(($) => $.admin.crud.delete)}
              onClick={() =>
                confirmDelete({
                  description: row.original.user.email ?? "",
                  onConfirm: () => deleteMutation.mutate({ path: { id: row.original.id } }),
                })
              }
            >
              ✕
            </ActionIcon>
          </Tooltip>
        </Group>
      ),
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.staffMembers)}</Title>
        <ButtonLink to="/{-$locale}/admin/management/staff_members/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <CrudList
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        pagination={pagination}
        onPaginationChange={setPagination}
        sorting={sorting}
        onSortingChange={setSorting}
        isLoading={isLoading}
      />
    </Stack>
  );
}

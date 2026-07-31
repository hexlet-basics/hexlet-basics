import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { ColumnDef, PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteRoleMutation,
  adminListRolesOptions,
  adminListRolesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { StaffRole } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/")({
  component: RolesList,
});

function RolesList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListRolesOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteRoleMutation(),
    invalidate: [adminListRolesQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: ColumnDef<StaffRole>[] = [
    {
      accessorKey: "name",
      header: t(($) => $.models.attributes.base.name),
      enableSorting: false,
    },
    {
      accessorKey: "description",
      header: t(($) => $.models.attributes.base.description),
      enableSorting: false,
    },
    {
      accessorKey: "permissionsCount",
      header: t(($) => $.models.attributes.staff_role.permissions_count),
      enableSorting: false,
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/management/roles/$id"
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
                  description: row.original.name,
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
        <Title order={2}>{t(($) => $.admin.resources.roles)}</Title>
        <ButtonLink to="/{-$locale}/admin/management/roles/new">
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

import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBannerMutation,
  adminListBannersOptions,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { Banner } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Banners admin list — second resource through the CRUD engine; it exercises the
// select and datetime field types the engine gained for this resource.
export const Route = createFileRoute("/{-$locale}/admin/banners/")({
  component: BannersList,
});

function BannersList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListBannersOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteBannerMutation(),
    invalidate: [adminListBannersQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: CrudColumnDef<Banner>[] = [
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.banner.state),
      enableSorting: false,
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.banner.locale),
      enableSorting: false,
    },
    {
      accessorKey: "background",
      header: t(($) => $.models.attributes.banner.background),
      enableSorting: false,
    },
    {
      accessorKey: "body",
      header: t(($) => $.models.attributes.banner.body),
      enableSorting: false,
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/banners/$id"
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
                  description: row.original.body,
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
        <Title order={2}>{t(($) => $.admin.resources.banners)}</Title>
        <ButtonLink to="/{-$locale}/admin/banners/new">{t(($) => $.admin.crud.new)}</ButtonLink>
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

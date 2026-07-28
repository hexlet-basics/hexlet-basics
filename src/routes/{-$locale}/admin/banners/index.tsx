import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { modals } from "@mantine/modals";
import { notifications } from "@mantine/notifications";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  createColumnHelper,
  type PaginationState,
  type SortingState,
} from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBannerMutation,
  adminListBannersOptions,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { Banner } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";

// Banners admin list — second resource through the CRUD engine; it exercises the
// select and datetime field types the engine gained for this resource.
export const Route = createFileRoute("/{-$locale}/admin/banners/")({
  component: BannersList,
});

const columnHelper = createColumnHelper<Banner>();

function BannersList() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();

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

  const deleteMutation = useMutation({
    ...adminDeleteBannerMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminListBannersQueryKey() });
      notifications.show({
        color: "green",
        message: t(($) => $.admin.crud.deleted),
      });
    },
    onError: () =>
      notifications.show({
        color: "red",
        message: t(($) => $.admin.crud.deleteError),
      }),
  });

  const confirmDelete = (banner: Banner) =>
    modals.openConfirmModal({
      title: t(($) => $.admin.crud.confirmDeleteTitle),
      children: banner.body,
      labels: {
        confirm: t(($) => $.admin.crud.delete),
        cancel: t(($) => $.admin.crud.cancel),
      },
      confirmProps: { color: "red" },
      onConfirm: () => deleteMutation.mutate({ path: { id: banner.id } }),
    });

  const columns = [
    columnHelper.accessor("state", {
      header: t(($) => $.models.attributes.banner.state),
      enableSorting: false,
    }),
    columnHelper.accessor("locale", {
      header: t(($) => $.models.attributes.banner.locale),
      enableSorting: false,
    }),
    columnHelper.accessor("background", {
      header: t(($) => $.models.attributes.banner.background),
      enableSorting: false,
    }),
    columnHelper.accessor("body", {
      header: t(($) => $.models.attributes.banner.body),
      enableSorting: false,
    }),
    columnHelper.display({
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
              onClick={() => confirmDelete(row.original)}
            >
              ✕
            </ActionIcon>
          </Tooltip>
        </Group>
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.banners)}</Title>
        <ButtonLink to="/{-$locale}/admin/banners/new">
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

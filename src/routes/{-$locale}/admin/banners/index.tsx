import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBannerMutation,
  adminListBannersOptions,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListBannersQuery } from "@/client/zod.gen";
import type { Banner } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { bannerColumns } from "@/components/admin/generated/banner";
import { useGeneratedColumns } from "@/components/admin/useGeneratedFields";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import type { ReactElement } from "react";

// Banners admin list, rendered through the CRUD engine.
export const Route = createFileRoute("/{-$locale}/admin/banners/")({
  validateSearch: zAdminListBannersQuery,
  component: BannersList,
});

function BannersList(): ReactElement {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListBannersOptions,
    listKey: adminListBannersQueryKey,
    remove: adminDeleteBannerMutation,
    searchable: true,
  });

  // The state shows as its label: it moves only through the transition buttons
  // of the edit page, so the list is where it is read.
  const labelState = (column: CrudColumnDef<Banner>): CrudColumnDef<Banner> =>
    "accessorKey" in column && column.accessorKey === "state"
      ? {
          ...column,
          cell: ({ row }) =>
            t(($) => $.models.attributes.banner["state/values"][row.original.state]),
        }
      : column;

  const columns: CrudColumnDef<Banner>[] = [
    ...useGeneratedColumns<Banner>(bannerColumns, ["state", "locale", "background", "body"]).map(
      (column) => labelState(column),
    ),
    list.actionsColumn({
      describe: (row) => row.body,
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/banners/$id" params={{ id: String(row.id) }} />
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.banners)}</Title>
        <ButtonLink to="/{-$locale}/admin/banners/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

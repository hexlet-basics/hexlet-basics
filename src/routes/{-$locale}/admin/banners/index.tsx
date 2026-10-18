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
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { bannerColumns } from "@/components/admin/generated/banner";
import { useGeneratedColumns } from "@/components/admin/useGeneratedFields";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

// Banners admin list, rendered through the CRUD engine.
export const Route = createFileRoute("/{-$locale}/admin/banners/")({
  validateSearch: zAdminListBannersQuery,
  component: BannersList,
});

function BannersList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListBannersOptions,
    searchable: true,
    listKey: adminListBannersQueryKey,
    remove: adminDeleteBannerMutation,
  });

  const columns: CrudColumnDef<Banner>[] = [
    // The state shows as its label: it moves only through the transition buttons
    // of the edit page, so the list is where it is read.
    ...useGeneratedColumns<Banner>(bannerColumns, ["state", "locale", "background", "body"]).map(
      (column): CrudColumnDef<Banner> =>
        "accessorKey" in column && column.accessorKey === "state"
          ? {
              ...column,
              cell: ({ row }) =>
                t(($) => $.models.attributes.banner["state/values"][row.original.state]),
            }
          : column,
    ),
    list.actionsColumn({
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/banners/$id" params={{ id: String(row.id) }} />
      ),
      describe: (row) => row.body,
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

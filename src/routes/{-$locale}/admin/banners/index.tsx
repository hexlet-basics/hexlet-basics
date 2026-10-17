import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteBannerMutation,
  adminListBannersOptions,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { Banner } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { bannerColumns } from "@/components/admin/generated/banner";
import { useGeneratedColumns } from "@/components/admin/useGeneratedFields";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

// Banners admin list — second resource through the CRUD engine; it exercises the
// select and datetime field types the engine gained for this resource.
export const Route = createFileRoute("/{-$locale}/admin/banners/")({
  component: BannersList,
});

function BannersList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListBannersOptions,
    listKey: adminListBannersQueryKey,
    remove: adminDeleteBannerMutation,
  });

  const columns: CrudColumnDef<Banner>[] = [
    ...useGeneratedColumns<Banner>(bannerColumns, ["state", "locale", "background", "body"]),
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

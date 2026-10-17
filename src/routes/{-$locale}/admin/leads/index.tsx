import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListLeadsOptions } from "@/client/@tanstack/react-query.gen";
import type { Lead } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";

// Sales leads (legacy `/admin/leads`) — read-only list.
export const Route = createFileRoute("/{-$locale}/admin/leads/")({
  component: LeadsList,
});

function LeadsList() {
  const { t } = useTranslation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListLeadsOptions({
      query: pager.query,
    }),
  );

  const columns: CrudColumnDef<Lead>[] = [
    {
      accessorKey: "email",
      header: t(($) => $.models.attributes.user.email),
    },
    {
      accessorKey: "phone",
      header: t(($) => $.models.attributes.lead.phone),
    },
    {
      accessorKey: "whatsapp",
      header: t(($) => $.models.attributes.lead.whatsapp),
    },
    {
      accessorKey: "telegram",
      header: t(($) => $.models.attributes.lead.telegram),
    },
    {
      id: "createdAt",
      header: t(($) => $.models.attributes.base.created_at),
      cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.leads)}</Title>

      <CrudList
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        {...pager.props(data?.perPage)}
        isLoading={isLoading}
      />
    </Stack>
  );
}

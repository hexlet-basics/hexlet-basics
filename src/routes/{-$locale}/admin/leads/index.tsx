import { Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListLeadsOptions, adminListLeadsQueryKey } from "@/client/@tanstack/react-query.gen";
import { zAdminListLeadsQuery } from "@/client/zod.gen";
import type { Lead } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudList } from "@/components/admin/useCrudList";
import { useDateTime } from "@/hooks/useDateTime";

// Sales leads (legacy `/admin/leads`) — read-only list.
export const Route = createFileRoute("/{-$locale}/admin/leads/")({
  validateSearch: zAdminListLeadsQuery,
  component: LeadsList,
});

function LeadsList() {
  const { t } = useTranslation();
  const formatDateTime = useDateTime();
  const list = useCrudList({
    list: adminListLeadsOptions,
    listKey: adminListLeadsQueryKey,
    searchable: true,
  });

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
      cell: ({ row }) => formatDateTime(row.original.createdAt),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.leads)}</Title>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

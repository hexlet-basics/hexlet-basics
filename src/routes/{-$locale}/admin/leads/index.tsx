import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { ColumnDef, PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { adminListLeadsOptions } from "@/client/@tanstack/react-query.gen";
import type { Lead } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";

// Sales leads (legacy `/admin/leads`) — read-only list.
export const Route = createFileRoute("/{-$locale}/admin/leads/")({
  component: LeadsList,
});

function LeadsList() {
  const { t } = useTranslation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListLeadsOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const columns: ColumnDef<Lead>[] = [
    {
      accessorKey: "email",
      header: t(($) => $.models.attributes.user.email),
      enableSorting: false,
    },
    {
      accessorKey: "phone",
      header: t(($) => $.models.attributes.lead.phone),
      enableSorting: false,
    },
    {
      accessorKey: "whatsapp",
      header: t(($) => $.models.attributes.lead.whatsapp),
      enableSorting: false,
    },
    {
      accessorKey: "telegram",
      header: t(($) => $.models.attributes.lead.telegram),
      enableSorting: false,
    },
    {
      id: "createdAt",
      header: t(($) => $.models.attributes.base.created_at),
      cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
      enableSorting: false,
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.leads)}</Title>

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

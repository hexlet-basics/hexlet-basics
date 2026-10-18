import { Stack, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListLeadsOptions } from "@/client/@tanstack/react-query.gen";
import { zAdminListLeadsQuery } from "@/client/zod.gen";
import type { Lead } from "@/client/types.gen";
import { useDateTime } from "@/hooks/useDateTime";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// Sales leads (legacy `/admin/leads`) — read-only list.
export const Route = createFileRoute("/{-$locale}/admin/leads/")({
  component: LeadsList,
  validateSearch: zAdminListLeadsQuery,
});

function LeadsList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListLeadsOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );
  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);
  const formatDateTime = useDateTime();

  const columns: DataTableColumn<Lead>[] = [
    {
      accessor: "email",
      title: t(($) => $.models.attributes.user.email),
    },
    {
      accessor: "phone",
      title: t(($) => $.models.attributes.lead.phone),
    },
    {
      accessor: "whatsapp",
      title: t(($) => $.models.attributes.lead.whatsapp),
    },
    {
      accessor: "telegram",
      title: t(($) => $.models.attributes.lead.telegram),
    },
    {
      accessor: "createdAt",
      render: (record) => formatDateTime(record.createdAt),
      title: t(($) => $.models.attributes.base.created_at),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.leads)}</Title>

      <TextInput
        maw={360}
        placeholder={t(($) => $.admin.crud.search)}
        value={term}
        onChange={(event) => {
          setTerm(event.currentTarget.value);
          applyTerm(event.currentTarget.value.trim());
        }}
      />

      <DataTable
        withTableBorder
        striped
        highlightOnHover
        fetching={isFetching}
        records={data?.items}
        totalRecords={data?.total}
        recordsPerPage={data?.perPage ?? search.perPage}
        page={search.page ?? 1}
        onPageChange={(page) => {
          void navigate({ search: (prev) => ({ ...prev, page }) });
        }}
        noRecordsText={t(($) => $.admin.crud.empty)}
        columns={columns}
      />
    </Stack>
  );
}

import { Stack, Text, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListMessagesOptions } from "@/client/@tanstack/react-query.gen";
import { zAdminListMessagesQuery } from "@/client/zod.gen";
import type { LessonAssistantMessage } from "@/client/types.gen";
import { useDateTime } from "@/hooks/useDateTime";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// In-lesson assistant messages (legacy `/admin/messages`) — read-only list for
// mining student questions; the per-lesson summaries live under lesson reviews.
export const Route = createFileRoute("/{-$locale}/admin/messages/")({
  component: MessagesList,
  validateSearch: zAdminListMessagesQuery,
});

function MessagesList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListMessagesOptions({
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

  const columns: DataTableColumn<LessonAssistantMessage>[] = [
    {
      accessor: "role",
      title: t(($) => $.models.attributes.assistant_message.role),
    },
    {
      accessor: "content",
      render: (record) => (
        <Text size="sm" lineClamp={2} maw={480}>
          {record.content}
        </Text>
      ),
      title: t(($) => $.models.attributes.assistant_message.content),
    },
    {
      accessor: "courseSlug",
      title: t(($) => $.models.attributes.assistant_message.course),
    },
    {
      accessor: "courseLessonName",
      title: t(($) => $.models.attributes.assistant_message.lesson),
    },
    {
      accessor: "createdAt",
      render: (record) => formatDateTime(record.createdAt),
      title: t(($) => $.models.attributes.base.created_at),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.messages)}</Title>

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

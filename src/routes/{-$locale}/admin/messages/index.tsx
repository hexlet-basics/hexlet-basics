import { Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListMessagesOptions } from "@/client/@tanstack/react-query.gen";
import type { LessonAssistantMessage } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";

// In-lesson assistant messages (legacy `/admin/messages`) — read-only list for
// mining student questions; the per-lesson summaries live under lesson reviews.
export const Route = createFileRoute("/{-$locale}/admin/messages/")({
  component: MessagesList,
});

function MessagesList() {
  const { t } = useTranslation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListMessagesOptions({
      query: pager.query,
    }),
  );

  const columns: CrudColumnDef<LessonAssistantMessage>[] = [
    {
      accessorKey: "role",
      header: t(($) => $.models.attributes.assistant_message.role),
    },
    {
      id: "content",
      header: t(($) => $.models.attributes.assistant_message.content),
      cell: ({ row }) => (
        <Text size="sm" lineClamp={2} maw={480}>
          {row.original.content}
        </Text>
      ),
    },
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.assistant_message.course),
    },
    {
      accessorKey: "courseLessonName",
      header: t(($) => $.models.attributes.assistant_message.lesson),
    },
    {
      id: "createdAt",
      header: t(($) => $.models.attributes.base.created_at),
      cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.messages)}</Title>

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

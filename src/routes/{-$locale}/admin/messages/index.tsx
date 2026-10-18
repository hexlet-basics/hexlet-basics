import { Stack, Text, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminListMessagesOptions,
  adminListMessagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListMessagesQuery } from "@/client/zod.gen";
import type { LessonAssistantMessage } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { useCrudList } from "@/components/admin/useCrudList";
import { useDateTime } from "@/hooks/useDateTime";

// In-lesson assistant messages (legacy `/admin/messages`) — read-only list for
// mining student questions; the per-lesson summaries live under lesson reviews.
export const Route = createFileRoute("/{-$locale}/admin/messages/")({
  validateSearch: zAdminListMessagesQuery,
  component: MessagesList,
});

function MessagesList() {
  const { t } = useTranslation();

  const formatDateTime = useDateTime();
  const list = useCrudList({
    list: adminListMessagesOptions,
    listKey: adminListMessagesQueryKey,
    searchable: true,
  });

  const columns: CrudColumnDef<LessonAssistantMessage>[] = [
    {
      accessorKey: "role",
      header: t(($) => $.models.attributes.assistant_message.role),
    },
    {
      cell: ({ row }) => (
        <Text size="sm" lineClamp={2} maw={480}>
          {row.original.content}
        </Text>
      ),
      header: t(($) => $.models.attributes.assistant_message.content),
      id: "content",
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
      cell: ({ row }) => formatDateTime(row.original.createdAt),
      header: t(($) => $.models.attributes.base.created_at),
      id: "createdAt",
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.messages)}</Title>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

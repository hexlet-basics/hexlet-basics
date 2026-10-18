import { Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminListLessonProgressOptions,
  adminListLessonProgressQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListLessonProgressQuery } from "@/client/zod.gen";
import type { LessonProgress } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { useCrudList } from "@/components/admin/useCrudList";
import { useDateTime } from "@/hooks/useDateTime";
import type { ReactElement } from "react";

// Per-lesson progress (route keeps the legacy table name) — read-only.
export const Route = createFileRoute("/{-$locale}/admin/lesson_progress/")({
  component: LessonProgressList,
  validateSearch: zAdminListLessonProgressQuery,
});

function LessonProgressList(): ReactElement {
  const { t } = useTranslation();

  const formatDateTime = useDateTime();
  const list = useCrudList({
    list: adminListLessonProgressOptions,
    listKey: adminListLessonProgressQueryKey,
  });

  const columns: CrudColumnDef<LessonProgress>[] = [
    {
      accessorKey: "userId",
      header: t(($) => $.models.attributes.lesson_progress.user_id),
    },
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.lesson_progress.course),
    },
    {
      accessorKey: "courseLessonName",
      header: t(($) => $.models.attributes.lesson_progress.lesson),
    },
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.lesson_progress.state),
    },
    {
      accessorKey: "messagesCount",
      header: t(($) => $.models.attributes.lesson_progress.messages_count),
    },
    {
      cell: ({ row }) => formatDateTime(row.original.createdAt),
      header: t(($) => $.models.attributes.base.created_at),
      id: "createdAt",
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.lessonProgress)}</Title>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

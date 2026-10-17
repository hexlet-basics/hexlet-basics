import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListLessonProgressOptions } from "@/client/@tanstack/react-query.gen";
import type { LessonProgress } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";

// Per-lesson progress (route keeps the legacy table name) — read-only.
export const Route = createFileRoute("/{-$locale}/admin/lesson_progress/")({
  component: LessonProgressList,
});

function LessonProgressList() {
  const { t } = useTranslation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListLessonProgressOptions({
      query: pager.query,
    }),
  );

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
      id: "createdAt",
      header: t(($) => $.models.attributes.base.created_at),
      cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.lessonProgress)}</Title>

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

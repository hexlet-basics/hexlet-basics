import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { adminListLessonProgressOptions } from "@/client/@tanstack/react-query.gen";
import type { LessonProgress } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";

// Per-lesson progress (route keeps the legacy table name) — read-only.
export const Route = createFileRoute("/{-$locale}/admin/language_lesson_members/")({
  component: LessonProgressList,
});

function LessonProgressList() {
  const { t } = useTranslation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListLessonProgressOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const columns: CrudColumnDef<LessonProgress>[] = [
    {
      accessorKey: "userId",
      header: t(($) => $.models.attributes.lesson_progress.user_id),
      enableSorting: false,
    },
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.lesson_progress.course),
      enableSorting: false,
    },
    {
      accessorKey: "courseLessonName",
      header: t(($) => $.models.attributes.lesson_progress.lesson),
      enableSorting: false,
    },
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.lesson_progress.state),
      enableSorting: false,
    },
    {
      accessorKey: "messagesCount",
      header: t(($) => $.models.attributes.lesson_progress.messages_count),
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
      <Title order={2}>{t(($) => $.admin.resources.lessonProgress)}</Title>

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

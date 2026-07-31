import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { ColumnDef, PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { adminListCourseLessonMembersOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseLessonMember } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";

// Per-lesson memberships (legacy `/admin/language_lesson_members`) — read-only.
export const Route = createFileRoute("/{-$locale}/admin/language_lesson_members/")({
  component: LessonMembersList,
});

function LessonMembersList() {
  const { t } = useTranslation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListCourseLessonMembersOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const columns: ColumnDef<CourseLessonMember>[] = [
    {
      accessorKey: "userId",
      header: t(($) => $.models.attributes.lesson_member.user_id),
      enableSorting: false,
    },
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.lesson_member.course),
      enableSorting: false,
    },
    {
      accessorKey: "courseLessonName",
      header: t(($) => $.models.attributes.lesson_member.lesson),
      enableSorting: false,
    },
    {
      accessorKey: "state",
      header: t(($) => $.models.attributes.lesson_member.state),
      enableSorting: false,
    },
    {
      accessorKey: "messagesCount",
      header: t(($) => $.models.attributes.lesson_member.messages_count),
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
      <Title order={2}>{t(($) => $.admin.resources.lessonMembers)}</Title>

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

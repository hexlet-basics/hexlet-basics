import { Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { adminListCourseLessonReviewsOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseLessonReview } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";

// AI lesson-review summaries (legacy `/admin/language_lesson_reviews`) —
// read-only; new summaries are produced by the course/lesson review actions.
export const Route = createFileRoute("/{-$locale}/admin/language_lesson_reviews/")({
  component: LessonReviewsList,
});

function LessonReviewsList() {
  const { t } = useTranslation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListCourseLessonReviewsOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const columns: CrudColumnDef<CourseLessonReview>[] = [
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.lesson_review.course),
      enableSorting: false,
    },
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.lesson_review.lesson),
      enableSorting: false,
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.lesson_review.locale),
      enableSorting: false,
    },
    {
      id: "summary",
      header: t(($) => $.models.attributes.lesson_review.summary),
      cell: ({ row }) => (
        <Text size="sm" lineClamp={3} maw={560}>
          {row.original.summary}
        </Text>
      ),
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
      <Title order={2}>{t(($) => $.admin.resources.lessonReviews)}</Title>

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

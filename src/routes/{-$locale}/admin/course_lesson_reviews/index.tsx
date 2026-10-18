import { Stack, Text, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminListCourseLessonReviewsOptions,
  adminListCourseLessonReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListCourseLessonReviewsQuery } from "@/client/zod.gen";
import type { CourseLessonReview } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudList } from "@/components/admin/useCrudList";
import { useDateTime } from "@/hooks/useDateTime";

// AI lesson-review summaries (legacy `/admin/course_lesson_reviews`) —
// read-only; new summaries are produced by the course/lesson review actions.
export const Route = createFileRoute("/{-$locale}/admin/course_lesson_reviews/")({
  validateSearch: zAdminListCourseLessonReviewsQuery,
  component: LessonReviewsList,
});

function LessonReviewsList() {
  const { t } = useTranslation();

  const formatDateTime = useDateTime();
  const list = useCrudList({
    list: adminListCourseLessonReviewsOptions,
    listKey: adminListCourseLessonReviewsQueryKey,
  });

  const columns: CrudColumnDef<CourseLessonReview>[] = [
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.lesson_review.course),
    },
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.lesson_review.lesson),
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.lesson_review.locale),
    },
    {
      id: "summary",
      header: t(($) => $.models.attributes.lesson_review.summary),
      cell: ({ row }) => (
        <Text size="sm" lineClamp={3} maw={560}>
          {row.original.summary}
        </Text>
      ),
    },
    {
      id: "createdAt",
      header: t(($) => $.models.attributes.base.created_at),
      cell: ({ row }) => formatDateTime(row.original.createdAt),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.lessonReviews)}</Title>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

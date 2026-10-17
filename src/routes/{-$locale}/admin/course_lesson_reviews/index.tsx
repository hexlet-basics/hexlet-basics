import { Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListCourseLessonReviewsOptions } from "@/client/@tanstack/react-query.gen";
import type { CourseLessonReview } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";

// AI lesson-review summaries (legacy `/admin/course_lesson_reviews`) —
// read-only; new summaries are produced by the course/lesson review actions.
export const Route = createFileRoute("/{-$locale}/admin/course_lesson_reviews/")({
  component: LessonReviewsList,
});

function LessonReviewsList() {
  const { t } = useTranslation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListCourseLessonReviewsOptions({
      query: pager.query,
    }),
  );

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
      cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.lessonReviews)}</Title>

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

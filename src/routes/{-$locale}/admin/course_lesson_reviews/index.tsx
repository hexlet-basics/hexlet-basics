import { Stack, Text, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListCourseLessonReviewsOptions } from "@/client/@tanstack/react-query.gen";
import { zAdminListCourseLessonReviewsQuery } from "@/client/zod.gen";
import type { CourseLessonReview } from "@/client/types.gen";
import { useDateTime } from "@/hooks/useDateTime";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// AI lesson-review summaries (legacy `/admin/course_lesson_reviews`) —
// read-only; new summaries are produced by the course/lesson review actions.
export const Route = createFileRoute("/{-$locale}/admin/course_lesson_reviews/")({
  component: LessonReviewsList,
  validateSearch: zAdminListCourseLessonReviewsQuery,
});

function LessonReviewsList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListCourseLessonReviewsOptions({
      query: { page: search.page ?? 1, perPage: search.perPage },
    }),
  );

  const formatDateTime = useDateTime();

  const columns: DataTableColumn<CourseLessonReview>[] = [
    {
      accessor: "courseSlug",
      title: t(($) => $.models.attributes.lesson_review.course),
    },
    {
      accessor: "slug",
      title: t(($) => $.models.attributes.lesson_review.lesson),
    },
    {
      accessor: "locale",
      title: t(($) => $.models.attributes.lesson_review.locale),
    },
    {
      accessor: "summary",
      render: (record) => (
        <Text size="sm" lineClamp={3} maw={560}>
          {record.summary}
        </Text>
      ),
      title: t(($) => $.models.attributes.lesson_review.summary),
    },
    {
      accessor: "createdAt",
      render: (record) => formatDateTime(record.createdAt),
      title: t(($) => $.models.attributes.base.created_at),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.lessonReviews)}</Title>

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

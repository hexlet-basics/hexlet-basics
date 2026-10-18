import { Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { adminListLessonProgressOptions } from "@/client/@tanstack/react-query.gen";
import { zAdminListLessonProgressQuery } from "@/client/zod.gen";
import type { LessonProgress } from "@/client/types.gen";
import { useDateTime } from "@/hooks/useDateTime";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// Per-lesson progress (route keeps the legacy table name) — read-only.
export const Route = createFileRoute("/{-$locale}/admin/lesson_progress/")({
  component: LessonProgressList,
  validateSearch: zAdminListLessonProgressQuery,
});

function LessonProgressList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListLessonProgressOptions({ query: { page: search.page ?? 1, perPage: search.perPage } }),
  );

  const formatDateTime = useDateTime();

  const columns: DataTableColumn<LessonProgress>[] = [
    {
      accessor: "userId",
      title: t(($) => $.models.attributes.lesson_progress.user_id),
    },
    {
      accessor: "courseSlug",
      title: t(($) => $.models.attributes.lesson_progress.course),
    },
    {
      accessor: "courseLessonName",
      title: t(($) => $.models.attributes.lesson_progress.lesson),
    },
    {
      accessor: "state",
      title: t(($) => $.models.attributes.lesson_progress.state),
    },
    {
      accessor: "messagesCount",
      title: t(($) => $.models.attributes.lesson_progress.messages_count),
    },
    {
      accessor: "createdAt",
      render: (record) => formatDateTime(record.createdAt),
      title: t(($) => $.models.attributes.base.created_at),
    },
  ];

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.resources.lessonProgress)}</Title>

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

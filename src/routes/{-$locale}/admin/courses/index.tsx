import { Button, Group, Stack, TextInput, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseVersionMutation,
  adminListCoursesOptions,
  adminListCoursesQueryKey,
  adminReviewCourseMutation,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListCoursesQuery } from "@/client/zod.gen";
import type { Course } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { SEARCH_DEBOUNCE_MS } from "@/lib/duration";
import { useState } from "react";
import { useDebouncedCallback } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// Courses admin (legacy `/admin/courses`; the domain concept is Course).
// No delete, parity with legacy. Each row carries the two lifecycle actions:
// build a new exercise version (river job, auto-promoted on success) and
// enqueue AI review of every current lesson.
export const Route = createFileRoute("/{-$locale}/admin/courses/")({
  component: CoursesList,
  validateSearch: zAdminListCoursesQuery,
});

function CoursesList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListCoursesOptions({
      query: { page: search.page ?? 1, perPage: search.perPage, q: search.q },
    }),
  );
  const [term, setTerm] = useState(search.q ?? "");
  const applyTerm = useDebouncedCallback((text: string) => {
    void navigate({
      search: (prev) => ({ ...prev, page: undefined, q: text === "" ? undefined : text }),
    });
  }, SEARCH_DEBOUNCE_MS);

  const versionMutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: [adminListCoursesQueryKey()],
    mutation: adminCreateCourseVersionMutation(),
    successMessage: t(($) => $.admin.courses.versionQueued),
  });

  const reviewMutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: [adminListCoursesQueryKey()],
    mutation: adminReviewCourseMutation(),
    successMessage: t(($) => $.admin.courses.reviewQueued),
  });

  const columns: DataTableColumn<Course>[] = [
    {
      accessor: "slug",
      title: t(($) => $.models.attributes.course.slug),
    },
    {
      accessor: "readiness",
      render: (record) =>
        t(($) => $.models.attributes.course["readiness/values"][record.readiness ?? "draft"]),
      title: t(($) => $.models.attributes.course.readiness),
    },
    {
      accessor: "learnAs",
      title: t(($) => $.models.attributes.course.learn_as),
    },
    {
      accessor: "currentVersion",
      render: (record) => record.currentVersion?.state ?? "",
      title: t(($) => $.admin.courses.currentVersion),
    },
    {
      accessor: "lessonsCount",
      title: t(($) => $.admin.courses.lessonsCount),
    },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions describe={record.slug} id={record.id}>
          <EditLink to="/{-$locale}/admin/courses/$id" params={{ id: String(record.id) }} />

          <Button
            size="xs"
            variant="light"
            loading={versionMutation.isPending && versionMutation.variables?.path.id === record.id}
            onClick={() => {
              versionMutation.mutate({ path: { id: record.id } });
            }}
          >
            {t(($) => $.admin.courses.createVersion)}
          </Button>
          <Button
            size="xs"
            variant="light"
            color="grape"
            loading={reviewMutation.isPending && reviewMutation.variables?.path.id === record.id}
            onClick={() => {
              reviewMutation.mutate({ path: { id: record.id } });
            }}
          >
            {t(($) => $.admin.courses.review)}
          </Button>
        </RowActions>
      ),
      textAlign: "right",
      title: "",
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.courses)}</Title>
        <ButtonLink to="/{-$locale}/admin/courses/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <TextInput
        maw={360}
        placeholder={t(($) => $.admin.crud.search)}
        value={term}
        onChange={(event) => {
          setTerm(event.currentTarget.value);
          applyTerm(event.currentTarget.value.trim());
        }}
      />

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

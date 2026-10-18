import { Button, Group, Stack, Title } from "@mantine/core";
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
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";
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

  const list = useCrudList({
    list: adminListCoursesOptions,
    listKey: adminListCoursesQueryKey,
    searchable: true,
  });

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

  const columns: CrudColumnDef<Course>[] = [
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.course.slug),
    },
    {
      cell: ({ row }) =>
        t(($) => $.models.attributes.course["readiness/values"][row.original.readiness ?? "draft"]),
      header: t(($) => $.models.attributes.course.readiness),
      id: "readiness",
    },
    {
      accessorKey: "learnAs",
      header: t(($) => $.models.attributes.course.learn_as),
    },
    {
      cell: ({ row }) => row.original.currentVersion?.state ?? "",
      header: t(($) => $.admin.courses.currentVersion),
      id: "currentVersion",
    },
    {
      accessorKey: "lessonsCount",
      header: t(($) => $.admin.courses.lessonsCount),
    },
    list.actionsColumn({
      describe: (row) => row.slug,
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/courses/$id" params={{ id: String(row.id) }} />
      ),
      extra: (row) => (
        <>
          <Button
            size="xs"
            variant="light"
            loading={versionMutation.isPending && versionMutation.variables?.path.id === row.id}
            onClick={() => {
              versionMutation.mutate({ path: { id: row.id } });
            }}
          >
            {t(($) => $.admin.courses.createVersion)}
          </Button>
          <Button
            size="xs"
            variant="light"
            color="grape"
            loading={reviewMutation.isPending && reviewMutation.variables?.path.id === row.id}
            onClick={() => {
              reviewMutation.mutate({ path: { id: row.id } });
            }}
          >
            {t(($) => $.admin.courses.review)}
          </Button>
        </>
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.courses)}</Title>
        <ButtonLink to="/{-$locale}/admin/courses/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

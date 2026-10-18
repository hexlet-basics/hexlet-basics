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
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Courses admin (legacy `/admin/courses`; the domain concept is Course).
// No delete, parity with legacy. Each row carries the two lifecycle actions:
// build a new exercise version (river job, auto-promoted on success) and
// enqueue AI review of every current lesson.
export const Route = createFileRoute("/{-$locale}/admin/courses/")({
  validateSearch: zAdminListCoursesQuery,
  component: CoursesList,
});

function CoursesList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListCoursesOptions,
    listKey: adminListCoursesQueryKey,
    searchable: true,
  });

  const versionMutation = useResourceMutation({
    mutation: adminCreateCourseVersionMutation(),
    invalidate: [adminListCoursesQueryKey()],
    successMessage: t(($) => $.admin.courses.versionQueued),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  const reviewMutation = useResourceMutation({
    mutation: adminReviewCourseMutation(),
    invalidate: [adminListCoursesQueryKey()],
    successMessage: t(($) => $.admin.courses.reviewQueued),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  const columns: CrudColumnDef<Course>[] = [
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.course.slug),
    },
    {
      id: "readiness",
      header: t(($) => $.models.attributes.course.readiness),
      cell: ({ row }) =>
        t(($) => $.models.attributes.course["readiness/values"][row.original.readiness ?? "draft"]),
    },
    {
      accessorKey: "learnAs",
      header: t(($) => $.models.attributes.course.learn_as),
    },
    {
      id: "currentVersion",
      header: t(($) => $.admin.courses.currentVersion),
      cell: ({ row }) => row.original.currentVersion?.state ?? "",
    },
    {
      accessorKey: "lessonsCount",
      header: t(($) => $.admin.courses.lessonsCount),
    },
    list.actionsColumn({
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/courses/$id" params={{ id: String(row.id) }} />
      ),
      extra: (row) => (
        <>
          <Button
            size="xs"
            variant="light"
            loading={versionMutation.isPending && versionMutation.variables?.path.id === row.id}
            onClick={() => versionMutation.mutate({ path: { id: row.id } })}
          >
            {t(($) => $.admin.courses.createVersion)}
          </Button>
          <Button
            size="xs"
            variant="light"
            color="grape"
            loading={reviewMutation.isPending && reviewMutation.variables?.path.id === row.id}
            onClick={() => reviewMutation.mutate({ path: { id: row.id } })}
          >
            {t(($) => $.admin.courses.review)}
          </Button>
        </>
      ),
      describe: (row) => row.slug,
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

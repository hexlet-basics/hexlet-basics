import { Button, Group, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseVersionMutation,
  adminListCoursesOptions,
  adminListCoursesQueryKey,
  adminReviewCourseMutation,
} from "@/client/@tanstack/react-query.gen";
import type { Course } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";
import { ButtonLink } from "@/components/RouterLink";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Courses admin (legacy `/admin/courses`; the domain concept is Course).
// No delete, parity with legacy. Each row carries the two lifecycle actions:
// build a new exercise version (river job, auto-promoted on success) and
// enqueue AI review of every current lesson.
export const Route = createFileRoute("/{-$locale}/admin/courses/")({
  component: CoursesList,
});

function CoursesList() {
  const { t } = useTranslation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListCoursesOptions({
      query: pager.query,
    }),
  );

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
      accessorKey: "progress",
      header: t(($) => $.models.attributes.course.readiness),
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
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/courses/$id"
            params={{ id: String(row.original.id) }}
            size="xs"
            variant="light"
          >
            {t(($) => $.admin.crud.edit)}
          </ButtonLink>
          <Button
            size="xs"
            variant="light"
            loading={
              versionMutation.isPending && versionMutation.variables?.path.id === row.original.id
            }
            onClick={() => versionMutation.mutate({ path: { id: row.original.id } })}
          >
            {t(($) => $.admin.courses.createVersion)}
          </Button>
          <Button
            size="xs"
            variant="light"
            color="grape"
            loading={
              reviewMutation.isPending && reviewMutation.variables?.path.id === row.original.id
            }
            onClick={() => reviewMutation.mutate({ path: { id: row.original.id } })}
          >
            {t(($) => $.admin.courses.review)}
          </Button>
        </Group>
      ),
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.courses)}</Title>
        <ButtonLink to="/{-$locale}/admin/courses/new">{t(($) => $.admin.crud.new)}</ButtonLink>
      </Group>

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

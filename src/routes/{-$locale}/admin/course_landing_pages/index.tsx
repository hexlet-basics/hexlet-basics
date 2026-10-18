import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteCourseLandingPageMutation,
  adminListCourseLandingPagesOptions,
  adminListCourseLandingPagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListCourseLandingPagesQuery } from "@/client/zod.gen";
import type { CourseLandingPage } from "@/client/types.gen";
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/")({
  validateSearch: zAdminListCourseLandingPagesQuery,
  component: CourseLandingPagesList,
});

function CourseLandingPagesList(): ReactElement {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListCourseLandingPagesOptions,
    listKey: adminListCourseLandingPagesQueryKey,
    remove: adminDeleteCourseLandingPageMutation,
    searchable: true,
  });

  const columns: CrudColumnDef<CourseLandingPage>[] = [
    {
      accessorKey: "name",
      header: t(($) => $.models.attributes.course_landing_page.name),
    },
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.course_landing_page.slug),
    },
    {
      accessorKey: "courseSlug",
      header: t(($) => $.models.attributes.course_landing_page.course_id),
    },
    {
      accessorKey: "state",
      cell: ({ row }) => {
        const { state } = row.original;
        return state
          ? t(($) => $.models.attributes.course_landing_page["state/values"][state])
          : "";
      },
      header: t(($) => $.models.attributes.course_landing_page.state),
    },
    {
      cell: ({ row }) => (row.original.main ? "✓" : ""),
      header: t(($) => $.models.attributes.course_landing_page.main),
      id: "main",
    },
    list.actionsColumn({
      describe: (row) => row.name,
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/course_landing_pages/$id" params={{ id: String(row.id) }} />
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.landingPages)}</Title>
        <ButtonLink to="/{-$locale}/admin/course_landing_pages/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

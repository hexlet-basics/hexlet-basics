import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteCourseLandingPageMutation,
  adminListCourseLandingPagesOptions,
  adminListCourseLandingPagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { CourseLandingPage } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/")({
  component: CourseLandingPagesList,
});

function CourseLandingPagesList() {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListCourseLandingPagesOptions,
    listKey: adminListCourseLandingPagesQueryKey,
    remove: adminDeleteCourseLandingPageMutation,
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
      header: t(($) => $.models.attributes.course_landing_page.state),
      cell: ({ row }) => {
        const state = row.original.state;
        return state
          ? t(($) => $.models.attributes.course_landing_page["state/values"][state])
          : "";
      },
    },
    {
      id: "main",
      header: t(($) => $.models.attributes.course_landing_page.main),
      cell: ({ row }) => (row.original.main ? "✓" : ""),
    },
    list.actionsColumn({
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/course_landing_pages/$id" params={{ id: String(row.id) }} />
      ),
      describe: (row) => row.name,
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

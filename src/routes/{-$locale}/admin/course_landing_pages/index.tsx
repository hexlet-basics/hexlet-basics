import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteCourseLandingPageMutation,
  adminListCourseLandingPagesOptions,
  adminListCourseLandingPagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { CourseLandingPage } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { useCrudPager } from "@/components/admin/useCrudPager";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/")({
  component: CourseLandingPagesList,
});

function CourseLandingPagesList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const pager = useCrudPager();

  const { data, isLoading } = useQuery(
    adminListCourseLandingPagesOptions({
      query: pager.query,
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteCourseLandingPageMutation(),
    invalidate: [adminListCourseLandingPagesQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
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
    },
    {
      id: "main",
      header: t(($) => $.models.attributes.course_landing_page.main),
      cell: ({ row }) => (row.original.main ? "✓" : ""),
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/course_landing_pages/$id"
            params={{ id: String(row.original.id) }}
            size="xs"
            variant="light"
          >
            {t(($) => $.admin.crud.edit)}
          </ButtonLink>
          <Tooltip label={t(($) => $.admin.crud.delete)}>
            <ActionIcon
              color="red"
              variant="light"
              aria-label={t(($) => $.admin.crud.delete)}
              onClick={() =>
                confirmDelete({
                  description: row.original.name,
                  onConfirm: () => deleteMutation.mutate({ path: { id: row.original.id } }),
                })
              }
            >
              ✕
            </ActionIcon>
          </Tooltip>
        </Group>
      ),
    },
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.landingPages)}</Title>
        <ButtonLink to="/{-$locale}/admin/course_landing_pages/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
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

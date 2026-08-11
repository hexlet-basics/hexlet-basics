import { ActionIcon, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { PaginationState, SortingState } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminDeleteCourseCategoryMutation,
  adminListCourseCategoriesOptions,
  adminListCourseCategoriesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import type { CourseCategory } from "@/client/types.gen";
import { type CrudColumnDef, CrudList } from "@/components/admin/CrudList";
import { ButtonLink } from "@/components/RouterLink";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// Course categories admin list — the first resource wired through the CRUD engine
// (Wave 1). Adding another resource is this file plus new.tsx / $id.tsx, each a
// thin binding of the generated hooks to the shared CrudList/CrudForm.
export const Route = createFileRoute("/{-$locale}/admin/course_categories/")({
  component: CourseCategoriesList,
});

export function CourseCategoriesList() {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });
  // Sorting state is held for the engine, but this resource's columns keep it off
  // until the generic `listPage` honors sortField/sortOrder server-side.
  const [sorting, setSorting] = useState<SortingState>([]);

  const { data, isLoading } = useQuery(
    adminListCourseCategoriesOptions({
      query: { page: pagination.pageIndex + 1, perPage: pagination.pageSize },
    }),
  );

  const deleteMutation = useResourceMutation({
    mutation: adminDeleteCourseCategoryMutation(),
    invalidate: [adminListCourseCategoriesQueryKey()],
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  const columns: CrudColumnDef<CourseCategory>[] = [
    {
      accessorKey: "name",
      header: t(($) => $.models.attributes.course_category.name),
      enableSorting: false,
    },
    {
      accessorKey: "slug",
      header: t(($) => $.models.attributes.course_category.slug),
      enableSorting: false,
    },
    {
      accessorKey: "locale",
      header: t(($) => $.models.attributes.course_category.locale),
      enableSorting: false,
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <Group gap="xs" justify="flex-end" wrap="nowrap">
          <ButtonLink
            to="/{-$locale}/admin/course_categories/$id"
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
                  description: row.original.name ?? String(row.original.id),
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
        <Title order={2}>{t(($) => $.admin.resources.courseCategories)}</Title>
        <ButtonLink to="/{-$locale}/admin/course_categories/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <CrudList
        columns={columns}
        data={data?.items ?? []}
        total={data?.total ?? 0}
        pagination={pagination}
        onPaginationChange={setPagination}
        sorting={sorting}
        onSortingChange={setSorting}
        isLoading={isLoading}
      />
    </Stack>
  );
}

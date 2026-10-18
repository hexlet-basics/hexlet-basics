import { Group, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminDeleteCourseCategoryMutation,
  adminListCourseCategoriesOptions,
  adminListCourseCategoriesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zAdminListCourseCategoriesQuery } from "@/client/zod.gen";
import type { CourseCategory } from "@/client/types.gen";
import { EditLink } from "@/components/admin/EditLink";
import { ButtonLink } from "@/components/RouterLink";
import { DataTable } from "mantine-datatable";
import type { DataTableColumn } from "mantine-datatable";
import { RowActions } from "@/components/admin/RowActions";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement } from "react";

// Course categories admin list — the first resource wired through the CRUD engine
// (Wave 1). Adding another resource is this file plus new.tsx / $id.tsx, each a
export const Route = createFileRoute("/{-$locale}/admin/course_categories/")({
  component: CourseCategoriesList,
  validateSearch: zAdminListCourseCategoriesQuery,
});

function CourseCategoriesList(): ReactElement {
  const { t } = useTranslation();
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isFetching } = useQuery(
    adminListCourseCategoriesOptions({
      query: { page: search.page ?? 1, perPage: search.perPage },
    }),
  );

  const columns: DataTableColumn<CourseCategory>[] = [
    { accessor: "name", title: t(($) => $.models.attributes.course_category.name) },
    { accessor: "slug", title: t(($) => $.models.attributes.course_category.slug) },
    { accessor: "locale", title: t(($) => $.models.attributes.course_category.locale) },
    {
      accessor: "actions",
      render: (record) => (
        <RowActions
          describe={record.name ?? String(record.id)}
          id={record.id}
          listKey={adminListCourseCategoriesQueryKey()}
          remove={adminDeleteCourseCategoryMutation()}
        >
          <EditLink
            to="/{-$locale}/admin/course_categories/$id"
            params={{ id: String(record.id) }}
          />
        </RowActions>
      ),
      textAlign: "right",
      title: "",
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

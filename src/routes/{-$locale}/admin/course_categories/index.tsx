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
import { CrudList } from "@/components/admin/CrudList";
import type { CrudColumnDef } from "@/components/admin/CrudList";
import { EditLink } from "@/components/admin/EditLink";
import { courseCategoryColumns } from "@/components/admin/generated/courseCategory";
import { useGeneratedColumns } from "@/components/admin/useGeneratedFields";
import { useCrudList } from "@/components/admin/useCrudList";
import { ButtonLink } from "@/components/RouterLink";
import type { ReactElement } from "react";

// Course categories admin list — the first resource wired through the CRUD engine
// (Wave 1). Adding another resource is this file plus new.tsx / $id.tsx, each a
// thin binding of the generated hooks to the shared CrudList/CrudForm.
export const Route = createFileRoute("/{-$locale}/admin/course_categories/")({
  validateSearch: zAdminListCourseCategoriesQuery,
  component: CourseCategoriesList,
});

function CourseCategoriesList(): ReactElement {
  const { t } = useTranslation();

  const list = useCrudList({
    list: adminListCourseCategoriesOptions,
    listKey: adminListCourseCategoriesQueryKey,
    remove: adminDeleteCourseCategoryMutation,
  });

  const columns: CrudColumnDef<CourseCategory>[] = [
    ...useGeneratedColumns<CourseCategory>(courseCategoryColumns, ["name", "slug", "locale"]),
    list.actionsColumn({
      describe: (row) => row.name ?? String(row.id),
      edit: (row) => (
        <EditLink to="/{-$locale}/admin/course_categories/$id" params={{ id: String(row.id) }} />
      ),
    }),
  ];

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{t(($) => $.admin.resources.courseCategories)}</Title>
        <ButtonLink to="/{-$locale}/admin/course_categories/new">
          {t(($) => $.admin.crud.new)}
        </ButtonLink>
      </Group>

      <CrudList columns={columns} {...list.props} />
    </Stack>
  );
}

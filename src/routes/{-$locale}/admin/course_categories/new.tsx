import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseCategoryMutation,
  adminListCourseCategoriesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { CourseCategoryForm } from "@/components/admin/forms/CourseCategoryForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { CourseCategoryInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/course_categories/new")({
  component: NewCourseCategory,
});

const emptyCourseCategory: CourseCategoryInput = {
  description: null,
  header: "",
  name: "",
  slug: "",
};

function NewCourseCategory(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/course_categories" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListCourseCategoriesQueryKey()],
    mutation: adminCreateCourseCategoryMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CourseCategoryForm
        defaultValues={emptyCourseCategory}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}

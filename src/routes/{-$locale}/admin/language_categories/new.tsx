import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseCategoryMutation,
  adminListCourseCategoriesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zCourseCategoryInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  emptyCourseCategory,
  useCourseCategoryFields,
} from "@/components/admin/resources/courseCategory";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute(
  "/{-$locale}/admin/language_categories/new",
)({
  component: NewCourseCategory,
});

function NewCourseCategory() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseCategoryFields();

  const backToList = () =>
    navigate({ to: "/{-$locale}/admin/language_categories" });

  const mutation = useResourceMutation({
    mutation: adminCreateCourseCategoryMutation(),
    invalidate: [adminListCourseCategoriesQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zCourseCategoryInput}
        defaultValues={emptyCourseCategory}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}

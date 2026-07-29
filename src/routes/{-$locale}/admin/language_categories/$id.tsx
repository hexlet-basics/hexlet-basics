import { Center, Loader, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetCourseCategoryOptions,
  adminGetCourseCategoryQueryKey,
  adminListCourseCategoriesQueryKey,
  adminUpdateCourseCategoryMutation,
} from "@/client/@tanstack/react-query.gen";
import { zCourseCategoryInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  courseCategoryToForm,
  useCourseCategoryFields,
} from "@/components/admin/resources/courseCategory";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/language_categories/$id")({
  component: EditCourseCategory,
});

function EditCourseCategory() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseCategoryFields();
  const { id } = Route.useParams();
  const categoryId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseCategoryOptions({ path: { id: categoryId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/language_categories" });

  const mutation = useResourceMutation({
    mutation: adminUpdateCourseCategoryMutation(),
    invalidate: [
      adminListCourseCategoriesQueryKey(),
      adminGetCourseCategoryQueryKey({ path: { id: categoryId } }),
    ],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <Center py="xl">
          <Loader />
        </Center>
      ) : (
        <CrudForm
          // Remount per row so the form re-seeds its default values from the
          // loaded category instead of keeping a stale initial snapshot.
          key={data.id}
          fields={fields}
          schema={zCourseCategoryInput}
          defaultValues={courseCategoryToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: categoryId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
